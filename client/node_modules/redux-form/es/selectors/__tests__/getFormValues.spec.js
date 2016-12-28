import createGetFormValues from '../getFormValues';
import plain from '../../structure/plain';
import plainExpectations from '../../structure/plain/expectations';
import immutable from '../../structure/immutable';
import immutableExpectations from '../../structure/immutable/expectations';
import addExpectations from '../../__tests__/addExpectations';

var describeGetFormValues = function describeGetFormValues(name, structure, expect) {
  var getFormValues = createGetFormValues(structure);

  var fromJS = structure.fromJS,
      getIn = structure.getIn;


  describe(name, function () {
    it('should return a function', function () {
      expect(getFormValues('foo')).toBeA('function');
    });

    it('should get the form values from state', function () {
      expect(getFormValues('foo')(fromJS({
        form: {
          foo: {
            values: {
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

    it('should use getFormState if provided', function () {
      expect(getFormValues('foo', function (state) {
        return getIn(state, 'someOtherSlice');
      })(fromJS({
        someOtherSlice: {
          foo: {
            values: {
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

describeGetFormValues('getFormValues.plain', plain, addExpectations(plainExpectations));
describeGetFormValues('getFormValues.immutable', immutable, addExpectations(immutableExpectations));