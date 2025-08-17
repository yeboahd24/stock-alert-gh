import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    open: true
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // Vendor chunks
          'react-vendor': ['react', 'react-dom'],
          'mui-core': ['@mui/material', '@emotion/react', '@emotion/styled'],
          'mui-icons': ['@mui/icons-material'],
          'mui-charts': ['@mui/x-charts'],
          'charts': ['recharts'],
          
          // App chunks
          'auth': [
            './src/contexts/AuthContext.tsx',
            './components/auth/LoginPage.tsx',
            './components/auth/ProtectedRoute.tsx'
          ],
          'dashboard': [
            './components/dashboard/Dashboard.tsx',
            './components/charts/StockPriceChart.tsx'
          ],
          'forms': [
            './components/forms/AlertForm.tsx',
            './components/forms/NotificationSettings.tsx'
          ],
          'tables': ['./components/tables/AlertsTable.tsx'],
          'common': [
            './components/common/StockCard.tsx',
            './components/common/StockDetailsModal.tsx',
            './components/common/UserMenu.tsx',
            './components/common/AlertStatusChip.tsx'
          ]
        }
      }
    },
    chunkSizeWarningLimit: 1000, // Increase warning limit to 1MB
    sourcemap: false // Disable sourcemaps in production for smaller builds
  }
})