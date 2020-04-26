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

func saveConvertedData(ctx context.Context, db *database) error {
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

	db, err := newDatabase()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = saveOriginalData(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = saveConvertedData(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
}
