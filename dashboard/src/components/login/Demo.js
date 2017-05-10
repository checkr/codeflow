import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { authCallback } from '../../actions'

class DevLoginPage extends Component {
  static propTypes = {
    authCallback: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object
  }

  componentDidMount() {
    const { router } = this.context
    this.props.authCallback('/auth/callback/demo', {}).then((resp) => {
      if (resp.type === "AUTH_SUCCESS") {
        router.push(this.props.next)
      }
    })
  }

  render() {
    return (
      <div className="row">
        <div className="col-sm-12 content">
          Devlopment Login
        </div>
      </div>
    )
  }
}

const mapStateToProps = (_state, _ownProps) => ({})

export default connect(mapStateToProps, {
  authCallback
})(DevLoginPage)
