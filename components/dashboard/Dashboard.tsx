import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Stack,
  Button,
  Paper,
  Tabs,
  Tab,
  Box,
  Fab,
  Snackbar,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import DashboardIcon from '@mui/icons-material/Dashboard';
import NotificationsIcon from '@mui/icons-material/Notifications';
import SettingsIcon from '@mui/icons-material/Settings';
import PersonIcon from '@mui/icons-material/Person';
import NotificationAddIcon from '@mui/icons-material/NotificationAdd';
import StockCard from '../common/StockCard';
import StockPriceChart from '../charts/StockPriceChart';
import AlertForm, { AlertFormData } from '../forms/AlertForm';
import NotificationSettings from '../forms/NotificationSettings';
import AlertsTable from '../tables/AlertsTable';
import { AlertType, AlertStatus } from '../../types/enums';
import { stockApi, alertApi, Stock, Alert as ApiAlert, setAuthToken } from '../../src/services/api';
import { useAuth } from '../../src/contexts/AuthContext';
import { useWebSocket, useStockUpdates } from '../../src/hooks/useWebSocket';
import { useStockSearch, useAlertFilter } from '../../src/hooks/useSearch';
import { useTechnicalIndicators } from '../../src/hooks/useTechnicalIndicators';
import UserMenu from '../common/UserMenu';
import UserProfile from '../profile/UserProfile';
import SearchBar from '../common/SearchBar';
import FilterChips from '../common/FilterChips';
import TechnicalIndicators from '../charts/TechnicalIndicators';
import MarketSummary from '../market/MarketSummary';
import TopMovers from '../market/TopMovers';
import Footer from '../common/Footer';
import BrandHeader from '../common/BrandHeader';

const StyledContainer = styled(Container)(({ theme }) => ({
  paddingTop: theme.spacing(3),
  paddingBottom: theme.spacing(3),
}));

