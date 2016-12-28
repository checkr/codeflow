import getValue from './getValue';
import isReactNative from '../isReactNative';

var createOnBlur = function createOnBlur(blur) {
  var _ref = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {},
      after = _ref.after,
      normalize = _ref.normalize,
      parse = _ref.parse;

  return function (event) {
    // read value from input
    var value = getValue(event, isReactNative);

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

export default createOnBlur;