import React, { Component } from 'react'
import { connect } from 'react-redux'
import _ from 'underscore'

class ProjectAccess extends Component {
  render() {
    const { project } = this.props

    if (_.isEmpty(project)) {
      return null
    }

    return (
      <div>
        <div className="hr-divider m-t-md m-b">
          <h3 className="hr-divider-content hr-divider-heading">Access</h3>
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => ({})

export default connect(mapStateToProps, {

})(ProjectAccess)
