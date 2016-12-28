'use strict';

var _deleteInWithCleanUp = require('../deleteInWithCleanUp');

var _deleteInWithCleanUp2 = _interopRequireDefault(_deleteInWithCleanUp);

var _plain = require('../structure/plain');

var _plain2 = _interopRequireDefault(_plain);

var _expectations = require('../structure/plain/expectations');

var _expectations2 = _interopRequireDefault(_expectations);

var _immutable = require('../structure/immutable');

var _immutable2 = _interopRequireDefault(_immutable);

var _expectations3 = require('../structure/immutable/expectations');

var _expectations4 = _interopRequireDefault(_expectations3);

var _addExpectations = require('./addExpectations');

var _addExpectations2 = _interopRequireDefault(_addExpectations);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var describeDeleteInWithCleanUp = function describeDeleteInWithCleanUp(name, structure, expect) {
  var fromJS = structure.fromJS;

  var deleteInWithCleanUp = (0, _deleteInWithCleanUp2.default)(structure);

  describe(name, function () {
    it('should delete from a flat structure', function () {
      expect(deleteInWithCleanUp(fromJS({
        dog: 'Scooby',
        cat: 'Garfield'
      }), 'dog')).toEqualMap({
        cat: 'Garfield'
      });
    });

    it('should not delete parent if has other children', function () {
      expect(deleteInWithCleanUp(fromJS({
        a: {
          b: 1,
          c: 2
        },
        d: {
          e: 3
        }
      }), 'a.b')).toEqualMap({
        a: {
          c: 2
        },
        d: {
          e: 3
        }
      });
    });

    it('should just set to undefined if leaf structure is an array', function () {
      expect(deleteInWithCleanUp(fromJS({
        a: [42]
      }), 'a[0]')).toEqualMap({
        a: [undefined]
      });
      expect(deleteInWithCleanUp(fromJS({
        a: [42]
      }), 'b[0]')).toEqualMap({
        a: [42]
      });
      expect(deleteInWithCleanUp(fromJS({
        a: [41, 42, 43]
      }), 'a[1]')).toEqualMap({
        a: [41, undefined, 43]
      });
      expect(deleteInWithCleanUp(fromJS({
        a: {
          b: 1,
          c: [2]
        },
        d: {
          e: 3
        }
      }), 'a.c[0]')).toEqualMap({
        a: {
          b: 1,
          c: [undefined]
        },
        d: {
          e: 3
        }
      });
    });

    it('should delete parent if no other children', function () {
      expect(deleteInWithCleanUp(fromJS({
        a: {
          b: 1,
          c: 2
        },
        d: {
          e: 3
        }
      }), 'd.e')).toEqualMap({
        a: {
          b: 1,
          c: 2
        }
      });
      expect(deleteInWithCleanUp(fromJS({
        a: {
          b: {
            c: {
              d: {
                e: {
                  f: 'That\'s DEEP!'
                }
              }
            }
          }
        }
      }), 'a.b.c.d.e.f')).toEqualMap({});
    });
  });
};

describeDeleteInWithCleanUp('deleteInWithCleanUp.plain', _plain2.default, (0, _addExpectations2.default)(_expectations2.default));
describeDeleteInWithCleanUp('deleteInWithCleanUp.immutable', _immutable2.default, (0, _addExpectations2.default)(_expectations4.default));