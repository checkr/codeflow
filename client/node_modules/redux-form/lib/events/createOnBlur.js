'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _getValue = require('./getValue');

var _getValue2 = _interopRequireDefault(_getValue);

var _isReactNative = require('../isReactNative');

var _isReactNative2 = _interopRequireDefault(_isReactNative);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var createOnBlur = function createOnBlur(blur) {
  var _ref = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {},
      after = _ref.after,
      normalize = _ref.normalize,
      parse = _ref.parse;

  return function (event) {
    // read value from input
    var value = (0, _getValue2.default)(event, _isReactNative2.default);

    // parse value if we have a parser
    if (parse) {
      value = parse(value);
    }

    // normalize value
    if (normalize) {
      value = normalize(value);
    }

    // dispatch blur action
    blur(value);

    // call after callback
    if (after) {
      after(value);
    }
  };
};

exports.default = createOnBlur;