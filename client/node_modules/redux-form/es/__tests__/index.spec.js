import expect from 'expect';
import * as expectedActionTypes from '../actionTypes';
import expectedPropTypes from '../propTypes';
import { actionTypes, arrayInsert, arrayMove, arrayPop, arrayPush, arrayRemove, arrayRemoveAll, arrayShift, arraySplice, arraySwap, arrayUnshift, autofill, blur, change, destroy, Field, FieldArray, FormSection, focus, formValueSelector, initialize, propTypes, reducer, reduxForm, registerField, reset, setSubmitFailed, setSubmitSucceeded, startAsyncValidation, startSubmit, stopAsyncValidation, stopSubmit, SubmissionError, touch, unregisterField, untouch, values } from '../index';

describe('index', function () {
  it('should export actionTypes', function () {
    expect(actionTypes).toEqual(expectedActionTypes);
  });
  it('should export arrayInsert', function () {
    expect(arrayInsert).toExist().toBeA('function');
  });
  it('should export arrayMove', function () {
    expect(arrayMove).toExist().toBeA('function');
  });
  it('should export arrayPop', function () {
    expect(arrayPop).toExist().toBeA('function');
  });
  it('should export arrayPush', function () {
    expect(arrayPush).toExist().toBeA('function');
  });
  it('should export arrayRemove', function () {
    expect(arrayRemove).toExist().toBeA('function');
  });
  it('should export arrayRemoveAll', function () {
    expect(arrayRemoveAll).toExist().toBeA('function');
  });
  it('should export arrayShift', function () {
    expect(arrayShift).toExist().toBeA('function');
  });
  it('should export arraySplice', function () {
    expect(arraySplice).toExist().toBeA('function');
  });
  it('should export arraySwap', function () {
    expect(arraySwap).toExist().toBeA('function');
  });
  it('should export arrayUnshift', function () {
    expect(arrayUnshift).toExist().toBeA('function');
  });
  it('should export autofill', function () {
    expect(autofill).toExist().toBeA('function');
  });
  it('should export blur', function () {
    expect(blur).toExist().toBeA('function');
  });
  it('should export change', function () {
    expect(change).toExist().toBeA('function');
  });
  it('should export destroy', function () {
    expect(destroy).toExist().toBeA('function');
  });
  it('should export Field', function () {
    expect(Field).toExist().toBeA('function');
  });
  it('should export FieldArray', function () {
    expect(FieldArray).toExist().toBeA('function');
  });
  it('should export FormSection', function () {
    expect(FormSection).toExist().toBeA('function');
  });
  it('should export focus', function () {
    expect(focus).toExist().toBeA('function');
  });
  it('should export formValueSelector', function () {
    expect(formValueSelector).toExist().toBeA('function');
  });
  it('should export initialize', function () {
    expect(initialize).toExist().toBeA('function');
  });
  it('should export propTypes', function () {
    expect(propTypes).toEqual(expectedPropTypes);
  });
  it('should export reducer', function () {
    expect(reducer).toExist().toBeA('function');
  });
  it('should export reduxForm', function () {
    expect(reduxForm).toExist().toBeA('function');
  });
  it('should export registerField', function () {
    expect(registerField).toExist().toBeA('function');
  });
  it('should export reset', function () {
    expect(reset).toExist().toBeA('function');
  });
  it('should export startAsyncValidation', function () {
    expect(startAsyncValidation).toExist().toBeA('function');
  });
  it('should export startSubmit', function () {
    expect(startSubmit).toExist().toBeA('function');
  });
  it('should export setSubmitFailed', function () {
    expect(setSubmitFailed).toExist().toBeA('function');
  });
  it('should export setSubmitSucceeded', function () {
    expect(setSubmitSucceeded).toExist().toBeA('function');
  });
  it('should export stopAsyncValidation', function () {
    expect(stopAsyncValidation).toExist().toBeA('function');
  });
  it('should export stopSubmit', function () {
    expect(stopSubmit).toExist().toBeA('function');
  });
  it('should export SubmissionError', function () {
    expect(SubmissionError).toExist().toBeA('function');
  });
  it('should export touch', function () {
    expect(touch).toExist().toBeA('function');
  });
  it('should export unregisterField', function () {
    expect(unregisterField).toExist().toBeA('function');
  });
  it('should export untouch', function () {
    expect(untouch).toExist().toBeA('function');
  });
  it('should export values', function () {
    expect(values).toExist().toBeA('function');
  });
});