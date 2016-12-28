'use strict';

var _getFormValues = require('../getFormValues');

var _getFormValues2 = _interopRequireDefault(_getFormValues);

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

var describeGetFormValues = function describeGetFormValues(name, structure, expect) {
  var getFormValues = (0, _getFormValues2.default)(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(getFormValues('foo')).toBeA('function');
    });

    it('should get the form values from state', function () {
      expect(getFormValues('foo')(fromJS({
        form: {
          foo: {
            values: {
              dog: 'Snoopy',
              cat: 'Garfield'
            }
          }
        }
      }))).toEqualMap({
        dog: 'Snoopy',
        cat: 'Garfield'
      });
    });

    it('should use getFormState if provided', function () {
      expect(getFormValues('foo', function (state) {
        return getIn(state, 'someOtherSlice');
      })(fromJS({
        someOtherSlice: {
          foo: {
            values: {
              dog: 'Snoopy',
              cat: 'Garfield'
            }
          }
        }
      }))).toEqualMap({
        dog: 'Snoopy',
        cat: 'Garfield'
      });
    });
  });
};

describeGetFormValues('getFormValues.plain', _plain2.default, (0, _addExpectations2.default)(_expectations2.default));
describeGetFormValues('getFormValues.immutable', _immutable2.default, (0, _addExpectations2.default)(_expectations4.default));