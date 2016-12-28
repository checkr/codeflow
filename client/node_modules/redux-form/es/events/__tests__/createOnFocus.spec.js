import expect, { createSpy } from 'expect';
import createOnFocus from '../createOnFocus';

describe('createOnFocus', function () {
  it('should return a function', function () {
    expect(createOnFocus()).toExist().toBeA('function');
  });

  it('should return a function that calls focus with name', function () {
    var focus = createSpy();
    createOnFocus('foo', focus)();
    expect(focus).toHaveBeenCalled().toHaveBeenCalledWith('foo');
  });
});