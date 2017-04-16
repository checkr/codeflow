import React, { PropTypes } from 'react'
import { connect } from 'react-redux'
import { fetchRefreshToken } from '../actions'

export function RequireAuthentication(Component) {
  class AuthenticatedComponent extends React.Component {
    static propTypes = {
      fetchRefreshToken: PropTypes.func.isRequired
    }
    static contextTypes = {
      router: PropTypes.object
    }

    componentWillMount() {
      this.checkAuth(this.props.isAuthenticated)
    }

    componentWillReceiveProps(nextProps) {
      this.checkAuth(nextProps.isAuthenticated)
    }

    checkAuth(isAuthenticated) {
      const { router } = this.context
      if (!isAuthenticated) {
        let redirectAfterLogin = this.props.location.pathname
        router.push(`/login?next=${redirectAfterLogin}`)
      }
    }

    render() {
      return (
        <div>
          {this.props.isAuthenticated === true ? <Component {...this.props}/>: null}
        </div>
      )
    }
  }

  const mapStateToProps = (state, _ownProps) => ({
    token: state.auth.token,
    userName: state.auth.userName,
    isAuthenticated: state.auth.isAuthenticated,
    refreshToken: state.auth.refreshToken
  })

  return connect(mapStateToProps, {
    fetchRefreshToken
  })(AuthenticatedComponent)
}
