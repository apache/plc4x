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

package profinet

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Tag is one data item of a submodule's cyclic data, plc4j profinet-ng's address syntax:
//
//	<slot>.<subslot>.<INPUT|OUTPUT>[.<index>]:<TYPE>[<count>]
//
// index is the data item's position in the submodule's GSDML Input or Output (default 0); TYPE is
// the value type to read it as (BOOL, USINT, UINT, INT, REAL, BYTE, WORD, ...), and count the
// number of elements, which for BOOL is the number of bits of an item used as bits. The count may
// also come before the colon, as plc4j writes it.
type Tag struct {
	Slot, Subslot uint16
	Output        bool
	Index         int
	ValueType     apiValues.PlcValueType
	Count         int
}

var tagPattern = regexp.MustCompile(`^(\d{1,5})\.(\d{1,5})\.(INPUT|OUTPUT)(?:\.(\d{1,5}))?(?:\[(\d{1,5})\])?:([A-Za-z_]+)(?:\[(\d{1,5})\])?$`)

func parseTag(address string) (Tag, error) {
	m := tagPattern.FindStringSubmatch(address)
	if m == nil {
		return Tag{}, errors.Errorf("unable to parse Profinet address %q: expected <slot>.<subslot>.<INPUT|OUTPUT>[.<index>]:<TYPE>[<count>]", address)
	}
	slot, err := strconv.ParseUint(m[1], 10, 16)
	if err != nil {
		return Tag{}, errors.Errorf("slot %s is out of range", m[1])
	}
	subslot, err := strconv.ParseUint(m[2], 10, 16)
	if err != nil {
		return Tag{}, errors.Errorf("subslot %s is out of range", m[2])
	}
	tag := Tag{Slot: uint16(slot), Subslot: uint16(subslot), Output: m[3] == "OUTPUT", Count: 1}
	if m[4] != "" {
		tag.Index, _ = strconv.Atoi(m[4])
	}
	valueType, ok := apiValues.PlcValueTypeByName(m[6])
	if !ok {
		return Tag{}, errors.Errorf("unknown value type %s", m[6])
	}
	tag.ValueType = valueType
	for _, count := range []string{m[5], m[7]} {
		if count != "" {
			n, _ := strconv.Atoi(count)
			if n < 1 {
				return Tag{}, errors.New("an element count is at least 1")
			}
			tag.Count = n
		}
	}
	return tag, nil
}

func (t Tag) direction() string {
	if t.Output {
		return "OUTPUT"
	}
	return "INPUT"
}

func (t Tag) GetAddressString() string {
	address := fmt.Sprintf("%d.%d.%s.%d:%s", t.Slot, t.Subslot, t.direction(), t.Index, t.ValueType)
	if t.Count > 1 {
		address += fmt.Sprintf("[%d]", t.Count)
	}
	return address
}

func (t Tag) GetValueType() apiValues.PlcValueType {
	return t.ValueType
}

func (t Tag) GetArrayInfo() []apiModel.ArrayInfo {
	if t.Count > 1 {
		return []apiModel.ArrayInfo{&arrayInfo{upperBound: uint32(t.Count)}}
	}
	return []apiModel.ArrayInfo{}
}

func (t Tag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionCyclic
}

func (t Tag) GetDuration() time.Duration {
	return 0
}

func (t Tag) String() string {
	return "profinet.Tag{" + t.GetAddressString() + "}"
}

func (t Tag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := t.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (t Tag) SerializeWithWriteBuffer(_ context.Context, wb utils.WriteBuffer) error {
	for _, c := range []string{"PlcTagItem", "tag", "ProfinetTag"} {
		if err := wb.PushContext(c); err != nil {
			return err
		}
	}
	address := t.GetAddressString()
	if err := wb.WriteString("address", uint32(len(address)*8), address, utils.WithEncoding("UTF-8")); err != nil {
		return err
	}
	for _, c := range []string{"ProfinetTag", "tag", "PlcTagItem"} {
		if err := wb.PopContext(c); err != nil {
			return err
		}
	}
	return nil
}

type arrayInfo struct {
	upperBound uint32
}

func (a *arrayInfo) GetSize() uint32       { return a.upperBound }
func (a *arrayInfo) GetLowerBound() uint32 { return 0 }
func (a *arrayInfo) GetUpperBound() uint32 { return a.upperBound - 1 }
func (a *arrayInfo) GetBase() uint32       { return 0 }
func (a *arrayInfo) IsRange() bool         { return true }
func (a *arrayInfo) String() string        { return fmt.Sprintf("[%d]", a.upperBound) }

type TagHandler struct{}

func NewTagHandler() TagHandler {
	return TagHandler{}
}

func (TagHandler) ParseTag(address string) (apiModel.PlcTag, error) {
	return parseTag(address)
}

func (TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("the profinet driver lists a device's data with a browse request, not a query")
}

func castToTag(tag apiModel.PlcTag) (Tag, bool) {
	if tag == nil {
		return Tag{}, false
	}
	if t, ok := tag.(Tag); ok {
		return t, true
	}
	parsed, err := parseTag(tag.GetAddressString())
	return parsed, err == nil
}
