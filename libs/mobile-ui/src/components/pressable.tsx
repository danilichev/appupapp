import { useCallback } from "react";
import {
  Pressable as RNPressable,
  PressableProps as RNPressableProps,
  PressableStateCallbackType,
} from "react-native";

export type PressableProps = Omit<RNPressableProps, "disabled"> & {
  activeOpacity?: number;
  isDisabled?: boolean;
};

export const Pressable = ({
  activeOpacity = 0.8,
  isDisabled,
  style,
  ...rest
}: PressableProps) => {
  const getStyles = useCallback(
    (state: PressableStateCallbackType) => {
      return [
        typeof style === "function" ? style(state) : style,
        state.pressed && { opacity: activeOpacity },
      ];
    },
    [activeOpacity, style],
  );

  return <RNPressable {...rest} disabled={isDisabled} style={getStyles} />;
};
