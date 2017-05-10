import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { authCallback, resetErrorMessage } from '../../actions'
import loadConfig from "../../config"

import OktaSignIn from '../../../public/okta-signin-widget-1.9.0/js/okta-sign-in.min.js';
import '../../../public/okta-signin-widget-1.9.0/css/okta-sign-in.css';
import '../../../public/okta-signin-widget-1.9.0/css/okta-theme.css';


class OktaLoginPage extends Component {
  static propTypes = {
    authCallback: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object
  }

  componentWillMount() {
    if (window.Backbone && window.Backbone.history) {
      window.Backbone.history.stop()
    }
  }

  componentWillUnmount() {
    window.Backbone.history.stop()
  }

  componentDidMount() {
    const { router } = this.context
    const CONFIG = loadConfig()
    var orgUrl = CONFIG.REACT_APP_OKTA_URL
    var oktaSignIn = new OktaSignIn({
      logo: CONFIG.REACT_APP_OKTA_LOGO,
      baseUrl: orgUrl,
      redirectUri: window.location.origin,
      // OpenID Connect options
      clientId: CONFIG.REACT_APP_OKTA_CLIENT_ID,
      features: {
        rememberMe: true,
        autoPush: true,
        selfServiceUnlock: true,
      },
      authParams: {
        responseType: 'id_token',
        responseMode: 'okta_post_message',
        scopes: [
          'openid',
          'email',
          'profile',
          'address',
          'phone',
          'groups'
        ]
      }
    })

    var that = this
    oktaSignIn.renderEl(
      { el: '#okta-login-container' },
      function (res) {
        // res.idToken - id_token generated
        // res.claims - decoded id_token information

        that.props.authCallback('/auth/callback/okta', { idToken: res.idToken }).then(() => {
          resetErrorMessage()
          router.push(that.props.next)
        })
      },
      function (err) {
        console.log(err) // eslint-disable-line no-console
      }
    )
  }

  render() {
    return (
      <div className="row">
        <div className="col-sm-12 content">
          <div id="okta-login-container" style={{ paddingTop: '50px' }} />
        </div>
      </div>
    )
  }
}

const mapStateToProps = (_state, _ownProps) => ({})

export default connect(mapStateToProps, {
  authCallback,
  OktaSignIn
})(OktaLoginPage)
