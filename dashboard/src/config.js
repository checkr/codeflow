/* eslint-disable no-console */

import _ from 'underscore'

const loadConfig = () => {
  let config = {}
  if (process.env.NODE_ENV === 'development') {
    config = {
      "REACT_APP_ROOT": "http://localhost:3000",
      "REACT_APP_API_ROOT": "http://localhost:3001",
      "REACT_APP_WS_ROOT": "ws://localhost:3003",
      "REACT_APP_WEBHOOKS_ROOT": "http://localhost:3002",
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
