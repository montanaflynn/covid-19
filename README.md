# Covid-19 Data

Current covid-19 data segmented into countries and states.

Includes `confirmed`, `recovered` and `deaths`.

The actual data comes from a https://interaktiv.morgenpost.de [csv file](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/data/Coronavirus.current.v2.csv) which is converted to JSON and translated into english.

The current JSON is hosted on GitHub at https://raw.githubusercontent.com/montanaflynn/covid-19/master/data/current.json

The data checks for updates and is commited every 15 minutes using a [GitHub action](https://github.com/montanaflynn/covid-19/blob/master/.github/workflows/main.yml).

### Example Usage

```
go run *.go
{
  "countries": {
    "Afghanistan": {
      "confirmed": 21,
      "recovered": 1,
      "deaths": 0
    },...
  },
  "states": {
    "Alabama": {
      "confirmed": 29,
      "recovered": 0,
      "deaths": 0
    },...
  }
}
```
