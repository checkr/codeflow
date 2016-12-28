import _noop from 'lodash-es/noop';
import expect from 'expect';

import isEvent from '../isEvent';

describe('isEvent', function () {
  it('should return false if event is undefined', function () {
    expect(isEvent()).toBe(false);
  });

  it('should return false if event is null', function () {
    expect(isEvent(null)).toBe(false);
  });

  it('should return false if event is not an object', function () {
    expect(isEvent(42)).toBe(false);
    expect(isEvent(true)).toBe(false);
    expect(isEvent(false)).toBe(false);
    expect(isEvent('not an event')).toBe(false);
  });

  it('should return false if event has no stopPropagation', function () {
    expect(isEvent({
      preventDefault: _noop
    })).toBe(false);
  });

  it('should return false if event has no preventDefault', function () {
    expect(isEvent({
      stopPropagation: _noop
    })).toBe(false);
  });

  it('should return true if event has stopPropagation, and preventDefault', function () {
    expect(isEvent({
      stopPropagation: _noop,
      preventDefault: _noop
    })).toBe(true);
  });
});