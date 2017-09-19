import React, { Component } from 'react'
import { connect } from 'react-redux'
import ButtonConfirmAction from '../ButtonConfirmAction'
import { projectDeployAll, loadBalancerDeployAll } from '../../actions'

class AdminActionsList extends Component {
    render() {
        return(
            <div>
                <div className="hr-divider m-t-md m-b">
                    <h3 className="hr-divider-content hr-divider-heading">Admin Actions</h3>
                </div>
                <ButtonConfirmAction btnLabel="Deploy All Projects" btnIconClass="fa fa-rocket" size="lg" btnClass="btn btn-warning btn-lg float-xs-center" onConfirm={() => this.onDeployAllProjectsClick()}>
                   Are you sure you want to deploy <b>every</b> project?
                </ButtonConfirmAction>
                <ButtonConfirmAction btnLabel="Update All Services" btnIconClass="fa fa-rocket" size="lg" btnClass="btn btn-danger btn-lg float-xs-center" onConfirm={() => this.onDeployAllLoadBalancersClick()}>
                   Are you sure you want to deploy <b>every</b> load balancer?
                </ButtonConfirmAction>

            </div>
        )
    }

    onDeployAllProjectsClick() {
        this.props.projectDeployAll()
    }
    onDeployAllLoadBalancersClick() {
        this.props.loadBalancerDeployAll()
    }
}

const mapStateToProps = (_state, _ownProps) => ({
})

export default connect(mapStateToProps, {
    projectDeployAll,
    loadBalancerDeployAll,
})(AdminActionsList)