import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import AlertForm from '../components/forms/AlertForm';
import DividendYieldDashboard from '../components/dividend/DividendYieldDashboard';
import DividendYieldAlerts from '../components/dividend/DividendYieldAlerts';
import { AlertType } from '../types/enums';
import { dividendApi, alertApi } from '../src/services/api';

// Mock the API modules
vi.mock('../src/services/api', () => ({
  dividendApi: {
    getGSEDividendStocks: vi.fn(),
    getDividendStockBySymbol: vi.fn(),
    getHighDividendYieldStocks: vi.fn(),
    getAllDividends: vi.fn(),
    getUpcomingDividends: vi.fn(),
  },
  alertApi: {
    getAlerts: vi.fn(),
    createAlert: vi.fn(),
    updateAlert: vi.fn(),
    deleteAlert: vi.fn(),
  },
}));

// Mock data
const mockGSEDividendResponse = {
  success: true,
  data: {
    timestamp: '2025-08-23T12:30:46.725685629Z',
    source: 'https://simplywall.st/stocks/gh/dividend-yield-high',
    count: 3,
    stocks: [
      {
        symbol: 'GCB',
        name: 'GCB Bank',
        dividend_yield: 2.5,
        price: 'GH₵9.85',
        market_cap: 'GH₵2.5b',
        country: 'Ghana',
        exchange: 'GSE',
        sector: 'Banks',
        url: 'https://simplywall.st/stocks/gh/banks/ghse-gcb/gcb-bank-shares',
      },
      {
        symbol: 'GOIL',
        name: 'GOIL Company',
        dividend_yield: 3.7,
        price: 'GH₵2.50',
        market_cap: 'GH₵1.2b',
        country: 'Ghana',
        exchange: 'GSE',
        sector: 'Energy',
        url: 'https://simplywall.st/stocks/gh/energy/ghse-goil/goil-company-shares',
      },
      {
        symbol: 'RBGH',
        name: 'Republic Bank Ghana',
        dividend_yield: 3.8,
        price: 'GH₵5.20',
        market_cap: 'GH₵800m',
        country: 'Ghana',
        exchange: 'GSE',
        sector: 'Banks',
        url: 'https://simplywall.st/stocks/gh/banks/ghse-rbgh/republic-bank-ghana-shares',
      },
    ],
  },
};

const mockAlerts = [
  {
    id: 'alert-1',
    userId: 'user-1',
    stockSymbol: 'GCB',
    stockName: 'GCB Bank',
    alertType: 'high_dividend_yield',
    thresholdYield: 3.0,
    currentYield: 2.5,
    status: 'active',
    createdAt: '2025-01-01T00:00:00Z',
    updatedAt: '2025-01-01T00:00:00Z',
  },
  {
    id: 'alert-2',
    userId: 'user-1',
    stockSymbol: 'GOIL',
    stockName: 'GOIL Company',
    alertType: 'target_dividend_yield',
    targetYield: 4.0,
    currentYield: 3.7,
    status: 'active',
    createdAt: '2025-01-01T00:00:00Z',
    updatedAt: '2025-01-01T00:00:00Z',
  },
];

const mockStocks = [
  { symbol: 'GCB', name: 'GCB Bank' },
  { symbol: 'GOIL', name: 'GOIL Company' },
  { symbol: 'RBGH', name: 'Republic Bank Ghana' },
];

