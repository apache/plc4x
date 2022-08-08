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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
)

type DefaultPlcBrowseQueryResult struct {
	Field             model.PlcField
	Name              string
	Readable          bool
	Writable          bool
	Subscribable      bool
	PossibleDataTypes []string
	Attributes        map[string]values.PlcValue
}

func (d *DefaultPlcBrowseQueryResult) GetField() model.PlcField {
	return d.Field
}

func (d *DefaultPlcBrowseQueryResult) GetName() string {
	return d.Name
}

func (d *DefaultPlcBrowseQueryResult) IsReadable() bool {
	return d.Readable
}

func (d *DefaultPlcBrowseQueryResult) IsWritable() bool {
	return d.Writable
}

func (d *DefaultPlcBrowseQueryResult) IsSubscribable() bool {
	return d.Subscribable
}

func (d *DefaultPlcBrowseQueryResult) GetPossibleDataTypes() []string {
	return d.PossibleDataTypes
}

func (d *DefaultPlcBrowseQueryResult) GetAttributes() map[string]values.PlcValue {
	return d.Attributes
}

func (d *DefaultPlcBrowseQueryResult) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcBrowseQueryResult"); err != nil {
		return err
	}

	var fieldAsString string
	fieldAsString = fmt.Sprintf("%s", d.Field)
	if err := writeBuffer.WriteString("field", uint32(len(fieldAsString)*8), "UTF-8", fieldAsString); err != nil {
		return err
	}
	if err := writeBuffer.WriteString("name", uint32(len(d.Name)*8), "UTF-8", d.Name); err != nil {
		return err
	}
	if err := writeBuffer.WriteBit("readable", d.Readable); err != nil {
		return err
	}
	if err := writeBuffer.WriteBit("writable", d.Writable); err != nil {
		return err
	}
	if err := writeBuffer.WriteBit("subscribable", d.Subscribable); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("possibleDataTypes", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for _, dataType := range d.PossibleDataTypes {
		if err := writeBuffer.WriteString("", uint32(len(dataType)*8), "UTF-8", dataType); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("possibleDataTypes", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("attributes", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for name, plcValue := range d.Attributes {
		if serializable, ok := plcValue.(utils.Serializable); ok {
			if err := writeBuffer.PushContext(name); err != nil {
				return err
			}
			if err := serializable.Serialize(writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext(name); err != nil {
				return err
			}
		} else {
			plcValueAsString := fmt.Sprintf("%v", plcValue)
			if err := writeBuffer.WriteString(name, uint32(len(plcValueAsString)*8), "UTF-8", plcValueAsString); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PopContext("attributes", utils.WithRenderAsList(true)); err != nil {
		return err
	}

	if err := writeBuffer.PopContext("PlcBrowseQueryResult"); err != nil {
		return err
	}
	return nil
}
