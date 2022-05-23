/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package options

// WithDiscoveryOption is a marker interface for options regarding discovery
type WithDiscoveryOption interface {
	isDiscoveryOption() bool
}

// DiscoveryOption is a marker struct which can be used to mark an option
type DiscoveryOption struct {
}

func (_ DiscoveryOption) isDiscoveryOption() bool {
	return true
}

func WithDiscoveryOptionProtocol(protocolName string) WithDiscoveryOption {
	return discoveryOptionProtocol{
		protocolName: protocolName,
	}
}

func WithDiscoveryOptionTransport(transportName string) WithDiscoveryOption {
	return discoveryOptionTransport{
		transportName: transportName,
	}
}

func WithDiscoveryOptionDeviceName(deviceName string) WithDiscoveryOption {
	return discoveryOptionDeviceName{
		deviceName: deviceName,
	}
}

func WithDiscoveryOptionLocalAddress(localAddress string) WithDiscoveryOption {
	return discoveryOptionLocalAddress{
		localAddress: localAddress,
	}
}

func WithDiscoveryOptionRemoteAddress(remoteAddress string) WithDiscoveryOption {
	return discoveryOptionRemoteAddress{
		remoteAddress: remoteAddress,
	}
}

func FilterDiscoveryOptionsProtocol(options []WithDiscoveryOption) []DiscoveryOptionProtocol {
	var filtered []DiscoveryOptionProtocol
	for _, option := range options {
		switch option.(type) {
		case DiscoveryOptionProtocol:
			filtered = append(filtered, option.(DiscoveryOptionProtocol))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsTransport(options []WithDiscoveryOption) []DiscoveryOptionTransport {
	var filtered []DiscoveryOptionTransport
	for _, option := range options {
		switch option.(type) {
		case DiscoveryOptionTransport:
			filtered = append(filtered, option.(DiscoveryOptionTransport))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsDeviceName(options []WithDiscoveryOption) []DiscoveryOptionDeviceName {
	var filtered []DiscoveryOptionDeviceName
	for _, option := range options {
		switch option.(type) {
		case DiscoveryOptionDeviceName:
			filtered = append(filtered, option.(DiscoveryOptionDeviceName))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsLocalAddress(options []WithDiscoveryOption) []DiscoveryOptionLocalAddress {
	var filtered []DiscoveryOptionLocalAddress
	for _, option := range options {
		switch option.(type) {
		case DiscoveryOptionLocalAddress:
			filtered = append(filtered, option.(DiscoveryOptionLocalAddress))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsRemoteAddress(options []WithDiscoveryOption) []DiscoveryOptionRemoteAddress {
	var filtered []DiscoveryOptionRemoteAddress
	for _, option := range options {
		switch option.(type) {
		case DiscoveryOptionRemoteAddress:
			filtered = append(filtered, option.(DiscoveryOptionRemoteAddress))
		}
	}
	return filtered
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type DiscoveryOptionProtocol interface {
	GetProtocolName() string
}

type discoveryOptionProtocol struct {
	DiscoveryOption
	protocolName string
}

func (d *discoveryOptionProtocol) GetProtocolName() string {
	return d.protocolName
}

type DiscoveryOptionTransport interface {
	GetTransportName() string
}

type discoveryOptionTransport struct {
	DiscoveryOption
	transportName string
}

func (d *discoveryOptionTransport) GetTransportName() string {
	return d.transportName
}

type DiscoveryOptionDeviceName interface {
	GetDeviceName() string
}

type discoveryOptionDeviceName struct {
	DiscoveryOption
	deviceName string
}

func (d *discoveryOptionDeviceName) GetDeviceName() string {
	return d.deviceName
}

type DiscoveryOptionLocalAddress interface {
	GetLocalAddress() string
}

type discoveryOptionLocalAddress struct {
	DiscoveryOption
	localAddress string
}

func (d *discoveryOptionLocalAddress) GetLocalAddress() string {
	return d.localAddress
}

type DiscoveryOptionRemoteAddress interface {
	GetRemoteAddress() string
}

type discoveryOptionRemoteAddress struct {
	DiscoveryOption
	remoteAddress string
}

func (d *discoveryOptionRemoteAddress) GetRemoteAddress() string {
	return d.remoteAddress
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////
