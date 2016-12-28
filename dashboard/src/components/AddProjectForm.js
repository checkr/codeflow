import React, { Component } from 'react'
import { Button, Form, FormGroup, Label, Input } from 'reactstrap'
import { Field, reduxForm } from 'redux-form'


class AddProjectForm extends Component {
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

  render() {
    const { onSubmit } = this.props
    return (
      <Form onSubmit={onSubmit}>
        <FormGroup>
          <Label for="gitSshUrl">Git SSH Url</Label>
          <Field name="gitSshUrl" component={this.renderInput} type="text" placeholder="git@github.com:checkr/codeflow.git"/>
        </FormGroup>
        <FormGroup check>
          <Label check>
            <Field name="bookmarked" component={this.renderCheckbox} type="checkbox"/> Add to my bookmarks
          </Label>
        </FormGroup>
        <br/>
        <Button>Create</Button>
      </Form>
    )
  }
}

export default reduxForm({
  form: 'addProject'
})(AddProjectForm)
