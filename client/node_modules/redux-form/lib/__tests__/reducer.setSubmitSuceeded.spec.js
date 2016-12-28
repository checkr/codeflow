'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeSetSubmitSucceeded = function describeSetSubmitSucceeded(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set submitSucceeded flag on submitSucceeded', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'change'
        }
      }), (0, _actions.setSubmitSucceeded)('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'change',
          submitSucceeded: true
        }
      });
    });

    it('should clear submitting flag on submitSucceeded', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'change',
          submitting: true
        }
      }), (0, _actions.setSubmitSucceeded)('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'change',
          submitSucceeded: true,
          submitting: true
        }
      });
    });

    it('should clear submitFailed flag on submitSucceeded', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          submitting: true,
          submitFailed: true
        }
      }), (0, _actions.setSubmitSucceeded)('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          submitSucceeded: true,
          submitting: true
        }
      });
    });
  };
};

exports.default = describeSetSubmitSucceeded;