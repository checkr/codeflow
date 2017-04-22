import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Nav, NavLink, NavItem, ButtonDropdown, DropdownToggle, DropdownMenu, DropdownItem } from 'reactstrap'
import NavItemCustom from '../components/NavItem'
import ProjectInit from '../components/ProjectInit'
import ProjectResources from '../components/ProjectResources'
import ProjectSettings from '../components/ProjectSettings'
import ProjectAccess from '../components/ProjectAccess'
import ProjectDeploy from '../components/ProjectDeploy'
import { fetchProject } from '../actions'
import _ from 'underscore'

const loadData = props => {
  props.fetchProject(props.params.project_slug)
}

class ProjectPage extends Component {
  constructor(props) {
    super(props)
    this.toggle = this.toggle.bind(this)
    this.state = {
      dropdownOpen: false,
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
    console.log(this.state.dropdownOpen)
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

    if (_.isEmpty(project)) {
      return null
    }

    return (
      <div>
        <Nav pills>
          <NavItemCustom to={'/projects/' + project_slug + '/deploy'} className="nav-item">Deploy</NavItemCustom>
          <NavItemCustom to={'/projects/' + project_slug + '/resources'} className="nav-item">Resources</NavItemCustom>
          <NavItemCustom to={'/projects/' + project_slug + '/settings'} className="nav-item">Settings</NavItemCustom>
          <NavItem to="" className="nav-item">
            <NavLink href={project.logsUrl} target="_blank" className="nav-item">Logs</NavLink>
          </NavItem>
          <ButtonDropdown isOpen={this.state.dropdownOpen} toggle={this.toggle} className="float-xs-right">
            <DropdownToggle caret>
              PRODUCTION
            </DropdownToggle>
            <DropdownMenu>
              <DropdownItem>Development</DropdownItem>
              <DropdownItem>Staging</DropdownItem>
              <DropdownItem>Production</DropdownItem>
            </DropdownMenu>
          </ButtonDropdown>
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
})(ProjectPage)
