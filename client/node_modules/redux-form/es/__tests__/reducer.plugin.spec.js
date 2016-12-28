import { CHANGE } from '../actionTypes';
import { change } from '../actions';

var describePlugin = function describePlugin(vanillaReducer, expect, _ref) {
  var fromJS = _ref.fromJS,
      deleteIn = _ref.deleteIn,
      getIn = _ref.getIn,
      setIn = _ref.setIn;
  return function () {
    it('should initialize state when a plugin is given', function () {
      var reducer = vanillaReducer.plugin({
        foo: function foo(state) {
          return state;
        }
      });
      var state = reducer();
      expect(state).toExist().toBeAMap().toBeSize(0);
    });

    it('should let plugin reducer respond to another action type', function () {
      var state1 = fromJS({
        foo: {
          values: {
            cat: 'dog',
            rat: 'hog'
          },
          fields: {
            cat: { touched: true },
            rat: { touched: true }
          }
        }
      });

      var plugin = function plugin(state, action) {
        if (action.type === 'RAT_POISON') {
          var result = state;
          result = deleteIn(result, 'values.rat');
          result = deleteIn(result, 'fields.rat');
          return result;
        }
        return state;
      };

      var reducer = vanillaReducer.plugin({ foo: plugin });

      var state2 = reducer(state1, { type: 'MILK', form: 'foo' });
      expect(state2).toBe(state1); // no change

      var state3 = reducer(state2, { type: 'RAT_POISON', form: 'foo' });
      expect(state3).toEqualMap({
        foo: {
          values: {
            cat: 'dog'
          },
          fields: {
            cat: { touched: true }
          }
        }
      });
    });

    it('should only respond to form specified', function () {
      var state1 = fromJS({
        foo: {
          values: {
            cat: 'dog',
            rat: 'hog'
          },
          fields: {
            cat: { touched: true },
            rat: { touched: true }
          }
        },
        bar: {
          values: {
            cat: 'dog',
            rat: 'hog'
          },
          fields: {
            cat: { touched: true },
            rat: { touched: true }
          }
        }
      });

      var plugin = function plugin(state, action) {
        if (action.type === 'RAT_POISON') {
          var result = state;
          result = deleteIn(result, 'values.rat');
          result = deleteIn(result, 'fields.rat');
          return result;
        }
        return state;
      };

      var reducer = vanillaReducer.plugin({ foo: plugin });

      var state2 = reducer(state1, { type: 'MILK', form: 'foo' });
      expect(state2).toBe(state1); // no change

      var state3 = reducer(state2, { type: 'RAT_POISON', form: 'foo' });
      expect(state3).toEqualMap({
        foo: {
          values: {
            cat: 'dog'
          },
          fields: {
            cat: { touched: true }
          }
        },
        bar: {
          values: {
            cat: 'dog',
            rat: 'hog'
          },
          fields: {
            cat: { touched: true },
            rat: { touched: true }
          }
        }
      });
    });

    it('should be provided the state from before the vanillaReducer', function () {
      var state1 = fromJS({
        foo: {
          values: {
            cat: 'beta',
            lastCat: 'alpha'
          },
          fields: {
            cat: { touched: false },
            lastCat: { touched: false }
          }
        }
      });

      // this plugin will change the value we are after so we can confirm we get the real starting state
      var intermediatePlugin = function intermediatePlugin(state) {
        return setIn(state, 'values.cat', 'zed');
      };

      var plugin = function plugin(state, action, startingState) {
        if (action.type === CHANGE && action.meta.field === 'cat') {
          var result = state;
          result = setIn(result, 'values.lastCat', getIn(startingState, 'values.cat'));
          result = setIn(result, 'fields.lastCat.touched', action.meta.touch);
          return result;
        }
        return state;
      };

      var reducer = vanillaReducer.plugin({ foo: intermediatePlugin }).plugin({ foo: plugin });

      var state2 = reducer(state1, change('foo', 'cat', 'charlie', true, false));

      expect(state2).toEqualMap({
        foo: {
          anyTouched: true,
          values: {
            cat: 'zed',
            lastCat: 'beta'
          },
          fields: {
            cat: { touched: true },
            lastCat: { touched: true }
          }
        }
      });
    });
  };
};

export default describePlugin;