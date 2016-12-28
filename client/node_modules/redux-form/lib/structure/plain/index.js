'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _some2 = require('lodash/some');

var _some3 = _interopRequireDefault(_some2);

var _splice = require('./splice');

var _splice2 = _interopRequireDefault(_splice);

var _getIn = require('./getIn');

var _getIn2 = _interopRequireDefault(_getIn);

var _setIn = require('./setIn');

var _setIn2 = _interopRequireDefault(_setIn);

var _deepEqual = require('./deepEqual');

var _deepEqual2 = _interopRequireDefault(_deepEqual);

var _deleteIn = require('./deleteIn');

var _deleteIn2 = _interopRequireDefault(_deleteIn);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var structure = {
  empty: {},
  emptyList: [],
  getIn: _getIn2.default,
  setIn: _setIn2.default,
  deepEqual: _deepEqual2.default,
  deleteIn: _deleteIn2.default,
  fromJS: function fromJS(value) {
    return value;
  },
  size: function size(array) {
    return array ? array.length : 0;
  },
  some: _some3.default,
  splice: _splice2.default
};

exports.default = structure;