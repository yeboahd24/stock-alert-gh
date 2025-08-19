import { useEffect, useState } from 'react';
import { wsService } from '../services/websocket';
import { Stock } from '../services/api';

export const useWebSocket = () => {
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    wsService.connect();
    
    const handleOpen = () => setIsConnected(true);
    const handleClose = () => setIsConnected(false);
    
    wsService.on('open', handleOpen);
    wsService.on('close', handleClose);

    return () => {
      wsService.off('open', handleOpen);
      wsService.off('close', handleClose);
    };
  }, []);

  return { isConnected };
};

export const useStockUpdates = (onUpdate: (stocks: Stock[]) => void) => {
  useEffect(() => {
    wsService.on('stock_update', onUpdate);
    
    return () => {
      wsService.off('stock_update', onUpdate);
    };
  }, [onUpdate]);
};