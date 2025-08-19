export interface PriceData {
  date: string;
  price: number;
  volume?: number;
}

export interface TechnicalIndicators {
  rsi: number;
  sma20: number;
  ema12: number;
  volumeAvg: number;
  volumeRatio: number;
}

export const calculateRSI = (prices: number[], period: number = 14): number => {
  if (prices.length < period + 1) return 50;
  
  let gains = 0;
  let losses = 0;
  
  for (let i = 1; i <= period; i++) {
    const change = prices[i] - prices[i - 1];
    if (change > 0) gains += change;
    else losses -= change;
  }
  
  const avgGain = gains / period;
  const avgLoss = losses / period;
  
  if (avgLoss === 0) return 100;
  const rs = avgGain / avgLoss;
  return 100 - (100 / (1 + rs));
};

export const calculateSMA = (prices: number[], period: number): number => {
  if (prices.length < period) return prices[prices.length - 1] || 0;
  
  const sum = prices.slice(-period).reduce((a, b) => a + b, 0);
  return sum / period;
};

export const calculateEMA = (prices: number[], period: number): number => {
  if (prices.length < period) return prices[prices.length - 1] || 0;
  
  const multiplier = 2 / (period + 1);
  let ema = prices.slice(0, period).reduce((a, b) => a + b, 0) / period;
  
  for (let i = period; i < prices.length; i++) {
    ema = (prices[i] * multiplier) + (ema * (1 - multiplier));
  }
  
  return ema;
};

export const calculateVolumeIndicators = (volumes: number[], currentVolume: number) => {
  if (volumes.length === 0) return { volumeAvg: currentVolume, volumeRatio: 1 };
  
  const volumeAvg = volumes.reduce((a, b) => a + b, 0) / volumes.length;
  const volumeRatio = currentVolume / volumeAvg;
  
  return { volumeAvg, volumeRatio };
};

export const generateMockHistoricalData = (currentPrice: number, days: number = 30): PriceData[] => {
  const data: PriceData[] = [];
  let price = currentPrice;
  
  for (let i = days; i >= 0; i--) {
    const date = new Date();
    date.setDate(date.getDate() - i);
    
    // Generate realistic price movement
    const change = (Math.random() - 0.5) * 0.05 * price;
    price = Math.max(0.01, price + change);
    
    data.push({
      date: date.toISOString().split('T')[0],
      price: parseFloat(price.toFixed(2)),
      volume: Math.floor(Math.random() * 100000) + 10000,
    });
  }
  
  return data;
};