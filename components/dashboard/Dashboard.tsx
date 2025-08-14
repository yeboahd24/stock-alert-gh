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
import { mockStore } from '../../data/sharesAlertMockData';
import { AlertType, AlertStatus } from '../../types/enums';
import { stockApi, alertApi, Stock, Alert as ApiAlert } from '../../src/services/api';

const StyledContainer = styled(Container)(({ theme }) => ({
  paddingTop: theme.spacing(3),
  paddingBottom: theme.spacing(3),
}));

const TabPanel = ({ children, value, index }: { children: React.ReactNode; value: number; index: number }) => (
  <div hidden={value !== index}>
    {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
  </div>
);

const Dashboard: React.FC = () => {
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

  // Load data from API
  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        const [stocksData, alertsData] = await Promise.all([
          stockApi.getAllStocks(),
          alertApi.getAllAlerts('user-123') // Using mock user ID
        ]);
        setStocks(stocksData);
        setAlerts(alertsData);
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
    <StyledContainer maxWidth="lg">
      <Stack spacing={3}>
        {/* Header */}
        <Stack direction="row" justifyContent="space-between" alignItems="center">
          <Stack>
            <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
              Shares Alert Ghana
            </Typography>
            <Typography variant="subtitle1" color="text.secondary">
              Welcome back, {mockStore.user.name}
            </Typography>
          </Stack>
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
            {/* Stock Cards Grid */}
            <Stack spacing={2}>
              <Typography variant="h6" component="h2">
                Your Watchlist
              </Typography>
              <Stack direction="row" spacing={2} sx={{ overflowX: 'auto', pb: 1 }}>
                {loading ? (
                  <Typography>Loading stocks...</Typography>
                ) : (
                  stocks.map((stock) => (
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
                )}
              </Stack>
            </Stack>

            {/* Chart */}
            <Stack spacing={2}>
              <Stack direction="row" justifyContent="space-between" alignItems="center">
                <Typography variant="h6" component="h2">
                  Price Chart
                </Typography>
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
            <Stack direction="row" justifyContent="space-between" alignItems="center">
              <Typography variant="h6" component="h2">
                My Alerts
              </Typography>
              <Button
                variant="contained"
                startIcon={<NotificationAddIcon />}
                onClick={() => setAlertFormOpen(true)}
              >
                Create Alert
              </Button>
            </Stack>
            <AlertsTable
              alerts={alerts.map(alert => ({
                ...alert,
                status: alert.status as AlertStatus,
                alertType: alert.alertType as AlertType,
              }))}
              onEdit={handleEditAlert}
              onDelete={handleDeleteAlert}
            />
          </Stack>
        </TabPanel>

        <TabPanel value={currentTab} index={2}>
          <Stack spacing={3}>
            <Typography variant="h6" component="h2">
              Notification Settings
            </Typography>
            <NotificationSettings onSave={handleSaveNotificationSettings} />
          </Stack>
        </TabPanel>

        <TabPanel value={currentTab} index={3}>
          <Stack spacing={3}>
            <Typography variant="h6" component="h2">
              Profile Settings
            </Typography>
            <Paper sx={{ p: 3 }}>
              <Typography variant="body1">
                Profile management features coming soon...
              </Typography>
            </Paper>
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
  );
};

export default Dashboard;