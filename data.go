package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

var (
	httpClient             = http.Client{}
	baseURL                = "https://funkeinteraktiv.b-cdn.net"
	currentDataEndpoint    = "/current.v4.csv"
	historicalDataEndpoint = "/history.light.v4.csv"
	currentDataURL         = fmt.Sprintf("%s%s", baseURL, currentDataEndpoint)
	historicalDataURL      = fmt.Sprintf("%s%s", baseURL, historicalDataEndpoint)
)

func saveData(ctx context.Context, url, file string) error {
	errChan := make(chan error)

	go func() {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errChan <- err
			return
		}

		req.Header.Add("user-agent", "Mozilla/5.0")

		res, err := httpClient.Do(req)
		if err != nil {
			errChan <- err
			return
		}

		defer res.Body.Close()

		f, err := os.Create(file)
		if err != nil {
			errChan <- err
			return
		}
		defer f.Close()

		_, err = io.Copy(f, res.Body)
		if err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			if err != nil {
				return fmt.Errorf("saveData: %w", err)
			}
			return nil
		}
	}
}

func saveOriginalData(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	fn := func(ctx context.Context, url, file string) func() error {
		return func() error {
			return saveData(ctx, url, file)
		}
	}

	g.Go(fn(ctx, currentDataURL, "./data/current.csv"))
	g.Go(fn(ctx, historicalDataURL, "./data/historical.csv"))

	return g.Wait()
}
