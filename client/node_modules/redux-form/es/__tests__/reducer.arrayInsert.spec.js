import { arrayInsert } from '../actions';

var describeArrayInsert = function describeArrayInsert(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should work with empty state', function () {
      var state = reducer(undefined, arrayInsert('foo', 'myField', 0, 'myValue'));
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
      }), arrayInsert('foo', 'myField.subField', 0, 'newValue'));
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
      }), arrayInsert('foo', 'myField.subField', 3, 'newValue'));
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
      }), arrayInsert('foo', 'myField.subField', 1, 'newValue'));
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
  };
};

export default describeArrayInsert;