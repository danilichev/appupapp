import { BaseButton, BaseButtonProps } from "./base";

export { BaseButton };

export type ButtonProps = BaseButtonProps & {
  size?: "md" | "lg" | "sm";
  variant?: "outline" | "primary" | "secondary";
};

export const Button = ({ size, variant, ...rest }: ButtonProps) => {
  return <BaseButton {...rest} />;
};
