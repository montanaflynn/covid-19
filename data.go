package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

func (d *database) createTable() {
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
	_, err := d.sqlite.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (d *database) saveData(data map[string]map[string]cases) {
	for primaryRegion, secondaryRegions := range data {
		for secondaryRegion, caseData := range secondaryRegions {
			tx, err := d.sqlite.Begin()
			if err != nil {
				log.Fatal(err)
			}

			stmt, err := d.sqlite.Prepare(`
			SELECT confirmed, recovered, deaths, active 
			FROM data WHERE primary_region = $1 AND secondary_region = $2
			ORDER by timestamp desc;
			`)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			var (
				confirmed sql.NullInt32
				recovered sql.NullInt32
				deaths    sql.NullInt32
				active    sql.NullInt32
			)
			err = stmt.QueryRow(primaryRegion, secondaryRegion).Scan(&confirmed, &recovered, &deaths, &active)
			if err != nil {
				log.Fatal(err)
			}

			if confirmed.Valid && recovered.Valid && deaths.Valid && active.Valid {
				confirmedMatch := int(confirmed.Int32) != caseData.Confirmed
				recoveredMatch := int(recovered.Int32) != caseData.Recovered
				deathsMatch := int(deaths.Int32) != caseData.Deaths
				activeMatch := int(active.Int32) != caseData.Active
				if confirmedMatch && recoveredMatch && deathsMatch && activeMatch {
					stmt, err = tx.Prepare(`
					INSERT INTO data(timestamp, primary_region, secondary_region, confirmed, recovered, deaths, active)
					VALUES($1, $2, $3, $4, $5, $6, $7);
					`)
					if err != nil {
						log.Fatal(err)
					}
					defer stmt.Close()

					_, err = stmt.Exec(
						time.Now().Unix(),
						primaryRegion,
						secondaryRegion,
						caseData.Confirmed,
						caseData.Recovered,
						caseData.Deaths,
						caseData.Active,
					)
					if err != nil {
						log.Fatalf("%q: %v\n", err, stmt)
					}
				}
			}
			tx.Commit()
		}
	}
}
