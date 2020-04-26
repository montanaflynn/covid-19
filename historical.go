package main

import (
	"context"
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

var (
	historicalDataBaseURL  = "https://funkeinteraktiv.b-cdn.net"
	historicalDataEndpoint = "/history.light.v4.csv"
	historicalDataURL      = fmt.Sprintf("%s%s", historicalDataBaseURL, historicalDataEndpoint)
)

func (d *database) insertHistoricalData(ctx context.Context, data datum) error {
	sqlStatement := `
	INSERT OR IGNORE INTO historical_data(
		timestamp, 
		date,
		primary_region, 
		secondary_region, 
		confirmed, 
		recovered, 
		deaths, 
		active,
		population,
		longitude,
		latitude
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	res, err := d.sqlite.ExecContext(ctx, sqlStatement,
		data.updated,
		data.date.Unix(),
		data.parent,
		data.label,
		data.confirmed,
		data.recovered,
		data.deaths,
		data.active,
		data.population,
		data.longitude,
		data.latitude,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 1 {
		return fmt.Errorf("rows affected: %d", rowsAffected)
	}

	return nil
}

func getHistoricalData(ctx context.Context, db database) error {
	res, err := http.Get(historicalDataURL)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	r := csv.NewReader(res.Body)

	// skip the headers
	_, err = r.Read()
	if err == io.EOF {
		log.Fatal("missing csv headers")
	}

	var data []datum

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

		date, err := time.Parse("20060102", record[9])
		if err != nil {
			log.Printf("%s\n%s\n", err, record)
			continue
		}

		updated, err := strconv.Atoi(record[11])
		if err != nil {
			log.Printf("%s\n%s\n", err, record)
			continue
		}

		updatedDate := time.Unix(int64(updated), 0)

		confirmed, err := strconv.Atoi(record[12])
		if err != nil {
			log.Fatal(err)
		}

		recovered, err := strconv.Atoi(record[13])
		if err != nil {
			log.Fatal(err)
		}

		deaths, err := strconv.Atoi(record[14])
		if err != nil {
			log.Fatal(err)
		}

		longitude, err := strconv.ParseFloat(record[6], 32)
		if err != nil {
			continue
		}

		latitude, err := strconv.ParseFloat(record[7], 32)
		if err != nil {
			continue
		}

		population, err := strconv.Atoi(record[8])
		if err != nil {
			continue
		}

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
			updated:    int(updatedDate.Unix()),
			date:       date,
			confirmed:  confirmed,
			recovered:  recovered,
			deaths:     deaths,
			longitude:  longitude,
			latitude:   latitude,
			population: population,
		}

		data = append(data, d)
	}

	countryCounts := map[string][]cases{}
	usaCounts := map[string][]cases{}
	canadaCounts := map[string][]cases{}
	germanyCounts := map[string][]cases{}
	chinaCounts := map[string][]cases{}

	results := map[string]map[string][]cases{
		"global":  countryCounts,
		"usa":     usaCounts,
		"canada":  canadaCounts,
		"germany": germanyCounts,
		"china":   chinaCounts,
	}

	for _, d := range data {
		err := db.insertHistoricalData(ctx, d)
		if err != nil {
			log.Fatal(err)
		}

		countryName := d.parent
		if d.parent == "global" {
			countryName = d.label
		}

		if d.parent == "global" {
			countryName = strings.Replace(countryName, "USA", "United States", -1)
			countryName = strings.Replace(countryName, "Austia", "Austria", -1)

			date := int(d.date.Unix())
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
				countryCounts[countryName] = append(c, cases{
					Updated:    updated,
					Date:       date,
					Confirmed:  confirmed,
					Recovered:  recovered,
					Deaths:     deaths,
					Active:     active,
					Population: population,
					Longitude:  longitude,
					Latitude:   latitude,
				})
			} else {
				countryCounts[countryName] = append([]cases{}, cases{
					Updated:    updated,
					Date:       date,
					Confirmed:  confirmed,
					Recovered:  recovered,
					Deaths:     deaths,
					Active:     active,
					Population: population,
					Longitude:  longitude,
					Latitude:   latitude,
				})
			}

		} else if countryName == "USA" || countryName == "Canada" || countryName == "Germany" || countryName == "China" {
			countMap := map[string][]cases{}
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

			date := int(d.date.Unix())
			updated := d.updated
			confirmed := d.confirmed
			recovered := d.recovered
			deaths := d.deaths
			active := d.active
			population := d.population
			longitude := d.longitude
			latitude := d.latitude

			c, ok := countMap[d.label]
			if ok {
				countMap[d.label] = append(c, cases{
					Updated:    updated,
					Date:       date,
					Confirmed:  confirmed,
					Recovered:  recovered,
					Deaths:     deaths,
					Active:     active,
					Population: population,
					Longitude:  longitude,
					Latitude:   latitude,
				})
			} else {
				countMap[d.label] = append([]cases{}, cases{
					Updated:    updated,
					Date:       date,
					Confirmed:  confirmed,
					Recovered:  recovered,
					Deaths:     deaths,
					Active:     active,
					Population: population,
					Longitude:  longitude,
					Latitude:   latitude,
				})
			}
		}

	}

	// output data in JSON format
	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./data/historical.json")
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
