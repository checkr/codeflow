'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeArrayPush = function describeArrayPush(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should work with empty state', function () {
      var state = reducer(undefined, (0, _actions.arrayPush)('foo', 'myField', 'myValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: ['myValue']
          }
        }
      });
    });

    it('should work pushing undefined to empty state', function () {
      var state = reducer(undefined, (0, _actions.arrayPush)('foo', 'myField'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: [undefined]
          }
        }
      });
    });

    it('should work with existing empty array', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: []
          }
        }
      }), (0, _actions.arrayPush)('foo', 'myField', 'myValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: ['myValue']
          }
        }
      });
    });

    it('should push at end', function () {
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
      }), (0, _actions.arrayPush)('foo', 'myField.subField', 'newValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['a', 'b', 'c', 'newValue']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true }, { touched: true, visited: true }, { touched: true }, {}]
            }
          }
        }
      });
    });

    it('should handle a null leaf value', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            steps: null
          }
        }
      }), (0, _actions.arrayPush)('foo', 'steps', 'value'));
      expect(state).toEqualMap({
        foo: {
          values: {
            steps: ['value']
          }
        }
      });
    });
  };
};

exports.default = describeArrayPush;