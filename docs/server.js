/**
 * HTTP Server
 *
 * Create a http server instance
 */

/* eslint-disable no-console */

var server = require('pushstate-server');
var build = "./_book"

var port = process.env.DOCS_PORT || 4000
var app = server.start({
  port: port,
  directory: build
});

process.on('SIGTERM', function () {
  console.log("SIGTERM recieved! Shutting down...");
  app.close();
});

console.log('HTTP Server listening on http://0.0.0.0:'+port)
