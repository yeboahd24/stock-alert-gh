# Deployment Guide - Shares Alert Ghana

## 🚀 Ready for Render Deployment!

Your application is now fully dockerized and ready for deployment on Render. Here's everything you need to know:

## 📁 Project Structure
```
shares-alert-ghana/
├── backend/
│   ├── Dockerfile              # Backend container
│   ├── main.go                 # Complete Go application
│   ├── go.mod                  # Go dependencies
│   └── .dockerignore           # Docker ignore rules
├── components/                 # React components
├── src/                        # Frontend source
├── Dockerfile                  # Frontend container
├── docker-compose.yml          # Local development
├── render.yaml                 # Render deployment config
├── nginx.conf                  # Frontend web server
├── deploy.sh                   # Local deployment script
├── test-deployment.sh          # Test builds
└── README.md                   # Documentation
```

## 🌐 Render Deployment Options

### Option 1: Blueprint Deployment (Recommended)
1. **Push to GitHub**
   ```bash
   git add .
   git commit -m "Ready for Render deployment"
   git push origin main
   ```

2. **Deploy on Render**
   - Go to [Render Dashboard](https://dashboard.render.com)
   - Click "New" → "Blueprint"
   - Connect your GitHub repository
   - Render will automatically detect `render.yaml`
   - Click "Apply" to deploy both services

### Option 2: Manual Service Creation
1. **Backend Service**
   - Type: Web Service
   - Environment: Go
   - Build Command: `go build -o main .`
   - Start Command: `./main`
   - Port: 8080

2. **Frontend Service**
   - Type: Static Site
   - Build Command: `npm ci && npm run build`
   - Publish Directory: `dist`
   - Environment Variable: `VITE_API_URL=https://your-backend-url.onrender.com/api/v1`

## 🔧 Environment Variables

### Backend (Auto-configured)
- `PORT`: Automatically set by Render
- `GO_ENV`: Set to "production"

### Frontend
- `VITE_API_URL`: Points to your backend service URL

## 🧪 Testing Before Deployment

Run the test script to ensure everything builds correctly:
```bash
./test-deployment.sh
```

## 🐳 Local Docker Testing

Test the full stack locally:
```bash
# Build and run all services
docker-compose up --build

# Access the application
# Frontend: http://localhost
# Backend: http://localhost:8080/api/v1
```

## 📊 API Endpoints

Your deployed backend will provide:

### Stock Data
- `GET /api/v1/stocks` - All stocks
- `GET /api/v1/stocks/{symbol}` - Specific stock
- `GET /api/v1/stocks/{symbol}/details` - Detailed stock info

### Alerts
- `GET /api/v1/alerts` - All alerts
- `POST /api/v1/alerts` - Create alert
- `PUT /api/v1/alerts/{id}` - Update alert
- `DELETE /api/v1/alerts/{id}` - Delete alert

### Health Check
- `GET /api/v1/health` - Service status

## 🔍 Monitoring

After deployment, monitor your services:
- **Logs**: Available in Render dashboard
- **Health**: Check `/api/v1/health` endpoint
- **Metrics**: View in Render service dashboard

## 🚨 Troubleshooting

### Common Issues:
1. **Build Failures**: Check logs in Render dashboard
2. **CORS Errors**: Backend is configured for Render domains
3. **API Connection**: Verify `VITE_API_URL` environment variable

### Debug Commands:
```bash
# Test backend locally
cd backend && go run .

# Test frontend locally
npm run dev

# Check Docker builds
docker-compose build
```

## 🎯 Expected Deployment URLs

After successful deployment:
- **Frontend**: `https://shares-alert-frontend.onrender.com`
- **Backend**: `https://shares-alert-backend.onrender.com`
- **API Health**: `https://shares-alert-backend.onrender.com/api/v1/health`

## 📈 Features Deployed

✅ **Real-time Ghana Stock Exchange data**
✅ **Interactive stock dashboard**
✅ **Stock price alerts system**
✅ **Detailed company information**
✅ **Responsive Material-UI design**
✅ **RESTful API backend**
✅ **Docker containerization**
✅ **Production-ready configuration**

## 🔄 Updates and Maintenance

To update your deployment:
1. Make changes to your code
2. Push to GitHub
3. Render will automatically redeploy

## 💡 Next Steps

After deployment, consider:
- Setting up custom domain
- Adding database for persistent alerts
- Implementing user authentication
- Adding email/SMS notifications
- Setting up monitoring and analytics

---

🎉 **Your Ghana Stock Exchange application is ready for the world!**