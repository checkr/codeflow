import { autofill, change } from '../actions';

var describeBlur = function describeBlur(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set value on autofill with empty state', function () {
      var state = reducer(undefined, autofill('foo', 'myField', 'myValue'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'myValue'
          },
          fields: {
            myField: {
              autofilled: true
            }
          }
        }
      });
    });

    it('should overwrite value on autofill', function () {
      var state = reducer(fromJS({
        foo: {
          anyTouched: true,
          values: {
            myField: 'before'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      }), autofill('foo', 'myField', 'after'));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: 'after'
          },
          fields: {
            myField: {
              touched: true,
              autofilled: true
            }
          }
        }
      });
    });

    it('should set value on change and remove autofilled', function () {
      var state = reducer(fromJS({
        foo: {
          anyTouched: true,
          values: {
            myField: 'autofilled value'
          },
          fields: {
            myField: {
              autofilled: true,
              touched: true
            }
          }
        }
      }), change('foo', 'myField', 'after change', true));
      expect(state).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            myField: 'after change'
          },
          fields: {
            myField: {
              touched: true
            }
          }
        }
      });
    });
  };
};

export default describeBlur;