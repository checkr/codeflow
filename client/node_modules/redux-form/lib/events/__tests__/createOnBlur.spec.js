'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _createOnBlur = require('../createOnBlur');

var _createOnBlur2 = _interopRequireDefault(_createOnBlur);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('createOnBlur', function () {
  it('should return a function', function () {
    (0, _expect2.default)((0, _createOnBlur2.default)()).toExist().toBeA('function');
  });

  it('should parse the value before dispatching action', function () {
    var blur = (0, _expect.createSpy)();
    var parse = (0, _expect.createSpy)(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    (0, _createOnBlur2.default)(blur, { parse: parse })('bar');
    (0, _expect2.default)(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    (0, _expect2.default)(blur).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
  });

  it('should normalize the value before dispatching action', function () {
    var blur = (0, _expect.createSpy)();
    var normalize = (0, _expect.createSpy)(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    (0, _createOnBlur2.default)(blur, { normalize: normalize })('bar');
    (0, _expect2.default)(normalize).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    (0, _expect2.default)(blur).toHaveBeenCalled().toHaveBeenCalledWith('normalized-bar');
  });

  it('should parse before normalize', function () {
    var blur = (0, _expect.createSpy)();
    var parse = (0, _expect.createSpy)(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    var normalize = (0, _expect.createSpy)(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    (0, _createOnBlur2.default)(blur, { normalize: normalize, parse: parse })('bar');
    (0, _expect2.default)(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    (0, _expect2.default)(normalize).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
    (0, _expect2.default)(blur).toHaveBeenCalled().toHaveBeenCalledWith('normalized-parsed-bar');
  });

  it('should call blur then after', function () {
    var blur = (0, _expect.createSpy)();
    var parse = (0, _expect.createSpy)(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    var normalize = (0, _expect.createSpy)(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    var after = (0, _expect.createSpy)();
    (0, _createOnBlur2.default)(blur, { parse: parse, normalize: normalize, after: after })('bar');
    (0, _expect2.default)(blur).toHaveBeenCalled();
    (0, _expect2.default)(normalize).toHaveBeenCalled();
    (0, _expect2.default)(after).toHaveBeenCalled().toHaveBeenCalledWith('normalized-parsed-bar');
  });
});