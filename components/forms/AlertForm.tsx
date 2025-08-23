import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Stack,
  IconButton,
  Autocomplete,
  Typography,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import { AlertType } from '../../types/enums';
import { formatAlertType } from '../../utils/formatters';

interface AlertFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (alertData: AlertFormData) => void;
  stocks: Array<{ symbol: string; name: string }>;
}

export interface AlertFormData {
  stockSymbol: string;
  stockName: string;
  alertType: AlertType;
  thresholdPrice?: number;
  thresholdYield?: number;
  targetYield?: number;
  yieldChangeThreshold?: number;
}

const AlertForm: React.FC<AlertFormProps> = ({ open, onClose, onSubmit, stocks }) => {
  const [formData, setFormData] = useState<AlertFormData>({
    stockSymbol: '',
    stockName: '',
    alertType: AlertType.PRICE_THRESHOLD,
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.stockSymbol && formData.alertType !== AlertType.HIGH_DIVIDEND_YIELD) {
      newErrors.stockSymbol = 'Stock selection is required';
    }

    if (formData.alertType === AlertType.PRICE_THRESHOLD && !formData.thresholdPrice) {
      newErrors.thresholdPrice = 'Threshold price is required';
    }

    if (formData.alertType === AlertType.PRICE_THRESHOLD && formData.thresholdPrice && formData.thresholdPrice <= 0) {
      newErrors.thresholdPrice = 'Threshold price must be greater than 0';
    }

    if (formData.alertType === AlertType.HIGH_DIVIDEND_YIELD && !formData.thresholdYield) {
      newErrors.thresholdYield = 'Minimum dividend yield is required';
    }

    if (formData.alertType === AlertType.HIGH_DIVIDEND_YIELD && formData.thresholdYield && formData.thresholdYield <= 0) {
      newErrors.thresholdYield = 'Minimum dividend yield must be greater than 0';
    }

    if (formData.alertType === AlertType.TARGET_DIVIDEND_YIELD && !formData.targetYield) {
      newErrors.targetYield = 'Target dividend yield is required';
    }

    if (formData.alertType === AlertType.TARGET_DIVIDEND_YIELD && formData.targetYield && formData.targetYield <= 0) {
      newErrors.targetYield = 'Target dividend yield must be greater than 0';
    }

    if (formData.alertType === AlertType.DIVIDEND_YIELD_CHANGE && !formData.yieldChangeThreshold) {
      newErrors.yieldChangeThreshold = 'Yield change threshold is required';
    }

    if (formData.alertType === AlertType.DIVIDEND_YIELD_CHANGE && formData.yieldChangeThreshold && formData.yieldChangeThreshold <= 0) {
      newErrors.yieldChangeThreshold = 'Yield change threshold must be greater than 0';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = () => {
    if (validateForm()) {
      onSubmit(formData);
      handleClose();
    }
  };

  const handleClose = () => {
    setFormData({
      stockSymbol: '',
      stockName: '',
      alertType: AlertType.PRICE_THRESHOLD,
    });
    setErrors({});
    onClose();
  };

  const handleStockChange = (value: { symbol: string; name: string } | null) => {
    if (value) {
      setFormData(prev => ({
        ...prev,
        stockSymbol: value.symbol,
        stockName: value.name,
      }));
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Stack direction="row" justifyContent="space-between" alignItems="center">
          Create New Alert
          <IconButton onClick={handleClose} size="small">
            <CloseIcon />
          </IconButton>
        </Stack>
      </DialogTitle>

      <DialogContent>
        <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
          Set up alerts to get notified about stock price changes, new IPO listings, or dividend announcements.
        </Typography>
        <Stack spacing={3} sx={{ mt: 1 }}>
          <Autocomplete
            options={stocks}
            getOptionLabel={(option) => `${option.symbol} - ${option.name}`}
            onChange={(_, value) => handleStockChange(value)}
            renderInput={(params) => (
              <TextField
                {...params}
                label="Select Stock"
                error={!!errors.stockSymbol}
                helperText={errors.stockSymbol}
                required
              />
            )}
          />

          <FormControl fullWidth>
            <InputLabel>Alert Type</InputLabel>
            <Select
              value={formData.alertType}
              label="Alert Type"
              onChange={(e) => setFormData(prev => ({ ...prev, alertType: e.target.value as AlertType }))}
            >
              {Object.values(AlertType).map((type) => (
                <MenuItem key={type} value={type}>
                  {formatAlertType(type)}
                </MenuItem>
              ))}
            </Select>
            <Typography variant=\"caption\" color=\"textSecondary\" sx={{ mt: 1 }}>
              • Price Alert: Get notified when stock reaches your target price<br/>
              • IPO Alert: Get notified about new company listings<br/>
              • Dividend Alert: Get notified about dividend announcements<br/>
              • High Dividend Yield: Get notified when stocks exceed yield thresholds<br/>
              • Target Dividend Yield: Get notified when a stock reaches your target yield<br/>
              • Dividend Yield Change: Get notified when yields change significantly
            </Typography>
          </FormControl>

          {formData.alertType === AlertType.PRICE_THRESHOLD && (
            <TextField
              label="Threshold Price (GHS)"
              type="number"
              value={formData.thresholdPrice || ''}
              onChange={(e) => setFormData(prev => ({ ...prev, thresholdPrice: parseFloat(e.target.value) }))}
              error={!!errors.thresholdPrice}
              helperText={errors.thresholdPrice || "You'll be notified when the stock price reaches or exceeds this amount"}
              required
              inputProps={{ min: 0, step: 0.01 }}
            />
          )}
          
          {formData.alertType === AlertType.IPO_ALERT && (
            <Typography variant="body2" color="textSecondary" sx={{ p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
              📈 You'll receive notifications when new companies are announced for listing on the Ghana Stock Exchange.
            </Typography>
          )}
          
          {formData.alertType === AlertType.DIVIDEND_ANNOUNCEMENT && (
            <Typography variant="body2" color="textSecondary" sx={{ p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
              💰 You'll receive notifications when companies announce dividend payments for this stock.
            </Typography>
          )}

          {formData.alertType === AlertType.HIGH_DIVIDEND_YIELD && (
            <>
              <TextField
                label=\"Minimum Dividend Yield (%)\"\
                type=\"number\"
                value={formData.thresholdYield || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, thresholdYield: parseFloat(e.target.value) }))}
                error={!!errors.thresholdYield}
                helperText={errors.thresholdYield || \"You'll be notified when any stock (or selected stock) reaches this dividend yield\"}
                required
                inputProps={{ min: 0, step: 0.1, max: 100 }}
              />
              <Typography variant=\"body2\" color=\"textSecondary\" sx={{ p: 2, bgcolor: 'green.50', borderRadius: 1 }}>
                📈 Monitor the market for high dividend yield opportunities. Leave stock selection empty to monitor all stocks.
              </Typography>
            </>
          )}

          {formData.alertType === AlertType.TARGET_DIVIDEND_YIELD && (
            <>
              <TextField
                label=\"Target Dividend Yield (%)\"\
                type=\"number\"
                value={formData.targetYield || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, targetYield: parseFloat(e.target.value) }))}
                error={!!errors.targetYield}
                helperText={errors.targetYield || \"You'll be notified when this specific stock reaches your target dividend yield\"}
                required
                inputProps={{ min: 0, step: 0.1, max: 100 }}
              />
              <Typography variant=\"body2\" color=\"textSecondary\" sx={{ p: 2, bgcolor: 'blue.50', borderRadius: 1 }}>
                🎯 Perfect for timing your entry into dividend-paying stocks. Set your ideal yield and wait for the right moment.
              </Typography>
            </>
          )}

          {formData.alertType === AlertType.DIVIDEND_YIELD_CHANGE && (
            <>
              <TextField
                label=\"Yield Change Threshold (%)\"\
                type=\"number\"
                value={formData.yieldChangeThreshold || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, yieldChangeThreshold: parseFloat(e.target.value) }))}
                error={!!errors.yieldChangeThreshold}
                helperText={errors.yieldChangeThreshold || \"You'll be notified when the dividend yield changes by this amount or more\"}
                required
                inputProps={{ min: 0, step: 0.1, max: 10 }}
              />
              <Typography variant=\"body2\" color=\"textSecondary\" sx={{ p: 2, bgcolor: 'orange.50', borderRadius: 1 }}>
                📊 Track significant dividend yield movements that might indicate opportunities or risks in your holdings.
              </Typography>
            </>
          )}
        </Stack>
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose} color="inherit">
          Cancel
        </Button>
        <Button onClick={handleSubmit} variant="contained">
          Create Alert
        </Button>
        <Typography variant="caption" color="textSecondary" align="center" sx={{ mt: 1 }}>
          You'll receive email notifications when your alert conditions are met.
        </Typography>
      </DialogActions>
    </Dialog>
  );
};

export default AlertForm;