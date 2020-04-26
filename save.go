package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

func saveOriginalCurrentData() error {
	res, err := http.Get(fmt.Sprintf("%s?t=%d", currentDataURL, time.Now().Unix()))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	f, err := os.Create("./data/current.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func saveOriginalHistoricalData() error {
	res, err := http.Get(historicalDataURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	f, err := os.Create("./data/historical.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func saveOriginalData() error {
	ctx, done := context.WithCancel(context.Background())
	g, _ := errgroup.WithContext(ctx)

	g.Go(saveOriginalCurrentData)
	g.Go(saveOriginalHistoricalData)

	time.AfterFunc(60*time.Second, func() {
		fmt.Printf("force finished after 10s")
		done()
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}
