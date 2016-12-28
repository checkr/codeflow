import React, { PropTypes } from 'react'
import { Link, IndexLink } from 'react-router'

export default class NavItem extends React.Component {
  static contextTypes = {
    router: PropTypes.object
  }

  render() {
    const { router } = this.context
    const { index, onlyActiveOnIndex, to, children, hideWhenActive, classNames, ...props } = this.props
    const isActive = router.isActive(to, onlyActiveOnIndex)
    const LinkComponent = index ? IndexLink : Link
    if (isActive && hideWhenActive) {
      return null
    }
    return (
      <li className={(isActive ? 'active ' : '') + classNames} key={index}> <LinkComponent to={to} {...props}>{children}</LinkComponent>
      </li>
    )
  }
}
