import React from 'react'
import moment from 'moment'

const renderFeatureHash = (feature) => {
  if(feature.externalLink && feature.externalLink !== '' && feature.externalLink.startsWith('http')) {
    return (<a href={feature.externalLink} target="_blank">{feature.hash.substring(0,8)}</a>)
  }

  return (<span>{feature.hash.substring(0,8)}</span>)
}

const Feature = ({feature, handleDeploy, includedClass, isFeatureHovered, onMouseEnter, onMouseLeave}) => (
  <li className={"list-group-item" + includedClass} key={feature.hash} onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
    <div className="feed-element">
      <div className="row media-body">
        <div className="col-xs-10">
          <strong>{renderFeatureHash(feature)} - {feature.message}</strong> <br/>
          <small className="text-muted">by <strong>{feature.user}</strong> {moment(feature.created).fromNow() } - {moment(feature.created).format('MMMM Do YYYY, h:mm:ss A')} </small>
        </div>
        <div className="col-xs-2 flex-xs-middle">
          {isFeatureHovered &&
          <button type="button" className="btn btn-secondary btn-sm float-xs-right" onClick={handleDeploy}>Deploy</button> }
        </div>
      </div>
    </div>
  </li>
)

export default Feature
