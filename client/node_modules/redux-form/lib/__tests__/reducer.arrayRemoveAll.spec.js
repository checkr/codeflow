'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeArrayRemoveAll = function describeArrayRemoveAll(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should do nothing with undefined', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: {}
          },
          fields: {
            myField: {}
          }
        }
      }), (0, _actions.arrayRemoveAll)('foo', 'myField.subField'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {}
          },
          fields: {
            myField: {}
          }
        }
      });
    });

    it('should do nothing if already empty', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: {
              subField: []
            }
          },
          fields: {
            myField: {
              subField: []
            }
          }
        }
      }), (0, _actions.arrayRemoveAll)('foo', 'myField.subField'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: []
            }
          },
          fields: {
            myField: {
              subField: []
            }
          }
        }
      });
    });

    it('should remove all the elements', function () {
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
      }), (0, _actions.arrayRemoveAll)('foo', 'myField.subField', 1));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: []
            }
          },
          fields: {
            myField: {
              subField: []
            }
          }
        }
      });
    });
  };
};

exports.default = describeArrayRemoveAll;