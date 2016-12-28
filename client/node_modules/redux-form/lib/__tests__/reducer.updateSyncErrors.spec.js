'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeUpdateSyncErrors = function describeUpdateSyncErrors(reducer, expect, _ref) {
  var fromJS = _ref.fromJS,
      setIn = _ref.setIn;
  return function () {
    it('should update sync errors', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          }
        }
      }), (0, _actions.updateSyncErrors)('foo', {
        myField: 'myField error',
        myOtherField: 'myOtherField error'
      }));
      expect(state).toEqual(setIn(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          }
        }
      }), 'foo.syncErrors', {
        myField: 'myField error',
        myOtherField: 'myOtherField error'
      }));
    });

    it('should update form-wide error', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          }
        }
      }), (0, _actions.updateSyncErrors)('foo', {
        myField: 'myField error'
      }, 'form wide error'));
      expect(state).toEqualMap(setIn(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          },
          syncError: true,
          error: 'form wide error'
        }
      }), 'foo.syncErrors', {
        myField: 'myField error'
      }));
    });

    it('should update complex sync errors', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          }
        }
      }), (0, _actions.updateSyncErrors)('foo', {
        myField: { complex: true, text: 'myField error' },
        myOtherField: { complex: true, text: 'myOtherField error' }
      }));
      expect(state).toEqual(setIn(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          }
        }
      }), 'foo.syncErrors', {
        myField: { complex: true, text: 'myField error' },
        myOtherField: { complex: true, text: 'myOtherField error' }
      }));
    });

    it('should clear sync errors', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          },
          syncErrors: {
            myField: 'myField error',
            myOtherField: 'myOtherField error'
          }
        }
      }), (0, _actions.updateSyncErrors)('foo', {}));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          }
        }
      });
    });
  };
};

exports.default = describeUpdateSyncErrors;