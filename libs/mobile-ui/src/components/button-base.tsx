import { ReactNode, useCallback, useMemo } from "react";
import {
  ActivityIndicator,
  StyleProp,
  StyleSheet,
  Text,
  TextStyle,
  ViewStyle,
} from "react-native";

import { colors } from "../theme/tokens";
import { Pressable } from "./pressable";

export type ButtonBaseProps = {
  backgroundColor?: string;
  children?: ReactNode;
  color?: string;
  isDisabled?: boolean;
  isLoading?: boolean;
  leadingComponent?: ReactNode;
  loadingComponent?: ReactNode;
  loadingPosition?: "start" | "center" | "end";
  onPress?: () => void;
  style?: StyleProp<ViewStyle>;
  textStyle?: StyleProp<TextStyle>;
  trailingComponent?: ReactNode;
};

export const ButtonBase = ({
  backgroundColor = colors.primary,
  children,
  color = colors.primaryForeground,
  isDisabled,
  isLoading,
  leadingComponent,
  loadingComponent,
  loadingPosition = "center",
  onPress,
  style,
  textStyle,
  trailingComponent,
}: ButtonBaseProps) => {
  const loadingComp = useMemo(
    () => loadingComponent || <ActivityIndicator color={color} />,
    [color, loadingComponent],
  );

  const renderChild = useCallback(
    (
      position: Required<ButtonBaseProps>["loadingPosition"],
      component?: ReactNode,
    ) =>
      isLoading && loadingPosition === position ? (
        loadingComp
      ) : typeof component === "string" ? (
        <Text style={[styles.text, { color }, textStyle]}>{component}</Text>
      ) : (
        component
      ),
    [color, isLoading, loadingComp, loadingPosition, textStyle],
  );

  return (
    <Pressable
      isDisabled={isDisabled || isLoading}
      onPress={onPress}
      style={[styles.container, { backgroundColor }, style]}
      withFeedback
    >
      {renderChild("start", leadingComponent)}
      {renderChild("center", children)}
      {renderChild("end", trailingComponent)}
    </Pressable>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: "center",
    alignSelf: "stretch",
    flexDirection: "row",
    justifyContent: "center",
    padding: 8,
    width: "auto",
  },
  text: {
    color: colors.primaryForeground,
    fontSize: 16,
  },
});
