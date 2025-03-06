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
};

export const Text = ({ size, style, variant, ...rest }: TextProps) => {
  styles.useVariants({ size, variant });

  return <RNText {...rest} style={[styles.text as TextStyle, style]} />;
};

const styles = StyleSheet.create((theme) => ({
  text: {
    fontFamily: theme.fontFamily,
    fontSize: theme.fontSizes.md,
    variants: theme.components.text.text.variants,
  },
}));
