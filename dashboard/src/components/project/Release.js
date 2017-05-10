import React from 'react'
import { Link } from 'react-router'
import moment from 'moment'
import {
  ListGroupItem,
  ListGroupItemHeading
} from 'reactstrap'

const renderFeatureHash = (feature) => {
  if(feature.externalLink && feature.externalLink !== '' && feature.externalLink.startsWith('http')) {
    return (<a href={feature.externalLink} target="_blank">{feature.hash.substring(0,8)}</a>)
  }

  return (<span>{feature.hash.substring(0,8)}</span>)
}

const Release = ({className, projectName, projectSlug, showHeader, release, actionBtn}) => {
  if (!release) {
    return null
  }

  return (
    <ListGroupItem className={className} >
      <div className="feed-element media-body">
        { showHeader && <div className="row">
          <div className="col-xs-12">
            <ListGroupItemHeading><Link to={`projects/${projectSlug}/deploy`}>{projectName}</Link></ListGroupItemHeading>
          </div>
        </div> }
        <div className="row">
          <div className="col-xs-10">
            <div className="row">
              <div className="col-xs-12">
                <i className="fa fa-code-fork" aria-hidden="true" /> <strong>
                   <i className="fa fa-angle-double-right" aria-hidden="true" /> { renderFeatureHash(release.tailFeature) } - {release.headFeature.message}
                </strong> <br/>
              </div>
              <div className="col-xs-12">
                <small className="text-muted">by <strong>{release.headFeature.user}</strong> {moment(release._created).fromNow() } - {moment(release._created).format('MMMM Do YYYY, h:mm:ss A')} </small>
              </div>
              <div className="col-xs-12">
              </div>
            </div>
          </div>
          <div className="col-xs-2 flex-xs-middle">
            {actionBtn}
          </div>
        </div>
      </div>
    </ListGroupItem>
  )
}

export default Release
