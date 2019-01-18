import React from 'react'
import { Route, Switch } from 'react-router' // react-router v4

import VivoSearchApp from './containers/VivoSearchApp'

const routes =
<Switch>
  <Route exact path='/' component={VivoSearchApp}/>
</Switch>

export default routes
