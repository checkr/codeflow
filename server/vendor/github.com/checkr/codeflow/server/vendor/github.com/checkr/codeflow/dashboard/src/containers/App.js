import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { resetErrorMessage, wsConnect } from '../actions'

class App extends Component {
  static propTypes = {
    errorMessage: PropTypes.string,
    resetErrorMessage: PropTypes.func.isRequired,
    children: PropTypes.node
  }

  componentWillMount() {
    this.props.wsConnect()
  }

  handleDismissClick = e => {
    this.props.resetErrorMessage()
    e.preventDefault()
  }

  renderErrorMessage() {
    const { errorMessage } = this.props

    if (!errorMessage) {
      return null
    }

    return (
      <p style={{ backgroundColor: '#e99', padding: 10 }}>
        <b>{errorMessage}</b>
        {' '}
        (<a href="#"
          onClick={this.handleDismissClick}
        >
          Dismiss
        </a>)
      </p>
    )
  }

  render() {
    const { children } = this.props
    return (
      <div>
        <nav className="navbar navbar-full navbar-light bg-faded">
          <a className="navbar-brand" href="#">
            <img src="/images/code.svg" width="30" height="30" className="d-inline-block align-top" alt=""/> Codeflow
          </a>
        </nav>
        <div className="container">
          <div>
            {this.renderErrorMessage()}
            {children}
          </div>
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => ({
  errorMessage: state.errorMessage,
  inputValue: ownProps.location.pathname.substring(1)
})

export default connect(mapStateToProps, {
  resetErrorMessage,
  wsConnect
})(App)
