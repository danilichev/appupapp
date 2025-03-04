import { UnistylesRuntime } from "react-native-unistyles";

import { Themes, themes } from "@appupapp/mobile-ui";

declare module "react-native-unistyles" {
  export interface UnistylesThemes {
    dark: Themes["dark"];
    light: Themes["light"];
  }
}

(Object.keys(themes) as (keyof Themes)[]).forEach((key) => {
  UnistylesRuntime.updateTheme(key, (theme) => ({
    ...theme,
    fontFamily: "Inter-Bold",
  }));
});
