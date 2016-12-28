'use strict';

var _formValueSelector = require('../formValueSelector');

var _formValueSelector2 = _interopRequireDefault(_formValueSelector);

var _plain = require('../structure/plain');

var _plain2 = _interopRequireDefault(_plain);

var _expectations = require('../structure/plain/expectations');

var _expectations2 = _interopRequireDefault(_expectations);

var _immutable = require('../structure/immutable');

var _immutable2 = _interopRequireDefault(_immutable);

var _expectations3 = require('../structure/immutable/expectations');

var _expectations4 = _interopRequireDefault(_expectations3);

var _addExpectations = require('./addExpectations');

var _addExpectations2 = _interopRequireDefault(_addExpectations);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

/* eslint react/no-multi-comp:0 */
var describeFormValueSelector = function describeFormValueSelector(name, structure, expect) {
  var fromJS = structure.fromJS,
      getIn = structure.getIn;

  var formValueSelector = (0, _formValueSelector2.default)(structure);

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

describeFormValueSelector('formValueSelector.plain', _plain2.default, (0, _addExpectations2.default)(_expectations2.default));
describeFormValueSelector('formValueSelector.immutable', _immutable2.default, (0, _addExpectations2.default)(_expectations4.default));