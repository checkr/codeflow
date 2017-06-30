import React, { Component } from 'react'
import { connect } from 'react-redux'
import { forEach, isEmpty } from 'lodash'
import { fetchServiceSpecSettings, updateServiceSpec, deleteServiceSpec, fetchServiceSpecServices, fetchProjects } from '../../actions'
import ServiceSpecForm from "./ServiceSpecForm"
import { SubmissionError } from 'redux-form'

class ServiceSpecList extends Component {

  constructor(props) {
    super(props)
    this.state = {
      edit: null
    }
    this.renderSpecs = this.renderSpecs.bind(this)
    this.onSaveSpec = this.onSaveSpec.bind(this)
    this.onAddSpec = this.onAddSpec.bind(this)
    this.onDeleteSpec = this.onDeleteSpec.bind(this)
    this.onEditServiceSpec = this.onEditServiceSpec.bind(this)
  }

  componentWillMount() {
    this.props.fetchServiceSpecSettings()
    this.props.fetchProjects()
  }

  renderField(name, value, icon) {
    return(
      <div className="row">
        <div className="input-group">
          <div className="col-xs-6">
            <div className="col-xs-right"><i className={icon} aria-hidden="true"></i> {name}</div>
          </div>
          <div className="col-xs-6">
            <input className="form-control" value={value} disabled="true"/>
          </div>
        </div>
      </div>
    )
  }

  renderSpec(specSetting) {
    return (
      <li className="list-group-item" key={specSetting._id}>
        <div className="feed-element">
          <div className="media-body">
            <div className="row">
              <div className="col-xs-12">
                <div className="row">
                  <div className="col-xs-10">
                    <h5>
                      <i className="fa fa-server" aria-hidden="true"></i> {specSetting.name}
                    </h5>
                  </div>
                  <div className="col-xs-2">
                    <button disabled={(this.state.edit) !== null} type="button" className="btn btn-secondary btn-sm float-xs-right btn-edit-resource" onClick={(e) => this.onEditServiceSpec(e, specSetting._id)}>
                      <i className="fa fa-pencil" aria-hidden="true" /> Edit
                    </button>
                  </div>
                </div>
                <div className="row form-group">
                  <div className="col-xs-12">
                    { this.renderField("CPU Request", specSetting.cpuRequest, "fa fa-tachometer")}
                    { this.renderField("CPU Limit", specSetting.cpuLimit, "fa fa-tachometer")}
                    { this.renderField("Memory Request", specSetting.memoryRequest, "fa fa-tachometer")}
                    { this.renderField("Memory Limit", specSetting.memoryLimit, "fa fa-tachometer")}
                    { this.renderField("Termination Grace Period (seconds)", specSetting.terminationGracePeriodSeconds, "fa fa-clock-o")}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
    )
  }

  renderEditSpec(specSetting) {
    return(
      <li className="list-group-item" key={specSetting._id}>
        <div className="feed-element">
          <div className="media-body">
            <ServiceSpecForm serviceSpecServices={this.props.serviceSpecServices} initialValues={specSetting} edit={specSetting.edit} onSubmit={this.onSaveSpec} onCancel={ () => this.onCancel()} onDelete={this.onDeleteSpec} projects={this.props.projects}/>
          </div>
        </div>
      </li>
    )
  }

  renderNewSpec() {
    const specSetting = {}
    return(
      <li className="list-group-item" key="newSpecSetting">
        <div className="feed-element">
          <div className="media-body">
            <ServiceSpecForm serviceSpecServices={[]} initialValues={specSetting} edit={specSetting.edit} onSubmit={this.onSaveSpec} onCancel={ () => this.onCancel()} onDelete={ () => this.onDelete()}/>
          </div>
        </div>
      </li>
    )
  }

  renderSpecs(serviceSpecSettings) {
    let spec_settings_jsx = []
    if (isEmpty(serviceSpecSettings)) {
      if(this.state.edit === 0) {
        spec_settings_jsx.push(this.renderNewSpec())
        return spec_settings_jsx
      }
      return
    }
    forEach(serviceSpecSettings, specSetting => {
      if(this.state.edit && this.state.edit === specSetting._id) {
        spec_settings_jsx.push(this.renderEditSpec(specSetting))
      } else {
        spec_settings_jsx.push(this.renderSpec(specSetting))
      }
    })
    if(this.state.edit === 0) {
      spec_settings_jsx.push(this.renderNewSpec())
    }
    return(spec_settings_jsx)
  }

  onAddSpec(e) {
    e.preventDefault()
    this.setState({ edit: 0 })
  }

  renderAddButton() {
    if (this.state.edit === 0) {
      return null
    } else {
      return (<div><br/><button type="submit" className="btn btn-primary float-xs-right" onClick={(e) => this.onAddSpec(e)}>Add new service spec</button><br/></div>)
    }
  }

  onCancel() {
    this.setState({ edit: null })
  }

  onEditServiceSpec(e, id) {
    e.preventDefault()
    this.props.fetchServiceSpecServices(id)
    this.setState({ edit: id })
  }

  onDeleteSpec() {
    this.props.deleteServiceSpec(this.props.serviceSpecForm.values._id)
      .then(action => {
        if(action.error) {
          const errorMessage = action.payload.response.Error
          throw new SubmissionError({ _error: errorMessage })
        } else {
          this.setState( {edit: null} )
          this.props.fetchServiceSpecSettings()
        }
      })
  }

  onSaveSpec(e) {
    e.preventDefault()
    this.props.updateServiceSpec(this.props.serviceSpecForm.values)
      .then(action => {
          if(action.error) {
            const errorMessage = action.payload.response.Error
            throw new SubmissionError({ _error: errorMessage })
          } else {
            this.setState( {edit: null} )
            this.props.fetchServiceSpecSettings()
          }
        })
  }


  render() {
    const { serviceSpecSettings } = this.props
    return (
      <div>
        <div className="hr-divider m-t-md m-b">
          <h3 className="hr-divider-content hr-divider-heading">Service Specs</h3>
        </div>
        <ul className="list-group">
          {this.renderSpecs(serviceSpecSettings)}
        </ul>
        {this.renderAddButton()}
      </div>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  projects: state.projects,
  serviceSpecSettings: state.serviceSpecsSettings,
  serviceSpecServices: state.serviceSpecServices,
  serviceSpecForm: state.form.serviceSpec
})

export default connect(mapStateToProps, {
  fetchServiceSpecSettings,
  fetchServiceSpecServices,
  fetchProjects,
  updateServiceSpec,
  deleteServiceSpec
})(ServiceSpecList)