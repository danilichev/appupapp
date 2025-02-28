import { createStaticNavigation } from "@react-navigation/native";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import * as React from "react";
import { Text, View } from "react-native";

import { Button } from "@appupapp/mobile-ui";

function HomeScreen() {
  return (
    <View style={{ alignItems: "center", justifyContent: "center", flex: 1 }}>
      <Text>Home Screen</Text>
      <Button
        color="pink"
        loadingPosition="end"
        onPress={() => console.log("Pressed")}
      >
        Press me
      </Button>
    </View>
  );
}

const RootStack = createNativeStackNavigator({
  screens: {
    Home: HomeScreen,
  },
});

export const Router = createStaticNavigation(RootStack);
