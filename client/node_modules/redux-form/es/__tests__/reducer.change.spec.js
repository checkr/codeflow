import { change } from '../actions';

var describeChange = function describeChange(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set value on change with empty state', function () {
      var state = reducer(undefined, change('foo', 'myField', 'myValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'myValue'
          }
        }
      });
    });

    it('should set value on change and touch with empty state', function () {
      var state = reducer(undefined, change('foo', 'myField', 'myValue', true));
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

    it('should set value on change and touch with initial value', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          }
        }
      }), change('foo', 'myField', 'myValue', true));
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
      }), change('foo', 'myField', '', true));
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

    it('should remove a value if on change is set with \'\'', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          }
        }
      }), change('foo', 'myField', '', true));
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

    it('should NOT remove a value if on change is set with \'\' if it\'s an array field', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: ['initialValue']
          }
        }
      }), change('foo', 'myField[0]', '', true));
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

    it('should remove nested value container if on change clears all values', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            nested: {
              myField: 'initialValue'
            }
          }
        }
      }), change('foo', 'nested.myField', '', true));
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

    it('should not modify a value if called with undefined', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          }
        }
      }), change('foo', 'myField', undefined, true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
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

    it('should set value on change and remove field-level submit and async errors', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initial'
          },
          asyncErrors: {
            myField: 'async error'
          },
          submitErrors: {
            myField: 'submit error'
          },
          error: 'some global error'
        }
      }), change('foo', 'myField', 'different', false));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'different'
          },
          error: 'some global error'
        }
      });
    });

    it('should NOT remove field-level submit errors and global errors if persistentSubmitErrors is enabled', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initial'
          },
          asyncErrors: {
            myField: 'async error' // only this will be removed
          },
          submitErrors: {
            myField: 'submit error'
          },
          error: 'some global error'
        }
      }), change('foo', 'myField', 'different', false, true));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'different'
          },
          submitErrors: {
            myField: 'submit error'
          },
          error: 'some global error'
        }
      });
    });

    it('should set nested value on change with empty state', function () {
      var state = reducer(undefined, change('foo', 'myField.mySubField', 'myValue', false));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              mySubField: 'myValue'
            }
          }
        }
      });
    });
  };
};

export default describeChange;