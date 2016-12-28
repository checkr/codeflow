import expect, { createSpy } from 'expect';
import createOnBlur from '../createOnBlur';

describe('createOnBlur', function () {
  it('should return a function', function () {
    expect(createOnBlur()).toExist().toBeA('function');
  });

  it('should parse the value before dispatching action', function () {
    var blur = createSpy();
    var parse = createSpy(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    createOnBlur(blur, { parse: parse })('bar');
    expect(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    expect(blur).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
  });

  it('should normalize the value before dispatching action', function () {
    var blur = createSpy();
    var normalize = createSpy(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    createOnBlur(blur, { normalize: normalize })('bar');
    expect(normalize).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    expect(blur).toHaveBeenCalled().toHaveBeenCalledWith('normalized-bar');
  });

  it('should parse before normalize', function () {
    var blur = createSpy();
    var parse = createSpy(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    var normalize = createSpy(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    createOnBlur(blur, { normalize: normalize, parse: parse })('bar');
    expect(parse).toHaveBeenCalled().toHaveBeenCalledWith('bar');
    expect(normalize).toHaveBeenCalled().toHaveBeenCalledWith('parsed-bar');
    expect(blur).toHaveBeenCalled().toHaveBeenCalledWith('normalized-parsed-bar');
  });

  it('should call blur then after', function () {
    var blur = createSpy();
    var parse = createSpy(function (value) {
      return 'parsed-' + value;
    }).andCallThrough();
    var normalize = createSpy(function (value) {
      return 'normalized-' + value;
    }).andCallThrough();
    var after = createSpy();
    createOnBlur(blur, { parse: parse, normalize: normalize, after: after })('bar');
    expect(blur).toHaveBeenCalled();
    expect(normalize).toHaveBeenCalled();
    expect(after).toHaveBeenCalled().toHaveBeenCalledWith('normalized-parsed-bar');
  });
});