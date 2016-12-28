'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeReset = function describeReset(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should reset values on reset on with previous state', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'dirtyValue',
            myOtherField: 'otherDirtyValue'
          },
          initial: {
            myField: 'initialValue',
            myOtherField: 'otherInitialValue'
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
      }), (0, _actions.reset)('foo'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'initialValue',
            myOtherField: 'otherInitialValue'
          },
          initial: {
            myField: 'initialValue',
            myOtherField: 'otherInitialValue'
          }
        }
      });
    });

    it('should reset deep values on reset on with previous state', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            deepField: {
              myField: 'dirtyValue',
              myOtherField: 'otherDirtyValue'
            }
          },
          initial: {
            deepField: {
              myField: 'initialValue',
              myOtherField: 'otherInitialValue'
            }
          },
          fields: {
            deepField: {
              myField: {
                touched: true
              },
              myOtherField: {
                touched: true
              }
            }
          },
          active: 'myField'
        }
      }), (0, _actions.reset)('foo'));
      expect(state).toEqualMap({
        foo: {
          values: {
            deepField: {
              myField: 'initialValue',
              myOtherField: 'otherInitialValue'
            }
          },
          initial: {
            deepField: {
              myField: 'initialValue',
              myOtherField: 'otherInitialValue'
            }
          }
        }
      });
    });

    it('should erase values if reset called with no initial values', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'bar'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      }), (0, _actions.reset)('foo'));
      expect(state).toEqualMap({
        foo: {}
      });
    });

    it('should not destroy registered fields', function () {
      var state = reducer(fromJS({
        foo: {
          registeredFields: [{ name: 'username', type: 'Field' }, { name: 'password', type: 'Field' }],
          values: {
            myField: 'bar'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      }), (0, _actions.reset)('foo'));
      expect(state).toEqualMap({
        foo: {
          registeredFields: [{ name: 'username', type: 'Field' }, { name: 'password', type: 'Field' }]
        }
      });
    });
  };
};

exports.default = describeReset;