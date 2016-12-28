import _noop from 'lodash-es/noop';
import expect, { createSpy } from 'expect';

import silenceEvent from '../silenceEvent';

describe('silenceEvent', function () {
  it('should return false if not an event', function () {
    expect(silenceEvent(undefined)).toBe(false);
    expect(silenceEvent(null)).toBe(false);
    expect(silenceEvent(true)).toBe(false);
    expect(silenceEvent(42)).toBe(false);
    expect(silenceEvent({})).toBe(false);
    expect(silenceEvent([])).toBe(false);
    expect(silenceEvent(_noop)).toBe(false);
  });

  it('should return true if an event', function () {
    expect(silenceEvent({
      preventDefault: _noop,
      stopPropagation: _noop
    })).toBe(true);
  });

  it('should call preventDefault and stopPropagation', function () {
    var preventDefault = createSpy();
    var stopPropagation = createSpy();

    silenceEvent({
      preventDefault: preventDefault,
      stopPropagation: stopPropagation
    });
    expect(preventDefault).toHaveBeenCalled();
    expect(stopPropagation).toNotHaveBeenCalled();
  });
});