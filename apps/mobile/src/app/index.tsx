import * as React from 'react';
import { registerRootComponent } from 'expo';

import { Router } from './router';

const App = () => {
  return <Router />;
};

registerRootComponent(App);
