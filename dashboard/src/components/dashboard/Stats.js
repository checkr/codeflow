import React, { Component } from 'react'
import { connect } from 'react-redux'
import { fetchStats } from '../../actions'

class DashboardStats extends Component {
  loadData = props => {
    props.fetchStats()
  }

  componentWillMount() {
    this.loadData(this.props)
  }

  render() {
    let { projects, features, releases, users } = this.props.stats
    return (
      <div>
        <div className="row statcards">
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-success">
              <div className="p-a">
                <span className="statcard-desc">Projects</span>
                <h2 className="statcard-number">
                  {projects}
                </h2>
              </div>
            </div>
          </div>
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-info">
              <div className="p-a">
                <span className="statcard-desc">Features</span>
                <h2 className="statcard-number">
                  {features}
                </h2>
              </div>
            </div>
          </div>
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-danger">
              <div className="p-a">
                <span className="statcard-desc">Releases</span>
                <h2 className="statcard-number">
                  {releases}
                </h2>
              </div>
            </div>
          </div>
          <div className="col-sm-3 m-b">
            <div className="statcard statcard-warning">
              <div className="p-a">
                <span className="statcard-desc">Users</span>
                <h2 className="statcard-number">
                  {users}
                </h2>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

const mapStateToProps = (state, _ownProps) => ({
  stats: state.stats,
})

export default connect(mapStateToProps, {
  fetchStats,
})(DashboardStats)