const TabPanel = ({ children, value, index }: { children: React.ReactNode; value: number; index: number }) => (
  <div hidden={value !== index}>
    {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
  </div>
);

const TechnicalIndicatorsWrapper: React.FC<{ stock: Stock }> = ({ stock }) => {
  const indicators = useTechnicalIndicators(
    stock.currentPrice,
    stock.volume,
    stock.symbol
  );
  
  return (
    <TechnicalIndicators
      indicators={indicators}
      symbol={stock.symbol}
    />
  );
};

const Dashboard: React.FC = () => {
  const { user, token } = useAuth();
  const { isConnected } = useWebSocket();
  const [currentTab, setCurrentTab] = useState(0);
  const [alertFormOpen, setAlertFormOpen] = useState(false);
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [alerts, setAlerts] = useState<ApiAlert[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedChartStock, setSelectedChartStock] = useState<string>('');
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({
    open: false,
    message: '',
    severity: 'success',
  });
  const [stockSearch, setStockSearch] = useState('');
  const [alertFilter, setAlertFilter] = useState('all');

  // Handle real-time stock updates
  useStockUpdates((updatedStocks: Stock[]) => {
    setStocks(updatedStocks);
  });

  // Filter stocks and alerts using optimized hooks
  const filteredStocks = useStockSearch(stocks, stockSearch);
  const filteredAlerts = useAlertFilter(alerts, alertFilter);

  // Alert filter options
  const alertFilters = [
    { label: 'All', value: 'all', active: alertFilter === 'all' },
    { label: 'Active', value: 'active', active: alertFilter === 'active' },
    { label: 'Triggered', value: 'triggered', active: alertFilter === 'triggered' },
    { label: 'Price Above', value: 'price_above', active: alertFilter === 'price_above' },
    { label: 'Price Below', value: 'price_below', active: alertFilter === 'price_below' },
  ];

  // Set auth token when component mounts
  useEffect(() => {
    if (token) {
      setAuthToken(token);
    }
  }, [token]);

  // Load data from API
  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        const [stocksData, alertsData] = await Promise.all([
          stockApi.getAllStocks(),
          alertApi.getAllAlerts() // Now uses authenticated user
        ]);
        setStocks(stocksData || []);
        setAlerts(alertsData || []);
        // Set the first stock as default for chart
        if (stocksData.length > 0 && !selectedChartStock) {
          setSelectedChartStock(stocksData[0].symbol);
        }
      } catch (error) {
        console.error('Failed to load data:', error);
        console.error('Error details:', error);
        setSnackbar({
          open: true,
          message: `Failed to load data from server: ${error instanceof Error ? error.message : 'Unknown error'}`,
          severity: 'error',
        });
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  // Generate chart data based on the selected stock
  const getChartData = () => {
    if (stocks.length === 0 || !selectedChartStock) return [];
    
    const selectedStock = stocks.find(stock => stock.symbol === selectedChartStock) || stocks[0];
    const currentPrice = selectedStock.currentPrice;
    const change = selectedStock.change;
    
    // Generate 7 days of mock historical data based on current price
    const dates: string[] = [];
    const prices: number[] = [];
    
    for (let i = 6; i >= 0; i--) {
      const date = new Date();
      date.setDate(date.getDate() - i);
      dates.push(date.toISOString().split('T')[0]);
      
      // Generate realistic price variations around current price
      const variation = (Math.random() - 0.5) * 0.1 * currentPrice;
      const price = i === 0 ? currentPrice : currentPrice - change + variation;
      prices.push(Math.max(0.01, price)); // Ensure positive price
    }
    
    return dates.map((date, index) => ({
      date,
      price: prices[index]
    }));
  };

  const chartData = getChartData();

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  const handleCreateAlert = async (alertData: AlertFormData) => {
    try {
      const newAlert = await alertApi.createAlert({
        stockSymbol: alertData.stockSymbol,
        stockName: alertData.stockName,
        alertType: alertData.alertType,
        thresholdPrice: alertData.thresholdPrice,
      });
      
      // Add the new alert to the local state
      setAlerts(prev => [...prev, newAlert]);
      
      setSnackbar({
        open: true,
        message: `Alert created for ${alertData.stockSymbol}`,
        severity: 'success',
      });
    } catch (error) {
      console.error('Failed to create alert:', error);
      setSnackbar({
        open: true,
        message: 'Failed to create alert',
        severity: 'error',
      });
    }
  };

  const handleSaveNotificationSettings = (settings: any) => {
    // Mock settings save
    console.log('Saving notification settings:', settings);
    setSnackbar({
      open: true,
      message: 'Notification settings saved successfully',
      severity: 'success',
    });
  };

  const handleEditAlert = (alert: any) => {
    console.log('Editing alert:', alert);
    // Would open edit form
  };

  const handleDeleteAlert = async (alertId: string) => {
    try {
      await alertApi.deleteAlert(alertId);
      
      // Remove the alert from local state
      setAlerts(prev => prev.filter(alert => alert.id !== alertId));
      
      setSnackbar({
        open: true,
        message: 'Alert deleted successfully',
        severity: 'success',
      });
    } catch (error) {
      console.error('Failed to delete alert:', error);
      setSnackbar({
        open: true,
        message: 'Failed to delete alert',
        severity: 'error',
      });
    }
  };

  const handleStockMenuClick = (symbol: string) => {
    console.log('Stock menu clicked for:', symbol);
    // Would show stock menu options
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <StyledContainer maxWidth="lg" sx={{ flex: 1 }}>
        <Stack spacing={3}>
        {/* Brand Header */}
        <Box sx={{ mx: -3, mb: 3 }}>
          <BrandHeader 
            showLiveIndicator={isConnected}
            subtitle={`Welcome back, ${user?.name} • Ghana Stock Exchange Platform`}
          />
        </Box>
        
        {/* User Menu */}
        <Stack direction="row" justifyContent="flex-end" sx={{ mb: 2 }}>
          <UserMenu 
            onOpenSettings={() => setCurrentTab(2)}
            onOpenNotifications={() => setCurrentTab(2)}
          />
        </Stack>

        {/* Navigation Tabs */}
        <Paper sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={currentTab} onChange={handleTabChange}>
            <Tab icon={<DashboardIcon />} label="Dashboard" />
            <Tab icon={<NotificationsIcon />} label="My Alerts" />
            <Tab icon={<SettingsIcon />} label="Notifications" />
            <Tab icon={<PersonIcon />} label="Profile" />
          </Tabs>
        </Paper>

        {/* Tab Panels */}
        <TabPanel value={currentTab} index={0}>
          <Stack spacing={3}>
            {/* Welcome Section */}
            <Paper 
              sx={{ 
                p: 4, 
                background: 'linear-gradient(135deg, #006B3F 0%, #004d2e 100%)',
                color: 'white',
                position: 'relative',
                overflow: 'hidden',
                '&::before': {
                  content: '""',
                  position: 'absolute',
                  top: -50,
                  right: -50,
                  width: 100,
                  height: 100,
                  background: 'rgba(252, 209, 22, 0.1)',
                  borderRadius: '50%',
                },
                '&::after': {
                  content: '""',
                  position: 'absolute',
                  bottom: -30,
                  left: -30,
                  width: 80,
                  height: 80,
                  background: 'rgba(206, 17, 38, 0.1)',
                  borderRadius: '50%',
                }
              }}
            >
              <Box sx={{ position: 'relative', zIndex: 1 }}>
                <Typography variant="h5" gutterBottom sx={{ fontWeight: 700 }}>
                  📈 Ghana's Premier Stock Alert Platform
                </Typography>
                <Typography variant="body1" sx={{ mb: 3, opacity: 0.95, lineHeight: 1.6 }}>
                  Stay ahead of the market with real-time alerts for Ghana Stock Exchange (GSE) stocks. 
                  Never miss price movements, IPO launches, or dividend announcements again.
                </Typography>
                <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                  <Button 
                    variant="contained" 
                    color="secondary"
                    onClick={() => setAlertFormOpen(true)}
                    startIcon={<NotificationAddIcon />}
                    sx={{ 
                      fontWeight: 600,
                      py: 1.5,
                      px: 3,
                      boxShadow: '0 4px 15px rgba(252, 209, 22, 0.4)'
                    }}
                  >
                    Create Your First Alert
                  </Button>
                  <Button 
                    variant="outlined" 
                    sx={{ 
                      color: 'white', 
                      borderColor: 'rgba(255,255,255,0.5)',
                      fontWeight: 600,
                      py: 1.5,
                      px: 3,
                      '&:hover': {
                        borderColor: 'white',
                        backgroundColor: 'rgba(255,255,255,0.1)'
                      }
                    }}
                    onClick={() => setCurrentTab(1)}
                  >
                    View My Alerts
                  </Button>
                </Stack>
              </Box>
            </Paper>

            {/* How It Works Section */}
            {alerts.length === 0 && (
              <Paper sx={{ p: 3 }}>
                <Typography variant="h6" gutterBottom>
                  🚀 How to Get Started
                </Typography>
                <Stack spacing={2}>
                  <Box display="flex" alignItems="center" gap={2}>
                    <Typography variant="h6" color="primary">1.</Typography>
                    <Typography>Browse stocks below and click "Create Alert" on any stock card</Typography>
                  </Box>
                  <Box display="flex" alignItems="center" gap={2}>
                    <Typography variant="h6" color="primary">2.</Typography>
                    <Typography>Choose alert type: Price alerts, IPO notifications, or dividend announcements</Typography>
                  </Box>
                  <Box display="flex" alignItems="center" gap={2}>
                    <Typography variant="h6" color="primary">3.</Typography>
                    <Typography>Receive instant email notifications when your conditions are met</Typography>
                  </Box>
                </Stack>
              </Paper>
            )}

            {/* Market Summary */}
            {stocks.length > 0 && (
              <MarketSummary stocks={stocks} />
            )}

            {/* Top Movers */}
            {stocks.length > 0 && (
              <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                <Box flex={1}>
                  <TopMovers stocks={stocks} type="gainers" />
                </Box>
                <Box flex={1}>
                  <TopMovers stocks={stocks} type="losers" />
                </Box>
              </Stack>
            )}
            {/* Stock Cards Grid */}
            <Stack spacing={2}>
              <Stack direction="row" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography variant="h6" component="h2">
                    Ghana Stock Exchange - Live Prices
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Real-time stock prices from the Ghana Stock Exchange. Click on any stock to create alerts.
                  </Typography>
                </Box>
                <Box sx={{ width: 300 }}>
                  <SearchBar
                    value={stockSearch}
                    onChange={setStockSearch}
                    placeholder="Search stocks..."
                  />
                </Box>
              </Stack>
              <Stack direction="row" spacing={2} sx={{ overflowX: 'auto', pb: 1 }}>
                {loading ? (
                  <Typography>Loading stocks...</Typography>
                ) : filteredStocks.length > 0 ? (
                  filteredStocks.map((stock) => (
                    <Box key={stock.symbol} sx={{ minWidth: 300 }}>
                      <StockCard
                        symbol={stock.symbol}
                        name={stock.name}
                        currentPrice={stock.currentPrice}
                        change={stock.change}
                        changePercent={stock.changePercent}
                        volume={stock.volume}
                        hasAlert={alerts.some(alert => alert.stockSymbol === stock.symbol)}
                        onMenuClick={handleStockMenuClick}
                      />
                    </Box>
                  ))
                ) : stockSearch ? (
                  <Typography color="text.secondary">No stocks found matching "{stockSearch}"</Typography>
                ) : (
                  <Paper sx={{ p: 3, textAlign: 'center' }}>
                    <Typography variant="h6" color="textSecondary" gutterBottom>
                      📊 Loading Ghana Stock Exchange Data...
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      We're fetching the latest stock prices from GSE. This may take a moment.
                    </Typography>
                  </Paper>
                )}
              </Stack>
            </Stack>

            {/* Technical Indicators */}
            {stocks.length > 0 && selectedChartStock && (() => {
              const selectedStock = stocks.find(s => s.symbol === selectedChartStock) || stocks[0];
              return (
                <TechnicalIndicatorsWrapper
                  stock={selectedStock}
                />
              );
            })()}

            {/* Chart */}
            <Stack spacing={2}>
              <Stack direction="row" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography variant="h6" component="h2">
                    Stock Price Trends
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    7-day price movement for selected stock. Use this to identify trends and set better alert prices.
                  </Typography>
                </Box>
                {stocks.length > 0 && (
                  <FormControl size="small" sx={{ minWidth: 200 }}>
                    <InputLabel>Select Stock</InputLabel>
                    <Select
                      value={selectedChartStock}
                      label="Select Stock"
                      onChange={(e) => setSelectedChartStock(e.target.value)}
                    >
                      {stocks.map((stock) => (
                        <MenuItem key={stock.symbol} value={stock.symbol}>
                          {stock.symbol} - {stock.name}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                )}
              </Stack>
              {stocks.length > 0 && selectedChartStock ? (
                <StockPriceChart
                  symbol={selectedChartStock}
                  data={chartData}
                  height={400}
                />
              ) : (
                <Typography>Loading chart data...</Typography>
              )}
            </Stack>
          </Stack>
        </TabPanel>

        <TabPanel value={currentTab} index={1}>
          <Stack spacing={3}>
            {/* Alerts Header */}
            <Paper sx={{ p: 2, bgcolor: 'grey.50' }}>
              <Typography variant="h6" gutterBottom>
                🔔 Your Stock Alerts ({filteredAlerts.length})
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Manage your personalized alerts for Ghana Stock Exchange stocks. 
                Get notified via email when your conditions are met.
              </Typography>
            </Paper>
            
            <Stack direction="row" justifyContent="space-between" alignItems="center">
              <Typography variant="subtitle1" component="h3">
                Active Alerts
              </Typography>
              <Button
                variant="contained"
                startIcon={<NotificationAddIcon />}
                onClick={() => setAlertFormOpen(true)}
              >
                Create Alert
              </Button>
            </Stack>
            <FilterChips
              filters={alertFilters}
              onFilterChange={setAlertFilter}
            />
            
            {filteredAlerts.length === 0 ? (
              <Paper sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="h6" color="textSecondary" gutterBottom>
                  🚨 No Alerts Yet
                </Typography>
                <Typography variant="body1" color="textSecondary" sx={{ mb: 2 }}>
                  Create your first alert to start monitoring Ghana Stock Exchange stocks.
                </Typography>
                <Button 
                  variant="contained" 
                  onClick={() => setAlertFormOpen(true)}
                  startIcon={<NotificationAddIcon />}
                >
                  Create Your First Alert
                </Button>
              </Paper>
            ) : (
              <AlertsTable
                alerts={filteredAlerts.map(alert => ({
                  ...alert,
                  status: alert.status as AlertStatus,
                  alertType: alert.alertType as AlertType,
                }))}
                onEdit={handleEditAlert}
                onDelete={handleDeleteAlert}
              />
            )}
          </Stack>
        </TabPanel>

        <TabPanel value={currentTab} index={2}>
          <Stack spacing={3}>
            <Paper sx={{ p: 2, bgcolor: 'grey.50' }}>
              <Typography variant="h6" gutterBottom>
                ⚙️ Notification Preferences
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Customize how and when you receive stock alerts. Choose your preferred notification methods and frequency.
              </Typography>
            </Paper>
            <NotificationSettings onSave={handleSaveNotificationSettings} />
          </Stack>
        </TabPanel>

        <TabPanel value={currentTab} index={3}>
          <Stack spacing={3}>
            <Paper sx={{ p: 2, bgcolor: 'grey.50' }}>
              <Typography variant="h6" gutterBottom>
                👤 Your Profile
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Manage your account information and preferences for the Shares Alert Ghana platform.
              </Typography>
            </Paper>
            <UserProfile />
          </Stack>
        </TabPanel>
      </Stack>

      {/* Floating Action Button */}
      <Fab
        color="primary"
        sx={{ position: 'fixed', bottom: 16, right: 16 }}
        onClick={() => setAlertFormOpen(true)}
      >
        <NotificationAddIcon />
      </Fab>

      {/* Alert Form Dialog */}
      <AlertForm
        open={alertFormOpen}
        onClose={() => setAlertFormOpen(false)}
        onSubmit={handleCreateAlert}
        stocks={stocks.map(stock => ({
          symbol: stock.symbol,
          name: stock.name,
        }))}
      />

      {/* Snackbar for notifications */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={() => setSnackbar(prev => ({ ...prev, open: false }))}
      >
        <Alert
          onClose={() => setSnackbar(prev => ({ ...prev, open: false }))}
          severity={snackbar.severity}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
      </StyledContainer>
      <Footer />
    </Box>
  );
};

export default Dashboard;