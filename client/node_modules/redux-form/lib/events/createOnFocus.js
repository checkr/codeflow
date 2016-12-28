"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
var createOnFocus = function createOnFocus(name, focus) {
  return function () {
    return focus(name);
  };
};
exports.default = createOnFocus;