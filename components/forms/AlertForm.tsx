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

    if (!formData.stockSymbol) {
      newErrors.stockSymbol = 'Stock selection is required';
    }

    if (formData.alertType === AlertType.PRICE_THRESHOLD && !formData.thresholdPrice) {
      newErrors.thresholdPrice = 'Threshold price is required';
    }

    if (formData.alertType === AlertType.PRICE_THRESHOLD && formData.thresholdPrice && formData.thresholdPrice <= 0) {
      newErrors.thresholdPrice = 'Threshold price must be greater than 0';
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
          </FormControl>

          {formData.alertType === AlertType.PRICE_THRESHOLD && (
            <TextField
              label="Threshold Price (GHS)"
              type="number"
              value={formData.thresholdPrice || ''}
              onChange={(e) => setFormData(prev => ({ ...prev, thresholdPrice: parseFloat(e.target.value) }))}
              error={!!errors.thresholdPrice}
              helperText={errors.thresholdPrice}
              required
              inputProps={{ min: 0, step: 0.01 }}
            />
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
      </DialogActions>
    </Dialog>
  );
};

export default AlertForm;