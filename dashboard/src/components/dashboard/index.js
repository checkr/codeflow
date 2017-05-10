import React, { Component } from 'react'
import { connect } from 'react-redux'

import DashboardStats    from './Stats'
import DashboardReleases from './Releases'

class DashboardPage extends Component {
  render() {
    return (
      <div>
        <DashboardStats/>
        <DashboardReleases/>
      </div>
    )
  }
}

const mapStateToProps = (_state, _ownProps) => ({})

export default connect(mapStateToProps, {
})(DashboardPage)
