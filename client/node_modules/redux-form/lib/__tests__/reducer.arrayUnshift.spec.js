'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeArrayUnshift = function describeArrayUnshift(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should work with empty state', function () {
      var state = reducer(undefined, (0, _actions.arrayUnshift)('foo', 'myField', 'myValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: ['myValue']
          }
        }
      });
    });

    it('should insert at beginning', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: {
              subField: ['a', 'b', 'c']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true }, { touched: true, visited: true }, { touched: true }]
            }
          }
        }
      }), (0, _actions.arrayUnshift)('foo', 'myField.subField', 'newValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['newValue', 'a', 'b', 'c']
            }
          },
          fields: {
            myField: {
              subField: [{}, { touched: true }, { touched: true, visited: true }, { touched: true }]
            }
          }
        }
      });
    });
  };
};

exports.default = describeArrayUnshift;