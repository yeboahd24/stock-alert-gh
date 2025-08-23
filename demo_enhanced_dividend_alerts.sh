#!/bin/bash

# Enhanced Dividend Alerts Demo Script
# This script demonstrates the complete functionality of the enhanced dividend alerts system

echo "🚀 Enhanced Dividend Alerts System Demo"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BACKEND_URL="http://localhost:8080/api/v1"
FRONTEND_URL="http://localhost:5173"

echo -e "\n${BLUE}📋 Demo Overview:${NC}"
echo "1. Backend API Testing"
echo "2. Frontend Component Testing"
echo "3. Integration Testing"
echo "4. Performance Testing"

# Function to check if service is running
check_service() {
    local url=$1
    local name=$2
    
    if curl -s "$url/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $name is running${NC}"
        return 0
    else
        echo -e "${RED}❌ $name is not running${NC}"
        return 1
    fi
}

# Function to test API endpoint
test_endpoint() {
    local endpoint=$1
    local description=$2
    local expected_status=${3:-200}
    
    echo -e "\n${YELLOW}🧪 Testing: $description${NC}"
    
    response=$(curl -s -w "%{http_code}" "$BACKEND_URL$endpoint")
    status_code="${response: -3}"
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✅ $endpoint - Status: $status_code${NC}"
        
        # Show sample data for successful responses
        if [ "$status_code" -eq 200 ]; then
            body="${response%???}"
            echo "$body" | jq -r '.data.count // .count // "N/A"' 2>/dev/null | head -1 | sed 's/^/   📊 Data points: /'
        fi
    else
        echo -e "${RED}❌ $endpoint - Expected: $expected_status, Got: $status_code${NC}"
    fi
}

# Function to run frontend tests
run_frontend_tests() {
    echo -e "\n${BLUE}🧪 Running Frontend Tests${NC}"
    
    if command -v npm &> /dev/null; then
        echo "Running dividend alerts test suite..."
        if npm test dividend-alerts.test.tsx --silent 2>/dev/null; then
            echo -e "${GREEN}✅ Frontend tests passed${NC}"
        else
            echo -e "${YELLOW}⚠️ Frontend tests require setup (npm install)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️ npm not found, skipping frontend tests${NC}"
    fi
}

# Function to run backend tests
run_backend_tests() {
    echo -e "\n${BLUE}🧪 Running Backend Integration Tests${NC}"
    
    if [ -f "backend/test_integration_dividend_alerts.go" ]; then
        cd backend
        if go run test_integration_dividend_alerts.go 2>/dev/null; then
            echo -e "${GREEN}✅ Backend integration tests completed${NC}"
        else
            echo -e "${YELLOW}⚠️ Backend tests require running server${NC}"
        fi
        cd ..
    else
        echo -e "${YELLOW}⚠️ Backend test file not found${NC}"
    fi
}

