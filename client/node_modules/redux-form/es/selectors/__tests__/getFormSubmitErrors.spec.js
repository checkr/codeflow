import createGetFormSubmitErrors from '../getFormSubmitErrors';
import plain from '../../structure/plain';
import plainExpectations from '../../structure/plain/expectations';
import immutable from '../../structure/immutable';
import immutableExpectations from '../../structure/immutable/expectations';
import addExpectations from '../../__tests__/addExpectations';

var describeGetFormSubmitErrors = function describeGetFormSubmitErrors(name, structure, expect) {
  var getFormSubmitErrors = createGetFormSubmitErrors(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(createGetFormSubmitErrors('foo')).toBeA('function');
    });

    it('should get the form values from state', function () {
      expect(getFormSubmitErrors('foo')(fromJS({
        form: {
          foo: {
            submitErrors: {
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

    it('should return undefined if there are no submitErrors', function () {
      expect(getFormSubmitErrors('foo')(fromJS({
        form: {
          foo: {}
        }
      }))).toEqual(undefined);
    });

    it('should use getFormState if provided', function () {
      expect(getFormSubmitErrors('foo', function (state) {
        return getIn(state, 'someOtherSlice');
      })(fromJS({
        someOtherSlice: {
          foo: {
            submitErrors: {
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

describeGetFormSubmitErrors('getFormSubmitErrors.plain', plain, addExpectations(plainExpectations));
describeGetFormSubmitErrors('getFormSubmitErrors.immutable', immutable, addExpectations(immutableExpectations));