JWT Middleware for Go-Json-Rest
==================================

[![Build Status](https://travis-ci.org/StephanDollberg/go-json-rest-middleware-jwt.svg?branch=master)](https://travis-ci.org/StephanDollberg/go-json-rest-middleware-jwt) [![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/StephanDollberg/go-json-rest-middleware-jwt) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/StephanDollberg/go-json-rest-middleware-jwt/master/LICENSE)

This is a middleware for [Go-Json-Rest](https://github.com/ant0ine/go-json-rest).

It uses [jwt-go](https://github.com/dgrijalva/jwt-go) to provide a jwt authentication middleware. It provides additional handler functions to provide the login api that will generate the token and an additional refresh handler that can be used to refresh tokens.

An example can be found in the [Go-Json-Rest Examples](https://github.com/ant0ine/go-json-rest-examples/tree/master/jwt) repo.
