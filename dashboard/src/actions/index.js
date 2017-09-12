import { CALL_API } from 'redux-api-middleware'
import loadConfig from "../config"

const CONFIG = loadConfig()
const API_ROOT = CONFIG.REACT_APP_API_ROOT;
const WS_ROOT = CONFIG.REACT_APP_WS_ROOT;

export const APP_CONFIG = 'APP_CONFIG'
export const appConfig = () => {
  return {
    type: APP_CONFIG,
    payload: CONFIG,
  }
}

export const postBody = (payload) => {
  if (typeof payload === 'object') {
    return JSON.stringify(payload)
  }
  return payload
}

export const authToken = () => {
  if (localStorage.authState) {
    return JSON.parse(localStorage.authState).token
  }
  return ''
}

export const ME_FETCH_REQUEST = 'ME_FETCH_REQUEST'
export const ME_FETCH_SUCCESS = 'ME_FETCH_SUCCESS'
export const ME_FETCH_FAILURE = 'ME_FETCH_FAILURE'
export const USER_FETCH_REQUEST = 'USER_FETCH_REQUEST'
export const USER_FETCH_SUCCESS = 'USER_FETCH_SUCCESS'
export const USER_FETCH_FAILURE = 'USER_FETCH_FAILURE'
export const USERS_FETCH_REQUEST = 'USERS_FETCH_REQUEST'
export const USERS_FETCH_SUCCESS = 'USERS_FETCH_SUCCESS'
export const USERS_FETCH_FAILURE = 'USERS_FETCH_FAILURE'

export const fetchUsers = (username = '') => {
  let endpoint = (username === '') ? '/users' : `users/${username}`
  let types = [ USERS_FETCH_REQUEST, USERS_FETCH_SUCCESS, USERS_FETCH_FAILURE ]

  if (username === 'me') {
    types = [ ME_FETCH_REQUEST, ME_FETCH_SUCCESS, ME_FETCH_FAILURE ]
  } else if (username !== '') {
    types = [ USER_FETCH_REQUEST, USER_FETCH_SUCCESS, USER_FETCH_FAILURE ]
  }

  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/${endpoint}`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: types
    }
  }
}

export const AUTH_HANDLER_REQUEST = 'AUTH_HANDLER_REQUEST'
export const AUTH_HANDLER_SUCCESS = 'AUTH_HANDLER_SUCCESS'
export const AUTH_HANDLER_FAILURE = 'AUTH_HANDLER_FAILURE'

export const authHandler = () => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/auth/handler`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      types: [ AUTH_HANDLER_REQUEST, AUTH_HANDLER_SUCCESS, AUTH_HANDLER_FAILURE ]
    }
  }
}
export const AUTH_REQUEST = 'AUTH_REQUEST'
export const AUTH_SUCCESS = 'AUTH_SUCCESS'
export const AUTH_FAILURE = 'AUTH_FAILURE'
export const AUTH_REQUIRED = 'AUTH_REQUIRED'

export const authCallback = (endpoint, payload) => {
  endpoint = (endpoint.indexOf('/') === -1) ? '/' + endpoint : endpoint
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}${endpoint}`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ AUTH_REQUEST, AUTH_SUCCESS, AUTH_FAILURE ]
    }
  }
}

export const NEW_PROJECT_REQUEST = 'NEW_PROJECT_REQUEST'
export const NEW_PROJECT_SUCCESS = 'NEW_PROJECT_SUCCESS'
export const NEW_PROJECT_FAILURE = 'NEW_PROJECT_FAILURE'

export const createNewProject = (payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ NEW_PROJECT_REQUEST, NEW_PROJECT_SUCCESS, NEW_PROJECT_FAILURE ]
    }
  }
}

export const DELETE_PROJECT_REQUEST = 'DELETE_PROJECT_REQUEST'
export const DELETE_PROJECT_SUCCESS = 'DELETE_PROJECT_SUCCESS'
export const DELETE_PROJECT_FAILURE = 'DELETE_PROJECT_FAILURE'

export const deleteProject = (slug) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${slug}`,
      method: 'DELETE',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ DELETE_PROJECT_REQUEST, DELETE_PROJECT_SUCCESS, DELETE_PROJECT_FAILURE ]
    }
  }
}

