import React, { Component } from 'react'
import { connect } from 'react-redux'
import { ListGroup } from 'reactstrap'
import { fetchBookmarkCurrentRelease, removeBookmarkCurrentReleases } from '../../actions'
import { map, tap, isEqual } from 'lodash'

import Release from '../project/Release'

const getProjectIds = (bookmarks) => {
  return tap({}, (obj) => {
    map(bookmarks, 'projectId').forEach((projectId) => obj[projectId] = true)
  })
}

class DashboardReleases extends Component {
  componentDidMount() {
    const { bookmarks, fetchBookmarkCurrentRelease } = this.props
    bookmarks.forEach(({slug}) => fetchBookmarkCurrentRelease({slug}))
  }

  componentWillReceiveProps({bookmarks}) {
    const { bookmarks: currentBookmarks, fetchBookmarkCurrentRelease } = this.props

    const currentProjectIds = getProjectIds(currentBookmarks)
    const nextProjectIds = getProjectIds(bookmarks)

    if (isEqual(currentProjectIds, nextProjectIds)) {
      return
    }
    bookmarks.forEach(({slug}) => fetchBookmarkCurrentRelease({slug}))
  }

  componentWillUnmount() {
    this.props.removeBookmarkCurrentReleases()
  }

  render() {
    const {bookmarks, bookmarkReleases} = this.props
    return (
      <div>
        <div className="hr-divider m-t-md m-b">
          <h3 className="hr-divider-content hr-divider-heading">Recent releases</h3>
        </div>
        <ListGroup>
          {bookmarks.map(({projectId, slug, name}) => <Release key={projectId} projectName={name} projectSlug={slug} release={bookmarkReleases[projectId]} showHeader />)}
        </ListGroup>
      </div>
    )
  }
}

const mapStateToProps = ({bookmarks, bookmarkReleases}) => ({
  bookmarks,
  bookmarkReleases
})

const ConnectedDashboardReleases = connect(mapStateToProps, {
  fetchBookmarkCurrentRelease,
  removeBookmarkCurrentReleases
})(DashboardReleases)

export default ConnectedDashboardReleases
