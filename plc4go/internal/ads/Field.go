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

package ads

import (
	"encoding/xml"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	adsModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

const NONE = int32(-1)

type PlcField struct {
	model.PlcField

	NumElements  int32
	StartElement int32
	EndElement   int32
}

func needsResolving(plcField model.PlcField) bool {
	switch plcField.(type) {
	case SymbolicPlcField:
		return true
	case DirectPlcField:
		return false
	default:
		return false
	}
}

type DirectPlcField struct {
	PlcField

	IndexGroup   uint32
	IndexOffset  uint32
	AdsDatatype  adsModel.AdsDataType
	StringLength int32
}

func newDirectAdsPlcField(indexGroup uint32, indexOffset uint32, adsDatatype adsModel.AdsDataType, stringLength int32, numElements int32, startElement int32, endElement int32) (model.PlcField, error) {
	return DirectPlcField{
		IndexGroup:   indexGroup,
		IndexOffset:  indexOffset,
		AdsDatatype:  adsDatatype,
		StringLength: stringLength,
		PlcField: PlcField{
			NumElements:  numElements,
			StartElement: startElement,
			EndElement:   endElement,
		},
	}, nil
}

func castToDirectAdsFieldFromPlcField(plcField model.PlcField) (DirectPlcField, error) {
	if adsField, ok := plcField.(DirectPlcField); ok {
		return adsField, nil
	}
	return DirectPlcField{}, errors.Errorf("couldn't %T cast to DirectPlcField", plcField)
}

func (m DirectPlcField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("DirectPlcField"); err != nil {
		return err
	}

	if err := writeBuffer.WriteUint32("indexGroup", 32, m.IndexGroup); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint32("indexOffset", 32, m.IndexOffset); err != nil {
		return err
	}
	if err := writeBuffer.WriteString("adsDatatypeName", uint32(len([]rune(m.AdsDatatype.String()))*8), "UTF-8", m.AdsDatatype.String()); err != nil {
		return err
	}
	if (m.AdsDatatype == adsModel.AdsDataType_STRING || m.AdsDatatype == adsModel.AdsDataType_WSTRING) && (m.StringLength != NONE) {
		if err := writeBuffer.WriteInt32("stringLength", 32, m.StringLength); err != nil {
			return err
		}
	}
	if m.NumElements != NONE {
		if err := writeBuffer.WriteInt32("numElements", 32, m.NumElements); err != nil {
			return err
		}
	}
	if m.StartElement != NONE && m.EndElement != NONE {
		if err := writeBuffer.WriteInt32("startElement", 32, m.StartElement); err != nil {
			return err
		}
		if err := writeBuffer.WriteInt32("endElement", 32, m.EndElement); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext("DirectPlcField"); err != nil {
		return err
	}
	return nil
}

func (m DirectPlcField) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	panic(name)
}

type SymbolicPlcField struct {
	PlcField

	SymbolicAddress string
}

func newAdsSymbolicPlcField(symbolicAddress string, numElements int32, startElement int32, endElement int32) (model.PlcField, error) {
	return SymbolicPlcField{
		SymbolicAddress: symbolicAddress,
		PlcField: PlcField{
			NumElements:  numElements,
			StartElement: startElement,
			EndElement:   endElement,
		},
	}, nil
}

func castToSymbolicPlcFieldFromPlcField(plcField model.PlcField) (SymbolicPlcField, error) {
	if adsField, ok := plcField.(SymbolicPlcField); ok {
		return adsField, nil
	}
	return SymbolicPlcField{}, errors.Errorf("couldn't cast %T to SymbolicPlcField", plcField)
}

func (m SymbolicPlcField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("SymbolicPlcField"); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("symbolicAddress", uint32(len([]rune(m.SymbolicAddress))*8), "UTF-8", m.SymbolicAddress); err != nil {
		return err
	}
	if m.NumElements != NONE {
		if err := writeBuffer.WriteInt32("numElements", 32, m.NumElements); err != nil {
			return err
		}
	}
	if m.StartElement != NONE && m.EndElement != NONE {
		if err := writeBuffer.WriteInt32("startElement", 32, m.StartElement); err != nil {
			return err
		}
		if err := writeBuffer.WriteInt32("endElement", 32, m.EndElement); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext("SymbolicPlcField"); err != nil {
		return err
	}
	return nil
}
