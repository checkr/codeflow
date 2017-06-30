import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import NavItem from './Item'
import { logoutUser, fetchUsers, fetchBookmarks } from '../../actions'
import { isEmpty } from 'lodash'

const loadData = props => {
  props.fetchUsers('me')
  props.fetchBookmarks('me')
}

class Navigation extends Component {
  static propTypes = {
    logoutUser: PropTypes.func.isRequired
  }

  componentWillMount() {
    loadData(this.props)
  }

  handleLogoutUserClick = e => {
    e.preventDefault()
    this.props.logoutUser()
  }

  renderBookmarks() {
    let bookmarks_jsx = []
    let { bookmarks } = this.props

    if (bookmarks.length > 0) {
      bookmarks_jsx.push(<li key="header-bookmarks" className="nav-header">Bookmarks</li>)
      bookmarks.forEach(bookmark =>
        bookmarks_jsx.push(<NavItem to={'/projects/'+bookmark.slug} classNames="" key={bookmark.projectId}>{bookmark.name}</NavItem>)
      )
    } else {
      return null
    }

    if (bookmarks_jsx.length === 0) {
      return null
    }

    return bookmarks_jsx
  }

  renderUser() {
    let { me } = this.props
    let user_jsx = []

    if (isEmpty(me)) {
      return null
    }

    user_jsx.push(<li className="nav-header" key="header-me">{me.name}</li>)
    user_jsx.push(<NavItem to="/logout" onClick={this.handleLogoutUserClick} key="me-logout">Logout</NavItem>)
    return user_jsx
  }

  renderAdmin() {
    //return null // not implemented
    let admin_jsx = []
    admin_jsx.push(<li className="nav-header" key="liadminsettings">Admin</li>)
    admin_jsx.push(<NavItem to="/admin" key="adminsettings">Settings</NavItem>)
    //admin_jsx.push(<NavItem to="/admin/users">Users</NavItem>)
    //admin_jsx.push(<NavItem to="/admin/teams">Teams</NavItem>)
    return admin_jsx
  }

  render() {
    return (
      <nav className="sidebar-nav">
        <div className="collapse nav-toggleable-sm" id="nav-toggleable-sm">
          <ul className="nav nav-pills nav-stacked">
            <NavItem to="/" onlyActiveOnIndex>Dashboard</NavItem>
            {this.renderAdmin()}

            <li className="nav-header">Projects</li>
            <NavItem to="/add">Add new</NavItem>
            <NavItem to="/browse">Browse</NavItem>

            {this.renderBookmarks()}

            {this.renderUser()}

          </ul>
        </div>
      </nav>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  me: state.me,
  bookmarks: state.bookmarks,
})

export default connect(mapStateToProps, {
  logoutUser,
  fetchUsers,
  fetchBookmarks,
})(Navigation)
