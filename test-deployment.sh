#!/bin/bash

# Test deployment script

echo "🧪 Testing local deployment..."

# Test backend build
echo "📦 Testing backend build..."
cd backend
if go build -o main .; then
    echo "✅ Backend builds successfully"
    rm -f main
else
    echo "❌ Backend build failed"
    exit 1
fi
cd ..

# Test frontend build
echo "📦 Testing frontend build..."
if npm run build; then
    echo "✅ Frontend builds successfully"
    rm -rf dist
else
    echo "❌ Frontend build failed"
    exit 1
fi

# Test Docker builds
echo "🐳 Testing Docker builds..."

# Test backend Docker build
echo "🔨 Building backend Docker image..."
if docker build -t shares-alert-backend-test ./backend; then
    echo "✅ Backend Docker image builds successfully"
    docker rmi shares-alert-backend-test
else
    echo "❌ Backend Docker build failed"
    exit 1
fi

# Test frontend Docker build
echo "🔨 Building frontend Docker image..."
if docker build -t shares-alert-frontend-test .; then
    echo "✅ Frontend Docker image builds successfully"
    docker rmi shares-alert-frontend-test
else
    echo "❌ Frontend Docker build failed"
    exit 1
fi

echo "🎉 All tests passed! Ready for deployment."
echo ""
echo "📋 Next steps:"
echo "1. Push your code to GitHub"
echo "2. Connect your repository to Render"
echo "3. Deploy using the render.yaml configuration"
echo ""
echo "🔗 Render deployment URL: https://dashboard.render.com"