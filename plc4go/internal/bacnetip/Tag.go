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

package bacnetip

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type BacNetPlcTag interface {
	apiModel.PlcTag

	GetObjectId() objectId
	GetProperties() []property
}

type plcTag struct {
	ObjectId objectId
	// Properties 1..N identifiers
	Properties []property
}

type objectId struct {
	// ObjectIdType defines the object type. If not defined ObjectIdTypeProprietary must be defined
	ObjectIdType *readWriteModel.BACnetObjectType
	// ObjectIdTypeProprietary is only defined if ObjectIdType is not defined
	ObjectIdTypeProprietary *uint16
	// ObjectIdInstance is the instance of this object
	ObjectIdInstance uint32
}

func (o objectId) getId() uint16 {
	if o.ObjectIdType != nil {
		return uint16(*o.ObjectIdType)
	} else {
		return *o.ObjectIdTypeProprietary
	}
}

func (o objectId) String() string {
	var result string
	if o.ObjectIdType != nil {
		result += fmt.Sprint(*o.ObjectIdType)
	} else {
		result += fmt.Sprint(*o.ObjectIdTypeProprietary)
	}
	result += fmt.Sprintf(":%d", o.ObjectIdInstance)
	return result
}

type property struct {
	// PropertyIdentifier defines the property type. If not defined PropertyIdentifierProprietary must be defined
	PropertyIdentifier *readWriteModel.BACnetPropertyIdentifier
	// PropertyIdentifierProprietary is only defined if PropertyIdentifier is not defined
	PropertyIdentifierProprietary *uint32
	// ArrayIndex Optional index of property
	ArrayIndex *uint
}

func (p property) getId() uint32 {
	if p.PropertyIdentifier != nil {
		return uint32(*p.PropertyIdentifier)
	} else {
		return *p.PropertyIdentifierProprietary
	}
}

func (p property) String() string {
	var result string
	if p.PropertyIdentifier != nil {
		result += fmt.Sprint(*p.PropertyIdentifier)
	} else {
		result += fmt.Sprint(*p.PropertyIdentifierProprietary)
	}
	if p.ArrayIndex != nil {
		result += fmt.Sprintf(":[%d]", p.ArrayIndex)
	}
	return result
}

func (m plcTag) GetAddressString() string {
	var properties []string
	for _, p := range m.Properties {
		properties = append(properties, fmt.Sprint(p))
	}
	propertiesString := strings.Join(properties, "&")
	return fmt.Sprintf("%v/%s", m.ObjectId, propertiesString)
}

func (m plcTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (m plcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (m plcTag) GetObjectId() objectId {
	return m.ObjectId
}

func (m plcTag) GetProperties() []property {
	return m.Properties
}

func (m plcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("BacNetPlcTag"); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("objectId", uint32(len([]rune(m.ObjectId.String()))*8), "UTF-8", m.ObjectId.String()); err != nil {
		return err
	}

	if err := writeBuffer.PushContext("properties"); err != nil {
		return err
	}
	for _, p := range m.Properties {
		if err := writeBuffer.WriteString("property", uint32(len([]rune(p.String()))*8), "UTF-8", p.String()); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("properties"); err != nil {
		return err
	}

	if err := writeBuffer.PopContext("BacNetPlcTag"); err != nil {
		return err
	}
	return nil
}

func (m plcTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
