import _noop from 'lodash-es/noop';
import expect, { createSpy } from 'expect';
import isPromise from 'is-promise';

import asyncValidation from '../asyncValidation';

describe('asyncValidation', function () {
  var field = 'myField';

  it('should throw an error if fn does not return a promise', function () {
    var fn = _noop;
    var start = _noop;
    var stop = _noop;
    expect(function () {
      return asyncValidation(fn, start, stop, field);
    }).toThrow(/promise/);
  });

  it('should return a promise', function () {
    var fn = function fn() {
      return Promise.resolve();
    };
    var start = _noop;
    var stop = _noop;
    expect(isPromise(asyncValidation(fn, start, stop, field))).toBe(true);
  });

  it('should call start, fn, and stop on promise resolve', function () {
    var fn = createSpy().andReturn(Promise.resolve());
    var start = createSpy();
    var stop = createSpy();
    var promise = asyncValidation(fn, start, stop, field);
    expect(fn).toHaveBeenCalled();
    expect(start).toHaveBeenCalled().toHaveBeenCalledWith(field);
    return promise.then(function () {
      expect(stop).toHaveBeenCalled();
    });
  });

  it('should throw when promise rejected with no errors', function () {
    var fn = createSpy().andReturn(Promise.reject());
    var start = createSpy();
    var stop = createSpy();
    var promise = asyncValidation(fn, start, stop, field);
    expect(fn).toHaveBeenCalled();
    expect(start).toHaveBeenCalled().toHaveBeenCalledWith(field);
    return promise.catch(function () {
      expect(stop).toHaveBeenCalled();
    });
  });

  it('should call start, fn, and stop on promise reject', function () {
    var errors = { foo: 'error' };
    var fn = createSpy().andReturn(Promise.reject(errors));
    var start = createSpy();
    var stop = createSpy();
    var promise = asyncValidation(fn, start, stop, field);
    expect(fn).toHaveBeenCalled();
    expect(start).toHaveBeenCalled().toHaveBeenCalledWith(field);
    return promise.catch(function () {
      expect(stop).toHaveBeenCalled().toHaveBeenCalledWith(errors);
    });
  });
});