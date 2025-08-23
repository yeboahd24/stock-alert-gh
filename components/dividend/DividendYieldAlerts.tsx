import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Chip,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  CircularProgress,
  Stack,
  Divider,
} from '@mui/material';
import {
  Edit,
  Delete,
  Pause,
  PlayArrow,
  TrendingUp,
  Target,
  ChangeCircle,
  NotificationsActive,
} from '@mui/icons-material';
import { AlertType, AlertStatus } from '../../types/enums';
import { Alert as AlertModel } from '../../types/schema';
import { alertApi } from '../../src/services/api';
import { formatPercentage, formatAlertType, formatDate } from '../../utils/formatters';

interface DividendYieldAlertsProps {
  onEditAlert?: (alert: AlertModel) => void;
  onCreateAlert?: () => void;
}

const DividendYieldAlerts: React.FC<DividendYieldAlertsProps> = ({ onEditAlert, onCreateAlert }) => {
  const [alerts, setAlerts] = useState<AlertModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [alertToDelete, setAlertToDelete] = useState<AlertModel | null>(null);

  const fetchAlerts = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await alertApi.getAlerts();
      // Filter for dividend yield related alerts
      const dividendYieldAlerts = data.filter((alert: AlertModel) => 
        [
          AlertType.HIGH_DIVIDEND_YIELD,
          AlertType.TARGET_DIVIDEND_YIELD,
          AlertType.DIVIDEND_YIELD_CHANGE,
          AlertType.DIVIDEND_ANNOUNCEMENT
        ].includes(alert.alertType as AlertType)
      );
      setAlerts(dividendYieldAlerts);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch alerts');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAlerts();
  }, []);

  const handleDeleteAlert = async () => {
    if (!alertToDelete) return;

    try {
      await alertApi.deleteAlert(alertToDelete.id);
      setAlerts(alerts.filter(alert => alert.id !== alertToDelete.id));
      setDeleteDialogOpen(false);
      setAlertToDelete(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete alert');
    }
  };

  const handleToggleAlert = async (alert: AlertModel) => {
    try {
      const newStatus = alert.status === AlertStatus.ACTIVE ? AlertStatus.INACTIVE : AlertStatus.ACTIVE;
      await alertApi.updateAlert(alert.id, { status: newStatus });
      setAlerts(alerts.map(a => 
        a.id === alert.id ? { ...a, status: newStatus } : a
      ));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update alert');
    }
  };

  const getAlertIcon = (alertType: string) => {
    switch (alertType) {
      case AlertType.HIGH_DIVIDEND_YIELD:
        return <TrendingUp color="success" />;
      case AlertType.TARGET_DIVIDEND_YIELD:
        return <Target color="primary" />;
      case AlertType.DIVIDEND_YIELD_CHANGE:
        return <ChangeCircle color="warning" />;
      case AlertType.DIVIDEND_ANNOUNCEMENT:
        return <NotificationsActive color="info" />;
      default:
        return <NotificationsActive />;
    }
  };

  const getAlertDescription = (alert: AlertModel) => {
    switch (alert.alertType) {
      case AlertType.HIGH_DIVIDEND_YIELD:
        return `Alert when ${alert.stockSymbol || 'any stock'} yield ≥ ${formatPercentage(alert.thresholdYield || 0)}`;
      case AlertType.TARGET_DIVIDEND_YIELD:
        return `Alert when ${alert.stockSymbol} yield reaches ${formatPercentage(alert.targetYield || 0)}`;
      case AlertType.DIVIDEND_YIELD_CHANGE:
        return `Alert when ${alert.stockSymbol} yield changes by ±${formatPercentage(alert.yieldChangeThreshold || 0)}`;
      case AlertType.DIVIDEND_ANNOUNCEMENT:
        return `Alert for dividend announcements from ${alert.stockSymbol || 'any stock'}`;
      default:
        return 'Dividend alert';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case AlertStatus.ACTIVE:
        return 'success';
      case AlertStatus.TRIGGERED:
        return 'warning';
      case AlertStatus.INACTIVE:
        return 'default';
      default:
        return 'default';
    }
  };

  const activeAlerts = alerts.filter(alert => alert.status === AlertStatus.ACTIVE);
  const triggeredAlerts = alerts.filter(alert => alert.status === AlertStatus.TRIGGERED);
  const inactiveAlerts = alerts.filter(alert => alert.status === AlertStatus.INACTIVE);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
        <Typography variant="body1" sx={{ ml: 2 }}>
          Loading dividend alerts...
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">
          💰 Dividend Yield Alerts
        </Typography>
        <Button
          variant="contained"
          onClick={onCreateAlert}
          startIcon={<NotificationsActive />}
        >
          Create Alert
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Summary Cards */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Active Alerts
              </Typography>
              <Typography variant="h4" color="success.main">
                {activeAlerts.length}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Triggered Alerts
              </Typography>
              <Typography variant="h4" color="warning.main">
                {triggeredAlerts.length}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Alerts
              </Typography>
              <Typography variant="h4" color="primary">
                {alerts.length}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Alerts Table */}
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell><strong>Type</strong></TableCell>
              <TableCell><strong>Stock</strong></TableCell>
              <TableCell><strong>Condition</strong></TableCell>
              <TableCell><strong>Current Yield</strong></TableCell>
              <TableCell><strong>Status</strong></TableCell>
              <TableCell><strong>Created</strong></TableCell>
              <TableCell align="center"><strong>Actions</strong></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {alerts.map((alert) => (
              <TableRow key={alert.id} hover>
                <TableCell>
                  <Box display="flex" alignItems="center">
                    {getAlertIcon(alert.alertType)}
                    <Typography variant="body2" sx={{ ml: 1 }}>
                      {formatAlertType(alert.alertType as AlertType)}
                    </Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Box>
                    <Typography variant="subtitle2" fontWeight="bold">
                      {alert.stockSymbol || 'All Stocks'}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {alert.stockName || 'Market-wide'}
                    </Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Typography variant="body2">
                    {getAlertDescription(alert)}
                  </Typography>
                </TableCell>
                <TableCell>
                  {alert.currentYield !== undefined ? (
                    <Typography 
                      variant="subtitle2" 
                      color={alert.currentYield >= 3.0 ? 'success.main' : 'text.primary'}
                    >
                      {formatPercentage(alert.currentYield)}
                    </Typography>
                  ) : (
                    <Typography variant="body2" color="textSecondary">
                      N/A
                    </Typography>
                  )}
                </TableCell>
                <TableCell>
                  <Chip
                    label={alert.status}
                    color={getStatusColor(alert.status) as any}
                    size="small"
                    variant={alert.status === AlertStatus.ACTIVE ? 'filled' : 'outlined'}
                  />
                </TableCell>
                <TableCell>
                  <Typography variant="body2">
                    {formatDate(alert.createdAt)}
                  </Typography>
                </TableCell>
                <TableCell align="center">
                  <Stack direction="row" spacing={1} justifyContent="center">
                    <Tooltip title={alert.status === AlertStatus.ACTIVE ? 'Pause Alert' : 'Activate Alert'}>
                      <IconButton
                        size="small"
                        color={alert.status === AlertStatus.ACTIVE ? 'warning' : 'success'}
                        onClick={() => handleToggleAlert(alert)}
                      >
                        {alert.status === AlertStatus.ACTIVE ? <Pause /> : <PlayArrow />}
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Edit Alert">
                      <IconButton
                        size="small"
                        color="primary"
                        onClick={() => onEditAlert?.(alert)}
                      >
                        <Edit />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Delete Alert">
                      <IconButton
                        size="small"
                        color="error"
                        onClick={() => {
                          setAlertToDelete(alert);
                          setDeleteDialogOpen(true);
                        }}
                      >
                        <Delete />
                      </IconButton>
                    </Tooltip>
                  </Stack>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {alerts.length === 0 && (
        <Box textAlign="center" py={6}>
          <Typography variant="h6" color="textSecondary" gutterBottom>
            No dividend alerts yet
          </Typography>
          <Typography variant="body2" color="textSecondary" paragraph>
            Create your first dividend yield alert to start monitoring opportunities
          </Typography>
          <Button
            variant="contained"
            onClick={onCreateAlert}
            startIcon={<NotificationsActive />}
          >
            Create Your First Alert
          </Button>
        </Box>
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Delete Alert</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete this alert? This action cannot be undone.
          </Typography>
          {alertToDelete && (
            <Box mt={2} p={2} bgcolor="grey.50" borderRadius={1}>
              <Typography variant="subtitle2">
                {formatAlertType(alertToDelete.alertType as AlertType)}
              </Typography>
              <Typography variant="body2" color="textSecondary">
                {getAlertDescription(alertToDelete)}
              </Typography>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleDeleteAlert} color="error" variant="contained">
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default DividendYieldAlerts;