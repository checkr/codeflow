import { startSubmit } from '../actions';

var describeStartSubmit = function describeStartSubmit(reducer, expect, _ref) {
  var fromJS = _ref.fromJS;
  return function () {
    it('should set submitting on startSubmit', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'notchange'
        }
      }), startSubmit('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          submitting: true
        }
      });
    });

    it('should set submitting on startSubmit, and NOT reset submitFailed', function () {
      var state = reducer(fromJS({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          submitFailed: true
        }
      }), startSubmit('foo'));
      expect(state).toEqualMap({
        foo: {
          doesnt: 'matter',
          should: 'notchange',
          submitting: true,
          submitFailed: true
        }
      });
    });
  };
};

export default describeStartSubmit;