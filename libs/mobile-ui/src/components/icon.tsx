import { Feather } from "@expo/vector-icons";
import { ComponentProps } from "react";

export type IconProps = Pick<
  ComponentProps<typeof Feather>,
  "color" | "name" | "size" | "style"
>;

export const Icon = (props: IconProps) => <Feather {...props} />;
