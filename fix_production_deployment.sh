#!/bin/bash

echo "🚨 PRODUCTION FIX: Adding missing dividend yield columns"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if we're in the right directory
if [ ! -f "backend/go.mod" ]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

print_warning "This script will fix the production database schema issue"
print_warning "Error: pq: column \"threshold_yield\" does not exist"
echo ""

# Step 1: Build the migration tool
print_status "Building migration tool..."
cd backend
if go build -o migrate_tool ./cmd/migrate; then
    print_status "Migration tool built successfully"
else
    print_error "Failed to build migration tool"
    exit 1
fi

# Step 2: Run the migration
print_status "Running dividend yield migration..."
if ./migrate_tool; then
    print_status "Migration completed successfully"
else
    print_error "Migration failed"
    print_warning "You may need to run the SQL manually:"
    echo "psql -d your_database -f ../fix_production_migration.sql"
    exit 1
fi

# Clean up
rm -f migrate_tool
cd ..

# Step 3: Test the backend build
print_status "Testing backend build..."
cd backend
if go build -o test_build .; then
    print_status "Backend builds successfully"
    rm -f test_build
else
    print_error "Backend build failed"
    exit 1
fi
cd ..

# Step 4: Commit and deploy the fix
if [[ -n $(git status -s) ]]; then
    print_status "Committing production fix..."
    git add .
    git commit -m "HOTFIX: Add missing dividend yield columns to production database

- Fixes error: pq: column \"threshold_yield\" does not exist
- Adds migration script for dividend yield features
- Ensures backward compatibility with IF NOT EXISTS clauses"
    
    print_status "Pushing to trigger deployment..."
    if git push origin main; then
        print_status "Successfully pushed to GitHub!"
        echo ""
        print_status "🔄 Render will automatically deploy the fix"
        print_warning "⏱️  Deployment usually takes 3-5 minutes"
        echo ""
        echo "🔗 Monitor deployment at:"
        echo "   https://dashboard.render.com"
        echo ""
        echo "🧪 After deployment, test these endpoints:"
        echo "   Health: https://stock-alert-gh-backend.onrender.com/api/v1/health"
        echo "   Alerts: https://stock-alert-gh-backend.onrender.com/api/v1/alerts"
    else
        print_error "Failed to push to GitHub"
        exit 1
    fi
else
    print_warning "No changes to commit"
fi

echo ""
print_status "🎉 Production fix deployment initiated!"
echo ""
echo "📋 What was fixed:"
echo "   ✅ Added missing dividend yield columns"
echo "   ✅ Applied database migration"
echo "   ✅ Added proper indexes for performance"
echo "   ✅ Ensured backward compatibility"
echo ""
echo "🔍 Next steps:"
echo "   1. Monitor Render deployment logs"
echo "   2. Test the health endpoint"
echo "   3. Verify alerts functionality"
echo "   4. Check application logs for any remaining issues"