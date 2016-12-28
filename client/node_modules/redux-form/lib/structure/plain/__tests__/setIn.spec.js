'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _setIn = require('../setIn');

var _setIn2 = _interopRequireDefault(_setIn);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('structure.plain.setIn', function () {
  it('should create a map if state is undefined and key is string', function () {
    (0, _expect2.default)((0, _setIn2.default)(undefined, 'dog', 'fido')).toBeA('object').toEqual({ dog: 'fido' });
  });

  it('should create an array if state is undefined and key is string', function () {
    (0, _expect2.default)((0, _setIn2.default)(undefined, '[0]', 'fido')).toBeA(Array).toEqual(['fido']);
    var result = (0, _setIn2.default)(undefined, '[1]', 'second');
    (0, _expect2.default)(result).toBeA(Array);
    (0, _expect2.default)(result.length).toBe(2);
    (0, _expect2.default)(result[0]).toBe(undefined);
    (0, _expect2.default)(result[1]).toBe('second');
  });

  it('should set and shallow keys without mutating state', function () {
    var state = { foo: 'bar' };
    (0, _expect2.default)((0, _setIn2.default)(state, 'foo', 'baz')).toNotBe(state).toEqual({ foo: 'baz' });
    (0, _expect2.default)((0, _setIn2.default)(state, 'cat', 'fluffy')).toNotBe(state).toEqual({ foo: 'bar', cat: 'fluffy' });
    (0, _expect2.default)((0, _setIn2.default)(state, 'age', 42)).toNotBe(state).toEqual({ foo: 'bar', age: 42 });
  });

  it('should set and deep keys without mutating state', function () {
    var state = {
      foo: {
        bar: ['baz', { dog: 42 }]
      }
    };
    var result1 = (0, _setIn2.default)(state, 'tv.best.canines[0]', 'scooby');
    (0, _expect2.default)(result1).toNotBe(state).toEqual({
      foo: {
        bar: ['baz', { dog: 42 }]
      },
      tv: {
        best: {
          canines: ['scooby']
        }
      }
    });
    (0, _expect2.default)(result1.foo).toBe(state.foo);

    var result2 = (0, _setIn2.default)(state, 'foo.bar[0]', 'cat');
    (0, _expect2.default)(result2).toNotBe(state).toEqual({
      foo: {
        bar: ['cat', { dog: 42 }]
      }
    });
    (0, _expect2.default)(result2.foo).toNotBe(state.foo);
    (0, _expect2.default)(result2.foo.bar).toNotBe(state.foo.bar);
    (0, _expect2.default)(result2.foo.bar[1]).toBe(state.foo.bar[1]);

    var result3 = (0, _setIn2.default)(state, 'foo.bar[1].dog', 7);
    (0, _expect2.default)(result3).toNotBe(state).toEqual({
      foo: {
        bar: ['baz', { dog: 7 }]
      }
    });
    (0, _expect2.default)(result3.foo).toNotBe(state.foo);
    (0, _expect2.default)(result3.foo.bar).toNotBe(state.foo.bar);
    (0, _expect2.default)(result3.foo.bar[0]).toBe(state.foo.bar[0]);
    (0, _expect2.default)(result3.foo.bar[1]).toNotBe(state.foo.bar[1]);
  });
});