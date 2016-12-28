import _some from 'lodash-es/some';
import splice from './splice';
import getIn from './getIn';
import setIn from './setIn';
import deepEqual from './deepEqual';
import deleteIn from './deleteIn';


var structure = {
  empty: {},
  emptyList: [],
  getIn: getIn,
  setIn: setIn,
  deepEqual: deepEqual,
  deleteIn: deleteIn,
  fromJS: function fromJS(value) {
    return value;
  },
  size: function size(array) {
    return array ? array.length : 0;
  },
  some: _some,
  splice: splice
};

export default structure;