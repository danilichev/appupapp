import { Text, View } from "react-native";
import { StyleSheet } from "react-native-unistyles";

import { Button } from "@appupapp/mobile-ui";

export const HomeScreen = () => {
  return (
    <View style={styles.container}>
      <Text>Home Screen</Text>
      <Button
        containerStyle={styles.button}
        loadingPosition="end"
        onPress={() => console.log("Pressed")}
        textStyle={styles.buttonText}
        variant="outline"
        size="lg"
        isLoading
      >
        Press me
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
  button: {
    // backgroundColor: theme.colors.secondary,
  },
  buttonText: {
    // color: theme.colors.secondaryForeground,
  },
}));
