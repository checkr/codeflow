import React from 'react'
import { Router, Route, IndexRoute, hashHistory, IndexRedirect } from 'react-router'
import App from './containers/App'
import DefaultLayout from './containers/DefaultLayout'
import EmptyLayout from './containers/EmptyLayout'
import Dashboard from './components/Dashboard'
import Projects from './components/Projects'
import AddProject from './components/AddProject'
import Project from './components/Project'
import OktaLogin from './components/OktaLogin'
import { RequireAuthentication } from './components/AuthenticatedComponent'

export default <Router history={hashHistory}>
  <Route path="/" component={App}>
    <Route component={RequireAuthentication(DefaultLayout)}>
      <IndexRoute component={Dashboard} />
      <Route path="projects" component={Projects} />
      <Route path="projects/add" component={AddProject} />
      <Route path="projects/:project_slug" component={Project}>
        <Route path=":section" component={Project} />
        <IndexRedirect to="deploy" />
      </Route>
    </Route>
    <Route component={EmptyLayout}>
      <Route path="login" component={OktaLogin} />
    </Route>
  </Route>
</Router>
