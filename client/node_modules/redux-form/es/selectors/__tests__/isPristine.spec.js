import createIsPristine from '../isPristine';
import plain from '../../structure/plain';
import plainExpectations from '../../structure/plain/expectations';
import immutable from '../../structure/immutable';
import immutableExpectations from '../../structure/immutable/expectations';
import addExpectations from '../../__tests__/addExpectations';

var describeIsPristine = function describeIsPristine(name, structure, expect) {
  var isPristine = createIsPristine(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(isPristine('foo')).toBeA('function');
    });

    it('should return true when values not present', function () {
      expect(isPristine('foo')(fromJS({
        form: {}
      }))).toBe(true);
    });

    it('should return true when values are pristine', function () {
      expect(isPristine('foo')(fromJS({
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
      }))).toBe(true);
    });

    it('should return true when values are dirty', function () {
      expect(isPristine('foo')(fromJS({
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
      }))).toBe(false);
    });

    it('should use getFormState if provided', function () {
      expect(isPristine('foo', function (state) {
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
      }))).toBe(false);
    });
  });
};

describeIsPristine('isPristine.plain', plain, addExpectations(plainExpectations));
describeIsPristine('isPristine.immutable', immutable, addExpectations(immutableExpectations));