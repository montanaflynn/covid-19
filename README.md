# Covid-19 Data ![](https://github.com/montanaflynn/covid-19/workflows/Update%20Data/badge.svg)

Current covid-19 data segmented into countries and states including `confirmed`, `recovered` and `deaths`.

A [GitHub action](https://github.com/montanaflynn/covid-19/blob/master/.github/workflows/main.yml) checks for updates every 15 minutes and updates the [current.json](https://raw.githubusercontent.com/montanaflynn/covid-19/master/data/current.json) file.

The [JSON file](https://montanaflynn.github.io/covid-19/data/current.json) and [website](https://montanaflynn.github.io/covid-19) are both hosted on GitHub so there is no running costs associated with this project.

## Website

If you want to see the current cases on maps worldwide or by country:

https://montanaflynn.github.io/covid-19

## Data JSON

If you just want the current data in `JSON` format it's available here:

https://montanaflynn.github.io/covid-19/data/current.json

## Data Sources

The data comes from a [https://interaktiv.morgenpost.de](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/) [csv file](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/data/Coronavirus.current.v2.csv) which is converted to JSON and translated into english.

Additional Vietnamese province level data comes from a [https://ncov.moh.gov.vn/ban-do-vn](https://ncov.moh.gov.vn/ban-do-vn) [json api](https://maps.vnpost.vn/app/api/democoronas/).

## TODO

- Better location translations
- Add data tables to website
