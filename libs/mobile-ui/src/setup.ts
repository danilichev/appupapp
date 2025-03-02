import { StyleSheet } from "react-native-unistyles";

import { Themes, themes } from "./theme";

declare module "react-native-unistyles" {
  export interface UnistylesThemes {
    light: Themes["light"];
  }
}

StyleSheet.configure({
  settings: { initialTheme: "light" },
  themes: themes,
});