export const PROJECTS_ALL_REQUEST = 'PROJECTS_ALL_REQUEST'
export const PROJECTS_ALL_SUCCESS = 'PROJECTS_ALL_SUCCESS'
export const PROJECTS_ALL_FAILURE = 'PROJECTS_ALL_FAILURE'

export const fetchProjectsAll = () => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects?max_items_per_page=9999`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [
        {
          type: 'PROJECTS_ALL_REQUEST',
        },
        PROJECTS_ALL_SUCCESS, PROJECTS_ALL_FAILURE
      ]
    }
  }
}

export const PROJECTS_REQUEST = 'PROJECTS_REQUEST'
export const PROJECTS_SUCCESS = 'PROJECTS_SUCCESS'
export const PROJECTS_FAILURE = 'PROJECTS_FAILURE'

export const fetchProjects = (query = '', callee = '') => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects${query}`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [
        {
          type: 'PROJECTS_REQUEST',
          meta: { callee: callee }
        },
        PROJECTS_SUCCESS, PROJECTS_FAILURE
      ]
    }
  }
}

export const PROJECT_REQUEST = 'PROJECT_REQUEST'
export const PROJECT_SUCCESS = 'PROJECT_SUCCESS'
export const PROJECT_FAILURE = 'PROJECT_FAILURE'

export const fetchProject = (slug) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${slug}`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_REQUEST, PROJECT_SUCCESS, PROJECT_FAILURE ]
    }
  }
}

export const BOOKMARKS_FETCH_REQUEST = 'BOOKMARKS_FETCH_REQUEST'
export const BOOKMARKS_FETCH_SUCCESS = 'BOOKMARKS_FETCH_SUCCESS'
export const BOOKMARKS_FETCH_FAILURE = 'BOOKMARKS_FETCH_FAILURE'

export const fetchBookmarks = (_payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/bookmarks`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ BOOKMARKS_FETCH_REQUEST, BOOKMARKS_FETCH_SUCCESS, BOOKMARKS_FETCH_FAILURE ]
    }
  }
}

export const BOOKMARKS_POST_REQUEST = 'BOOKMARKS_POST_REQUEST'
export const BOOKMARKS_POST_SUCCESS = 'BOOKMARKS_POST_SUCCESS'
export const BOOKMARKS_POST_FAILURE = 'BOOKMARKS_POST_FAILURE'

export const createBookmarks = (payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/bookmarks`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ BOOKMARKS_POST_REQUEST, BOOKMARKS_POST_SUCCESS, BOOKMARKS_POST_FAILURE ]
    }
  }
}

export const BOOKMARKS_DELETE_REQUEST = 'BOOKMARKS_DELETE_REQUEST'
export const BOOKMARKS_DELETE_SUCCESS = 'BOOKMARKS_DELETE_SUCCESS'
export const BOOKMARKS_DELETE_FAILURE = 'BOOKMARKS_DELETE_FAILURE'

export const deleteBookmarks = (payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/bookmarks`,
      method: 'DELETE',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ BOOKMARKS_DELETE_REQUEST, BOOKMARKS_DELETE_SUCCESS, BOOKMARKS_DELETE_FAILURE ]
    }
  }
}

export const FEATURES_REQUEST = 'FEATURES_REQUEST'
export const FEATURES_SUCCESS = 'FEATURES_SUCCESS'
export const FEATURES_FAILURE = 'FEATURES_FAILURE'

