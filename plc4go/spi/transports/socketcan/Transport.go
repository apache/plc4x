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

// Package socketcan is a CAN transport over Linux SocketCAN raw sockets, in pure Go (no cgo, no
// libpcap), so it works in CGO_ENABLED=0 builds.
//
// The connection string names the network interface:
//
//	canopen:socketcan://can0
//
// The instance delivers frames in the layout of package can: one 16 byte record per frame.
//
// The kernel delivers error frames (bus off, error passive, lost arbitration, ...) only when asked.
// This option asks for them; its value is a CAN error class mask, as in Linux's linux/can/error.h:
//
//	socketcan.error-mask   e.g. 0x1FFFFFFF for every class, 0 (the default) for none
//
// They arrive as frames with can.Frame.Error set. On
// anything but Linux the transport is still registered, so a connection string naming it fails with
// a clear error rather than "unknown transport".
package socketcan

import (
	"net/url"
	"regexp"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

const (
	// TransportCode is the scheme of a connection string using this transport.
	TransportCode = "socketcan"

	// ErrorMaskOption is the option naming the error classes the kernel delivers as error frames.
	ErrorMaskOption = "socketcan.error-mask"

	// errorClassMask is every error class bit Linux defines (CAN_ERR_MASK).
	errorClassMask = uint32(0x1FFFFFFF)
)

// interfaceNamePattern is what Linux accepts as a network interface name: at most 15 characters,
// no slash and no whitespace.
var interfaceNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.\-:]{1,15}$`)

type Transport struct {
	log zerolog.Logger
}

var _ transports.Transport = (*Transport)(nil)

func NewTransport(_options ...options.WithOption) *Transport {
	return &Transport{
		log: options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

func (m *Transport) GetTransportCode() string {
	return TransportCode
}

func (m *Transport) GetTransportName() string {
	return "Linux SocketCAN Transport"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, transportOptions map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	interfaceName := transportUrl.Host
	if interfaceName == "" {
		// socketcan:can0 (no slashes) parses into Opaque.
		interfaceName = transportUrl.Opaque
	}
	if !interfaceNamePattern.MatchString(interfaceName) {
		return nil, errors.Errorf("socketcan needs a network interface name such as can0, got %q", interfaceName)
	}
	errorMask, err := parseErrorMask(transportOptions[ErrorMaskOption])
	if err != nil {
		return nil, err
	}
	instance := NewTransportInstance(interfaceName, m, _options...)
	instance.ErrorMask = errorMask
	return instance, nil
}

// parseErrorMask reads the error-mask option: absent is 0, otherwise one number, in decimal or with
// a 0x prefix, inside CAN_ERR_MASK.
func parseErrorMask(values []string) (uint32, error) {
	if len(values) == 0 {
		return 0, nil
	}
	if len(values) > 1 {
		return 0, errors.Errorf("%s takes one value, got %d", ErrorMaskOption, len(values))
	}
	mask, err := strconv.ParseUint(values[0], 0, 32)
	if err != nil {
		return 0, errors.Wrapf(err, "%s must be a number such as 0x1FFFFFFF, got %q", ErrorMaskOption, values[0])
	}
	if uint32(mask)&^errorClassMask != 0 {
		return 0, errors.Errorf("%s 0x%X has bits outside the error classes (0x%X)", ErrorMaskOption, mask, errorClassMask)
	}
	return uint32(mask), nil
}

func (m *Transport) Close() error {
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}
