var createOnFocus = function createOnFocus(name, focus) {
  return function () {
    return focus(name);
  };
};
export default createOnFocus;