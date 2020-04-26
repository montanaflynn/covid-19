package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	currentDataBaseURL  = "https://funkeinteraktiv.b-cdn.net"
	currentDataEndpoint = "/current.v4.csv"
	currentDataURL      = fmt.Sprintf("%s%s", currentDataBaseURL, currentDataEndpoint)
)

type datum struct {
	parent     string
	label      string
	updated    int
	date       time.Time
	confirmed  int
	recovered  int
	deaths     int
	active     int
	population int
	latitude   float64
	longitude  float64
	source     string
	sourceURL  string
	scraper    string
}

type cases struct {
	Date       int     `json:"date,omitempty"`
	Updated    int     `json:"updated"`
	Confirmed  int     `json:"confirmed"`
	Recovered  int     `json:"recovered"`
	Deaths     int     `json:"deaths"`
	Active     int     `json:"active"`
	Population int     `json:"population,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
}

func main() {

	log.SetFlags(log.Llongfile)
	db := newDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		return db.createCurrentDataTable(ctx)
	})
	g.Go(func() error {
		return db.createHistoricalDataTable(ctx)
	})

	g, _ = errgroup.WithContext(ctx)

	g.Go(saveOriginalData)

	g.Go(func() error {
		return getHistoricalData(ctx, db)
	})
	g.Go(func() error {
		return getCurrentData(ctx, db)
	})

	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}

	err = ctx.Err()
	if err != nil {
		log.Fatal(err)
	}
}
