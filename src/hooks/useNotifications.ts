import { useState, useEffect } from 'react';
import { userApi, UserPreferences } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

export const useNotifications = () => {
  const { isAuthenticated } = useAuth();
  const [preferences, setPreferences] = useState<UserPreferences | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadPreferences = async () => {
    if (!isAuthenticated) return;
    
    try {
      setLoading(true);
      setError(null);
      const prefs = await userApi.getPreferences();
      setPreferences(prefs);
    } catch (err) {
      console.error('Failed to load preferences:', err);
      setError('Failed to load notification preferences');
      // Set default preferences if loading fails
      setPreferences({
        id: '',
        userId: '',
        emailNotifications: true,
        pushNotifications: true,
        notificationFrequency: 'immediate',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      });
    } finally {
      setLoading(false);
    }
  };

  const updatePreferences = async (updates: Partial<UserPreferences>) => {
    if (!isAuthenticated || !preferences) return;

    try {
      setLoading(true);
      setError(null);
      const updatedPrefs = await userApi.updatePreferences(updates);
      setPreferences(updatedPrefs);
      return updatedPrefs;
    } catch (err) {
      console.error('Failed to update preferences:', err);
      setError('Failed to update notification preferences');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPreferences();
  }, [isAuthenticated]);

  return {
    preferences,
    loading,
    error,
    updatePreferences,
    refreshPreferences: loadPreferences,
  };
};