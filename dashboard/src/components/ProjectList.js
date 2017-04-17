import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'
import { fetchProjects } from '../actions'
import _ from 'underscore'
import merge from 'lodash/merge'
import queryString  from 'query-string'
import FontAwesome from 'react-fontawesome'
import { createBookmarks, deleteBookmarks } from '../actions'
import Pagination from './Pagination'

const loadData = props => {
  props.fetchProjects(props.routing.search)
}

class ProjectList extends Component {
  static contextTypes = {
    router: PropTypes.object
  }

  static propTypes = {
    createBookmarks: PropTypes.func.isRequired,
    deleteBookmarks: PropTypes.func.isRequired
  }

  componentWillMount() {
    loadData(this.props)
  }

  search(term) {
    const { pathname, query } = this.props.routing
    const { router } = this.context
    const urlParams = merge({}, query, { q: term })

    router.push({
      pathname,
      query: urlParams
    })

    this.paginate(queryString.stringify(urlParams))
  }

  paginate(search) {
    this.props.fetchProjects(`?${search}`, 'pagination')
  }

  handleBookmarkCreate = (project, e) => {
    e.preventDefault()
    this.props.createBookmarks({ projectId: project._id })
  }

  handleBookmarkDelete = (project, e) => {
    e.preventDefault()
    this.props.deleteBookmarks({ projectId: project._id })
  }

  render() {
    const { projects } = this.props

    if (_.isEmpty(projects)) {
      return null
    }

    const { records, pagination } = this.props.projects
    const { bookmarks } = this.props
    let projects_jsx = []
    let that = this
    records.forEach(function (project) {
      let star = 'star-o'
      let action = that.handleBookmarkCreate
      let bookmark = _.findWhere(bookmarks, { projectId: project._id })
      if (!_.isEmpty(bookmark)) {
        star = 'star'
        action = that.handleBookmarkDelete
      }
      projects_jsx.push(
        <li className="list-group-item" key={project._id}>
          <div className="feed-element">
            <Link to={'/projects/'+project.slug}>
              <div className="media-body ">
                <strong>{project.name}</strong> <FontAwesome onClick={(e) => action(project, e)} className="float-xs-right" name={star}/>
              </div>
            </Link>
          </div>
        </li>
      )
    })

    return (
      <div>
        <div className="clearfix row mb-2">
          <div className="col-md-5 pull-right">
            <div className="input-group">
              <span className="input-group-addon"><i className="fa fa-search fa-fw"></i></span>
              <input type="text"
                placeholder="Search projects"
                autoFocus="true"
                className="form-control"
                onChange={({ target }) => this.search(target.value)}
                value={this.props.routing.query['q']}/>
            </div>
          </div>
        </div>
        <ul className="list-group search-results">
          {projects_jsx}
        </ul>
        <Pagination onChange={(_p, s) => this.paginate(s)} totalPages={pagination.totalPages} page={pagination.current} count={pagination.recordsOnPage} queryParam="projects_page"/>
      </div>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  projects: state.projects,
  bookmarks: state.bookmarks,
  routing: state.routing.locationBeforeTransitions
})

export default connect(mapStateToProps, {
  createBookmarks,
  deleteBookmarks,
  fetchProjects
})(ProjectList)
