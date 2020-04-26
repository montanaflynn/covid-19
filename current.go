package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func insertCurrentData(ctx context.Context, tx *sql.Tx, primaryRegion, secondaryRegion string, data cases) {
	stmt, err := tx.Prepare(`
	INSERT INTO data(timestamp, primary_region, secondary_region, confirmed, recovered, deaths, active)
	VALUES($1, $2, $3, $4, $5, $6, $7);
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		time.Now().Unix(),
		primaryRegion,
		secondaryRegion,
		data.Confirmed,
		data.Recovered,
		data.Deaths,
		data.Active,
	)
	if err != nil {
		log.Fatalf("%q: %v\n", err, stmt)
	}
}

func (db *database) saveCurrentData(ctx context.Context, data map[string]map[string]cases) {
	for primaryRegion, secondaryRegions := range data {
		for secondaryRegion, caseData := range secondaryRegions {
			tx, err := db.sqlite.Begin()
			if err != nil {
				log.Fatal(err)
			}

			stmt, err := db.sqlite.Prepare(`
			SELECT confirmed, recovered, deaths, active 
			FROM data WHERE primary_region = $1 AND secondary_region = $2
			ORDER by timestamp desc;
			`)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			var (
				confirmed sql.NullInt32
				recovered sql.NullInt32
				deaths    sql.NullInt32
				active    sql.NullInt32
			)
			err = stmt.QueryRow(primaryRegion, secondaryRegion).Scan(&confirmed, &recovered, &deaths, &active)
			if err != nil && err != sql.ErrNoRows {
				log.Fatal(err)
			}

			if confirmed.Valid && recovered.Valid && deaths.Valid && active.Valid {
				confirmedMatch := int(confirmed.Int32) != caseData.Confirmed
				recoveredMatch := int(recovered.Int32) != caseData.Recovered
				deathsMatch := int(deaths.Int32) != caseData.Deaths
				activeMatch := int(active.Int32) != caseData.Active
				if confirmedMatch && recoveredMatch && deathsMatch && activeMatch {
					insertCurrentData(ctx, tx, primaryRegion, secondaryRegion, caseData)
				}
			} else {
				insertCurrentData(ctx, tx, primaryRegion, secondaryRegion, caseData)
			}
			tx.Commit()
		}
	}
}

func getCurrentData(ctx context.Context, db database) error {

	res, err := http.Get(fmt.Sprintf("%s?t=%d", currentDataURL, time.Now().Unix()))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	r := csv.NewReader(res.Body)

	var data []datum

	// skip the headers
	_, err = r.Read()
	if err == io.EOF {
		log.Fatal("missing csv headers")
	}

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		updated, err := strconv.Atoi(record[11])
		if err != nil {
			log.Printf("%s\n%s\n", err, record)
			continue
		}

		date := time.Unix(int64(updated), 0)

		confirmed, err := strconv.Atoi(record[13])
		if err != nil {
			log.Fatal(err)
		}
		recovered, err := strconv.Atoi(record[14])
		if err != nil {
			log.Fatal(err)
		}
		deaths, err := strconv.Atoi(record[15])
		if err != nil {
			log.Fatal(err)
		}

		longitude, _ := strconv.ParseFloat(record[6], 64)
		latitude, _ := strconv.ParseFloat(record[7], 64)
		population, _ := strconv.Atoi(record[8])

		parent := record[5]
		if parent == "null" {
			parent = "global"
		}

		label := record[4]
		if parent == "Denmark" && label == "Greenland" {
			parent = "global"
		}

		if parent == "Germany" {
			label = record[2]
		}

		if parent == "Canada" && label == "Recovered" {
			continue
		}

		d := datum{
			parent:     parent,
			label:      label,
			updated:    updated,
			date:       date,
			confirmed:  confirmed,
			recovered:  recovered,
			deaths:     deaths,
			active:     confirmed - (recovered + deaths),
			longitude:  longitude,
			latitude:   latitude,
			population: population,
			source:     record[15],
			sourceURL:  record[16],
			scraper:    record[17],
		}

		data = append(data, d)
	}

	countryCounts := make(map[string]cases)
	usaCounts := make(map[string]cases)
	chinaCounts := make(map[string]cases)
	germanyCounts := make(map[string]cases)
	canadaCounts := make(map[string]cases)

	for _, d := range data {

		labelName := d.label
		labelName = strings.Replace(labelName, "Peking (Beijing)", "Beijing", -1)
		labelName = strings.Replace(labelName, "Innere Mongolei", "Nei Mongol", -1)
		labelName = strings.Replace(labelName, "Hubei (Wuhan)", "Hubei", -1)
		labelName = strings.Replace(labelName, "Xinjiang", "Xinjiang Uygur", -1)

		countryName := d.parent
		if d.parent == "global" {
			countryName = d.label
		}

		if d.parent == "global" {
			countryName = strings.Replace(countryName, "USA", "United States", -1)
			countryName = strings.Replace(countryName, "Austia", "Austria", -1)

			updated := d.updated
			confirmed := d.confirmed
			recovered := d.recovered
			deaths := d.deaths
			active := d.active
			population := d.population
			longitude := d.longitude
			latitude := d.latitude

			c, ok := countryCounts[countryName]
			if ok {
				confirmed = c.Confirmed + d.confirmed
				recovered = c.Recovered + d.recovered
				deaths = c.Deaths + d.deaths
				active = c.Active + d.active
				if c.Updated > d.updated {
					updated = c.Updated
				}
			}
			countryCounts[countryName] = cases{
				Updated:    updated,
				Confirmed:  confirmed,
				Recovered:  recovered,
				Deaths:     deaths,
				Active:     active,
				Population: population,
				Longitude:  longitude,
				Latitude:   latitude,
			}

		} else if countryName == "USA" || countryName == "Canada" || countryName == "Germany" || countryName == "China" {
			countMap := map[string]cases{}
			switch countryName {
			case "USA":
				countMap = usaCounts
			case "Canada":
				countMap = canadaCounts
			case "Germany":
				countMap = germanyCounts
			case "China":
				countMap = chinaCounts
			}

			updated := d.updated
			confirmed := d.confirmed
			recovered := d.recovered
			deaths := d.deaths
			active := d.active
			population := d.population
			longitude := d.longitude
			latitude := d.latitude

			c, ok := countMap[labelName]
			if ok {
				confirmed = c.Confirmed + d.confirmed
				recovered = c.Recovered + d.recovered
				deaths = c.Deaths + d.deaths
				active = c.Active + d.active
				if c.Updated > d.updated {
					updated = c.Updated
				}
			}

			countMap[labelName] = cases{
				Updated:    updated,
				Confirmed:  confirmed,
				Recovered:  recovered,
				Deaths:     deaths,
				Active:     active,
				Population: population,
				Longitude:  longitude,
				Latitude:   latitude,
			}
		}
	}

	// get vietnam province data
	vietnamCounts, err := getVietnamData()
	if err != nil {
		log.Fatal(err)
	}

	results := map[string]map[string]cases{
		"global":  countryCounts,
		"usa":     usaCounts,
		"canada":  canadaCounts,
		"germany": germanyCounts,
		"china":   chinaCounts,
		"vietnam": vietnamCounts,
	}

	// save data to database
	db.saveCurrentData(ctx, results)

	// output data in JSON format
	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./data/current.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
