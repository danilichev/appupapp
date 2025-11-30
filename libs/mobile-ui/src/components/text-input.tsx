import {
  createContext,
  forwardRef,
  useCallback,
  useContext,
  useState,
} from "react";
import {
  ColorValue,
  NativeSyntheticEvent,
  TextInputFocusEventData,
} from "react-native";
import { UnistylesRuntime } from "react-native-unistyles";
import { StyleSheet } from "react-native-unistyles";

import { Icon, IconProps } from "./icon";
import { Text } from "./text";
import {
  TextInputBase,
  TextInputBaseProps,
  TextInputBaseRef,
} from "./text-input-base";

export type TextInputProps = TextInputBaseProps & {
  error?: string;
  label?: string;
  iconColor?: ColorValue;
};

export const TextInputContext = createContext<
  Pick<TextInputProps, "iconColor">
>({});

export const TextInput = Object.assign(
  forwardRef<TextInputBaseRef, TextInputProps>(
    (
      {
        containerStyle,
        error,
        iconColor,
        label,
        style,
        ...props
      }: TextInputProps,
      ref,
    ) => {
      const [isFocused, setFocused] = useState(false);

      const onBlur = useCallback(
        (e: NativeSyntheticEvent<TextInputFocusEventData>) => {
          setFocused(false);
          props.onBlur?.(e);
        },
        [props],
      );

      const onFocus = useCallback(
        (e: NativeSyntheticEvent<TextInputFocusEventData>) => {
          setFocused(true);
          props.onFocus?.(e);
        },
        [props],
      );

      const theme = UnistylesRuntime.getTheme();
      styles.useVariants({
        state: error ? "error" : isFocused ? "focused" : undefined,
      });

      return (
        <>
          {label ? (
            <Text style={styles.label} variant="label" weight="bold">
              {label}
            </Text>
          ) : null}
          <TextInputContext.Provider
            value={{
              iconColor:
                iconColor ||
                props.placeholderTextColor ||
                theme.colors.mutedForeground,
            }}
          >
            <TextInputBase
              placeholderTextColor={theme.colors.mutedForeground}
              {...props}
              containerStyle={[styles.container, containerStyle]}
              onBlur={onBlur}
              onFocus={onFocus}
              ref={ref}
              style={[styles.input, style]}
            />
          </TextInputContext.Provider>
          {error ? <Text style={styles.error}>{error}</Text> : null}
        </>
      );
    },
  ),
  {
    Icon: TextInputIcon,
  },
);

function TextInputIcon(props: IconProps) {
  const context = useContext(TextInputContext);
  const theme = UnistylesRuntime.getTheme();

  return (
    <Icon
      {...props}
      color={props.color || context.iconColor}
      size={props.size || theme.spaces.md}
    />
  );
}

const styles = StyleSheet.create((theme) => ({
  container: {
    alignItems: "center",
    backgroundColor: theme.colors.background,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.md,
    borderWidth: theme.borders.sm,
    flexDirection: "row",
    gap: theme.spaces.sm,
    height: theme.spaces.xl_4,
    marginBottom: theme.spaces.sm,
    paddingHorizontal: theme.spaces.md,
    width: "100%",
    variants: {
      state: {
        error: {
          borderColor: theme.colors.destructive,
          marginBottom: theme.spaces.xs_2,
        },
        focused: {
          borderColor: theme.colors.primary,
        },
      },
    },
  },
  error: {
    alignSelf: "flex-start",
    color: theme.colors.destructive,
    fontSize: theme.fontSizes.xs,
    marginBottom: theme.spaces.sm,
  },
  label: {
    alignSelf: "flex-start",
    marginBottom: theme.spaces.xs,
  },
  input: {
    color: theme.colors.foreground,
    flex: 1,
    fontSize: theme.fontSizes.md,
    height: "100%",
  },
}));
