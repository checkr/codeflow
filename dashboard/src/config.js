/* eslint-disable no-console */

import _ from 'underscore'

const loadConfig = () => {
  let config = {}

  if (process.env.NODE_ENV === 'development') {
    config = {
      "REACT_APP_ROOT": "http://"+location.hostname+":"+process.env.REACT_APP_ROOT_PORT,
      "REACT_APP_API_ROOT": "http://"+location.hostname+":"+process.env.REACT_APP_API_ROOT_PORT,
      "REACT_APP_WS_ROOT": "ws://"+location.hostname+":"+process.env.REACT_APP_WS_ROOT_PORT,
      "REACT_APP_WEBHOOKS_ROOT": "http://"+location.hostname+":"+process.env.REACT_APP_WEBHOOKS_ROOT_PORT,
      "REACT_APP_OKTA_CLIENT_ID": process.env.REACT_APP_OKTA_CLIENT_ID,
      "REACT_APP_OKTA_URL": process.env.REACT_APP_OKTA_URL,
      "REACT_APP_OKTA_LOGO": process.env.REACT_APP_OKTA_LOGO,
    }
  } else {
    try{
      config = JSON.parse(unescape("##JSON_STRING_CONFIG##"))
    }catch(e){
      console.log(e);
    }
  }
  return _.extend(config, process.env)
}

export default loadConfig
