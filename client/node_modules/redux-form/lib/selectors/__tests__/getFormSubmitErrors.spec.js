'use strict';

var _getFormSubmitErrors = require('../getFormSubmitErrors');

var _getFormSubmitErrors2 = _interopRequireDefault(_getFormSubmitErrors);

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

var describeGetFormSubmitErrors = function describeGetFormSubmitErrors(name, structure, expect) {
  var getFormSubmitErrors = (0, _getFormSubmitErrors2.default)(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect((0, _getFormSubmitErrors2.default)('foo')).toBeA('function');
    });

    it('should get the form values from state', function () {
      expect(getFormSubmitErrors('foo')(fromJS({
        form: {
          foo: {
            submitErrors: {
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

    it('should return undefined if there are no submitErrors', function () {
      expect(getFormSubmitErrors('foo')(fromJS({
        form: {
          foo: {}
        }
      }))).toEqual(undefined);
    });

    it('should use getFormState if provided', function () {
      expect(getFormSubmitErrors('foo', function (state) {
        return getIn(state, 'someOtherSlice');
      })(fromJS({
        someOtherSlice: {
          foo: {
            submitErrors: {
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

describeGetFormSubmitErrors('getFormSubmitErrors.plain', _plain2.default, (0, _addExpectations2.default)(_expectations2.default));
describeGetFormSubmitErrors('getFormSubmitErrors.immutable', _immutable2.default, (0, _addExpectations2.default)(_expectations4.default));