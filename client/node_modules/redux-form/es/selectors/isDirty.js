import createIsPristine from './isPristine';

var createIsDirty = function createIsDirty(structure) {
  return function (form, getFormState) {
    var isPristine = createIsPristine(structure)(form, getFormState);
    return function (state) {
      return !isPristine(state);
    };
  };
};

export default createIsDirty;