package transports

import (
	stdErrors "errors"
	"syscall"
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
	if err == nil {
		return false
	}
	var errno syscall.Errno
	if !stdErrors.As(err, &errno) {
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
