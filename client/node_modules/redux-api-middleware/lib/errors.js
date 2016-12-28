/**
 * Error class for an RSAA that does not conform to the RSAA definition
 *
 * @class InvalidRSAA
 * @access public
 * @param {array} validationErrors - an array of validation errors
 */
'use strict';

var _inherits = require('babel-runtime/helpers/inherits')['default'];

var _classCallCheck = require('babel-runtime/helpers/class-call-check')['default'];

exports.__esModule = true;

var InvalidRSAA = (function (_Error) {
  _inherits(InvalidRSAA, _Error);

  function InvalidRSAA(validationErrors) {
    _classCallCheck(this, InvalidRSAA);

    _Error.call(this);
    this.name = 'InvalidRSAA';
    this.message = 'Invalid RSAA';
    this.validationErrors = validationErrors;
  }

  /**
   * Error class for a custom `payload` or `meta` function throwing
   *
   * @class InternalError
   * @access public
   * @param {string} message - the error message
   */
  return InvalidRSAA;
})(Error);

var InternalError = (function (_Error2) {
  _inherits(InternalError, _Error2);

  function InternalError(message) {
    _classCallCheck(this, InternalError);

    _Error2.call(this);
    this.name = 'InternalError';
    this.message = message;
  }

  /**
   * Error class for an error raised trying to make an API call
   *
   * @class RequestError
   * @access public
   * @param {string} message - the error message
   */
  return InternalError;
})(Error);

var RequestError = (function (_Error3) {
  _inherits(RequestError, _Error3);

  function RequestError(message) {
    _classCallCheck(this, RequestError);

    _Error3.call(this);
    this.name = 'RequestError';
    this.message = message;
  }

  /**
   * Error class for an API response outside the 200 range
   *
   * @class ApiError
   * @access public
   * @param {number} status - the status code of the API response
   * @param {string} statusText - the status text of the API response
   * @param {object} response - the parsed JSON response of the API server if the
   *  'Content-Type' header signals a JSON response
   */
  return RequestError;
})(Error);

var ApiError = (function (_Error4) {
  _inherits(ApiError, _Error4);

  function ApiError(status, statusText, response) {
    _classCallCheck(this, ApiError);

    _Error4.call(this);
    this.name = 'ApiError';
    this.status = status;
    this.statusText = statusText;
    this.response = response;
    this.message = status + ' - ' + statusText;
  }

  return ApiError;
})(Error);

exports.InvalidRSAA = InvalidRSAA;
exports.InternalError = InternalError;
exports.RequestError = RequestError;
exports.ApiError = ApiError;