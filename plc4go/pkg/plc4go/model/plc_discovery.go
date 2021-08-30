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

package model

import "net/url"

type PlcDiscoveryEvent struct {
	ProtocolCode  string
	TransportCode string
	TransportUrl  url.URL
	Options       map[string][]string
	Name          string
}

func NewPlcDiscoveryEvent(protocolCode string, transportCode string, transportUrl url.URL, options map[string][]string, name string) PlcDiscoveryEvent {
	return PlcDiscoveryEvent{
		ProtocolCode:  protocolCode,
		TransportCode: transportCode,
		TransportUrl:  transportUrl,
		Options:       options,
		Name:          name,
	}
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

/////////////////////////////////////////////////////////////////////////////////////////7
// Internal
/////////////////////////////////////////////////////////////////////////////////////////7

type WithDiscoveryOption interface {
	IsDiscoveryOption() bool
}

type discoveryOption struct {
}

func (_ discoveryOption) IsDiscoveryOption() bool {
	return true
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

type DiscoveryOptionProtocol interface {
	GetProtocolName() string
}

type discoveryOptionProtocol struct {
	discoveryOption
	protocolName string
}

func (d *discoveryOptionProtocol) GetProtocolName() string {
	return d.protocolName
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

type DiscoveryOptionTransport interface {
	GetTransportName() string
}

type discoveryOptionTransport struct {
	discoveryOption
	transportName string
}

func (d *discoveryOptionTransport) GetTransportName() string {
	return d.transportName
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

type DiscoveryOptionDeviceName interface {
	GetDeviceName() string
}

type discoveryOptionDeviceName struct {
	discoveryOption
	deviceName string
}

func (d *discoveryOptionDeviceName) GetDeviceName() string {
	return d.deviceName
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

type DiscoveryOptionLocalAddress interface {
	GetLocalAddress() string
}

type discoveryOptionLocalAddress struct {
	discoveryOption
	localAddress string
}

func (d *discoveryOptionLocalAddress) GetLocalAddress() string {
	return d.localAddress
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

type DiscoveryOptionRemoteAddress interface {
	GetRemoteAddress() string
}

type discoveryOptionRemoteAddress struct {
	discoveryOption
	remoteAddress string
}

func (d *discoveryOptionRemoteAddress) GetRemoteAddress() string {
	return d.remoteAddress
}
