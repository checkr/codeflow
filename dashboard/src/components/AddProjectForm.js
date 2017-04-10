import React, { Component } from 'react'
import { Button, Form, FormGroup, Label, Input, FormText } from 'reactstrap'
import { Field, reduxForm } from 'redux-form'
import { connect } from 'react-redux'


class AddProject extends Component {
  renderInput(field) {   
    return (
      <Input {...field.input} type={field.type} placeholder={field.placeholder} />
    )  
  }

  renderCheckbox(field) {   
    return (
      <Input {...field.input} type={field.type} placeholder={field.placeholder} />
    )  
  }

  renderGitUrl() {
    let { formValues } = this.props
    if (formValues && formValues.values.gitProtocol === "HTTPS") {
      return (
        <FormGroup>
          <Label for="gitUrl">Git HTTPS Url</Label>
          <Field name="gitUrl" component={this.renderInput} type="text" placeholder="https://github.com/checkr/codeflow.git"/>
        </FormGroup>
        )
    } else {
      return (
        <FormGroup>
          <Label for="gitUrl">Git SSH Url</Label>
          <Field name="gitUrl" component={this.renderInput} type="text" placeholder="git@github.com:checkr/codeflow.git"/>
        </FormGroup>
        )
    }
  }

  render() {
    const { onSubmit } = this.props
    return (
      <Form onSubmit={onSubmit}>
        <FormGroup>
          <Label>Protocol</Label>
          <FormGroup tag="fieldset">
            <FormGroup check>
              <Label check>
                <Field className="form-check-input" name={'gitProtocol'} component={this.renderInput} type="radio" value="SSH"/> SSH
              </Label>
              <FormText color="muted">Use SSH for private repositories</FormText>
            </FormGroup>
            <FormGroup check>
              <Label check>
                <Field className="form-check-input" name={'gitProtocol'} component={this.renderInput} type="radio" value="HTTPS"/> HTTPS
              </Label>
              <FormText color="muted">Use HTTPS for public repositories</FormText>
            </FormGroup>
          </FormGroup>
        </FormGroup>
        {this.renderGitUrl()}
        <FormGroup>
          <Label>
            <Field name="bookmarked" component={this.renderCheckbox} type="checkbox"/> Add to my bookmarks
          </Label>
        </FormGroup>
        <br/>
        <Button>Create</Button>
      </Form>
    )
  }
}

const AddProjectForm = reduxForm({
  form: 'addProject'
})(AddProject)

export default connect(
  state => {
    const formValues = state.form.addProject
    return { formValues: formValues, initialValues: {gitProtocol: "SSH", bookmarked: true} }
  }
)(AddProjectForm)
