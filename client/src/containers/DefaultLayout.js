import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import Navigation from '../components/Navigation'

class Default extends Component {
  static propTypes = {
    children: PropTypes.node
  }

  render() {
    let { children } = this.props
    return (
      <div className="row">
        <div className="col-sm-3 sidebar">
          <Navigation/>
        </div>
        <div className="col-sm-9 content">
          {children}
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => ({
})

export default connect(mapStateToProps, {
})(Default)
