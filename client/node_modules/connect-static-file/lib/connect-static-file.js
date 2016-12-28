'use strict';

var send = require('send');
var accepts = require('accepts');
var mime = require('mime');

module.exports = function connectStaticFile(path, options) {
        options = options || {};

        var sendOptions = {
                dotfiles    : 'allow',
                etag        : options.etag,
                extensions  : options.extensions,
                index       : false,
                lastModified: options.lastModified,
                maxAge      : options.maxAge,
                root        : path
        };

        var encoded = options.encoded;
        var headers = [];

        if (options.headers) {
                Object.keys(options.headers).forEach(function(name) {
                        headers.push({
                                name : name,
                                value: options.headers[name]
                        });
                });
        }

        function connectStaticFileMiddleware(request, response, next) {
                function onError(err) {
                        if (err.code === 'ENOENT') {
                                // file not found, go to the next middleware without error
                                next();

                                return;
                        }

                        next(err);
                }

                function onDirectory() {
                        next();
                }

                function onHeaders(response, path, stat) {
                        for (var i = 0; i < headers.length; i++) {
                                response.setHeader(headers[i].name, headers[i].value);
                        }

                        if (encoded) {
                                response.setHeader('Content-Encoding', encoded);
                        }

                        if (!response.getHeader('Content-Type') && encoded) {
                                // foo.css.gz -> foo.css
                                var encodedPath = path;
                                encodedPath = encodedPath.replace(/\.(?:gz|gzip|zlib|bz2|xz)$/i, '');

                                var type = mime.lookup(encodedPath);
                                var charset = mime.charsets.lookup(type);

                                response.setHeader('Content-Type', type + (charset ? '; charset=' + charset : ''));
                        }
                }

                if (encoded) {
                        var accept = accepts(request);

                        var method = accept.encodings([encoded]);

                        if (method !== encoded) {
                                next();

                                return;
                        }
                }

                send(request, '', sendOptions)
                        .on('error', onError)
                        .on('directory', onDirectory)
                        .on('headers', onHeaders)
                        .pipe(response);
        }

        return connectStaticFileMiddleware;
};
