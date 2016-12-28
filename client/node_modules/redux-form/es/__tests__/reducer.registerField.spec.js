import { registerField } from '../actions';

var describeRegisterField = function describeRegisterField(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should create registeredFields if it does not exist and a field', function () {
      var state = reducer(fromJS({
        foo: {}
      }), registerField('foo', 'bar', 'Field'));
      expect(state).toEqualMap({
        foo: {
          registeredFields: [{ name: 'bar', type: 'Field' }]
        }
      });
    });

    it('should add a field to registeredFields', function () {
      var state = reducer(fromJS({
        foo: {
          registeredFields: [{ name: 'baz', type: 'FieldArray' }]
        }
      }), registerField('foo', 'bar', 'Field'));
      expect(state).toEqualMap({
        foo: {
          registeredFields: [{ name: 'baz', type: 'FieldArray' }, { name: 'bar', type: 'Field' }]
        }
      });
    });

    it('should do nothing if the field already exists', function () {
      var initialState = fromJS({
        foo: {
          registeredFields: [{ name: 'bar', type: 'FieldArray' }]
        }
      });
      var state = reducer(initialState, registerField('foo', 'bar', 'Field'));
      expect(state).toEqual(initialState);
    });
  };
};

export default describeRegisterField;