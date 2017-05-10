import React, { Component } from 'react'
import { connect } from 'react-redux'

import DashboardStats    from '../components/DashboardStats'
import DashboardReleases from '../components/DashboardReleases'

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
