import expect from 'expect';
import prefixName from '../prefixName';

describe('prefixName', function () {
  it('should concat sectionPrefix and name', function () {
    var context = {
      _reduxForm: {
        sectionPrefix: 'foo'
      }
    };
    expect(prefixName(context, 'bar')).toBe('foo.bar');
  });

  it('should ignore empty sectionPrefix', function () {
    var context = {
      _reduxForm: {
        sectionPrefix: undefined
      }
    };
    expect(prefixName(context, 'bar')).toBe('bar');
  });

  it('should not prefix array fields', function () {
    var context = {
      _reduxForm: {
        sectionPrefix: 'foo'
      }
    };
    expect(prefixName(context, 'bar.bar[0]')).toBe('bar.bar[0]');
  });
});