import React, { Component } from 'react'
import { connect } from 'react-redux'
import ServiceSpecList from './ServiceSpecList'

class Admin extends Component {
  render() {
    return (
      <div>
        <ServiceSpecList specs={this.props.serviceSpecs}/>
      </div>
    )
  }
}

const mapStateToProps = (_state, _ownProps) => ({})

export default connect(mapStateToProps, {
})(Admin)