'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeDestroy = function describeDestroy(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should destroy form state', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          active: 'myField'
        },
        otherThing: {
          touchThis: false
        }
      }), (0, _actions.destroy)('foo'));
      expect(state).toEqualMap({
        otherThing: {
          touchThis: false
        }
      });
    });
  };
};

exports.default = describeDestroy;