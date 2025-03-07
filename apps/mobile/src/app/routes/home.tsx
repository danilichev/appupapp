import { View } from "react-native";
import { StyleSheet } from "react-native-unistyles";

import { Button, ButtonIcon, Text } from "@appupapp/mobile-ui";

export const HomeScreen = () => {
  return (
    <View style={styles.container}>
      <Text variant="h1">Log In</Text>
      <Button onPress={() => console.log("Pressed")} variant="primary">
        <ButtonIcon name="user" /> Log In
      </Button>
    </View>
  );
};

const styles = StyleSheet.create((theme) => ({
  container: {
    alignItems: "center",
    backgroundColor: theme.colors.background,
    flex: 1,
    justifyContent: "center",
  },
}));
