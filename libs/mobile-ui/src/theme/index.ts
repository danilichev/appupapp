import { TextStyle } from "react-native";
import { UnistylesValues } from "react-native-unistyles/lib/typescript/src/types";

import { darkColors, Tokens, tokens } from "./tokens";

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
            h1: { fontSize: tokens.fontSizes.xl_5 },
            h2: { fontSize: tokens.fontSizes.xl_4 },
            h3: { fontSize: tokens.fontSizes.xl_3 },
            h4: { fontSize: tokens.fontSizes.xl_2 },
            h5: { fontSize: tokens.fontSizes.xl },
            h6: { fontSize: tokens.fontSizes.lg },
            p: { fontSize: tokens.fontSizes.md },
            label: { fontSize: tokens.fontSizes.sm, fontWeight: "700" },
          },
          weight: {
            light: { fontWeight: "200" },
            regular: { fontWeight: "400" },
            medium: { fontWeight: "500" },
            bold: { fontWeight: "700" },
          },
        },
      } as UnistylesValues & {
        variants: {
          size: Record<keyof Tokens["fontSizes"], TextStyle>;
          variant: Record<
            "h1" | "h2" | "h3" | "h4" | "h5" | "h6" | "p" | "label",
            TextStyle
          >;
          weight: Record<"light" | "regular" | "medium" | "bold", TextStyle>;
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
