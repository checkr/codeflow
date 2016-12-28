'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeArrayMove = function describeArrayMove(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should do nothing with empty state', function () {
      var state = reducer(undefined, (0, _actions.arrayMove)('foo', 'myField', 0, 1));
      expect(state).toEqualMap({
        foo: {}
      });
    });

    it('should move values and field state', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: {
              subField: ['a', 'b', 'c']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true }, { touched: true, visited: false }, { touched: true, visited: true }]
            }
          },
          submitErrors: {
            myField: {
              subField: ['Invalid']
            }
          }
        }
      }), (0, _actions.arrayMove)('foo', 'myField.subField', 0, 2));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['b', 'c', 'a']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true, visited: false }, { touched: true, visited: true }, { touched: true }]
            }
          },
          submitErrors: {
            myField: {
              subField: [,, 'Invalid'] // eslint-disable-line no-sparse-arrays
            }
          }
        }
      });
    });

    it('should move to overflow indexes', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: {
              subField: ['a', 'b', 'c']
            }
          }
        }
      }), (0, _actions.arrayMove)('foo', 'myField.subField', 0, 3));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['b', 'c',, 'a'] // eslint-disable-line no-sparse-arrays
            }
          }
        }
      });
    });
  };
};

exports.default = describeArrayMove;