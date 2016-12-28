var createGetFormSyncErrors = function createGetFormSyncErrors(_ref) {
  var getIn = _ref.getIn;
  return function (form) {
    var getFormState = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : function (state) {
      return getIn(state, 'form');
    };
    return function (state) {
      return getIn(getFormState(state), form + '.syncErrors');
    };
  };
};

export default createGetFormSyncErrors;