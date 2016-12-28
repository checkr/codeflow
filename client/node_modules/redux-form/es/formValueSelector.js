import invariant from 'invariant';
import plain from './structure/plain';

var createFormValueSelector = function createFormValueSelector(_ref) {
  var getIn = _ref.getIn;
  return function (form) {
    var getFormState = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : function (state) {
      return getIn(state, 'form');
    };

    invariant(form, 'Form value must be specified');
    return function (state) {
      for (var _len = arguments.length, fields = Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
        fields[_key - 1] = arguments[_key];
      }

      invariant(fields.length, 'No fields specified');
      return fields.length === 1 ?
      // only selecting one field, so return its value
      getIn(getFormState(state), form + '.values.' + fields[0]) :
      // selecting many fields, so return an object of field values
      fields.reduce(function (accumulator, field) {
        var value = getIn(getFormState(state), form + '.values.' + field);
        return value === undefined ? accumulator : plain.setIn(accumulator, field, value);
      }, {});
    };
  };
};

export default createFormValueSelector;