'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _noop2 = require('lodash/noop');

var _noop3 = _interopRequireDefault(_noop2);

var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

var _createOnBlur = require('./events/createOnBlur');

var _createOnBlur2 = _interopRequireDefault(_createOnBlur);

var _createOnChange = require('./events/createOnChange');

var _createOnChange2 = _interopRequireDefault(_createOnChange);

var _createOnDragStart = require('./events/createOnDragStart');

var _createOnDragStart2 = _interopRequireDefault(_createOnDragStart);

var _createOnDrop = require('./events/createOnDrop');

var _createOnDrop2 = _interopRequireDefault(_createOnDrop);

var _createOnFocus = require('./events/createOnFocus');

var _createOnFocus2 = _interopRequireDefault(_createOnFocus);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _objectWithoutProperties(obj, keys) { var target = {}; for (var i in obj) { if (keys.indexOf(i) >= 0) continue; if (!Object.prototype.hasOwnProperty.call(obj, i)) continue; target[i] = obj[i]; } return target; }

var processProps = function processProps(type, props, _value) {
  var value = props.value;

  if (type === 'checkbox') {
    return _extends({}, props, {
      checked: !!value
    });
  }
  if (type === 'radio') {
    return _extends({}, props, {
      checked: value === _value,
      value: _value
    });
  }
  if (type === 'select-multiple') {
    return _extends({}, props, {
      value: value || []
    });
  }
  if (type === 'file') {
    return _extends({}, props, {
      value: undefined
    });
  }
  return props;
};

var createFieldProps = function createFieldProps(getIn, name, _ref) {
  var asyncError = _ref.asyncError,
      asyncValidating = _ref.asyncValidating,
      blur = _ref.blur,
      change = _ref.change,
      dirty = _ref.dirty,
      dispatch = _ref.dispatch,
      focus = _ref.focus,
      format = _ref.format,
      normalize = _ref.normalize,
      parse = _ref.parse,
      pristine = _ref.pristine,
      props = _ref.props,
      state = _ref.state,
      submitError = _ref.submitError,
      submitting = _ref.submitting,
      value = _ref.value,
      _value = _ref._value,
      syncError = _ref.syncError,
      syncWarning = _ref.syncWarning,
      custom = _objectWithoutProperties(_ref, ['asyncError', 'asyncValidating', 'blur', 'change', 'dirty', 'dispatch', 'focus', 'format', 'normalize', 'parse', 'pristine', 'props', 'state', 'submitError', 'submitting', 'value', '_value', 'syncError', 'syncWarning']);

  var asyncValidate = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : _noop3.default;

  var error = syncError || asyncError || submitError;
  var warning = syncWarning;
  var boundParse = parse && function (value) {
    return parse(value, name);
  };
  var boundNormalize = normalize && function (value) {
    return normalize(name, value);
  };
  var boundChange = function boundChange(value) {
    return dispatch(change(name, value));
  };
  var onChange = (0, _createOnChange2.default)(boundChange, {
    normalize: boundNormalize,
    parse: boundParse
  });

  var formatFieldValue = function formatFieldValue(value, format) {
    if (format === null) {
      return value;
    }
    var defaultFormattedValue = value == null ? '' : value;
    return format ? format(value, name) : defaultFormattedValue;
  };

  var formattedFieldValue = formatFieldValue(value, format);

  return {
    input: processProps(custom.type, {
      name: name,
      onBlur: (0, _createOnBlur2.default)(function (value) {
        return dispatch(blur(name, value));
      }, {
        normalize: boundNormalize,
        parse: boundParse,
        after: asyncValidate.bind(null, name)
      }),
      onChange: onChange,
      onDragStart: (0, _createOnDragStart2.default)(name, formattedFieldValue),
      onDrop: (0, _createOnDrop2.default)(name, boundChange),
      onFocus: (0, _createOnFocus2.default)(name, function () {
        return dispatch(focus(name));
      }),
      value: formattedFieldValue
    }, _value),
    meta: _extends({}, state, {
      active: !!(state && getIn(state, 'active')),
      asyncValidating: asyncValidating,
      autofilled: !!(state && getIn(state, 'autofilled')),
      dirty: dirty,
      dispatch: dispatch,
      error: error,
      warning: warning,
      invalid: !!error,
      pristine: pristine,
      submitting: !!submitting,
      touched: !!(state && getIn(state, 'touched')),
      valid: !error,
      visited: !!(state && getIn(state, 'visited'))
    }),
    custom: _extends({}, custom, props)
  };
};

exports.default = createFieldProps;