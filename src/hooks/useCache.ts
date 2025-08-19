import { useCallback } from 'react';
import { cache } from '../services/cache';

export const useCache = () => {
  const clearCache = useCallback(() => {
    cache.clear();
  }, []);

  const clearStockCache = useCallback(() => {
    cache.delete('stocks:all');
  }, []);

  const clearAlertsCache = useCallback(() => {
    cache.delete('alerts:all:all');
  }, []);

  return {
    clearCache,
    clearStockCache,
    clearAlertsCache,
  };
};