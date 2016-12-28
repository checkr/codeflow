import { submit } from '../actions';

var describeSubmit = function describeSubmit(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set triggerSubmit with no previous state', function () {
      var state = reducer(undefined, submit('foo'));
      expect(state).toEqualMap({
        foo: {
          triggerSubmit: true
        }
      });
    });

    it('should set triggerSubmit with previous state', function () {
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
          active: 'otherField'
        }
      }), submit('foo'));
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
          active: 'otherField',
          triggerSubmit: true
        }
      });
    });
  };
};

export default describeSubmit;