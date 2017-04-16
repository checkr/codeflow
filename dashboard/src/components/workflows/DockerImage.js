import React, { Component } from 'react'
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap'
import { connect } from 'react-redux'
import { fetchBuild, updateBuild } from '../../actions'

class DockerImage extends Component {
  loadData = props => {
    props.fetchBuild(props.project.slug, props.workflow.releaseId)
  }

  constructor(props) {
    super(props)
    this.state = {
      modal: false
    }

    this.toggle = this.toggle.bind(this)
  }

  toggle(e) {
    e.preventDefault()

    this.loadData(this.props)

    this.setState({
      modal: !this.state.modal
    })
  }

  onConfirm(e) {
    e.preventDefault()
    this.props.updateBuild(this.props.project.slug, this.props.workflow.releaseId)
    this.setState({
      modal: false
    })
  }

  onCancel(e) {
    e.preventDefault()
    this.setState({
      modal: false
    })
  }

  renderBtn() {
    let { workflow } = this.props

    switch(workflow.state) {
      case 'running':
        return <Button className="tag tag-warning flow" key={workflow._id} onClick={(e) => this.toggle(e)}><i className="fa fa-refresh fa-spin fa-lg fa-fw" aria-hidden="true"></i> {workflow.type.toUpperCase()}:{workflow.name}</Button>
      case 'building':
        return <Button className="tag tag-warning flow" key={workflow._id} onClick={(e) => this.toggle(e)}><i className="fa fa-refresh fa-spin fa-lg fa-fw" aria-hidden="true"></i> {workflow.type.toUpperCase()}:{workflow.name}</Button>
      case 'complete':
        return <span className="tag tag-success flow" key={workflow._id}><i className="fa fa-check fa-lg" aria-hidden="true"></i> {workflow.type.toUpperCase()}:{workflow.name}</span>
      case 'failed':
        return <Button className="tag tag-danger flow" key={workflow._id} onClick={(e) => this.toggle(e)}><i className="fa fa-times fa-lg" aria-hidden="true"></i> {workflow.type.toUpperCase()}:{workflow.name}</Button>
      default:
        return <Button className="tag tag-info flow" key={workflow._id} onClick={(e) => this.toggle(e)}><i className="fa fa-cog fa-spin fa-lg fa-fw" aria-hidden="true"></i> {workflow.type.toUpperCase()}:{workflow.name}</Button>
    }
  }

  render() {
    let { build } = this.props
    return (
      <span>
        {this.renderBtn()}
        <Modal isOpen={this.state.modal} toggle={this.toggle} className={this.props.className}>
          <ModalHeader toggle={this.toggle}>Docker Image Build Log</ModalHeader>
          <ModalBody>
            { build.buildLog !== "" && <pre>{build.buildLog}</pre>}
            { build.buildError !== "" && <code>{build.buildError}</code>}
          </ModalBody>
          <ModalFooter>
            <Button color="primary" onClick={(e) => this.onConfirm(e)}>Rebuild & Redeploy</Button>{' '}
            <Button color="secondary" onClick={(e) => this.onCancel(e)}>Close</Button>
          </ModalFooter>
        </Modal>
      </span>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  project: state.project,
  build: state.build,
})

export default connect(mapStateToProps, {
  fetchBuild, updateBuild
})(DockerImage)
