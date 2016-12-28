'use strict';

var _isPristine = require('../isPristine');

var _isPristine2 = _interopRequireDefault(_isPristine);

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

var describeIsPristine = function describeIsPristine(name, structure, expect) {
  var isPristine = (0, _isPristine2.default)(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(isPristine('foo')).toBeA('function');
    });

    it('should return true when values not present', function () {
      expect(isPristine('foo')(fromJS({
        form: {}
      }))).toBe(true);
    });

    it('should return true when values are pristine', function () {
      expect(isPristine('foo')(fromJS({
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
      }))).toBe(true);
    });

    it('should return true when values are dirty', function () {
      expect(isPristine('foo')(fromJS({
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
      }))).toBe(false);
    });

    it('should use getFormState if provided', function () {
      expect(isPristine('foo', function (state) {
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
      }))).toBe(false);
    });
  });
};

describeIsPristine('isPristine.plain', _plain2.default, (0, _addExpectations2.default)(_expectations2.default));
describeIsPristine('isPristine.immutable', _immutable2.default, (0, _addExpectations2.default)(_expectations4.default));