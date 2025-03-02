import { StyleSheet } from "react-native-unistyles";

import { Theme } from "../../theme";
import { BaseButton, BaseButtonProps } from "./base";

export type ButtonProps = BaseButtonProps & {
  size?: "lg" | "md" | "sm";
  variant?: "outline" | "primary" | "secondary";
};

export const Button = ({
  containerStyle,
  size = "md",
  textStyle,
  variant = "primary",
  ...rest
}: ButtonProps) => {
  styles.useVariants({ size, variant });

  return (
    <BaseButton
      {...rest}
      color={styles.text.color}
      containerStyle={[styles.container, containerStyle]}
      textStyle={[styles.text, textStyle]}
    />
  );
};

const styles = StyleSheet.create((theme: Theme) => ({
  container: {
    variants: {
      size: {
        lg: {
          borderRadius: theme.radii.md,
          gap: theme.spaces.lg,
          paddingHorizontal: theme.spaces.xl,
          paddingVertical: theme.spaces.lg,
        },
        md: {
          borderRadius: theme.radii.sm,
          gap: theme.spaces.md,
          paddingHorizontal: theme.spaces.lg,
          paddingVertical: theme.spaces.md,
        },
        sm: {
          borderRadius: theme.radii.xs,
          gap: theme.spaces.sm,
          paddingHorizontal: theme.spaces.md,
          paddingVertical: theme.spaces.sm,
        },
      },
      variant: {
        outline: {
          backgroundColor: "transparent",
          borderColor: theme.colors.primary,
          borderWidth: theme.borders.md,
        },
        primary: { backgroundColor: theme.colors.primary },
        secondary: { backgroundColor: theme.colors.secondary },
      },
    },
  },
  text: {
    fontSize: theme.fontSizes.md,
    fontWeight: "bold",
    variants: {
      variant: {
        outline: { color: theme.colors.primary },
        primary: { color: theme.colors.primaryForeground },
        secondary: { color: theme.colors.secondaryForeground },
      },
    },
  },
}));
