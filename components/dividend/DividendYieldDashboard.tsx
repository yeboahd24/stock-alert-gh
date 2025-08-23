import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Chip,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TextField,
  InputAdornment,
  IconButton,
  Tooltip,
  Alert,
  CircularProgress,
  Stack,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  Search,
  Refresh,
  NotificationsActive,
  Info,
} from '@mui/icons-material';
import { dividendApi } from '../../src/services/api';
import { GSEDividendStock, GSEDividendResponse } from '../../types/schema';
import { formatPercentage, formatCurrency } from '../../utils/formatters';

interface DividendYieldDashboardProps {
  onCreateAlert?: (alertData: any) => void;
}

const DividendYieldDashboard: React.FC<DividendYieldDashboardProps> = ({ onCreateAlert }) => {
  const [dividendData, setDividendData] = useState<GSEDividendResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [minYieldFilter, setMinYieldFilter] = useState<number | ''>('');

  const fetchDividendData = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await dividendApi.getGSEDividendStocks();
      setDividendData(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch dividend data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDividendData();
  }, []);

  const filteredStocks = dividendData?.data.stocks.filter(stock => {
    const matchesSearch = stock.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         stock.symbol.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         stock.sector.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesYield = minYieldFilter === '' || stock.dividend_yield >= Number(minYieldFilter);
    
    return matchesSearch && matchesYield;
  }) || [];

  const highYieldStocks = filteredStocks.filter(stock => stock.dividend_yield >= 3.0);
  const averageYield = filteredStocks.length > 0 
    ? filteredStocks.reduce((sum, stock) => sum + stock.dividend_yield, 0) / filteredStocks.length 
    : 0;

  const handleCreateYieldAlert = (stock: GSEDividendStock, alertType: string) => {
    if (onCreateAlert) {
      const alertData = {
        stockSymbol: stock.symbol,
        stockName: stock.name,
        alertType,
        ...(alertType === 'high_dividend_yield' && { thresholdYield: stock.dividend_yield }),
        ...(alertType === 'target_dividend_yield' && { targetYield: stock.dividend_yield + 0.5 }),
        ...(alertType === 'dividend_yield_change' && { yieldChangeThreshold: 0.5 }),
      };
      onCreateAlert(alertData);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="body1" sx={{ ml: 2 }}>
          Loading dividend yield data...
        </Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" action={
        <Button color="inherit" size="small" onClick={fetchDividendData}>
          Retry
        </Button>
      }>
        {error}
      </Alert>
    );
  }

  return (
    <Box>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">
          📈 Dividend Yield Dashboard
        </Typography>
        <Button
          variant="outlined"
          startIcon={<Refresh />}
          onClick={fetchDividendData}
          disabled={loading}
        >
          Refresh Data
        </Button>
      </Box>

      {/* Summary Cards */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Stocks
              </Typography>
              <Typography variant="h4">
                {dividendData?.data.count || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                High Yield Stocks (≥3%)
              </Typography>
              <Typography variant="h4" color="success.main">
                {highYieldStocks.length}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Average Yield
              </Typography>
              <Typography variant="h4" color="primary">
                {formatPercentage(averageYield)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Last Updated
              </Typography>
              <Typography variant="body1">
                {dividendData?.data.timestamp ? 
                  new Date(dividendData.data.timestamp).toLocaleString() : 
                  'Unknown'
                }
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems="center">
            <TextField
              placeholder="Search stocks, symbols, or sectors..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search />
                  </InputAdornment>
                ),
              }}
              sx={{ flexGrow: 1 }}
            />
            <TextField
              label="Min Yield (%)"
              type="number"
              value={minYieldFilter}
              onChange={(e) => setMinYieldFilter(e.target.value === '' ? '' : Number(e.target.value))}
              inputProps={{ min: 0, step: 0.1 }}
              sx={{ width: 150 }}
            />
            <Typography variant="body2" color="textSecondary">
              Showing {filteredStocks.length} stocks
            </Typography>
          </Stack>
        </CardContent>
      </Card>

      {/* Data Source Info */}
      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body2">
          Data sourced from <strong>SimplyWall.St</strong> via GSE Dividends API. 
          Last updated: {dividendData?.data.timestamp ? new Date(dividendData.data.timestamp).toLocaleString() : 'Unknown'}
        </Typography>
      </Alert>

      {/* Dividend Stocks Table */}
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell><strong>Stock</strong></TableCell>
              <TableCell><strong>Sector</strong></TableCell>
              <TableCell align="right"><strong>Current Price</strong></TableCell>
              <TableCell align="right"><strong>Market Cap</strong></TableCell>
              <TableCell align="right"><strong>Dividend Yield</strong></TableCell>
              <TableCell align="center"><strong>Actions</strong></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredStocks.map((stock) => (
              <TableRow key={stock.symbol} hover>
                <TableCell>
                  <Box>
                    <Typography variant="subtitle2" fontWeight="bold">
                      {stock.symbol}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {stock.name}
                    </Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Chip 
                    label={stock.sector} 
                    size="small" 
                    variant="outlined"
                  />
                </TableCell>
                <TableCell align="right">
                  <Typography variant="body2">
                    {stock.price}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Typography variant="body2">
                    {stock.market_cap}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Box display="flex" alignItems="center" justifyContent="flex-end">
                    {stock.dividend_yield >= 3.0 ? (
                      <TrendingUp color="success" fontSize="small" />
                    ) : stock.dividend_yield > 0 ? (
                      <TrendingUp color="primary" fontSize="small" />
                    ) : (
                      <TrendingDown color="disabled" fontSize="small" />
                    )}
                    <Typography 
                      variant="subtitle2" 
                      fontWeight="bold"
                      color={stock.dividend_yield >= 3.0 ? 'success.main' : 'text.primary'}
                      sx={{ ml: 1 }}
                    >
                      {formatPercentage(stock.dividend_yield)}
                    </Typography>
                  </Box>
                </TableCell>
                <TableCell align="center">
                  <Stack direction="row" spacing={1} justifyContent="center">
                    <Tooltip title="Create High Yield Alert">
                      <IconButton
                        size="small"
                        color="success"
                        onClick={() => handleCreateYieldAlert(stock, 'high_dividend_yield')}
                      >
                        <NotificationsActive />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="View Details">
                      <IconButton
                        size="small"
                        color="primary"
                        onClick={() => window.open(stock.url, '_blank')}
                      >
                        <Info />
                      </IconButton>
                    </Tooltip>
                  </Stack>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {filteredStocks.length === 0 && (
        <Box textAlign="center" py={4}>
          <Typography variant="h6" color="textSecondary">
            No stocks match your current filters
          </Typography>
          <Typography variant="body2" color="textSecondary">
            Try adjusting your search terms or minimum yield filter
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default DividendYieldDashboard;