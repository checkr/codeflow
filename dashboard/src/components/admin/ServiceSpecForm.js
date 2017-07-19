
import React, { Component } from 'react'
import { Form, Button, Modal, ModalBody, ModalHeader, ModalFooter, ListGroup, ListGroupItem } from 'reactstrap'
import { Field, reduxForm } from 'redux-form'
import { forEach, isEmpty } from 'lodash'
import { connect } from 'react-redux'

//validations
const validate = values => {
  const errors = {}
  if (!values.cpuLimit) {
    errors.cpuLimit = '* Required ^'
  }
  if (!values.cpuRequest) {
    errors.cpuRequest = '* Required ^'
  }
  if (!values.memoryLimit) {
    errors.memoryRequest = '* Required ^'
  }
  if (!values.memoryRequest) {
    errors.memoryRequest = '* Required ^'
  }
  if (!/[0-9]+m$/.test(values.cpuLimit)) {
    errors.cpuLimit = '* Invalid value for CPU. Example: 500m'
  }
  if (!/[0-9]+m$/.test(values.cpuRequest)) {
    errors.cpuRequest = '* Invalid value for CPU. Example 500m'
  }

  if ((!/[0-9]+Mi$/.test(values.memoryRequest)) && (!/[0-9]+Gi$/.test(values.memoryRequest))) {
    errors.memoryRequest = '* Invalid. Must end with "Mi" or Gi"'
  }
  if ((!/[0-9]+Mi$/.test(values.memoryLimit)) && (!/[0-9]+Gi$/.test(values.memoryLimit))) {
    errors.memoryLimit = '* Invalid. Must end with "Mi" or Gi"'
  }

  if (!/^[0-9]+$/.test(values.terminationGracePeriodSeconds)) {
    errors.terminationGracePeriodSeconds = "Invalid. Must be numeric only."
  }
  return errors
}

class ServiceSpec extends Component {
  constructor(props) {
    super(props);
    this.state = {
      serviceToggle: false
    }
    this.serviceToggle = this.serviceToggle.bind(this)
    this.renderServices = this.renderServices.bind(this)
  }

  renderField({label, icon, input, type, meta: {touched, error, warning}}) { 
    return (
      <div className="input-group">
        <div className="col-xs-6">
          <i className={icon} aria-hidden="true"></i> {label}
        </div>
        <div className="col-xs-6">
          <input className="form-control" {...input} type={type} />
          {touched &&
            ((error &&
              <span>
                {error}
              </span>) ||
              (warning &&
                <span>
                  {warning}
                </span>))}
        </div>
      </div>
    )
  }

  serviceToggle() {
    this.setState({
      serviceToggle: !this.state.serviceToggle
    })
  }

  render() {
      const { initialValues, onCancel, onDelete, onSubmit, invalid, serviceSpecServices, projects } = this.props
      const deleteEnabled = (serviceSpecServices.length === 0)
      return(
        <div>
          <Form onSubmit={onSubmit} noValidate >
            <div className="col-xs-12">
              <div className="col-xs-12">
                <div className="form-group">
                  <div className="row">
                    <h5>
                    <Field name="name" label="Name" component={this.renderField} type="text" icon="fa fa-server" />
                    </h5>
                  </div>
                  <div className="row">
                    <Field name="cpuRequest" label="CPU Request" component={this.renderField} type="text" icon="fa fa-tachometer" />
                  </div>
                  <div className="row">
                    <Field name="cpuLimit" label="CPU Limit" component={this.renderField} type="text" icon="fa fa-tachometer" />
                  </div>
                  <div className="row">
                    <Field name="memoryRequest" label="Memory Request" component={this.renderField} type="text" icon="fa fa-tachometer" />
                  </div>
                  <div className="row">
                    <Field name="memoryLimit" label="Memory Limit" component={this.renderField} type="text" icon="fa fa-tachometer" />
                  </div>
                  <div className="row">
                    <Field name="terminationGracePeriodSeconds" label="Termination Grace Period (seconds)" component={this.renderField} type="text" icon="fa fa-clock-o" />
                  </div>
                </div>
              </div>
            </div>
            <div className="col-xs-12 input-group form-group">
              <button type="button" className="btn btn-secondary btn-sm float-xs-right btn-service-action-right" onClick={() => onCancel()}>
                <i className="fa fa-times" aria-hidden="true" /> Cancel
              </button>
              <button type="button" className="btn btn-danger btn-sm float-xs-right btn-service-action-right" disabled={!deleteEnabled} onClick={() => onDelete()}>
                <i className="fa fa-trash" aria-hidden="true" /> Delete
              </button>
              <button type="submit" className="btn btn-success btn-sm float-xs-right btn-service-action-right" disabled={invalid}>
                <i className="fa fa-check" aria-hidden="true" /> Save
              </button>
              <button type="button" className="btn btn-warning" onClick={() => this.serviceToggle() }>Used by {serviceSpecServices.length} running services</button>
            </div>
          </Form>
          <Modal isOpen={this.state.serviceToggle} toggle={this.serviceToggle}>
            <ModalHeader>
              <i className="fa fa-server" /> {initialValues.name} usage
            </ModalHeader>
            <ModalBody>
              <ListGroup>
                { this.renderServices(serviceSpecServices, projects) }
              </ListGroup>
            </ModalBody>
            <ModalFooter>
              <Button color="secondary" onClick={this.serviceToggle}>Done</Button>
            </ModalFooter>
          </Modal>
        </div>
      )
  }
  
  renderServices(services, projects) {
    let jsx = []
    if (isEmpty(services) || isEmpty(projects)) {
      return jsx
    }
    let projectMap = {}
    forEach(projects.records, p => {
      projectMap[p._id] = p.name
    })
    forEach(services, service => {
      if (!projectMap[service.projectId]) {
        return
      }
      const newName = projectMap[service.projectId].replace(/\//g, "-")
      const href = "/projects/" + newName + "/resources"
      jsx.push(
        <ListGroupItem key={service._id} className="justify-content-between" tag="a" href={href} action>
          <div className="row">
            <div className="col-xs-10">
              {projectMap[service.projectId]} {service.name}
            </div>
            <div className="col-xs-2">
              <i className="fa fa-times" /><b>{service.count}</b>
            </div>
          </div>
        </ListGroupItem>
      )
    })
    return (jsx)
  }
  
}

const ServiceSpecForm = reduxForm({
  form: 'serviceSpec',
  validate
})(ServiceSpec)

export default connect( state => {
  const formValues = state.form.serviceSpec
  return { formValues }
})(ServiceSpecForm)
