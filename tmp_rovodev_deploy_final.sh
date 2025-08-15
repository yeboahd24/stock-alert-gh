#!/bin/bash

echo "🚀 Final Deployment Fix for Stock Alert Ghana"
echo "=============================================="

# Test builds first
echo "🔧 Testing builds..."

# Test backend build
echo "📦 Testing backend build..."
cd backend
if go build -o main .; then
    echo "✅ Backend build successful"
    rm -f main
else
    echo "❌ Backend build failed"
    exit 1
fi
cd ..

# Test frontend build
echo "📦 Testing frontend build..."
if npm run build; then
    echo "✅ Frontend build successful"
else
    echo "❌ Frontend build failed"
    exit 1
fi

echo ""
echo "🎯 Deployment Configuration Summary:"
echo "====================================="
echo "Backend Service Name: stock-alert-gh-backend"
echo "Frontend Service Name: stock-alert-gh"
echo "API URL: https://stock-alert-gh-backend.onrender.com/api/v1"
echo ""

# Check git status
if [[ -n $(git status -s) ]]; then
    echo "📝 Committing deployment fixes..."
    git add .
    git commit -m "Fix: Update deployment configuration for correct service names"
    
    echo "📤 Pushing to trigger deployment..."
    if git push origin main; then
        echo "✅ Successfully pushed to GitHub!"
        echo ""
        echo "🔄 Render will automatically deploy both services:"
        echo "   1. Backend: stock-alert-gh-backend"
        echo "   2. Frontend: stock-alert-gh"
        echo ""
        echo "⏱️  Deployment usually takes 3-5 minutes."
        echo ""
        echo "🔗 After deployment, check these URLs:"
        echo "   Frontend: https://stock-alert-gh.onrender.com"
        echo "   Backend Health: https://stock-alert-gh-backend.onrender.com/api/v1/health"
        echo "   Backend API: https://stock-alert-gh-backend.onrender.com/api/v1/stocks"
    else
        echo "❌ Failed to push to GitHub"
        exit 1
    fi
else
    echo "✅ No changes to commit - configuration is up to date"
    echo ""
    echo "🔗 Your deployment URLs:"
    echo "   Frontend: https://stock-alert-gh.onrender.com"
    echo "   Backend Health: https://stock-alert-gh-backend.onrender.com/api/v1/health"
fi

echo ""
echo "🎉 Deployment configuration is ready!"