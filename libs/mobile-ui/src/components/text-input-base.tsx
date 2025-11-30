import { forwardRef, ReactNode } from "react";
import {
  StyleProp,
  StyleSheet,
  TextInput as RNTextInput,
  TextInputProps as RNTextInputProps,
  View,
  ViewStyle,
} from "react-native";

import { colors } from "../theme/tokens";

export type TextInputBaseProps = RNTextInputProps & {
  containerStyle?: StyleProp<ViewStyle>;
  leadingComponent?: ReactNode;
  trailingComponent?: ReactNode;
};

export type TextInputBaseRef = RNTextInput;

export const TextInputBase = forwardRef<TextInputBaseRef, TextInputBaseProps>(
  (
    { containerStyle, leadingComponent, trailingComponent, style, ...props },
    ref,
  ) => (
    <View style={[styles.container, containerStyle]}>
      {leadingComponent}
      <RNTextInput {...props} ref={ref} style={[styles.input, style]} />
      {trailingComponent}
    </View>
  ),
);

const styles = StyleSheet.create({
  container: {
    alignItems: "center",
    flexDirection: "row",
    gap: 8,
    marginBottom: 8,
    paddingHorizontal: 8,
    width: "100%",
  },
  input: {
    color: colors.foreground,
    flex: 1,
    fontSize: 16,
    height: "100%",
  },
});
