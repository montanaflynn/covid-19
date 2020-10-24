# Covid-19 [![](https://github.com/montanaflynn/covid-19/workflows/Update/badge.svg)](https://github.com/montanaflynn/covid-19/actions)

Current and historical covid-19 `confirmed`, `recovered`, `deaths` and `active` case counts by region. 

If you want to see the current cases on maps worldwide or by country:

https://montanaflynn.github.io/covid-19

If you want to see the historical cases by country or region over time:

https://montanaflynn.github.io/covid-19/historical.html

If you want to use the data for your own project:

Dataset | Format | URL
--------|---------|----
Current | JSON | https://montanaflynn.github.io/covid-19/data/current.json
Current | CSV  | https://montanaflynn.github.io/covid-19/data/current.csv
Historical | JSON | https://montanaflynn.github.io/covid-19/data/historical.json
Historical | CSV  | https://montanaflynn.github.io/covid-19/data/current.csv

The data is also available in a [sqlite database](./data/covid.db) with the following tables:

Table | Schema URL
------|--------
`current_data` | [database.go#L14-L24](https://github.com/montanaflynn/covid-19/blob/master/database.go#L12-L22)
`historical_data` | [database.go#L28-L41](https://github.com/montanaflynn/covid-19/blob/master/database.go#L28-L41)

## Architecture

The current and historical data comes from [https://interaktiv.morgenpost.de](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/) [current](https://interaktiv.morgenpost.de/data/corona/current.v4.csv) and [historical](https://interaktiv.morgenpost.de/data/corona/history.light.v4.csv) CSV files which are converted to JSON and commited to this repo along with the original format. The data is also saved in [sqlite database](./data/covid.db).

Additional Vietnamese province level data comes from [wikipedia](https://vi.wikipedia.org/wiki/%C4%90%E1%BA%A1i_d%E1%BB%8Bch_COVID-19_t%E1%BA%A1i_Vi%E1%BB%87t_Nam).

A [GitHub action](https://github.com/montanaflynn/covid-19/blob/master/.github/workflows/main.yml) checks for updates every 15 minutes and updates the [JSON files](./data) and [sqlite database](./data/covid.db).

The JSON files [current.json](https://raw.githubusercontent.com/montanaflynn/covid-19/master/data/current.json) and [historical.json](https://montanaflynn.github.io/covid-19/data/historical.json) are hosted on GitHub so there is no running costs associated with this project.

The website maps and tables are rendered in the browser using [map.js](https://github.com/montanaflynn/covid-19/blob/master/assets/map.js).

## Screenshots

[![](https://i.imgur.com/z370DBE.png)](https://montanaflynn.github.io/covid-19/)
[![](https://i.imgur.com/c4AHfNb.png)](https://montanaflynn.github.io/covid-19/historical.html)

## TODO

- Improve performance of webpage to only load data once for all maps
- Add responsive styling to work for all screen sizes
