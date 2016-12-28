'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _splice = require('../splice');

var _splice2 = _interopRequireDefault(_splice);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('structure.plain.splice', function () {
  var testInsertWithValue = function testInsertWithValue(value) {
    it('should insert even when initial array is undefined', function () {
      (0, _expect2.default)((0, _splice2.default)(undefined, 2, 0, value)) // really goes to index 0
      .toBeA('array').toEqual([,, value]); // eslint-disable-line no-sparse-arrays
    });

    it('should insert ' + value + ' at start', function () {
      (0, _expect2.default)((0, _splice2.default)(['b', 'c', 'd'], 0, 0, value)).toBeA('array').toEqual([value, 'b', 'c', 'd']);
    });

    it('should insert ' + value + ' at end', function () {
      (0, _expect2.default)((0, _splice2.default)(['a', 'b', 'c'], 3, 0, value)).toBeA('array').toEqual(['a', 'b', 'c', value]);
    });

    it('should insert ' + value + ' in middle', function () {
      (0, _expect2.default)((0, _splice2.default)(['a', 'b', 'd'], 2, 0, value)).toBeA('array').toEqual(['a', 'b', value, 'd']);
    });

    it('should insert ' + value + ' when index is out of range', function () {
      (0, _expect2.default)((0, _splice2.default)(['a', 'b', 'c'], 5, 0, value)).toBeA('array').toEqual(['a', 'b', 'c',,, value]); // eslint-disable-line no-sparse-arrays
    });
  };

  testInsertWithValue('value');
  testInsertWithValue(undefined);

  it('should return empty array when removing and initial array is undefined', function () {
    (0, _expect2.default)((0, _splice2.default)(undefined, 2, 1)).toBeA('array').toEqual([]);
  });

  it('should remove at start', function () {
    (0, _expect2.default)((0, _splice2.default)(['a', 'b', 'c', 'd'], 0, 1)).toBeA('array').toEqual(['b', 'c', 'd']);
  });

  it('should remove in the middle then insert in that position', function () {
    (0, _expect2.default)((0, _splice2.default)(['a', 'b', 'c', 'd'], 1, 1, 'e')).toBeA('array').toEqual(['a', 'e', 'c', 'd']);
  });

  it('should remove at end', function () {
    (0, _expect2.default)((0, _splice2.default)(['a', 'b', 'c', 'd'], 3, 1)).toBeA('array').toEqual(['a', 'b', 'c']);
  });

  it('should remove in middle', function () {
    (0, _expect2.default)((0, _splice2.default)(['a', 'b', 'c', 'd'], 2, 1)).toBeA('array').toEqual(['a', 'b', 'd']);
  });
});