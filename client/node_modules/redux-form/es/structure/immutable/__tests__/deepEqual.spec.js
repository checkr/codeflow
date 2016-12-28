import expect from 'expect';
import { fromJS, List } from 'immutable';
import deepEqual from '../deepEqual';

describe('structure.immutable.deepEqual', function () {
  var testBothWays = function testBothWays(a, b, expectation) {
    expect(deepEqual(a, b)).toBe(expectation);
    expect(deepEqual(b, a)).toBe(expectation);
  };

  it('should work with nested Immutable.Maps', function () {
    testBothWays(fromJS({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }), fromJS({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }), true);
    testBothWays(fromJS({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }), fromJS({
      a: {
        b: {
          c: 42
        },
        d: 2,
        e: 3
      },
      f: 4
    }), false);
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

  it('should work with plain objects inside Immutable.Maps', function () {
    testBothWays(fromJS({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }).setIn('a.b.g', { h: { i: 29 } }), fromJS({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }).setIn('a.b.g', { h: { i: 29 } }), true);
    testBothWays(fromJS({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }).setIn('a.b.g', { h: { i: 29 } }), fromJS({
      a: {
        b: {
          c: 1
        },
        d: 2,
        e: 3
      },
      f: 4
    }).setIn('a.b.g', { h: { i: 30 } }), false);
  });

  it('should work with Immutable.Maps inside plain objects', function () {
    testBothWays({
      a: {
        b: {
          c: fromJS({
            h: {
              i: 29
            }
          })
        },
        d: 2,
        e: 3
      },
      f: 4
    }, {
      a: {
        b: {
          c: fromJS({
            h: {
              i: 29
            }
          })
        },
        d: 2,
        e: 3
      },
      f: 4
    }, true);
    testBothWays({
      a: {
        b: {
          c: fromJS({
            h: {
              i: 29
            }
          })
        },
        d: 2,
        e: 3
      },
      f: 4
    }, {
      a: {
        b: {
          c: fromJS({
            h: {
              i: 30
            }
          })
        },
        d: 2,
        e: 3
      },
      f: 4
    }, false);
  });

  it('should work with Immutable.Lists', function () {
    var firstObj = { a: 1 };
    var secondObj = { a: 1 };
    var thirdObj = { c: 1 };

    testBothWays(List(['a', 'b']), List(['a', 'b', 'c']), false);
    testBothWays(List(['a', 'b', 'c']), List(['a', 'b', 'c']), true);
    testBothWays(List(['a', 'b', firstObj]), List(['a', 'b', secondObj]), true);
    testBothWays(List(['a', 'b', firstObj]), List(['a', 'b', thirdObj]), false);
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
    testBothWays(fromJS({
      a: {
        b: ''
      }
    }), fromJS({
      a: {
        b: undefined
      }
    }), true);
  });

  it('should treat null and \'\' as equal', function () {
    testBothWays(fromJS({
      a: {
        b: ''
      }
    }), fromJS({
      a: {
        b: null
      }
    }), true);
  });

  it('should treat null and undefined as equal', function () {
    testBothWays(fromJS({
      a: {
        b: undefined
      }
    }), fromJS({
      a: {
        b: null
      }
    }), true);
  });

  it('should treat false and undefined as equal', function () {
    testBothWays(fromJS({
      a: {
        b: false
      }
    }), fromJS({
      a: {
        b: undefined
      }
    }), true);
  });
});