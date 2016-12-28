import createGetFormSyncErrors from '../getFormSyncErrors';
import plain from '../../structure/plain';
import plainExpectations from '../../structure/plain/expectations';
import immutable from '../../structure/immutable';
import immutableExpectations from '../../structure/immutable/expectations';
import addExpectations from '../../__tests__/addExpectations';

var describeGetFormSyncErrors = function describeGetFormSyncErrors(name, structure, expect) {
  var getFormSyncErrors = createGetFormSyncErrors(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(createGetFormSyncErrors('foo')).toBeA('function');
    });

    it('should get the form values from state', function () {
      expect(getFormSyncErrors('foo')(fromJS({
        form: {
          foo: {
            syncErrors: {
              dog: 'Snoopy',
              cat: 'Garfield'
            }
          }
        }
      }))).toEqualMap({
        dog: 'Snoopy',
        cat: 'Garfield'
      });
    });

    it('should return undefined if there are no syncErrors', function () {
      expect(getFormSyncErrors('foo')(fromJS({
        form: {
          foo: {}
        }
      }))).toEqual(undefined);
    });

    it('should use getFormState if provided', function () {
      expect(getFormSyncErrors('foo', function (state) {
        return getIn(state, 'someOtherSlice');
      })(fromJS({
        someOtherSlice: {
          foo: {
            syncErrors: {
              dog: 'Snoopy',
              cat: 'Garfield'
            }
          }
        }
      }))).toEqualMap({
        dog: 'Snoopy',
        cat: 'Garfield'
      });
    });
  });
};

describeGetFormSyncErrors('getFormSyncErrors.plain', plain, addExpectations(plainExpectations));
describeGetFormSyncErrors('getFormSyncErrors.immutable', immutable, addExpectations(immutableExpectations));