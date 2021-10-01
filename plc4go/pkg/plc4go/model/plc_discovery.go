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

type PlcDiscoveryEvent interface {
	GetProtocolCode() string
	GetTransportCode() string
	GetTransportUrl() url.URL
	GetOptions() map[string][]string
	GetName() string
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

type WithDiscoveryOption interface {
	IsDiscoveryOption() bool
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type discoveryOption struct {
}

func (_ discoveryOption) IsDiscoveryOption() bool {
	return true
}

func FilterDiscoveryOptionsProtocol(options []WithDiscoveryOption) []DiscoveryOptionProtocol {
	castFunc := func(typ interface{}) DiscoveryOptionProtocol {
		if casted, ok := typ.(discoveryOptionProtocol); ok {
			return &casted
		}
		return nil
	}
	var filtered []DiscoveryOptionProtocol
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionProtocol:
			filtered = append(filtered, castFunc(option))
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
	castFunc := func(typ interface{}) DiscoveryOptionTransport {
		if casted, ok := typ.(discoveryOptionTransport); ok {
			return &casted
		}
		return nil
	}
	var filtered []DiscoveryOptionTransport
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionTransport:
			filtered = append(filtered, castFunc(option))
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
	castFunc := func(typ interface{}) DiscoveryOptionDeviceName {
		if casted, ok := typ.(discoveryOptionDeviceName); ok {
			return &casted
		}
		return nil
	}
	var filtered []DiscoveryOptionDeviceName
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionDeviceName:
			filtered = append(filtered, castFunc(option))
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
	castFunc := func(typ interface{}) DiscoveryOptionLocalAddress {
		if casted, ok := typ.(discoveryOptionLocalAddress); ok {
			return &casted
		}
		return nil
	}
	var filtered []DiscoveryOptionLocalAddress
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionLocalAddress:
			filtered = append(filtered, castFunc(option))
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
	castFunc := func(typ interface{}) DiscoveryOptionRemoteAddress {
		if casted, ok := typ.(discoveryOptionRemoteAddress); ok {
			return &casted
		}
		return nil
	}
	var filtered []DiscoveryOptionRemoteAddress
	for _, option := range options {
		switch option.(type) {
		case discoveryOptionRemoteAddress:
			filtered = append(filtered, castFunc(option))
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

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////
