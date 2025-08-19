import React from 'react';
import { Paper, Typography, List, ListItem, ListItemText, Box, Chip } from '@mui/material';
import { Stock } from '../../src/services/api';

interface TopMoversProps {
  stocks: Stock[];
  type: 'gainers' | 'losers';
  limit?: number;
}

const TopMovers: React.FC<TopMoversProps> = ({ stocks, type, limit = 5 }) => {
  const sortedStocks = [...stocks]
    .sort((a, b) => type === 'gainers' ? b.changePercent - a.changePercent : a.changePercent - b.changePercent)
    .filter(s => type === 'gainers' ? s.change > 0 : s.change < 0)
    .slice(0, limit);

  const title = type === 'gainers' ? 'Top Gainers' : 'Top Losers';
  const color = type === 'gainers' ? 'success' : 'error';

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        {title}
      </Typography>
      
      <List dense>
        {sortedStocks.length === 0 ? (
          <ListItem>
            <ListItemText primary={`No ${type} today`} />
          </ListItem>
        ) : (
          sortedStocks.map((stock) => (
            <ListItem key={stock.symbol} sx={{ px: 0 }}>
              <ListItemText
                primary={
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Box>
                      <Typography variant="body2" fontWeight="medium">
                        {stock.symbol}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        ₵{stock.currentPrice.toFixed(2)}
                      </Typography>
                    </Box>
                    <Chip
                      label={`${stock.changePercent > 0 ? '+' : ''}${stock.changePercent.toFixed(2)}%`}
                      color={color}
                      size="small"
                      variant="outlined"
                    />
                  </Box>
                }
              />
            </ListItem>
          ))
        )}
      </List>
    </Paper>
  );
};

export default TopMovers;