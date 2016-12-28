'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _createOnChange = require('../createOnChange');

var _createOnChange2 = _interopRequireDefault(_createOnChange);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('createOnChange', function () {
  it('should return a function', function () {
    (0, _expect2.default)((0, _createOnChange2.default)()).toExist().toBeA('function');
  });

  it('should parse the value before dispatching action', function () {
    var change = (0, _expect.createSpy)();
    var parse = (0, _expect.createSpy)(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    (0, _createOnChange2.default)(change, { parse: parse })('bar');
    (0, _expect2.default)(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    (0, _expect2.default)(change).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
  });

  it('should normalize the value before dispatching action', function () {
    var change = (0, _expect.createSpy)();
    var normalize = (0, _expect.createSpy)(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    (0, _createOnChange2.default)(change, { normalize: normalize })('bar');
    (0, _expect2.default)(normalize).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    (0, _expect2.default)(change).toHaveBeenCalled().toHaveBeenCalledWith('normalized-bar');
  });

  it('should parse before normalize', function () {
    var change = (0, _expect.createSpy)();
    var parse = (0, _expect.createSpy)(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    var normalize = (0, _expect.createSpy)(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    (0, _createOnChange2.default)(change, { normalize: normalize, parse: parse })('bar');
    (0, _expect2.default)(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    (0, _expect2.default)(normalize).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
    (0, _expect2.default)(change).toHaveBeenCalled().toHaveBeenCalledWith('normalized-parsed-bar');
  });
});