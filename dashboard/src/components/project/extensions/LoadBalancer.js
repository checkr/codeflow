import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Form, FormGroup, Input, Alert } from 'reactstrap'
import { Field, FieldArray, reduxForm } from 'redux-form'
import { isEmpty, get, find } from 'lodash'
import ButtonConfirmAction from '../../ButtonConfirmAction'
import { Tooltip } from 'reactstrap';

const renderInput = field => {
  return (
    <Input {...field.input} type={field.type} placeholder={field.placeholder} />
  )
}

const renderSelect = field => {
  return (
    <select {...field.input} name={field.name} className="form-control">
      <option disabled value="">Choose service</option>
      {field.services.map(function (service) {
        return <option key={service._id} value={service._id}>{service.name}</option>
      })}
    </select>
  )
}

const renderProtocolSelect = field => {
  return (
    <select {...field.input} name={field.name} className="form-control">
      <option disabled value="">Choose protocol</option>
      <option key="http" value="HTTP">HTTP</option>
      <option key="https" value="HTTPS">HTTPS</option>
      <option key="ssl" value="SSL">SSL</option>
      <option key="tcp" value="TCP">TCP</option>
      <option key="udp" value="UDP">UDP</option>
    </select>
  )
}

const renderListenerSelect = field => {
  if (isEmpty(field.service.listeners)) {
    return null
  }
  return (
    <select {...field.input} name={field.name} className="form-control">
      <option disabled value="" key="">Choose service listener</option>
      {field.service.listeners.map(function (listener) {
        return <option key={listener.port} value={listener.port}>{listener.port}</option>
      })}
    </select>
  )
}

const normalizeInt = (value, _previousValue) => {
  return parseInt(value, 10)
}

