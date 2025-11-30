import React, { forwardRef } from "react";
import {
  FieldValues,
  useController,
  UseControllerProps,
} from "react-hook-form";

import { TextInput, TextInputProps } from "./text-input";
import { TextInputBaseRef } from "./text-input-base";

export type TextInputContollerProps<V extends FieldValues> = TextInputProps &
  UseControllerProps<V>;

const TextInputControllerView = <V extends FieldValues>(
  { name, control, ...props }: TextInputContollerProps<V>,
  ref: React.Ref<TextInputBaseRef>,
) => {
  const controller = useController<V>({ name, control });

  return (
    <TextInput
      ref={ref}
      error={controller.fieldState.error?.message}
      {...props}
      onBlur={controller.field.onBlur}
      onChangeText={controller.field.onChange}
      value={controller.field.value}
    />
  );
};

export const TextInputController = forwardRef(TextInputControllerView) as <
  V extends FieldValues,
>(
  props: TextInputContollerProps<V> & { ref?: React.Ref<TextInputBaseRef> },
) => React.ReactElement;
