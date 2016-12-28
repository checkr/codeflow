'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeRegisterField = function describeRegisterField(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should create registeredFields if it does not exist and a field', function () {
      var state = reducer(fromJS({
        foo: {}
      }), (0, _actions.registerField)('foo', 'bar', 'Field'));
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
      }), (0, _actions.registerField)('foo', 'bar', 'Field'));
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
      var state = reducer(initialState, (0, _actions.registerField)('foo', 'bar', 'Field'));
      expect(state).toEqual(initialState);
    });
  };
};

exports.default = describeRegisterField;