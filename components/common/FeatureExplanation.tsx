import React from 'react';
import { Box, Typography, Card, CardContent, Grid } from '@mui/material';

const FeatureExplanation: React.FC = () => {
  const features = [
    {
      icon: "📈",
      title: "Price Alerts",
      description: "Get notified when stocks reach your target price. Perfect for buying opportunities or stop-loss triggers."
    },
    {
      icon: "🎆",
      title: "IPO Alerts", 
      description: "Be the first to know when new companies list on the Ghana Stock Exchange. Don't miss investment opportunities."
    },
    {
      icon: "💰",
      title: "Dividend Alerts",
      description: "Stay informed about dividend announcements and payment dates for your favorite stocks."
    }
  ];

  return (
    <Box sx={{ mb: 4 }}>
      <Typography variant="h5" gutterBottom align="center">
        How Shares Alert Ghana Works
      </Typography>
      <Typography variant="body1" color="textSecondary" align="center" sx={{ mb: 3 }}>
        Set up personalized alerts for Ghana Stock Exchange stocks and never miss important market events.
      </Typography>
      
      <Grid container spacing={3}>
        {features.map((feature, index) => (
          <Grid item xs={12} md={4} key={index}>
            <Card sx={{ height: '100%', textAlign: 'center' }}>
              <CardContent>
                <Typography variant="h3" sx={{ mb: 2 }}>
                  {feature.icon}
                </Typography>
                <Typography variant="h6" gutterBottom>
                  {feature.title}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  {feature.description}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default FeatureExplanation;