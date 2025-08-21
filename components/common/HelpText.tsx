import React from 'react';
import { Box, Typography, Paper } from '@mui/material';

interface HelpTextProps {
  title: string;
  description: string;
  icon?: string;
}

const HelpText: React.FC<HelpTextProps> = ({ title, description, icon = "ℹ️" }) => {
  return (
    <Paper elevation={0} sx={{ p: 2, bgcolor: 'grey.50', mb: 2 }}>
      <Box display="flex" alignItems="flex-start" gap={1}>
        <Typography variant="body1">{icon}</Typography>
        <Box>
          <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
            {title}
          </Typography>
          <Typography variant="body2" color="textSecondary">
            {description}
          </Typography>
        </Box>
      </Box>
    </Paper>
  );
};

export default HelpText;