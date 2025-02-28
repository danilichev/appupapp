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

type ButtonChild = ReactNode | FC<{ color?: string }>;

export type BaseButtonProps = {
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
    () =>
      Array.isArray(children)
        ? ([...children, ...Array.from({ length: 3 }).fill(null)].slice(
            0,
            3,
          ) as ButtonChild[])
        : [null, children, null],
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
      style={[styles.container, containerStyle]}
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
    backgroundColor: colors.primary,
    flexDirection: "row",
    padding: 8,
  },
  text: {
    color: colors.primaryForeground,
    fontSize: 16,
  },
});
