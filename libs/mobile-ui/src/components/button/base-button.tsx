import { ReactNode, useCallback, useMemo } from "react";
import {
  ActivityIndicator,
  StyleProp,
  StyleSheet,
  Text,
  ViewStyle,
} from "react-native";

import { Pressable } from "../pressable";

export type BaseButtonProps = {
  children?: ReactNode;
  containerStyle?: StyleProp<ViewStyle>;
  isDisabled?: boolean;
  isLoading?: boolean;
  loadingComponent?: ReactNode;
  loadingPosition?: "start" | "center" | "end";
  onPress?: () => void;
  textStyle?: StyleProp<ViewStyle>;
};

export const BaseButton = ({
  children,
  containerStyle,
  isDisabled,
  isLoading,
  loadingComponent = <ActivityIndicator />,
  loadingPosition = "center",
  onPress,
  textStyle,
}: BaseButtonProps) => {
  const [startComponent, centerComponent, endComponent] = useMemo(
    () =>
      typeof children === "string"
        ? [null, children, null]
        : Array.isArray(children)
        ? children
        : [null, null, null],
    [children],
  );

  const maybeLoader = useCallback(
    (
      position: Required<BaseButtonProps>["loadingPosition"],
      component?: ReactNode,
    ) =>
      isLoading && loadingPosition === position ? loadingComponent : component,
    [isLoading, loadingComponent, loadingPosition],
  );

  return (
    <Pressable
      isDisabled={isDisabled || isLoading}
      onPress={onPress}
      style={[styles.container, containerStyle]}
    >
      {maybeLoader("start", startComponent)}
      {maybeLoader(
        "center",
        typeof centerComponent === "string" ? (
          <Text style={[styles.title, textStyle]}>{centerComponent}</Text>
        ) : (
          centerComponent
        ),
      )}
      {maybeLoader("end", endComponent)}
    </Pressable>
  );
};

const styles = StyleSheet.create({
  container: { alignItems: "center", flexDirection: "row", padding: 8 },
  title: { fontSize: 16 },
});
