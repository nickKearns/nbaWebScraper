package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// Player struct
type Player struct {
	Name         string `json:"player"`
	InjuryReport string `json:"injuryReport"`
}

var injuredPlayers []Player

func writeToJSON(players []Player, fileName string) {

	// Alex Gray helped me with this function
	jsonData, err := json.MarshalIndent(players, "", "    ")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	f.Write(jsonData)

	f.Close()

}

func visitSite() {
	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("a[href$='roster/']", func(e *colly.HTMLElement) {
		// Find link using an attribute selector
		// Matches any element that includes href=""
		link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		// Visit link
		e.Request.Visit(link)
	})

	c.OnHTML("span.CellPlayerName--long > span ", func(e *colly.HTMLElement) {

		playerName := e.ChildText("a")

		injuryReport := e.ChildText("span.CellPlayerName-icon.icon-moon-injury > div > div")

		if injuryReport != "" {

			injuredPlayers = append(injuredPlayers, Player{playerName, injuryReport})

			playerName = strings.TrimSpace(playerName)
			println(playerName)
			println(injuryReport)
		}

	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		// fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		// fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://www.cbssports.com/nba/teams/")
}

func main() {

	visitSite()

	writeToJSON(injuredPlayers, "injuredPlayers.json")

}
