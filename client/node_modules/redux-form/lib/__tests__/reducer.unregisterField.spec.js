'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeUnregisterField = function describeUnregisterField(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should remove a field from registeredFields', function () {
      var state = reducer(fromJS({
        foo: {
          registeredFields: [{ name: 'bar', type: 'field' }]
        }
      }), (0, _actions.unregisterField)('foo', 'bar'));
      expect(state).toEqualMap({
        foo: {}
      });
    });

    it('should do nothing if there are no registered fields', function () {
      var initialState = fromJS({
        foo: {}
      });
      var state = reducer(initialState, (0, _actions.unregisterField)('foo', 'bar'));
      expect(state).toEqual(initialState);
    });

    it('should do nothing if the field is not registered', function () {
      var state = reducer(fromJS({
        foo: {
          registeredFields: [{ name: 'bar', type: 'field' }]
        }
      }), (0, _actions.unregisterField)('foo', 'baz'));
      expect(state).toEqualMap({
        foo: {
          registeredFields: [{ name: 'bar', type: 'field' }]
        }
      });
    });
  };
};

exports.default = describeUnregisterField;