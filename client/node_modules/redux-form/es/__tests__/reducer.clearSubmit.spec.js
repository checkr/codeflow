import { clearSubmit } from '../actions';

var describeClearSubmit = function describeClearSubmit(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should do nothing on clear submit with no previous state', function () {
      var state = reducer(undefined, clearSubmit('foo'));
      expect(state).toEqualMap({
        foo: {}
      });
    });

    it('should clear triggerSubmit with previous state', function () {
      var state = reducer(fromJS({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {},
            otherField: {
              visited: true,
              active: true
            }
          },
          active: 'otherField',
          triggerSubmit: true
        }
      }), clearSubmit('foo'));
      expect(state).toEqualMap({
        foo: {
          values: {
            myField: 'initialValue'
          },
          initial: {
            myField: 'initialValue'
          },
          fields: {
            myField: {},
            otherField: {
              visited: true,
              active: true
            }
          },
          active: 'otherField'
        }
      });
    });
  };
};

export default describeClearSubmit;