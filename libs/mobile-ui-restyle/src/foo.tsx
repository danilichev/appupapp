import { StyleSheet, Text, View } from 'react-native';

export const Foo = () => (
  <View style={styles.root}>
    <Text>Foo</Text>
  </View>
);

const styles = StyleSheet.create({
  root: {
    alignItems: 'center',
    backgroundColor: 'pink',
    flex: 1,
    justifyContent: 'center',
  },
});
