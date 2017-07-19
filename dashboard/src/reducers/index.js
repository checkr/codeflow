import * as ActionTypes from '../actions'
import { routerReducer as routing } from 'react-router-redux'
import { combineReducers } from 'redux'
import { reducer as form } from 'redux-form'
import {
  extend,
  find,
  isEqual,
  union,
  without
} from 'lodash'

const appConfig = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.APP_CONFIG:
      return action.payload
    default:
      return state
  }
}

const me = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.ME_FETCH_REQUEST:
      return {}
    case ActionTypes.ME_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const projectsAll = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECTS_ALL_REQUEST:
      return []
    case ActionTypes.PROJECTS_ALL_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const projects = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECTS_REQUEST:
      if (action.meta.callee === 'pagination') {
        return state
      }
      return []
    case ActionTypes.PROJECTS_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const project = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECT_REQUEST:
      return {}
    case ActionTypes.PROJECT_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const serviceSpecsSettings = (state = [], action = {}) => {
  switch (action.type) {
    case ActionTypes.SERVICE_SPECS_SETTINGS_FETCH_REQUEST:
      return []
    case ActionTypes.SERVICE_SPECS_SETTINGS_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const serviceSpecServices = (state = [], action = {}) => {
  switch (action.type) {
  case ActionTypes.SERVICE_SPEC_SERVICES_FETCH_REQUEST:
    return []
  case ActionTypes.SERVICE_SPEC_SERVICES_FETCH_SUCCESS:
    return action.payload
  default:
    return state
  }
}

const serviceSpecs = (state = [], action = {}) => {
  switch (action.type) {
    case ActionTypes.SERVICE_SPECS_FETCH_REQUEST:
      return []
    case ActionTypes.SERVICE_SPECS_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const projectServices = (state = [], action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECT_SERVICE_FETCH_REQUEST:
      return []
    case ActionTypes.PROJECT_SERVICE_UPDATE_SUCCESS:
      extend(find(state, { _id: action.payload._id }), action.payload)
      return extend([], state)
    case ActionTypes.PROJECT_SERVICE_CREATE_SUCCESS: {
      return union([], state, [action.payload])
    }
    case ActionTypes.PROJECT_SERVICE_DELETE_SUCCESS: {
      return without(state, find(state, { _id: action.payload._id }))
    }
    case ActionTypes.PROJECT_SERVICE_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const projectExtensions = (state = [], action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECT_EXTENSION_FETCH_REQUEST:
      return []
    case ActionTypes.PROJECT_EXTENSION_UPDATE_SUCCESS:
      extend(find(state, { _id: action.payload._id }), action.payload)
      return extend([], state)
    case ActionTypes.PROJECT_EXTENSION_CREATE_SUCCESS: {
      return union([], state, [action.payload])
    }
    case ActionTypes.PROJECT_EXTENSION_DELETE_SUCCESS: {
      return without(state, find(state, { _id: action.payload._id }))
    }
    case ActionTypes.PROJECT_EXTENSION_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const projectSettings = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECT_SETTINGS_FETCH_REQUEST:
      return {}
    case ActionTypes.PROJECT_SETTINGS_UPDATE_SUCCESS:
      return action.payload
    case ActionTypes.PROJECT_SETTINGS_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const bookmarks = (state = [], action = {}) => {
  switch (action.type) {
    case ActionTypes.BOOKMARKS_FETCH_SUCCESS:
      return action.payload
    case ActionTypes.BOOKMARKS_POST_SUCCESS:
      return action.payload
    case ActionTypes.BOOKMARKS_DELETE_SUCCESS:
      return action.payload
    case ActionTypes.WS_MESSAGE_RECEIVED:
      if (action.message.channel === `bookmarks/${action.meta.me._id}`) {
        return action.message.data
      }
      return state
    default:
      return state
  }
}

const features = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.FEATURES_REQUEST:
      if (action.meta.callee === 'pagination') {
        return Object.assign({dirty: false}, state)
      }
      return {}
    case ActionTypes.FEATURES_SUCCESS:
      return action.payload
    case ActionTypes.WS_MESSAGE_RECEIVED: {
      if (action.message.channel === 'features') {
        if (isEqual(action.meta.project._id, action.message.data.projectId)) {
          return Object.assign({dirty: true}, state)
        }
      }
      return state
    }
    default:
      return state
  }
}

const currentRelease = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECT_CURRENT_RELEASE_FETCH_REQUEST:
      return {}
    case ActionTypes.PROJECT_CURRENT_RELEASE_FETCH_SUCCESS:
      return action.payload
    case ActionTypes.WS_MESSAGE_RECEIVED: {
      if (action.message.channel === 'releases/promote') {
        if (isEqual(action.meta.project._id, action.message.data.projectId)) {
          return action.message.data
        }
      }
      if (action.message.channel === 'releases') {
        if (isEqual(action.meta.project._id, action.message.data.projectId) && isEqual(action.message.data._id, state._id)) {
          return action.message.data
        }
      }
      return state
    }
    default:
      return state
  }
}

const bookmarkReleases = (state = {}, action = {}) => {
  switch(action.type) {
    case ActionTypes.BOOKMARK_CURRENT_RELEASE_FETCH_SUCCESS:
      return (() => {
        const o = {}
        if (action.payload.projectId) {
          o[action.payload.projectId] = action.payload
        }
        return  Object.assign({}, state, o)
      })()
    case ActionTypes.REMOVE_BOOKMARK_CURRENT_RELEASES:
      return {}
    case ActionTypes.BOOKMARK_CURRENT_RELEASE_DISMISS:
      return {}
    default:
      return state
  }
}

const releases = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.PROJECT_RELEASES_FETCH_REQUEST:
      if (action.meta.callee === 'pagination') {
        return Object.assign({dirty: false}, state)
      }
      return {}
    case ActionTypes.PROJECT_RELEASES_FETCH_SUCCESS:
      return action.payload
    case ActionTypes.PROJECT_RELEASE_CREATE_SUCCESS:
    case ActionTypes.PROJECT_ROLLBACK_TO_CREATE_SUCCESS: {
      return Object.assign({dirty: true}, state)
    }
    case ActionTypes.WS_MESSAGE_RECEIVED: {
      if (action.message.channel === 'releases') {
        if (isEqual(action.meta.project._id, action.message.data.projectId)) {
          return Object.assign({dirty: true}, state)
        }
      }
      return state
    }
    default:
      return state
  }
}

