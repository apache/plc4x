package spi

import "github.com/apache/plc4x/plc4go/spi/transports"

// TransportErrorHandlerSetter exposes the ability to receive notifications when
// the underlying transport reports an error classification.
type TransportErrorHandlerSetter interface {
	SetTransportErrorHandler(handler transports.TransportErrorHandler)
}
