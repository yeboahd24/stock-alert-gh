import React, { useState } from 'react';
import { Card, CardContent, Typography, Chip, Stack, IconButton, CardActionArea } from '@mui/material';
import { styled } from '@mui/material/styles';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import { formatPrice, formatPercentage } from '../../utils/formatters';
import StockDetailsModal from './StockDetailsModal';

const StyledCard = styled(Card)(({ theme }) => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
  background: 'linear-gradient(145deg, #ffffff 0%, #fafbfc 100%)',
  border: `1px solid ${theme.palette.grey[200]}`,
  '&:hover': {
    transform: 'translateY(-4px) scale(1.02)',
    boxShadow: `0 12px 40px rgba(206, 17, 38, 0.15)`,
    borderColor: theme.palette.primary.light,
  },
}));

const PriceChangeChip = styled(Chip)<{ isPositive: boolean }>(({ theme, isPositive }) => ({
  backgroundColor: isPositive ? theme.palette.success.main : theme.palette.error.main,
  color: theme.palette.common.white,
  fontWeight: 600,
  borderRadius: 8,
  fontSize: '0.75rem',
  height: 28,
  '& .MuiChip-label': {
    padding: '0 8px',
  },
}));

interface StockCardProps {
  symbol: string;
  name: string;
  currentPrice: number;
  change: number;
  changePercent: number;
  volume: number;
  hasAlert?: boolean;
  onMenuClick?: (symbol: string) => void;
}

const StockCard: React.FC<StockCardProps> = ({
  symbol,
  name,
  currentPrice,
  change,
  changePercent,
  volume,
  hasAlert = false,
  onMenuClick
}) => {
  const [detailsModalOpen, setDetailsModalOpen] = useState(false);
  const isPositive = change >= 0;

  return (
    <>
      <StyledCard>
        <CardActionArea onClick={() => setDetailsModalOpen(true)} sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
          <CardContent sx={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
        <Stack direction="row" justifyContent="space-between" alignItems="flex-start" sx={{ mb: 2 }}>
          <Stack>
            <Typography 
              variant="h6" 
              component="h3" 
              sx={{ 
                fontWeight: 700,
                color: 'primary.main',
                fontSize: '1.1rem'
              }}
            >
              {symbol}
            </Typography>
            <Typography 
              variant="body2" 
              color="text.secondary"
              sx={{ 
                fontSize: '0.85rem',
                fontWeight: 500
              }}
            >
              {name}
            </Typography>
          </Stack>
          <Stack direction="row" alignItems="center" spacing={1}>
            {hasAlert && (
              <Chip 
                label="🔔 Alert" 
                size="small" 
                sx={{
                  backgroundColor: 'secondary.main',
                  color: 'secondary.contrastText',
                  fontWeight: 600,
                  fontSize: '0.7rem'
                }}
              />
            )}
            <IconButton 
              size="small" 
              onClick={(e) => {
                e.stopPropagation();
                onMenuClick?.(symbol);
              }}
            >
              <MoreVertIcon />
            </IconButton>
          </Stack>
        </Stack>

        <Stack spacing={2}>
          <Stack direction="row" justifyContent="space-between" alignItems="center">
            <Typography 
              variant="h5" 
              sx={{ 
                fontWeight: 700,
                color: 'text.primary',
                fontSize: '1.4rem'
              }}
            >
              {formatPrice(currentPrice)}
            </Typography>
            <Stack direction="row" alignItems="center" spacing={0.5}>
              {isPositive ? (
                <TrendingUpIcon color="success" fontSize="small" />
              ) : (
                <TrendingDownIcon color="error" fontSize="small" />
              )}
              <PriceChangeChip
                isPositive={isPositive}
                label={`${isPositive ? '+' : ''}${formatPrice(change)} (${formatPercentage(changePercent)})`}
                size="small"
              />
            </Stack>
          </Stack>

          <Stack direction="row" justifyContent="space-between">
            <Typography 
              variant="body2" 
              color="text.secondary"
              sx={{ fontSize: '0.8rem' }}
            >
              Volume
            </Typography>
            <Typography 
              variant="body2" 
              sx={{ 
                fontWeight: 600,
                color: 'text.primary',
                fontSize: '0.85rem'
              }}
            >
              {volume.toLocaleString()}
            </Typography>
          </Stack>
          </Stack>
        </CardContent>
      </CardActionArea>
    </StyledCard>

    {/* Stock Details Modal */}
    <StockDetailsModal
      open={detailsModalOpen}
      onClose={() => setDetailsModalOpen(false)}
      symbol={symbol}
    />
  </>
  );
};

export default StockCard;