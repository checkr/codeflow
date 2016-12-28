import createHasError from '../hasError';
import plain from '../structure/plain';
import plainExpectations from '../structure/plain/expectations';
import immutable from '../structure/immutable';
import immutableExpectations from '../structure/immutable/expectations';
import addExpectations from './addExpectations';

var describeHasError = function describeHasError(name, structure, expect) {
  var fromJS = structure.fromJS,
      getIn = structure.getIn;

  var hasError = createHasError(structure);

  describe(name, function () {
    it('should return false for falsy values', function () {
      var field = fromJS({ name: 'foo', type: 'Field' });
      expect(hasError(field, undefined)).toBe(false);
      expect(hasError(field, null)).toBe(false);
      expect(hasError(field, '')).toBe(false);
      expect(hasError(field, 0)).toBe(false);
      expect(hasError(field, false)).toBe(false);
    });

    it('should return false for empty structures', function () {
      var field = fromJS({ name: 'foo', type: 'Field' });
      var obj = fromJS({});
      var array = fromJS([]);
      expect(hasError(field, obj, obj, obj)).toBe(false);
      expect(hasError(field, array, array, array)).toBe(false);
    });

    it('should return false for deeply nested structures with undefined values', function () {
      var field1 = fromJS({ name: 'nested.myArrayField', type: 'FieldArray' });
      expect(hasError(field1, fromJS({
        nested: {
          myArrayField: [undefined, undefined]
        }
      }))).toBe(false);
      var field2 = fromJS({ name: 'nested.deeper.foo', type: 'Field' });
      expect(hasError(field2, fromJS({
        nested: {
          deeper: {
            foo: undefined,
            bar: undefined
          }
        }
      }))).toBe(false);
    });

    it('should return true for errors that match a Field', function () {
      var field = fromJS({ name: 'foo.bar', type: 'Field' });
      var plainError = {
        foo: {
          bar: 'An error'
        }
      };
      var error = fromJS(plainError);
      expect(hasError(field, plainError)).toBe(true);
      expect(hasError(field, null, error)).toBe(true);
      expect(hasError(field, null, null, error)).toBe(true);
    });

    it('should return false for errors that don\'t match a Field', function () {
      var field = fromJS({ name: 'foo.baz', type: 'Field' });
      var error = fromJS({
        foo: {
          bar: 'An error'
        }
      });
      expect(hasError(field, error)).toBe(false);
      expect(hasError(field, null, error)).toBe(false);
      expect(hasError(field, null, null, error)).toBe(false);
    });

    it('should return true for errors that match a FieldArray', function () {
      var field = fromJS({ name: 'foo.bar', type: 'FieldArray' });
      var plainError = {
        foo: {
          bar: ['An error']
        }
      };
      plainError.foo.bar._error = 'An error';

      expect(hasError(field, plainError)).toBe(true);

      var error = fromJS(plainError);
      if (getIn(error, 'foo.bar._error') === 'An error') {
        // cannot work for Immutable Lists because you can not set a value under a string key
        expect(hasError(field, null, error)).toBe(true);
        expect(hasError(field, null, null, error)).toBe(true);
      }
    });

    it('should return false for errors that don\'t match a FieldArray', function () {
      var field = fromJS({ name: 'foo.baz', type: 'FieldArray' });
      var plainError = {
        foo: {
          bar: ['An error']
        }
      };
      plainError.foo.bar._error = 'An error';

      expect(hasError(field, plainError)).toBe(false);

      var error = fromJS(plainError);
      if (getIn(error, 'foo.bar._error') === 'An error') {
        // cannot work for Immutable Lists because you can not set a value under a string key
        expect(hasError(field, null, error)).toBe(false);
        expect(hasError(field, null, null, error)).toBe(false);
      }
    });

    it('should return true if a Field that has an object value has an _error', function () {
      var field = fromJS({ name: 'foo', type: 'Field' });
      var plainError = {
        foo: {
          _error: 'An error'
        }
      };

      expect(hasError(field, plainError)).toBe(true);

      var error = fromJS(plainError);
      expect(hasError(field, null, error)).toBe(true);
      expect(hasError(field, null, null, error)).toBe(true);
    });
  });
};

describeHasError('hasError.plain', plain, addExpectations(plainExpectations));
describeHasError('hasError.immutable', immutable, addExpectations(immutableExpectations));