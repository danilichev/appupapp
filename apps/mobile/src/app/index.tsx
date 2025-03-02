import "./theme";

import { registerRootComponent } from "expo";

import { Router } from "./router";

export const App = () => {
  return <Router />;
};

registerRootComponent(App);
