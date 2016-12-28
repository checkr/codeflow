var isFieldArrayRegx = /\[\d+\]$/;

export default function formatName(context, name) {
  var sectionPrefix = context._reduxForm.sectionPrefix;

  return !sectionPrefix || isFieldArrayRegx.test(name) ? name : sectionPrefix + "." + name;
}