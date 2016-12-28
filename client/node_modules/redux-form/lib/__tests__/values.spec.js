'use strict';

var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; }; /* eslint react/no-multi-comp:0 */


var _react = require('react');

var _react2 = _interopRequireDefault(_react);

var _expect = require('expect');

var _reactRedux = require('react-redux');

var _redux = require('redux');

var _reduxImmutablejs = require('redux-immutablejs');

var _reactAddonsTestUtils = require('react-addons-test-utils');

var _reactAddonsTestUtils2 = _interopRequireDefault(_reactAddonsTestUtils);

var _reducer = require('../reducer');

var _reducer2 = _interopRequireDefault(_reducer);

var _values = require('../values');

var _values2 = _interopRequireDefault(_values);

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

var describeValues = function describeValues(name, structure, combineReducers, expect) {
  var values = (0, _values2.default)(structure);
  var reducer = (0, _reducer2.default)(structure);
  var fromJS = structure.fromJS;

  var makeStore = function makeStore(initial) {
    return (0, _redux.createStore)(combineReducers({ form: reducer }), fromJS({ form: initial }));
  };

  var testProps = function testProps(state) {
    var config = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};

    var store = makeStore({ testForm: state });
    var spy = (0, _expect.createSpy)(function () {
      return _react2.default.createElement('div', null);
    }).andCallThrough();

    var Decorated = values(_extends({ form: 'testForm' }, config))(spy);
    _reactAddonsTestUtils2.default.renderIntoDocument(_react2.default.createElement(
      _reactRedux.Provider,
      { store: store },
      _react2.default.createElement(Decorated, null)
    ));
    expect(spy).toHaveBeenCalled();
    return spy.calls[0].arguments[0];
  };

  describe(name, function () {
    it('should get values from Redux state', function () {
      var values = {
        cat: 'rat',
        dog: 'cat'
      };
      var props = testProps({ values: values });
      expect(props.values).toEqualMap(values);
    });

    it('should use values prop', function () {
      var values = {
        cat: 'rat',
        dog: 'cat'
      };
      var props = testProps({ values: values }, { prop: 'foo' });
      expect(props.foo).toEqualMap(values);
    });
  });
};

describeValues('values.plain', _plain2.default, _redux.combineReducers, (0, _addExpectations2.default)(_expectations2.default));
describeValues('values.immutable', _immutable2.default, _reduxImmutablejs.combineReducers, (0, _addExpectations2.default)(_expectations4.default));