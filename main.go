package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	baseURL  = "https://funkeinteraktiv.b-cdn.net"
	endpoint = "/current.v4.csv"
)

type datum struct {
	parent    string
	label     string
	updated   int
	date      time.Time
	confirmed int
	recovered int
	deaths    int
	source    string
	sourceURL string
	scraper   string
}

type cases struct {
	Updated   int `json:"updated"`
	Confirmed int `json:"confirmed"`
	Recovered int `json:"recovered"`
	Deaths    int `json:"deaths"`
	Active    int `json:"active"`
}

type results struct {
	Global  map[string]cases `json:"global"`
	USA     map[string]cases `json:"usa"`
	Canada  map[string]cases `json:"canada"`
	Germany map[string]cases `json:"germany"`
	China   map[string]cases `json:"china"`
	Vietnam map[string]cases `json:"vietnam"`
}

func main() {

	log.SetFlags(log.Llongfile)

	res, err := http.Get(fmt.Sprintf("%s%s?t=%d", baseURL, endpoint, time.Now().Unix()))
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
			parent:    parent,
			label:     label,
			updated:   updated,
			date:      date,
			confirmed: confirmed,
			recovered: recovered,
			deaths:    deaths,
			source:    record[15],
			sourceURL: record[16],
			scraper:   record[17],
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

			c, ok := countryCounts[countryName]
			if !ok {
				active := d.confirmed - d.recovered - d.deaths
				countryCounts[countryName] = cases{d.updated, d.confirmed, d.recovered, d.deaths, active}
			} else {
				updatedConfirmed := c.Confirmed + d.confirmed
				updatedRecovered := c.Recovered + d.recovered
				updatedDeaths := c.Deaths + d.deaths
				updatedActive := c.Active + (d.confirmed - d.recovered - d.deaths)
				updated := d.updated
				if c.Updated > d.updated {
					updated = c.Updated
				}
				countryCounts[countryName] = cases{updated, updatedConfirmed, updatedRecovered, updatedDeaths, updatedActive}
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

			c, ok := countMap[labelName]
			if !ok {
				active := d.confirmed - d.recovered - d.deaths
				countMap[labelName] = cases{d.updated, d.confirmed, d.recovered, d.deaths, active}
			} else {
				updatedConfirmed := c.Confirmed + d.confirmed
				updatedRecovered := c.Recovered + d.recovered
				updatedDeaths := c.Deaths + d.deaths
				updatedActive := c.Active + (d.confirmed - d.recovered - d.deaths)
				updated := d.updated
				if c.Updated > d.updated {
					updated = c.Updated
				}
				countMap[labelName] = cases{updated, updatedConfirmed, updatedRecovered, updatedDeaths, updatedActive}
			}
		}
	}

	// get vietnam province data
	vietnamCounts, err := getVietnamData()
	if err != nil {
		log.Fatal(err)
	}

	output := results{
		Global:  countryCounts,
		USA:     usaCounts,
		Canada:  canadaCounts,
		Germany: germanyCounts,
		China:   chinaCounts,
		Vietnam: vietnamCounts,
	}
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", jsonBytes)
	return
}
