import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Stack,
  Switch,
  FormControlLabel,
  Divider,
  TextField,
  Button,
} from '@mui/material';
import { NotificationChannel } from '../../types/enums';
import { formatNotificationChannel } from '../../utils/formatters';

interface NotificationSettingsProps {
  initialSettings?: NotificationSettings;
  onSave: (settings: NotificationSettings) => void;
}

export interface NotificationSettings {
  channels: Record<NotificationChannel, boolean>;
  phoneNumber?: string;
  telegramUsername?: string;
  whatsappNumber?: string;
}

const NotificationSettings: React.FC<NotificationSettingsProps> = ({ 
  initialSettings,
  onSave 
}) => {
  const [settings, setSettings] = useState<NotificationSettings>(
    initialSettings || {
      channels: {
        [NotificationChannel.EMAIL]: true,
        [NotificationChannel.SMS]: false,
        [NotificationChannel.WHATSAPP]: false,
        [NotificationChannel.TELEGRAM]: false,
        [NotificationChannel.MOBILE_PUSH]: true,
      },
    }
  );

  const [phoneError, setPhoneError] = useState('');

  const validatePhone = (phone: string): boolean => {
    const phoneRegex = /^\+233\d{9}$/;
    if (phone && !phoneRegex.test(phone)) {
      setPhoneError('Phone number must be in format +233XXXXXXXXX');
      return false;
    }
    setPhoneError('');
    return true;
  };

  const handleChannelToggle = (channel: NotificationChannel) => {
    setSettings(prev => ({
      ...prev,
      channels: {
        ...prev.channels,
        [channel]: !prev.channels[channel],
      },
    }));
  };

  const handleSave = () => {
    if (settings.phoneNumber && !validatePhone(settings.phoneNumber)) {
      return;
    }
    onSave(settings);
  };

  return (
    <Card>
      <CardContent>
        <Stack spacing={3}>
          <Typography variant="h6" component="h3">
            Notification Settings
          </Typography>

          <Stack spacing={2}>
            <Typography variant="subtitle2" color="text.secondary">
              Notification Channels
            </Typography>
            
            {Object.values(NotificationChannel).map((channel) => (
              <FormControlLabel
                key={channel}
                control={
                  <Switch
                    checked={settings.channels[channel]}
                    onChange={() => handleChannelToggle(channel)}
                  />
                }
                label={formatNotificationChannel(channel)}
              />
            ))}
          </Stack>

          <Divider />

          <Stack spacing={2}>
            <Typography variant="subtitle2" color="text.secondary">
              Contact Information
            </Typography>

            {(settings.channels[NotificationChannel.SMS] || settings.channels[NotificationChannel.WHATSAPP]) && (
              <TextField
                label="Phone Number"
                placeholder="+233244123456"
                value={settings.phoneNumber || ''}
                onChange={(e) => {
                  setSettings(prev => ({ ...prev, phoneNumber: e.target.value }));
                  if (e.target.value) validatePhone(e.target.value);
                }}
                error={!!phoneError}
                helperText={phoneError || 'Required for SMS and WhatsApp notifications'}
                required
              />
            )}

            {settings.channels[NotificationChannel.TELEGRAM] && (
              <TextField
                label="Telegram Username"
                placeholder="@username"
                value={settings.telegramUsername || ''}
                onChange={(e) => setSettings(prev => ({ ...prev, telegramUsername: e.target.value }))}
                helperText="Your Telegram username (including @)"
              />
            )}
          </Stack>

          <Button
            variant="contained"
            onClick={handleSave}
            sx={{ alignSelf: 'flex-start' }}
          >
            Save Settings
          </Button>
        </Stack>
      </CardContent>
    </Card>
  );
};

export default NotificationSettings;