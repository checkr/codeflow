'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _createOnDragStart = require('./createOnDragStart');

var createOnDrop = function createOnDrop(name, change) {
  return function (event) {
    change(event.dataTransfer.getData(_createOnDragStart.dataKey));
    event.preventDefault();
  };
};
exports.default = createOnDrop;