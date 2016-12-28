import React, { Component } from 'react'

export default class DashboardStats extends Component {
  render() {
    return (
      <div>
        <div className="row statcards">
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-success">
              <div className="p-a">
                <span className="statcard-desc">Projects</span>
                <h2 className="statcard-number">
                  12
                  <small className="delta-indicator delta-positive">5%</small>
                </h2>
                <hr className="statcard-hr m-b-0"/>
              </div>
            </div>
          </div>
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-danger">
              <div className="p-a">
                <span className="statcard-desc">Deploys</span>
                <h2 className="statcard-number">
                  758
                  <small className="delta-indicator delta-negative">1.3%</small>
                </h2>
                <hr className="statcard-hr m-b-0"/>
              </div>
            </div>
          </div>
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-info">
              <div className="p-a">
                <span className="statcard-desc">Code Pushes</span>
                <h2 className="statcard-number">
                  100
                  <small className="delta-indicator delta-positive">6.75%</small>
                </h2>
                <hr className="statcard-hr m-b-0"/>
              </div>
            </div>
          </div>
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-warning">
              <div className="p-a">
                <span className="statcard-desc">Active Users</span>
                <h2 className="statcard-number">
                  25
                  <small className="delta-indicator delta-negative">1.3%</small>
                </h2>
                <hr className="statcard-hr m-b-0"/>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}
