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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcBrowseItem
type DefaultPlcBrowseItem struct {
	Field        model.PlcField
	Address      string
	Name         string
	PlcValueType values.PlcValueType
	ArrayInfo    []model.PlcBrowseItemArrayInfo
	Readable     bool
	Writable     bool
	Subscribable bool
	Children     map[string]model.PlcBrowseItem
	Options      map[string]values.PlcValue
}

func (d *DefaultPlcBrowseItem) GetField() model.PlcField {
	return d.Field
}

func (d *DefaultPlcBrowseItem) GetAddress() string {
	return d.Address
}

func (d *DefaultPlcBrowseItem) GetName() string {
	return d.Name
}

func (d *DefaultPlcBrowseItem) GetPlcValueType() values.PlcValueType {
	return d.PlcValueType
}

func (d *DefaultPlcBrowseItem) GetArrayInfo() []model.PlcBrowseItemArrayInfo {
	return d.ArrayInfo
}

func (d *DefaultPlcBrowseItem) IsReadable() bool {
	return d.Readable
}

func (d *DefaultPlcBrowseItem) IsWritable() bool {
	return d.Writable
}

func (d *DefaultPlcBrowseItem) IsSubscribable() bool {
	return d.Subscribable
}

func (d *DefaultPlcBrowseItem) GetChildren() map[string]model.PlcBrowseItem {
	return d.Children
}

func (d *DefaultPlcBrowseItem) GetOptions() map[string]values.PlcValue {
	return d.Options
}
