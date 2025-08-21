package services

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"

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
	// Check if scraping is enabled via environment variable
	if os.Getenv("ENABLE_SCRAPING") != "true" {
		log.Println("Scraping disabled via ENABLE_SCRAPING env var, using mock data only")
		return
	}

	// Use longer interval to reduce memory pressure - scrape every 6 hours instead of daily
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	log.Println("Starting dividend scraping service (every 6 hours)...")

	// Wait 5 minutes after startup before first scrape to let app stabilize
	time.Sleep(5 * time.Minute)
	
	// Run initial scrape
	log.Println("Running initial dividend scraping...")
	if err := s.ScrapeDividends(); err != nil {
		log.Printf("Error in initial dividend scraping: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			// Force garbage collection before scraping to free memory
			runtime.GC()
			
			log.Println("Starting scheduled dividend scraping...")
			if err := s.ScrapeDividends(); err != nil {
				log.Printf("Error scraping dividends: %v", err)
			}
			
			// Force garbage collection after scraping
			runtime.GC()
			log.Println("Dividend scraping cycle completed")
		}
	}
}

func (s *DividendScraperService) ScrapeDividends() error {
	// Check memory before starting
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.Alloc > 300*1024*1024 { // 300MB limit before starting
		log.Printf("Memory usage too high (%d MB), using mock data instead", m.Alloc/1024/1024)
		dividendData := s.getMockDividendData()
		return s.processDividendData(dividendData)
	}
	
	// Try real scraping first, fallback to mock if Chrome not available
	dividendData, err := s.scrapeRealData()
	if err != nil {
		log.Printf("Real scraping failed, using mock data: %v", err)
		dividendData = s.getMockDividendData()
	}
	return s.processDividendData(dividendData)
}

