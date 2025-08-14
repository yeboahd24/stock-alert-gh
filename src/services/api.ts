// API service for communicating with the Go backend
const API_BASE_URL = (import.meta as any).env.VITE_API_URL || 'http://localhost:8080/api/v1';

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
}

export interface CreateAlertRequest {
  stockSymbol: string;
  stockName: string;
  alertType: string;
  thresholdPrice?: number;
}

// Stock API functions
export const stockApi = {
  // Get all stocks
  getAllStocks: async (): Promise<Stock[]> => {
    console.log('Fetching stocks from:', `${API_BASE_URL}/stocks`);
    try {
      const response = await fetch(`${API_BASE_URL}/stocks`);
      console.log('Response status:', response.status);
      if (!response.ok) {
        throw new Error(`Failed to fetch stocks: ${response.status} ${response.statusText}`);
      }
      const data = await response.json();
      console.log('Stocks data received:', data.length, 'stocks');
      return data;
    } catch (error) {
      console.error('Error in getAllStocks:', error);
      throw error;
    }
  },

  // Get specific stock
  getStock: async (symbol: string): Promise<Stock> => {
    const response = await fetch(`${API_BASE_URL}/stocks/${symbol}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch stock ${symbol}`);
    }
    return response.json();
  },

  // Get stock details
  getStockDetails: async (symbol: string): Promise<DetailedStock> => {
    const response = await fetch(`${API_BASE_URL}/stocks/${symbol}/details`);
    if (!response.ok) {
      throw new Error(`Failed to fetch stock details for ${symbol}`);
    }
    return response.json();
  },
};

// Alert API functions
export const alertApi = {
  // Get all alerts
  getAllAlerts: async (userId?: string, status?: string): Promise<Alert[]> => {
    const params = new URLSearchParams();
    if (userId) params.append('userId', userId);
    if (status) params.append('status', status);
    
    const url = `${API_BASE_URL}/alerts${params.toString() ? `?${params.toString()}` : ''}`;
    console.log('Fetching alerts from:', url);
    try {
      const response = await fetch(url);
      console.log('Alerts response status:', response.status);
      if (!response.ok) {
        throw new Error(`Failed to fetch alerts: ${response.status} ${response.statusText}`);
      }
      const data = await response.json();
      console.log('Alerts data received:', data.length, 'alerts');
      return data;
    } catch (error) {
      console.error('Error in getAllAlerts:', error);
      throw error;
    }
  },

  // Create alert
  createAlert: async (alertData: CreateAlertRequest): Promise<Alert> => {
    const response = await fetch(`${API_BASE_URL}/alerts`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(alertData),
    });
    if (!response.ok) {
      throw new Error('Failed to create alert');
    }
    return response.json();
  },

  // Get specific alert
  getAlert: async (id: string): Promise<Alert> => {
    const response = await fetch(`${API_BASE_URL}/alerts/${id}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch alert ${id}`);
    }
    return response.json();
  },

  // Update alert
  updateAlert: async (id: string, updates: Partial<CreateAlertRequest & { status: string }>): Promise<Alert> => {
    const response = await fetch(`${API_BASE_URL}/alerts/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updates),
    });
    if (!response.ok) {
      throw new Error(`Failed to update alert ${id}`);
    }
    return response.json();
  },

  // Delete alert
  deleteAlert: async (id: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/alerts/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error(`Failed to delete alert ${id}`);
    }
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