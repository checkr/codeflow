/* eslint react/no-multi-comp:0 */
import createFormValueSelector from '../formValueSelector';
import plain from '../structure/plain';
import plainExpectations from '../structure/plain/expectations';
import immutable from '../structure/immutable';
import immutableExpectations from '../structure/immutable/expectations';
import addExpectations from './addExpectations';

var describeFormValueSelector = function describeFormValueSelector(name, structure, expect) {
  var fromJS = structure.fromJS,
      getIn = structure.getIn;

  var formValueSelector = createFormValueSelector(structure);

  describe(name, function () {
    it('should throw an error if no form specified', function () {
      expect(function () {
        return formValueSelector();
      }).toThrow('Form value must be specified');
    });

    it('should return a function', function () {
      expect(formValueSelector('myForm')).toBeA('function');
    });

    it('should throw an error if no fields specified', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({});
      expect(function () {
        return selector(state);
      }).toThrow('No fields specified');
    });

    it('should return undefined for a single value when no redux-form state found', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({});
      expect(selector(state, 'foo')).toBe(undefined);
    });

    it('should return undefined for a single value when no form slice found', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {}
      });
      expect(selector(state, 'foo')).toBe(undefined);
    });

    it('should return undefined for a single value when no values found', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            // no values
          }
        }
      });
      expect(selector(state, 'foo')).toBe(undefined);
    });

    it('should get a single value', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            values: {
              foo: 'bar'
            }
          }
        }
      });
      expect(selector(state, 'foo')).toBe('bar');
    });

    it('should get a single deep value', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            values: {
              dog: {
                cat: {
                  ewe: {
                    pig: 'Napoleon'
                  }
                }
              }
            }
          }
        }
      });
      expect(selector(state, 'dog.cat.ewe.pig')).toBe('Napoleon');
    });

    it('should return {} for multiple values when no redux-form state found', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({});
      expect(selector(state, 'foo', 'bar')).toEqual({});
    });

    it('should return {} for multiple values when no form slice found', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {}
      });
      expect(selector(state, 'foo', 'bar')).toEqual({});
    });

    it('should return {} for multiple values when no values found', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            // no values
          }
        }
      });
      expect(selector(state, 'foo', 'bar')).toEqual({});
    });

    it('should get multiple values', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            values: {
              foo: 'bar',
              dog: 'cat',
              another: 'do not read'
            }
          }
        }
      });
      expect(selector(state, 'foo', 'dog')).toEqual({
        foo: 'bar',
        dog: 'cat'
      });
    });

    it('should get multiple deep values', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            values: {
              dog: {
                cat: {
                  ewe: {
                    pig: 'Napoleon'
                  }
                },
                rat: {
                  hog: 'Wilbur'
                }
              }
            }
          }
        }
      });
      expect(selector(state, 'dog.cat.ewe.pig', 'dog.rat.hog')).toEqual({
        dog: {
          cat: {
            ewe: {
              pig: 'Napoleon'
            }
          },
          rat: {
            hog: 'Wilbur'
          }
        }
      });
    });

    it('should get an array', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            values: {
              mice: ['Jaq', 'Gus', 'Major', 'Bruno']
            }
          }
        }
      });
      expect(selector(state, 'mice')).toEqualMap(['Jaq', 'Gus', 'Major', 'Bruno']);
    });

    it('should get a deep array', function () {
      var selector = formValueSelector('myForm');
      var state = fromJS({
        form: {
          myForm: {
            values: {
              rodent: {
                rat: {
                  hog: 'Wilbur'
                },
                mice: ['Jaq', 'Gus', 'Major', 'Bruno']
              }
            }
          }
        }
      });
      expect(selector(state, 'rodent.rat.hog', 'rodent.mice')).toEqual({
        rodent: {
          rat: {
            hog: 'Wilbur'
          },
          mice: fromJS(['Jaq', 'Gus', 'Major', 'Bruno'])
        }
      });
    });

    it('should get a single value using a different mount point', function () {
      var selector = formValueSelector('myForm', function (state) {
        return getIn(state, 'otherMountPoint');
      });
      var state = fromJS({
        otherMountPoint: {
          myForm: {
            values: {
              foo: 'bar'
            }
          }
        }
      });
      expect(selector(state, 'foo')).toBe('bar');
    });
  });
};

describeFormValueSelector('formValueSelector.plain', plain, addExpectations(plainExpectations));
describeFormValueSelector('formValueSelector.immutable', immutable, addExpectations(immutableExpectations));