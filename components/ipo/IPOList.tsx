import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Grid,
  Chip,
  Box,
  CircularProgress,
} from '@mui/material';
import { formatCurrency, formatDate } from '../../utils/formatters';

interface IPOAnnouncement {
  id: string;
  companyName: string;
  symbol: string;
  sector: string;
  offerPrice: number;
  listingDate: string;
  status: 'announced' | 'listed' | 'cancelled';
  createdAt: string;
  updatedAt: string;
}

const IPOList: React.FC = () => {
  const [ipos, setIpos] = useState<IPOAnnouncement[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchIPOs();
  }, []);

  const fetchIPOs = async () => {
    try {
      const response = await fetch('/api/v1/ipos');
      if (!response.ok) throw new Error('Failed to fetch IPOs');
      const data = await response.json();
      setIpos(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch IPOs');
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'announced': return 'primary';
      case 'listed': return 'success';
      case 'cancelled': return 'error';
      default: return 'default';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'announced': return 'Announced';
      case 'listed': return 'Listed';
      case 'cancelled': return 'Cancelled';
      default: return status;
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" p={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Typography color="error" align="center">
        {error}
      </Typography>
    );
  }

  if (ipos.length === 0) {
    return (
      <Box textAlign="center" py={4}>
        <Typography variant="h6" color="textSecondary" gutterBottom>
          No IPO announcements available
        </Typography>
        <Typography variant="body2" color="textSecondary">
          IPOs (Initial Public Offerings) are when private companies first sell shares to the public.<br/>
          Set up an IPO alert to be notified when new companies list on the Ghana Stock Exchange.
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="body1" color="textSecondary" sx={{ mb: 3 }}>
        🎆 Track upcoming Initial Public Offerings (IPOs) on the Ghana Stock Exchange. 
        These are companies going public for the first time, offering investment opportunities.
      </Typography>
      <Grid container spacing={3}>
      {ipos.map((ipo) => (
        <Grid item xs={12} md={6} lg={4} key={ipo.id}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={2}>
                <Typography variant="h6" component="h3">
                  {ipo.companyName}
                </Typography>
                <Chip
                  label={getStatusLabel(ipo.status)}
                  color={getStatusColor(ipo.status) as any}
                  size="small"
                />
              </Box>
              
              <Typography variant="body2" color="textSecondary" gutterBottom>
                Symbol: {ipo.symbol}
              </Typography>
              
              <Typography variant="body2" color="textSecondary" gutterBottom>
                Sector: {ipo.sector}
              </Typography>
              
              <Typography variant="body1" gutterBottom>
                Offer Price: {formatCurrency(ipo.offerPrice)}
              </Typography>
              
              <Typography variant="body2" color="textSecondary">
                {ipo.status === 'listed' ? 'Listed on' : 'Expected listing'}: {formatDate(ipo.listingDate)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      ))}
      </Grid>
    </Box>
  );
};

export default IPOList;