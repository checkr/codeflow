/* eslint-disable no-console */

import React, { Component } from 'react'
import { connect } from 'react-redux'
import { isEmpty } from 'lodash'
import ProjectSettingsForm from './SettingsForm'
import { fetchProjectSettings, updateProjectSettings, deleteProject } from '../../actions'

class ProjectSettings extends Component {
  componentWillMount() {
    this.loadData(this.props)
  }

  loadData = props => {
    props.fetchProjectSettings(props.project.slug)
  }

  handleSubmit = () => {
    const { project } = this.props

    return this.props.updateProjectSettings(project.slug, this.props.projectSettingsForm.values).then((fetch) => {
      console.log(fetch)
    })
  }

  render() {
    const { project, projectSettings, deleteProject } = this.props

    if (isEmpty(project)) {
      return null
    }

    return (
      <div>
        <ProjectSettingsForm initialValues={projectSettings} onSubmit={this.handleSubmit} deleteProject={deleteProject} project={project}/>
      </div>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  projectSettingsForm: state.form.projectSettings,
  projectSettings: state.projectSettings
})

export default connect(mapStateToProps, {
  fetchProjectSettings,
  updateProjectSettings,
  deleteProject
})(ProjectSettings)
