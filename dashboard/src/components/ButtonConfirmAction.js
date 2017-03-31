import React, { Component } from 'react'
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap'

class ButtonConfirmAction extends Component {
  constructor(props) {
    super(props)
    this.state = {
      modal: false
    }

    this.toggle = this.toggle.bind(this)
  }

  toggle(e) {
    e.preventDefault()
    this.setState({
      modal: !this.state.modal
    })
  }

  onConfirm(e) {
    e.preventDefault()
    this.setState({
      modal: false
    })
    this.props.onConfirm()
  }

  onCancel(e) {
    e.preventDefault()
    this.setState({
      modal: false
    })
    this.props.onCancel()
  }
  
  renderBtn() {
    if (this.props.btnIconClass !== '') {
      return (
        <Button className={this.props.btnClass} onClick={(e) => this.toggle(e)}>
          <i className={this.props.btnIconClass} aria-hidden="true" /> {this.props.btnLabel}
        </Button>
      ) 
    }

    return (
      <Button className={this.props.btnClass} onClick={(e) => this.toggle(e)}>
        <i className={this.props.btnIconClass} aria-hidden="true" /> {this.props.btnLabel}
      </Button>
    )
  }

  render() {
    return (
      <div>
        {this.renderBtn()}
        <Modal isOpen={this.state.modal} toggle={this.toggle} className={this.props.className}>
          <ModalHeader toggle={this.toggle}>Confirmation needed</ModalHeader>
          <ModalBody>
            {this.props.children}
          </ModalBody>
          <ModalFooter>
            <Button color="primary" onClick={(e) => this.onConfirm(e)}>Confirm</Button>{' '}
            <Button color="secondary" onClick={(e) => this.onCancel(e)}>Cancel</Button>
          </ModalFooter>
        </Modal>
      </div>
    )
  }
}

export default ButtonConfirmAction
