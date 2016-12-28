import _toPath from 'lodash-es/toPath';
import { Map, Iterable, List, fromJS as _fromJS } from 'immutable';

import deepEqual from './deepEqual';
import setIn from './setIn';
import splice from './splice';
import plainGetIn from '../plain/getIn';

var structure = {
  empty: Map(),
  emptyList: List(),
  getIn: function getIn(state, field) {
    return Map.isMap(state) || List.isList(state) ? state.getIn(_toPath(field)) : plainGetIn(state, field);
  },
  setIn: setIn,
  deepEqual: deepEqual,
  deleteIn: function deleteIn(state, field) {
    return state.deleteIn(_toPath(field));
  },
  fromJS: function fromJS(jsValue) {
    return _fromJS(jsValue, function (key, value) {
      return Iterable.isIndexed(value) ? value.toList() : value.toMap();
    });
  },
  size: function size(list) {
    return list ? list.size : 0;
  },
  some: function some(iterable, callback) {
    return Iterable.isIterable(iterable) ? iterable.some(callback) : false;
  },
  splice: splice
};

export default structure;