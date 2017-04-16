import { createStore, applyMiddleware } from 'redux'
import thunk from 'redux-thunk'
import rootReducer from '../reducers'
import { apiMiddleware } from 'redux-api-middleware'
import socketMiddleware from '../middleware/socket'

const configureStore = preloadedState => createStore(
  rootReducer,
  preloadedState,
  applyMiddleware(thunk, socketMiddleware, apiMiddleware)
)

export default configureStore
