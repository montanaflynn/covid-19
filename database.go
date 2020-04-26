package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type database struct {
	sqlite *sql.DB
}

func newDatabase() database {
	db, err := sql.Open("sqlite3", "./covid.db")
	if err != nil {
		log.Fatal(err)
	}
	return database{db}
}

func (db *database) createCurrentDataTable(ctx context.Context) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS data (
		timestamp INTEGER NOT NULL,
		primary_region STRING NOT NULL,
		secondary_region STRING,
		confirmed INTEGER,
		recovered INTEGER,
		deaths INTEGER,
		active INTEGER
	);`
	_, err := db.sqlite.ExecContext(ctx, sqlStmt)
	if err != nil {
		return errors.Wrap(err, sqlStmt)
	}
	return nil
}

func (db *database) createHistoricalDataTable(ctx context.Context) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS historical_data (
		timestamp INTEGER NOT NULL,
		date INTEGER NOT NULL,
		primary_region STRING NOT NULL,
		secondary_region STRING,
		confirmed INTEGER,
		recovered INTEGER,
		deaths INTEGER,
		active INTEGER,
		population INTEGER,
		longitude REAL,
		latitude REAL,
		UNIQUE(date, primary_region, secondary_region)
	);`
	_, err := db.sqlite.ExecContext(ctx, sqlStmt)
	if err != nil {
		return errors.Wrap(err, sqlStmt)
	}
	return nil
}
