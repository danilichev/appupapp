import {
  Text as RNText,
  TextProps as RNTextProps,
  TextStyle,
} from "react-native";
import { StyleSheet } from "react-native-unistyles";

import { Theme } from "../theme";

type TextProps = RNTextProps & {
  size?: keyof Theme["components"]["text"]["text"]["variants"]["size"];
  variant?: keyof Theme["components"]["text"]["text"]["variants"]["variant"];
  weight?: keyof Theme["components"]["text"]["text"]["variants"]["weight"];
};

export const Text = ({
  size = "md",
  style,
  variant = "p",
  weight = "regular",
  ...rest
}: TextProps) => {
  styles.useVariants({ size, variant, weight });

  return <RNText {...rest} style={[styles.text as TextStyle, style]} />;
};

const styles = StyleSheet.create((theme) => ({
  text: {
    fontFamily: theme.fontFamily,
    fontSize: theme.fontSizes.md,
    variants: theme.components.text.text.variants,
  },
}));
