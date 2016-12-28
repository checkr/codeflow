"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = formatName;
var isFieldArrayRegx = /\[\d+\]$/;

function formatName(context, name) {
  var sectionPrefix = context._reduxForm.sectionPrefix;

  return !sectionPrefix || isFieldArrayRegx.test(name) ? name : sectionPrefix + "." + name;
}