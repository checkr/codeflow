import _noop from 'lodash-es/noop';
import expect, { createSpy } from 'expect';

import silenceEvents from '../silenceEvents';

describe('silenceEvents', function () {
  it('should return a function', function () {
    expect(silenceEvents()).toExist().toBeA('function');
  });

  it('should return pass all args if first arg is not event', function () {
    var spy = createSpy();
    var silenced = silenceEvents(spy);

    silenced(1, 2, 3);
    expect(spy).toHaveBeenCalled().toHaveBeenCalledWith(1, 2, 3);
    spy.restore();

    silenced('foo', 'bar');
    expect(spy).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'bar');
    spy.restore();

    silenced({ value: 10 }, false);
    expect(spy).toHaveBeenCalled().toHaveBeenCalledWith({ value: 10 }, false);
    spy.restore();
  });

  it('should return pass other args if first arg is event', function () {
    var spy = createSpy();
    var silenced = silenceEvents(spy);
    var event = {
      preventDefault: _noop,
      stopPropagation: _noop
    };

    silenced(event, 1, 2, 3);
    expect(spy).toHaveBeenCalled().toHaveBeenCalledWith(1, 2, 3);
    spy.restore();

    silenced(event, 'foo', 'bar');
    expect(spy).toHaveBeenCalled().toHaveBeenCalledWith('foo', 'bar');
    spy.restore();

    silenced(event, { value: 10 }, false);
    expect(spy).toHaveBeenCalled().toHaveBeenCalledWith({ value: 10 }, false);
    spy.restore();
  });

  it('should silence event', function () {
    var spy = createSpy();
    var preventDefault = createSpy();
    var stopPropagation = createSpy();
    var event = {
      preventDefault: preventDefault,
      stopPropagation: stopPropagation
    };

    silenceEvents(spy)(event);
    expect(preventDefault).toHaveBeenCalled();
    expect(stopPropagation).toNotHaveBeenCalled();
    expect(spy).toHaveBeenCalled();
  });
});