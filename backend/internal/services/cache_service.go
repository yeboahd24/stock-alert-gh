package services

import (
	"fmt"
	"log"
	"time"

	"shares-alert-backend/internal/cache"
)

type CacheService struct {
	cache *cache.RedisCache
}

func NewCacheService(redisCache *cache.RedisCache) *CacheService {
	return &CacheService{
		cache: redisCache,
	}
}

// InvalidateStockCache invalidates all stock-related cache entries
func (s *CacheService) InvalidateStockCache() error {
	patterns := []string{
		"stocks:*",
		"stock:*",
	}

	for _, pattern := range patterns {
		if err := s.cache.DeletePattern(pattern); err != nil {
			log.Printf("Failed to invalidate cache pattern %s: %v", pattern, err)
			return err
		}
	}

	log.Println("Stock cache invalidated successfully")
	return nil
}

// InvalidateStockSymbol invalidates cache for a specific stock symbol
func (s *CacheService) InvalidateStockSymbol(symbol string) error {
	keys := []string{
		fmt.Sprintf("stock:live:%s", symbol),
		fmt.Sprintf("stock:details:%s", symbol),
	}

	for _, key := range keys {
		if err := s.cache.Delete(key); err != nil {
			log.Printf("Failed to invalidate cache key %s: %v", key, err)
		}
	}

	// Also invalidate the all stocks cache since it contains this stock
	if err := s.cache.Delete("stocks:all"); err != nil {
		log.Printf("Failed to invalidate all stocks cache: %v", err)
	}

	log.Printf("Cache invalidated for stock symbol: %s", symbol)
	return nil
}

// GetCacheStats returns cache statistics
func (s *CacheService) GetCacheStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Check if some common keys exist
	stockKeys := []string{
		"stocks:all",
		"stock:live:MTN",
		"stock:live:ACCESS",
		"stock:live:GCB",
	}

	existingKeys := 0
	for _, key := range stockKeys {
		if s.cache.Exists(key) {
			existingKeys++
			if ttl, err := s.cache.GetTTL(key); err == nil {
				stats[key+"_ttl"] = ttl.String()
			}
		}
	}

	stats["cached_keys"] = existingKeys
	stats["total_checked"] = len(stockKeys)
	stats["timestamp"] = time.Now().Format(time.RFC3339)

	return stats
}

// WarmupCache pre-loads frequently accessed data
func (s *CacheService) WarmupCache(stockService *StockService) error {
	log.Println("Starting cache warmup...")

	// Warmup all stocks
	if _, err := stockService.GetAllStocks(); err != nil {
		log.Printf("Failed to warmup all stocks cache: %v", err)
	}

	// Warmup popular stocks
	popularStocks := []string{"MTN", "ACCESS", "GCB", "TOTAL", "GOIL"}
	for _, symbol := range popularStocks {
		if _, err := stockService.GetStock(symbol); err != nil {
			log.Printf("Failed to warmup stock %s: %v", symbol, err)
		}
		
		if _, err := stockService.GetStockDetails(symbol); err != nil {
			log.Printf("Failed to warmup stock details %s: %v", symbol, err)
		}
	}

	log.Println("Cache warmup completed")
	return nil
}