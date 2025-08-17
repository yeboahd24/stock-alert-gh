import React, { useEffect, useState } from 'react';
import {
  Container,
  Paper,
  Typography,
  Button,
  Box,
  Stack,
  Alert,
  CircularProgress,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import GoogleIcon from '@mui/icons-material/Google';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import { useAuth } from '../../src/contexts/AuthContext';
import { authApi } from '../../src/services/api';

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(4),
  maxWidth: 400,
  margin: '0 auto',
  marginTop: theme.spacing(8),
  textAlign: 'center',
}));

const LoginPage: React.FC = () => {
  const { login, isLoading } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isAuthenticating, setIsAuthenticating] = useState(false);

  // Handle OAuth callback
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    const error = urlParams.get('error');

    if (error) {
      setError('Authentication was cancelled or failed');
      // Clean up URL
      window.history.replaceState({}, document.title, window.location.pathname);
      return;
    }

    if (code) {
      setIsAuthenticating(true);
      login(code)
        .then(() => {
          // Clean up URL
          window.history.replaceState({}, document.title, window.location.pathname);
        })
        .catch((err) => {
          console.error('Login failed:', err);
          setError('Login failed. Please try again.');
          // Clean up URL
          window.history.replaceState({}, document.title, window.location.pathname);
        })
        .finally(() => {
          setIsAuthenticating(false);
        });
    }
  }, [login]);

  const handleGoogleLogin = async () => {
    try {
      setError(null);
      setIsAuthenticating(true);
      
      // Generate a random state for security
      const state = Math.random().toString(36).substring(2, 15);
      
      // Get Google OAuth URL from backend
      const { authUrl } = await authApi.getGoogleAuthUrl(state);
      
      // Redirect to Google OAuth
      window.location.href = authUrl;
    } catch (err) {
      console.error('Failed to initiate Google login:', err);
      setError('Failed to start authentication. Please try again.');
      setIsAuthenticating(false);
    }
  };

  if (isLoading || isAuthenticating) {
    return (
      <Container maxWidth="sm">
        <StyledPaper>
          <Stack spacing={3} alignItems="center">
            <CircularProgress />
            <Typography variant="h6">
              {isAuthenticating ? 'Authenticating...' : 'Loading...'}
            </Typography>
          </Stack>
        </StyledPaper>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm">
      <StyledPaper elevation={3}>
        <Stack spacing={3} alignItems="center">
          {/* Logo and Title */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <TrendingUpIcon sx={{ fontSize: 40, color: 'primary.main' }} />
            <Typography variant="h4" component="h1" fontWeight="bold">
              Shares Alert
            </Typography>
          </Box>
          
          <Typography variant="h6" color="text.secondary">
            Ghana Stock Exchange
          </Typography>

          <Typography variant="body1" color="text.secondary" sx={{ textAlign: 'center' }}>
            Track your favorite stocks and get notified when they hit your target prices.
          </Typography>

          {error && (
            <Alert severity="error" sx={{ width: '100%' }}>
              {error}
            </Alert>
          )}

          {/* Google Sign In Button */}
          <Button
            variant="contained"
            size="large"
            startIcon={<GoogleIcon />}
            onClick={handleGoogleLogin}
            disabled={isAuthenticating}
            sx={{
              width: '100%',
              py: 1.5,
              backgroundColor: '#4285f4',
              '&:hover': {
                backgroundColor: '#3367d6',
              },
            }}
          >
            {isAuthenticating ? 'Signing in...' : 'Sign in with Google'}
          </Button>

          {/* Features */}
          <Box sx={{ mt: 4, textAlign: 'left', width: '100%' }}>
            <Typography variant="subtitle2" color="text.secondary" gutterBottom>
              Features:
            </Typography>
            <Stack spacing={1}>
              <Typography variant="body2" color="text.secondary">
                • Real-time stock price tracking
              </Typography>
              <Typography variant="body2" color="text.secondary">
                • Custom price alerts
              </Typography>
              <Typography variant="body2" color="text.secondary">
                • Email notifications
              </Typography>
              <Typography variant="body2" color="text.secondary">
                • Portfolio monitoring
              </Typography>
            </Stack>
          </Box>
        </Stack>
      </StyledPaper>
    </Container>
  );
};

export default LoginPage;