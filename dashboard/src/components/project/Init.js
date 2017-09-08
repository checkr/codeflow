import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'
import { Form, FormGroup, Input, Button } from 'reactstrap'
import { Alert } from 'reactstrap'
import { loadProject, appConfig } from '../../actions'
import { isEmpty } from 'lodash'

class ProjectInit extends Component {
  componentWillMount() {
    this.props.appConfig()
  }

  render() {
    const { project, config } = this.props

    if (isEmpty(project) || project.pinged) {
      return null
    }

    const gitDeployKeyURL = `https://github.com/${project.name}/settings/keys`
    const gitWebhookSettingsURL = `https://github.com/${project.name}/settings/hooks/new`

    if (/^https:/.test(project.gitUrl)) {
      return (
        <div>
          <div className="hr-divider m-t-md m-b">
            <h3 className="hr-divider-content hr-divider-heading">Initial Setup</h3>
          </div>
          No additional setup required.  Polling for commits...  Refresh the page or continue to setup Resources.
        </div>
      )
    }

    return (
      <div>
        <div className="hr-divider m-t-md m-b">
          <h3 className="hr-divider-content hr-divider-heading">Initial Setup</h3>
        </div>
        <Form>
          <h4>1. Setup the Git Deploy Key</h4>
          <FormGroup>
            <Input readOnly="readOnly" type="textarea" name="text" id="rsaPublicKey" rows="6" value={project.rsaPublicKey.replace(/^\s+|\s+$/g, '')}/>
            <Alert color="info">
              <Button color="link" target="_blank" href={gitDeployKeyURL}>Copy this key and <b>click here</b> to add it to the Github deploy key list</Button>
            </Alert>
          </FormGroup>
          <h4>2. Setup the Github webhook</h4>
          <Alert color="info">
            <Button color="link" target="_blank" href={gitWebhookSettingsURL}><b>Click here</b> to open the github webhook settings</Button>
          </Alert>
          <FormGroup>
            <Input readOnly="readOnly" type="text" name="content_type" id="githubcontenttype" value="Content Type: application/json" />
            <Alert color="info">
              <icon className="fa fa-info-circle"/> Make sure to choose Content Type: <b>application/json</b> webhook settings.
            </Alert>
          </FormGroup>
          <FormGroup>
            <Input readOnly="readOnly" type="text" name="github_webhook" id="githubWebhook" value={"Payload URL: " + config.REACT_APP_WEBHOOKS_ROOT + "/github"} />
            <Alert color="info">
              <icon className="fa fa-info-circle"/> Copy and paste the Payload URL into the webhook settings page.
            </Alert>
          </FormGroup>
          <FormGroup>
            <Input readOnly="readOnly" type="text" name="github_webhook_secret" id="githubWebhookSecret" value={"Secret: " + project.secret} />
            <Alert color="info">
              <icon className="fa fa-info-circle"/> Copy and paste the Secret into the webhook settings page.
            </Alert>
          </FormGroup>
        </Form>
        <hr/>
        <Alert color="warning">
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
