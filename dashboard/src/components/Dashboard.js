import React, { Component } from 'react'
import { connect } from 'react-redux'

import DashboardStats from '../components/DashboardStats'

class DashboardPage extends Component {
  render() {
    return (
      <div>
        <DashboardStats/>
      </div>
    )
  }
}

const mapStateToProps = (_state, _ownProps) => ({})

export default connect(mapStateToProps, {
})(DashboardPage)
