// find_links_in_page.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	wikiBaseURL  = "https://vi.wikipedia.org"
	wikiEndpoint = "/wiki/%C4%90%E1%BA%A1i_d%E1%BB%8Bch_COVID-19_t%E1%BA%A1i_Vi%E1%BB%87t_Nam"
)

func getVietnamData() (map[string]cases, error) {
	t := newTimer("getting vietnamese wikipedia data")
	defer t.end()
	response, err := http.Get(fmt.Sprintf("%s%s", wikiBaseURL, wikiEndpoint))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	var rows [][]string

	vietnamCounts := make(map[string]cases)

	now := time.Now().Unix() * 1000

	doc.Find("table").Each(func(index int, table *goquery.Selection) {
		table.Find("caption span.nowrap").Each(func(indextr int, element *goquery.Selection) {
			if element.Text() == "Số ca nhiễm theo tỉnh thành tại Việt Nam" {
				table.Find("tbody tr").Each(func(index int, rowhtml *goquery.Selection) {
					row := []string{}
					rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
						row = append(row, tablecell.Text())
					})
					rows = append(rows, row)
				})
				for _, r := range rows {
					if len(r) == 5 {
						province := strings.TrimSpace(r[0])
						if province == "TP. Hồ Chí Minh" {
							province = "Hồ Chí Minh"
						}
						confirmed, err := strconv.Atoi(strings.TrimSpace(r[1]))
						if err != nil {
							log.Fatal(err)
						}
						recovered, err := strconv.Atoi(strings.TrimSpace(r[3]))
						if err != nil {
							log.Fatal(err)
						}
						deaths, err := strconv.Atoi(strings.TrimSpace(r[4]))
						if err != nil {
							log.Fatal(err)
						}
						active := confirmed - recovered - deaths
						vietnamCounts[province] = cases{
							Updated:   int(now),
							Confirmed: confirmed,
							Recovered: recovered,
							Deaths:    deaths,
							Active:    active,
						}

					}
				}
			}
		})
	})

	return vietnamCounts, nil
}
