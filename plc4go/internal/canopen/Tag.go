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

package canopen

import (
	"context"
	"fmt"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// TagKind is what part of CANopen a tag addresses.
type TagKind uint8

const (
	// TagKindSDO is one object dictionary entry of one node, read and written with SDO transfers.
	TagKindSDO TagKind = iota
	// TagKindPDO is one process data object, written by sending it and subscribed to by listening.
	TagKindPDO
	// TagKindNMT is network management: writing sends a state command, subscribing reports the
	// commands seen on the bus.
	TagKindNMT
	// TagKindHeartbeat is the heartbeat (and boot-up) messages nodes produce.
	TagKindHeartbeat
)

func (k TagKind) String() string {
	switch k {
	case TagKindSDO:
		return "SDO"
	case TagKindPDO:
		return "PDO"
	case TagKindNMT:
		return "NMT"
	case TagKindHeartbeat:
		return "HEARTBEAT"
	}
	return fmt.Sprintf("TagKind(%d)", uint8(k))
}

// Tag is a parsed CANopen address. The syntax is plc4j's:
//
//	SDO:<node>[/<answerNode>]:<index>/<subindex>:<type>     e.g. SDO:1:0x1018/1:UNSIGNED32
//	<RECEIVE|TRANSMIT>_PDO_<1-4>:<node>:<type>              e.g. TRANSMIT_PDO_1:5:UNSIGNED16
//	NMT  |  NMT:<node>
//	HEARTBEAT  |  HEARTBEAT:<node>
//
// Index and subindex are decimal, or hexadecimal with a 0x prefix. <type> is a CANOpenDataType name
// (UNSIGNED8, INTEGER32, REAL32, VISIBLE_STRING, ...). A node of 0, or none, on NMT and HEARTBEAT
// means every node.
//
// A Tag is also an apiModel.PlcSubscriptionTag so it can be handed to the subscription request
// builder directly.
type Tag struct {
	Kind TagKind
	// Node is the node the tag addresses, 0 for every node (NMT and HEARTBEAT only).
	Node uint8
	// AnswerNode is the node expected to answer an SDO request, when it is not Node itself.
	AnswerNode uint8
	Index      uint16
	SubIndex   uint8
	DataType   readWriteModel.CANOpenDataType
	// Service is the PDO a PDO tag addresses.
	Service readWriteModel.CANOpenService
}

var (
	_ apiModel.PlcSubscriptionTag = Tag{}
	_ utils.Serializable          = Tag{}
)

// IsWildcard reports whether an NMT or HEARTBEAT tag covers every node.
func (t Tag) IsWildcard() bool {
	return (t.Kind == TagKindNMT || t.Kind == TagKindHeartbeat) && t.Node == 0
}

// ResponderNode is the node an SDO transfer of this tag expects its answers from.
func (t Tag) ResponderNode() uint8 {
	if t.AnswerNode != 0 {
		return t.AnswerNode
	}
	return t.Node
}

// GetAddressString spells the tag so that the tag handler parses it back into the same tag.
func (t Tag) GetAddressString() string {
	switch t.Kind {
	case TagKindSDO:
		node := fmt.Sprintf("%d", t.Node)
		if t.AnswerNode != 0 && t.AnswerNode != t.Node {
			node = fmt.Sprintf("%d/%d", t.Node, t.AnswerNode)
		}
		return fmt.Sprintf("SDO:%s:0x%04X/%d:%s", node, t.Index, t.SubIndex, t.DataType.PLC4XEnumName())
	case TagKindPDO:
		return fmt.Sprintf("%s:%d:%s", t.Service.PLC4XEnumName(), t.Node, t.DataType.PLC4XEnumName())
	case TagKindNMT, TagKindHeartbeat:
		if t.Node == 0 {
			return t.Kind.String()
		}
		return fmt.Sprintf("%s:%d", t.Kind, t.Node)
	}
	return ""
}

func (t Tag) GetValueType() apiValues.PlcValueType {
	switch t.Kind {
	case TagKindSDO, TagKindPDO:
		return valueTypeOf(t.DataType)
	case TagKindNMT:
		return apiValues.USINT
	case TagKindHeartbeat:
		return apiValues.Struct
	}
	return apiValues.NULL
}

// valueTypeOf is the value type a data type's values are written as and decode into. It follows the
// data item spec rather than the enum's plcValueName, which disagrees with it for the odd-width
// integers (named RAW_BYTE_ARRAY, decoded as UDINT/ULINT/DINT/LINT) and for VISIBLE_STRING (named
// CHAR, decoded as STRING) - a write built from the plcValueName would not serialize.
func valueTypeOf(dataType readWriteModel.CANOpenDataType) apiValues.PlcValueType {
	switch dataType {
	case readWriteModel.CANOpenDataType_BOOLEAN:
		// plcValueName says BIT, which is not a plc4go value type.
		return apiValues.BOOL
	case readWriteModel.CANOpenDataType_UNSIGNED24:
		return apiValues.UDINT
	case readWriteModel.CANOpenDataType_UNSIGNED40, readWriteModel.CANOpenDataType_UNSIGNED48, readWriteModel.CANOpenDataType_UNSIGNED56:
		return apiValues.ULINT
	case readWriteModel.CANOpenDataType_INTEGER24:
		return apiValues.DINT
	case readWriteModel.CANOpenDataType_INTEGER40, readWriteModel.CANOpenDataType_INTEGER48, readWriteModel.CANOpenDataType_INTEGER56:
		return apiValues.LINT
	case readWriteModel.CANOpenDataType_VISIBLE_STRING, readWriteModel.CANOpenDataType_OCTET_STRING, readWriteModel.CANOpenDataType_UNICODE_STRING:
		return apiValues.STRING
	case readWriteModel.CANOpenDataType_RECORD:
		return apiValues.RAW_BYTE_ARRAY
	}
	if valueType, ok := apiValues.PlcValueTypeByName(dataType.PlcValueName()); ok {
		return valueType
	}
	return apiValues.RAW_BYTE_ARRAY
}

func (t Tag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

// GetPlcSubscriptionType is event: CANopen pushes, nothing is polled.
func (t Tag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionEvent
}

func (t Tag) GetDuration() time.Duration {
	return 0
}

func (t Tag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := t.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

// SerializeWithWriteBuffer renders the tag wrapped as PlcTagItem/tag/CANOpen<Kind>Tag, the way the
// other drivers wrap a tag inside a request.
func (t Tag) SerializeWithWriteBuffer(_ context.Context, wb utils.WriteBuffer) error {
	writeString := func(name, value string) error {
		return wb.WriteString(name, uint32(len(value)*8), value, utils.WithEncoding("UTF-8"))
	}
	typeName := map[TagKind]string{
		TagKindSDO:       "CANOpenSDOTag",
		TagKindPDO:       "CANOpenPDOTag",
		TagKindNMT:       "CANOpenNMTTag",
		TagKindHeartbeat: "CANOpenHeartbeatTag",
	}[t.Kind]
	for _, context := range []string{"PlcTagItem", "tag", typeName} {
		if err := wb.PushContext(context); err != nil {
			return err
		}
	}
	if err := wb.WriteUint8("node", 8, t.Node); err != nil {
		return err
	}
	switch t.Kind {
	case TagKindSDO:
		if err := wb.WriteUint8("answerNode", 8, t.ResponderNode()); err != nil {
			return err
		}
		if err := wb.WriteUint16("index", 16, t.Index); err != nil {
			return err
		}
		if err := wb.WriteUint8("subIndex", 8, t.SubIndex); err != nil {
			return err
		}
		if err := writeString("dataType", t.DataType.PLC4XEnumName()); err != nil {
			return err
		}
	case TagKindPDO:
		if err := writeString("service", t.Service.PLC4XEnumName()); err != nil {
			return err
		}
		if err := writeString("dataType", t.DataType.PLC4XEnumName()); err != nil {
			return err
		}
	}
	for _, context := range []string{typeName, "tag", "PlcTagItem"} {
		if err := wb.PopContext(context); err != nil {
			return err
		}
	}
	return nil
}

func (t Tag) String() string {
	return "canopen.Tag{" + t.GetAddressString() + "}"
}

// castToTag narrows a tag out of a request back to a CANopen tag. A tag wrapped by the subscription
// request builder doesn't expose its inner tag, so it is re-parsed from its address, which works
// because GetAddressString round-trips.
func castToTag(tag apiModel.PlcTag) (Tag, bool) {
	if tag == nil {
		return Tag{}, false
	}
	if canopenTag, ok := tag.(Tag); ok {
		return canopenTag, true
	}
	parsed, err := NewTagHandler().ParseTag(tag.GetAddressString())
	if err != nil {
		return Tag{}, false
	}
	canopenTag, ok := parsed.(Tag)
	return canopenTag, ok
}
