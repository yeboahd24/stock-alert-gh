package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"shares-alert-backend/internal/models"
)

func main() {
	// Test IPO creation
	testCreateIPO()
	
	// Wait a moment
	time.Sleep(2 * time.Second)
	
	// Test fetching IPOs
	testGetIPOs()
}

func testCreateIPO() {
	fmt.Println("Testing IPO creation...")
	
	ipoData := models.CreateIPORequest{
		CompanyName: "Ghana Tech Solutions Ltd",
		Symbol:      "GTECH",
		Sector:      "Technology",
		OfferPrice:  15.50,
		ListingDate: time.Now().Add(30 * 24 * time.Hour), // 30 days from now
	}
	
	jsonData, err := json.Marshal(ipoData)
	if err != nil {
		log.Printf("Error marshaling IPO data: %v", err)
		return
	}
	
	resp, err := http.Post("http://localhost:8080/api/v1/ipos", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating IPO: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ IPO created successfully!")
	} else {
		fmt.Printf("❌ Failed to create IPO. Status: %d\n", resp.StatusCode)
	}
}

func testGetIPOs() {
	fmt.Println("Testing IPO retrieval...")
	
	resp, err := http.Get("http://localhost:8080/api/v1/ipos")
	if err != nil {
		log.Printf("Error fetching IPOs: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		var ipos []models.IPOAnnouncement
		if err := json.NewDecoder(resp.Body).Decode(&ipos); err != nil {
			log.Printf("Error decoding IPOs: %v", err)
			return
		}
		
		fmt.Printf("✅ Retrieved %d IPO(s):\n", len(ipos))
		for _, ipo := range ipos {
			fmt.Printf("  - %s (%s): GH₵%.2f, listing %s\n", 
				ipo.CompanyName, ipo.Symbol, ipo.OfferPrice, ipo.ListingDate.Format("2006-01-02"))
		}
	} else {
		fmt.Printf("❌ Failed to fetch IPOs. Status: %d\n", resp.StatusCode)
	}
}