import React from 'react';
import { Paper, Typography, Grid, Box, Chip } from '@mui/material';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import { Stock } from '../../src/services/api';

interface MarketSummaryProps {
  stocks: Stock[];
}

const MarketSummary: React.FC<MarketSummaryProps> = ({ stocks }) => {
  const gainers = stocks.filter(s => s.change > 0).length;
  const losers = stocks.filter(s => s.change < 0).length;
  const unchanged = stocks.filter(s => s.change === 0).length;
  
  const totalVolume = stocks.reduce((sum, s) => sum + s.volume, 0);
  const avgChange = stocks.reduce((sum, s) => sum + s.changePercent, 0) / stocks.length;

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Market Summary
      </Typography>
      
      <Grid container spacing={2}>
        <Grid item xs={6} md={2}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              Gainers
            </Typography>
            <Box display="flex" alignItems="center" justifyContent="center">
              <TrendingUpIcon color="success" fontSize="small" />
              <Typography variant="h6" color="success.main" sx={{ ml: 0.5 }}>
                {gainers}
              </Typography>
            </Box>
          </Box>
        </Grid>

        <Grid item xs={6} md={2}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              Losers
            </Typography>
            <Box display="flex" alignItems="center" justifyContent="center">
              <TrendingDownIcon color="error" fontSize="small" />
              <Typography variant="h6" color="error.main" sx={{ ml: 0.5 }}>
                {losers}
              </Typography>
            </Box>
          </Box>
        </Grid>

        <Grid item xs={6} md={2}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              Unchanged
            </Typography>
            <Typography variant="h6">
              {unchanged}
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={6} md={3}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              Total Volume
            </Typography>
            <Typography variant="h6">
              {(totalVolume / 1000000).toFixed(1)}M
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={12} md={3}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              Market Trend
            </Typography>
            <Box>
              <Chip
                label={`${avgChange > 0 ? '+' : ''}${avgChange.toFixed(2)}%`}
                color={avgChange > 0 ? 'success' : avgChange < 0 ? 'error' : 'default'}
                variant="filled"
              />
            </Box>
          </Box>
        </Grid>
      </Grid>
    </Paper>
  );
};

export default MarketSummary;