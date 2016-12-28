'use strict';

var _expect = require('expect');

var _expect2 = _interopRequireDefault(_expect);

var _getDisplayName = require('../getDisplayName');

var _getDisplayName2 = _interopRequireDefault(_getDisplayName);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

describe('getDisplayName', function () {
  it('should read displayName', function () {
    (0, _expect2.default)((0, _getDisplayName2.default)({ displayName: 'foo' })).toBe('foo');
  });

  it('should read name', function () {
    (0, _expect2.default)((0, _getDisplayName2.default)({ name: 'foo' })).toBe('foo');
  });

  it('should default to Component', function () {
    (0, _expect2.default)((0, _getDisplayName2.default)({})).toBe('Component');
  });
});