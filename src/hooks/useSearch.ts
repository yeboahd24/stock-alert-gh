import { useMemo } from 'react';

export const useStockSearch = <T extends { symbol: string; name: string }>(
  items: T[],
  searchTerm: string
) => {
  return useMemo(() => {
    if (!searchTerm.trim()) return items;
    
    const term = searchTerm.toLowerCase();
    return items.filter(item => 
      item.symbol.toLowerCase().includes(term) ||
      item.name.toLowerCase().includes(term)
    );
  }, [items, searchTerm]);
};

export const useAlertFilter = <T extends { status: string; alertType: string }>(
  items: T[],
  filter: string
) => {
  return useMemo(() => {
    if (filter === 'all') return items;
    
    return items.filter(item => 
      item.status === filter || item.alertType === filter
    );
  }, [items, filter]);
};