export const fetchProjectFeatures = (slug = '', query = '', callee = '' ) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${slug}/features${query}`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [
        {
          type: 'FEATURES_REQUEST',
          meta: { callee: callee }
        },
        FEATURES_SUCCESS, FEATURES_FAILURE ]
    }
  }
}

// Resets the currently visible error message.
export const RESET_ERROR_MESSAGE = 'RESET_ERROR_MESSAGE'

export const resetErrorMessage = () => ({
  type: RESET_ERROR_MESSAGE
})

export const LOGOUT_USER = 'LOGOUT_USER'

// Destroy auth object.
export const logoutUser = () => ({
  type: LOGOUT_USER
})

export const REFRESH_TOKEN = 'REFRESH_TOKEN'
export const REFRESH_TOKEN_REQUEST = 'REFRESH_TOKEN_REQUEST'
export const REFRESH_TOKEN_SUCCESS = 'REFRESH_TOKEN_SUCCESS'

export const fetchRefreshToken = () => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/auth/refresh_token`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ REFRESH_TOKEN_REQUEST, REFRESH_TOKEN_SUCCESS, AUTH_REQUIRED ]
    }
  }
}

export const PROJECT_SERVICE_FETCH_REQUEST = 'PROJECT_SERVICE_FETCH_REQUEST'
export const PROJECT_SERVICE_FETCH_SUCCESS = 'PROJECT_SERVICE_FETCH_SUCCESS'
export const PROJECT_SERVICE_FETCH_FAILURE = 'PROJECT_SERVICE_FETCH_FAILURE'

export const fetchProjectServices = (project_slug) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/services`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_SERVICE_FETCH_REQUEST, PROJECT_SERVICE_FETCH_SUCCESS, PROJECT_SERVICE_FETCH_FAILURE ]
    }
  }
}

export const PROJECT_SERVICE_CREATE_REQUEST = 'PROJECT_SERVICE_CREATE_REQUEST'
export const PROJECT_SERVICE_CREATE_SUCCESS = 'PROJECT_SERVICE_CREATE_SUCCESS'
export const PROJECT_SERVICE_CREATE_FAILURE = 'PROJECT_SERVICE_CREATE_FAILURE'

export const createProjectService = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/services`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_SERVICE_CREATE_REQUEST, PROJECT_SERVICE_CREATE_SUCCESS, PROJECT_SERVICE_CREATE_FAILURE ]
    }
  }
}

export const PROJECT_SERVICE_UPDATE_REQUEST = 'PROJECT_SERVICE_UPDATE_REQUEST'
export const PROJECT_SERVICE_UPDATE_SUCCESS = 'PROJECT_SERVICE_UPDATE_SUCCESS'
export const PROJECT_SERVICE_UPDATE_FAILURE = 'PROJECT_SERVICE_UPDATE_FAILURE'

export const updateProjectService = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/services`,
      method: 'PUT',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_SERVICE_UPDATE_REQUEST, PROJECT_SERVICE_UPDATE_SUCCESS, PROJECT_SERVICE_UPDATE_FAILURE ]
    }
  }
}

export const PROJECT_SERVICE_DELETE_REQUEST = 'PROJECT_SERVICE_DELETE_REQUEST'
export const PROJECT_SERVICE_DELETE_SUCCESS = 'PROJECT_SERVICE_DELETE_SUCCESS'
export const PROJECT_SERVICE_DELETE_FAILURE = 'PROJECT_SERVICE_DELETE_FAILURE'

export const deleteProjectService = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/services`,
      method: 'DELETE',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_SERVICE_DELETE_REQUEST, PROJECT_SERVICE_DELETE_SUCCESS, PROJECT_SERVICE_DELETE_FAILURE ]
    }
  }
}

export const PROJECT_RELEASE_CREATE_REQUEST = 'PROJECT_RELEASE_CREATE_REQUEST'
export const PROJECT_RELEASE_CREATE_SUCCESS = 'PROJECT_RELEASE_CREATE_SUCCESS'
export const PROJECT_RELEASE_CREATE_FAILURE = 'PROJECT_RELEASE_CREATE_FAILURE'

