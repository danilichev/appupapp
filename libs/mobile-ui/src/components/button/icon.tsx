import { Icon, IconProps } from "../icon";
import { ButtonChildProps } from "./base";

export const renderButtonIcon =
  (props: IconProps) =>
  ({ color }: ButtonChildProps) =>
    <Icon {...props} color={props.color || color} />;
