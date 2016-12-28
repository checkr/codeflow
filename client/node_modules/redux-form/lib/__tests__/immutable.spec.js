'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _actionTypes = require('../actionTypes');

var expectedActionTypes = _interopRequireWildcard(_actionTypes);

var _propTypes = require('../propTypes');

var _propTypes2 = _interopRequireDefault(_propTypes);

var _immutable = require('../immutable');

function _interopRequireWildcard(obj) { if (obj && obj.__esModule) { return obj; } else { var newObj = {}; if (obj != null) { for (var key in obj) { if (Object.prototype.hasOwnProperty.call(obj, key)) newObj[key] = obj[key]; } } newObj.default = obj; return newObj; } }

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('immutable', function () {
  it('should export actionTypes', function () {
    (0, _expect2.default)(_immutable.actionTypes).toEqual(expectedActionTypes);
  });
  it('should export arrayInsert', function () {
    (0, _expect2.default)(_immutable.arrayInsert).toExist().toBeA('function');
  });
  it('should export arrayPop', function () {
    (0, _expect2.default)(_immutable.arrayPop).toExist().toBeA('function');
  });
  it('should export arrayPush', function () {
    (0, _expect2.default)(_immutable.arrayPush).toExist().toBeA('function');
  });
  it('should export arrayRemove', function () {
    (0, _expect2.default)(_immutable.arrayRemove).toExist().toBeA('function');
  });
  it('should export arrayShift', function () {
    (0, _expect2.default)(_immutable.arrayShift).toExist().toBeA('function');
  });
  it('should export arraySplice', function () {
    (0, _expect2.default)(_immutable.arraySplice).toExist().toBeA('function');
  });
  it('should export arraySwap', function () {
    (0, _expect2.default)(_immutable.arraySwap).toExist().toBeA('function');
  });
  it('should export arrayUnshift', function () {
    (0, _expect2.default)(_immutable.arrayUnshift).toExist().toBeA('function');
  });
  it('should export blur', function () {
    (0, _expect2.default)(_immutable.blur).toExist().toBeA('function');
  });
  it('should export change', function () {
    (0, _expect2.default)(_immutable.change).toExist().toBeA('function');
  });
  it('should export destroy', function () {
    (0, _expect2.default)(_immutable.destroy).toExist().toBeA('function');
  });
  it('should export Field', function () {
    (0, _expect2.default)(_immutable.Field).toExist().toBeA('function');
  });
  it('should export FieldArray', function () {
    (0, _expect2.default)(_immutable.FieldArray).toExist().toBeA('function');
  });
  it('should export focus', function () {
    (0, _expect2.default)(_immutable.focus).toExist().toBeA('function');
  });
  it('should export formValueSelector', function () {
    (0, _expect2.default)(_immutable.formValueSelector).toExist().toBeA('function');
  });
  it('should export initialize', function () {
    (0, _expect2.default)(_immutable.initialize).toExist().toBeA('function');
  });
  it('should export propTypes', function () {
    (0, _expect2.default)(_immutable.propTypes).toEqual(_propTypes2.default);
  });
  it('should export reducer', function () {
    (0, _expect2.default)(_immutable.reducer).toExist().toBeA('function');
  });
  it('should export reduxForm', function () {
    (0, _expect2.default)(_immutable.reduxForm).toExist().toBeA('function');
  });
  it('should export reset', function () {
    (0, _expect2.default)(_immutable.reset).toExist().toBeA('function');
  });
  it('should export startAsyncValidation', function () {
    (0, _expect2.default)(_immutable.startAsyncValidation).toExist().toBeA('function');
  });
  it('should export startSubmit', function () {
    (0, _expect2.default)(_immutable.startSubmit).toExist().toBeA('function');
  });
  it('should export setSubmitFailed', function () {
    (0, _expect2.default)(_immutable.setSubmitFailed).toExist().toBeA('function');
  });
  it('should export setSubmitSucceeded', function () {
    (0, _expect2.default)(_immutable.setSubmitSucceeded).toExist().toBeA('function');
  });
  it('should export stopAsyncValidation', function () {
    (0, _expect2.default)(_immutable.stopAsyncValidation).toExist().toBeA('function');
  });
  it('should export stopSubmit', function () {
    (0, _expect2.default)(_immutable.stopSubmit).toExist().toBeA('function');
  });
  it('should export SubmissionError', function () {
    (0, _expect2.default)(_immutable.SubmissionError).toExist().toBeA('function');
  });
  it('should export touch', function () {
    (0, _expect2.default)(_immutable.touch).toExist().toBeA('function');
  });
  it('should export untouch', function () {
    (0, _expect2.default)(_immutable.untouch).toExist().toBeA('function');
  });
  it('should export values', function () {
    (0, _expect2.default)(_immutable.values).toExist().toBeA('function');
  });
});