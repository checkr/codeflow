var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

/* eslint react/no-multi-comp:0 */
import React from 'react';
import { createSpy } from 'expect';
import { Provider } from 'react-redux';
import { combineReducers as plainCombineReducers, createStore } from 'redux';
import { combineReducers as immutableCombineReducers } from 'redux-immutablejs';
import TestUtils from 'react-addons-test-utils';
import createReducer from '../reducer';
import createValues from '../values';
import plain from '../structure/plain';
import plainExpectations from '../structure/plain/expectations';
import immutable from '../structure/immutable';
import immutableExpectations from '../structure/immutable/expectations';
import addExpectations from './addExpectations';

var describeValues = function describeValues(name, structure, combineReducers, expect) {
  var values = createValues(structure);
  var reducer = createReducer(structure);
  var fromJS = structure.fromJS;

  var makeStore = function makeStore(initial) {
    return createStore(combineReducers({ form: reducer }), fromJS({ form: initial }));
  };

  var testProps = function testProps(state) {
    var config = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};

    var store = makeStore({ testForm: state });
    var spy = createSpy(function () {
      return React.createElement('div', null);
    }).andCallThrough();

    var Decorated = values(_extends({ form: 'testForm' }, config))(spy);
    TestUtils.renderIntoDocument(React.createElement(
      Provider,
      { store: store },
      React.createElement(Decorated, null)
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

describeValues('values.plain', plain, plainCombineReducers, addExpectations(plainExpectations));
describeValues('values.immutable', immutable, immutableCombineReducers, addExpectations(immutableExpectations));