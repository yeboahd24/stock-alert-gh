import React from 'react';
import { Box, Typography, Stack, Link, Divider, Chip } from '@mui/material';
import { styled } from '@mui/material/styles';
import FavoriteIcon from '@mui/icons-material/Favorite';
import PublicIcon from '@mui/icons-material/Public';

const FooterContainer = styled(Box)(({ theme }) => ({
  background: `linear-gradient(135deg, ${theme.palette.grey[900]} 0%, ${theme.palette.grey[800]} 100%)`,
  padding: theme.spacing(4, 0),
  marginTop: 'auto',
  position: 'relative',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    height: '4px',
    background: `linear-gradient(90deg, ${theme.palette.primary.main} 0%, ${theme.palette.secondary.main} 50%, ${theme.palette.success.main} 100%)`,
  }
}));

const BrandSection = styled(Stack)(({ theme }) => ({
  alignItems: 'center',
  marginBottom: theme.spacing(3),
}));

const GhanaFlag = styled(Box)(() => ({
  width: 24,
  height: 16,
  position: 'relative',
  borderRadius: 2,
  overflow: 'hidden',
  marginRight: 8,
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    width: '100%',
    height: '33.33%',
    backgroundColor: '#CE1126',
  },
  '&::after': {
    content: '""',
    position: 'absolute',
    bottom: 0,
    left: 0,
    width: '100%',
    height: '33.33%',
    backgroundColor: '#006B3F',
  },
  '& > div': {
    position: 'absolute',
    top: '33.33%',
    left: 0,
    width: '100%',
    height: '33.33%',
    backgroundColor: '#FCD116',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    '&::before': {
      content: '"★"',
      fontSize: '6px',
      color: '#000',
      fontWeight: 'bold',
    }
  }
}));

const Footer: React.FC = () => {
  return (
    <FooterContainer component="footer">
      <Stack
        spacing={3}
        sx={{ maxWidth: 'lg', mx: 'auto', px: 3 }}
      >
        <BrandSection direction="row" spacing={1}>
          <GhanaFlag>
            <div />
          </GhanaFlag>
          <Typography 
            variant="h6" 
            sx={{ 
              color: 'white',
              fontWeight: 700,
              fontSize: '1.1rem'
            }}
          >
            Shares Alert Ghana
          </Typography>
          <Chip
            icon={<PublicIcon />}
            label="GSE Official"
            size="small"
            sx={{
              backgroundColor: 'success.main',
              color: 'white',
              fontSize: '0.7rem',
              height: 24
            }}
          />
        </BrandSection>
        
        <Divider sx={{ borderColor: 'grey.700' }} />
        
        <Stack
          direction={{ xs: 'column', md: 'row' }}
          justifyContent="space-between"
          alignItems="center"
          spacing={2}
        >
          <Stack direction="row" alignItems="center" spacing={1}>
            <Typography variant="body2" sx={{ color: 'grey.400' }}>
              © 2025 Dominic Kofi Yeboah. Made with
            </Typography>
            <FavoriteIcon sx={{ color: 'primary.main', fontSize: 16 }} />
            <Typography variant="body2" sx={{ color: 'grey.400' }}>
              for Ghana's investors
            </Typography>
          </Stack>
          
          <Stack direction="row" spacing={3}>
            <Link 
              href="#" 
              variant="body2" 
              sx={{ 
                color: 'grey.300',
                textDecoration: 'none',
                '&:hover': {
                  color: 'secondary.main',
                  textDecoration: 'underline'
                }
              }}
            >
              Privacy Policy
            </Link>
            <Link 
              href="#" 
              variant="body2" 
              sx={{ 
                color: 'grey.300',
                textDecoration: 'none',
                '&:hover': {
                  color: 'secondary.main',
                  textDecoration: 'underline'
                }
              }}
            >
              Terms of Service
            </Link>
            <Link 
              href="mailto:yeboahd24@gmail.com" 
              variant="body2" 
              sx={{ 
                color: 'grey.300',
                textDecoration: 'none',
                '&:hover': {
                  color: 'secondary.main',
                  textDecoration: 'underline'
                }
              }}
            >
              Contact Us
            </Link>
          </Stack>
          
          <Typography 
            variant="body2" 
            sx={{ 
              color: 'grey.500',
              fontSize: '0.8rem'
            }}
          >
            Powered by Ghana Stock Exchange API
          </Typography>
        </Stack>
      </Stack>
    </FooterContainer>
  );
};

export default Footer;