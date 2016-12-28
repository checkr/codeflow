import expect from 'expect';
import deepEqual from 'deep-equal';
import { Map, List, Iterable, fromJS } from 'immutable';

var deepEqualValues = function deepEqualValues(a, b) {
  if (Iterable.isIterable(a)) {
    return Iterable.isIterable(b) && a.count() === b.count() && a.every(function (value, key) {
      return deepEqualValues(value, b.get(key));
    });
  }
  return deepEqual(a, b); // neither are immutable iterables
};

var api = {
  toBeAMap: function toBeAMap() {
    expect.assert(Map.isMap(this.actual), 'expected %s to be an immutable Map', this.actual);
    return this;
  },
  toBeAList: function toBeAList() {
    expect.assert(List.isList(this.actual), 'expected %s to be an immutable List', this.actual);
    return this;
  },
  toBeSize: function toBeSize(size) {
    expect.assert(Iterable.isIterable(this.actual) && this.actual.count() === size, 'expected %s to contain %s elements', this.actual, size);
    return this;
  },
  toEqualMap: function toEqualMap(expected) {
    expect.assert(deepEqualValues(this.actual, fromJS(expected)), 'expected...\n%s\n...but found...\n%s', fromJS(expected), this.actual);
    return this;
  }
};

export default api;