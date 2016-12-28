'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeArrayPop = function describeArrayPop(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should do nothing with no array', function () {
      var state = reducer(fromJS({
        foo: {}
      }), (0, _actions.arrayPop)('foo', 'myField.subField'));
      expect(state).toEqualMap({
        foo: {}
      });
    });

    it('should pop from end', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: {
              subField: ['a', 'b', 'c', 'd']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true, visited: true }, { touched: true }, { touched: true, visited: true }, { touched: true }]
            }
          }
        }
      }), (0, _actions.arrayPop)('foo', 'myField.subField'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['a', 'b', 'c']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true, visited: true }, { touched: true }, { touched: true, visited: true }]
            }
          }
        }
      });
    });
  };
};

exports.default = describeArrayPop;