'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _deleteIn = require('../deleteIn');

var _deleteIn2 = _interopRequireDefault(_deleteIn);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('structure.plain.deleteIn', function () {
  it('should not return state if path not found', function () {
    var state = { foo: 'bar' };
    (0, _expect2.default)((0, _deleteIn2.default)(state, undefined)).toBe(state);
    (0, _expect2.default)((0, _deleteIn2.default)(state, 'dog')).toBe(state);
    (0, _expect2.default)((0, _deleteIn2.default)(state, 'cat.rat.pig')).toBe(state);
  });

  it('should do nothing if array index out of bounds', function () {
    var state = {
      foo: [{
        bar: ['dog']
      }]
    };
    (0, _expect2.default)((0, _deleteIn2.default)(state, 'foo[2].bar[0]')).toEqual(state);
    (0, _expect2.default)((0, _deleteIn2.default)(state, 'foo[0].bar[2]')).toEqual(state);
  });

  it('should throw exception for non-numerical array indexes', function () {
    var state = {
      foo: ['dog']
    };
    (0, _expect2.default)(function () {
      return (0, _deleteIn2.default)(state, 'foo[bar]');
    }).toThrow(/non\-numerical index/);
  });

  it('should delete shallow keys without mutating state', function () {
    var state = { foo: 'bar', dog: 'fido' };
    (0, _expect2.default)((0, _deleteIn2.default)(state, 'foo')).toNotBe(state).toEqual({ dog: 'fido' });
    (0, _expect2.default)((0, _deleteIn2.default)(state, 'dog')).toNotBe(state).toEqual({ foo: 'bar' });
  });

  it('should delete shallow array indexes without mutating state', function () {
    var state = ['the', 'quick', 'brown', 'fox'];
    (0, _expect2.default)((0, _deleteIn2.default)(state, 4)).toBe(state); // index not found
    (0, _expect2.default)((0, _deleteIn2.default)(state, 0)).toNotBe(state).toEqual(['quick', 'brown', 'fox']);
    (0, _expect2.default)((0, _deleteIn2.default)(state, 1)).toNotBe(state).toEqual(['the', 'brown', 'fox']);
    (0, _expect2.default)((0, _deleteIn2.default)(state, 2)).toNotBe(state).toEqual(['the', 'quick', 'fox']);
    (0, _expect2.default)((0, _deleteIn2.default)(state, 3)).toNotBe(state).toEqual(['the', 'quick', 'brown']);
  });

  it('should delete deep keys without mutating state', function () {
    var state = {
      foo: {
        bar: ['baz', { dog: 42 }]
      }
    };

    var result1 = (0, _deleteIn2.default)(state, 'foo.bar[0]');
    (0, _expect2.default)(result1).toNotBe(state).toEqual({
      foo: {
        bar: [{ dog: 42 }]
      }
    });
    (0, _expect2.default)(result1.foo).toNotBe(state.foo);
    (0, _expect2.default)(result1.foo.bar).toNotBe(state.foo.bar);
    (0, _expect2.default)(result1.foo.bar.length).toBe(1);
    (0, _expect2.default)(result1.foo.bar[0]).toBe(state.foo.bar[1]);

    var result2 = (0, _deleteIn2.default)(state, 'foo.bar[1].dog');
    (0, _expect2.default)(result2).toNotBe(state).toEqual({
      foo: {
        bar: ['baz', {}]
      }
    });
    (0, _expect2.default)(result2.foo).toNotBe(state.foo);
    (0, _expect2.default)(result2.foo.bar).toNotBe(state.foo.bar);
    (0, _expect2.default)(result2.foo.bar[0]).toBe(state.foo.bar[0]);
    (0, _expect2.default)(result2.foo.bar[1]).toNotBe(state.foo.bar[1]);

    var result3 = (0, _deleteIn2.default)(state, 'foo.bar');
    (0, _expect2.default)(result3).toNotBe(state).toEqual({
      foo: {}
    });
    (0, _expect2.default)(result3.foo).toNotBe(state.foo);
  });

  it('should not mutate deep state if can\'t find final key', function () {
    var state = {
      foo: {
        bar: [{}]
      }
    };
    var result = (0, _deleteIn2.default)(state, 'foo.bar[0].dog');
    (0, _expect2.default)(result).toBe(state).toEqual({
      foo: {
        bar: [{}]
      }
    });
    (0, _expect2.default)(result.foo).toBe(state.foo);
    (0, _expect2.default)(result.foo.bar).toBe(state.foo.bar);
    (0, _expect2.default)(result.foo.bar[0]).toBe(state.foo.bar[0]);
  });
});