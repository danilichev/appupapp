import { colors, darkColors } from "./colors";

export const theme = {
  borders: {
    sm: 1,
    md: 2,
    lg: 3,
  },
  colors,
  fontFamily: "System",
  fontSizes: {
    xs: 12,
    sm: 14,
    md: 16,
    lg: 18,
    xl: 20,
    xl_2: 24,
    xl_3: 28,
    xl_4: 32,
    xl_5: 40,
    xl_6: 48,
  },
  radii: {
    xs: 4,
    sm: 6,
    md: 8,
    lg: 12,
  },
  spaces: {
    xs_3: 2,
    xs_2: 4,
    xs: 8,
    sm: 12,
    md: 16,
    lg: 20,
    xl: 24,
    xl_2: 32,
    xl_3: 40,
    xl_4: 48,
    xl_5: 56,
    xl_6: 64,
  },
};

export type Theme = typeof theme;

export const themes = {
  dark: { ...theme, colors: darkColors },
  light: theme,
};

export type Themes = typeof themes;
