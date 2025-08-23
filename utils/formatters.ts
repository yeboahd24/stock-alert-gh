import { AlertType, NotificationChannel } from '../types/enums';

// Formatting functions for stock alert system
export const formatPrice = (price: number): string => {
  return `GHS ${price.toFixed(2)}`;
};

export const formatCurrency = (amount: number): string => {
  return `GH₵ ${amount.toFixed(2)}`;
};

export const formatPercentage = (value: number): string => {
  return `${value.toFixed(2)}%`;
};

export const formatNumber = (value: number): string => {
  return value.toLocaleString();
};

export const formatDate = (date: Date | string): string => {
  const dateObj = typeof date === 'string' ? new Date(date) : date;
  return dateObj.toLocaleDateString('en-GB', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

export const formatDateTime = (date: Date): string => {
  return date.toLocaleString('en-GB');
};

export const formatAlertType = (type: AlertType): string => {
  switch (type) {
    case AlertType.PRICE_THRESHOLD:
      return 'Price Threshold';
    case AlertType.IPO_ALERT:
      return 'IPO Alert';
    case AlertType.DIVIDEND_ANNOUNCEMENT:
      return 'Dividend Announcement';
    case AlertType.HIGH_DIVIDEND_YIELD:
      return 'High Dividend Yield';
    case AlertType.DIVIDEND_YIELD_CHANGE:
      return 'Dividend Yield Change';
    case AlertType.TARGET_DIVIDEND_YIELD:
      return 'Target Dividend Yield';
    default:
      return type;
  }
};

export const formatNotificationChannel = (channel: NotificationChannel): string => {
  switch (channel) {
    case NotificationChannel.SMS:
      return 'SMS';
    case NotificationChannel.WHATSAPP:
      return 'WhatsApp';
    case NotificationChannel.TELEGRAM:
      return 'Telegram';
    case NotificationChannel.EMAIL:
      return 'Email';
    case NotificationChannel.MOBILE_PUSH:
      return 'Mobile Push';
    default:
      return channel;
  }
};