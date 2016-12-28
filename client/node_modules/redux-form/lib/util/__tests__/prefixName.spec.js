'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _prefixName = require('../prefixName');

var _prefixName2 = _interopRequireDefault(_prefixName);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('prefixName', function () {
  it('should concat sectionPrefix and name', function () {
    var context = {
      _reduxForm: {
        sectionPrefix: 'foo'
      }
    };
    (0, _expect2.default)((0, _prefixName2.default)(context, 'bar')).toBe('foo.bar');
  });

  it('should ignore empty sectionPrefix', function () {
    var context = {
      _reduxForm: {
        sectionPrefix: undefined
      }
    };
    (0, _expect2.default)((0, _prefixName2.default)(context, 'bar')).toBe('bar');
  });

  it('should not prefix array fields', function () {
    var context = {
      _reduxForm: {
        sectionPrefix: 'foo'
      }
    };
    (0, _expect2.default)((0, _prefixName2.default)(context, 'bar.bar[0]')).toBe('bar.bar[0]');
  });
});