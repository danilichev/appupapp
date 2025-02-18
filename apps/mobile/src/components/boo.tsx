import { StyleSheet, Text, View } from "react-native";

export const Boo = () => (
  <View style={styles.root}>
    <Text>Boo</Text>
  </View>
);

const styles = StyleSheet.create({
  root: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
  },
});
