import getValue from './getValue';
import isReactNative from '../isReactNative';

var createOnChange = function createOnChange(change) {
  var _ref = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {},
      parse = _ref.parse,
      normalize = _ref.normalize;

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

    // dispatch change action
    change(value);
  };
};

export default createOnChange;