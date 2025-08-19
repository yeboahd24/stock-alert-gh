import React from 'react';
import { Box, Typography, Link } from '@mui/material';

const Footer: React.FC = () => {
  return (
    <Box
      component="footer"
      sx={{
        py: 3,
        px: 2,
        mt: 'auto',
        backgroundColor: 'grey.100',
        textAlign: 'center',
      }}
    >
      <Typography variant="body2" color="text.secondary">
        © 2025 Dominic Kofi Yeboah. All rights reserved.
      </Typography>
      <Typography variant="body2" color="text.secondary">
        Contact: <Link href="mailto:yeboahd24@gmail.com">yeboahd24@gmail.com</Link>
      </Typography>
    </Box>
  );
};

export default Footer;