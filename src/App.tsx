import React, { Suspense } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Box } from '@mui/material';
import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from '../components/auth/ProtectedRoute';
import BrandedLoader from '../components/common/BrandedLoader';
import theme from '../theme/theme';
import './styles/brand.css';

// Lazy load the Dashboard component
const Dashboard = React.lazy(() => import('../components/dashboard/Dashboard'));

// Loading component
const LoadingFallback = () => (
  <Box 
    display="flex" 
    justifyContent="center" 
    alignItems="center" 
    minHeight="100vh"
    sx={{ backgroundColor: 'background.default' }}
  >
    <BrandedLoader 
      message="Initializing Shares Alert Ghana..."
      size="large"
    />
  </Box>
);

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <ProtectedRoute>
          <Suspense fallback={<LoadingFallback />}>
            <Dashboard />
          </Suspense>
        </ProtectedRoute>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
