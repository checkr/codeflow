import expect from 'expect';
import getDisplayName from '../getDisplayName';

describe('getDisplayName', function () {
  it('should read displayName', function () {
    expect(getDisplayName({ displayName: 'foo' })).toBe('foo');
  });

  it('should read name', function () {
    expect(getDisplayName({ name: 'foo' })).toBe('foo');
  });

  it('should default to Component', function () {
    expect(getDisplayName({})).toBe('Component');
  });
});