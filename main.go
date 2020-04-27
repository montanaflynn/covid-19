package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	historicalJSONFilePath = "./data/historical.json"
	currentJSONFilePath    = "./data/current.json"
	historicalCSVFilePath  = "./data/historical.csv"
	currentCSVFilePath     = "./data/current.csv"
)

func getData(ctx context.Context, db *database) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return getHistoricalData(ctx, db, historicalCSVFilePath)
	})

	g.Go(func() error {
		return getCurrentData(ctx, db, currentCSVFilePath)
	})

	return g.Wait()
}

func main() {
	log.SetFlags(log.Llongfile)

	t := newTimer("setting up sqlite database")

	db, err := newDatabase()
	if err != nil {
		log.Fatal(err)
	}

	t.end().reset("getting all the data")
	defer t.end()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = getData(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
}
