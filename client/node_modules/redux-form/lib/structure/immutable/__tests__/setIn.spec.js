'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _immutable = require('immutable');

var _setIn = require('../setIn');

var _setIn2 = _interopRequireDefault(_setIn);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('structure.immutable.setIn', function () {
  it('should handle undefined', function () {
    var result = (0, _setIn2.default)(new _immutable.Map(), undefined, 'success');
    (0, _expect2.default)(result).toEqual('success');
  });
  it('should handle dot paths', function () {
    var result = (0, _setIn2.default)(new _immutable.Map(), 'a.b.c', 'success');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var c = b.get('c');
    (0, _expect2.default)(c).toEqual('success');
  });
  it('should handle array paths', function () {
    var result = (0, _setIn2.default)(new _immutable.Map(), 'a.b[0]', 'success');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing').toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['success']));
  });
  it('should handle array paths with successive sets', function () {
    var result = (0, _setIn2.default)(new _immutable.Map(), 'a.b[2]', 'success');
    result = (0, _setIn2.default)(result, 'a.b[0]', 'success');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing').toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['success', undefined, 'success']));
  });
  it('should handle array paths with existing array', function () {
    var result = (0, _setIn2.default)(new _immutable.Map({
      a: new _immutable.Map({
        b: new _immutable.List(['first'])
      })
    }), 'a.b[1].value', 'success');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing').toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['first', { value: 'success' }]));
  });
  it('should handle array paths with existing array with undefined', function () {
    var result = (0, _setIn2.default)(new _immutable.Map({
      a: new _immutable.Map({
        b: new _immutable.List(['first', undefined])
      })
    }), 'a.b[1].value', 'success');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing').toBeA(_immutable.List).toEqual((0, _immutable.fromJS)(['first', { value: 'success' }]));
  });
  it('should handle multiple array paths', function () {
    var result = (0, _setIn2.default)(new _immutable.Map(), 'a.b[0].c.d[13].e', 'success');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing').toBeA(_immutable.Map);

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing').toBeA(_immutable.List);

    var b0 = b.get(0);
    (0, _expect2.default)(b0).toExist('b[0] missing').toBeA(_immutable.Map);

    var c = b0.get('c');
    (0, _expect2.default)(c).toExist('c missing').toBeA(_immutable.Map);

    var d = c.get('d');
    (0, _expect2.default)(d).toExist('d missing').toBeA(_immutable.List);

    var d13 = d.get(13);
    (0, _expect2.default)(d13).toExist('d[13] missing').toBeA(_immutable.Map);

    var e = d13.get('e');
    (0, _expect2.default)(e).toExist('e missing').toEqual('success');
  });
  it('should handle indexer paths', function () {
    var result = (0, _setIn2.default)(new _immutable.Map(), 'a.b[c].d[e]', 'success');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var c = b.get('c');
    (0, _expect2.default)(c).toExist('c missing');

    var d = c.get('d');
    (0, _expect2.default)(d).toExist('d missing');

    var e = d.get('e');
    (0, _expect2.default)(e).toExist('e missing').toEqual('success');
  });
  it('should update existing Map', function () {
    var initial = (0, _immutable.fromJS)({
      a: {
        b: { c: 'one' }
      }
    });

    var result = (0, _setIn2.default)(initial, 'a.b.c', 'two');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var c = b.get('c');
    (0, _expect2.default)(c).toEqual('two');
  });
  it('should update existing List', function () {
    var initial = (0, _immutable.fromJS)({
      a: {
        b: [{ c: 'one' }]
      }
    });

    var result = (0, _setIn2.default)(initial, 'a.b[0].c', 'two');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var b0 = b.get(0);
    (0, _expect2.default)(b0).toExist();

    var b0c = b0.get('c');
    (0, _expect2.default)(b0c).toEqual('two');
  });
  it('should not break existing Map', function () {
    var initial = (0, _immutable.fromJS)({
      a: {
        b: { c: 'one' }
      }
    });

    var result = (0, _setIn2.default)(initial, 'a.b.d', 'two');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var c = b.get('c');
    (0, _expect2.default)(c).toEqual('one');

    var d = b.get('d');
    (0, _expect2.default)(d).toEqual('two');
  });
  it('should not break existing List', function () {
    var initial = (0, _immutable.fromJS)({
      a: {
        b: [{ c: 'one' }, { c: 'two' }]
      }
    });

    var result = (0, _setIn2.default)(initial, 'a.b[0].c', 'changed');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var b0 = b.get(0);
    (0, _expect2.default)(b0).toExist();

    var b0c = b0.get('c');
    (0, _expect2.default)(b0c).toEqual('changed');

    var b1 = b.get(1);
    (0, _expect2.default)(b1).toExist();

    var b1c = b1.get('c');
    (0, _expect2.default)(b1c).toEqual('two');
  });
  it('should add to an existing List', function () {
    var initial = (0, _immutable.fromJS)({
      a: {
        b: [{ c: 'one' }, { c: 'two' }]
      }
    });

    var result = (0, _setIn2.default)(initial, 'a.b[2].c', 'three');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var b0 = b.get(0);
    (0, _expect2.default)(b0).toExist();

    var b0c = b0.get('c');
    (0, _expect2.default)(b0c).toEqual('one');

    var b1 = b.get(1);
    (0, _expect2.default)(b1).toExist();

    var b1c = b1.get('c');
    (0, _expect2.default)(b1c).toEqual('two');

    var b2 = b.get(2);
    (0, _expect2.default)(b2).toExist();

    var b2c = b2.get('c');
    (0, _expect2.default)(b2c).toEqual('three');
  });
  it('should set a value directly on new list', function () {
    var result = (0, _setIn2.default)(new _immutable.Map(), 'a.b[2]', 'three');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var b0 = b.get(0);
    (0, _expect2.default)(b0).toEqual(undefined);

    var b1 = b.get(1);
    (0, _expect2.default)(b1).toEqual(undefined);

    var b2 = b.get(2);
    (0, _expect2.default)(b2).toEqual('three');
  });
  it('should add to an existing List item', function () {
    var initial = (0, _immutable.fromJS)({
      a: {
        b: [{
          c: '123'
        }]
      }
    });

    var result = (0, _setIn2.default)(initial, 'a.b[0].d', '12');

    var a = result.get('a');
    (0, _expect2.default)(a).toExist('a missing');

    var b = a.get('b');
    (0, _expect2.default)(b).toExist('b missing');

    var b0 = b.get(0);
    (0, _expect2.default)(b0).toExist();

    var b0d = b0.get('d');
    (0, _expect2.default)(b0d).toEqual('12');

    var b0c = b0.get('c');
    (0, _expect2.default)(b0c).toEqual('123');
  });
});