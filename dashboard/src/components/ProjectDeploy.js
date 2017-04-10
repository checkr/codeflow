import React, { Component } from 'react'
import _ from 'underscore'
import { connect } from 'react-redux'
import { Alert } from 'reactstrap'
import { fetchProjectFeatures, createProjectRelease, createProjectRollbackTo, fetchProjectReleases, fetchProjectCurrentRelease } from '../actions'
import moment from 'moment'
import Pagination from './Pagination'
import DockerImage from './workflows/DockerImage'

const loadData = props => {
  if (props.project.slug) {
    props.fetchProjectFeatures(props.project.slug, props.routing.search)
    props.fetchProjectReleases(props.project.slug, props.routing.search)
    props.fetchProjectCurrentRelease(props.project.slug)
  }
}

class ProjectDeploy extends Component {
  constructor(props) {
    super(props)
    this.state = {
      featureHover: null,
      releaseHover: null,
    }
  }

  componentWillMount() {
    loadData(this.props)
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.project.slug !== '' && nextProps.project.slug !== this.props.project.slug) {
      loadData(nextProps)
    }
    if (nextProps.features.dirty) {
      nextProps.fetchProjectFeatures(nextProps.project.slug, nextProps.routing.search)
    }
    if (nextProps.releases.dirty) {
      nextProps.fetchProjectReleases(nextProps.project.slug, nextProps.routing.search)
    }
  }

  paginateFeatures(pathname, search) {
    this.props.fetchProjectFeatures(this.props.project.slug, '?'+search, 'pagination')
  }

  paginateReleases(pathname, search) {
    this.props.fetchProjectReleases(this.props.project.slug, '?'+search, 'pagination')
  }
  
  onDeployFeature(feature, e) {
    e.preventDefault()
    this.props.createProjectRelease(this.props.project.slug, feature)
  }

  onRollbackTo(release, e) {
    e.preventDefault()
    this.props.createProjectRollbackTo(this.props.project.slug, release)
  }
  
  renderFeatureHash(feature) {
    if(feature.externalLink && feature.externalLink !== '' && feature.externalLink.startsWith('http')) {
      return (<a href={feature.externalLink} target="_blank">{feature.hash.substring(0,8)}</a>)
    }

    return (<span>{feature.hash.substring(0,8)}</span>)
  }
  
  onFeatureMouseEnterHandler(id) {
    this.setState({featureHover: id});
  }

  onFeatureMouseLeaveHandler() {
    this.setState({featureHover: null});
  }

  renderFeatures() {
    let { records, pagination } = this.props.features
    let features_jsx = [] 

    if (_.isEmpty(records)) {
      return(
        <Alert color="info">
          This project has no deployable features. Do some work and come back!
        </Alert>
      )
    }

    let includedClass = ""
    records.forEach(feature => {
      if (this.state.featureHover === feature._id) {
        includedClass = " feature-included"
      }
      features_jsx.push(
        <li className={"list-group-item" + includedClass} key={feature.hash} onMouseEnter={() => this.onFeatureMouseEnterHandler(feature._id)} onMouseLeave={() => this.onFeatureMouseLeaveHandler()}>
          <div className="feed-element">
            <div className="row media-body">
              <div className="col-xs-10">
                <strong>{this.renderFeatureHash(feature)} - {feature.message}</strong> <br/>
                <small className="text-muted">by <strong>{feature.user}</strong> {moment(feature.created).fromNow() } - {moment(feature.created).format('MMMM Do YYYY, h:mm:ss A')} </small>
              </div>
              <div className="col-xs-2 flex-xs-middle">
                {this.state.featureHover === feature._id && 
                <button type="button" className="btn btn-secondary btn-sm float-xs-right" onClick={(e) => this.onDeployFeature(feature, e)}>Deploy</button> }
              </div>
            </div>
          </div>
        </li>
      )
    })

    return (
      <div>
        <ul className="list-group">{features_jsx}</ul>
        <Pagination onChange={(p,s) => this.paginateFeatures(p,s)} totalPages={pagination.totalPages} page={pagination.current} count={pagination.recordsOnPage} queryParam="features_page"/>
      </div>
    )
  }
  

  renderCurrentReleaseActions(release) {
    return <button type="button" key="btn" className="btn btn-secondary btn-sm float-xs-right" onClick={(e) => this.onDeployFeature(release.headFeature, e)}>Redeploy</button>
  }

  renderReleaseActions(release) {
    let jsx = [] 
    switch(release.state) {
      case 'waiting':
        jsx.push(<i key="waiting" className="fa fa-circle-o-notch fa-spin fa-fw float-xs-right" />)
        break
      case 'running':
        jsx.push(<i key="running" className="fa fa-refresh fa-spin fa-fw float-xs-right" />)
        break
      case 'failed':
        jsx.push(<i key="failed" className="fa fa-exclamation-triangle fa-fw float-xs-right" style={{ color: 'red' }} aria-hidden="true" />)
        break
      default:
        jsx.push(<button type="button" key="btn" className="btn btn-secondary btn-sm float-xs-right" onClick={(e) => this.onRollbackTo(release, e)}>Rollback</button>)
    }
    return jsx
  }

  renderReleaseWorkflow(release) {
    if (release.state === "complete") {
      return null 
    }

    let flows = _.map(release.workflow, (wf) => {
      if (wf.name === "DockerImage") {
        return <DockerImage key={wf._id} workflow={wf}/>
      }

      switch(wf.state) {
        case 'running':
          return <span className="tag tag-warning flow" key={wf._id}><i className="fa fa-refresh fa-spin fa-lg fa-fw" aria-hidden="true"></i> {wf.type.toUpperCase()}:{wf.name}</span> 
        case 'building':
          return <span className="tag tag-warning flow" key={wf._id}><i className="fa fa-refresh fa-spin fa-lg fa-fw" aria-hidden="true"></i> {wf.type.toUpperCase()}:{wf.name}</span> 
        case 'complete':
          return <span className="tag tag-success flow" key={wf._id}><i className="fa fa-check fa-lg" aria-hidden="true"></i> {wf.type.toUpperCase()}:{wf.name}</span> 
        case 'failed':
          return <span className="tag tag-danger flow" key={wf._id}><i className="fa fa-times fa-lg" aria-hidden="true"></i> {wf.type.toUpperCase()}:{wf.name}</span> 
        default:
          return <span className="tag tag-info flow" key={wf._id}><i className="fa fa-cog fa-spin fa-lg fa-fw" aria-hidden="true"></i> {wf.type.toUpperCase()}:{wf.name}</span> 
      }
    })
    
    return (
      <div>
        <div className="row">
          <div className="col-xs-12 flows">
            {flows}
          </div>
        </div>
      </div>
    )
  }

  onReleaseMouseEnterHandler(release) {
    this.setState({releaseHover: release._id});
  }

  onReleaseMouseLeaveHandler() {
    this.setState({releaseHover: null});
  }

  renderReleases() {
    let { records, pagination } = this.props.releases
    let releases_jsx = [] 
    let that = this 
    
    if (_.isEmpty(records)) {
      return (<li>
        <Alert color="info">
          This project has no releases. Build some features and deploy them!
        </Alert>
      </li>) 
    }

    records.forEach((release, i) => {
      releases_jsx.push(
        <li className="list-group-item" key={release._id} onMouseEnter={() => this.onReleaseMouseEnterHandler(release)} onMouseLeave={() => this.onReleaseMouseLeaveHandler()}>
          <div className="feed-element">
            <div className="row media-body">
              <div className="col-xs-10">
                <div className="row">
                  <div className="col-xs-12">
                    <i className="fa fa-code-fork" aria-hidden="true" /> <strong>
                      {this.renderFeatureHash(release.tailFeature)} <i className="fa fa-angle-double-right" aria-hidden="true" /> {this.renderFeatureHash(release.headFeature)} - {release.headFeature.message}
                    </strong> <br/>
                  </div>
                  <div className="col-xs-12">
                    <small className="text-muted">by <strong>{release.headFeature.user}</strong> {moment(release._created).fromNow() } - {moment(release._created).format('MMMM Do YYYY, h:mm:ss A')} </small>
                  </div>
                  <div className="col-xs-12">
                    {this.renderReleaseWorkflow(release)}
                  </div>
                </div>
              </div>
              <div className="col-xs-2 flex-xs-middle">
                {that.renderReleaseActions(release)}
              </div>
            </div>
          </div>
        </li>
      )
    })

    return (
      <div>
        <ul className="list-group">{releases_jsx}</ul>
        <Pagination onChange={(p,s) => this.paginateReleases(p,s)} totalPages={pagination.totalPages} page={pagination.current} count={pagination.recordsOnPage} queryParam="releases_page"/>
      </div>
    )
  }

  renderCurrentRelease() {
    let { currentRelease } = this.props
    let releases_jsx = [] 
    releases_jsx.push(
      <li className="list-group-item" key={currentRelease._id} onMouseEnter={() => this.onReleaseMouseEnterHandler(currentRelease)} onMouseLeave={() => this.onReleaseMouseLeaveHandler()}>
        <div className="feed-element">
          <div className="row media-body">
            <div className="col-xs-10">
              <div className="row">
                <div className="col-xs-12">
                  <i className="fa fa-code-fork" aria-hidden="true" /> <strong>
                    {this.renderFeatureHash(currentRelease.tailFeature)} <i className="fa fa-angle-double-right" aria-hidden="true" /> {this.renderFeatureHash(currentRelease.headFeature)} - {currentRelease.headFeature.message}
                  </strong> <br/>
                </div>
                <div className="col-xs-12">
                  <small className="text-muted">by <strong>{currentRelease.headFeature.user}</strong> {moment(currentRelease._created).fromNow() } - {moment(currentRelease._created).format('MMMM Do YYYY, h:mm:ss A')} </small>
                </div>
                <div className="col-xs-12">
                  {this.renderReleaseWorkflow(currentRelease)}
                </div>
              </div>
            </div>
            <div className="col-xs-2 flex-xs-middle">
              {this.renderCurrentReleaseActions(currentRelease)}
            </div>
          </div>
        </div>
      </li>
      )

    return (
      <div>
        <ul className="list-group">{releases_jsx}</ul>
      </div>
    )
  }

  render() {
    const { project, currentRelease } = this.props

    if (_.isEmpty(project) || !project.pinged) {
      return null
    }

    return (
      <div>
        <div className="clearfix">
          <div className="hr-divider m-t-md m-b">
            <h3 className="hr-divider-content hr-divider-heading">Features</h3>
          </div>
          {this.renderFeatures()}
        </div>

        { !_.isEmpty(currentRelease._id) && <div className="clearfix">
          <div className="hr-divider m-t-md m-b">
            <h3 className="hr-divider-content hr-divider-heading">Current Release</h3>
          </div>
          {this.renderCurrentRelease()}
        </div> }

        <div className="clearfix">
          <div className="hr-divider m-t-md m-b">
            <h3 className="hr-divider-content hr-divider-heading">Releases</h3>
          </div>

          <ul className="list-group">
            {this.renderReleases()}
          </ul>
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => ({
  features: state.features,
  releases: state.releases,
  currentRelease: state.currentRelease,
  routing: state.routing.locationBeforeTransitions
})

export default connect(mapStateToProps, {
  fetchProjectFeatures,
  createProjectRelease,
  createProjectRollbackTo,
  fetchProjectReleases,
  fetchProjectCurrentRelease
})(ProjectDeploy)
