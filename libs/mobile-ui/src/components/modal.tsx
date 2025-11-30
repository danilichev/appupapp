import {
  Modal as RNModal,
  StyleProp,
  useWindowDimensions,
  View,
  ViewStyle,
} from "react-native";
import Animated, {
  runOnJS,
  SlideInUp,
  SlideOutUp,
} from "react-native-reanimated";
import { StyleSheet, useUnistyles } from "react-native-unistyles";

import { useToggle } from "@linksholder/react-hooks";

import { Button } from "./button";
import { Pressable } from "./pressable";
import { Text } from "./text";

export type ModalProps = {
  children: React.ReactNode;
  isOpen: boolean;
  onClose: () => void;
  slot?: React.ReactNode;
  style?: StyleProp<ViewStyle>;
};

export const Modal = ({
  children,
  isOpen,
  onClose,
  slot,
  style,
}: ModalProps) => {
  const contentToggler = useToggle();
  const window = useWindowDimensions();
  const { theme } = useUnistyles();

  return (
    <RNModal
      animationType="none"
      onRequestClose={onClose}
      onShow={contentToggler.on}
      transparent
      visible={isOpen}
    >
      <Pressable
        onPress={contentToggler.off}
        style={[
          styles.backdrop,
          { width: window.width, height: window.height },
        ]}
      />
      {contentToggler.isOn ? (
        <Animated.View
          entering={SlideInUp}
          exiting={SlideOutUp.withCallback(() => {
            runOnJS(onClose)();
          })}
          style={[
            styles.contentContainer,
            {
              height: window.height * 0.5,
              left: theme.spaces.lg,
              top: window.height * 0.25,
              width: window.width - theme.spaces.lg * 2,
            },
            style,
          ]}
        >
          {children}
        </Animated.View>
      ) : null}
      {slot}
    </RNModal>
  );
};

type ModalContentProps = {
  children: React.ReactNode;
  style?: StyleProp<ViewStyle>;
};

Modal.Content = function ModalContent({ children, style }: ModalContentProps) {
  return <View style={[styles.content, style]}>{children}</View>;
};

export type ModalFooterProps = {
  cancelText?: string;
  isDestructive?: boolean;
  isSubmitting?: boolean;
  onCancel?: () => void;
  onSubmit: () => void;
  submitText?: string;
};

Modal.Footer = function ModalFooter({
  cancelText,
  isDestructive,
  isSubmitting,
  onCancel,
  onSubmit,
  submitText,
}: ModalFooterProps) {
  return (
    <View style={styles.footer}>
      {onCancel ? (
        <Button
          isDisabled={isSubmitting}
          onPress={onCancel}
          style={styles.footerButton}
          variant="secondary"
        >
          {cancelText || "Cancel"}
        </Button>
      ) : null}
      <Button
        isLoading={isSubmitting}
        onPress={onSubmit}
        style={[
          styles.footerButton,
          isDestructive && styles.footerButtonDestructive,
        ]}
        textStyle={[isDestructive && styles.footerButtonTextDestructive]}
        variant="primary"
      >
        {submitText || "Submit"}
      </Button>
    </View>
  );
};

type ModalHeaderProps = {
  title: string;
};

Modal.Header = function ModalHeader({ title }: ModalHeaderProps) {
  return (
    <Text variant="h6" weight="bold" style={styles.title}>
      {title}
    </Text>
  );
};

const styles = StyleSheet.create((theme) => ({
  backdrop: {
    flex: 1,
    backgroundColor: "rgba(0, 0, 0, 0.2)",
  },
  content: {
    flex: 1,
  },
  contentContainer: {
    position: "absolute",
    backgroundColor: "#fff",
    borderRadius: theme.radii.lg,
    padding: theme.spaces.lg,
  },
  footer: {
    flexDirection: "row",
    marginTop: theme.spaces.md,
    gap: theme.spaces.md,
  },
  footerButton: {
    flex: 1,
  },
  footerButtonDestructive: {
    backgroundColor: theme.colors.destructive,
  },
  footerButtonTextDestructive: {
    color: theme.colors.destructiveForeground,
  },
  title: {
    marginBottom: theme.spaces.md,
    textAlign: "center",
  },
}));
