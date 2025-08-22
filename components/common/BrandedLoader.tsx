import React from 'react';
import { Box, CircularProgress, Typography, Stack } from '@mui/material';
import { styled, keyframes } from '@mui/material/styles';

const pulse = keyframes`
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
`;

const LoaderContainer = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '200px',
  padding: theme.spacing(4),
}));

const BrandIcon = styled(Box)(({ theme }) => ({
  width: 60,
  height: 60,
  borderRadius: '16px',
  background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.primary.dark})`,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  marginBottom: theme.spacing(2),
  animation: `${pulse} 2s ease-in-out infinite`,
  boxShadow: `0 8px 25px rgba(206, 17, 38, 0.3)`,
}));

const GhanaFlag = styled(Box)(() => ({
  width: 32,
  height: 22,
  position: 'relative',
  borderRadius: 3,
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
      fontSize: '12px',
      color: '#000',
      fontWeight: 'bold',
    }
  }
}));

interface BrandedLoaderProps {
  message?: string;
  size?: 'small' | 'medium' | 'large';
}

const BrandedLoader: React.FC<BrandedLoaderProps> = ({ 
  message = "Loading Ghana Stock Exchange data...",
  size = 'medium'
}) => {
  const getSize = () => {
    switch (size) {
      case 'small': return 30;
      case 'large': return 50;
      default: return 40;
    }
  };

  return (
    <LoaderContainer>
      <Stack alignItems="center" spacing={2}>
        <Box sx={{ position: 'relative' }}>
          <BrandIcon>
            <GhanaFlag>
              <div />
            </GhanaFlag>
          </BrandIcon>
          <CircularProgress
            size={getSize()}
            thickness={3}
            sx={{
              position: 'absolute',
              top: '50%',
              left: '50%',
              marginTop: `-${getSize() / 2}px`,
              marginLeft: `-${getSize() / 2}px`,
              color: 'secondary.main',
            }}
          />
        </Box>
        <Typography 
          variant="body1" 
          color="text.secondary"
          sx={{ 
            fontWeight: 500,
            textAlign: 'center',
            maxWidth: 300
          }}
        >
          {message}
        </Typography>
      </Stack>
    </LoaderContainer>
  );
};

export default BrandedLoader;