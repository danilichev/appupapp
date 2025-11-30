import { createContext, FC, useContext, useMemo } from "react";
import { StyleSheet, UnistylesRuntime } from "react-native-unistyles";

import { Theme } from "../theme";
import { ButtonBase, ButtonBaseProps } from "./button-base";
import { Icon, IconProps } from "./icon";

export type ButtonProps = ButtonBaseProps & {
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
  styles.useVariants({ isDisabled: rest.isDisabled, size, variant });

  const color = useMemo(() => rest.color || styles.text.color, [rest.color]);

  return (
    <ButtonContext.Provider value={{ color, size, variant }}>
      <ButtonBase
        {...rest}
        color={color}
        style={[styles.container, style]}
        textStyle={[styles.text, textStyle]}
      />
    </ButtonContext.Provider>
  );
};

Button.Icon = ((props: IconProps) => {
  const button = useContext(ButtonContext);
  const theme = UnistylesRuntime.getTheme();

  return (
    <Icon
      {...props}
      color={props.color || button.color}
      size={props.size || theme.spaces.lg}
    />
  );
}) as FC<IconProps>;

const styles = StyleSheet.create((theme: Theme) => ({
  container: {
    variants: {
      isDisabled: {
        false: {},
        true: { opacity: 0.8 },
      },
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
