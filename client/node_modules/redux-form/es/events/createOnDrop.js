import { dataKey } from './createOnDragStart';
var createOnDrop = function createOnDrop(name, change) {
  return function (event) {
    change(event.dataTransfer.getData(dataKey));
    event.preventDefault();
  };
};
export default createOnDrop;