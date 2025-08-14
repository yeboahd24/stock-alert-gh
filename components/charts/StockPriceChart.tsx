import React from 'react';
import { LineChart } from '@mui/x-charts/LineChart';
import { Card, CardContent, Typography, Stack } from '@mui/material';
import { formatPrice } from '../../utils/formatters';

interface StockPriceChartProps {
  symbol: string;
  data: Array<{
    date: string;
    price: number;
  }>;
  height?: number;
}

const StockPriceChart: React.FC<StockPriceChartProps> = ({ 
  symbol, 
  data, 
  height = 300 
}) => {
  // Transform data for the chart
  const xAxisData = data.map(item => new Date(item.date));
  const seriesData = data.map(item => item.price);

  return (
    <Card>
      <CardContent>
        <Stack spacing={2}>
          <Typography variant="h6" component="h3">
            {symbol} Price Chart
          </Typography>
          
          <LineChart
            width={undefined}
            height={height}
            series={[
              {
                data: seriesData,
                label: `${symbol} Price`,
                color: '#1976d2',
                curve: 'linear',
              },
            ]}
            xAxis={[
              {
                data: xAxisData,
                scaleType: 'time',
                valueFormatter: (value: Date) => {
                  return value.toLocaleDateString('en-GB', {
                    month: 'short',
                    day: 'numeric',
                  });
                },
              },
            ]}
            yAxis={[
              {
                valueFormatter: (value: number) => formatPrice(value),
              },
            ]}
            grid={{ horizontal: true, vertical: true }}
            margin={{ left: 80, right: 20, top: 20, bottom: 60 }}
            slotProps={{
              legend: {
                position: { vertical: 'bottom', horizontal: 'center' },
              },
            }}
          />
        </Stack>
      </CardContent>
    </Card>
  );
};

export default StockPriceChart;