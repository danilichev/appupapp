import { UnistylesRuntime } from "react-native-unistyles";

import { Theme, Themes, themes } from "@appupapp/mobile-ui";

declare module "react-native-unistyles" {
  export interface UnistylesThemes {
    dark: Themes["dark"];
    light: Themes["light"];
  }
}

(Object.keys(themes) as (keyof Themes)[]).forEach((key) => {
  UnistylesRuntime.updateTheme(key, (theme) => ({
    ...theme,
    fontFamily: "Inter-Medium",
    components: {
      ...theme.components,
      button: {
        text: {
          fontFamily: "Inter-Bold",
          letterSpacing: 0.5,
        },
      },
      text: {
        ...theme.components.text,
        text: {
          ...theme.components.text.text,
          variants: {
            ...theme.components.text.text.variants,
            variant: {
              ...theme.components.text.text.variants.variant,
              h1: {
                color: "blue",
                fontSize: 32,
                lineHeight: 40,
              },
            },
          } as Theme["components"]["text"]["text"]["variants"],
        },
      },
    },
  }));
});
