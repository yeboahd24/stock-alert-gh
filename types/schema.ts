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
    alertType: 'price_threshold' | 'ipo_alert' | 'dividend_announcement' | 'high_dividend_yield' | 'dividend_yield_change' | 'target_dividend_yield';
    thresholdPrice?: number;
    currentPrice?: number;
    thresholdYield?: number;
    currentYield?: number;
    targetYield?: number;
    yieldChangeThreshold?: number;
    lastYield?: number;
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

// GSE Dividend API types
export interface GSEDividendStock {
  symbol: string;
  name: string;
  dividend_yield: number;
  price: string;
  market_cap: string;
  country: string;
  exchange: string;
  sector: string;
  url: string;
}

export interface GSEDividendData {
  timestamp: string;
  source: string;
  count: number;
  stocks: GSEDividendStock[];
}

export interface GSEDividendResponse {
  success: boolean;
  data: GSEDividendData;
}

// Enhanced Alert types
export interface Alert {
  id: string;
  userId: string;
  stockSymbol: string;
  stockName: string;
  alertType: string;
  thresholdPrice?: number;
  currentPrice?: number;
  thresholdYield?: number;
  currentYield?: number;
  targetYield?: number;
  yieldChangeThreshold?: number;
  lastYield?: number;
  status: string;
  createdAt: string;
  updatedAt: string;
  triggeredAt?: string;
}

export interface CreateAlertRequest {
  stockSymbol: string;
  stockName: string;
  alertType: string;
  thresholdPrice?: number;
  thresholdYield?: number;
  targetYield?: number;
  yieldChangeThreshold?: number;
}