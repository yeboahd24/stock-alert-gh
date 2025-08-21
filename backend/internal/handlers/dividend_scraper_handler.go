package handlers

import (
	"net/http"

	"github.com/go-chi/render"

	"shares-alert-backend/internal/services"
)

type DividendScraperHandler struct {
	scraperService *services.DividendScraperService
}

func NewDividendScraperHandler(scraperService *services.DividendScraperService) *DividendScraperHandler {
	return &DividendScraperHandler{
		scraperService: scraperService,
	}
}

func (h *DividendScraperHandler) TriggerScraping(w http.ResponseWriter, r *http.Request) {
	go func() {
		// Run scraping in background
		if err := h.scraperService.ScrapeDividends(); err != nil {
			// Log error but don't fail the response
			return
		}
	}()

	render.JSON(w, r, map[string]string{
		"message": "Dividend scraping triggered successfully",
		"status":  "running",
	})
}