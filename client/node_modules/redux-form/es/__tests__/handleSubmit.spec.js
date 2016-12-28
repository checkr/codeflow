import _noop from 'lodash-es/noop';
import expect, { createSpy } from 'expect';
import handleSubmit from '../handleSubmit';
import SubmissionError from '../SubmissionError';


describe('handleSubmit', function () {
  it('should stop if sync validation fails', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy();
    var props = { startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    handleSubmit(submit, props, false, asyncValidate, ['foo', 'baz']);

    expect(submit).toNotHaveBeenCalled();
    expect(startSubmit).toNotHaveBeenCalled();
    expect(stopSubmit).toNotHaveBeenCalled();
    expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
    expect(asyncValidate).toNotHaveBeenCalled();
    expect(setSubmitSucceeded).toNotHaveBeenCalled();
    expect(setSubmitFailed).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
  });

  it('should stop and return errors if sync validation fails', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var syncErrors = { foo: 'error' };
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy();
    var props = {
      startSubmit: startSubmit,
      stopSubmit: stopSubmit,
      touch: touch,
      setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded,
      syncErrors: syncErrors,
      values: values
    };

    var result = handleSubmit(submit, props, false, asyncValidate, ['foo', 'baz']);

    expect(asyncValidate).toNotHaveBeenCalled();
    expect(submit).toNotHaveBeenCalled();
    expect(startSubmit).toNotHaveBeenCalled();
    expect(stopSubmit).toNotHaveBeenCalled();
    expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
    expect(setSubmitSucceeded).toNotHaveBeenCalled();
    expect(setSubmitFailed).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
    expect(result).toBe(syncErrors);
  });

  it('should return result of sync submit', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = undefined;
    var props = { dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    expect(handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz'])).toBe(69);

    expect(submit).toHaveBeenCalled().toHaveBeenCalledWith(values, dispatch, props);
    expect(startSubmit).toNotHaveBeenCalled();
    expect(stopSubmit).toNotHaveBeenCalled();
    expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
    expect(setSubmitFailed).toNotHaveBeenCalled();
    expect(setSubmitSucceeded).toHaveBeenCalled();
  });

  it('should not submit if async validation fails', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy().andReturn(Promise.resolve(values));
    var props = { dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(function () {
      throw new Error('Expected to fail');
    }).catch(function (result) {
      expect(result).toBe(values);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toNotHaveBeenCalled();
      expect(startSubmit).toNotHaveBeenCalled();
      expect(stopSubmit).toNotHaveBeenCalled();
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitSucceeded).toNotHaveBeenCalled();
      expect(setSubmitFailed).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
    });
  });

  it('should call onSubmitFail with async errors and dispatch if async validation fails and onSubmitFail is defined', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var onSubmitFail = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy().andReturn(Promise.resolve(values));
    var props = { dispatch: dispatch, onSubmitFail: onSubmitFail, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(function () {
      throw new Error('Expected to fail');
    }).catch(function (result) {
      expect(result).toBe(values);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toNotHaveBeenCalled();
      expect(startSubmit).toNotHaveBeenCalled();
      expect(stopSubmit).toNotHaveBeenCalled();
      expect(onSubmitFail).toHaveBeenCalled();
      expect(onSubmitFail.calls[0].arguments[0]).toEqual(values);
      expect(onSubmitFail.calls[0].arguments[1]).toEqual(dispatch);
      expect(onSubmitFail.calls[0].arguments[2]).toBe(null);
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitSucceeded).toNotHaveBeenCalled();
      expect(setSubmitFailed).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
    });
  });

  it('should not submit if async validation fails and return rejected promise', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncErrors = { foo: 'async error' };
    var asyncValidate = createSpy().andReturn(Promise.reject(asyncErrors));
    var props = {
      dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values
    };

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(function () {
      throw new Error('Expected to fail');
    }).catch(function (result) {
      expect(result).toBe(asyncErrors);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toNotHaveBeenCalled();
      expect(startSubmit).toNotHaveBeenCalled();
      expect(stopSubmit).toNotHaveBeenCalled();
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitSucceeded).toNotHaveBeenCalled();
      expect(setSubmitFailed).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
    });
  });

  it('should sync submit if async validation passes', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy().andReturn(Promise.resolve());
    var props = { dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(function (result) {
      expect(result).toBe(69);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toHaveBeenCalled().toHaveBeenCalledWith(values, dispatch, props);
      expect(startSubmit).toNotHaveBeenCalled();
      expect(stopSubmit).toNotHaveBeenCalled();
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitFailed).toNotHaveBeenCalled();
      expect(setSubmitSucceeded).toHaveBeenCalled();
    });
  });

  it('should async submit if async validation passes', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(Promise.resolve(69));
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy().andReturn(Promise.resolve());
    var props = { dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(function (result) {
      expect(result).toBe(69);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toHaveBeenCalled().toHaveBeenCalledWith(values, dispatch, props);
      expect(startSubmit).toHaveBeenCalled();
      expect(stopSubmit).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitFailed).toNotHaveBeenCalled();
      expect(setSubmitSucceeded).toHaveBeenCalled();
    });
  });

  it('should set submit errors if async submit fails', function () {
    var values = { foo: 'bar', baz: 42 };
    var submitErrors = { foo: 'submit error' };
    var submit = createSpy().andReturn(Promise.reject(new SubmissionError(submitErrors)));
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy().andReturn(Promise.resolve());
    var props = { dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(function (error) {
      expect(error).toBe(submitErrors);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toHaveBeenCalled().toHaveBeenCalledWith(values, dispatch, props);
      expect(startSubmit).toHaveBeenCalled();
      expect(stopSubmit).toHaveBeenCalled().toHaveBeenCalledWith(submitErrors);
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitFailed).toHaveBeenCalled();
      expect(setSubmitSucceeded).toNotHaveBeenCalled();
    });
  });

  it('should not set errors if rejected value not a SubmissionError', function () {
    var values = { foo: 'bar', baz: 42 };
    var submitErrors = { foo: 'submit error' };
    var submit = createSpy().andReturn(Promise.reject(submitErrors));
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy().andReturn(Promise.resolve());
    var props = { dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    var resolveSpy = createSpy();
    var errorSpy = createSpy();

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(resolveSpy, errorSpy).then(function () {
      expect(resolveSpy).toNotHaveBeenCalled();
      expect(errorSpy).toHaveBeenCalledWith(submitErrors);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toHaveBeenCalled().toHaveBeenCalledWith(values, dispatch, props);
      expect(startSubmit).toHaveBeenCalled();
      expect(stopSubmit).toHaveBeenCalled().toHaveBeenCalledWith(undefined); // not wrapped in SubmissionError
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitFailed).toHaveBeenCalled();
      expect(setSubmitSucceeded).toNotHaveBeenCalled();
    });
  });

  it('should set submit errors if async submit fails and return rejected promise', function () {
    var values = { foo: 'bar', baz: 42 };
    var submitErrors = { foo: 'submit error' };
    var submit = createSpy().andReturn(Promise.reject(new SubmissionError(submitErrors)));
    var dispatch = _noop;
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy().andReturn(Promise.resolve());
    var props = {
      dispatch: dispatch, startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values
    };

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(function (error) {
      expect(error).toBe(submitErrors);
      expect(asyncValidate).toHaveBeenCalled().toHaveBeenCalledWith();
      expect(submit).toHaveBeenCalled().toHaveBeenCalledWith(values, dispatch, props);
      expect(startSubmit).toHaveBeenCalled();
      expect(stopSubmit).toHaveBeenCalled().toHaveBeenCalledWith(submitErrors);
      expect(touch).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'baz');
      expect(setSubmitFailed).toHaveBeenCalled();
      expect(setSubmitSucceeded).toNotHaveBeenCalled();
    });
  });

  it('should submit when there are old submit errors and persistentSubmitErrors is enabled', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(69);
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy();
    var props = { startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values, persistentSubmitErrors: true };

    handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']);

    expect(submit).toHaveBeenCalled();
  });

  it('should not swallow errors', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andThrow(new Error('spline reticulation failed'));
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy();
    var props = { startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    expect(function () {
      return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']);
    }).toThrow('spline reticulation failed');
    expect(submit).toHaveBeenCalled();
  });

  it('should not swallow async errors', function () {
    var values = { foo: 'bar', baz: 42 };
    var submit = createSpy().andReturn(Promise.reject(new Error('spline reticulation failed')));
    var startSubmit = createSpy();
    var stopSubmit = createSpy();
    var touch = createSpy();
    var setSubmitFailed = createSpy();
    var setSubmitSucceeded = createSpy();
    var asyncValidate = createSpy();
    var props = { startSubmit: startSubmit, stopSubmit: stopSubmit, touch: touch, setSubmitFailed: setSubmitFailed, setSubmitSucceeded: setSubmitSucceeded, values: values };

    var resultSpy = createSpy();
    var errorSpy = createSpy();

    return handleSubmit(submit, props, true, asyncValidate, ['foo', 'baz']).then(resultSpy, errorSpy).then(function () {
      expect(submit).toHaveBeenCalled();
      expect(resultSpy).toNotHaveBeenCalled('promise should not have resolved');
      expect(errorSpy).toHaveBeenCalled('promise should have rejected');
    });
  });
});