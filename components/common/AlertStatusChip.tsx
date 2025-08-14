import React from 'react';
import { Chip } from '@mui/material';
import { AlertStatus } from '../../types/enums';

interface AlertStatusChipProps {
  status: AlertStatus;
  size?: 'small' | 'medium';
}

const AlertStatusChip: React.FC<AlertStatusChipProps> = ({ status, size = 'small' }) => {
  const getChipProps = (status: AlertStatus) => {
    switch (status) {
      case AlertStatus.ACTIVE:
        return {
          label: 'Active',
          color: 'success' as const,
          variant: 'filled' as const,
        };
      case AlertStatus.INACTIVE:
        return {
          label: 'Inactive',
          color: 'default' as const,
          variant: 'outlined' as const,
        };
      case AlertStatus.TRIGGERED:
        return {
          label: 'Triggered',
          color: 'warning' as const,
          variant: 'filled' as const,
        };
      default:
        return {
          label: status,
          color: 'default' as const,
          variant: 'outlined' as const,
        };
    }
  };

  const chipProps = getChipProps(status);

  return (
    <Chip
      {...chipProps}
      size={size}
    />
  );
};

export default AlertStatusChip;