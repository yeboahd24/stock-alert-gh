import { cache } from './cache';

// API service for communicating with the Go backend
const getApiBaseUrl = (): string => {
  // Check for environment variable first
  if (import.meta.env?.VITE_API_URL) {
    return import.meta.env.VITE_API_URL;
  }
  
  // Fallback for production deployment
  if (import.meta.env?.PROD) {
    return 'https://stock-alert-gh-backend.onrender.com/api/v1';
  }
  
  // Development fallback
  return 'http://localhost:10000/api/v1';
};

const API_BASE_URL = getApiBaseUrl();

// Debug: Log the API URL being used
console.log('API Base URL:', API_BASE_URL);
console.log('VITE_API_URL env var:', import.meta.env?.VITE_API_URL);
console.log('Production mode:', import.meta.env?.PROD);

// Auth token management
let authToken: string | null = null;

export const setAuthToken = (token: string | null) => {
  authToken = token;
};

export const getAuthToken = () => authToken;

// Helper function to make authenticated requests
const makeAuthenticatedRequest = async (url: string, options: RequestInit = {}) => {
  const token = authToken || localStorage.getItem('auth_token');
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (response.status === 401) {
    // Token expired or invalid
    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_user');
    window.location.href = '/login';
    throw new Error('Authentication required');
  }

  return response;
};

// Types matching the backend API
export interface Stock {
  symbol: string;
  name: string;
  currentPrice: number;
  previousClose: number;
  change: number;
  changePercent: number;
  volume: number;
  lastUpdated: string;
  marketCap?: number;
  sector?: string;
  industry?: string;
}

export interface Company {
  address: string;
  directors: string[];
  email: string;
  facsimile?: string;
  industry: string;
  name: string;
  sector: string;
  telephone: string;
  website: string;
}

export interface DetailedStock {
  symbol: string;
  name: string;
  currentPrice: number;
  previousClose: number;
  change: number;
  changePercent: number;
  volume: number;
  lastUpdated: string;
  marketCap: number;
  shares: number;
  sector: string;
  industry: string;
  dps?: number;
  eps?: number;
  company: Company;
}

export interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
  googleId: string;
  emailVerified: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UserPreferences {
  id: string;
  userId: string;
  emailNotifications: boolean;
  pushNotifications: boolean;
  notificationFrequency: string;
  createdAt: string;
  updatedAt: string;
}

