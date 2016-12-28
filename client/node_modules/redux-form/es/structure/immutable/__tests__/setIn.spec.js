import expect from 'expect';
import { fromJS, Map, List } from 'immutable';
import setIn from '../setIn';

describe('structure.immutable.setIn', function () {
  it('should handle undefined', function () {
    var result = setIn(new Map(), undefined, 'success');
    expect(result).toEqual('success');
  });
  it('should handle dot paths', function () {
    var result = setIn(new Map(), 'a.b.c', 'success');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var c = b.get('c');
    expect(c).toEqual('success');
  });
  it('should handle array paths', function () {
    var result = setIn(new Map(), 'a.b[0]', 'success');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing').toBeA(List).toEqual(fromJS(['success']));
  });
  it('should handle array paths with successive sets', function () {
    var result = setIn(new Map(), 'a.b[2]', 'success');
    result = setIn(result, 'a.b[0]', 'success');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing').toBeA(List).toEqual(fromJS(['success', undefined, 'success']));
  });
  it('should handle array paths with existing array', function () {
    var result = setIn(new Map({
      a: new Map({
        b: new List(['first'])
      })
    }), 'a.b[1].value', 'success');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing').toBeA(List).toEqual(fromJS(['first', { value: 'success' }]));
  });
  it('should handle array paths with existing array with undefined', function () {
    var result = setIn(new Map({
      a: new Map({
        b: new List(['first', undefined])
      })
    }), 'a.b[1].value', 'success');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing').toBeA(List).toEqual(fromJS(['first', { value: 'success' }]));
  });
  it('should handle multiple array paths', function () {
    var result = setIn(new Map(), 'a.b[0].c.d[13].e', 'success');

    var a = result.get('a');
    expect(a).toExist('a missing').toBeA(Map);

    var b = a.get('b');
    expect(b).toExist('b missing').toBeA(List);

    var b0 = b.get(0);
    expect(b0).toExist('b[0] missing').toBeA(Map);

    var c = b0.get('c');
    expect(c).toExist('c missing').toBeA(Map);

    var d = c.get('d');
    expect(d).toExist('d missing').toBeA(List);

    var d13 = d.get(13);
    expect(d13).toExist('d[13] missing').toBeA(Map);

    var e = d13.get('e');
    expect(e).toExist('e missing').toEqual('success');
  });
  it('should handle indexer paths', function () {
    var result = setIn(new Map(), 'a.b[c].d[e]', 'success');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var c = b.get('c');
    expect(c).toExist('c missing');

    var d = c.get('d');
    expect(d).toExist('d missing');

    var e = d.get('e');
    expect(e).toExist('e missing').toEqual('success');
  });
  it('should update existing Map', function () {
    var initial = fromJS({
      a: {
        b: { c: 'one' }
      }
    });

    var result = setIn(initial, 'a.b.c', 'two');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var c = b.get('c');
    expect(c).toEqual('two');
  });
  it('should update existing List', function () {
    var initial = fromJS({
      a: {
        b: [{ c: 'one' }]
      }
    });

    var result = setIn(initial, 'a.b[0].c', 'two');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var b0 = b.get(0);
    expect(b0).toExist();

    var b0c = b0.get('c');
    expect(b0c).toEqual('two');
  });
  it('should not break existing Map', function () {
    var initial = fromJS({
      a: {
        b: { c: 'one' }
      }
    });

    var result = setIn(initial, 'a.b.d', 'two');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var c = b.get('c');
    expect(c).toEqual('one');

    var d = b.get('d');
    expect(d).toEqual('two');
  });
  it('should not break existing List', function () {
    var initial = fromJS({
      a: {
        b: [{ c: 'one' }, { c: 'two' }]
      }
    });

    var result = setIn(initial, 'a.b[0].c', 'changed');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var b0 = b.get(0);
    expect(b0).toExist();

    var b0c = b0.get('c');
    expect(b0c).toEqual('changed');

    var b1 = b.get(1);
    expect(b1).toExist();

    var b1c = b1.get('c');
    expect(b1c).toEqual('two');
  });
  it('should add to an existing List', function () {
    var initial = fromJS({
      a: {
        b: [{ c: 'one' }, { c: 'two' }]
      }
    });

    var result = setIn(initial, 'a.b[2].c', 'three');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var b0 = b.get(0);
    expect(b0).toExist();

    var b0c = b0.get('c');
    expect(b0c).toEqual('one');

    var b1 = b.get(1);
    expect(b1).toExist();

    var b1c = b1.get('c');
    expect(b1c).toEqual('two');

    var b2 = b.get(2);
    expect(b2).toExist();

    var b2c = b2.get('c');
    expect(b2c).toEqual('three');
  });
  it('should set a value directly on new list', function () {
    var result = setIn(new Map(), 'a.b[2]', 'three');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var b0 = b.get(0);
    expect(b0).toEqual(undefined);

    var b1 = b.get(1);
    expect(b1).toEqual(undefined);

    var b2 = b.get(2);
    expect(b2).toEqual('three');
  });
  it('should add to an existing List item', function () {
    var initial = fromJS({
      a: {
        b: [{
          c: '123'
        }]
      }
    });

    var result = setIn(initial, 'a.b[0].d', '12');

    var a = result.get('a');
    expect(a).toExist('a missing');

    var b = a.get('b');
    expect(b).toExist('b missing');

    var b0 = b.get(0);
    expect(b0).toExist();

    var b0d = b0.get('d');
    expect(b0d).toEqual('12');

    var b0c = b0.get('c');
    expect(b0c).toEqual('123');
  });
});