export const createProjectRelease = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/releases`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_RELEASE_CREATE_REQUEST, PROJECT_RELEASE_CREATE_SUCCESS, PROJECT_RELEASE_CREATE_FAILURE ]
    }
  }
}

export const PROJECT_ROLLBACK_TO_CREATE_REQUEST = 'PROJECT_ROLLBACK_TO_CREATE_REQUEST'
export const PROJECT_ROLLBACK_TO_CREATE_SUCCESS = 'PROJECT_ROLLBACK_TO_CREATE_SUCCESS'
export const PROJECT_ROLLBACK_TO_CREATE_FAILURE = 'PROJECT_ROLLBACK_TO_CREATE_FAILURE'

export const createProjectRollbackTo = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/releases/rollback`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_ROLLBACK_TO_CREATE_REQUEST, PROJECT_ROLLBACK_TO_CREATE_SUCCESS, PROJECT_ROLLBACK_TO_CREATE_FAILURE ]
    }
  }
}

export const PROJECT_RELEASES_FETCH_REQUEST = 'PROJECT_RELEASES_FETCH_REQUEST'
export const PROJECT_RELEASES_FETCH_SUCCESS = 'PROJECT_RELEASES_FETCH_SUCCESS'
export const PROJECT_RELEASES_FETCH_FAILURE = 'PROJECT_RELEASES_FETCH_FAILURE'

export const fetchProjectReleases = (slug = '', query = '', callee = '') => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${slug}/releases${query}`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [
        {
          type: PROJECT_RELEASES_FETCH_REQUEST,
          meta: { callee: callee }
        },
        PROJECT_RELEASES_FETCH_SUCCESS, PROJECT_RELEASES_FETCH_FAILURE
      ]
    }
  }
}

export const PROJECT_CURRENT_RELEASE_FETCH_REQUEST = 'PROJECT_CURRENT_RELEASE_FETCH_REQUEST'
export const PROJECT_CURRENT_RELEASE_FETCH_SUCCESS = 'PROJECT_CURRENT_RELEASE_FETCH_SUCCESS'
export const PROJECT_CURRENT_RELEASE_FETCH_FAILURE = 'PROJECT_CURRENT_RELEASE_FETCH_FAILURE'

const fetchProjectCurrentReleaseFactory = ({slug, types}) => ({
  [CALL_API]: {
    endpoint: `${API_ROOT}/projects/${slug}/releases/current`,
    method: 'GET',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authToken()}`
    },
    types
  }
})

export const fetchProjectCurrentRelease = (slug = '') => {
  const types = [
    PROJECT_CURRENT_RELEASE_FETCH_REQUEST,
    PROJECT_CURRENT_RELEASE_FETCH_SUCCESS,
    PROJECT_CURRENT_RELEASE_FETCH_FAILURE
  ]

  return fetchProjectCurrentReleaseFactory({slug, types})
}

export const BOOKMARK_CURRENT_RELEASE_FETCH_REQUEST = 'BOOKMARK_CURRENT_RELEASE_FETCH_REQUEST'
export const BOOKMARK_CURRENT_RELEASE_FETCH_SUCCESS = 'BOOKMARK_CURRENT_RELEASE_FETCH_SUCCESS'
export const BOOKMARK_CURRENT_RELEASE_FETCH_FAILURE = 'BOOKMARK_CURRENT_RELEASE_FETCH_FAILURE'
export const BOOKMARK_CURRENT_RELEASE_DISMISS = 'BOOKMARK_CURRENT_RELEASE_DISMISS'

export const fetchBookmarkCurrentRelease = ({slug = ''}) => {
  const types = [
    BOOKMARK_CURRENT_RELEASE_FETCH_REQUEST,
    BOOKMARK_CURRENT_RELEASE_FETCH_SUCCESS,
    BOOKMARK_CURRENT_RELEASE_FETCH_FAILURE
  ]

  return fetchProjectCurrentReleaseFactory({slug, types})
}

