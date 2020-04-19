# Covid-19 [![](https://github.com/montanaflynn/covid-19/workflows/Update/badge.svg)](https://github.com/montanaflynn/covid-19/actions)

Current covid-19 `confirmed`, `recovered`, `deaths` and `active` case counts by region.

## Website

If you want to see the current cases on maps worldwide or by country:

https://montanaflynn.github.io/covid-19

## Data JSON

If you want the current data in `JSON` format it's available here:

https://montanaflynn.github.io/covid-19/data/current.json

## Data Sources

The data comes from a [https://interaktiv.morgenpost.de](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/) [csv file](https://interaktiv.morgenpost.de/corona-virus-karte-infektionen-deutschland-weltweit/data/Coronavirus.current.v2.csv) which is converted to JSON and translated into english.

Additional Vietnamese province level data comes from [wikipedia](https://vi.wikipedia.org/wiki/%C4%90%E1%BA%A1i_d%E1%BB%8Bch_COVID-19_t%E1%BA%A1i_Vi%E1%BB%87t_Nam).

## Architecture

A [GitHub action](https://github.com/montanaflynn/covid-19/blob/master/.github/workflows/main.yml) checks for updates every 15 minutes and updates the [current.json](https://raw.githubusercontent.com/montanaflynn/covid-19/master/data/current.json) file.

The [JSON file](https://montanaflynn.github.io/covid-19/data/current.json) and [website](https://montanaflynn.github.io/covid-19) are both hosted on GitHub so there is no running costs associated with this project.

The website maps and tables are rendered in the browser using [map.js](https://github.com/montanaflynn/covid-19/blob/master/assets/map.js).

## Screenshot

[![](https://i.imgur.com/z370DBE.png)](https://montanaflynn.github.io/covid-19/)

## TODO

- Better location translations
- Improve performance of webpage to only load data once for all maps
- Add responsive styling to work for all screen sizes
