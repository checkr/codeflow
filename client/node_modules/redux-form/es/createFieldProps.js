import _noop from 'lodash-es/noop';

var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

function _objectWithoutProperties(obj, keys) { var target = {}; for (var i in obj) { if (keys.indexOf(i) >= 0) continue; if (!Object.prototype.hasOwnProperty.call(obj, i)) continue; target[i] = obj[i]; } return target; }

import createOnBlur from './events/createOnBlur';
import createOnChange from './events/createOnChange';
import createOnDragStart from './events/createOnDragStart';
import createOnDrop from './events/createOnDrop';
import createOnFocus from './events/createOnFocus';


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

  var asyncValidate = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : _noop;

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
  var onChange = createOnChange(boundChange, {
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
      onBlur: createOnBlur(function (value) {
        return dispatch(blur(name, value));
      }, {
        normalize: boundNormalize,
        parse: boundParse,
        after: asyncValidate.bind(null, name)
      }),
      onChange: onChange,
      onDragStart: createOnDragStart(name, formattedFieldValue),
      onDrop: createOnDrop(name, boundChange),
      onFocus: createOnFocus(name, function () {
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

export default createFieldProps;