export const REMOVE_BOOKMARK_CURRENT_RELEASES = 'REMOVE_BOOKMARK_CURRENT_RELEASES'
export const removeBookmarkCurrentReleases = () => ({
  type: REMOVE_BOOKMARK_CURRENT_RELEASES
})

export const PROJECT_EXTENSION_FETCH_REQUEST = 'PROJECT_EXTENSION_FETCH_REQUEST'
export const PROJECT_EXTENSION_FETCH_SUCCESS = 'PROJECT_EXTENSION_FETCH_SUCCESS'
export const PROJECT_EXTENSION_FETCH_FAILURE = 'PROJECT_EXTENSION_FETCH_FAILURE'

export const fetchProjectExtensions = (project_slug) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/extensions`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_EXTENSION_FETCH_REQUEST, PROJECT_EXTENSION_FETCH_SUCCESS, PROJECT_EXTENSION_FETCH_FAILURE ]
    }
  }
}

export const PROJECT_EXTENSION_CREATE_REQUEST = 'PROJECT_EXTENSION_CREATE_REQUEST'
export const PROJECT_EXTENSION_CREATE_SUCCESS = 'PROJECT_EXTENSION_CREATE_SUCCESS'
export const PROJECT_EXTENSION_CREATE_FAILURE = 'PROJECT_EXTENSION_CREATE_FAILURE'

export const createProjectExtension = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/extensions`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_EXTENSION_CREATE_REQUEST, PROJECT_EXTENSION_CREATE_SUCCESS, PROJECT_EXTENSION_CREATE_FAILURE ]
    }
  }
}

export const PROJECT_EXTENSION_UPDATE_REQUEST = 'PROJECT_EXTENSION_UPDATE_REQUEST'
export const PROJECT_EXTENSION_UPDATE_SUCCESS = 'PROJECT_EXTENSION_UPDATE_SUCCESS'
export const PROJECT_EXTENSION_UPDATE_FAILURE = 'PROJECT_EXTENSION_UPDATE_FAILURE'

export const updateProjectExtension = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/extensions`,
      method: 'PUT',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_EXTENSION_UPDATE_REQUEST, PROJECT_EXTENSION_UPDATE_SUCCESS, PROJECT_EXTENSION_UPDATE_FAILURE ]
    }
  }
}

export const PROJECT_EXTENSION_DELETE_REQUEST = 'PROJECT_EXTENSION_DELETE_REQUEST'
export const PROJECT_EXTENSION_DELETE_SUCCESS = 'PROJECT_EXTENSION_DELETE_SUCCESS'
export const PROJECT_EXTENSION_DELETE_FAILURE = 'PROJECT_EXTENSION_DELETE_FAILURE'

export const deleteProjectExtension = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/extensions`,
      method: 'DELETE',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_EXTENSION_DELETE_REQUEST, PROJECT_EXTENSION_DELETE_SUCCESS, PROJECT_EXTENSION_DELETE_FAILURE ]
    }
  }
}


export const WS_CONNECT = 'WS_CONNECT'
export const WS_CONNECTING = 'WS_CONNECTING'
export const WS_CONNECTED = 'WS_CONNECTED'
export const WS_DISCONNECT = 'WS_DISCONNECT'
export const WS_DISCONNECTED = 'WS_DISCONNECTED'
export const WS_MESSAGE_RECEIVED = 'WS_MESSAGE_RECEIVED'
export const WS_SEND_MESSAGE = 'WS_SEND_MESSAGE'

export const wsConnect = () => ({
  type: WS_CONNECT,
  url: `${WS_ROOT}`
})

export const wsConnecting = () => ({
  type: WS_CONNECTING
})

export const wsConnected = () => ({
  type: WS_CONNECTED
})

export const wsDisconnect = () => ({
  type: WS_DISCONNECT
})

export const wsDisconnected = () => ({
  type: WS_DISCONNECTED
})

export const wsMessageReceived = (msg, meta) => ({
  type: WS_MESSAGE_RECEIVED,
  message: msg,
  meta: meta
})

