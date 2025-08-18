package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"shares-alert-backend/internal/cache"
	"shares-alert-backend/internal/config"
	"shares-alert-backend/internal/httpclient"
	"shares-alert-backend/internal/models"
)

type StockService struct {
	config     *config.ExternalConfig
	cache      *cache.RedisCache
	httpClient *http.Client
	cacheTTL   time.Duration
}

func NewStockService(cfg *config.ExternalConfig, redisCache *cache.RedisCache, cacheTTL time.Duration) *StockService {
	return &StockService{
		config:     cfg,
		cache:      redisCache,
		cacheTTL:   cacheTTL,
		httpClient: httpclient.CreateClientWithTimeout(20 * time.Second),
	}
}

func (s *StockService) GetAllStocks() ([]models.EnhancedStock, error) {
	cacheKey := "stocks:all"
	
	// Try to get from cache first
	var cachedStocks []models.EnhancedStock
	if err := s.cache.Get(cacheKey, &cachedStocks); err == nil {
		log.Printf("Cache hit for all stocks")
		return cachedStocks, nil
	}

	log.Printf("Cache miss for all stocks, fetching from API")
	
	resp, err := s.fetchWithProxy("/live")
	if err != nil {
		log.Printf("External API unavailable, using mock data. Error: %v", err)
		mockStocks := s.getMockStocks()
		// Cache mock data for a shorter time
		s.cache.Set(cacheKey, mockStocks, 1*time.Minute)
		return mockStocks, nil
	}
	defer resp.Body.Close()

	var stocks []models.StockLive
	if err := json.NewDecoder(resp.Body).Decode(&stocks); err != nil {
		log.Printf("Failed to parse external API response, using mock data. Error: %v", err)
		mockStocks := s.getMockStocks()
		s.cache.Set(cacheKey, mockStocks, 1*time.Minute)
		return mockStocks, nil
	}

	enhancedStocks := s.convertToEnhancedStocks(stocks)
	
	// Cache the successful response
	if err := s.cache.Set(cacheKey, enhancedStocks, s.cacheTTL); err != nil {
		log.Printf("Failed to cache all stocks: %v", err)
	}

	return enhancedStocks, nil
}

func (s *StockService) GetStock(symbol string) (*models.EnhancedStock, error) {
	symbol = strings.ToUpper(symbol)
	cacheKey := fmt.Sprintf("stock:live:%s", symbol)
	
	// Try to get from cache first
	var cachedStock models.EnhancedStock
	if err := s.cache.Get(cacheKey, &cachedStock); err == nil {
		log.Printf("Cache hit for stock %s", symbol)
		return &cachedStock, nil
	}

	log.Printf("Cache miss for stock %s, fetching from API", symbol)

	resp, err := s.fetchWithProxy("/live/" + symbol)
	if err != nil {
		// Fallback to mock data
		mockStocks := s.getMockStocks()
		for _, stock := range mockStocks {
			if strings.ToUpper(stock.Symbol) == symbol {
				// Cache mock data for shorter time
				s.cache.Set(cacheKey, stock, 1*time.Minute)
				return &stock, nil
			}
		}
		return nil, fmt.Errorf("stock not found")
	}
	defer resp.Body.Close()

	var stock models.StockLive
	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
		// Fallback to mock data
		mockStocks := s.getMockStocks()
		for _, mockStock := range mockStocks {
			if strings.ToUpper(mockStock.Symbol) == symbol {
				s.cache.Set(cacheKey, mockStock, 1*time.Minute)
				return &mockStock, nil
			}
		}
		return nil, fmt.Errorf("stock not found")
	}

	enhanced := s.convertToEnhancedStock(stock)
	
	// Cache the successful response
	if err := s.cache.Set(cacheKey, enhanced, s.cacheTTL); err != nil {
		log.Printf("Failed to cache stock %s: %v", symbol, err)
	}

	return &enhanced, nil
}

func (s *StockService) GetStockDetails(symbol string) (*models.DetailedStock, error) {
	symbol = strings.ToUpper(symbol)
	cacheKey := fmt.Sprintf("stock:details:%s", symbol)
	
	// Try to get from cache first
	var cachedStock models.DetailedStock
	if err := s.cache.Get(cacheKey, &cachedStock); err == nil {
		log.Printf("Cache hit for stock details %s", symbol)
		return &cachedStock, nil
	}

	log.Printf("Cache miss for stock details %s, fetching from API", symbol)

	resp, err := s.fetchWithProxy("/equities/" + symbol)
	if err != nil {
		mockStock := s.getMockDetailedStock(symbol)
		if mockStock != nil {
			// Cache mock data for shorter time
			s.cache.Set(cacheKey, *mockStock, 1*time.Minute)
			return mockStock, nil
		}
		return nil, fmt.Errorf("stock not found")
	}
	defer resp.Body.Close()

	var equity models.StockEquity
	if err := json.NewDecoder(resp.Body).Decode(&equity); err != nil {
		mockStock := s.getMockDetailedStock(symbol)
		if mockStock != nil {
			s.cache.Set(cacheKey, *mockStock, 1*time.Minute)
			return mockStock, nil
		}
		return nil, fmt.Errorf("stock not found")
	}

	// Get live data for current price, change, and volume
	var currentPrice float64 = equity.Price
	var change float64 = 0
	var volume int64 = 0
	var changePercent float64 = 0

	if liveResp, err := s.fetchWithProxy("/live/" + symbol); err == nil {
		defer liveResp.Body.Close()
		if liveResp.StatusCode == 200 {
			var liveStock models.StockLive
			if json.NewDecoder(liveResp.Body).Decode(&liveStock) == nil {
				currentPrice = liveStock.Price
				change = liveStock.Change
				volume = liveStock.Volume

				if currentPrice > 0 {
					previousClose := currentPrice - change
					if previousClose > 0 {
						changePercent = (change / previousClose) * 100
					}
				}
			}
		}
	}

	detailedStock := &models.DetailedStock{
		Symbol:        equity.Name,
		Name:          equity.Company.Name,
		CurrentPrice:  currentPrice,
		PreviousClose: currentPrice - change,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        volume,
		LastUpdated:   time.Now(),
		MarketCap:     equity.Capital,
		Shares:        equity.Shares,
		Sector:        equity.Company.Sector,
		Industry:      equity.Company.Industry,
		DPS:           equity.DPS,
		EPS:           equity.EPS,
		Company:       equity.Company,
	}

	// Cache the successful response
	if err := s.cache.Set(cacheKey, *detailedStock, s.cacheTTL); err != nil {
		log.Printf("Failed to cache stock details %s: %v", symbol, err)
	}

	return detailedStock, nil
}

