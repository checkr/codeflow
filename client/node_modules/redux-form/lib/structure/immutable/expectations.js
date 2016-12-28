'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _deepEqual = require('deep-equal');

var _deepEqual2 = _interopRequireDefault(_deepEqual);

var _immutable = require('immutable');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var deepEqualValues = function deepEqualValues(a, b) {
  if (_immutable.Iterable.isIterable(a)) {
    return _immutable.Iterable.isIterable(b) && a.count() === b.count() && a.every(function (value, key) {
      return deepEqualValues(value, b.get(key));
    });
  }
  return (0, _deepEqual2.default)(a, b); // neither are immutable iterables
};

var api = {
  toBeAMap: function toBeAMap() {
    _expect2.default.assert(_immutable.Map.isMap(this.actual), 'expected %s to be an immutable Map', this.actual);
    return this;
  },
  toBeAList: function toBeAList() {
    _expect2.default.assert(_immutable.List.isList(this.actual), 'expected %s to be an immutable List', this.actual);
    return this;
  },
  toBeSize: function toBeSize(size) {
    _expect2.default.assert(_immutable.Iterable.isIterable(this.actual) && this.actual.count() === size, 'expected %s to contain %s elements', this.actual, size);
    return this;
  },
  toEqualMap: function toEqualMap(expected) {
    _expect2.default.assert(deepEqualValues(this.actual, (0, _immutable.fromJS)(expected)), 'expected...\n%s\n...but found...\n%s', (0, _immutable.fromJS)(expected), this.actual);
    return this;
  }
};

exports.default = api;