# slug

See the [API docs](http://go.pkgdoc.org/github.com/extemporalgenome/slug).

Latin-ish inputs should have very stable output. All inputs are passed through
an NFKD transform, and anything still in the unicode Letter and Number
categories are passed through intact. Anything in the Mark or Lm/Sk categories
(modifiers) are skipped, and runs of characters from any other categories are
collapsed to a single hyphen.
