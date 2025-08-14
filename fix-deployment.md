# Fix Deployment Instructions

## Current Issue
Your frontend is deployed at `https://stock-alert-gh.onrender.com/` but it's trying to call a backend at `https://shares-alert-backend.onrender.com/` which doesn't exist.

## Solution Options

### Option 1: Deploy Both Services Separately (Recommended)
1. **Deploy the backend first:**
   - Go to Render Dashboard
   - Create a new Web Service
   - Connect your GitHub repo
   - Set the root directory to `backend`
   - Use these settings:
     - Name: `shares-alert-backend`
     - Environment: Go
     - Build Command: `go build -o main .`
     - Start Command: `./main`
     - Port: 8080

2. **Update frontend to use the backend URL:**
   - Your frontend will automatically use the correct backend URL after the backend is deployed

### Option 2: Single Service Deployment
If you prefer to run both frontend and backend on the same service, you'll need to:
1. Modify the nginx.conf to proxy API calls to the backend
2. Use a different deployment approach

## Files Updated
- ✅ `.env.production` - Updated API URL
- ✅ `render.yaml` - Fixed service names
- ✅ `backend/main.go` - Fixed CORS configuration

## Next Steps
1. Commit these changes:
   ```bash
   git add .
   git commit -m "Fix CORS and deployment configuration"
   git push origin main
   ```

2. Deploy the backend service on Render using the settings above

3. Once both services are deployed, your app should work correctly!

## Testing
After deployment, test these URLs:
- Frontend: https://stock-alert-gh.onrender.com/
- Backend Health: https://shares-alert-backend.onrender.com/api/v1/health
- Backend API: https://shares-alert-backend.onrender.com/api/v1/stocks