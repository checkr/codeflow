import React, { Component } from 'react'
import { connect } from 'react-redux'
import ServiceSpecList from './ServiceSpecList'
import AdminActionsList from './AdminActionsList'

class Admin extends Component {
  render() {
    return (
      <div>
        <ServiceSpecList specs={this.props.serviceSpecs}/>
        <AdminActionsList/>
      </div>
    )
  }
}

const mapStateToProps = (_state, _ownProps) => ({})

export default connect(mapStateToProps, {
})(Admin)