package main

import "time"

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