func (s *DividendScraperService) scrapeRealData() ([]ScrapedDividendData, error) {
	// Launch browser with rod - Ultra memory-optimized configuration
	l := launcher.New().
		Headless(true).
		NoSandbox(true).
		Set("disable-gpu").
		Set("disable-dev-shm-usage").
		Set("disable-extensions").
		Set("disable-background-timer-throttling").
		Set("disable-backgrounding-occluded-windows").
		Set("disable-renderer-backgrounding").
		Set("disable-setuid-sandbox").
		Set("no-first-run").
		Set("no-zygote").
		Set("single-process").
		Set("disable-web-security").
		Set("memory-pressure-off").
		Set("max_old_space_size=128").
		Set("disable-background-networking").
		Set("disable-default-apps").
		Set("disable-sync").
		Set("disable-plugins").
		Set("disable-images").
		Set("disable-javascript").
		Set("disable-plugins-discovery").
		Set("disable-preconnect").
		Set("disable-translate").
		Set("no-pings").
		Set("no-referrers").
		Set("disable-client-side-phishing-detection").
		Set("disable-component-extensions-with-background-pages").
		Set("disable-ipc-flooding-protection")

	// Use system Chromium if available (Docker environment)
	if chromiumPath := os.Getenv("ROD_LAUNCHER_BIN"); chromiumPath != "" {
		l = l.Bin(chromiumPath)
	}

	// Try to launch browser with error handling
	url, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(url)
	err = browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}
	
	// Ensure browser cleanup with explicit close
	defer func() {
		if browser != nil {
			browser.MustClose()
		}
		// Force garbage collection after browser cleanup
		runtime.GC()
	}()

	// Navigate to page with better error handling
	page := browser.MustPage()
	
	// Set user agent to avoid bot detection
	page = page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})
	
	// Ensure page cleanup
	defer func() {
		if page != nil {
			page.MustClose()
		}
	}()

	err = page.Navigate("https://simplywall.st/stocks/gh/dividend-yield-high")
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to page: %w", err)
	}

	// Wait for page to load completely
	err = page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for page load: %w", err)
	}

	// Minimal wait time to save memory
	time.Sleep(2 * time.Second)

	// Try multiple selectors to find stock cards
	var cards rod.Elements
	selectors := []string{
		"div[data-testid='screener-card']",
		"[data-testid='screener-card']",
		"div[data-cy='stock-card']",
		"article[data-testid*='stock']",
		".stock-card",
		"div[class*='card'][class*='stock']",
		"div[class*='screener']",
	}

	for _, selector := range selectors {
		cards, err = page.Elements(selector)
		if err == nil && len(cards) > 0 {
			log.Printf("Found %d cards using selector: %s", len(cards), selector)
			break
		}
	}

	if len(cards) == 0 {
		// Try to get page content for debugging
		html, _ := page.HTML()
		log.Printf("No cards found. Page HTML length: %d", len(html))
		
		// Look for any elements that might contain stock data
		allDivs, _ := page.Elements("div")
		log.Printf("Total div elements found: %d", len(allDivs))
		
		return nil, fmt.Errorf("no stock cards found on page")
	}

	var data []ScrapedDividendData

	for i, card := range cards {
		if i >= 5 { // Further reduced limit to save memory - only top 5 stocks
			break
		}
		
		// Check memory usage more frequently
		if i%2 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			if m.Alloc > 350*1024*1024 { // Lower 350MB limit
				log.Printf("Memory usage too high (%d MB), stopping scraping early", m.Alloc/1024/1024)
				break
			}
		}

		// Try multiple ways to extract company name
		var companyName string
		nameSelectors := []string{"h2", "h3", "h1", "[data-testid='company-name']", ".company-name"}
		for _, nameSelector := range nameSelectors {
			if nameEl, err := card.Element(nameSelector); err == nil {
				companyName = nameEl.MustText()
				break
			}
		}

		if companyName == "" {
			continue
		}

		// Extract ticker from company name or data attributes
		ticker := s.extractTicker(companyName)
		if ticker == "" {
			// Try to extract from data attributes
			if tickerAttr, err := card.Attribute("data-ticker"); err == nil && *tickerAttr != "" {
				ticker = *tickerAttr
			}
		}

		if ticker == "" {
			continue
		}

		// Extract sector with multiple approaches
		sector := ""
		sectorSelectors := []string{
			"div:contains('Sector')",
			"span:contains('Sector')", 
			"[data-testid='sector']",
			".sector",
		}
		for _, sectorSelector := range sectorSelectors {
			if sectorEl, err := card.ElementR("*", sectorSelector); err == nil {
				sector = sectorEl.MustText()
				break
			}
		}

		// Extract dividend yield with multiple approaches
		yield := ""
		yieldSelectors := []string{
			"div:contains('Dividend')",
			"span:contains('Dividend')",
			"div:contains('%')",
			"[data-testid='dividend-yield']",
			".dividend-yield",
		}
		for _, yieldSelector := range yieldSelectors {
			if yieldEl, err := card.ElementR("*", yieldSelector); err == nil {
				yieldText := yieldEl.MustText()
				if strings.Contains(yieldText, "%") {
					yield = yieldText
					break
				}
			}
		}

		// Extract price with multiple approaches
		price := ""
		priceSelectors := []string{
			"div:contains('Price')",
			"span:contains('GH₵')",
			"div:contains('GH₵')",
			"[data-testid='price']",
			".price",
		}
		for _, priceSelector := range priceSelectors {
			if priceEl, err := card.ElementR("*", priceSelector); err == nil {
				priceText := priceEl.MustText()
				if strings.Contains(priceText, "GH₵") || strings.Contains(priceText, "$") {
					price = priceText
					break
				}
			}
		}

		if yield != "" && yield != "0%" {
			data = append(data, ScrapedDividendData{
				Ticker:    ticker,
				Name:      companyName,
				DivYield:  yield,
				LastPrice: price,
				Industry:  sector,
			})
			
			log.Printf("Scraped: %s (%s) - Yield: %s, Price: %s, Sector: %s", 
				companyName, ticker, yield, price, sector)
		}
	}

	log.Printf("Successfully scraped %d dividend stocks", len(data))
	return data, nil
}

// extractTicker extracts stock ticker from company name
func (s *DividendScraperService) extractTicker(companyName string) string {
	// Common Ghana stock tickers mapping
	tickerMap := map[string]string{
		"GCB Bank":           "GCB",
		"Access Bank":        "ACCESS",
		"CAL Bank":           "CAL",
		"Total Petroleum":    "TOTAL",
		"MTN Ghana":          "MTN",
		"Ecobank":            "EBG",
		"Standard Chartered": "SCB",
		"Societe Generale":   "SG-GH",
	}

	// Try exact match first
	for name, ticker := range tickerMap {
		if strings.Contains(strings.ToLower(companyName), strings.ToLower(name)) {
			return ticker
		}
	}

	// Extract ticker from parentheses if present
	tickerRegex := regexp.MustCompile(`\(([A-Z0-9\.]+)\)`)
	if match := tickerRegex.FindStringSubmatch(companyName); len(match) > 1 {
		return match[1]
	}

	// Extract first word if it looks like a ticker
	words := strings.Fields(companyName)
	if len(words) > 0 {
		firstWord := strings.ToUpper(words[0])
		if len(firstWord) <= 6 && regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(firstWord) {
			return firstWord
		}
	}

	return ""
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

