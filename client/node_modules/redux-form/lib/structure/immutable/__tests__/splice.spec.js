'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _immutable = require('immutable');

var _splice = require('../splice');

var _splice2 = _interopRequireDefault(_splice);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('structure.immutable.splice', function () {
  var testInsertWithValue = function testInsertWithValue(value) {
    it('should insert even when initial array is undefined', function () {
      (0, _expect2.default)((0, _splice2.default)(undefined, 2, 0, value)) // really goes to index 0
      .toBeA(_immutable.List).toEqual((0, _immutable.fromJS)([,, value])); // eslint-disable-line no-sparse-arrays
    });

    it('should insert at start', function () {
      (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['b', 'c', 'd']), 0, 0, value)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)([value, 'b', 'c', 'd']));
    });

    it('should insert at end', function () {
      (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['a', 'b', 'c']), 3, 0, value)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['a', 'b', 'c', value]));
    });

    it('should insert in middle', function () {
      (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['a', 'b', 'd']), 2, 0, value)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['a', 'b', value, 'd']));
    });

    it('should insert in out of range', function () {
      (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['a', 'b', 'c']), 5, 0, value)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['a', 'b', 'c',,, value])); // eslint-disable-line no-sparse-arrays
    });
  };

  testInsertWithValue('value');
  testInsertWithValue(undefined);

  it('should return empty array when removing and initial array is undefined', function () {
    (0, _expect2.default)((0, _splice2.default)(undefined, 2, 1)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)([]));
  });

  it('should remove at start', function () {
    (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['a', 'b', 'c', 'd']), 0, 1)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['b', 'c', 'd']));
  });

  it('should remove in the middle then insert in that position', function () {
    (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['a', 'b', 'c', 'd']), 1, 1, 'e')).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['a', 'e', 'c', 'd']));
  });

  it('should remove at end', function () {
    (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['a', 'b', 'c', 'd']), 3, 1)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['a', 'b', 'c']));
  });

  it('should remove in middle', function () {
    (0, _expect2.default)((0, _splice2.default)((0, _immutable.fromJS)(['a', 'b', 'c', 'd']), 2, 1)).toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['a', 'b', 'd']));
  });
});