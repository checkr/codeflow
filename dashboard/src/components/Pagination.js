import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { Pagination as Pagination1, PaginationItem, PaginationLink } from 'reactstrap'
import merge from 'lodash/merge'
import queryString  from 'query-string'

class Pagination extends Component {
  static contextTypes = {
    router: PropTypes.object
  }

  // handler for clicks on page buttons
  // calls back to owner with requested page
  handlePageClick(e) {
    e.preventDefault()
    const { router } = this.context
    let query = {}

    query[this.props.queryParam] = e.target.text
    let urlParams = merge({}, this.props.routing.query, query)

    router.push({
      pathname: this.props.routing.pathname,
      query: urlParams
    })

    this.props.onChange(this.props.routing.pathname, queryString.stringify(urlParams))
  }

  // handler for clicks on Previous button
  // calls back to owner with requested page
  handlePrevClick(e) {
    e.preventDefault()
    const { router } = this.context
    let pageIndex = this.props.routing.query[this.props.queryParam] ? parseInt(this.props.routing.query[this.props.queryParam],10) : 1
    let query = {}

    query[this.props.queryParam] = pageIndex - 1
    let urlParams = merge({}, this.props.routing.query, query)

    router.push({
      pathname: this.props.routing.pathname,
      query: urlParams
    })

    this.props.onChange(this.props.routing.pathname, queryString.stringify(urlParams))
  }

  // handler for clicks on Next button
  // calls back to owner with requested page
  handleNextClick(e) {
    e.preventDefault()
    const { router } = this.context
    let pageIndex = this.props.routing.query[this.props.queryParam] ? parseInt(this.props.routing.query[this.props.queryParam],10) : 1
    let query = {}
    query[this.props.queryParam] = pageIndex + 1
    let urlParams = merge({}, this.props.routing.query, query)

    router.push({
      pathname: this.props.routing.pathname,
      query: urlParams
    })

    this.props.onChange(this.props.routing.pathname, queryString.stringify(urlParams))
  }

  // render Previous button enabled or disabled based on props
  renderPreviousButton() {
    return this.props.page === 1 ?
      <PaginationItem className="disabled"><PaginationLink previous disabled href="#" /></PaginationItem>:
      <PaginationItem><PaginationLink previous onClick={(e) => this.handlePrevClick(e)} href="#"/></PaginationItem>
  }

  // render all page buttons with active page based on props
  renderPageButtons() {
    let buttons = []
    for (var i=1; i<=this.props.totalPages; i++) {
      if (this.props.totalPages > 7) {
        if (i === 4) {
          if (this.props.page > 3 && this.props.page < this.props.totalPages - 3) {
            buttons.push(<PaginationItem key={i} className="active"><PaginationLink disabled href="#">... {this.props.page} ...</PaginationLink></PaginationItem>)
          } else {
            buttons.push(<PaginationItem key={i}><PaginationLink disabled href="#">...</PaginationLink></PaginationItem>)
          }
        }
        if (i > 3 && i < this.props.totalPages - 3) {
          continue
        }
      }
      buttons.push(i === (this.props.page) ?
        <PaginationItem key={i} className="active"><PaginationLink disabled href="#">{i}</PaginationLink></PaginationItem>:
        <PaginationItem key={i}><PaginationLink href="#" onClick={(e) => this.handlePageClick(e)}>{i}</PaginationLink></PaginationItem>
      )
    }

    return buttons
  }

  // render Next button enabled or disabled based on props
  renderNextButton() {
    return this.props.page === this.props.totalPages ?
        <PaginationItem className="disabled"><PaginationLink next href="#"/></PaginationItem>:
        <PaginationItem><PaginationLink next href="#" onClick={(e) => this.handleNextClick(e)}/></PaginationItem>
  }

  // main render for component
  render() {
    if (this.props.page === 1 && this.props.totalPages === 1) {
      return null
    }

    return (
      <Pagination1 size="sm float-xs-right">
        {this.renderPreviousButton()}
        {this.renderPageButtons()}
        {this.renderNextButton()}
      </Pagination1>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  routing: state.routing.locationBeforeTransitions
})

export default connect(mapStateToProps, {
})(Pagination)
