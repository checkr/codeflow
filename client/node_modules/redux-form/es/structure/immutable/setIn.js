import _toPath from 'lodash-es/toPath';
import { List, Map } from 'immutable';


var arrayPattern = /\[(\d+)\]/;

var undefinedArrayMerge = function undefinedArrayMerge(previous, next) {
  return next !== undefined ? next : previous;
};

var mergeLists = function mergeLists(original, value) {
  return original && List.isList(original) ? original.mergeDeepWith(undefinedArrayMerge, value) : value;
};

/*
 * ImmutableJS' setIn function doesn't support array (List) creation
 * so we must pre-insert all arrays in the path ahead of time.
 * 
 * Additionally we must also pre-set a dummy Map at the location
 * of an array index if there's parts that come afterwards because 
 * the setIn function uses `{}` to mark an unset value instead of 
 * undefined (which is the case for list / arrays).
 */
export default function setIn(state, field, value) {
  if (!field || typeof field !== 'string' || !arrayPattern.test(field)) {
    return state.setIn(_toPath(field), value);
  }

  return state.withMutations(function (mutable) {
    var arraySafePath = field.split('.');
    var pathSoFar = null;

    var _loop = function _loop(partIndex) {
      var part = arraySafePath[partIndex];
      var match = arrayPattern.exec(part);

      pathSoFar = pathSoFar === null ? part : pathSoFar + '.' + part;

      if (!match) return 'continue';

      var arr = [];
      arr[parseInt(match[1])] = partIndex + 1 >= arraySafePath.length ? new Map() : undefined;

      mutable = mutable.updateIn(_toPath(pathSoFar).slice(0, -1), function (value) {
        return mergeLists(value, new List(arr));
      });
    };

    for (var partIndex in arraySafePath) {
      var _ret = _loop(partIndex);

      if (_ret === 'continue') continue;
    }

    return mutable.setIn(_toPath(field), value);
  });
}