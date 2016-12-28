'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _actions = require('../actions');

var describeArraySplice = function describeArraySplice(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should work with empty state', function () {
      var state = reducer(undefined, (0, _actions.arraySplice)('foo', 'myField', 0, 0, 'myValue'));
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
      }), (0, _actions.arraySplice)('foo', 'myField.subField', 0, 0, 'newValue'));
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

    it('should insert at end', function () {
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
      }), (0, _actions.arraySplice)('foo', 'myField.subField', 3, 0, 'newValue'));
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

    it('should insert in middle', function () {
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
      }), (0, _actions.arraySplice)('foo', 'myField.subField', 1, 0, 'newValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['a', 'newValue', 'b', 'c']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true }, {}, { touched: true, visited: true }, { touched: true }]
            }
          }
        }
      });
    });

    it('should remove from beginning', function () {
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
      }), (0, _actions.arraySplice)('foo', 'myField.subField', 0, 1));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['b', 'c', 'd']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true }, { touched: true, visited: true }, { touched: true }]
            }
          }
        }
      });
    });

    it('should remove from end', function () {
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
      }), (0, _actions.arraySplice)('foo', 'myField.subField', 3, 1));
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

    it('should remove from middle', function () {
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
      }), (0, _actions.arraySplice)('foo', 'myField.subField', 1, 1));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: {
              subField: ['a', 'c', 'd']
            }
          },
          fields: {
            myField: {
              subField: [{ touched: true, visited: true }, { touched: true, visited: true }, { touched: true }]
            }
          }
        }
      });
    });
  };
};

exports.default = describeArraySplice;