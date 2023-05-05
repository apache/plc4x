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

package options

// WithDiscoveryOption is a marker interface for options regarding discovery
type WithDiscoveryOption interface {
	isDiscoveryOption() bool
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

func WithDiscoveryOptionProtocolSpecific(key string, value any) WithDiscoveryOption {
	return discoveryOptionProtocolSpecific{
		key:   key,
		value: value,
	}
}

func FilterDiscoveryOptionsProtocol(options []WithDiscoveryOption) []DiscoveryOptionProtocol {
	var filtered []DiscoveryOptionProtocol
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionProtocol:
			filtered = append(filtered, option.(DiscoveryOptionProtocol))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsTransport(options []WithDiscoveryOption) []DiscoveryOptionTransport {
	var filtered []DiscoveryOptionTransport
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionTransport:
			filtered = append(filtered, option.(DiscoveryOptionTransport))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsDeviceName(options []WithDiscoveryOption) []DiscoveryOptionDeviceName {
	var filtered []DiscoveryOptionDeviceName
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionDeviceName:
			filtered = append(filtered, option.(DiscoveryOptionDeviceName))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsLocalAddress(options []WithDiscoveryOption) []DiscoveryOptionLocalAddress {
	var filtered []DiscoveryOptionLocalAddress
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionLocalAddress:
			filtered = append(filtered, option.(DiscoveryOptionLocalAddress))
		}
	}
	return filtered
}

func FilterDiscoveryOptionsRemoteAddress(options []WithDiscoveryOption) []DiscoveryOptionRemoteAddress {
	var filtered []DiscoveryOptionRemoteAddress
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionRemoteAddress:
			filtered = append(filtered, option.(DiscoveryOptionRemoteAddress))
		}
	}
	return filtered
}

func FilterDiscoveryOptionProtocolSpecific(options []WithDiscoveryOption) []DiscoveryOptionProtocolSpecific {
	var filtered []DiscoveryOptionProtocolSpecific
	for _, option := range options {
		switch option := option.(type) {
		case discoveryOptionProtocolSpecific:
			filtered = append(filtered, option)
		}
	}
	return filtered
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type discoveryOption struct {
}

func (_ discoveryOption) isDiscoveryOption() bool {
	return true
}

type DiscoveryOptionProtocol interface {
	GetProtocolName() string
}

type discoveryOptionProtocol struct {
	discoveryOption
	protocolName string
}

func (d discoveryOptionProtocol) GetProtocolName() string {
	return d.protocolName
}

type DiscoveryOptionTransport interface {
	GetTransportName() string
}

type discoveryOptionTransport struct {
	discoveryOption
	transportName string
}

func (d discoveryOptionTransport) GetTransportName() string {
	return d.transportName
}

type DiscoveryOptionDeviceName interface {
	GetDeviceName() string
}

type discoveryOptionDeviceName struct {
	discoveryOption
	deviceName string
}

func (d discoveryOptionDeviceName) GetDeviceName() string {
	return d.deviceName
}

type DiscoveryOptionLocalAddress interface {
	GetLocalAddress() string
}

type discoveryOptionLocalAddress struct {
	discoveryOption
	localAddress string
}

func (d discoveryOptionLocalAddress) GetLocalAddress() string {
	return d.localAddress
}

type DiscoveryOptionRemoteAddress interface {
	GetRemoteAddress() string
}

type discoveryOptionRemoteAddress struct {
	discoveryOption
	remoteAddress string
}

func (d discoveryOptionRemoteAddress) GetRemoteAddress() string {
	return d.remoteAddress
}

type DiscoveryOptionProtocolSpecific interface {
	GetKey() string
	GetValue() any
}

type discoveryOptionProtocolSpecific struct {
	discoveryOption
	key   string
	value any
}

func (d discoveryOptionProtocolSpecific) GetKey() string {
	return d.key
}

func (d discoveryOptionProtocolSpecific) GetValue() any {
	return d.value
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////
