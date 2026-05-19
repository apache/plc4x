/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package errors provides error-handling primitives compatible with the now-archived
// github.com/pkg/errors. It exposes the same API surface (New, Errorf, Wrap, Wrapf,
// WithStack, WithMessage, WithMessagef, Cause, Is, As, Unwrap) so callers can swap
// imports without behavioural changes:
//
//   - New / Errorf attach a stack trace at the point of creation.
//   - Wrap / Wrapf prepend a message AND attach a stack trace.
//   - WithStack adds a stack trace without changing the message.
//   - WithMessage / WithMessagef prepend a message without a stack trace.
//   - Cause walks the chain via the (unexported) causer interface, matching pkg/errors.
//   - Is, As, Unwrap delegate to the stdlib errors package for chain interop.
//
// Format verbs mirror pkg/errors: %s/%v render the message; %+v also dumps the stack.
package errors

import (
	stderrors "errors"
	"fmt"
	"io"
)

// New returns an error with the given message and a stack trace captured at the call site.
func New(message string) error {
	return &fundamental{msg: message, stack: callers()}
}

// Errorf returns a formatted error with a stack trace captured at the call site.
func Errorf(format string, args ...any) error {
	return &fundamental{msg: fmt.Sprintf(format, args...), stack: callers()}
}

// fundamental is an error with a message and a stack, but no underlying cause.
type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, f.msg)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", f.msg)
	}
}

// WithStack annotates err with a stack trace. Returns nil if err is nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{error: err, stack: callers()}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error  { return w.error }
func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// Wrap returns an error annotating err with a message and a stack trace.
// Returns nil if err is nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withStack{
		error: &withMessage{cause: err, msg: message},
		stack: callers(),
	}
}

// Wrapf returns an error annotating err with a formatted message and a stack trace.
// Returns nil if err is nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &withStack{
		error: &withMessage{cause: err, msg: fmt.Sprintf(format, args...)},
		stack: callers(),
	}
}

// WithMessage annotates err with a message but no stack trace. Returns nil if err is nil.
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{cause: err, msg: message}
}

// WithMessagef annotates err with a formatted message but no stack trace.
// Returns nil if err is nil.
func WithMessagef(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &withMessage{cause: err, msg: fmt.Sprintf(format, args...)}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}

// Cause walks err's chain via the causer interface (Cause() error) and returns
// the topmost error that does not implement it. Returns nil if err is nil.
func Cause(err error) error {
	type causer interface{ Cause() error }
	for err != nil {
		c, ok := err.(causer)
		if !ok {
			break
		}
		err = c.Cause()
	}
	return err
}

// Is reports whether any error in err's chain matches target. Delegates to stdlib errors.Is.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches target. Delegates to stdlib errors.As.
func As(err error, target any) bool { return stderrors.As(err, target) }

// Unwrap returns err's wrapped error, or nil if err does not wrap. Delegates to stdlib errors.Unwrap.
func Unwrap(err error) error { return stderrors.Unwrap(err) }
