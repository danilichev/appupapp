import { View, ViewProps } from "react-native";

export type BoxProps = ViewProps & {
  isCentered?: boolean;
  isHorizontal?: boolean;
};

export const Box = ({
  isCentered,
  isHorizontal,
  style,
  ...props
}: BoxProps) => (
  <View
    {...props}
    style={[
      isCentered && { alignItems: "center", justifyContent: "center" },
      isHorizontal && { flexDirection: "row" },
      style,
    ]}
  />
);
