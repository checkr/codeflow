'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _getIn = require('../getIn');

var _getIn2 = _interopRequireDefault(_getIn);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('structure.plain.getIn', function () {
  it('should return undefined if state is undefined', function () {
    (0, _expect2.default)((0, _getIn2.default)(undefined, 'dog')).toBe(undefined);
  });

  it('should return undefined if any step on the path is undefined', function () {
    (0, _expect2.default)((0, _getIn2.default)({
      a: {
        b: {}
      }
    }, 'a.b.c')).toBe(undefined);
  });

  it('should get shallow values', function () {
    (0, _expect2.default)((0, _getIn2.default)({ foo: 'bar' }, 'foo')).toBe('bar');
    (0, _expect2.default)((0, _getIn2.default)({ foo: 42 }, 'foo')).toBe(42);
    (0, _expect2.default)((0, _getIn2.default)({ foo: false }, 'foo')).toBe(false);
  });

  it('should get deep values', function () {
    var state = {
      foo: {
        bar: ['baz', { dog: 42 }]
      }
    };
    (0, _expect2.default)((0, _getIn2.default)(state, 'foo.bar[0]')).toBe('baz');
    (0, _expect2.default)((0, _getIn2.default)(state, 'foo.bar[1].dog')).toBe(42);
  });

  it('should get a value nested 1 level', function () {
    (0, _expect2.default)((0, _getIn2.default)({ foo: { bar: 42 } }, 'foo.bar')).toBe(42);
  });

  it('should get a value nested 2 levels', function () {
    (0, _expect2.default)((0, _getIn2.default)({ foo: { bar: { baz: 42 } } }, 'foo.bar.baz')).toBe(42);
  });

  it('should get a value nested 3 levels', function () {
    (0, _expect2.default)((0, _getIn2.default)({ foo: { bar: { baz: { yolanda: 42 } } } }, 'foo.bar.baz.yolanda')).toBe(42);
  });

  it('should return undefined if the requested level does not exist', function () {
    (0, _expect2.default)((0, _getIn2.default)({}, 'foo')).toBe(undefined);
    (0, _expect2.default)((0, _getIn2.default)({}, 'foo.bar')).toBe(undefined);
    (0, _expect2.default)((0, _getIn2.default)({}, 'foo.bar.baz')).toBe(undefined);
    (0, _expect2.default)((0, _getIn2.default)({}, 'foo.bar.baz.yolanda')).toBe(undefined);
  });

  it('should return undefined for invalid/empty path', function () {
    (0, _expect2.default)((0, _getIn2.default)({ foo: 42 }, undefined)).toBe(undefined);
    (0, _expect2.default)((0, _getIn2.default)({ foo: 42 }, null)).toBe(undefined);
    (0, _expect2.default)((0, _getIn2.default)({ foo: 42 }, '')).toBe(undefined);
  });

  it('should get string keys on arrays', function () {
    var array = [1, 2, 3];
    array.stringKey = 'hello';
    var state = {
      foo: {
        bar: array
      }
    };
    (0, _expect2.default)((0, _getIn2.default)(state, 'foo.bar[0]')).toBe(1);
    (0, _expect2.default)((0, _getIn2.default)(state, 'foo.bar[1]')).toBe(2);
    (0, _expect2.default)((0, _getIn2.default)(state, 'foo.bar[2]')).toBe(3);
    (0, _expect2.default)((0, _getIn2.default)(state, 'foo.bar.stringKey')).toBe('hello');
  });
});