/**
 * Codeflow HTTP Server
 *
 * Create a http server instance
 * Before serving it also replaces `##{REACT_APP_*}##` environment variables
 * and ouputs files to `build-env` directory
 */

/* eslint-disable no-console */

var server = require('pushstate-server');
var ncp = require('ncp').ncp;
var path = require('path');
var extname = path.extname;
var replaceStream = require('replacestream');
var build = __dirname + "/build"
var buildEnv = __dirname + "/build-env"
var config = {}

// config
Object.keys(process.env)
  .filter(key => (/^REACT_APP_/i).test(key))
  .map(env => config[env] = process.env[env])

var options = {
  transform: function (read, write, file) { 
    if (extname(file.name) === '.js') {
      read = read.pipe(replaceStream('##JSON_STRING_CONFIG##', escape(JSON.stringify(config))))
    }
    read.pipe(write)
  } 
}

ncp(build, buildEnv, options, function (err) {
 if (err) {
   console.log(err)
   process.exit();
 }
});

var port = process.env.REACT_APP_PORT || 3000
var app = server.start({
  port: port,
  directory: buildEnv
});

process.on('SIGTERM', function () {
  console.log("SIGTERM recieved! Shutting down...");
  app.close();
});

console.log('Codeflow HTTP Server listening on http://0.0.0.0:'+port)