export const wsSendMsg = () => ({
  type: WS_SEND_MESSAGE
})

export const PROJECT_SETTINGS_FETCH_REQUEST = 'PROJECT_SETTINGS_FETCH_REQUEST'
export const PROJECT_SETTINGS_FETCH_SUCCESS = 'PROJECT_SETTINGS_FETCH_SUCCESS'
export const PROJECT_SETTINGS_FETCH_FAILURE = 'PROJECT_SETTINGS_FETCH_FAILURE'

export const fetchProjectSettings = (project_slug, _payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/settings`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_SETTINGS_FETCH_REQUEST, PROJECT_SETTINGS_FETCH_SUCCESS, PROJECT_SETTINGS_FETCH_FAILURE ]
    }
  }
}

export const PROJECT_SETTINGS_UPDATE_REQUEST = 'PROJECT_SETTINGS_UPDATE_REQUEST'
export const PROJECT_SETTINGS_UPDATE_SUCCESS = 'PROJECT_SETTINGS_UPDATE_SUCCESS'
export const PROJECT_SETTINGS_UPDATE_FAILURE = 'PROJECT_SETTINGS_UPDATE_FAILURE'

export const updateProjectSettings = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/settings`,
      method: 'PUT',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_SETTINGS_UPDATE_REQUEST, PROJECT_SETTINGS_UPDATE_SUCCESS, PROJECT_SETTINGS_UPDATE_FAILURE ]
    }
  }
}

export const STATS_FETCH_REQUEST = 'STATS_FETCH_REQUEST'
export const STATS_FETCH_SUCCESS = 'STATS_FETCH_SUCCESS'
export const STATS_FETCH_FAILURE = 'STATS_FETCH_FAILURE'

export const fetchStats = (slug = "") => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/stats${slug}`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ STATS_FETCH_REQUEST, STATS_FETCH_SUCCESS, STATS_FETCH_FAILURE ]
    }
  }
}

export const BUILD_FETCH_REQUEST = 'BUILD_FETCH_REQUEST'
export const BUILD_FETCH_SUCCESS = 'BUILD_FETCH_SUCCESS'
export const BUILD_FETCH_FAILURE = 'BUILD_FETCH_FAILURE'

export const fetchBuild = (project_slug, id) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/releases/${id}/build`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ BUILD_FETCH_REQUEST, BUILD_FETCH_SUCCESS, BUILD_FETCH_FAILURE ]
    }
  }
}

export const BUILD_UPDATE_REQUEST = 'BUILD_UPDATE_REQUEST'
export const BUILD_UPDATE_SUCCESS = 'BUILD_UPDATE_SUCCESS'
export const BUILD_UPDATE_FAILURE = 'BUILD_UPDATE_FAILURE'

export const updateBuild = (project_slug, id) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/releases/${id}/build`,
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ BUILD_UPDATE_REQUEST, BUILD_UPDATE_SUCCESS, BUILD_UPDATE_FAILURE ]
    }
  }
}

export const SERVICE_SPECS_FETCH_REQUEST = 'SERVICE_SPECS_FETCH_REQUEST'
export const SERVICE_SPECS_FETCH_SUCCESS = 'SERVICE_SPECS_FETCH_SUCCESS'
export const SERVICE_SPECS_FETCH_FAILURE = 'SERVICE_SPECS_FETCH_FAILURE'

export const fetchServiceSpecs = (project_slug) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/serviceSpecs`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ SERVICE_SPECS_FETCH_REQUEST, SERVICE_SPECS_FETCH_SUCCESS, SERVICE_SPECS_FETCH_FAILURE ]
    }
  }
}

export const SERVICE_SPECS_SETTINGS_FETCH_REQUEST = 'SERVICE_SPECS_SETTINGS_FETCH_REQUEST'
export const SERVICE_SPECS_SETTINGS_FETCH_SUCCESS = 'SERVICE_SPECS_SETTINGS_FETCH_SUCCESS'
export const SERVICE_SPECS_SETTINGS_FETCH_FAILURE = 'SERVICE_SPECS_SETTINGS_FETCH_FAILURE'