describe('Enhanced Dividend Alerts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  describe('AlertForm - Dividend Yield Alerts', () => {
    const mockOnSubmit = vi.fn();
    const mockOnClose = vi.fn();

    beforeEach(() => {
      mockOnSubmit.mockClear();
      mockOnClose.mockClear();
    });

    it('should render high dividend yield alert form', async () => {
      const user = userEvent.setup();
      
      render(
        <AlertForm
          open={true}
          onClose={mockOnClose}
          onSubmit={mockOnSubmit}
          stocks={mockStocks}
        />
      );

      // Select alert type
      const alertTypeSelect = screen.getByLabelText('Alert Type');
      await user.click(alertTypeSelect);
      
      const highYieldOption = screen.getByText('High Dividend Yield');
      await user.click(highYieldOption);

      // Check if threshold yield field appears
      expect(screen.getByLabelText(/Minimum Dividend Yield/)).toBeInTheDocument();
      expect(screen.getByText(/Monitor the market for high dividend yield opportunities/)).toBeInTheDocument();
    });

    it('should validate high dividend yield alert form', async () => {
      const user = userEvent.setup();
      
      render(
        <AlertForm
          open={true}
          onClose={mockOnClose}
          onSubmit={mockOnSubmit}
          stocks={mockStocks}
        />
      );

      // Select high dividend yield alert type
      const alertTypeSelect = screen.getByLabelText('Alert Type');
      await user.click(alertTypeSelect);
      await user.click(screen.getByText('High Dividend Yield'));

      // Try to submit without required fields
      const submitButton = screen.getByText('Create Alert');
      await user.click(submitButton);

      // Should show validation error
      await waitFor(() => {
        expect(screen.getByText('Minimum dividend yield is required')).toBeInTheDocument();
      });

      expect(mockOnSubmit).not.toHaveBeenCalled();
    });

    it('should submit valid high dividend yield alert', async () => {
      const user = userEvent.setup();
      
      render(
        <AlertForm
          open={true}
          onClose={mockOnClose}
          onSubmit={mockOnSubmit}
          stocks={mockStocks}
        />
      );

      // Select alert type
      const alertTypeSelect = screen.getByLabelText('Alert Type');
      await user.click(alertTypeSelect);
      await user.click(screen.getByText('High Dividend Yield'));

      // Fill in threshold yield
      const thresholdYieldInput = screen.getByLabelText(/Minimum Dividend Yield/);
      await user.type(thresholdYieldInput, '3.5');

      // Submit form
      const submitButton = screen.getByText('Create Alert');
      await user.click(submitButton);

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith({
          stockSymbol: '',
          stockName: '',
          alertType: AlertType.HIGH_DIVIDEND_YIELD,
          thresholdYield: 3.5,
        });
      });
    });

    it('should render target dividend yield alert form', async () => {
      const user = userEvent.setup();
      
      render(
        <AlertForm
          open={true}
          onClose={mockOnClose}
          onSubmit={mockOnSubmit}
          stocks={mockStocks}
        />
      );

      // Select stock first
      const stockSelect = screen.getByLabelText('Select Stock');
      await user.click(stockSelect);
      await user.click(screen.getByText('GCB - GCB Bank'));

      // Select alert type
      const alertTypeSelect = screen.getByLabelText('Alert Type');
      await user.click(alertTypeSelect);
      await user.click(screen.getByText('Target Dividend Yield'));

      // Check if target yield field appears
      expect(screen.getByLabelText(/Target Dividend Yield/)).toBeInTheDocument();
      expect(screen.getByText(/Perfect for timing your entry/)).toBeInTheDocument();
    });

    it('should render dividend yield change alert form', async () => {
      const user = userEvent.setup();
      
      render(
        <AlertForm
          open={true}
          onClose={mockOnClose}
          onSubmit={mockOnSubmit}
          stocks={mockStocks}
        />
      );

      // Select stock first
      const stockSelect = screen.getByLabelText('Select Stock');
      await user.click(stockSelect);
      await user.click(screen.getByText('GOIL - GOIL Company'));

      // Select alert type
      const alertTypeSelect = screen.getByLabelText('Alert Type');
      await user.click(alertTypeSelect);
      await user.click(screen.getByText('Dividend Yield Change'));

      // Check if yield change threshold field appears
      expect(screen.getByLabelText(/Yield Change Threshold/)).toBeInTheDocument();
      expect(screen.getByText(/Track significant dividend yield movements/)).toBeInTheDocument();
    });
  });

  describe('DividendYieldDashboard', () => {
    const mockOnCreateAlert = vi.fn();

    beforeEach(() => {
      mockOnCreateAlert.mockClear();
      (dividendApi.getGSEDividendStocks as any).mockResolvedValue(mockGSEDividendResponse);
    });

    it('should render dividend yield dashboard', async () => {
      render(<DividendYieldDashboard onCreateAlert={mockOnCreateAlert} />);

      // Should show loading initially
      expect(screen.getByText('Loading dividend yield data...')).toBeInTheDocument();

      // Wait for data to load
      await waitFor(() => {
        expect(screen.getByText('📈 Dividend Yield Dashboard')).toBeInTheDocument();
      });

      // Should display summary cards
      expect(screen.getByText('Total Stocks')).toBeInTheDocument();
      expect(screen.getByText('High Yield Stocks (≥3%)')).toBeInTheDocument();
      expect(screen.getByText('Average Yield')).toBeInTheDocument();

      // Should display stock data
      expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      expect(screen.getByText('GOIL Company')).toBeInTheDocument();
      expect(screen.getByText('Republic Bank Ghana')).toBeInTheDocument();
    });

    it('should filter stocks by search term', async () => {
      const user = userEvent.setup();
      
      render(<DividendYieldDashboard onCreateAlert={mockOnCreateAlert} />);

      await waitFor(() => {
        expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      });

      // Search for specific stock
      const searchInput = screen.getByPlaceholderText('Search stocks, symbols, or sectors...');
      await user.type(searchInput, 'GCB');

      // Should only show GCB
      expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      expect(screen.queryByText('GOIL Company')).not.toBeInTheDocument();
    });

    it('should filter stocks by minimum yield', async () => {
      const user = userEvent.setup();
      
      render(<DividendYieldDashboard onCreateAlert={mockOnCreateAlert} />);

      await waitFor(() => {
        expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      });

      // Set minimum yield filter
      const minYieldInput = screen.getByLabelText('Min Yield (%)');
      await user.type(minYieldInput, '3.5');

      // Should only show stocks with yield >= 3.5%
      expect(screen.getByText('GOIL Company')).toBeInTheDocument(); // 3.7%
      expect(screen.getByText('Republic Bank Ghana')).toBeInTheDocument(); // 3.8%
      expect(screen.queryByText('GCB Bank')).not.toBeInTheDocument(); // 2.5%
    });

    it('should handle create alert action', async () => {
      const user = userEvent.setup();
      
      render(<DividendYieldDashboard onCreateAlert={mockOnCreateAlert} />);

      await waitFor(() => {
        expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      });

      // Click create alert button for a stock
      const alertButtons = screen.getAllByTitle('Create High Yield Alert');
      await user.click(alertButtons[0]);

      expect(mockOnCreateAlert).toHaveBeenCalledWith({
        stockSymbol: 'GCB',
        stockName: 'GCB Bank',
        alertType: 'high_dividend_yield',
        thresholdYield: 2.5,
      });
    });
  });

  describe('DividendYieldAlerts', () => {
    const mockOnEditAlert = vi.fn();
    const mockOnCreateAlert = vi.fn();

    beforeEach(() => {
      mockOnEditAlert.mockClear();
      mockOnCreateAlert.mockClear();
      (alertApi.getAlerts as any).mockResolvedValue(mockAlerts);
    });

    it('should render dividend yield alerts', async () => {
      render(
        <DividendYieldAlerts
          onEditAlert={mockOnEditAlert}
          onCreateAlert={mockOnCreateAlert}
        />
      );

      // Should show loading initially
      expect(screen.getByText('Loading dividend alerts...')).toBeInTheDocument();

      // Wait for data to load
      await waitFor(() => {
        expect(screen.getByText('💰 Dividend Yield Alerts')).toBeInTheDocument();
      });

      // Should display summary cards
      expect(screen.getByText('Active Alerts')).toBeInTheDocument();
      expect(screen.getByText('Triggered Alerts')).toBeInTheDocument();
      expect(screen.getByText('Total Alerts')).toBeInTheDocument();

      // Should display alert data
      expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      expect(screen.getByText('GOIL Company')).toBeInTheDocument();
    });

    it('should handle alert deletion', async () => {
      const user = userEvent.setup();
      (alertApi.deleteAlert as any).mockResolvedValue({});
      
      render(
        <DividendYieldAlerts
          onEditAlert={mockOnEditAlert}
          onCreateAlert={mockOnCreateAlert}
        />
      );

      await waitFor(() => {
        expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButtons = screen.getAllByTitle('Delete Alert');
      await user.click(deleteButtons[0]);

      // Should show confirmation dialog
      expect(screen.getByText('Delete Alert')).toBeInTheDocument();
      expect(screen.getByText('Are you sure you want to delete this alert?')).toBeInTheDocument();

      // Confirm deletion
      const confirmButton = screen.getByRole('button', { name: 'Delete' });
      await user.click(confirmButton);

      await waitFor(() => {
        expect(alertApi.deleteAlert).toHaveBeenCalledWith('alert-1');
      });
    });

    it('should handle alert toggle', async () => {
      const user = userEvent.setup();
      (alertApi.updateAlert as any).mockResolvedValue({});
      
      render(
        <DividendYieldAlerts
          onEditAlert={mockOnEditAlert}
          onCreateAlert={mockOnCreateAlert}
        />
      );

      await waitFor(() => {
        expect(screen.getByText('GCB Bank')).toBeInTheDocument();
      });

      // Click pause button
      const pauseButtons = screen.getAllByTitle('Pause Alert');
      await user.click(pauseButtons[0]);

      await waitFor(() => {
        expect(alertApi.updateAlert).toHaveBeenCalledWith('alert-1', { status: 'inactive' });
      });
    });

    it('should show empty state when no alerts', async () => {
      (alertApi.getAlerts as any).mockResolvedValue([]);
      
      render(
        <DividendYieldAlerts
          onEditAlert={mockOnEditAlert}
          onCreateAlert={mockOnCreateAlert}
        />
      );

      await waitFor(() => {
        expect(screen.getByText('No dividend alerts yet')).toBeInTheDocument();
        expect(screen.getByText('Create your first dividend yield alert')).toBeInTheDocument();
      });
    });
  });

  describe('API Integration', () => {
    it('should handle API errors gracefully', async () => {
      const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});
      (dividendApi.getGSEDividendStocks as any).mockRejectedValue(new Error('API Error'));
      
      render(<DividendYieldDashboard />);

      await waitFor(() => {
        expect(screen.getByText('API Error')).toBeInTheDocument();
      });

      consoleError.mockRestore();
    });

    it('should retry failed requests', async () => {
      const user = userEvent.setup();
      (dividendApi.getGSEDividendStocks as any)
        .mockRejectedValueOnce(new Error('Network Error'))
        .mockResolvedValueOnce(mockGSEDividendResponse);
      
      render(<DividendYieldDashboard />);

      await waitFor(() => {
        expect(screen.getByText('Network Error')).toBeInTheDocument();
      });

      // Click retry button
      const retryButton = screen.getByText('Retry');
      await user.click(retryButton);

      await waitFor(() => {
        expect(screen.getByText('📈 Dividend Yield Dashboard')).toBeInTheDocument();
      });
    });
  });

  describe('Data Formatting', () => {
    it('should format percentages correctly', () => {
      const { formatPercentage } = require('../utils/formatters');
      
      expect(formatPercentage(2.5)).toBe('2.50%');
      expect(formatPercentage(0)).toBe('0.00%');
      expect(formatPercentage(10.123)).toBe('10.12%');
    });

    it('should format alert types correctly', () => {
      const { formatAlertType } = require('../utils/formatters');
      
      expect(formatAlertType(AlertType.HIGH_DIVIDEND_YIELD)).toBe('High Dividend Yield');
      expect(formatAlertType(AlertType.TARGET_DIVIDEND_YIELD)).toBe('Target Dividend Yield');
      expect(formatAlertType(AlertType.DIVIDEND_YIELD_CHANGE)).toBe('Dividend Yield Change');
    });
  });
});