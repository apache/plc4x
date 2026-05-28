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

package transports

import (
	stdErrors "errors"
	"fmt"
	"syscall"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// TransportErrorKind represents the transport-level severity of an error.
type TransportErrorKind int

const (
	TransportErrorUnknown TransportErrorKind = iota
	// TransportErrorTransient indicates a transient issue (e.g. temporary IO hiccup) that may succeed on retry without reconnect.
	TransportErrorTransient
	// TransportErrorRetryable indicates the caller should retry the operation after resetting or reconnecting.
	TransportErrorRetryable
	// TransportErrorFatal indicates the transport connection can no longer be used.
	TransportErrorFatal
)

// TransportErrorHandler is invoked when callers classify a transport error.
type TransportErrorHandler func(kind TransportErrorKind, err error)

// TransportError wraps an underlying error with its classified transport severity.
type TransportError struct {
	kind TransportErrorKind
	err  error
}

// NewTransportError creates a new TransportError with the given kind and cause.
func NewTransportError(kind TransportErrorKind, err error) error {
	if utils.IsNil(err) {
		return nil
	}
	var existing *TransportError
	if ErrorAs(err, &existing) {
		return err
	}
	return &TransportError{kind: kind, err: err}
}

// AsTransportError retrieves a TransportError from the provided error chain.
func AsTransportError(err error) (*TransportError, bool) {
	if utils.IsNil(err) {
		return nil, false
	}
	var transportErr *TransportError
	if ErrorAs(err, &transportErr) {
		return transportErr, true
	}
	return nil, false
}

// Error implements the error interface.
func (t *TransportError) Error() string {
	if t.err == nil {
		return fmt.Sprintf("transport error (%s)", t.kind.String())
	}
	return fmt.Sprintf("transport error (%s): %v", t.kind.String(), t.err)
}

// Unwrap exposes the underlying cause.
func (t *TransportError) Unwrap() error {
	return t.err
}

// Kind reports the TransportErrorKind associated with the error.
func (t *TransportError) Kind() TransportErrorKind {
	return t.kind
}

// ErrorIs mirrors errors.Is but guards against panics triggered by improperly constructed error values.
func ErrorIs(err error, target error) (matched bool) {
	if err == nil || target == nil || utils.IsNil(err) {
		return false
	}
	defer func() {
		if r := recover(); r != nil {
			matched = false
		}
	}()
	return stdErrors.Is(err, target)
}

// ErrorAs mirrors errors.As but guards against panics triggered by improperly constructed error values.
func ErrorAs(err error, target any) (matched bool) {
	if err == nil || target == nil || utils.IsNil(err) {
		return false
	}
	defer func() {
		if r := recover(); r != nil {
			matched = false
		}
	}()
	return stdErrors.As(err, target)
}

// IsFatal reports whether the error kind signals an unusable transport.
func (k TransportErrorKind) IsFatal() bool {
	return k == TransportErrorFatal
}

// IsRetryable reports whether the operation may succeed if repeated.
func (k TransportErrorKind) IsRetryable() bool {
	return k == TransportErrorTransient || k == TransportErrorRetryable
}

// IsTransientSyscallError checks for errno values that are commonly treated as transient.
func IsTransientSyscallError(err error) bool {
	if err == nil || utils.IsNil(err) {
		return false
	}
	var errno syscall.Errno
	if !ErrorAs(err, &errno) {
		return false
	}
	switch errno {
	case syscall.EAGAIN,
		syscall.EINTR,
		syscall.EINPROGRESS,
		syscall.EALREADY,
		syscall.ENOBUFS:
		return true
	}
	switch errno {
	case syscall.Errno(10004), // WSAEINTR
		syscall.Errno(10035), // WSAEWOULDBLOCK
		syscall.Errno(10036), // WSAEINPROGRESS
		syscall.Errno(10037), // WSAEALREADY
		syscall.Errno(10055): // WSAENOBUFS
		return true
	}
	return false
}

// String returns a human readable representation of the error kind.
func (k TransportErrorKind) String() string {
	switch k {
	case TransportErrorUnknown:
		return "unknown"
	case TransportErrorTransient:
		return "transient"
	case TransportErrorRetryable:
		return "retryable"
	case TransportErrorFatal:
		return "fatal"
	default:
		return "invalid"
	}
}
