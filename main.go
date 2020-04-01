package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/montanaflynn/gountries"
)

var (
	baseURL    = "https://interaktiv.morgenpost.de"
	endpoint   = "/corona-virus-karte-infektionen-deutschland-weltweit/data/Coronavirus.current.v2.csv"
	vietnamAPI = "https://maps.vnpost.vn/app/api/democoronas/"
)

type datum struct {
	parent    string
	label     string
	updated   int
	date      time.Time
	confirmed int
	recovered int
	deaths    int
	lon       float64
	lat       float64
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

	res, err := http.Get(fmt.Sprintf("%s%s", baseURL, endpoint))
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

		updated, err := strconv.Atoi(record[2])
		if err != nil {
			log.Printf("%s\n%s\n", err, record)
			continue
		}

		date := time.Unix(int64(updated), 0)

		confirmed, err := strconv.Atoi(record[4])
		if err != nil {
			log.Fatal(err)
		}
		recovered, err := strconv.Atoi(record[5])
		if err != nil {
			log.Fatal(err)
		}
		deaths, err := strconv.Atoi(record[6])
		if err != nil {
			log.Fatal(err)
		}

		d := datum{
			parent:    record[0],
			label:     record[1],
			updated:   updated,
			date:      date,
			confirmed: confirmed,
			recovered: recovered,
			deaths:    deaths,
			source:    record[9],
			sourceURL: record[10],
			scraper:   record[11],
		}

		if record[7] != "null" {
			lon, err := strconv.ParseFloat(record[7], 64)
			if err != nil {
				log.Fatal(err)
			}
			d.lon = lon
		}

		if record[8] != "null" {
			lat, err := strconv.ParseFloat(record[8], 64)
			if err != nil {
				log.Fatal(err)
			}
			d.lat = lat
		}

		// fmt.Printf("%+v", d)

		data = append(data, d)
	}

	countryCounts := make(map[string]cases)
	usaCounts := make(map[string]cases)
	chinaCounts := make(map[string]cases)
	germanyCounts := make(map[string]cases)
	canadaCounts := make(map[string]cases)

	query := gountries.New()

	for _, d := range data {

		labelName := d.label
		labelName = strings.Replace(labelName, "Peking (Beijing)", "Beijing", -1)
		labelName = strings.Replace(labelName, "Innere Mongolei", "Nei Mongol", -1)
		labelName = strings.Replace(labelName, "Hubei (Wuhan)", "Hubei", -1)
		labelName = strings.Replace(labelName, "Xinjiang", "Xinjiang Uygur", -1)
		labelName = strings.Replace(labelName, "Amerikanische Jungferninseln", "U.S. Virgin Islands", -1)
		labelName = strings.Replace(labelName, "Nördliche Marianen", "Northern Mariana Islands", -1)
		labelName = strings.Replace(labelName, "Amerikanisch-Samoa", "American Samoa", -1)
		labelName = strings.Replace(labelName, "Neufundland und Labrador", "Newfoundland and Labrador", -1)
		labelName = strings.Replace(labelName, "Kreuzfahrtschiffe", "Cruiseship", -1)
		labelName = strings.Replace(labelName, "Kreuzfahrtschiff", "Cruiseship", -1)

		countryName := d.parent
		countryName = strings.Replace(countryName, "Kreuzfahrtschiffe und Sonstige", "Cruiseship and others", -1)

		country, err := query.FindCountryByAlpha(countries[d.parent])
		if err == nil {
			countryName = country.Name.Common
		}

		if d.parent == "global" {
			countryName = labelName

			country, err := query.FindCountryByAlpha(countries[labelName])
			if err == nil {
				countryName = country.Name.Common
			}
		}

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

		if countryName == "United States" || countryName == "Canada" || countryName == "Germany" || countryName == "China" {

			countMap := map[string]cases{}
			switch countryName {
			case "United States":
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
	res, err = http.Get(vietnamAPI)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	vietnamData := []vietnamDataSchema{}

	err = json.Unmarshal(resBody, &vietnamData)
	if err != nil {
		log.Fatal(err)
	}

	vietnamCounts := make(map[string]cases)

	for _, d := range vietnamData {
		c, ok := vietnamCounts[d.Address]
		if !ok {
			active := d.Number - d.Recovered
			vietnamCounts[d.Address] = cases{0, d.Number, d.Recovered, 0, active}
		} else {
			updatedConfirmed := c.Confirmed + d.Number
			updatedRecovered := c.Recovered + d.Recovered
			updatedDeaths := c.Deaths + 0
			updatedActive := c.Active + (d.Number - d.Recovered)
			updated := 0
			if c.Updated > 0 {
				updated = c.Updated
			}
			vietnamCounts[d.Address] = cases{updated, updatedConfirmed, updatedRecovered, updatedDeaths, updatedActive}
		}
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

type vietnamDataSchema struct {
	Recovered int     `json:"recovered"`
	Doubt     int     `json:"doubt"`
	Code      string  `json:"code"`
	Lat       float64 `json:"Lat"`
	Number    int     `json:"number"`
	Date      string  `json:"date"`
	Lng       float64 `json:"Lng"`
	ID        int     `json:"id"`
	Address   string  `json:"address"`
	Deaths    int     `json:"deaths"`
}
