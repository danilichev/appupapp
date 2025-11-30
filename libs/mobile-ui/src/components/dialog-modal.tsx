import { Modal, ModalFooterProps, ModalProps } from "./modal";
import { Text } from "./text";

export type DialogModalProps = Omit<ModalProps, "children"> &
  ModalFooterProps & {
    title?: string;
    children?: React.ReactNode;
  };

export const DialogModal = ({
  children,
  isDestructive,
  isSubmitting,
  onClose,
  onSubmit,
  submitText,
  title,
  ...props
}: DialogModalProps) => {
  return (
    <Modal {...props} onClose={onClose}>
      {title ? <Modal.Header title={title} /> : null}
      {children ? (
        <Modal.Content>
          {typeof children === "string" ? (
            <Text variant="p" size="md" weight="medium">
              {children}
            </Text>
          ) : (
            children
          )}
        </Modal.Content>
      ) : null}
      <Modal.Footer
        isDestructive={isDestructive}
        isSubmitting={isSubmitting}
        onCancel={onClose}
        onSubmit={onSubmit}
        submitText={submitText}
      />
    </Modal>
  );
};
