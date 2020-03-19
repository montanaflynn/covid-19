# Covid-19 Data ![](https://github.com/montanaflynn/covid-19/workflows/Update%20Data/badge.svg)

Current covid-19 data segmented into countries and states including `confirmed`, `recovered` and `deaths`.

The data comes from a [https://interaktiv.morgenpost.de](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/) [csv file](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/data/Coronavirus.current.v2.csv) which is converted to JSON and translated into english.

A [GitHub action](https://github.com/montanaflynn/covid-19/blob/master/.github/workflows/main.yml) checks for updates every 15 minutes and updates the [current.json](https://raw.githubusercontent.com/montanaflynn/covid-19/master/data/current.json) file.

### JSON Data

If you just want the current data in JSON format it's available here:

https://montanaflynn.github.io/covid-19/data/current.json

### Website

If you want to see the current cases by country, state or province on a map:

https://montanaflynn.github.io/covid-19

### Example Usage

```
go run *.go
{
  "global": {
    "Afghanistan": {
      "confirmed": 21,
      "recovered": 1,
      "deaths": 0
    },...
  },
  "usa": {
    "Alabama": {
      "confirmed": 29,
      "recovered": 0,
      "deaths": 0
    },...
  }
}
```

### TODO

- Better Country Name Translation
- Toggle between US and World in Website
- Add table with sorting and filtering to Website
