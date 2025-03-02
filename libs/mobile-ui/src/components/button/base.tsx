import { ReactNode, useCallback, useMemo } from "react";
import {
  ActivityIndicator,
  StyleProp,
  StyleSheet,
  Text,
  TextStyle,
  ViewStyle,
} from "react-native";

import { Pressable } from "../../components/pressable";
import { colors } from "../../theme/colors";

export type BaseButtonProps = {
  backgroundColor?: string;
  color?: string;
  children?: ReactNode;
  containerStyle?: StyleProp<ViewStyle>;
  isDisabled?: boolean;
  isLoading?: boolean;
  loadingComponent?: ReactNode;
  loadingPosition?: "start" | "center" | "end";
  onPress?: () => void;
  textStyle?: StyleProp<TextStyle>;
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
      ...Array.from({ length: 3 }).fill(null),
    ],
    [children],
  );

  const loadingComp = useMemo(
    () => loadingComponent || <ActivityIndicator color={color} />,
    [color, loadingComponent],
  );

  const renderChild = useCallback(
    (
      position: Required<BaseButtonProps>["loadingPosition"],
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
      style={[styles.container, { backgroundColor }, containerStyle]}
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
