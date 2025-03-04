import { useContext } from "react";
import { UnistylesRuntime } from "react-native-unistyles";

import { ButtonContext } from "./button";
import { Icon, IconProps } from "./icon";

export const ButtonIcon = (props: IconProps) => {
  const button = useContext(ButtonContext);
  const theme = UnistylesRuntime.getTheme();

  return (
    <Icon
      {...props}
      color={props.color || button.color}
      size={props.size || theme.spaces.lg}
    />
  );
};
