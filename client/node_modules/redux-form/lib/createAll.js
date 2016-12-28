'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

var _reducer = require('./reducer');

var _reducer2 = _interopRequireDefault(_reducer);

var _reduxForm = require('./reduxForm');

var _reduxForm2 = _interopRequireDefault(_reduxForm);

var _Field = require('./Field');

var _Field2 = _interopRequireDefault(_Field);

var _Fields = require('./Fields');

var _Fields2 = _interopRequireDefault(_Fields);

var _FieldArray = require('./FieldArray');

var _FieldArray2 = _interopRequireDefault(_FieldArray);

var _formValueSelector = require('./formValueSelector');

var _formValueSelector2 = _interopRequireDefault(_formValueSelector);

var _values = require('./values');

var _values2 = _interopRequireDefault(_values);

var _getFormValues = require('./selectors/getFormValues');

var _getFormValues2 = _interopRequireDefault(_getFormValues);

var _getFormSyncErrors = require('./selectors/getFormSyncErrors');

var _getFormSyncErrors2 = _interopRequireDefault(_getFormSyncErrors);

var _getFormSubmitErrors = require('./selectors/getFormSubmitErrors');

var _getFormSubmitErrors2 = _interopRequireDefault(_getFormSubmitErrors);

var _isDirty = require('./selectors/isDirty');

var _isDirty2 = _interopRequireDefault(_isDirty);

var _isInvalid = require('./selectors/isInvalid');

var _isInvalid2 = _interopRequireDefault(_isInvalid);

var _isPristine = require('./selectors/isPristine');

var _isPristine2 = _interopRequireDefault(_isPristine);

var _isValid = require('./selectors/isValid');

var _isValid2 = _interopRequireDefault(_isValid);

var _FormSection = require('./FormSection');

var _FormSection2 = _interopRequireDefault(_FormSection);

var _SubmissionError = require('./SubmissionError');

var _SubmissionError2 = _interopRequireDefault(_SubmissionError);

var _propTypes = require('./propTypes');

var _propTypes2 = _interopRequireDefault(_propTypes);

var _actions = require('./actions');

var actions = _interopRequireWildcard(_actions);

var _actionTypes = require('./actionTypes');

var actionTypes = _interopRequireWildcard(_actionTypes);

function _interopRequireWildcard(obj) { if (obj && obj.__esModule) { return obj; } else { var newObj = {}; if (obj != null) { for (var key in obj) { if (Object.prototype.hasOwnProperty.call(obj, key)) newObj[key] = obj[key]; } } newObj.default = obj; return newObj; } }

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var createAll = function createAll(structure) {
  return _extends({
    // separate out field actions
    actionTypes: actionTypes
  }, actions, {
    Field: (0, _Field2.default)(structure),
    Fields: (0, _Fields2.default)(structure),
    FieldArray: (0, _FieldArray2.default)(structure),
    FormSection: _FormSection2.default,
    formValueSelector: (0, _formValueSelector2.default)(structure),
    getFormValues: (0, _getFormValues2.default)(structure),
    getFormSyncErrors: (0, _getFormSyncErrors2.default)(structure),
    getFormSubmitErrors: (0, _getFormSubmitErrors2.default)(structure),
    isDirty: (0, _isDirty2.default)(structure),
    isInvalid: (0, _isInvalid2.default)(structure),
    isPristine: (0, _isPristine2.default)(structure),
    isValid: (0, _isValid2.default)(structure),
    propTypes: _propTypes2.default,
    reduxForm: (0, _reduxForm2.default)(structure),
    reducer: (0, _reducer2.default)(structure),
    SubmissionError: _SubmissionError2.default,
    values: (0, _values2.default)(structure)
  });
};

exports.default = createAll;