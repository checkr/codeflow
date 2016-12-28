'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeTouch = function describeTouch(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should mark fields as touched on touch', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'value',
            myOtherField: 'otherValue'
          }
        }
      }), (0, _actions.touch)('foo', 'myField', 'myOtherField'));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
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
      });
    });

    it('should mark deep fields as touched on touch', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            deep: {
              myField: 'value',
              myOtherField: 'otherValue'
            }
          }
        }
      }), (0, _actions.touch)('foo', 'deep.myField', 'deep.myOtherField'));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
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
      });
    });

    it('should mark array fields as touched on touch', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myFields: ['value', 'otherValue']
          }
        }
      }), (0, _actions.touch)('foo', 'myFields[0]', 'myFields[1]'));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myFields: ['value', 'otherValue']
          },
          fields: {
            myFields: [{
              touched: true
            }, {
              touched: true
            }]
          }
        }
      });
    });
  };
};

exports.default = describeTouch;