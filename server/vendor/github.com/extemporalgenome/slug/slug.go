// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package slug transforms strings into a normalized form well suited for use in URLs.
package slug

import (
	"golang.org/x/text/unicode/norm"
	"encoding/hex"
	"unicode"
	"unicode/utf8"
)

var lat = []*unicode.RangeTable{unicode.Letter, unicode.Number}
var nop = []*unicode.RangeTable{unicode.Mark, unicode.Sk, unicode.Lm}

// Slug replaces each run of characters which are not unicode letters or
// numbers with a single hyphen, except for leading or trailing runs. Letters
// will be stripped of diacritical marks and lowercased. Letter or number
// codepoints that do not have combining marks or a lower-cased variant will
// be passed through unaltered.
func Slug(s string) string {
	buf := make([]rune, 0, len(s))
	dash := false
	for _, r := range norm.NFKD.String(s) {
		switch {
		// unicode 'letters' like mandarin characters pass through
		case unicode.IsOneOf(lat, r):
			buf = append(buf, unicode.ToLower(r))
			dash = true
		case unicode.IsOneOf(nop, r):
			// skip
		case dash:
			buf = append(buf, '-')
			dash = false
		}
	}
	if i := len(buf) - 1; i >= 0 && buf[i] == '-' {
		buf = buf[:i]
	}
	return string(buf)
}

// SlugAscii is identical to Slug, except that runs of one or more unicode
// letters or numbers that still fall outside the ASCII range will have their
// UTF-8 representation hex encoded and delimited by hyphens. As with Slug, in
// no case will hyphens appear at either end of the returned string.
func SlugAscii(s string) string {
	const m = utf8.UTFMax
	var (
		ib    [m * 3]byte
		ob    []byte
		buf   = make([]byte, 0, len(s))
		dash  = false
		latin = true
	)
	for _, r := range norm.NFKD.String(s) {
		switch {
		case unicode.IsOneOf(lat, r):
			r = unicode.ToLower(r)
			n := utf8.EncodeRune(ib[:m], r)
			if r >= 128 {
				if latin && dash {
					buf = append(buf, '-')
				}
				n = hex.Encode(ib[m:], ib[:n])
				ob = ib[m : m+n]
				latin = false
			} else {
				if !latin {
					buf = append(buf, '-')
				}
				ob = ib[:n]
				latin = true
			}
			dash = true
			buf = append(buf, ob...)
		case unicode.IsOneOf(nop, r):
			// skip
		case dash:
			buf = append(buf, '-')
			dash = false
			latin = true
		}
	}
	if i := len(buf) - 1; i >= 0 && buf[i] == '-' {
		buf = buf[:i]
	}
	return string(buf)
}

// IsSlugAscii returns true only if SlugAscii(s) == s.
func IsSlugAscii(s string) bool {
	dash := true
	for _, r := range s {
		switch {
		case r == '-':
			if dash {
				return false
			}
			dash = true
		case 'a' <= r && r <= 'z', '0' <= r && r <= '9':
			dash = false
		default:
			return false
		}
	}
	return !dash
}
