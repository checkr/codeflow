import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { Button, Form, FormGroup, Label, Input } from 'reactstrap'
import { Field, FieldArray, reduxForm, change } from 'redux-form'
import ButtonConfirmAction from '../ButtonConfirmAction'

const normalizeBool = (value, _previousValue) => {
    if(value === ""){
        return false
    }
    return true
}

class ProjectSettings extends Component {
  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props)

    this.state = { collapsed: true }
  }

  renderInput = (field) => {
    if(field['data-type'] === 'file') {
      return (
        <Input {...field.input} type='textarea' wrap="off" placeholder={field.placeholder} />
      )
    } else {
      return (
        <Input {...field.input} type='text' placeholder={field.placeholder} />
      )
    }
  }

  renderCheckbox(field) {
    return (
      <Input {...field.input} type={field.type} placeholder={field.placeholder} />
    )
  }

  renderSelect = (field) => {
    return (
      <select {...field.input} className="form-control">
        <option value="env">Env</option>
        <option value="protected-env">Protected Env</option>
        <option value="build">Build Arg</option>
        <option value="file">File</option>
      </select>
    )
  }

  renderDisabledInput = (field) => {
    return (
      <Input {...field.input} type={field.type} placeholder={field.placeholder} disabled/>
    )
  }

  addConfigVar = (fields) => {
    let { formValues } = this.props
    if (formValues.values) {
      let { newConfigVarType, newConfigVarKey, newConfigVarValue } = formValues.values
      if(newConfigVarKey && newConfigVarValue && newConfigVarKey !== '' && newConfigVarValue !== '') {
        if (!newConfigVarType || newConfigVarType === "") {
          newConfigVarType = "env"
        }
        fields.push({ type: newConfigVarType, key: newConfigVarKey, value: newConfigVarValue })
        this.props.dispatch(change('projectSettings', 'newConfigVarType', 'env'))
        this.props.dispatch(change('projectSettings', 'newConfigVarKey', ''))
        this.props.dispatch(change('projectSettings', 'newConfigVarValue', ''))
      }
    }
  }

  randomKey() {
    return Math.random().toString(36).substring(7)
  }

  deleteConfigVar = (secret, fields, i) => {
    let { formValues } = this.props
    if (formValues.values) {
      let { deletedSecrets } = this.props.formValues.values
      if(!deletedSecrets) {
        deletedSecrets = []
      }
      let ds = [].concat(deletedSecrets).concat(secret)
      fields.remove(i)
      this.props.dispatch(change('projectSettings', 'deletedSecrets', ds))
    }
  }

  renderConfigVars = ({ fields }) => {
    return (
      <div>
        <div className="hr-divider m-t-md m-b" onClick={() => this.onToggleConfigVars()}>
          <h3 className="hr-divider-content hr-divider-heading config-vars__heading">
            Config Variables <i className={`fa fa-caret-${this.state.collapsed ? 'up' : 'down'}`} aria-hidden="true"/>
          </h3>
        </div>
        <FormGroup className={`config-vars ${this.state.collapsed ? 'config-vars--collapsed' : ''}`}>
          {fields.map((s, i) => {
            let secret = this.props.formValues.values.secrets[i]
            return (<div key={i} className={'row flex-items-xs-middle config-var'}>
              <div className="col-xs-2">
                <Field name={'secrets['+i+'].type'} component={this.renderSelect}/>
              </div>
              <div className="col-xs-4">
                <Field name={'secrets['+i+'].key'} component={this.renderInput} type="text" placeholder="SOME_KEY"/>
              </div>
              <div className="col-xs-4">
                <Field name={'secrets['+i+'].value'} data-type={secret.type} component={this.renderInput} type="textarea" wrap="off" placeholder="******"/>
              </div>
              <div className="col-xs-2">
                <button type="button" className="btn btn-secondary btn-sm" onClick={() => this.deleteConfigVar(secret, fields, i)}>
                  <i className="fa fa-times" aria-hidden="true" />
                </button>
              </div>
            </div>)
          }
          )}

          { (fields.length > 0) && <hr/>}

          <div className="row flex-items-xs-middle config-var">
            <div className="col-xs-11">
              <button type="button" className="btn btn-secondary btn-sm float-xs-right" onClick={() => fields.push({type: "env"})}>
                <i className="fa fa-plus" aria-hidden="true" /> Add
              </button>
            </div>
          </div>
        </FormGroup>
        <div className="text-sm-center mb-1" hidden={!this.state.collapsed}>
          <button type="button" className="btn btn-secondary btn-sm" onClick={() => this.onToggleConfigVars()}>Reveal Variables</button>
        </div>
      </div>
    )
  }

  onProjectDelete() {
    this.props.deleteProject(this.props.project.slug).then(() => {
      window.location = "/"
    })
  }

  onProjectDeleteCancel() {

  }

  render() {
    const { pristine, submitting, handleSubmit } = this.props

    return (
      <div>
        <Form onSubmit={handleSubmit}>
          <div className="hr-divider m-t-md m-b">
            <h3 className="hr-divider-content hr-divider-heading">Project</h3>
          </div>
          <FormGroup check>
            <Label check>
              <Field name="continuousIntegration" component={this.renderCheckbox} type="checkbox" normalize={normalizeBool} /> Continuous Integration (CircleCI)
            </Label>
          </FormGroup>
          <FormGroup check>
            <Label check>
              <Field name="continuousDelivery" component={this.renderCheckbox} type="checkbox" normalize={normalizeBool} /> Continuous Delivery (Deploy on green)
            </Label>
          </FormGroup>
          <FormGroup>
            <Label for="notifyChannels">Notify Channels (Slack)</Label>
            <Field name="notifyChannels" component={this.renderInput} type="text" placeholder="#eng,#devops"/>
          </FormGroup>
          <FieldArray name="secrets" collapsed={this.state.collapsed} component={this.renderConfigVars}/>
          <br/>
          <Button disabled={submitting || pristine} className="float-xs-right">{submitting ? 'Saving…' : 'Save'}</Button>
        </Form>
        <ButtonConfirmAction btnLabel="Delete Project" onConfirm={() => this.onProjectDelete()} onCancel={this.onProjectCancel} btnClass="btn btn-link text-danger">
          Are you sure you want to <b>DELETE</b> this project?
        </ButtonConfirmAction>
      </div>
    )
  }

  onToggleConfigVars() {
    this.setState({
      collapsed: !this.state.collapsed
    })
  }
}


const ProjectSettingsForm = reduxForm({
  enableReinitialize: true,
  destroyOnUnmount: false,
  form: 'projectSettings',
  onSubmitSuccess: (_component, _dispatch) => {}
})(ProjectSettings)

export default connect(
  state => {
    const formValues = state.form.projectSettings
    return { formValues: formValues}
  }
)(ProjectSettingsForm)
