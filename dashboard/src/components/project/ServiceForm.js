import React, { Component } from 'react'
import { Alert, Form, FormGroup, FormFeedback, Input } from 'reactstrap'
import { Field, FieldArray, reduxForm } from 'redux-form'
import { Tooltip } from 'reactstrap';
import ButtonConfirmAction from '../ButtonConfirmAction'
import { toSafeInteger, without } from 'lodash'
import { forEach } from 'lodash'

const MIN_PORT = 1
const MAX_PORT = 65535

const validatePortRange = value => {
  const number = toSafeInteger(value)

  if ( number < MIN_PORT || number > MAX_PORT ) {
    return `Invalid port. Must be between ${MIN_PORT} and less than ${MAX_PORT}`
  }
}

const renderInput = ({input, _meta, ...field}) => {
  let _field = without(field, ["meta"])
  return (
    <Input {..._field} {...input} type={field.type} placeholder={field.placeholder} disabled={field.disabled} />
  )
}

const normalizeInt = (value, _previousValue) => {
  return parseInt(value, 10)
}

const normalizeBool = (value, _previousValue) => {
    if(value === ""){
        return false
    }
    return true
}

const InputWithWarnings = props => {
  const { error, touched } = props.meta

  return (
    <FormGroup color={error ? "danger" : ""}>
      {renderInput(props)}
      {touched && error && <FormFeedback>{error}</FormFeedback>}
    </FormGroup>
  )
}

const renderListeners = ({ fields }) => (
  <div className="col-xs-12">
    { (fields.length > 0) && <label>Container ports</label>}
    {fields.map((service, i) =>
    <div className="row" key={i}>
      <div className="col-xs-4">
        <Field name={'listeners['+i+'].port'} component={InputWithWarnings} validate={validatePortRange} type="number" step="1" min={MIN_PORT} max={MAX_PORT} normalize={normalizeInt}/>
      </div>
      <div className="col-xs-4">
        <div className="form-group service-protocol">
          <label className="form-check-inline">
            <Field className="form-check-input" name={'listeners['+i+'].protocol'} component={renderInput} type="radio" value="TCP"/> TCP
          </label>
          <label className="form-check-inline">
            <Field className="form-check-input" name={'listeners['+i+'].protocol'} component={renderInput} type="radio" value="UDP"/> UDP
          </label>
        </div>
      </div>
      <div className="col-xs-4">
        <button type="button" className="btn btn-secondary btn-sm float-xs-right btn-service-action-right" onClick={() => fields.remove(i)}>
          <i className="fa fa-times" aria-hidden="true" />
        </button>
      </div>
    </div>
    )}
    <div className="row">
      <div className="col-xs-2" style={{ position: 'absolute', zIndex: 100 }}>
        <button type="button" className="btn btn-secondary btn-sm float-xs-left btn-service-action" onClick={() => fields.push({ protocol: 'TCP' })}>Add container port </button>
      </div>
    </div>
  </div>
  )

const renderSpecsSelect = field => {
  return (
    <select {...field.input} name={field.name} className="form-control">
      <option disabled value="">Choose service spec</option>
      {field.serviceSpecs.map(function (spec) {
        return <option key={spec._id} value={spec._id}>{spec.name}</option>
      })}
    </select>
  )
}

class ProjectServiceForm extends Component {
  constructor(props) {
    super(props);

    this.toggleContainerPort = this.toggleContainerPort.bind(this);
    this.toggleTooltipServiceSettings = this.toggleTooltipServiceSettings.bind(this);
    this.state = {
      tooltipContainerPortOpen: false,
      tooltipServiceSettingsOpen: false
    };
  }

  toggleContainerPort() {
    this.setState({
      tooltipContainerPortOpen: !this.state.tooltipContainerPortOpen
    })
  }

  toggleTooltipServiceSettings() {
    this.setState({
      tooltipServiceSettingsOpen: !this.state.tooltipServiceSettingsOpen
    })
  }

