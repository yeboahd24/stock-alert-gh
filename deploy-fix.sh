#!/bin/bash

echo "🚀 Deploying Stock Alert Ghana fixes..."

# Check if we have uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo "📝 Committing changes..."
    git add .
    git commit -m "Fix: Use correct Ghana Stock Exchange API endpoint"
fi

# Push to trigger deployment
echo "📤 Pushing to GitHub to trigger deployment..."
git push origin main

if [ $? -eq 0 ]; then
    echo "✅ Successfully pushed to GitHub!"
    echo "🔄 Render will automatically redeploy the application."
    echo "⏱️  This usually takes 2-3 minutes."
    echo ""
    echo "🔗 Check deployment status at:"
    echo "   Backend: https://stock-alert-gh-backend.onrender.com/api/v1/health"
    echo "   Frontend: https://stock-alert-gh.onrender.com"
else
    echo "❌ Failed to push to GitHub. Please check your authentication."
    echo "💡 You may need to:"
    echo "   1. Set up GitHub authentication (token or SSH key)"
    echo "   2. Or manually push the changes through GitHub web interface"
fi