import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Nav, NavLink } from 'reactstrap'
import NavItem from '../navigation/Item'
import { isEmpty } from 'lodash'

import ProjectInit      from './Init'
import ProjectResources from './Resources'
import ProjectSettings  from './Settings'
import ProjectAccess    from './Access'
import ProjectDeploy    from './Deploy'
import Release          from './Release'
import { fetchProject } from '../../actions'

const loadData = props => {
  props.fetchProject(props.params.project_slug)
}

class Project extends Component {
  constructor(props) {
    super(props)
    this.toggle = this.toggle.bind(this)
    this.state = {
      dropdownOpen: false
    }
  }

  componentWillMount() {
    loadData(this.props)
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.params.project_slug !== this.props.params.project_slug) {
      loadData(nextProps)
    }
  }

  toggle() {
    this.setState({
      dropdownOpen: !this.state.dropdownOpen
    })
  }

  renderSection(project) {
    switch(this.props.params.section) {
      case 'deploy':
        return (<div>
          <ProjectInit project={project}/>
          <ProjectDeploy project={project}/>
        </div>)
      case 'resources':
        return (<ProjectResources project={project}/>)
      case 'access':
        return (<ProjectAccess project={project}/>)
      case 'settings':
        return (<ProjectSettings project={project}/>)
      default:
        return null
    }
  }

  render() {
    const { project_slug } = this.props.params
    const { project } = this.props

    if (isEmpty(project)) {
      return null
    }

    return (
      <div>
        <Nav pills>
          <NavItem to={'/projects/' + project_slug + '/deploy'} classNames="nav-item">Deploy</NavItem>
          <NavItem to={'/projects/' + project_slug + '/resources'} classNames="nav-item">Resources</NavItem>
          <NavItem to={'/projects/' + project_slug + '/settings'} classNames="nav-item">Settings</NavItem>
          <NavLink href={project.logsUrl} target="_blank" className="float-xs-right">Logs</NavLink>
        </Nav>

        {this.renderSection(project)}
      </div>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  project: state.project
})

export default connect(mapStateToProps, {
  fetchProject
})(Project)

export {
  Release
}
