// Alert types for the stock alert system
export enum AlertType {
  PRICE_THRESHOLD = 'price_threshold',
  IPO_ALERT = 'ipo_alert',
  DIVIDEND_ANNOUNCEMENT = 'dividend_announcement',
  HIGH_DIVIDEND_YIELD = 'high_dividend_yield',
  DIVIDEND_YIELD_CHANGE = 'dividend_yield_change',
  TARGET_DIVIDEND_YIELD = 'target_dividend_yield'
}

// Notification channels available
export enum NotificationChannel {
  SMS = 'sms',
  WHATSAPP = 'whatsapp',
  TELEGRAM = 'telegram',
  EMAIL = 'email',
  MOBILE_PUSH = 'mobile_push'
}

// Alert status
export enum AlertStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  TRIGGERED = 'triggered'
}

// Language options
export enum Language {
  ENGLISH = 'english',
  TWI = 'twi'
}