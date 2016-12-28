'use strict';

var _extends = require('babel-runtime/helpers/extends')['default'];

var _regeneratorRuntime = require('babel-runtime/regenerator')['default'];

var _interopRequireDefault = require('babel-runtime/helpers/interop-require-default')['default'];

exports.__esModule = true;

var _isomorphicFetch = require('isomorphic-fetch');

var _isomorphicFetch2 = _interopRequireDefault(_isomorphicFetch);

var _lodashIsplainobject = require('lodash.isplainobject');

var _lodashIsplainobject2 = _interopRequireDefault(_lodashIsplainobject);

var _CALL_API = require('./CALL_API');

var _CALL_API2 = _interopRequireDefault(_CALL_API);

var _validation = require('./validation');

var _errors = require('./errors');

var _util = require('./util');

/**
 * A Redux middleware that processes RSAA actions.
 *
 * @type {ReduxMiddleware}
 * @access public
 */
function apiMiddleware(_ref) {
  var _this = this;

  var getState = _ref.getState;

  return function (next) {
    return function callee$2$0(action) {
      var validationErrors, _callAPI, _requestType, callAPI, endpoint, headers, method, body, credentials, bailout, types, _normalizeTypeDescriptors, requestType, successType, failureType, res;

      return _regeneratorRuntime.async(function callee$2$0$(context$3$0) {
        while (1) switch (context$3$0.prev = context$3$0.next) {
          case 0:
            if (_validation.isRSAA(action)) {
              context$3$0.next = 2;
              break;
            }

            return context$3$0.abrupt('return', next(action));

          case 2:
            validationErrors = _validation.validateRSAA(action);

            if (!validationErrors.length) {
              context$3$0.next = 7;
              break;
            }

            _callAPI = action[_CALL_API2['default']];

            if (_callAPI.types && Array.isArray(_callAPI.types)) {
              _requestType = _callAPI.types[0];

              if (_requestType && _requestType.type) {
                _requestType = _requestType.type;
              }
              next({
                type: _requestType,
                payload: new _errors.InvalidRSAA(validationErrors),
                error: true
              });
            }
            return context$3$0.abrupt('return');

          case 7:
            callAPI = action[_CALL_API2['default']];
            endpoint = callAPI.endpoint;
            headers = callAPI.headers;
            method = callAPI.method;
            body = callAPI.body;
            credentials = callAPI.credentials;
            bailout = callAPI.bailout;
            types = callAPI.types;
            _normalizeTypeDescriptors = _util.normalizeTypeDescriptors(types);
            requestType = _normalizeTypeDescriptors[0];
            successType = _normalizeTypeDescriptors[1];
            failureType = _normalizeTypeDescriptors[2];
            context$3$0.prev = 19;

            if (!(typeof bailout === 'boolean' && bailout || typeof bailout === 'function' && bailout(getState()))) {
              context$3$0.next = 22;
              break;
            }

            return context$3$0.abrupt('return');

          case 22:
            context$3$0.next = 30;
            break;

          case 24:
            context$3$0.prev = 24;
            context$3$0.t0 = context$3$0['catch'](19);
            context$3$0.next = 28;
            return _regeneratorRuntime.awrap(_util.actionWith(_extends({}, requestType, {
              payload: new _errors.RequestError('[CALL_API].bailout function failed'),
              error: true
            }), [action, getState()]));

          case 28:
            context$3$0.t1 = context$3$0.sent;
            return context$3$0.abrupt('return', next(context$3$0.t1));

          case 30:
            if (!(typeof endpoint === 'function')) {
              context$3$0.next = 41;
              break;
            }

            context$3$0.prev = 31;

            endpoint = endpoint(getState());
            context$3$0.next = 41;
            break;

          case 35:
            context$3$0.prev = 35;
            context$3$0.t2 = context$3$0['catch'](31);
            context$3$0.next = 39;
            return _regeneratorRuntime.awrap(_util.actionWith(_extends({}, requestType, {
              payload: new _errors.RequestError('[CALL_API].endpoint function failed'),
              error: true
            }), [action, getState()]));

          case 39:
            context$3$0.t3 = context$3$0.sent;
            return context$3$0.abrupt('return', next(context$3$0.t3));

          case 41:
            if (!(typeof headers === 'function')) {
              context$3$0.next = 52;
              break;
            }

            context$3$0.prev = 42;

            headers = headers(getState());
            context$3$0.next = 52;
            break;

          case 46:
            context$3$0.prev = 46;
            context$3$0.t4 = context$3$0['catch'](42);
            context$3$0.next = 50;
            return _regeneratorRuntime.awrap(_util.actionWith(_extends({}, requestType, {
              payload: new _errors.RequestError('[CALL_API].headers function failed'),
              error: true
            }), [action, getState()]));

          case 50:
            context$3$0.t5 = context$3$0.sent;
            return context$3$0.abrupt('return', next(context$3$0.t5));

          case 52:
            context$3$0.next = 54;
            return _regeneratorRuntime.awrap(_util.actionWith(requestType, [action, getState()]));

          case 54:
            context$3$0.t6 = context$3$0.sent;
            next(context$3$0.t6);
            context$3$0.prev = 56;
            context$3$0.next = 59;
            return _regeneratorRuntime.awrap(_isomorphicFetch2['default'](endpoint, { method: method, body: body, credentials: credentials, headers: headers }));

          case 59:
            res = context$3$0.sent;
            context$3$0.next = 68;
            break;

          case 62:
            context$3$0.prev = 62;
            context$3$0.t7 = context$3$0['catch'](56);
            context$3$0.next = 66;
            return _regeneratorRuntime.awrap(_util.actionWith(_extends({}, requestType, {
              payload: new _errors.RequestError(context$3$0.t7.message),
              error: true
            }), [action, getState()]));

          case 66:
            context$3$0.t8 = context$3$0.sent;
            return context$3$0.abrupt('return', next(context$3$0.t8));

          case 68:
            if (!res.ok) {
              context$3$0.next = 75;
              break;
            }

            context$3$0.next = 71;
            return _regeneratorRuntime.awrap(_util.actionWith(successType, [action, getState(), res]));

          case 71:
            context$3$0.t9 = context$3$0.sent;
            return context$3$0.abrupt('return', next(context$3$0.t9));

          case 75:
            context$3$0.next = 77;
            return _regeneratorRuntime.awrap(_util.actionWith(_extends({}, failureType, {
              error: true
            }), [action, getState(), res]));

          case 77:
            context$3$0.t10 = context$3$0.sent;
            return context$3$0.abrupt('return', next(context$3$0.t10));

          case 79:
          case 'end':
            return context$3$0.stop();
        }
      }, null, _this, [[19, 24], [31, 35], [42, 46], [56, 62]]);
    };
  };
}

exports.apiMiddleware = apiMiddleware;

// Do not process actions without a [CALL_API] property

// Try to dispatch an error request FSA for invalid RSAAs

// Parse the validated RSAA action

// Should we bail out?

// Process [CALL_API].endpoint function

// Process [CALL_API].headers function

// We can now dispatch the request FSA

// Make the API call

// The request was malformed, or there was a network error

// Process the server response