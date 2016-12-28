'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeBlur = function describeBlur(reducer, expect, _ref) {
  var fromJS = _ref.fromJS,
      setIn = _ref.setIn;
  return function () {
    it('should set value on blur with empty state', function () {
      var state = reducer(undefined, (0, _actions.blur)('foo', 'myField', 'myValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'myValue'
          }
        }
      });
    });

    it('should set value on blur and touch with empty state', function () {
      var state = reducer(undefined, (0, _actions.blur)('foo', 'myField', 'myValue', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: 'myValue'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      });
    });

    it('should set value on blur and touch with initial value', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {
              active: true
            }
          },
          active: 'myField'
        }
      }), (0, _actions.blur)('foo', 'myField', 'myValue', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: 'myValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      });
    });

    it('should not modify value if undefined is passed on blur (for android react native)', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'myValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {
              active: true
            }
          },
          active: 'myField'
        }
      }), (0, _actions.blur)('foo', 'myField', undefined, true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: 'myValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      });
    });

    it('should not modify value if undefined is passed on blur, even if no value existed (for android react native)', function () {
      var state = reducer(fromJS({
        foo: {
          fields: {
            myField: {
              active: true
            }
          },
          active: 'myField'
        }
      }), (0, _actions.blur)('foo', 'myField', undefined, true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          fields: {
            myField: {
              touched: true
            }
          }
        }
      });
    });

    it('should remove a value if on blur is set with \'\'', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          }
        }
      }), (0, _actions.blur)('foo', 'myField', '', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          fields: {
            myField: {
              touched: true
            }
          }
        }
      });
    });

    it('should allow setting an initialized field to \'\'', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          }
        }
      }), (0, _actions.blur)('foo', 'myField', '', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: ''
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      });
    });

    it('should NOT remove a value if on blur is set with \'\' if it\'s an array field', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: ['initialValue']
          }
        }
      }), (0, _actions.blur)('foo', 'myField[0]', '', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: [undefined]
          },
          fields: {
            myField: [{
              touched: true
            }]
          }
        }
      });
    });

    it('should remove nested value container if on blur clears all values', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            nested: {
              myField: 'initialValue'
            }
          }
        }
      }), (0, _actions.blur)('foo', 'nested.myField', '', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          fields: {
            nested: {
              myField: {
                touched: true
              }
            }
          }
        }
      });
    });

    it('should set nested value on blur', function () {
      var state = reducer(fromJS({
        foo: {
          fields: {
            myField: {
              mySubField: {
                active: true
              }
            }
          },
          active: 'myField.mySubField'
        }
      }), (0, _actions.blur)('foo', 'myField.mySubField', 'hello', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: {
              mySubField: 'hello'
            }
          },
          fields: {
            myField: {
              mySubField: {
                touched: true
              }
            }
          }
        }
      });
    });

    it('should set array value on blur', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myArray: []
          },
          fields: {
            myArray: [{ active: true }]
          },
          active: 'myArray[0]'
        }
      }), (0, _actions.blur)('foo', 'myArray[0]', 'hello', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myArray: ['hello']
          },
          fields: {
            myArray: [{
              touched: true
            }]
          }
        }
      });
    });

    it('should set complex value on blur', function () {
      var state = reducer(fromJS({
        foo: {
          fields: {
            myComplexField: {
              active: true
            }
          },
          active: 'myComplexField'
        }
      }), (0, _actions.blur)('foo', 'myComplexField', { id: 42, name: 'Bobby' }, true));
      expect(state).toEqualMap(setIn(fromJS({ // must use setIn to make sure complex value is js object
        foo: {
          anyTouched: true,
          fields: {
            myComplexField: {
              touched: true
            }
          }
        }
      }), 'foo.values.myComplexField', { id: 42, name: 'Bobby' }));
    });

    it('should NOT remove active field if the blurred field is not active', function () {
      var state = reducer(fromJS({
        foo: {
          fields: {
            myField: {
              active: true
            },
            myOtherField: {}
          },
          active: 'myField'
        }
      }), (0, _actions.blur)('foo', 'myOtherField'));
      expect(state).toEqualMap({
        foo: {
          fields: {
            myField: {
              active: true
            },
            myOtherField: {}
          },
          active: 'myField'
        }
      });
    });

    it('should NOT destroy an empty array field object on blur', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myArray: [{}]
          }
        }
      }), (0, _actions.blur)('foo', 'myArray[0].foo', '', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myArray: [{}]
          },
          fields: {
            myArray: [{
              foo: {
                touched: true
              }
            }]
          }
        }
      });
    });
  };
};

exports.default = describeBlur;