import { TextStyle } from "react-native";
import { UnistylesValues } from "react-native-unistyles/lib/typescript/src/types";

import { colors, darkColors } from "./colors";

const tokens = {
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

export type Tokens = typeof tokens;

export const theme = {
  ...tokens,
  components: {
    button: {} as
      | { container?: UnistylesValues; text?: UnistylesValues }
      | undefined,
    text: {
      text: {
        variants: {
          size: Object.fromEntries(
            Object.entries(tokens.fontSizes).map(([key, value]) => [
              key,
              { fontSize: value },
            ]),
          ),
          variant: {
            h1: { fontSize: tokens.fontSizes.xl_6 },
            h2: { fontSize: tokens.fontSizes.xl_4 },
            h3: { fontSize: tokens.fontSizes.xl_2 },
            h4: { fontSize: tokens.fontSizes.xl },
            p: { fontSize: tokens.fontSizes.md },
          },
        },
      } as UnistylesValues & {
        variants: {
          size: Record<keyof Tokens["fontSizes"], TextStyle>;
          variant: Record<"h1" | "h2" | "h3" | "h4" | "p", TextStyle>;
        };
      },
    },
  },
};

export type Theme = typeof theme;

export const themes = {
  dark: { ...theme, colors: darkColors },
  light: theme,
};

export type Themes = typeof themes;
