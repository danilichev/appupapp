import { FC, ReactNode, useCallback, useMemo } from "react";
import {
  ActivityIndicator,
  StyleProp,
  StyleSheet,
  Text,
  ViewStyle,
} from "react-native";

import { colors } from "../../theme/colors";
import { Pressable } from "../pressable";

export type ButtonChildProps = { color?: string };

export type ButtonChild = ReactNode | FC<ButtonChildProps>;

export type BaseButtonProps = {
  backgroundColor?: string;
  color?: string;
  children?: ButtonChild | ButtonChild[];
  containerStyle?: StyleProp<ViewStyle>;
  isDisabled?: boolean;
  isLoading?: boolean;
  loadingComponent?: ButtonChild;
  loadingPosition?: "start" | "center" | "end";
  onPress?: () => void;
  textStyle?: StyleProp<ViewStyle>;
};

export const BaseButton = ({
  backgroundColor = colors.primary,
  color = colors.primaryForeground,
  children,
  containerStyle,
  isDisabled,
  isLoading,
  loadingComponent,
  loadingPosition = "center",
  onPress,
  textStyle,
}: BaseButtonProps) => {
  const [startComponent, centerComponent, endComponent] = useMemo(
    () => [
      ...(Array.isArray(children) ? children : [null, children]),
      ...(Array.from({ length: 3 }).fill(null) as ButtonChild[]),
    ],
    [children],
  );

  const loadingComp = useMemo(
    () =>
      typeof loadingComponent === "function" ? (
        loadingComponent({ color })
      ) : (
        <ActivityIndicator color={color} />
      ),
    [color, loadingComponent],
  );

  const renderChild = useCallback(
    (
      position: Required<BaseButtonProps>["loadingPosition"],
      component?: ButtonChild,
    ) =>
      isLoading && loadingPosition === position ? (
        loadingComp
      ) : typeof component === "function" ? (
        component({ color })
      ) : typeof component === "string" ? (
        <Text style={[styles.text, textStyle, { color }]}>{component}</Text>
      ) : (
        component
      ),
    [color, isLoading, loadingComp, loadingPosition, textStyle],
  );

  return (
    <Pressable
      isDisabled={isDisabled || isLoading}
      onPress={onPress}
      style={[styles.container, containerStyle, { backgroundColor }]}
    >
      {renderChild("start", startComponent)}
      {renderChild("center", centerComponent)}
      {renderChild("end", endComponent)}
    </Pressable>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: "center",
    flexDirection: "row",
    padding: 8,
  },
  text: {
    color: colors.primaryForeground,
    fontSize: 16,
  },
});
