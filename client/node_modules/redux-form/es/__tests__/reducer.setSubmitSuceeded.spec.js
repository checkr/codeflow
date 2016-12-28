import { setSubmitSucceeded } from '../actions';

var describeSetSubmitSucceeded = function describeSetSubmitSucceeded(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set submitSucceeded flag on submitSucceeded', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'change'
        }
      }), setSubmitSucceeded('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'change',
          submitSucceeded: true
        }
      });
    });

    it('should clear submitting flag on submitSucceeded', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'change',
          submitting: true
        }
      }), setSubmitSucceeded('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'change',
          submitSucceeded: true,
          submitting: true
        }
      });
    });

    it('should clear submitFailed flag on submitSucceeded', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          submitting: true,
          submitFailed: true
        }
      }), setSubmitSucceeded('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          submitSucceeded: true,
          submitting: true
        }
      });
    });
  };
};

export default describeSetSubmitSucceeded;