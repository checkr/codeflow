'use strict';

var _noop2 = require('lodash/noop');

var _noop3 = _interopRequireDefault(_noop2);

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _getValue = require('../getValue');

var _getValue2 = _interopRequireDefault(_getValue);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('getValue', function () {
  it('should return value if non-event value is passed', function () {
    (0, _expect2.default)((0, _getValue2.default)(undefined, true)).toBe(undefined);
    (0, _expect2.default)((0, _getValue2.default)(undefined, false)).toBe(undefined);
    (0, _expect2.default)((0, _getValue2.default)(null, true)).toBe(null);
    (0, _expect2.default)((0, _getValue2.default)(null, false)).toBe(null);
    (0, _expect2.default)((0, _getValue2.default)(5, true)).toBe(5);
    (0, _expect2.default)((0, _getValue2.default)(5, false)).toBe(5);
    (0, _expect2.default)((0, _getValue2.default)(true, true)).toBe(true);
    (0, _expect2.default)((0, _getValue2.default)(true, false)).toBe(true);
    (0, _expect2.default)((0, _getValue2.default)(false, true)).toBe(false);
    (0, _expect2.default)((0, _getValue2.default)(false, false)).toBe(false);
    (0, _expect2.default)((0, _getValue2.default)('dog', true)).toBe('dog');
    (0, _expect2.default)((0, _getValue2.default)('dog', false)).toBe('dog');
  });

  it('should not unwrap value if non-event object containing value key is passed', function () {
    (0, _expect2.default)((0, _getValue2.default)({ value: 5 }, true)).toEqual({ value: 5 });
    (0, _expect2.default)((0, _getValue2.default)({ value: 5 }, false)).toEqual({ value: 5 });
    (0, _expect2.default)((0, _getValue2.default)({ value: true }, true)).toEqual({ value: true });
    (0, _expect2.default)((0, _getValue2.default)({ value: true }, false)).toEqual({ value: true });
    (0, _expect2.default)((0, _getValue2.default)({ value: false }, true)).toEqual({ value: false });
    (0, _expect2.default)((0, _getValue2.default)({ value: false }, false)).toEqual({ value: false });
  });

  it('should return value if object NOT containing value key is passed', function () {
    var foo = { bar: 5, baz: 8 };
    (0, _expect2.default)((0, _getValue2.default)(foo)).toBe(foo);
  });

  it('should return event.nativeEvent.text if defined and not react-native', function () {
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      nativeEvent: {
        text: 'foo'
      }
    }, false)).toBe('foo');
  });

  it('should return event.nativeEvent.text if react-native', function () {
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      nativeEvent: {
        text: 'foo'
      }
    }, true)).toBe('foo');
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      nativeEvent: {
        text: undefined
      }
    }, true)).toBe(undefined);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      nativeEvent: {
        text: null
      }
    }, true)).toBe(null);
  });

  it('should return event.target.checked if checkbox', function () {
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'checkbox',
        checked: true
      }
    }, true)).toBe(true);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'checkbox',
        checked: true
      }
    }, false)).toBe(true);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'checkbox',
        checked: false
      }
    }, true)).toBe(false);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'checkbox',
        checked: false
      }
    }, false)).toBe(false);
  });

  it('should return a number type for numeric inputs, when a value is set', function () {
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'number',
        value: '3.1415'
      }
    }, true)).toBe('3.1415');
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'range',
        value: '2.71828'
      }
    }, true)).toBe('2.71828');
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'number',
        value: '3'
      }
    }, false)).toBe('3');
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'range',
        value: '3.1415'
      }
    }, false)).toBe('3.1415');

    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'range',
        value: ''
      }
    }, false)).toBe('');
  });

  it('should return event.target.files if file', function () {
    var myFiles = ['foo', 'bar'];
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'file',
        files: myFiles
      }
    }, true)).toBe(myFiles);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'file',
        files: myFiles
      }
    }, false)).toBe(myFiles);
  });

  it('should return event.dataTransfer.files if file and files not in target.files', function () {
    var myFiles = ['foo', 'bar'];
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'file'
      },
      dataTransfer: {
        files: myFiles
      }
    }, true)).toBe(myFiles);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'file'
      },
      dataTransfer: {
        files: myFiles
      }
    }, false)).toBe(myFiles);
  });

  it('should return selected options if is a multiselect', function () {
    var options = [{ selected: true, value: 'foo' }, { selected: true, value: 'bar' }, { selected: false, value: 'baz' }];
    var expected = options.filter(function (option) {
      return option.selected;
    }).map(function (option) {
      return option.value;
    });
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'select-multiple',
        options: options
      }
    }, true)).toEqual(expected);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'select-multiple'
      }
    }, false)).toEqual([]); // no options specified
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        type: 'select-multiple',
        options: options
      }
    }, false)).toEqual(expected);
  });

  it('should return event.target.value if not file or checkbox', function () {
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: undefined
      }
    }, true)).toBe(undefined);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: undefined
      }
    }, false)).toBe(undefined);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: null
      }
    }, true)).toBe(null);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: null
      }
    }, false)).toBe(null);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: true
      }
    }, true)).toBe(true);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: true
      }
    }, false)).toBe(true);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: false
      }
    }, true)).toBe(false);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: false
      }
    }, false)).toBe(false);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: 42
      }
    }, true)).toBe(42);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: 42
      }
    }, false)).toBe(42);
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: 'foo'
      }
    }, true)).toBe('foo');
    (0, _expect2.default)((0, _getValue2.default)({
      preventDefault: _noop3.default,
      stopPropagation: _noop3.default,
      target: {
        value: 'foo'
      }
    }, false)).toBe('foo');
  });
});