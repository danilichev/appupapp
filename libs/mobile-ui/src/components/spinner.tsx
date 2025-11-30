import { ActivityIndicator, ActivityIndicatorProps } from "react-native";
import { StyleSheet, UnistylesRuntime } from "react-native-unistyles";

export type SpinnerProps = Pick<
  ActivityIndicatorProps,
  "color" | "size" | "style"
>;

export const Spinner = (props: SpinnerProps) => {
  const theme = UnistylesRuntime.getTheme();

  return (
    <ActivityIndicator
      color={theme.colors.foreground}
      {...props}
      style={[styles.spinner, props.style]}
    />
  );
};

const styles = StyleSheet.create((theme) => ({
  spinner: {
    margin: theme.spaces.md,
  },
}));