export interface Alert {
  id: string;
  userId: string;
  stockSymbol: string;
  stockName: string;
  alertType: string;
  thresholdPrice?: number;
  currentPrice?: number;
  status: string;
  createdAt: string;
  updatedAt: string;
  triggeredAt?: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface CreateAlertRequest {
  stockSymbol: string;
  stockName: string;
  alertType: string;
  thresholdPrice?: number;
}

// Authentication API functions
export const authApi = {
  // Get Google OAuth URL
  getGoogleAuthUrl: async (state?: string): Promise<{ authUrl: string }> => {
    const params = state ? `?state=${encodeURIComponent(state)}` : '';
    const response = await fetch(`${API_BASE_URL}/auth/google${params}`);
    if (!response.ok) {
      throw new Error('Failed to get Google auth URL');
    }
    return response.json();
  },

  // Handle Google OAuth callback
  googleCallback: async (code: string): Promise<AuthResponse> => {
    const response = await fetch(`${API_BASE_URL}/auth/google/callback`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ code }),
    });
    if (!response.ok) {
      const error = await response.text();
      throw new Error(`Authentication failed: ${error}`);
    }
    const data = await response.json();
    setAuthToken(data.token);
    return data;
  },

  // Get user profile
  getProfile: async (): Promise<User> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/auth/profile`);
    if (!response.ok) {
      throw new Error('Failed to get user profile');
    }
    return response.json();
  },

  // Logout
  logout: async (): Promise<void> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
    });
    if (!response.ok) {
      throw new Error('Failed to logout');
    }
    setAuthToken(null);
  },
};

// User API functions
export const userApi = {
  // Get user preferences
  getPreferences: async (): Promise<UserPreferences> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/user/preferences`);
    if (!response.ok) {
      throw new Error('Failed to get user preferences');
    }
    return response.json();
  },

  // Update user preferences
  updatePreferences: async (preferences: Partial<UserPreferences>): Promise<UserPreferences> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/user/preferences`, {
      method: 'PUT',
      body: JSON.stringify(preferences),
    });
    if (!response.ok) {
      throw new Error('Failed to update user preferences');
    }
    return response.json();
  },
};

// Stock API functions
export const stockApi = {
  // Get all stocks
  getAllStocks: async (): Promise<Stock[]> => {
    const cacheKey = 'stocks:all';
    const cached = cache.get<Stock[]>(cacheKey);
    if (cached) return cached;

    console.log('Fetching stocks from:', `${API_BASE_URL}/stocks`);
    try {
      const response = await fetch(`${API_BASE_URL}/stocks`);
      console.log('Response status:', response.status);
      if (!response.ok) {
        throw new Error(`Failed to fetch stocks: ${response.status} ${response.statusText}`);
      }
      const data = await response.json();
      console.log('Stocks data received:', data.length, 'stocks');
      cache.set(cacheKey, data, 2); // Cache for 2 minutes
      return data;
    } catch (error) {
      console.error('Error in getAllStocks:', error);
      throw error;
    }
  },

  // Get specific stock
  getStock: async (symbol: string): Promise<Stock> => {
    const cacheKey = `stock:${symbol}`;
    const cached = cache.get<Stock>(cacheKey);
    if (cached) return cached;

    const response = await fetch(`${API_BASE_URL}/stocks/${symbol}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch stock ${symbol}`);
    }
    const data = await response.json();
    cache.set(cacheKey, data, 1); // Cache for 1 minute
    return data;
  },

  // Get stock details
  getStockDetails: async (symbol: string): Promise<DetailedStock> => {
    const cacheKey = `stock:details:${symbol}`;
    const cached = cache.get<DetailedStock>(cacheKey);
    if (cached) return cached;

    const response = await fetch(`${API_BASE_URL}/stocks/${symbol}/details`);
    if (!response.ok) {
      throw new Error(`Failed to fetch stock details for ${symbol}`);
    }
    const data = await response.json();
    cache.set(cacheKey, data, 5); // Cache for 5 minutes
    return data;
  },
};

// Alert API functions
export const alertApi = {
  // Get all alerts (authenticated)
  getAllAlerts: async (status?: string, stockSymbol?: string): Promise<Alert[]> => {
    const cacheKey = `alerts:${status || 'all'}:${stockSymbol || 'all'}`;
    const cached = cache.get<Alert[]>(cacheKey);
    if (cached) return cached;

    const params = new URLSearchParams();
    if (status) params.append('status', status);
    if (stockSymbol) params.append('stockSymbol', stockSymbol);
    
    const url = `${API_BASE_URL}/alerts${params.toString() ? `?${params.toString()}` : ''}`;
    console.log('Fetching alerts from:', url);
    try {
      const response = await makeAuthenticatedRequest(url);
      console.log('Alerts response status:', response.status);
      if (!response.ok) {
        if (response.status === 401) {
          console.warn('Authentication failed for alerts - user may need to re-login');
          return []; // Return empty array instead of throwing error
        }
        throw new Error(`Failed to fetch alerts: ${response.status} ${response.statusText}`);
      }
      const data = await response.json();
      console.log('Alerts data received:', Array.isArray(data) ? data.length : 'not an array', 'alerts');
      const alerts = Array.isArray(data) ? data : [];
      cache.set(cacheKey, alerts, 1); // Cache for 1 minute
      return alerts;
    } catch (error) {
      console.error('Error in getAllAlerts:', error);
      // Return empty array instead of throwing to prevent app crash
      return [];
    }
  },

  // Create alert (authenticated)
  createAlert: async (alertData: CreateAlertRequest): Promise<Alert> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/alerts`, {
      method: 'POST',
      body: JSON.stringify(alertData),
    });
    if (!response.ok) {
      throw new Error('Failed to create alert');
    }
    const data = await response.json();
    // Clear alerts cache when creating new alert
    cache.delete('alerts:all:all');
    return data;
  },

  // Get specific alert (authenticated)
  getAlert: async (id: string): Promise<Alert> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/alerts/${id}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch alert ${id}`);
    }
    return response.json();
  },

  // Update alert (authenticated)
  updateAlert: async (id: string, updates: Partial<CreateAlertRequest & { status: string }>): Promise<Alert> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/alerts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updates),
    });
    if (!response.ok) {
      throw new Error(`Failed to update alert ${id}`);
    }
    const data = await response.json();
    // Clear alerts cache when updating
    cache.delete('alerts:all:all');
    return data;
  },

  // Delete alert (authenticated)
  deleteAlert: async (id: string): Promise<void> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/alerts/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error(`Failed to delete alert ${id}`);
    }
    // Clear alerts cache when deleting
    cache.delete('alerts:all:all');
  },
};

// GSE Dividend API functions
export const dividendApi = {
  // Get all GSE dividend stocks
  getGSEDividendStocks: async (): Promise<any> => {
    const response = await fetch(`${API_BASE_URL}/dividends/gse`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch GSE dividend stocks');
    }
    
    return response.json();
  },

  // Get dividend data for specific stock
  getDividendStockBySymbol: async (symbol: string): Promise<any> => {
    const response = await fetch(`${API_BASE_URL}/dividends/gse/${symbol}`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch dividend data for ${symbol}`);
    }
    
    return response.json();
  },

  // Get high dividend yield stocks
  getHighDividendYieldStocks: async (minYield?: number): Promise<any> => {
    const params = minYield ? `?minYield=${minYield}` : '';
    const response = await fetch(`${API_BASE_URL}/dividends/high-yield${params}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch high dividend yield stocks');
    }
    
    return response.json();
  },

  // Get traditional dividend announcements
  getAllDividends: async (): Promise<any> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/dividends`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch dividend announcements');
    }
    
    return response.json();
  },

  // Get upcoming dividend payments
  getUpcomingDividends: async (): Promise<any> => {
    const response = await makeAuthenticatedRequest(`${API_BASE_URL}/dividends/upcoming`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch upcoming dividends');
    }
    
    return response.json();
  },
};

// Health check
export const healthCheck = async (): Promise<{ status: string; timestamp: string }> => {
  const response = await fetch(`${API_BASE_URL}/health`);
  if (!response.ok) {
    throw new Error('Health check failed');
  }
  return response.json();
};