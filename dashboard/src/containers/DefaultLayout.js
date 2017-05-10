import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import Navigation from '../components/navigation'

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

const mapStateToProps = () => ({})

export default connect(mapStateToProps, {
})(Default)
