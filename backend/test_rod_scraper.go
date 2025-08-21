package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	// Launch browser with rod
	url := launcher.New().
		Headless(true).
		NoSandbox(true).
		Set("disable-gpu").
		Set("disable-dev-shm-usage").
		Set("disable-extensions").
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// Navigate to page
	page := browser.MustPage("https://simplywall.st/stocks/gh/dividend-yield-high")
	defer page.MustClose()

	// Wait for dividend cards to load
	page.Timeout(15 * time.Second).MustElement("div[data-testid='screener-card']")

	// Extract company data
	cards := page.MustElements("div[data-testid='screener-card']")
	
	fmt.Printf("Found %d dividend stocks:\n", len(cards))
	
	for i, card := range cards {
		if i >= 5 { // Limit to first 5 for testing
			break
		}
		
		company, err := card.Element("h2")
		if err != nil {
			continue
		}
		companyName := company.MustText()

		// Extract sector
		sector := ""
		if sectorEl, err := card.ElementR("div", "Sector"); err == nil {
			sector = sectorEl.MustText()
		}

		// Extract dividend yield
		yield := ""
		if yieldEl, err := card.ElementR("div", "Dividend"); err == nil {
			yield = yieldEl.MustText()
		}

		fmt.Printf("%d. %s | %s | %s\n", i+1, companyName, sector, yield)
	}

	log.Println("Rod scraping test completed successfully!")
}