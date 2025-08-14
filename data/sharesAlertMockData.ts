// Mock data for Shares Alert Ghana application

// Data for global state store
export const mockStore = {
  user: {
    id: "user-123",
    name: "Kwame Asante",
    email: "kwame.asante@gmail.com",
    phone: "+233244123456",
    language: "english" as const,
    subscriptionTier: "premium" as const
  },
  alerts: [
    {
      id: "alert-1",
      stockSymbol: "MTN",
      stockName: "MTN Ghana",
      alertType: "price_threshold" as const,
      thresholdPrice: 0.85,
      currentPrice: 0.82,
      status: "active" as const,
      createdAt: "2024-01-15T10:30:00Z"
    },
    {
      id: "alert-2",
      stockSymbol: "GCB",
      stockName: "GCB Bank Limited",
      alertType: "dividend_announcement" as const,
      status: "active" as const,
      createdAt: "2024-01-10T14:20:00Z"
    }
  ]
};

// Data returned by API queries
export const mockQuery = {
  stocks: [
    {
      symbol: "MTN",
      name: "MTN Ghana",
      currentPrice: 0.82,
      previousClose: 0.80,
      change: 0.02,
      changePercent: 2.5,
      volume: 150000,
      lastUpdated: "2024-01-20T16:00:00Z"
    },
    {
      symbol: "GCB",
      name: "GCB Bank Limited",
      currentPrice: 4.25,
      previousClose: 4.30,
      change: -0.05,
      changePercent: -1.16,
      volume: 75000,
      lastUpdated: "2024-01-20T16:00:00Z"
    },
    {
      symbol: "TOTAL",
      name: "Total Petroleum Ghana",
      currentPrice: 2.15,
      previousClose: 2.10,
      change: 0.05,
      changePercent: 2.38,
      volume: 45000,
      lastUpdated: "2024-01-20T16:00:00Z"
    }
  ],
  notifications: [
    {
      id: "notif-1",
      title: "Price Alert Triggered",
      message: "MTN Ghana has reached your target price of GHS 0.85",
      timestamp: "2024-01-20T15:45:00Z",
      read: false
    },
    {
      id: "notif-2",
      title: "IPO Announcement",
      message: "New IPO listing available for subscription",
      timestamp: "2024-01-19T09:30:00Z",
      read: true
    }
  ]
};

// Data passed as props to the root component
export const mockRootProps = {
  initialTab: "dashboard" as const,
  showWelcomeModal: false
};