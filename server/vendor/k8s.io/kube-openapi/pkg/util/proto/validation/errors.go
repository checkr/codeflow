/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"fmt"
)

type errors struct {
	errors []error
}

func (e *errors) Errors() []error {
	return e.errors
}

func (e *errors) AppendErrors(err ...error) {
	e.errors = append(e.errors, err...)
}

type validationError struct {
	Path string
	Err  error
}

func (e validationError) Error() string {
	return fmt.Sprintf("validationError(%s): %v", e.Path, e.Err)
}

type invalidTypeError struct {
	Path     string
	Expected string
	Actual   string
}

func (e invalidTypeError) Error() string {
	return fmt.Sprintf("invalid type for %s: got %q, expected %q", e.Path, e.Actual, e.Expected)
}

type missingRequiredFieldError struct {
	Path  string
	Field string
}

func (e missingRequiredFieldError) Error() string {
	return fmt.Sprintf("missing required field %q in %s", e.Field, e.Path)
}

type unknownFieldError struct {
	Path  string
	Field string
}

func (e unknownFieldError) Error() string {
	return fmt.Sprintf("unknown field %q in %s", e.Field, e.Path)
}

type invalidObjectTypeError struct {
	Path string
	Type string
}

func (e invalidObjectTypeError) Error() string {
	return fmt.Sprintf("unknown object type %q in %s", e.Type, e.Path)
}
