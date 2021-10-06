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

import (
	"net/url"
)

type DefaultPlcDiscoveryEvent struct {
	ProtocolCode  string
	TransportCode string
	TransportUrl  url.URL
	Options       map[string][]string
	Name          string
}

func (d *DefaultPlcDiscoveryEvent) GetProtocolCode() string {
	return d.TransportCode
}

func (d *DefaultPlcDiscoveryEvent) GetTransportCode() string {
	return d.TransportCode
}

func (d *DefaultPlcDiscoveryEvent) GetTransportUrl() url.URL {
	return d.TransportUrl
}

func (d *DefaultPlcDiscoveryEvent) GetOptions() map[string][]string {
	return d.Options
}

func (d *DefaultPlcDiscoveryEvent) GetName() string {
	return d.Name
}

func (d *DefaultPlcDiscoveryEvent) GetConnectionString() string {
	if d.Options != nil {
		panic("Not implemented")
	}
	return d.ProtocolCode + ":" + d.TransportCode + "//" + d.TransportUrl.Host
}
