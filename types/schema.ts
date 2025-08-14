// Type definitions for Shares Alert Ghana application

// Props types (data passed to components)
export interface PropTypes {
  initialTab: 'dashboard' | 'alerts' | 'notifications' | 'profile';
  showWelcomeModal: boolean;
}

// Store types (global state data)
export interface StoreTypes {
  user: {
    id: string;
    name: string;
    email: string;
    phone: string;
    language: 'english' | 'twi';
    subscriptionTier: 'free' | 'premium';
  };
  alerts: Array<{
    id: string;
    stockSymbol: string;
    stockName: string;
    alertType: 'price_threshold' | 'ipo_alert' | 'dividend_announcement';
    thresholdPrice?: number;
    currentPrice?: number;
    status: 'active' | 'inactive' | 'triggered';
    createdAt: string;
  }>;
}

// Query types (API response data)
export interface QueryTypes {
  stocks: Array<{
    symbol: string;
    name: string;
    currentPrice: number;
    previousClose: number;
    change: number;
    changePercent: number;
    volume: number;
    lastUpdated: string;
  }>;
  notifications: Array<{
    id: string;
    title: string;
    message: string;
    timestamp: string;
    read: boolean;
  }>;
}