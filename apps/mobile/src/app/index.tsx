import { registerRootComponent } from "expo";
import * as React from "react";

import { Router } from "./router";

export const App = () => {
  return <Router />;
};

registerRootComponent(App);
