'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _shallowequal = require('shallowequal');

var _shallowequal2 = _interopRequireDefault(_shallowequal);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var shallowCompare = function shallowCompare(instance, nextProps, nextState) {
  return !(0, _shallowequal2.default)(instance.props, nextProps) || !(0, _shallowequal2.default)(instance.state, nextState);
};

exports.default = shallowCompare;