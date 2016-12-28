import expect from 'expect';
import deepEqual from '../deepEqual';

describe('structure.plain.deepEqual', function () {
  var testBothWays = function testBothWays(a, b, expectation) {
    expect(deepEqual(a, b)).toBe(expectation);
    expect(deepEqual(b, a)).toBe(expectation);
  };

  it('should work with nested objects', function () {
    testBothWays({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }, {
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }, true);
    testBothWays({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }, {
      a: {
        b: {
          c: 42
        },
        d: 2,
        e: 3
      },
      f: 4
    }, false);
  });

  it('should work with plain objects', function () {
    testBothWays({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }, {
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }, true);
    testBothWays({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }, {
      a: {
        b: {
          c: 42
        },
        d: 2,
        e: 3
      },
      f: 4
    }, false);
  });

  it('should work with arrays', function () {
    var firstObj = { a: 1 };
    var secondObj = { a: 1 };
    var thirdObj = { c: 1 };

    testBothWays(['a', 'b'], ['a', 'b', 'c'], false);
    testBothWays(['a', 'b', 'c'], ['a', 'b', 'c'], true);
    testBothWays(['a', 'b', firstObj], ['a', 'b', secondObj], true);
    testBothWays(['a', 'b', firstObj], ['a', 'b', thirdObj], false);
  });

  it('should work with plain objects with cycles', function () {
    // Set up cyclical structures:
    //
    // base1, base2 {
    //   a: 1,
    //   deep: {
    //     b: 2,
    //     base: {
    //       a: 1,
    //       deep: { ... }
    //     }
    //   }
    // }

    var base1 = { a: 1 };
    var deep1 = { b: 2, base: base1 };
    base1.deep = deep1;

    var base2 = { a: 1 };
    var deep2 = { b: 2, base: base2 };
    base2.deep = deep2;

    testBothWays(base1, base2, true);
  });

  it('should treat undefined and \'\' as equal', function () {
    testBothWays({
      a: {
        b: ''
      }
    }, {
      a: {
        b: undefined
      }
    }, true);
  });

  it('should treat null and \'\' as equal', function () {
    testBothWays({
      a: {
        b: ''
      }
    }, {
      a: {
        b: null
      }
    }, true);
  });

  it('should treat null and undefined as equal', function () {
    testBothWays({
      a: {
        b: undefined
      }
    }, {
      a: {
        b: null
      }
    }, true);
  });

  it('should special case _error key for arrays', function () {
    var a = ['a', 'b'];
    var b = ['a', 'b'];
    b._error = 'something';
    var c = ['a', 'b'];
    c._error = 'something';

    testBothWays(a, b, false);
    testBothWays(b, c, true);
  });

  it('should treat false and undefined as equal', function () {
    testBothWays({
      a: {
        b: false
      }
    }, {
      a: {
        b: undefined
      }
    }, true);
  });
});