const renderListeners = ({ fields, service, tooltipServiceProtocolOpen, toggleServiceProtocol }) => {
  return (
  <div className="col-xs-12">
    { !isEmpty(service) && <hr />}
    { (fields.length > 0) &&
    <div className="row">
      <div className="col-xs-4">
        <label>Port</label>
      </div>
      <div className="col-xs-4">
        <label>Container Port</label>
      </div>
      <div className="col-xs-3">
        <label>Service Protocol</label> <i className="fa fa-question-circle" id="ToolTipServiceProtocol" aria-hidden="true"/>
        <Tooltip placement="right" isOpen={tooltipServiceProtocolOpen} target="ToolTipServiceProtocol" toggle={toggleServiceProtocol}>
          <b>HTTPS</b>: Serve HTTPS using a pre-existing certificate.<br/>
          <b>HTTP</b>:  Plain HTTP (no SSL).<br/>
          <b>SSL</b>:   Serve SSL encrypted TCP using a pre-existing certificate.  (Example: for use with websocket protocol wss://).<br/>
          <b>TCP</b>:   Plain TCP (no SSL).<br/>
          <b>UDP</b>:   Use UDP.  Only available for Internal service type.<br/>
        </Tooltip>
      </div>
      <div className="col-xs-1" />
    </div>
    }
    {fields.map((s, i) =>
    <div className="row" key={i}>
      <div className="col-xs-4">
        <div className="form-group">
          <Field name={'listenerPairs['+i+'].source.port'} component={renderInput} type="number" step="1" min="0" max="65535" normalize={normalizeInt}/>
        </div>
      </div>
      <div className="col-xs-4">
        <div className="form-group">
          <Field className="form-control" name={'listenerPairs['+i+'].destination.port'} service={service} component={renderListenerSelect} normalize={normalizeInt}/>
        </div>
      </div>
      <div className="col-xs-3">
        <Field className="form-control" name={'listenerPairs['+i+'].destination.protocol'} component={renderProtocolSelect}/>
      </div>
      <div className="col-xs-1">
        <button type="button" className="btn btn-secondary btn-sm float-xs-right btn-service-action-right" onClick={() => fields.remove(i)}>
          <i className="fa fa-times" aria-hidden="true" />
        </button>
      </div>
    </div>
    )}
    {!isEmpty(service) && <div className="row">
      <div className="col-xs-2">
        <button type="button" className="btn btn-secondary btn-sm float-xs-left btn-service-action" onClick={() => fields.push({destination: {protocol: "HTTPS"}})}>Add port map</button>
      </div>
    </div>}
  </div>
  )
}

class LoadBalancer extends Component {
  renderShow() {
    let { services, extension } = this.props
    let service = {}
    if (!isEmpty(extension) && extension.serviceId) {
      service = find(services, { _id: extension.serviceId })
    }
    let { dns } = ""

    if (!isEmpty(extension.dns)) {
      dns = <pre><i className="fa fa-globe" aria-hidden="true" /> {extension.dns}</pre>
    }

    if (!isEmpty(extension.subdomain) && !isEmpty(extension.fqdn)) {
      dns = <pre><i className="fa fa-globe" aria-hidden="true" /> {extension.subdomain}.{extension.fqdn}</pre>
    }

    return (
      <div className="extension-lb">
        { !isEmpty(service) && <div><strong>{extension.type}</strong> load balancer <i className="fa fa-angle-double-right" aria-hidden="true"></i> <strong>{service.name}</strong></div> }
        {dns}
      </div>
    )
  }

  renderEdit() {
    const { handleSubmit, onCancel, onDelete, formValues, error } = this.props
    let { services } = this.props
    let service = {}

    if (!isEmpty(formValues) && formValues.values.serviceId) {
      service = find(services, { _id: formValues.values.serviceId })
    }
    return (
      <Form onSubmit={handleSubmit} noValidate>
        { error &&
        <Alert color="danger">{error}</Alert>
        }
        <FormGroup>
          <div className="row">
            <div className="col-xs-12">
              <div className="row">
                <div className="col-xs-6">
                  <div className="form-group">
                    <label>Service</label>
                    <Field className="form-control" name="serviceId" services={services} component={renderSelect}/>
                  </div>
                </div>
                <div className="col-xs-6">
                  <div className="form-group">
                    <label>Subdomain (<strong>{get(formValues, 'values.subdomain')}</strong>.example.net)</label>
                    <Field name="subdomain" component={renderInput} type="text" placeholder="api"/>
                  </div>
                </div>
                <div className="col-xs-8">
                  <div className="form-group">
                    <label>Access</label>  <i className="fa fa-question-circle" id="ToolTipAccess" aria-hidden="true"></i>
                    <div className="form-group service-protocol">
                      <label className="form-check-inline">
                        <Field className="form-check-input" name="type" component={renderInput} type="radio" value="internal"/> Internal
                      </label>
                      <label className="form-check-inline">
                        <Field className="form-check-input" name="type" component={renderInput} type="radio" value="office"/> Office
                      </label>
                      <label className="form-check-inline">
                        <Field className="form-check-input" name="type" component={renderInput} type="radio" value="external"/> External
                      </label>
                      <div>
                        <Tooltip placement="right" isOpen={this.state.tooltipAccessOpen} target="ToolTipAccess" toggle={this.toggleAccess}>
                          <b>Internal</b>: Internal to kubernetes.  Normally used if other services need to connect to this service.<br/>
                          <b>Office</b>: Exposes the service to the Office network (or VPN).<br/>
                          <b>External</b>: Exposes the service to the public internet!<br/>
                        </Tooltip>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="row">
                <FieldArray name="listenerPairs" props={{ service: service, tooltipServiceProtocolOpen: this.state.tooltipServiceProtocolOpen, toggleServiceProtocol: this.toggleServiceProtocol }} component={renderListeners}/>
              </div>
            </div>
            <div className="col-xs-12">
              <button type="button" className="btn btn-secondary btn-sm float-xs-right btn-service-action-right" onClick={() => onCancel()}>
                <i className="fa fa-times" aria-hidden="true" /> Cancel
              </button>
              <ButtonConfirmAction btnLabel="Delete" btnIconClass="fa fa-check" onConfirm={onDelete} btnClass="btn btn-danger btn-sm float-xs-right btn-service-action-right">
                Are you sure?
              </ButtonConfirmAction>
              <ButtonConfirmAction btnLabel="Save" btnIconClass="fa fa-check" onConfirm={handleSubmit} btnClass="btn btn-success btn-sm float-xs-right btn-service-action-right">
                Are you sure?
              </ButtonConfirmAction>
            </div>
          </div>
        </FormGroup>
      </Form>
    )
  }

  constructor(props) {
    super(props);

    this.toggleAccess = this.toggleAccess.bind(this);
    this.toggleServiceProtocol = this.toggleServiceProtocol.bind(this);
    this.state = {
      tooltipServiceProtocolOpen: false,
      tooltipAccessOpen: false
    };
  }

  toggleAccess() {
    this.setState({
      tooltipAccessOpen: !this.state.tooltipAccessOpen
    })
  }

  toggleServiceProtocol() {
    this.setState({
      tooltipServiceProtocolOpen: !this.state.tooltipServiceProtocolOpen
    })
  }

  render() {
    if (this.props.edit) {
      return this.renderEdit()
    } else {
      return this.renderShow()
    }
  }
}

const LoadBalancerForm = reduxForm({
  enableReinitialize: true,
  destroyOnUnmount: false,
  form: 'projectExtension'
})(LoadBalancer)

export default connect(
  state => {
    const formValues = state.form.projectExtension
    return { formValues: formValues }
  }
)(LoadBalancerForm)
