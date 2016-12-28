'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeSubmit = function describeSubmit(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set triggerSubmit with no previous state', function () {
      var state = reducer(undefined, (0, _actions.submit)('foo'));
      expect(state).toEqualMap({
        foo: {
          triggerSubmit: true
        }
      });
    });

    it('should set triggerSubmit with previous state', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {},
            otherField: {
              visited: true,
              active: true
            }
          },
          active: 'otherField'
        }
      }), (0, _actions.submit)('foo'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {},
            otherField: {
              visited: true,
              active: true
            }
          },
          active: 'otherField',
          triggerSubmit: true
        }
      });
    });
  };
};

exports.default = describeSubmit;