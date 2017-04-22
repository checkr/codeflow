import React, { Component } from 'react'
import { connect } from 'react-redux'
import { isEmpty } from 'lodash'

class ProjectAccess extends Component {
  render() {
    const { project } = this.props

    if (isEmpty(project)) {
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

const mapStateToProps = (_state, _ownProps) => ({})

export default connect(mapStateToProps, {

})(ProjectAccess)
