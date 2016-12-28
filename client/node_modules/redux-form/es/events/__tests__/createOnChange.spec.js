import expect, { createSpy } from 'expect';
import createOnChange from '../createOnChange';

describe('createOnChange', function () {
  it('should return a function', function () {
    expect(createOnChange()).toExist().toBeA('function');
  });

  it('should parse the value before dispatching action', function () {
    var change = createSpy();
    var parse = createSpy(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    createOnChange(change, { parse: parse })('bar');
    expect(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    expect(change).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
  });

  it('should normalize the value before dispatching action', function () {
    var change = createSpy();
    var normalize = createSpy(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    createOnChange(change, { normalize: normalize })('bar');
    expect(normalize).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    expect(change).toHaveBeenCalled().toHaveBeenCalledWith('normalized-bar');
  });

  it('should parse before normalize', function () {
    var change = createSpy();
    var parse = createSpy(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    var normalize = createSpy(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    createOnChange(change, { normalize: normalize, parse: parse })('bar');
    expect(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    expect(normalize).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
    expect(change).toHaveBeenCalled().toHaveBeenCalledWith('normalized-parsed-bar');
  });
});