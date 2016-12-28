'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _expectPredicate = require('expect-predicate');

var _expectPredicate2 = _interopRequireDefault(_expectPredicate);

var _actionTypes = require('../actionTypes');

var _actions = require('../actions');

var _fluxStandardAction = require('flux-standard-action');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

_expect2.default.extend(_expectPredicate2.default);

describe('actions', function () {

  it('should create array insert action', function () {
    (0, _expect2.default)((0, _actions.arrayInsert)('myForm', 'myField', 0, 'foo')).toEqual({
      type: _actionTypes.ARRAY_INSERT,
      meta: {
        form: 'myForm',
        field: 'myField',
        index: 0
      },
      payload: 'foo'
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array move action', function () {
    (0, _expect2.default)((0, _actions.arrayMove)('myForm', 'myField', 2, 4)).toEqual({
      type: _actionTypes.ARRAY_MOVE,
      meta: {
        form: 'myForm',
        field: 'myField',
        from: 2,
        to: 4
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array pop action', function () {
    (0, _expect2.default)((0, _actions.arrayPop)('myForm', 'myField')).toEqual({
      type: _actionTypes.ARRAY_POP,
      meta: {
        form: 'myForm',
        field: 'myField'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array push action', function () {
    (0, _expect2.default)((0, _actions.arrayPush)('myForm', 'myField', 'foo')).toEqual({
      type: _actionTypes.ARRAY_PUSH,
      meta: {
        form: 'myForm',
        field: 'myField'
      },
      payload: 'foo'
    }).toPass(_fluxStandardAction.isFSA);

    (0, _expect2.default)((0, _actions.arrayPush)('myForm', 'myField')).toEqual({
      type: _actionTypes.ARRAY_PUSH,
      meta: {
        form: 'myForm',
        field: 'myField'
      },
      payload: undefined
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array remove action', function () {
    (0, _expect2.default)((0, _actions.arrayRemove)('myForm', 'myField', 3)).toEqual({
      type: _actionTypes.ARRAY_REMOVE,
      meta: {
        form: 'myForm',
        field: 'myField',
        index: 3
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array removeAll action', function () {
    (0, _expect2.default)((0, _actions.arrayRemoveAll)('myForm', 'myField')).toEqual({
      type: _actionTypes.ARRAY_REMOVE_ALL,
      meta: {
        form: 'myForm',
        field: 'myField'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array shift action', function () {
    (0, _expect2.default)((0, _actions.arrayShift)('myForm', 'myField')).toEqual({
      type: _actionTypes.ARRAY_SHIFT,
      meta: {
        form: 'myForm',
        field: 'myField'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array splice action', function () {
    (0, _expect2.default)((0, _actions.arraySplice)('myForm', 'myField', 1, 1)).toEqual({
      type: _actionTypes.ARRAY_SPLICE,
      meta: {
        form: 'myForm',
        field: 'myField',
        index: 1,
        removeNum: 1
      }
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.arraySplice)('myForm', 'myField', 2, 1)).toEqual({
      type: _actionTypes.ARRAY_SPLICE,
      meta: {
        form: 'myForm',
        field: 'myField',
        index: 2,
        removeNum: 1
      }
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.arraySplice)('myForm', 'myField', 2, 0, 'foo')).toEqual({
      type: _actionTypes.ARRAY_SPLICE,
      meta: {
        form: 'myForm',
        field: 'myField',
        index: 2,
        removeNum: 0
      },
      payload: 'foo'
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.arraySplice)('myForm', 'myField', 3, 2, { foo: 'bar' })).toEqual({
      type: _actionTypes.ARRAY_SPLICE,
      meta: {
        form: 'myForm',
        field: 'myField',
        index: 3,
        removeNum: 2
      },
      payload: { foo: 'bar' }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array unshift action', function () {
    (0, _expect2.default)((0, _actions.arrayUnshift)('myForm', 'myField', 'foo')).toEqual({
      type: _actionTypes.ARRAY_UNSHIFT,
      meta: {
        form: 'myForm',
        field: 'myField'
      },
      payload: 'foo'
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create array swap action', function () {
    (0, _expect2.default)((0, _actions.arraySwap)('myForm', 'myField', 0, 8)).toEqual({
      type: _actionTypes.ARRAY_SWAP,
      meta: {
        form: 'myForm',
        field: 'myField',
        indexA: 0,
        indexB: 8
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should throw an exception with illegal array swap indices', function () {
    (0, _expect2.default)(function () {
      return (0, _actions.arraySwap)('myForm', 'myField', 2, 2);
    }).toThrow('Swap indices cannot be equal');
    (0, _expect2.default)(function () {
      return (0, _actions.arraySwap)('myForm', 'myField', -2, 2);
    }).toThrow('Swap indices cannot be negative');
    (0, _expect2.default)(function () {
      return (0, _actions.arraySwap)('myForm', 'myField', 2, -2);
    }).toThrow('Swap indices cannot be negative');
  });

  it('should create blur action', function () {
    (0, _expect2.default)((0, _actions.blur)('myForm', 'myField', 'bar', false)).toEqual({
      type: _actionTypes.BLUR,
      meta: {
        form: 'myForm',
        field: 'myField',
        touch: false
      },
      payload: 'bar'
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.blur)('myForm', 'myField', 7, true)).toEqual({
      type: _actionTypes.BLUR,
      meta: {
        form: 'myForm',
        field: 'myField',
        touch: true
      },
      payload: 7
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create change action', function () {
    (0, _expect2.default)((0, _actions.change)('myForm', 'myField', 'bar', false, true)).toEqual({
      type: _actionTypes.CHANGE,
      meta: {
        form: 'myForm',
        field: 'myField',
        touch: false,
        persistentSubmitErrors: true
      },
      payload: 'bar'
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.change)('myForm', 'myField', 7, true, false)).toEqual({
      type: _actionTypes.CHANGE,
      meta: {
        form: 'myForm',
        field: 'myField',
        touch: true,
        persistentSubmitErrors: false
      },
      payload: 7
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create focus action', function () {
    (0, _expect2.default)((0, _actions.focus)('myForm', 'myField')).toEqual({
      type: _actionTypes.FOCUS,
      meta: {
        form: 'myForm',
        field: 'myField'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create clear submit action', function () {
    (0, _expect2.default)((0, _actions.clearSubmit)('myForm')).toEqual({
      type: _actionTypes.CLEAR_SUBMIT,
      meta: {
        form: 'myForm'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create initialize action', function () {
    var data = { a: 8, c: 9 };
    (0, _expect2.default)((0, _actions.initialize)('myForm', data)).toEqual({
      type: _actionTypes.INITIALIZE,
      meta: {
        form: 'myForm',
        keepDirty: undefined
      },
      payload: data
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create initialize action with a keepDirty value', function () {
    var data = { a: 8, c: 9 };
    (0, _expect2.default)((0, _actions.initialize)('myForm', data, true)).toEqual({
      type: _actionTypes.INITIALIZE,
      meta: {
        form: 'myForm',
        keepDirty: true
      },
      payload: data
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create registerField action', function () {
    (0, _expect2.default)((0, _actions.registerField)('myForm', 'foo', 'Field')).toEqual({
      type: _actionTypes.REGISTER_FIELD,
      meta: {
        form: 'myForm'
      },
      payload: {
        name: 'foo',
        type: 'Field'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create reset action', function () {
    (0, _expect2.default)((0, _actions.reset)('myForm')).toEqual({
      type: _actionTypes.RESET,
      meta: {
        form: 'myForm'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create destroy action', function () {
    (0, _expect2.default)((0, _actions.destroy)('myForm')).toEqual({
      type: _actionTypes.DESTROY,
      meta: {
        form: 'myForm'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create startAsyncValidation action', function () {
    (0, _expect2.default)((0, _actions.startAsyncValidation)('myForm', 'myField')).toEqual({
      type: _actionTypes.START_ASYNC_VALIDATION,
      meta: {
        form: 'myForm',
        field: 'myField'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create startSubmit action', function () {
    (0, _expect2.default)((0, _actions.startSubmit)('myForm')).toEqual({
      type: _actionTypes.START_SUBMIT,
      meta: {
        form: 'myForm'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create startSubmit action', function () {
    (0, _expect2.default)((0, _actions.startSubmit)('myForm')).toEqual({
      type: _actionTypes.START_SUBMIT,
      meta: {
        form: 'myForm'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create stopAsyncValidation action', function () {
    var errors = {
      foo: 'Foo error',
      bar: 'Error for bar'
    };
    (0, _expect2.default)((0, _actions.stopAsyncValidation)('myForm', errors)).toEqual({
      type: _actionTypes.STOP_ASYNC_VALIDATION,
      meta: {
        form: 'myForm'
      },
      payload: errors,
      error: true
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create stopSubmit action', function () {
    (0, _expect2.default)((0, _actions.stopSubmit)('myForm')).toEqual({
      type: _actionTypes.STOP_SUBMIT,
      meta: {
        form: 'myForm'
      },
      payload: undefined
    }).toPass(_fluxStandardAction.isFSA);
    var errors = {
      foo: 'Foo error',
      bar: 'Error for bar'
    };
    (0, _expect2.default)((0, _actions.stopSubmit)('myForm', errors)).toEqual({
      type: _actionTypes.STOP_SUBMIT,
      meta: {
        form: 'myForm'
      },
      payload: errors,
      error: true
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create submit action', function () {
    (0, _expect2.default)((0, _actions.submit)('myForm')).toEqual({
      type: _actionTypes.SUBMIT,
      meta: {
        form: 'myForm'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create setSubmitFailed action', function () {
    (0, _expect2.default)((0, _actions.setSubmitFailed)('myForm')).toEqual({
      type: _actionTypes.SET_SUBMIT_FAILED,
      meta: {
        form: 'myForm',
        fields: []
      },
      error: true
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.setSubmitFailed)('myForm', 'a', 'b', 'c')).toEqual({
      type: _actionTypes.SET_SUBMIT_FAILED,
      meta: {
        form: 'myForm',
        fields: ['a', 'b', 'c']
      },
      error: true
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create setSubmitSucceeded action', function () {
    (0, _expect2.default)((0, _actions.setSubmitSucceeded)('myForm')).toEqual({
      type: _actionTypes.SET_SUBMIT_SUCCEEDED,
      meta: {
        form: 'myForm',
        fields: []
      },
      error: false
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.setSubmitSucceeded)('myForm', 'a', 'b', 'c')).toEqual({
      type: _actionTypes.SET_SUBMIT_SUCCEEDED,
      meta: {
        form: 'myForm',
        fields: ['a', 'b', 'c']
      },
      error: false
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create touch action', function () {
    (0, _expect2.default)((0, _actions.touch)('myForm', 'foo', 'bar')).toEqual({
      type: _actionTypes.TOUCH,
      meta: {
        form: 'myForm',
        fields: ['foo', 'bar']
      }
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.touch)('myForm', 'cat', 'dog', 'pig')).toEqual({
      type: _actionTypes.TOUCH,
      meta: {
        form: 'myForm',
        fields: ['cat', 'dog', 'pig']
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create unregisterField action', function () {
    (0, _expect2.default)((0, _actions.unregisterField)('myForm', 'foo')).toEqual({
      type: _actionTypes.UNREGISTER_FIELD,
      meta: {
        form: 'myForm'
      },
      payload: {
        name: 'foo'
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create untouch action', function () {
    (0, _expect2.default)((0, _actions.untouch)('myForm', 'foo', 'bar')).toEqual({
      type: _actionTypes.UNTOUCH,
      meta: {
        form: 'myForm',
        fields: ['foo', 'bar']
      }
    }).toPass(_fluxStandardAction.isFSA);
    (0, _expect2.default)((0, _actions.untouch)('myForm', 'cat', 'dog', 'pig')).toEqual({
      type: _actionTypes.UNTOUCH,
      meta: {
        form: 'myForm',
        fields: ['cat', 'dog', 'pig']
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create updateSyncErrors action', function () {
    (0, _expect2.default)((0, _actions.updateSyncErrors)('myForm', { foo: 'foo error' })).toEqual({
      type: _actionTypes.UPDATE_SYNC_ERRORS,
      meta: {
        form: 'myForm'
      },
      payload: {
        error: undefined,
        syncErrors: {
          foo: 'foo error'
        }
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create updateSyncErrors action with no errors if none given', function () {
    (0, _expect2.default)((0, _actions.updateSyncErrors)('myForm')).toEqual({
      type: _actionTypes.UPDATE_SYNC_ERRORS,
      meta: {
        form: 'myForm'
      },
      payload: {
        error: undefined,
        syncErrors: {}
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create updateSyncWarnings action', function () {
    (0, _expect2.default)((0, _actions.updateSyncWarnings)('myForm', { foo: 'foo warning' })).toEqual({
      type: _actionTypes.UPDATE_SYNC_WARNINGS,
      meta: {
        form: 'myForm'
      },
      payload: {
        warning: undefined,
        syncWarnings: {
          foo: 'foo warning'
        }
      }
    }).toPass(_fluxStandardAction.isFSA);
  });

  it('should create updateSyncWarnings action with no warnings if none given', function () {
    (0, _expect2.default)((0, _actions.updateSyncWarnings)('myForm')).toEqual({
      type: _actionTypes.UPDATE_SYNC_WARNINGS,
      meta: {
        form: 'myForm'
      },
      payload: {
        warning: undefined,
        syncWarnings: {}
      }
    }).toPass(_fluxStandardAction.isFSA);
  });
});