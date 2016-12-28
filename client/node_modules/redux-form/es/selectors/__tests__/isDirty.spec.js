import createIsDirty from '../isDirty';
import plain from '../../structure/plain';
import plainExpectations from '../../structure/plain/expectations';
import immutable from '../../structure/immutable';
import immutableExpectations from '../../structure/immutable/expectations';
import addExpectations from '../../__tests__/addExpectations';

var describeIsDirty = function describeIsDirty(name, structure, expect) {
  var isDirty = createIsDirty(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(isDirty('foo')).toBeA('function');
    });

    it('should return false when values not present', function () {
      expect(isDirty('foo')(fromJS({
        form: {}
      }))).toBe(false);
    });

    it('should return false when values are pristine', function () {
      expect(isDirty('foo')(fromJS({
        form: {
          foo: {
            initial: {
              dog: 'Snoopy',
              cat: 'Garfield'
            },
            values: {
              dog: 'Snoopy',
              cat: 'Garfield'
            }
          }
        }
      }))).toBe(false);
    });

    it('should return true when values are dirty', function () {
      expect(isDirty('foo')(fromJS({
        form: {
          foo: {
            initial: {
              dog: 'Snoopy',
              cat: 'Garfield'
            },
            values: {
              dog: 'Odie',
              cat: 'Garfield'
            }
          }
        }
      }))).toBe(true);
    });

    it('should use getFormState if provided', function () {
      expect(isDirty('foo', function (state) {
        return getIn(state, 'someOtherSlice');
      })(fromJS({
        someOtherSlice: {
          foo: {
            initial: {
              dog: 'Snoopy',
              cat: 'Garfield'
            },
            values: {
              dog: 'Odie',
              cat: 'Garfield'
            }
          }
        }
      }))).toBe(true);
    });
  });
};

describeIsDirty('isDirty.plain', plain, addExpectations(plainExpectations));
describeIsDirty('isDirty.immutable', immutable, addExpectations(immutableExpectations));