import { Text, View } from "react-native";
import { StyleSheet } from "react-native-unistyles";

import { Button, ButtonIcon } from "@appupapp/mobile-ui";

export const HomeScreen = () => {
  return (
    <View style={styles.container}>
      <Button onPress={() => console.log("Pressed")} variant="outline">
        <ButtonIcon name="user" /> Log In
      </Button>
      <Text style={{ fontSize: 16, fontWeight: "bold" }}>Log In</Text>
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
