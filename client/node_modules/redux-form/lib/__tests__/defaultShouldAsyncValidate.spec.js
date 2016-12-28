'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _defaultShouldAsyncValidate = require('../defaultShouldAsyncValidate');

var _defaultShouldAsyncValidate2 = _interopRequireDefault(_defaultShouldAsyncValidate);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('defaultShouldAsyncValidate', function () {

  it('should not async validate if sync validation is not passing', function () {
    (0, _expect2.default)((0, _defaultShouldAsyncValidate2.default)({
      syncValidationPasses: false
    })).toBe(false);
  });

  it('should async validate if blur triggered and sync passes', function () {
    (0, _expect2.default)((0, _defaultShouldAsyncValidate2.default)({
      syncValidationPasses: true,
      trigger: 'blur'
    })).toBe(true);
  });

  it('should not async validate when pristine and initialized', function () {
    (0, _expect2.default)((0, _defaultShouldAsyncValidate2.default)({
      syncValidationPasses: true,
      trigger: 'submit',
      pristine: true,
      initialized: true
    })).toBe(false);
  });

  it('should async validate when submitting and dirty', function () {
    (0, _expect2.default)((0, _defaultShouldAsyncValidate2.default)({
      syncValidationPasses: true,
      trigger: 'submit',
      pristine: false,
      initialized: true
    })).toBe(true);
  });

  it('should async validate when submitting and not initialized', function () {
    (0, _expect2.default)((0, _defaultShouldAsyncValidate2.default)({
      syncValidationPasses: true,
      trigger: 'submit',
      pristine: true,
      initialized: false
    })).toBe(true);
  });

  it('should not async validate when unknown trigger', function () {
    (0, _expect2.default)((0, _defaultShouldAsyncValidate2.default)({
      syncValidationPasses: true,
      trigger: 'wtf'
    })).toBe(false);
  });
});