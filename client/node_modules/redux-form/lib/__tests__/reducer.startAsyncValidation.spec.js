'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeStartAsyncValidation = function describeStartAsyncValidation(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set asyncValidating on startAsyncValidation', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'notchange'
        }
      }), (0, _actions.startAsyncValidation)('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          asyncValidating: true
        }
      });
    });

    it('should set asyncValidating with field name on startAsyncValidation', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          }
        }
      }), (0, _actions.startAsyncValidation)('foo', 'myField'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          },
          asyncValidating: 'myField'
        }
      });
    });
  };
};

exports.default = describeStartAsyncValidation;