import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { authCallback, resetErrorMessage } from '../actions'

class OktaLoginPage extends Component {
  static propTypes = {
    authCallback: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object
  }

  componentWillUnmount() {
    window.Backbone.history.stop()
  }

  componentDidMount() {
    const { router } = this.context

    var orgUrl = 'https://checkr.okta.com'
    var oktaSignIn = new window.OktaSignIn({
      logo: 'https://ok4static.oktacdn.com/bc/image/fileStoreRecord?id=fs0pgyx2uFHva5qsw1t6',
      baseUrl: orgUrl,
      // OpenID Connect options
      clientId: 'TJxx1X61RTCF8uxNpxll',
      authParams: {
        responseType: 'id_token',
        responseMode: 'okta_post_message',
        scope: [
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

        that.props.authCallback('/oauth2/callback/okta', { idToken: res.idToken }).then(() => {
          var next = that.props.location.query.next ? that.props.location.query : '/'
          resetErrorMessage()
          router.push(next)
        })
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

const mapStateToProps = (state, ownProps) => ({})

export default connect(mapStateToProps, {
  authCallback
})(OktaLoginPage)