export const fetchServiceSpecSettings = () => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/admin/serviceSpecs`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ SERVICE_SPECS_SETTINGS_FETCH_REQUEST, SERVICE_SPECS_SETTINGS_FETCH_SUCCESS, SERVICE_SPECS_SETTINGS_FETCH_FAILURE ]
    }
  }
}

export const SERVICE_SPECS_SETTINGS_UPDATE_REQUEST = 'SERVICE_SPECS_SETTINGS_UPDATE_REQUEST'
export const SERVICE_SPECS_SETTINGS_UPDATE_SUCCESS = 'SERVICE_SPECS_SETTINGS_UPDATE_SUCCESS'
export const SERVICE_SPECS_SETTINGS_UPDATE_FAILURE = 'SERVICE_SPECS_SETTINGS_UPDATE_FAILURE'

export const updateServiceSpec = (serviceSpec) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/admin/serviceSpecs`,
      method: 'PUT',
      body: postBody(serviceSpec),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ SERVICE_SPECS_SETTINGS_UPDATE_REQUEST, SERVICE_SPECS_SETTINGS_UPDATE_SUCCESS, SERVICE_SPECS_SETTINGS_UPDATE_FAILURE ]
    }
  }
}

export const SERVICE_SPECS_SETTINGS_DELETE_REQUEST = 'SERVICE_SPECS_SETTINGS_DELETE_REQUEST'
export const SERVICE_SPECS_SETTINGS_DELETE_SUCCESS = 'SERVICE_SPECS_SETTINGS_DELETE_SUCCESS'
export const SERVICE_SPECS_SETTINGS_DELETE_FAILURE = 'SERVICE_SPECS_SETTINGS_DELETE_FAILURE'

export const deleteServiceSpec = (slug) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/admin/serviceSpecs/${slug}`,
      method: 'DELETE',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ SERVICE_SPECS_SETTINGS_DELETE_REQUEST, SERVICE_SPECS_SETTINGS_DELETE_SUCCESS, SERVICE_SPECS_SETTINGS_DELETE_FAILURE ]
    }
  }
}

export const SERVICE_SPEC_SERVICES_FETCH_REQUEST = 'SERVICE_SPEC_SERVICES_FETCH_REQUEST'
export const SERVICE_SPEC_SERVICES_FETCH_SUCCESS = 'SERVICE_SPEC_SERVICES_FETCH_SUCCESS'
export const SERVICE_SPEC_SERVICES_FETCH_FAILURE = 'SERVICE_SPEC_SERVICES_FETCH_FAILURE'

export const fetchServiceSpecServices = (slug) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/admin/serviceSpecs/${slug}/services`,
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ SERVICE_SPEC_SERVICES_FETCH_REQUEST, SERVICE_SPEC_SERVICES_FETCH_SUCCESS, SERVICE_SPEC_SERVICES_FETCH_FAILURE ]
    }
  }
}

export const PROJECT_CANCEL_RELEASE_REQUEST = 'PROJECT_CANCEL_RELEASE_REQUEST'
export const PROJECT_CANCEL_RELEASE_SUCCESS = 'PROJECT_CANCEL_RELEASE_SUCCESS'
export const PROJECT_CANCEL_RELEASE_FAILURE = 'PROJECT_CANCEL_RELEASE_FAILURE'

export const cancelRelease = (project_slug, payload) => {
  return {
    [CALL_API]: {
      endpoint: `${API_ROOT}/projects/${project_slug}/releases/cancel`,
      method: 'POST',
      body: postBody(payload),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken()}`
      },
      types: [ PROJECT_CANCEL_RELEASE_REQUEST, PROJECT_CANCEL_RELEASE_SUCCESS, PROJECT_CANCEL_RELEASE_FAILURE ]
    }
  }
}