// Updates error message to notify about the failed fetches.
const errorMessage = (state = null, action) => {
  const { error, payload } = action

  switch (action.type) {
    case ActionTypes.RESET_ERROR_MESSAGE:
      return null
    default:
      break
  }

  if (error && payload && payload.response) {
    return payload.message + ' :: ' + payload.name
  }

  return state
}

const initialAuthState = localStorage.getItem('authState') ? JSON.parse(localStorage.getItem('authState')) : {
  handler: null,
  token: null,
  userName: null,
  isAuthenticated: false,
  refreshToken: false,
  statusText: null
}

const auth = (state = initialAuthState, action) => {
  let auth = {}

  if (action.payload && action.payload.status === 401) {
    auth = Object.assign({}, state, {
      'isAuthenticated': false,
      'refreshToken': false,
      'token': null,
      'userName': null,
      'statusText': 'You have been logged out because your session expired'
    })
    localStorage.setItem('authState', JSON.stringify(auth))
    return auth
  }

  switch (action.type) {
    case ActionTypes.REFRESH_TOKEN:
      state.refreshToken = true
      return state

    case ActionTypes.AUTH_HANDLER_REQUEST:
      auth = Object.assign({}, state, {
        'handler': null,
      })
      localStorage.setItem('authState', JSON.stringify(auth))
      return auth

    case ActionTypes.AUTH_HANDLER_SUCCESS:
      auth = Object.assign({}, state, {
        'handler': action.payload.handler,
      })
      localStorage.setItem('authState', JSON.stringify(auth))
      return auth

    case ActionTypes.REFRESH_TOKEN_SUCCESS:
      auth = Object.assign({}, state, {
        'isAuthenticated': true,
        'refreshToken': false,
        'token': action.response.token,
        'userName': null,
        'statusText': 'You have been successfully logged in.'
      })
      localStorage.setItem('authState', JSON.stringify(auth))
      return auth

    case ActionTypes.AUTH_SUCCESS:
      auth = Object.assign({}, state, {
        'isAuthenticated': true,
        'refreshToken': false,
        'token': action.payload.token,
        'userName': null,
        'statusText': 'You have been successfully logged in.'
      })
      localStorage.setItem('authState', JSON.stringify(auth))
      return auth

    case ActionTypes.LOGOUT_USER:
      auth = Object.assign({}, state, {
        'isAuthenticated': false,
        'refreshToken': false,
        'token': null,
        'userName': null,
        'statusText': 'You have been successfully logged out.'
      })
      localStorage.setItem('authState', JSON.stringify(auth))
      return auth

    case ActionTypes.AUTH_REQUIRED:
      auth = Object.assign({}, state, {
        'isAuthenticated': false,
        'refreshToken': false,
        'token': null,
        'userName': null,
        'statusText': 'You have been logged out because your session expired'
      })
      localStorage.setItem('authState', JSON.stringify(auth))
      return auth
    default:
      return state
  }
}

const stats = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.STATS_FETCH_REQUEST:
      return {}
    case ActionTypes.STATS_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const build = (state = {}, action = {}) => {
  switch (action.type) {
    case ActionTypes.BUILD_FETCH_REQUEST:
      return {}
    case ActionTypes.BUILD_FETCH_SUCCESS:
      return action.payload
    default:
      return state
  }
}

const rootReducer = combineReducers({
  me,
  projects,
  project,
  projectServices,
  projectExtensions,
  projectSettings,
  projectsAll,
  serviceSpecs,
  serviceSpecsSettings,
  serviceSpecServices,
  bookmarks,
  bookmarkReleases,
  features,
  releases,
  currentRelease,
  errorMessage,
  auth,
  routing,
  form,
  appConfig,
  stats,
  build
})

export default rootReducer
