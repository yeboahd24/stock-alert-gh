import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Stack,
  Grid,
  Card,
  CardContent,
  Chip,
  IconButton,
  CircularProgress,
  Link,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import { stockApi, DetailedStock } from '../../src/services/api';
import { formatPrice, formatPercentage, formatNumber } from '../../utils/formatters';

interface StockDetailsModalProps {
  open: boolean;
  onClose: () => void;
  symbol: string;
}

const StockDetailsModal: React.FC<StockDetailsModalProps> = ({ open, onClose, symbol }) => {
  const [stockDetails, setStockDetails] = useState<DetailedStock | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (open && symbol) {
      loadStockDetails();
    }
  }, [open, symbol]);

  const loadStockDetails = async () => {
    try {
      setLoading(true);
      setError(null);
      const details = await stockApi.getStockDetails(symbol);
      setStockDetails(details);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load stock details');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setStockDetails(null);
    setError(null);
    onClose();
  };

  if (!open) return null;

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Stack direction="row" justifyContent="space-between" alignItems="center">
          <Typography variant="h6">
            {stockDetails ? `${stockDetails.symbol} - ${stockDetails.name}` : `Stock Details - ${symbol}`}
          </Typography>
          <IconButton onClick={handleClose} size="small">
            <CloseIcon />
          </IconButton>
        </Stack>
      </DialogTitle>

      <DialogContent>
        {loading && (
          <Stack alignItems="center" spacing={2} sx={{ py: 4 }}>
            <CircularProgress />
            <Typography>Loading stock details...</Typography>
          </Stack>
        )}

        {error && (
          <Stack alignItems="center" spacing={2} sx={{ py: 4 }}>
            <Typography color="error">{error}</Typography>
            <Button onClick={loadStockDetails} variant="outlined">
              Retry
            </Button>
          </Stack>
        )}

        {stockDetails && (
          <Stack spacing={3}>
            {/* Price Information */}
            <Card>
              <CardContent>
                <Grid container spacing={3}>
                  <Grid item xs={12} sm={6}>
                    <Stack spacing={1}>
                      <Typography variant="h4" component="div" sx={{ fontWeight: 'bold' }}>
                        {formatPrice(stockDetails.currentPrice)}
                      </Typography>
                      <Stack direction="row" alignItems="center" spacing={1}>
                        {stockDetails.change >= 0 ? (
                          <TrendingUpIcon color="success" />
                        ) : (
                          <TrendingDownIcon color="error" />
                        )}
                        <Typography
                          variant="body1"
                          color={stockDetails.change >= 0 ? 'success.main' : 'error.main'}
                          sx={{ fontWeight: 'medium' }}
                        >
                          {stockDetails.change >= 0 ? '+' : ''}{formatPrice(stockDetails.change)} 
                          ({stockDetails.changePercent >= 0 ? '+' : ''}{formatPercentage(stockDetails.changePercent)})
                        </Typography>
                      </Stack>
                    </Stack>
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <Stack spacing={2}>
                      <Stack direction="row" justifyContent="space-between">
                        <Typography variant="body2" color="text.secondary">Previous Close:</Typography>
                        <Typography variant="body2">{formatPrice(stockDetails.previousClose)}</Typography>
                      </Stack>
                      <Stack direction="row" justifyContent="space-between">
                        <Typography variant="body2" color="text.secondary">Volume:</Typography>
                        <Typography variant="body2">{formatNumber(stockDetails.volume || 0)}</Typography>
                      </Stack>
                      <Stack direction="row" justifyContent="space-between">
                        <Typography variant="body2" color="text.secondary">Market Cap:</Typography>
                        <Typography variant="body2">{formatPrice(stockDetails.marketCap || 0)}</Typography>
                      </Stack>
                      <Stack direction="row" justifyContent="space-between">
                        <Typography variant="body2" color="text.secondary">Shares Outstanding:</Typography>
                        <Typography variant="body2">{formatNumber(stockDetails.shares || 0)}</Typography>
                      </Stack>
                    </Stack>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>

            {/* Financial Metrics */}
            {(stockDetails.eps !== null || stockDetails.dps !== null) && (
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>Financial Metrics</Typography>
                  <Grid container spacing={2}>
                    {stockDetails.eps !== null && (
                      <Grid item xs={6}>
                        <Stack direction="row" justifyContent="space-between">
                          <Typography variant="body2" color="text.secondary">EPS:</Typography>
                          <Typography variant="body2">{formatPrice(stockDetails.eps || 0)}</Typography>
                        </Stack>
                      </Grid>
                    )}
                    {stockDetails.dps !== null && (
                      <Grid item xs={6}>
                        <Stack direction="row" justifyContent="space-between">
                          <Typography variant="body2" color="text.secondary">DPS:</Typography>
                          <Typography variant="body2">{formatPrice(stockDetails.dps || 0)}</Typography>
                        </Stack>
                      </Grid>
                    )}
                  </Grid>
                </CardContent>
              </Card>
            )}

            {/* Sector & Industry */}
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>Classification</Typography>
                <Stack direction="row" spacing={1}>
                  <Chip label={stockDetails.sector} color="primary" variant="outlined" />
                  <Chip label={stockDetails.industry} color="secondary" variant="outlined" />
                </Stack>
              </CardContent>
            </Card>

            {/* Company Information */}
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>Company Information</Typography>
                <Stack spacing={2}>
                  <Stack>
                    <Typography variant="subtitle2" color="text.secondary">Company Name</Typography>
                    <Typography variant="body1">{stockDetails.company.name}</Typography>
                  </Stack>
                  
                  <Stack>
                    <Typography variant="subtitle2" color="text.secondary">Address</Typography>
                    <Typography variant="body2">{stockDetails.company.address}</Typography>
                  </Stack>

                  <Grid container spacing={2}>
                    <Grid item xs={12} sm={6}>
                      <Stack>
                        <Typography variant="subtitle2" color="text.secondary">Phone</Typography>
                        <Typography variant="body2">{stockDetails.company.telephone}</Typography>
                      </Stack>
                    </Grid>
                    <Grid item xs={12} sm={6}>
                      <Stack>
                        <Typography variant="subtitle2" color="text.secondary">Email</Typography>
                        <Link href={`mailto:${stockDetails.company.email}`} variant="body2">
                          {stockDetails.company.email}
                        </Link>
                      </Stack>
                    </Grid>
                  </Grid>

                  {stockDetails.company.website && (
                    <Stack>
                      <Typography variant="subtitle2" color="text.secondary">Website</Typography>
                      <Link 
                        href={stockDetails.company.website.startsWith('http') ? stockDetails.company.website : `https://${stockDetails.company.website}`} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        variant="body2"
                      >
                        {stockDetails.company.website}
                      </Link>
                    </Stack>
                  )}
                </Stack>
              </CardContent>
            </Card>

            {/* Last Updated */}
            <Typography variant="caption" color="text.secondary" textAlign="center">
              Last updated: {new Date(stockDetails.lastUpdated).toLocaleString()}
            </Typography>
          </Stack>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

export default StockDetailsModal;