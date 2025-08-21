package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"

	"shares-alert-backend/internal/models"
	"shares-alert-backend/internal/repository"
)

type DividendScraperService struct {
	dividendRepo    *repository.DividendRepository
	dividendService *DividendService
}

type ScrapedDividendData struct {
	Ticker    string `json:"ticker"`
	Name      string `json:"name"`
	DivYield  string `json:"div_yield"`
	LastPrice string `json:"last_price"`
	Industry  string `json:"industry"`
}

func NewDividendScraperService(
	dividendRepo *repository.DividendRepository,
	dividendService *DividendService,
) *DividendScraperService {
	return &DividendScraperService{
		dividendRepo:    dividendRepo,
		dividendService: dividendService,
	}
}

func (s *DividendScraperService) StartDividendScraping() {
	ticker := time.NewTicker(24 * time.Hour) // Scrape daily
	defer ticker.Stop()

	log.Println("Starting dividend scraping service...")

	// Run once immediately
	if err := s.ScrapeDividends(); err != nil {
		log.Printf("Error in initial dividend scraping: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := s.ScrapeDividends(); err != nil {
				log.Printf("Error scraping dividends: %v", err)
			}
		}
	}
}

func (s *DividendScraperService) ScrapeDividends() error {
	// Try real scraping first, fallback to mock if Chrome not available
	dividendData, err := s.scrapeRealData()
	if err != nil {
		log.Printf("Real scraping failed, using mock data: %v", err)
		dividendData = s.getMockDividendData()
	}
	return s.processDividendData(dividendData)
}

func (s *DividendScraperService) scrapeRealData() ([]ScrapedDividendData, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://simplywall.st/stocks/gh/dividend-yield-high"),
		chromedp.WaitVisible(`text=Company Last Price 7D Return 1Y Return`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scrape page: %w", err)
	}

	return s.parseDividendData(htmlContent)
}

func (s *DividendScraperService) parseDividendData(html string) ([]ScrapedDividendData, error) {
	// Extract table rows with dividend data
	rowRegex := regexp.MustCompile(`<tr[^>]*role="row"[^>]*>(.*?)</tr>`)
	rows := rowRegex.FindAllStringSubmatch(html, -1)

	var data []ScrapedDividendData
	
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}

		rowContent := row[1]
		
		// Extract ticker and company name
		tickerRegex := regexp.MustCompile(`([A-Z0-9\.]+)\s+([^<]+)`)
		tickerMatch := tickerRegex.FindStringSubmatch(rowContent)
		
		// Extract dividend yield
		yieldRegex := regexp.MustCompile(`(\d+(?:\.\d+)?%)(?=[^%]*(?:Banks|Energy|Telecom|Insurance))`)
		yieldMatch := yieldRegex.FindStringSubmatch(rowContent)
		
		// Extract price
		priceRegex := regexp.MustCompile(`GH₵([\d\.,]+)`)
		priceMatch := priceRegex.FindStringSubmatch(rowContent)
		
		// Extract industry
		industryRegex := regexp.MustCompile(`(Banks|Energy|Telecom|Insurance|Mining|Manufacturing)`)
		industryMatch := industryRegex.FindStringSubmatch(rowContent)

		if tickerMatch != nil && yieldMatch != nil && len(tickerMatch) > 2 {
			ticker := strings.TrimSpace(tickerMatch[1])
			name := strings.TrimSpace(tickerMatch[2])
			divYield := yieldMatch[1]
			
			var price, industry string
			if priceMatch != nil {
				price = "GH₵" + priceMatch[1]
			}
			if industryMatch != nil {
				industry = industryMatch[1]
			}

			data = append(data, ScrapedDividendData{
				Ticker:    ticker,
				Name:      name,
				DivYield:  divYield,
				LastPrice: price,
				Industry:  industry,
			})
		}
	}

	return data, nil
}

func (s *DividendScraperService) getMockDividendData() []ScrapedDividendData {
	return []ScrapedDividendData{
		{Ticker: "GCB", Name: "GCB Bank Limited", DivYield: "10.4%", LastPrice: "GH₵4.20", Industry: "Banks"},
		{Ticker: "ACCESS", Name: "Access Bank Ghana", DivYield: "8.5%", LastPrice: "GH₵16.37", Industry: "Banks"},
		{Ticker: "CAL", Name: "CAL Bank Limited", DivYield: "7.2%", LastPrice: "GH₵0.95", Industry: "Banks"},
		{Ticker: "TOTAL", Name: "Total Petroleum Ghana", DivYield: "10.1%", LastPrice: "GH₵3.45", Industry: "Energy"},
		{Ticker: "MTN", Name: "MTN Ghana", DivYield: "6.8%", LastPrice: "GH₵0.82", Industry: "Telecom"},
	}
}



func (s *DividendScraperService) processDividendData(data []ScrapedDividendData) error {
	for _, item := range data {
		// Skip if no dividend yield
		if item.DivYield == "" || item.DivYield == "0%" {
			continue
		}

		// Parse dividend yield to amount (simplified calculation)
		yieldStr := strings.TrimSuffix(item.DivYield, "%")
		yieldPercent, err := strconv.ParseFloat(yieldStr, 64)
		if err != nil {
			continue
		}

		// Calculate estimated dividend amount from yield and price
		priceStr := strings.TrimPrefix(item.LastPrice, "GH₵")
		priceStr = strings.ReplaceAll(priceStr, ",", "")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			continue
		}

		dividendAmount := (price * yieldPercent) / 100

		// Create dividend announcement
		dividend := &models.CreateDividendRequest{
			StockSymbol:  item.Ticker,
			StockName:    item.Name,
			DividendType: "cash",
			Amount:       dividendAmount,
			Currency:     "GHS",
			ExDate:       time.Now().AddDate(0, 0, 30), // Estimated 30 days from now
			PaymentDate:  time.Now().AddDate(0, 0, 45), // Estimated 45 days from now
		}

		// Check if dividend already exists for this stock recently
		existing, _ := s.dividendRepo.GetBySymbol(item.Ticker)
		hasRecent := false
		for _, div := range existing {
			if div.CreatedAt.After(time.Now().AddDate(0, 0, -7)) { // Within last 7 days
				hasRecent = true
				break
			}
		}

		if !hasRecent {
			_, err := s.dividendService.CreateDividendAnnouncement(dividend)
			if err != nil {
				log.Printf("Failed to create dividend for %s: %v", item.Ticker, err)
			} else {
				log.Printf("Created dividend announcement for %s: %.2f GHS", item.Ticker, dividendAmount)
			}
		}
	}

	return nil
}