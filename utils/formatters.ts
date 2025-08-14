import { AlertType, NotificationChannel } from '../types/enums';

// Formatting functions for stock alert system
export const formatPrice = (price: number): string => {
  return `GHS ${price.toFixed(2)}`;
};

export const formatPercentage = (value: number): string => {
  return `${value.toFixed(2)}%`;
};

export const formatNumber = (value: number): string => {
  return value.toLocaleString();
};

export const formatDate = (date: Date): string => {
  return date.toLocaleDateString('en-GB');
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