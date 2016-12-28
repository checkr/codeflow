import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

class Empty extends Component {
  static propTypes = {
    children: PropTypes.node
  }

  render() {
    let { children } = this.props
    return (
      <div className="row">
        <div className="col-sm-12">
          {children}
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => ({
})

export default connect(mapStateToProps, {
})(Empty)
