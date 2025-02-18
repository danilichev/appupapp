import * as React from 'react';
import { View, Text } from 'react-native';
import { createStaticNavigation } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { Foo } from '@appupapp/mobile-ui-restyle';

function HomeScreen() {
  return (
    <View style={{ flex: 1, alignItems: 'center', justifyContent: 'center' }}>
      <Text>Home Screen</Text>
      <Foo />
    </View>
  );
}

const RootStack = createNativeStackNavigator({
  screens: {
    Home: HomeScreen,
  },
});

export const Router = createStaticNavigation(RootStack);