  renderServiceSpecsTip() {
    let { serviceSpecs } = this.props
    let specs_jsx = []

    forEach(serviceSpecs, spec => {
      specs_jsx.push(this.renderQuickSpec(spec))
    })
    return specs_jsx
  }

  renderQuickSpec(spec) {
    return (
      <div style={{ align: 'left' }}>

        <div className="hr-divider">
          <h3 className="hr-divider-content hr-divider-heading">{spec.name}</h3>
        </div>
          CPU: {spec.cpuRequest} / {spec.cpuLimit}<br/>
          Memory: {spec.memoryRequest} / {spec.memoryLimit}<br/>
          Termination: {spec.terminationGracePeriodSeconds}s
      <br/>
      </div>
    )
  }

  render() {
    const { edit, onCancel, onDelete, handleSubmit, error } = this.props
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
                      <label>Name</label>
                      <Field name="name" className="form-control" component={renderInput} disabled={edit} type="text"/>
                    </div>
                  </div>
                  <div className="col-xs-4">
                      <div className="form-group">
                        <label>Spec</label>
                        <i className="fa fa-question-circle" id="ToolTipServiceSettings" aria-hidden="true" style={{ position: 'left', zIndex: 100, bottom: '-20px', left: '175px' }}></i>

                        <Field className="form-control" name="specId" serviceSpecs={this.props.serviceSpecs} component={renderSpecsSelect}/>                      
                      </div>
                      <Tooltip placement="bottom" isOpen={this.state.tooltipServiceSettingsOpen} target="ToolTipServiceSettings" toggle={this.toggleTooltipServiceSettings}>
                        {this.renderServiceSpecsTip()}
                      </Tooltip>
                  </div>
                  <div className="col-xs-2">
                    <div className="form-group">
                      <label>Count</label>
                      <Field name="count" component={renderInput} type="number" step="1" min="0" max="100" normalize={normalizeInt}/>
                    </div>
                  </div>
                </div>
                <div className="form-group">
                  <label>Command</label>
                  <Field name="command" className="form-control" component={renderInput} type="text"/>
                </div>

                <div className="form-group">
                  <Field className="form-check-input" name="oneShot" component={renderInput} type="checkbox" value="false" normalize={normalizeBool} /> One-shot
                </div>

                <div className="row">
                  <div className="col-xs-12">
                    <i className="fa fa-question-circle" id="ToolTipContainerPort" aria-hidden="true" style={{ position: 'absolute', zIndex: 100, bottom: '-20px', left: '175px' }}></i>
                    <Tooltip placement="right" isOpen={this.state.tooltipContainerPortOpen} target="ToolTipContainerPort" toggle={this.toggleContainerPort}>
                      If your application is a webserver then add the port that it listens on.
                    </Tooltip>
                  </div>
                </div>
                <div className="row">
                  <FieldArray name="listeners" component={renderListeners}/>
                </div>
              </div>
              <div className="col-xs-12">
                <button type="button" className="btn btn-secondary btn-sm float-xs-right btn-service-action-right" onClick={() => onCancel()}>
                  <i className="fa fa-times" aria-hidden="true" /> Cancel
                </button>
                { edit &&
                <ButtonConfirmAction btnLabel="Delete" btnIconClass="fa fa-trash" onConfirm={onDelete} btnClass="btn btn-danger btn-sm float-xs-right btn-service-action-right">
                  Are you sure?
                </ButtonConfirmAction>}
                <button type="submit" className="btn btn-success btn-sm float-xs-right btn-service-action-right">
                  <i className="fa fa-check" aria-hidden="true" /> Save
                </button>
              </div>
            </div>
          </FormGroup>
        </Form>
    )
  }
}

export default reduxForm({
  enableReinitialize: true,
  destroyOnUnmount: true,
  form: 'projectService'
})(ProjectServiceForm)
