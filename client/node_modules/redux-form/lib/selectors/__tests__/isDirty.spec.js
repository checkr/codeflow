'use strict';

var _isDirty = require('../isDirty');

var _isDirty2 = _interopRequireDefault(_isDirty);

var _plain = require('../../structure/plain');

var _plain2 = _interopRequireDefault(_plain);

var _expectations = require('../../structure/plain/expectations');

var _expectations2 = _interopRequireDefault(_expectations);

var _immutable = require('../../structure/immutable');

var _immutable2 = _interopRequireDefault(_immutable);

var _expectations3 = require('../../structure/immutable/expectations');

var _expectations4 = _interopRequireDefault(_expectations3);

var _addExpectations = require('../../__tests__/addExpectations');

var _addExpectations2 = _interopRequireDefault(_addExpectations);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var describeIsDirty = function describeIsDirty(name, structure, expect) {
  var isDirty = (0, _isDirty2.default)(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(isDirty('foo')).toBeA('function');
    });

    it('should return false when values not present', function () {
      expect(isDirty('foo')(fromJS({
        form: {}
      }))).toBe(false);
    });

    it('should return false when values are pristine', function () {
      expect(isDirty('foo')(fromJS({
        form: {
          foo: {
            initial: {
              dog: 'Snoopy',
              cat: 'Garfield'
            },
            values: {
              dog: 'Snoopy',
              cat: 'Garfield'
            }
          }
        }
      }))).toBe(false);
    });

    it('should return true when values are dirty', function () {
      expect(isDirty('foo')(fromJS({
        form: {
          foo: {
            initial: {
              dog: 'Snoopy',
              cat: 'Garfield'
            },
            values: {
              dog: 'Odie',
              cat: 'Garfield'
            }
          }
        }
      }))).toBe(true);
    });

    it('should use getFormState if provided', function () {
      expect(isDirty('foo', function (state) {
        return getIn(state, 'someOtherSlice');
      })(fromJS({
        someOtherSlice: {
          foo: {
            initial: {
              dog: 'Snoopy',
              cat: 'Garfield'
            },
            values: {
              dog: 'Odie',
              cat: 'Garfield'
            }
          }
        }
      }))).toBe(true);
    });
  });
};

describeIsDirty('isDirty.plain', _plain2.default, (0, _addExpectations2.default)(_expectations2.default));
describeIsDirty('isDirty.immutable', _immutable2.default, (0, _addExpectations2.default)(_expectations4.default));