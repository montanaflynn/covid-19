package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	httpClient             = http.Client{}
	baseURL                = "https://interaktiv.morgenpost.de"
	currentDataEndpoint    = "/data/corona/current.v4.csv"
	historicalDataEndpoint = "/data/corona/history.light.v4.csv"
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
