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

package model

import (
	"net/url"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcDiscoveryItem
type DefaultPlcDiscoveryItem struct {
	ProtocolCode  string
	TransportCode string
	TransportUrl  url.URL `ignore:"true"` // TODO: find a way to render this as string (e.g. stringer annotation or something)
	Options       map[string][]string
	Name          string
	Attributes    map[string]apiValues.PlcValue
}

func NewDefaultPlcDiscoveryItem(
	ProtocolCode string,
	TransportCode string,
	TransportUrl url.URL,
	Options map[string][]string,
	Name string,
	Attributes map[string]apiValues.PlcValue,
) apiModel.PlcDiscoveryItem {
	return &DefaultPlcDiscoveryItem{
		ProtocolCode,
		TransportCode,
		TransportUrl,
		Options,
		Name,
		Attributes,
	}
}

func (d *DefaultPlcDiscoveryItem) GetProtocolCode() string {
	return d.ProtocolCode
}

func (d *DefaultPlcDiscoveryItem) GetTransportCode() string {
	return d.TransportCode
}

func (d *DefaultPlcDiscoveryItem) GetTransportUrl() url.URL {
	return d.TransportUrl
}

func (d *DefaultPlcDiscoveryItem) GetOptions() map[string][]string {
	return d.Options
}

func (d *DefaultPlcDiscoveryItem) GetName() string {
	return d.Name
}

func (d *DefaultPlcDiscoveryItem) GetAttributes() map[string]apiValues.PlcValue {
	return d.Attributes
}

func (d *DefaultPlcDiscoveryItem) GetConnectionUrl() string {
	options := ""
	if d.Options != nil {
		var flatOptions []string
		for k, vl := range d.Options {
			for _, v := range vl {
				flatOptions = append(flatOptions, url.QueryEscape(k)+"="+url.QueryEscape(v))
			}
		}
		options += "?" + strings.Join(flatOptions, "&")
	}
	return d.ProtocolCode + ":" + d.TransportCode + "//" + d.TransportUrl.Host + options
}
