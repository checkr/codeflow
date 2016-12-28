import _isObject from 'lodash-es/isObject';
import expect from 'expect';


var expectations = {
  toBeAMap: function toBeAMap() {
    expect.assert(_isObject(this.actual), 'expected %s to be an object', this.actual);
    return this;
  },
  toBeAList: function toBeAList() {
    expect.assert(Array.isArray(this.actual), 'expected %s to be an array', this.actual);
    return this;
  },
  toBeSize: function toBeSize(size) {
    expect.assert(this.actual && Object.keys(this.actual).length === size, 'expected %s to contain %s elements', this.actual, size);
    return this;
  },
  toEqualMap: function toEqualMap(expected) {
    return expect(this.actual).toEqual(expected);
  }
};

export default expectations;