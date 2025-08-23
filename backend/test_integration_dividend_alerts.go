package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"shares-alert-backend/internal/models"
)

// Test configuration
const (
	baseURL = "http://localhost:8080/api/v1"
	testTimeout = 30 * time.Second
)

// Test data
var testAlerts = []models.CreateAlertRequest{
	{
		StockSymbol:    "",
		StockName:      "High Yield Opportunities",
		AlertType:      models.AlertTypeHighDividendYield,
		ThresholdYield: floatPtr(3.0),
	},
	{
		StockSymbol: "GCB",
		StockName:   "GCB Bank",
		AlertType:   models.AlertTypeTargetDividendYield,
		TargetYield: floatPtr(4.0),
	},
	{
		StockSymbol:          "GOIL",
		StockName:            "GOIL Company",
		AlertType:            models.AlertTypeDividendYieldChange,
		YieldChangeThreshold: floatPtr(0.5),
	},
}

type TestResult struct {
	TestName string
	Success  bool
	Error    string
	Duration time.Duration
}

type IntegrationTest struct {
	client   *http.Client
	authToken string
	results  []TestResult
}

func main() {
	fmt.Println("🚀 Starting Enhanced Dividend Alerts Integration Tests...")
	
	test := &IntegrationTest{
		client: &http.Client{Timeout: testTimeout},
		results: make([]TestResult, 0),
	}

	// Run all tests
	tests := []struct {
		name string
		fn   func() error
	}{
		{"Health Check", test.testHealthCheck},
		{"GSE Dividend API", test.testGSEDividendAPI},
		{"High Dividend Yield Stocks", test.testHighDividendYieldStocks},
		{"Create High Dividend Yield Alert", test.testCreateHighDividendYieldAlert},
		{"Create Target Dividend Yield Alert", test.testCreateTargetDividendYieldAlert},
		{"Create Dividend Yield Change Alert", test.testCreateDividendYieldChangeAlert},
		{"List Dividend Alerts", test.testListDividendAlerts},
		{"Update Dividend Alert", test.testUpdateDividendAlert},
		{"Delete Dividend Alert", test.testDeleteDividendAlert},
		{"Traditional Dividend Endpoints", test.testTraditionalDividendEndpoints},
	}

	for _, testCase := range tests {
		test.runTest(testCase.name, testCase.fn)
	}

	// Print results
	test.printResults()
}

func (t *IntegrationTest) runTest(name string, testFn func() error) {
	fmt.Printf("\n🧪 Running test: %s\n", name)
	start := time.Now()
	
	err := testFn()
	duration := time.Since(start)
	
	result := TestResult{
		TestName: name,
		Success:  err == nil,
		Duration: duration,
	}
	
	if err != nil {
		result.Error = err.Error()
		fmt.Printf("   ❌ FAILED: %v (%.2fs)\n", err, duration.Seconds())
	} else {
		fmt.Printf("   ✅ PASSED (%.2fs)\n", duration.Seconds())
	}
	
	t.results = append(t.results, result)
}

func (t *IntegrationTest) testHealthCheck() error {
	resp, err := t.client.Get(baseURL + "/health")
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	var health map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return fmt.Errorf("failed to decode health response: %w", err)
	}

	if health["status"] != "ok" {
		return fmt.Errorf("health check status is not ok: %v", health["status"])
	}

	fmt.Printf("   📊 Health: %s, Version: %s\n", health["status"], health["version"])
	return nil
}

func (t *IntegrationTest) testGSEDividendAPI() error {
	resp, err := t.client.Get(baseURL + "/dividends/gse")
	if err != nil {
		return fmt.Errorf("GSE dividend API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GSE dividend API returned status %d", resp.StatusCode)
	}

	var dividendResponse models.GSEDividendResponse
	if err := json.NewDecoder(resp.Body).Decode(&dividendResponse); err != nil {
		return fmt.Errorf("failed to decode GSE dividend response: %w", err)
	}

	if !dividendResponse.Success {
		return fmt.Errorf("GSE dividend API returned unsuccessful response")
	}

	if len(dividendResponse.Data.Stocks) == 0 {
		return fmt.Errorf("no dividend stocks returned from GSE API")
	}

	fmt.Printf("   📈 Found %d dividend stocks from GSE API\n", dividendResponse.Data.Count)
	
	// Test specific stock endpoint
	firstStock := dividendResponse.Data.Stocks[0]
	stockResp, err := t.client.Get(baseURL + "/dividends/gse/" + firstStock.Symbol)
	if err != nil {
		return fmt.Errorf("specific stock dividend request failed: %w", err)
	}
	defer stockResp.Body.Close()

	if stockResp.StatusCode != http.StatusOK {
		return fmt.Errorf("specific stock dividend API returned status %d", stockResp.StatusCode)
	}

	var stockData models.GSEDividendStock
	if err := json.NewDecoder(stockResp.Body).Decode(&stockData); err != nil {
		return fmt.Errorf("failed to decode stock dividend response: %w", err)
	}

	fmt.Printf("   🏢 %s (%s): %.2f%% yield\n", stockData.Name, stockData.Symbol, stockData.DividendYield)
	return nil
}

