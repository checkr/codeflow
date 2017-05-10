import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { createNewProject } from '../../actions'
import AddProjectForm from './Form'

class AddProjectsPage extends Component {
  static propTypes = {
    createNewProject: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object
  }

  handleSubmit = e => {
    e.preventDefault()
    const { router } = this.context
    this.props.createNewProject(this.props.addProjectForm.values).then((fetch) => {
      if (fetch.payload.slug) {
        router.push('/projects/'+fetch.payload.slug)
      }
    })
  }

  render() {
    return (
      <div>
        <AddProjectForm initialValues={{"bookmarked": true}} onSubmit={this.handleSubmit}/>
      </div>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  addProjectForm: state.form.addProject
})

export default connect(mapStateToProps, {
  createNewProject
})(AddProjectsPage)
