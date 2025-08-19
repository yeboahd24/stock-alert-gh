import React from 'react';
import { Paper, Typography, Grid, Box, Chip } from '@mui/material';
import { TechnicalIndicators as TechIndicators } from '../../utils/technicalIndicators';

interface TechnicalIndicatorsProps {
  indicators: TechIndicators;
  symbol: string;
}

const TechnicalIndicators: React.FC<TechnicalIndicatorsProps> = ({ indicators, symbol }) => {
  const getRSIColor = (rsi: number) => {
    if (rsi > 70) return 'error';
    if (rsi < 30) return 'success';
    return 'default';
  };

  const getVolumeColor = (ratio: number) => {
    if (ratio > 1.5) return 'error';
    if (ratio > 1.2) return 'warning';
    return 'success';
  };

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Technical Indicators - {symbol}
      </Typography>
      
      <Grid container spacing={2}>
        <Grid item xs={6} md={3}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              RSI (14)
            </Typography>
            <Box>
              <Chip
                label={indicators.rsi.toFixed(1)}
                color={getRSIColor(indicators.rsi)}
                variant="filled"
              />
            </Box>
          </Box>
        </Grid>

        <Grid item xs={6} md={3}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              SMA (20)
            </Typography>
            <Typography variant="h6">
              ₵{indicators.sma20.toFixed(2)}
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={6} md={3}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              EMA (12)
            </Typography>
            <Typography variant="h6">
              ₵{indicators.ema12.toFixed(2)}
            </Typography>
          </Box>
        </Grid>

        <Grid item xs={6} md={3}>
          <Box textAlign="center">
            <Typography variant="caption" color="text.secondary">
              Volume Ratio
            </Typography>
            <Box>
              <Chip
                label={`${indicators.volumeRatio.toFixed(1)}x`}
                color={getVolumeColor(indicators.volumeRatio)}
                variant="outlined"
              />
            </Box>
          </Box>
        </Grid>
      </Grid>
    </Paper>
  );
};

export default TechnicalIndicators;