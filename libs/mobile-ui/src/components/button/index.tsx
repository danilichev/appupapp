import { createContext, useMemo } from "react";
import { StyleSheet } from "react-native-unistyles";

import { Theme } from "../../theme";
import { BaseButton, BaseButtonProps } from "./base";

export type ButtonProps = BaseButtonProps & {
  size?: "lg" | "md" | "sm";
  variant?: "outline" | "primary" | "secondary";
};

export const ButtonContext = createContext<
  Pick<ButtonProps, "color" | "size" | "variant">
>({});

export const Button = ({
  size = "md",
  style,
  textStyle,
  variant = "primary",
  ...rest
}: ButtonProps) => {
  styles.useVariants({ size, variant });

  const color = useMemo(() => rest.color || styles.text.color, [rest.color]);

  console.log("color", color);

  return (
    <ButtonContext.Provider value={{ color, size, variant }}>
      <BaseButton
        {...rest}
        color={color}
        style={[styles.container, style]}
        textStyle={[styles.text, textStyle]}
      />
    </ButtonContext.Provider>
  );
};

const styles = StyleSheet.create((theme: Theme) => ({
  container: {
    variants: {
      size: {
        lg: {
          borderRadius: theme.radii.md,
          gap: theme.spaces.lg,
          paddingHorizontal: theme.spaces.lg,
          paddingVertical: theme.spaces.md,
        },
        md: {
          borderRadius: theme.radii.sm,
          gap: theme.spaces.md,
          paddingHorizontal: theme.spaces.md,
          paddingVertical: theme.spaces.sm,
        },
        sm: {
          borderRadius: theme.radii.xs,
          gap: theme.spaces.sm,
          paddingHorizontal: theme.spaces.sm,
          paddingVertical: theme.spaces.xs,
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
    ...(theme.components.button?.container as object),
  },
  text: {
    fontFamily: theme.fontFamily,
    fontSize: theme.fontSizes.md,
    variants: {
      variant: {
        outline: { color: theme.colors.primary },
        primary: { color: theme.colors.primaryForeground },
        secondary: { color: theme.colors.secondaryForeground },
      },
    },
    ...(theme.components.button?.text as object),
  },
}));
