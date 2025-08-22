import React from 'react';
import { Box, Typography, Stack, Chip } from '@mui/material';
import { styled } from '@mui/material/styles';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';

const BrandContainer = styled(Box)(({ theme }) => ({
  background: `linear-gradient(135deg, ${theme.palette.primary.main} 0%, ${theme.palette.primary.dark} 100%)`,
  padding: theme.spacing(3, 0),
  position: 'relative',
  overflow: 'hidden',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    right: 0,
    width: '200px',
    height: '100%',
    background: `linear-gradient(45deg, ${theme.palette.secondary.main}20, transparent)`,
    borderRadius: '0 0 0 100px',
  },
  '&::after': {
    content: '""',
    position: 'absolute',
    bottom: 0,
    left: 0,
    width: '150px',
    height: '60%',
    background: `linear-gradient(-45deg, ${theme.palette.success.main}15, transparent)`,
    borderRadius: '0 100px 0 0',
  }
}));

const LogoContainer = styled(Stack)(({ theme }) => ({
  alignItems: 'center',
  position: 'relative',
  zIndex: 1,
}));

const BrandIcon = styled(Box)(() => ({
  width: 48,
  height: 48,
  borderRadius: '12px',
  background: 'linear-gradient(135deg, #FCD116, #fde047)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  marginRight: 16,
  boxShadow: '0 4px 12px rgba(252, 209, 22, 0.3)',
}));

const GhanaFlag = styled(Box)(() => ({
  width: 24,
  height: 16,
  position: 'relative',
  borderRadius: 2,
  overflow: 'hidden',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    width: '100%',
    height: '33.33%',
    backgroundColor: '#CE1126', // Ghana red
  },
  '&::after': {
    content: '""',
    position: 'absolute',
    bottom: 0,
    left: 0,
    width: '100%',
    height: '33.33%',
    backgroundColor: '#006B3F', // Ghana green
  },
  '& > div': {
    position: 'absolute',
    top: '33.33%',
    left: 0,
    width: '100%',
    height: '33.33%',
    backgroundColor: '#FCD116', // Ghana gold
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    '&::before': {
      content: '"★"',
      fontSize: '8px',
      color: '#000',
      fontWeight: 'bold',
    }
  }
}));

interface BrandHeaderProps {
  title?: string;
  subtitle?: string;
  showLiveIndicator?: boolean;
}

const BrandHeader: React.FC<BrandHeaderProps> = ({ 
  title = "Shares Alert Ghana",
  subtitle = "Ghana Stock Exchange Monitoring Platform",
  showLiveIndicator = false
}) => {
  return (
    <BrandContainer>
      <LogoContainer direction="row" spacing={2}>
        <BrandIcon>
          <GhanaFlag>
            <div />
          </GhanaFlag>
        </BrandIcon>
        <Stack>
          <Stack direction="row" alignItems="center" spacing={1}>
            <Typography 
              variant="h4" 
              component="h1" 
              sx={{ 
                color: 'white',
                fontWeight: 700,
                textShadow: '0 2px 4px rgba(0,0,0,0.3)'
              }}
            >
              {title}
            </Typography>
            {showLiveIndicator && (
              <Chip
                icon={<TrendingUpIcon />}
                label="LIVE"
                size="small"
                sx={{
                  backgroundColor: 'success.main',
                  color: 'white',
                  fontWeight: 600,
                  animation: 'pulse 2s infinite',
                  '@keyframes pulse': {
                    '0%': { opacity: 1 },
                    '50%': { opacity: 0.7 },
                    '100%': { opacity: 1 }
                  }
                }}
              />
            )}
          </Stack>
          <Typography 
            variant="subtitle1" 
            sx={{ 
              color: 'rgba(255,255,255,0.9)',
              fontWeight: 400
            }}
          >
            {subtitle}
          </Typography>
        </Stack>
      </LogoContainer>
    </BrandContainer>
  );
};

export default BrandHeader;