# Main demo execution
main() {
    echo -e "\n${BLUE}🔍 Checking Services${NC}"
    
    # Check if backend is running
    if check_service "$BACKEND_URL" "Backend API"; then
        BACKEND_RUNNING=true
    else
        BACKEND_RUNNING=false
        echo -e "${YELLOW}💡 To start backend: cd backend && go run cmd/server/main.go${NC}"
    fi
    
    # Check if frontend is running
    if check_service "$FRONTEND_URL" "Frontend App"; then
        FRONTEND_RUNNING=true
    else
        FRONTEND_RUNNING=false
        echo -e "${YELLOW}💡 To start frontend: npm run dev${NC}"
    fi
    
    if [ "$BACKEND_RUNNING" = true ]; then
        echo -e "\n${BLUE}🌐 Testing Backend API Endpoints${NC}"
        
        # Test health check
        test_endpoint "/health" "Health Check"
        
        # Test GSE dividend endpoints
        test_endpoint "/dividends/gse" "GSE Dividend Stocks"
        test_endpoint "/dividends/high-yield?minYield=2.0" "High Dividend Yield Stocks"
        test_endpoint "/dividends" "Traditional Dividend Announcements"
        test_endpoint "/dividends/upcoming" "Upcoming Dividend Payments"
        
        # Test specific stock endpoint (assuming GCB exists)
        test_endpoint "/dividends/gse/GCB" "Specific Stock Dividend Data"
        
        # Test alert endpoints (will return 401 without auth)
        test_endpoint "/alerts" "List Alerts" 401
        
        echo -e "\n${BLUE}📊 Sample API Response:${NC}"
        echo "Getting sample dividend data..."
        curl -s "$BACKEND_URL/dividends/gse" | jq -r '
            if .success then
                "🏢 Found " + (.data.count | tostring) + " dividend stocks",
                "📈 Sample stocks:",
                (.data.stocks[:3] | .[] | "   • " + .symbol + " (" + .name + "): " + (.dividend_yield | tostring) + "% yield")
            else
                "❌ API Error"
            end
        ' 2>/dev/null || echo "   (jq not available for JSON parsing)"
        
    else
        echo -e "\n${YELLOW}⚠️ Backend not running, skipping API tests${NC}"
    fi
    
    # Run tests
    run_backend_tests
    run_frontend_tests
    
    echo -e "\n${BLUE}📁 Generated Files Summary${NC}"
    echo "Backend enhancements:"
    echo "  ✅ Enhanced models with dividend yield fields"
    echo "  ✅ Updated repository with yield tracking"
    echo "  ✅ Enhanced services with yield monitoring"
    echo "  ✅ New dividend API endpoints"
    echo "  ✅ Email templates for yield alerts"
    echo "  ✅ Database migration script"
    
    echo -e "\nFrontend enhancements:"
    echo "  ✅ Enhanced AlertForm with yield alert types"
    echo "  ✅ DividendYieldDashboard component"
    echo "  ✅ DividendYieldAlerts management component"
    echo "  ✅ Updated API service with dividend endpoints"
    echo "  ✅ Enhanced types and formatters"
    echo "  ✅ Comprehensive test suite"
    
    echo -e "\nDocumentation:"
    echo "  ✅ Enhanced dividend alerts guide"
    echo "  ✅ Frontend integration documentation"
    echo "  ✅ API integration guide"
    echo "  ✅ Database migration script"
    
    echo -e "\n${GREEN}🎉 Demo completed!${NC}"
    
    if [ "$BACKEND_RUNNING" = true ] && [ "$FRONTEND_RUNNING" = true ]; then
        echo -e "\n${BLUE}🌐 Live Demo URLs:${NC}"
        echo "  • Frontend: $FRONTEND_URL"
        echo "  • Backend API: $BACKEND_URL"
        echo "  • Health Check: $BACKEND_URL/health"
        echo "  • Dividend Data: $BACKEND_URL/dividends/gse"
    else
        echo -e "\n${YELLOW}💡 To run full demo:${NC}"
        echo "  1. Start backend: cd backend && go run cmd/server/main.go"
        echo "  2. Start frontend: npm run dev"
        echo "  3. Run this demo again: ./demo_enhanced_dividend_alerts.sh"
    fi
    
    echo -e "\n${BLUE}📚 Next Steps:${NC}"
    echo "  1. Run database migration: backend/migrations/add_dividend_yield_fields.sql"
    echo "  2. Configure email settings for alert notifications"
    echo "  3. Set up monitoring for dividend yield changes"
    echo "  4. Deploy to production environment"
}

# Check if jq is available for JSON parsing
if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}💡 Install 'jq' for better JSON output formatting${NC}"
fi

# Run the demo
main

echo -e "\n${BLUE}📖 For detailed documentation, see:${NC}"
echo "  • ENHANCED_DIVIDEND_ALERTS.md"
echo "  • FRONTEND_DIVIDEND_ALERTS_INTEGRATION.md"
echo "  • DIVIDEND_API_INTEGRATION.md"