func (s *StockService) fetchWithProxy(endpoint string) (*http.Response, error) {
	// Try direct first
	resp, err := s.httpClient.Get(s.config.GSEBaseURL + endpoint)
	if err == nil && resp.StatusCode == 200 {
		log.Printf("Direct API success: %s", s.config.GSEBaseURL+endpoint)
		return resp, nil
	}
	if resp != nil {
		resp.Body.Close()
	}
	log.Printf("Direct API failed: %v, trying proxy...", err)

	// Try through CORS proxy
	proxyURL := s.config.ProxyURL + s.config.GSEBaseURL + endpoint
	resp, err = s.httpClient.Get(proxyURL)
	if err == nil && resp.StatusCode == 200 {
		log.Printf("Proxy API success: %s", proxyURL)
		return resp, nil
	}
	if resp != nil {
		resp.Body.Close()
	}
	log.Printf("Proxy API also failed: %v", err)

	return nil, err
}

func (s *StockService) convertToEnhancedStocks(stocks []models.StockLive) []models.EnhancedStock {
	enhanced := make([]models.EnhancedStock, len(stocks))
	for i, stock := range stocks {
		enhanced[i] = s.convertToEnhancedStock(stock)
	}
	return enhanced
}

func (s *StockService) convertToEnhancedStock(stock models.StockLive) models.EnhancedStock {
	changePercent := 0.0
	if stock.Price > 0 {
		previousClose := stock.Price - stock.Change
		if previousClose > 0 {
			changePercent = (stock.Change / previousClose) * 100
		}
	}

	return models.EnhancedStock{
		Symbol:        stock.Name,
		Name:          stock.Name,
		CurrentPrice:  stock.Price,
		PreviousClose: stock.Price - stock.Change,
		Change:        stock.Change,
		ChangePercent: changePercent,
		Volume:        stock.Volume,
		LastUpdated:   time.Now(),
	}
}

func (s *StockService) getMockStocks() []models.EnhancedStock {
	return []models.EnhancedStock{
		{
			Symbol:        "ACCESS",
			Name:          "Access Bank Ghana Plc",
			CurrentPrice:  16.37,
			PreviousClose: 16.37,
			Change:        0.0,
			ChangePercent: 0.0,
			Volume:        0,
			LastUpdated:   time.Now(),
			MarketCap:     func() *float64 { v := 2100000000.0; return &v }(),
			Sector:        func() *string { v := "Financial Services"; return &v }(),
			Industry:      func() *string { v := "Banking"; return &v }(),
		},
		{
			Symbol:        "GCB",
			Name:          "GCB Bank Limited",
			CurrentPrice:  4.20,
			PreviousClose: 4.15,
			Change:        0.05,
			ChangePercent: 1.20,
			Volume:        67000,
			LastUpdated:   time.Now(),
			MarketCap:     func() *float64 { v := 1800000000.0; return &v }(),
			Sector:        func() *string { v := "Financial Services"; return &v }(),
			Industry:      func() *string { v := "Banking"; return &v }(),
		},
		{
			Symbol:        "MTN",
			Name:          "MTN Ghana",
			CurrentPrice:  0.82,
			PreviousClose: 0.80,
			Change:        0.02,
			ChangePercent: 2.5,
			Volume:        125000,
			LastUpdated:   time.Now(),
			MarketCap:     func() *float64 { v := 1500000000.0; return &v }(),
			Sector:        func() *string { v := "Telecommunications"; return &v }(),
			Industry:      func() *string { v := "Mobile Networks"; return &v }(),
		},
	}
}

func (s *StockService) getMockDetailedStock(symbol string) *models.DetailedStock {
	mockStocks := map[string]models.DetailedStock{
		"MTN": {
			Symbol:        "MTN",
			Name:          "MTN Ghana",
			CurrentPrice:  0.82,
			PreviousClose: 0.80,
			Change:        0.02,
			ChangePercent: 2.5,
			Volume:        125000,
			LastUpdated:   time.Now(),
			MarketCap:     1500000000.0,
			Shares:        1829268293,
			Sector:        "Telecommunications",
			Industry:      "Mobile Networks",
			DPS:           func() *float64 { v := 0.05; return &v }(),
			EPS:           func() *float64 { v := 0.12; return &v }(),
			Company: models.Company{
				Name:      "MTN Ghana",
				Address:   "Accra, Ghana",
				Email:     "info@mtn.com.gh",
				Telephone: "+233-244-300-000",
				Website:   "https://www.mtn.com.gh",
				Sector:    "Telecommunications",
				Industry:  "Mobile Networks",
				Directors: []string{"Selorm Adadevoh", "Ebenezer Asante"},
			},
		},
	}

	if stock, exists := mockStocks[symbol]; exists {
		return &stock
	}
	return nil
}