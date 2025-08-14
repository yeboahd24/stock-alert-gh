import React, { useState } from 'react';
import { Card, CardContent, Typography, Chip, Stack, IconButton, CardActionArea } from '@mui/material';
import { styled } from '@mui/material/styles';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import { formatPrice, formatPercentage } from '../../utils/formatters';
import StockDetailsModal from './StockDetailsModal';

const StyledCard = styled(Card)(({ theme }) => ({
  transition: 'all 0.2s ease-in-out',
  '&:hover': {
    transform: 'translateY(-2px)',
    boxShadow: theme.shadows[4],
  },
}));

const PriceChangeChip = styled(Chip)<{ isPositive: boolean }>(({ theme, isPositive }) => ({
  backgroundColor: isPositive ? theme.palette.success.main : theme.palette.error.main,
  color: theme.palette.common.white,
  fontWeight: theme.typography.fontWeightMedium,
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
        <CardActionArea onClick={() => setDetailsModalOpen(true)}>
          <CardContent>
        <Stack direction="row" justifyContent="space-between" alignItems="flex-start" sx={{ mb: 2 }}>
          <Stack>
            <Typography variant="h6" component="h3" sx={{ fontWeight: 'bold' }}>
              {symbol}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {name}
            </Typography>
          </Stack>
          <Stack direction="row" alignItems="center" spacing={1}>
            {hasAlert && (
              <Chip 
                label="Alert Set" 
                size="small" 
                color="info" 
                variant="outlined"
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
            <Typography variant="h5" sx={{ fontWeight: 'bold' }}>
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
            <Typography variant="body2" color="text.secondary">
              Volume
            </Typography>
            <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
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