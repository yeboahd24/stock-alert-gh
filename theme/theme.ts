// Ghana-inspired theme for Shares Alert Ghana application
import { createTheme } from '@mui/material/styles';

// Ghana flag colors: Red, Gold (Yellow), Green
const ghanaColors = {
  red: '#CE1126',
  gold: '#FCD116', 
  green: '#006B3F',
  darkGreen: '#004d2e',
  lightGold: '#fde047',
  darkRed: '#9f0f1f'
};

const theme = createTheme({
  palette: {
    primary: {
      main: ghanaColors.red, // Ghana red for primary branding
      light: '#e53e3e',
      dark: ghanaColors.darkRed,
      contrastText: '#ffffff'
    },
    secondary: {
      main: ghanaColors.gold, // Ghana gold for accents
      light: ghanaColors.lightGold,
      dark: '#d4af37',
      contrastText: '#000000'
    },
    success: {
      main: ghanaColors.green, // Ghana green for positive movements
      light: '#38a169',
      dark: ghanaColors.darkGreen,
      contrastText: '#ffffff'
    },
    error: {
      main: '#dc2626', // Slightly different red for errors
      light: '#ef4444',
      dark: '#991b1b',
      contrastText: '#ffffff'
    },
    warning: {
      main: '#f59e0b', // Warm orange for warnings
      light: '#fbbf24',
      dark: '#d97706',
      contrastText: '#000000'
    },
    info: {
      main: '#0ea5e9', // Sky blue for info
      light: '#38bdf8',
      dark: '#0284c7',
      contrastText: '#ffffff'
    },
    text: {
      primary: '#1f2937', // Dark gray for better readability
      secondary: '#6b7280',
      disabled: '#9ca3af'
    },
    background: {
      default: '#f9fafb', // Very light gray
      paper: '#ffffff'
    },
    grey: {
      50: '#f9fafb',
      100: '#f3f4f6',
      200: '#e5e7eb',
      300: '#d1d5db',
      400: '#9ca3af',
      500: '#6b7280',
      600: '#4b5563',
      700: '#374151',
      800: '#1f2937',
      900: '#111827'
    },
    divider: '#e5e7eb'
  },
  typography: {
    fontFamily: '"Inter", "Poppins", "Roboto", "Helvetica", "Arial", sans-serif',
    fontSize: 14,
    fontWeightLight: 300,
    fontWeightRegular: 400,
    fontWeightMedium: 500,
    fontWeightBold: 700,
    h1: {
      fontWeight: 700,
      fontSize: '2.5rem',
      lineHeight: 1.2,
      letterSpacing: '-0.02em'
    },
    h2: {
      fontWeight: 700,
      fontSize: '2rem',
      lineHeight: 1.3,
      letterSpacing: '-0.01em'
    },
    h3: {
      fontWeight: 600,
      fontSize: '1.75rem',
      lineHeight: 1.3
    },
    h4: {
      fontWeight: 600,
      fontSize: '1.5rem',
      lineHeight: 1.4
    },
    h5: {
      fontWeight: 600,
      fontSize: '1.25rem',
      lineHeight: 1.4
    },
    h6: {
      fontWeight: 600,
      fontSize: '1.125rem',
      lineHeight: 1.4
    },
    subtitle1: {
      fontWeight: 500,
      fontSize: '1rem',
      lineHeight: 1.5
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.6
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.5
    },
    button: {
      fontWeight: 600,
      textTransform: 'none' as const,
      letterSpacing: '0.02em'
    }
  },
  shape: {
    borderRadius: 12
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontWeight: 600,
          borderRadius: 12,
          padding: '10px 24px',
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)'
          }
        },
        contained: {
          '&:hover': {
            boxShadow: '0 6px 20px rgba(0, 0, 0, 0.2)'
          }
        }
      }
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 16,
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
          border: '1px solid #f3f4f6',
          '&:hover': {
            boxShadow: '0 8px 25px rgba(0, 0, 0, 0.12)'
          }
        }
      }
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0 1px 3px rgba(0, 0, 0, 0.08)'
        }
      }
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          fontWeight: 500
        }
      }
    },
    MuiTab: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontWeight: 500,
          fontSize: '0.95rem'
        }
      }
    },
    MuiFab: {
      styleOverrides: {
        root: {
          boxShadow: '0 4px 20px rgba(206, 17, 38, 0.3)'
        }
      }
    }
  }
});

export default theme;