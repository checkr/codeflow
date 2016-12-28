import { wsConnect, wsConnected, wsDisconnect, wsDisconnected, wsMessageReceived, wsConnecting } from '../actions'

const socketMiddleware = (function() {

  var socket = null
  var reconnect = null
  var reconnectCount = 0


  const onOpen = (ws, store, _token) => _evt => {
    //Tell the store we're connected
    store.dispatch(wsConnected())
    clearTimeout(reconnect)
    if(reconnectCount > 5) {
      reconnectCount = 0
      location.reload()
    }
  }

  const onClose = (ws, store) => _evt => {
    //Tell the store we've disconnected
    store.dispatch(wsDisconnect())
  }

  const onMessage = (ws, store) => evt => {
    //Parse the JSON message received on the websocket
    let msg = JSON.parse(evt.data)
    const state = store.getState()
    let meta = {
      me: state.me,
      project: state.project,
    }
    store.dispatch(wsMessageReceived(msg, meta))
  }

  return store => next => action => {
    switch(action.type) {
      //The user wants us to connect
      case 'WS_CONNECT':
        //Start a new connection to the server
        if(socket != null) {
          socket.close()
        }

        //Send an action that shows a "connecting..." status for now
        store.dispatch(wsConnecting())

        //Attempt to connect (we could send a 'failed' action on error)
        socket = new WebSocket(action.url)
        socket.onmessage = onMessage(socket, store)
        socket.onclose = onClose(socket, store)
        socket.onopen = onOpen(socket, store, action.token)
        socket.onerror = function (_evt) {

        }
        break

      //The user wants us to disconnect
      case 'WS_DISCONNECT': {
        if(socket != null) {
          socket.close()
        }

        socket = null
        let timeout = 10000

        //Set our state to disconnected
        store.dispatch(wsDisconnected())

        // Fast check for the first 10 times then back-off
        if (reconnectCount < 10) {
          timeout = 1000
        }

        //Try to reconnect
        reconnect = setTimeout(() => {
          reconnectCount++
          store.dispatch(wsConnect())
        }, timeout)

        break
      }

      //Send the 'SEND_MESSAGE' action down the websocket to the server
      case 'WS_SEND_MESSAGE':
        socket.send(JSON.stringify(action))
        break

      //This action is irrelevant to us, pass it on to the next middleware
      default:
        return next(action)
    }
  }
})()

export default socketMiddleware
