import { colors, darkColors } from "./colors";

export const theme = {
  borders: {
    sm: 1,
    md: 2,
    lg: 3,
  },
  colors,
  fontSizes: {
    xs: 12,
    sm: 14,
    md: 16,
    lg: 18,
    xl: 20,
    xl_2x: 24,
    xl_3x: 28,
    xl_4x: 32,
    xl_5x: 40,
    xl_6x: 48,
  },
  radii: {
    xs: 4,
    sm: 6,
    md: 8,
    lg: 12,
  },
  spaces: {
    xs: 4,
    sm: 8,
    md: 12,
    lg: 16,
    xl: 24,
    xl_2x: 32,
    xl_3x: 40,
    xl_4x: 48,
    xl_5x: 56,
    xl_6x: 64,
  },
};

export type Theme = typeof theme;

export const themes = {
  dark: { ...theme, colors: darkColors },
  light: theme,
};

export type Themes = typeof themes;
