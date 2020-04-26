# Covid-19 [![](https://github.com/montanaflynn/covid-19/workflows/Update/badge.svg)](https://github.com/montanaflynn/covid-19/actions)

Current and historical covid-19 `confirmed`, `recovered`, `deaths` and `active` case counts by region. 

If you want to see the current cases on maps worldwide or by country:

https://montanaflynn.github.io/covid-19

## Data

### JSON Format:

Data | URL
-----|--------
Current | https://montanaflynn.github.io/covid-19/data/current.json
Historical | https://montanaflynn.github.io/covid-19/data/historical.json


### CSV Format:

Data | URL
-----|--------
Current | https://montanaflynn.github.io/covid-19/data/historical.csv
Historical | https://montanaflynn.github.io/covid-19/data/current.csv


Also available as [sqlite3 database](./data/covid.db) in the following tables:

Table | Schema URL
------|--------
`current_data` | [database.go#L14-L24](https://github.com/montanaflynn/covid-19/blob/master/database.go#L12-L22)
`historical_data` | [database.go#L28-L41](https://github.com/montanaflynn/covid-19/blob/master/database.go#L28-L41)

## Architecture

The current and historical data comes from [https://interaktiv.morgenpost.de](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/) [current](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/data/Coronavirus.current.v2.csv) and [historical](https://funkeinteraktiv.b-cdn.net/history.light.v4.csv) CSV files which are converted to JSON and commited to this repo along with the original format. The data is also saved in [sqlite3 database](./data/covid.db).

Additional Vietnamese province level data comes from [wikipedia](https://vi.wikipedia.org/wiki/%C4%90%E1%BA%A1i_d%E1%BB%8Bch_COVID-19_t%E1%BA%A1i_Vi%E1%BB%87t_Nam).

A [GitHub action](https://github.com/montanaflynn/covid-19/blob/master/.github/workflows/main.yml) checks for updates every 15 minutes and updates the [JSON files](./data) and [sqlite3 database](./data/covid.db).

The JSON files [current.json](https://raw.githubusercontent.com/montanaflynn/covid-19/master/data/current.json) and [historical.json](https://montanaflynn.github.io/covid-19/data/historical.json) are hosted on GitHub so there is no running costs associated with this project.

The website maps and tables are rendered in the browser using [map.js](https://github.com/montanaflynn/covid-19/blob/master/assets/map.js).

## Screenshot

[![](https://i.imgur.com/z370DBE.png)](https://montanaflynn.github.io/covid-19/)

## TODO

- Improve performance of webpage to only load data once for all maps
- Add responsive styling to work for all screen sizes
