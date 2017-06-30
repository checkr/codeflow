import React from 'react'
import { Router, Route, IndexRoute, hashHistory, IndexRedirect } from 'react-router'
import App from './containers/App'
import DefaultLayout from './containers/DefaultLayout'
import EmptyLayout from './containers/EmptyLayout'
import Dashboard from './components/dashboard'
import Browse from './components/browse'
import Add from './components/add'
import Project from './components/project'
import Login from './components/login'
import Admin from './components/admin'
import { RequireAuthentication } from './components/AuthenticatedComponent'

export default <Router history={hashHistory}>
  <Route path="/" component={App}>
    <Route component={RequireAuthentication(DefaultLayout)}>
      <Route path="admin" component={Admin} />
      <IndexRoute component={Dashboard} />
      <Route path="browse" component={Browse} />
      <Route path="add" component={Add} />
      <Route path="projects/:project_slug" component={Project}>
        <Route path=":section" component={Project} />
        <IndexRedirect to="deploy" />
      </Route>
    </Route>
    <Route component={EmptyLayout}>
      <Route path="login" component={Login} />
    </Route>
  </Route>
</Router>
