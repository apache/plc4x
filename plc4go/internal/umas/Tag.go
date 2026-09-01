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

package umas

import (
	"context"
	"regexp"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// symbolicAddressPattern is plc4j's SymbolicUmasTag.SYMBOLIC_ADDRESS_PATTERN, byte for byte: an
// IEC 61131-3 identifier, optionally indexed and optionally followed by dotted member accesses which
// may themselves be indexed. Go's regexp has no possessive or backreference constructs in play here,
// so the expression carries over unchanged.
//
//	g_r32                  a plain global variable
//	g_plant.meta.r32       a nested struct member
//	g_arrInt[3]            an array element
//	g_plant.items[2].value both mixed
var symbolicAddressPattern = regexp.MustCompile(`^([a-zA-Z_]\w*)(\[\d+])*(\.([a-zA-Z_]\w*)(\[\d+])*)*$`)

// PlcTag is a parsed UMAS tag address. UMAS has no numeric addressing in plc4x - every tag names a
// symbol of the PLC project, which the connection resolves against the symbol table it downloads
// during connect. Ported from plc4j's SymbolicUmasTag.
//
// It is also an apiModel.PlcSubscriptionTag, because UMAS subscriptions are emulated by polling the
// read path and the subscription-request builder only accepts tags which are one. plc4j gets this
// for free from PollingSubscriptionConnectionBase, which wraps the plain tag itself.
type PlcTag interface {
	apiModel.PlcSubscriptionTag
	utils.Serializable

	// GetSymbolicAddress is the symbol name this tag addresses, exactly as the user spelled it.
	GetSymbolicAddress() string
}

type plcTag struct {
	SymbolicAddress string
	// ValueType is what the symbol turned out to be. It is apiValues.NULL for a tag the user
	// asked for by name, because only the connection's symbol table knows the type; the browser
	// fills it in for the items it hands out.
	ValueType apiValues.PlcValueType
	ArrayInfo []apiModel.ArrayInfo
}

var _ PlcTag = plcTag{}

// NewTag is the tag a user address parses to: the type is unknown until the symbol table is
// consulted.
func NewTag(symbolicAddress string) PlcTag {
	return plcTag{
		SymbolicAddress: symbolicAddress,
		ValueType:       apiValues.NULL,
		ArrayInfo:       []apiModel.ArrayInfo{},
	}
}

// NewTagWithType is the tag the browser hands out, which knows what the symbol is.
func NewTagWithType(symbolicAddress string, valueType apiValues.PlcValueType, arrayInfo []apiModel.ArrayInfo) PlcTag {
	if arrayInfo == nil {
		arrayInfo = []apiModel.ArrayInfo{}
	}
	return plcTag{
		SymbolicAddress: symbolicAddress,
		ValueType:       valueType,
		ArrayInfo:       arrayInfo,
	}
}

func (m plcTag) GetSymbolicAddress() string {
	return m.SymbolicAddress
}

func (m plcTag) GetAddressString() string {
	return m.SymbolicAddress
}

func (m plcTag) GetValueType() apiValues.PlcValueType {
	return m.ValueType
}

func (m plcTag) GetArrayInfo() []apiModel.ArrayInfo {
	return m.ArrayInfo
}

// GetPlcSubscriptionType is what a tag which wasn't added through one of the typed builder methods
// defaults to. Polling can only emulate cyclic and change-of-state subscriptions, and reporting only
// what moved is the cheaper of the two.
func (m plcTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionChangeOfState
}

// GetDuration is not applicable: UMAS has no per-tag subscription duration, the poll interval comes
// from the subscription request.
func (m plcTag) GetDuration() time.Duration {
	return 0
}

func (m plcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("UmasTag"); err != nil {
		return err
	}
	if err := writeBuffer.WriteString("symbolicAddress", uint32(len(m.SymbolicAddress)*8), m.SymbolicAddress); err != nil {
		return err
	}
	if m.ValueType != apiValues.NULL {
		if err := writeBuffer.WriteString("valueType", uint32(len(m.ValueType.String())*8), m.ValueType.String()); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("UmasTag"); err != nil {
		return err
	}
	return nil
}

func (m plcTag) String() string {
	return m.SymbolicAddress
}