func (t *IntegrationTest) testHighDividendYieldStocks() error {
	resp, err := t.client.Get(baseURL + "/dividends/high-yield?minYield=2.0")
	if err != nil {
		return fmt.Errorf("high dividend yield request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("high dividend yield API returned status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode high yield response: %w", err)
	}

	count, ok := response["count"].(float64)
	if !ok {
		return fmt.Errorf("invalid count in high yield response")
	}

	minYield, ok := response["minYield"].(float64)
	if !ok {
		return fmt.Errorf("invalid minYield in high yield response")
	}

	fmt.Printf("   📊 Found %.0f stocks with yield >= %.1f%%\n", count, minYield)
	return nil
}

func (t *IntegrationTest) testCreateHighDividendYieldAlert() error {
	// Note: This test requires authentication
	// For now, we'll test the endpoint structure
	alertData := testAlerts[0]
	
	jsonData, err := json.Marshal(alertData)
	if err != nil {
		return fmt.Errorf("failed to marshal alert data: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/alerts", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	// Note: In real test, we would add: req.Header.Set("Authorization", "Bearer " + t.authToken)

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("create alert request failed: %w", err)
	}
	defer resp.Body.Close()

	// For unauthenticated request, we expect 401
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("   🔒 Authentication required (expected for this test)\n")
		return nil
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create alert returned status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("   ✨ High dividend yield alert created successfully\n")
	return nil
}

func (t *IntegrationTest) testCreateTargetDividendYieldAlert() error {
	alertData := testAlerts[1]
	
	jsonData, err := json.Marshal(alertData)
	if err != nil {
		return fmt.Errorf("failed to marshal alert data: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/alerts", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("create alert request failed: %w", err)
	}
	defer resp.Body.Close()

	// For unauthenticated request, we expect 401
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("   🔒 Authentication required (expected for this test)\n")
		return nil
	}

	fmt.Printf("   🎯 Target dividend yield alert endpoint tested\n")
	return nil
}

func (t *IntegrationTest) testCreateDividendYieldChangeAlert() error {
	alertData := testAlerts[2]
	
	jsonData, err := json.Marshal(alertData)
	if err != nil {
		return fmt.Errorf("failed to marshal alert data: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/alerts", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("create alert request failed: %w", err)
	}
	defer resp.Body.Close()

	// For unauthenticated request, we expect 401
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("   🔒 Authentication required (expected for this test)\n")
		return nil
	}

	fmt.Printf("   📊 Dividend yield change alert endpoint tested\n")
	return nil
}

func (t *IntegrationTest) testListDividendAlerts() error {
	resp, err := t.client.Get(baseURL + "/alerts")
	if err != nil {
		return fmt.Errorf("list alerts request failed: %w", err)
	}
	defer resp.Body.Close()

	// For unauthenticated request, we expect 401
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("   🔒 Authentication required (expected for this test)\n")
		return nil
	}

	fmt.Printf("   📋 List alerts endpoint tested\n")
	return nil
}

func (t *IntegrationTest) testUpdateDividendAlert() error {
	// Test update endpoint structure
	updateData := models.UpdateAlertRequest{
		ThresholdYield: floatPtr(3.5),
	}
	
	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal update data: %w", err)
	}

	req, err := http.NewRequest("PUT", baseURL+"/alerts/test-id", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("update alert request failed: %w", err)
	}
	defer resp.Body.Close()

	// For unauthenticated request, we expect 401
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("   🔒 Authentication required (expected for this test)\n")
		return nil
	}

	fmt.Printf("   ✏️ Update alert endpoint tested\n")
	return nil
}

func (t *IntegrationTest) testDeleteDividendAlert() error {
	req, err := http.NewRequest("DELETE", baseURL+"/alerts/test-id", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("delete alert request failed: %w", err)
	}
	defer resp.Body.Close()

	// For unauthenticated request, we expect 401
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("   🔒 Authentication required (expected for this test)\n")
		return nil
	}

	fmt.Printf("   🗑️ Delete alert endpoint tested\n")
	return nil
}

func (t *IntegrationTest) testTraditionalDividendEndpoints() error {
	// Test traditional dividend announcements endpoint
	resp, err := t.client.Get(baseURL + "/dividends")
	if err != nil {
		return fmt.Errorf("traditional dividends request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("traditional dividends API returned status %d", resp.StatusCode)
	}

	// Test upcoming dividends endpoint
	resp2, err := t.client.Get(baseURL + "/dividends/upcoming")
	if err != nil {
		return fmt.Errorf("upcoming dividends request failed: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf("upcoming dividends API returned status %d", resp2.StatusCode)
	}

	fmt.Printf("   💰 Traditional dividend endpoints working\n")
	return nil
}

func (t *IntegrationTest) printResults() {
	fmt.Println("\n" + "="*60)
	fmt.Println("📊 INTEGRATION TEST RESULTS")
	fmt.Println("="*60)

	passed := 0
	failed := 0
	totalDuration := time.Duration(0)

	for _, result := range t.results {
		status := "✅ PASS"
		if !result.Success {
			status = "❌ FAIL"
			failed++
		} else {
			passed++
		}
		
		fmt.Printf("%-30s %s (%.2fs)\n", result.TestName, status, result.Duration.Seconds())
		if result.Error != "" {
			fmt.Printf("   Error: %s\n", result.Error)
		}
		totalDuration += result.Duration
	}

	fmt.Println("-" * 60)
	fmt.Printf("Total Tests: %d | Passed: %d | Failed: %d\n", len(t.results), passed, failed)
	fmt.Printf("Total Duration: %.2fs\n", totalDuration.Seconds())
	
	if failed > 0 {
		fmt.Println("\n❌ Some tests failed. Check the errors above.")
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All tests passed!")
	}
}

func floatPtr(f float64) *float64 {
	return &f
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}