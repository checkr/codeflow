/* eslint-disable no-console */

import React, { Component } from 'react'
import { connect } from 'react-redux'
import _ from 'underscore'
import ProjectSettingsForm from '../components/ProjectSettingsForm'
import { fetchProjectSettings, updateProjectSettings } from '../actions'

class ProjectSettings extends Component {
  componentWillMount() {
    this.loadData(this.props)
  }

  loadData = props => {
    props.fetchProjectSettings(props.project.slug)
  }

  handleSubmit = e => {
    e.preventDefault()
    const { project } = this.props
    this.props.updateProjectSettings(project.slug, this.props.projectSettingsForm.values).then((fetch) => {
      console.log(fetch)
    })
  }

  render() {
    const { project, projectSettings } = this.props

    if (_.isEmpty(project)) {
      return null
    }

    return (
      <div>
        <ProjectSettingsForm initialValues={projectSettings} onSubmit={this.handleSubmit}/>
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => ({
  projectSettingsForm: state.form.projectSettings,
  projectSettings: state.projectSettings
})

export default connect(mapStateToProps, {
  fetchProjectSettings,
  updateProjectSettings
})(ProjectSettings)
