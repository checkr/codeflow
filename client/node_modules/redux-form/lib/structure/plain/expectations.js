'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _isObject2 = require('lodash/isObject');

var _isObject3 = _interopRequireDefault(_isObject2);

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var expectations = {
  toBeAMap: function toBeAMap() {
    _expect2.default.assert((0, _isObject3.default)(this.actual), 'expected %s to be an object', this.actual);
    return this;
  },
  toBeAList: function toBeAList() {
    _expect2.default.assert(Array.isArray(this.actual), 'expected %s to be an array', this.actual);
    return this;
  },
  toBeSize: function toBeSize(size) {
    _expect2.default.assert(this.actual && Object.keys(this.actual).length === size, 'expected %s to contain %s elements', this.actual, size);
    return this;
  },
  toEqualMap: function toEqualMap(expected) {
    return (0, _expect2.default)(this.actual).toEqual(expected);
  }
};

exports.default = expectations;