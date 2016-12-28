export var dataKey = 'text';
var createOnDragStart = function createOnDragStart(name, value) {
  return function (event) {
    event.dataTransfer.setData(dataKey, value);
  };
};

export default createOnDragStart;