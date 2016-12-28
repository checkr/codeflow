'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
var dataKey = exports.dataKey = 'text';
var createOnDragStart = function createOnDragStart(name, value) {
  return function (event) {
    event.dataTransfer.setData(dataKey, value);
  };
};

exports.default = createOnDragStart;