import { StyleSheet } from "react-native-unistyles";

import { Themes, themes } from "@appupapp/mobile-ui";

declare module "react-native-unistyles" {
  export interface UnistylesThemes {
    dark: Themes["dark"];
    light: Themes["light"];
  }
}

StyleSheet.configure({
  settings: { initialTheme: "light" },
  themes,
});
