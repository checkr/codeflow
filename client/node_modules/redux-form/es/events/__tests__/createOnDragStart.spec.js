import expect, { createSpy } from 'expect';
import createOnDragStart, { dataKey } from '../createOnDragStart';

describe('createOnDragStart', function () {
  it('should return a function', function () {
    expect(createOnDragStart()).toExist().toBeA('function');
  });

  it('should return a function that calls dataTransfer.setData with key and result from value', function () {
    var setData = createSpy();
    createOnDragStart('foo', 'bar')({
      dataTransfer: { setData: setData }
    });
    expect(setData).toHaveBeenCalled().toHaveBeenCalledWith(dataKey, 'bar');
  });
});