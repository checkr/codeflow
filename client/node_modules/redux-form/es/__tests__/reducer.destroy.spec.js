import { destroy } from '../actions';

var describeDestroy = function describeDestroy(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should destroy form state', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          active: 'myField'
        },
        otherThing: {
          touchThis: false
        }
      }), destroy('foo'));
      expect(state).toEqualMap({
        otherThing: {
          touchThis: false
        }
      });
    });
  };
};

export default describeDestroy;