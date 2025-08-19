import { useMemo } from 'react';
import { 
  calculateRSI, 
  calculateSMA, 
  calculateEMA, 
  calculateVolumeIndicators,
  generateMockHistoricalData,
  TechnicalIndicators
} from '../../utils/technicalIndicators';

export const useTechnicalIndicators = (
  currentPrice: number,
  currentVolume: number,
  symbol: string
): TechnicalIndicators => {
  return useMemo(() => {
    // Generate mock historical data for calculations
    const historicalData = generateMockHistoricalData(currentPrice, 30);
    const prices = historicalData.map(d => d.price);
    const volumes = historicalData.map(d => d.volume || 0);

    const rsi = calculateRSI(prices, 14);
    const sma20 = calculateSMA(prices, 20);
    const ema12 = calculateEMA(prices, 12);
    const { volumeAvg, volumeRatio } = calculateVolumeIndicators(volumes, currentVolume);

    return {
      rsi,
      sma20,
      ema12,
      volumeAvg,
      volumeRatio,
    };
  }, [currentPrice, currentVolume, symbol]);
};