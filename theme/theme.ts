// Theme configuration for Shares Alert Ghana application
import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2', // Blue for primary actions and branding
      light: '#42a5f5',
      dark: '#1565c0',
      contrastText: '#ffffff'
    },
    secondary: {
      main: '#388e3c', // Green for positive stock movements and success states
      light: '#66bb6a',
      dark: '#2e7d32',
      contrastText: '#ffffff'
    },
    error: {
      main: '#d32f2f', // Red for negative stock movements and alerts
      light: '#ef5350',
      dark: '#c62828',
      contrastText: '#ffffff'
    },
    warning: {
      main: '#f57c00', // Orange for warning states and pending alerts
      light: '#ff9800',
      dark: '#e65100',
      contrastText: '#ffffff'
    },
    info: {
      main: '#0288d1', // Light blue for informational elements
      light: '#03a9f4',
      dark: '#01579b',
      contrastText: '#ffffff'
    },
    success: {
      main: '#388e3c', // Green for successful operations
      light: '#4caf50',
      dark: '#2e7d32',
      contrastText: '#ffffff'
    },
    text: {
      primary: '#212121',
      secondary: '#757575',
      disabled: '#bdbdbd'
    },
    background: {
      default: '#fafafa',
      paper: '#ffffff'
    },
    grey: {
      50: '#fafafa',
      100: '#f5f5f5',
      200: '#eeeeee',
      300: '#e0e0e0',
      400: '#bdbdbd',
      500: '#9e9e9e',
      600: '#757575',
      700: '#616161',
      800: '#424242',
      900: '#212121'
    },
    divider: '#e0e0e0'
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    fontSize: 14,
    fontWeightLight: 300,
    fontWeightRegular: 400,
    fontWeightMedium: 500,
    fontWeightBold: 700
  },
  shape: {
    borderRadius: 8
  }
});

export default theme;