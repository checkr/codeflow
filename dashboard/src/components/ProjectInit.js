import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'
import { Form, FormGroup, Label, Input } from 'reactstrap'
import { Alert } from 'reactstrap'
import { loadProject, appConfig } from '../actions'
import _ from 'underscore'

class ProjectInit extends Component {
  componentWillMount() {
    this.props.appConfig()
  }

  render() {
    const { project, config } = this.props

    if (_.isEmpty(project) || project.pinged) {
      return null
    }

    return (
      <div>
        <div className="hr-divider m-t-md m-b">
          <h3 className="hr-divider-content hr-divider-heading">Initial Setup</h3>
        </div>
        <Form>
          <FormGroup>
            <Label for="github_webhook">Github Webhook</Label>
            <Input readOnly="readOnly" type="text" name="github_webhook" id="githubWebhook" value={config.REACT_APP_WEBHOOKS_ROOT + "/github"} />
          </FormGroup>
          <FormGroup>
            <Label for="github_webhook">Webhook Secret</Label>
            <Input readOnly="readOnly" type="text" name="github_webhook_secret" id="githubWebhookSecret" value={project.secret} />
          </FormGroup>
          <FormGroup>
            <Label for="rsa_pub_key">RSA Public Key</Label>
            <Input readOnly="readOnly" type="textarea" name="text" id="rsaPublicKey" rows="6" value={project.rsaPublicKey}/>
          </FormGroup>
        </Form>

        <hr/>
        <Alert color="info">
          This page will refresh once first webhook is received. You can still access this information later by visiting the <Link to={'/projects/' + project.slug + '/settings'}><strong>Settings</strong></Link> tab.
        </Alert>
      </div>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  config: state.appConfig
})

export default connect(mapStateToProps, {
  appConfig,
  loadProject
})(ProjectInit)
