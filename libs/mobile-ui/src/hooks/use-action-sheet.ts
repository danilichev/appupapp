import { useActionSheet as useRNActionSheet } from "@expo/react-native-action-sheet";
import { useCallback, useMemo } from "react";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { StyleSheet, UnistylesRuntime } from "react-native-unistyles";

type ActionSheetOption<T> = {
  isDisabled?: boolean | ((props?: T) => boolean);
  onPress?: (props?: T) => void;
  title: string | ((props?: T) => string);
};

type UseActionSheetProps<T> = {
  cancelOption?: Partial<ActionSheetOption<T>>;
  destructiveOption?: ActionSheetOption<T>;
  options?: ActionSheetOption<T>[];
};

export const useActionSheet = <T>({
  cancelOption,
  destructiveOption,
  options,
}: UseActionSheetProps<T>) => {
  const insets = useSafeAreaInsets();
  const theme = UnistylesRuntime.getTheme();

  const { showActionSheetWithOptions } = useRNActionSheet();

  const allOptions = useMemo(
    () =>
      [
        ...(options || []),
        destructiveOption,
        {
          title: cancelOption?.title || "Cancel",
          onPress: cancelOption?.onPress,
        },
      ].filter(Boolean) as ActionSheetOption<T>[],
    [options, destructiveOption, cancelOption],
  );

  const showActionSheet = useCallback(
    (props?: T) => {
      const titles = allOptions.map((opt) =>
        typeof opt.title === "function" ? opt.title(props) : opt.title,
      );

      const disabledIndices = allOptions.reduce<number[]>((acc, opt, idx) => {
        const disabled =
          typeof opt.isDisabled === "function"
            ? opt.isDisabled(props)
            : opt.isDisabled;
        return disabled ? [...acc, idx] : acc;
      }, []);

      showActionSheetWithOptions(
        {
          cancelButtonIndex: allOptions.length - 1,
          containerStyle: {
            ...styles.container,
            paddingBottom: insets.bottom,
          },
          destructiveButtonIndex: destructiveOption
            ? allOptions.findIndex(
                (opt) => opt.title === destructiveOption.title,
              )
            : undefined,
          destructiveColor: theme.colors.destructiveForeground,
          disabledButtonIndices: disabledIndices,
          //TODO: add support for disabled color
          options: titles,
          textStyle: styles.text,
        },
        (selectedIdx) => {
          allOptions[selectedIdx ?? -1]?.onPress?.(props);
        },
      );
    },
    [
      allOptions,
      destructiveOption,
      insets.bottom,
      showActionSheetWithOptions,
      theme.colors.destructiveForeground,
    ],
  );

  const onShowActionSheet = useCallback(
    (props?: T) => () => showActionSheet(props),
    [showActionSheet],
  );

  return { onShow: onShowActionSheet, show: showActionSheet };
};

const styles = StyleSheet.create((theme) => ({
  container: {
    borderTopEndRadius: theme.radii.xl,
    borderTopStartRadius: theme.radii.xl,
  },
  text: {
    color: theme.colors.foreground,
    fontFamily: theme.fontFamily,
    fontSize: theme.fontSizes.md,
  },
}));
