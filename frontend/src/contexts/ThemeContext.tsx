import React, { createContext, useState, useContext, useCallback, useEffect } from 'react';
import { api } from '../lib/api';

export type ThemeMode = 'dark' | 'light';

interface ThemeContextType {
  theme: ThemeMode;
  setTheme: (theme: ThemeMode) => Promise<void>;
  isLoading: boolean;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

const THEME_STORAGE_KEY = 'theme_preference';

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [theme, setThemeState] = useState<ThemeMode>('dark');
  const [isLoading, setIsLoading] = useState(true);

  // Load theme preference from storage
  useEffect(() => {
    const loadTheme = async () => {
      try {
        // Load theme preference
        const savedTheme = await api.getSetting(THEME_STORAGE_KEY);

        if (savedTheme) {
          const themeMode = savedTheme as ThemeMode;
          setThemeState(themeMode);
          await applyTheme(themeMode);
        } else {
          // No saved preference: apply dark as the default theme
          setThemeState('dark');
          await applyTheme('dark');
        }
      } catch (error) {
        console.error('Failed to load theme settings:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadTheme();
  }, []);

  // Apply theme to document
  const applyTheme = useCallback(async (themeMode: ThemeMode) => {
    const root = document.documentElement;

    // Remove all theme classes
    root.classList.remove('theme-dark', 'theme-light');

    // Add new theme class
    root.classList.add(`theme-${themeMode}`);

    // Note: Window theme updates removed since we're using custom titlebar
  }, []);

  const setTheme = useCallback(async (newTheme: ThemeMode) => {
    try {
      setIsLoading(true);

      // Apply theme immediately
      setThemeState(newTheme);
      await applyTheme(newTheme);

      // Save to storage
      await api.saveSetting(THEME_STORAGE_KEY, newTheme);
    } catch (error) {
      console.error('Failed to save theme preference:', error);
    } finally {
      setIsLoading(false);
    }
  }, [applyTheme]);

  
  const value: ThemeContextType = {
    theme,
    setTheme,
    isLoading,
  };

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useThemeContext = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useThemeContext must be used within a ThemeProvider');
  }
  return context;
};
