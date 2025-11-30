import { useMemo } from "react";
import { StyleProp, View, ViewStyle } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";

type SafeAreaInsetProps = {
  minSize?: number;
  placement: "bottom" | "left" | "right" | "top";
  style?: StyleProp<ViewStyle>;
};

export const SafeAreaInset = ({
  minSize,
  placement,
  style,
}: SafeAreaInsetProps) => {
  const insets = useSafeAreaInsets();

  const insetStyle = useMemo(
    () =>
      ({
        bottom: { height: getMinmumInset(insets.bottom, minSize) },
        left: { width: getMinmumInset(insets.left, minSize) },
        right: { width: getMinmumInset(insets.right, minSize) },
        top: { height: getMinmumInset(insets.top, minSize) },
      }[placement]),
    [insets, placement, minSize],
  );

  return <View style={[style, insetStyle]} />;
};

const getMinmumInset = (
  inset: number | undefined,
  minSize: number | undefined,
) => Math.max(inset || 0, minSize || 0);
