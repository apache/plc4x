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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcBrowseItem
type DefaultPlcBrowseItem struct {
	Tag          apiModel.PlcTag
	Name         string
	DataTypeName string
	Readable     bool
	Writable     bool
	Subscribable bool
	Children     map[string]apiModel.PlcBrowseItem
	Options      map[string]apiValues.PlcValue
}

func NewDefaultPlcBrowseItem(Tag apiModel.PlcTag, Name string, DataTypeName string, Readable bool, Writable bool, Subscribable bool, Children map[string]apiModel.PlcBrowseItem, Options map[string]apiValues.PlcValue) apiModel.PlcBrowseItem {
	return &DefaultPlcBrowseItem{
		Tag:          Tag,
		Name:         Name,
		DataTypeName: DataTypeName,
		Readable:     Readable,
		Writable:     Writable,
		Subscribable: Subscribable,
		Children:     Children,
		Options:      Options,
	}
}

func (d *DefaultPlcBrowseItem) GetTag() apiModel.PlcTag {
	return d.Tag
}

func (d *DefaultPlcBrowseItem) GetName() string {
	return d.Name
}

func (d *DefaultPlcBrowseItem) GetDataTypeName() string {
	return d.DataTypeName
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

func (d *DefaultPlcBrowseItem) GetChildren() map[string]apiModel.PlcBrowseItem {
	return d.Children
}

func (d *DefaultPlcBrowseItem) GetOptions() map[string]apiValues.PlcValue {
	return d.Options
}
