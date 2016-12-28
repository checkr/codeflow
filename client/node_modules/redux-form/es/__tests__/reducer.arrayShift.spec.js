import { arrayShift } from '../actions';

var describeArrayShift = function describeArrayShift(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
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
      }), arrayShift('foo', 'myField.subField'));
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
  };
};

export default describeArrayShift;