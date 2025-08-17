package handlers

import (
	"net/http"

	"github.com/go-chi/render"

	"shares-alert-backend/internal/services"
)

type CacheHandler struct {
	cacheService *services.CacheService
	stockService *services.StockService
}

func NewCacheHandler(cacheService *services.CacheService, stockService *services.StockService) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
		stockService: stockService,
	}
}

// GetCacheStats returns cache statistics
func (h *CacheHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := h.cacheService.GetCacheStats()
	render.JSON(w, r, map[string]interface{}{
		"status": "ok",
		"cache":  stats,
	})
}

// InvalidateCache clears all stock cache
func (h *CacheHandler) InvalidateCache(w http.ResponseWriter, r *http.Request) {
	if err := h.cacheService.InvalidateStockCache(); err != nil {
		http.Error(w, "Failed to invalidate cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]string{
		"status":  "ok",
		"message": "Cache invalidated successfully",
	})
}

// WarmupCache pre-loads frequently accessed data
func (h *CacheHandler) WarmupCache(w http.ResponseWriter, r *http.Request) {
	if err := h.cacheService.WarmupCache(h.stockService); err != nil {
		http.Error(w, "Failed to warmup cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]string{
		"status":  "ok",
		"message": "Cache warmup completed successfully",
	})
}