'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeUntouch = function describeUntouch(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should unmark fields as touched on untouch', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          },
          fields: {
            myField: {
              touched: true
            },
            myOtherField: {
              touched: true
            }
          }
        }
      }), (0, _actions.untouch)('foo', 'myField', 'myOtherField'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          },
          fields: {
            myField: {},
            myOtherField: {}
          }
        }
      });
    });

    it('should unmark deep fields as touched on untouch', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            deep: {
              myField: 'value',
              myOtherField: 'otherValue'
            }
          },
          fields: {
            deep: {
              myField: {
                touched: true
              },
              myOtherField: {
                touched: true
              }
            }
          }
        }
      }), (0, _actions.untouch)('foo', 'deep.myField', 'deep.myOtherField'));
      expect(state).toEqualMap({
        foo: {
          values: {
            deep: {
              myField: 'value',
              myOtherField: 'otherValue'
            }
          },
          fields: {
            deep: {
              myField: {},
              myOtherField: {}
            }
          }
        }
      });
    });

    it('should unmark array fields as touched on untouch', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myFields: ['value', 'otherValue']
          },
          fields: {
            myFields: [{ touched: true }, { touched: true }]
          }
        }
      }), (0, _actions.untouch)('foo', 'myFields[0]', 'myFields[1]'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myFields: ['value', 'otherValue']
          },
          fields: {
            myFields: [{}, {}]
          }
        }
      });
    });
  };
};

exports.default = describeUntouch;