import React from 'react';
import {
  Paper,
  Typography,
  Avatar,
  Box,
  Stack,
  Chip,
  Divider,
  Grid,
  Card,
  CardContent,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import { useAuth } from '../../src/contexts/AuthContext';
import EmailIcon from '@mui/icons-material/Email';
import PersonIcon from '@mui/icons-material/Person';
import VerifiedIcon from '@mui/icons-material/Verified';
import CalendarTodayIcon from '@mui/icons-material/CalendarToday';

const ProfilePaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(3),
  marginBottom: theme.spacing(2),
}));

const ProfileHeader = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(2),
  marginBottom: theme.spacing(3),
}));

const UserProfile: React.FC = () => {
  const { user } = useAuth();

  if (!user) {
    return (
      <ProfilePaper>
        <Typography variant="h6" color="text.secondary">
          No user data available
        </Typography>
      </ProfilePaper>
    );
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <Stack spacing={3}>
      {/* Profile Header */}
      <ProfilePaper elevation={2}>
        <ProfileHeader>
          <Avatar
            src={user.picture}
            alt={user.name}
            sx={{ width: 80, height: 80 }}
          >
            {user.name.charAt(0).toUpperCase()}
          </Avatar>
          <Box>
            <Typography variant="h4" gutterBottom>
              {user.name}
            </Typography>
            <Stack direction="row" spacing={1} alignItems="center">
              <EmailIcon color="action" fontSize="small" />
              <Typography variant="body1" color="text.secondary">
                {user.email}
              </Typography>
              {user.emailVerified && (
                <Chip
                  icon={<VerifiedIcon />}
                  label="Verified"
                  size="small"
                  color="success"
                  variant="outlined"
                />
              )}
            </Stack>
          </Box>
        </ProfileHeader>
      </ProfilePaper>

      {/* Profile Details */}
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Account Activity
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <Stack spacing={2}>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Member Since
                  </Typography>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <CalendarTodayIcon fontSize="small" color="action" />
                    <Typography variant="body2">
                      {formatDate(user.createdAt)}
                    </Typography>
                  </Stack>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Last Updated
                  </Typography>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <CalendarTodayIcon fontSize="small" color="action" />
                    <Typography variant="body2">
                      {formatDate(user.updatedAt)}
                    </Typography>
                  </Stack>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Email Status
                  </Typography>
                  <Chip
                    label={user.emailVerified ? 'Verified' : 'Not Verified'}
                    color={user.emailVerified ? 'success' : 'warning'}
                    size="small"
                  />
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Stack>
  );
};

export default UserProfile;