import _isEqualWith from 'lodash-es/isEqualWith';


var customizer = function customizer(obj, other) {
  if (obj == other) return true;
  if ((obj == null || obj === '' || obj === false) && (other == null || other === '' || other === false)) return true;

  if (obj && other && obj._error !== other._error) return false;
  if (obj && other && obj._warning !== other._warning) return false;
};

var deepEqual = function deepEqual(a, b) {
  return _isEqualWith(a, b, customizer);
};

export default deepEqual;