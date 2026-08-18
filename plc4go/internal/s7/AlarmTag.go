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

package s7

import (
	"context"
	"strings"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	alarmTagAddressAlarm  = "ALM"
	alarmTagAddressQueryS = "QUERY:ALARM_S"
	alarmTagAddressQuery8 = "QUERY:ALARM_8"
)

type AlarmTagKind uint8

const (
	AlarmTagPush AlarmTagKind = iota
	AlarmTagQuery
)

// AlarmTag addresses the S7 alarm/event services instead of a memory location. "ALM" subscribes
// to pushed alarm indications, "QUERY:ALARM_S"/"QUERY:ALARM_8" run a one-shot alarm query.
type AlarmTag struct {
	address   string
	kind      AlarmTagKind
	queryType readWriteModel.QueryType
}

func tryParseAlarmTag(tagAddress string) (*AlarmTag, bool) {
	switch {
	case strings.EqualFold(tagAddress, alarmTagAddressAlarm):
		return &AlarmTag{address: alarmTagAddressAlarm, kind: AlarmTagPush, queryType: readWriteModel.QueryType_ALARM_S}, true
	case strings.EqualFold(tagAddress, alarmTagAddressQueryS):
		return &AlarmTag{address: alarmTagAddressQueryS, kind: AlarmTagQuery, queryType: readWriteModel.QueryType_ALARM_S}, true
	case strings.EqualFold(tagAddress, alarmTagAddressQuery8):
		return &AlarmTag{address: alarmTagAddressQuery8, kind: AlarmTagQuery, queryType: readWriteModel.QueryType_ALARM_8}, true
	default:
		return nil, false
	}
}

func (m *AlarmTag) GetAddressString() string {
	return m.address
}

func (m *AlarmTag) GetValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (m *AlarmTag) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (m *AlarmTag) GetKind() AlarmTagKind {
	return m.kind
}

func (m *AlarmTag) GetQueryType() readWriteModel.QueryType {
	return m.queryType
}

func (m *AlarmTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionEvent
}

func (m *AlarmTag) GetDuration() time.Duration {
	return 0
}

func (m *AlarmTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *AlarmTag) SerializeWithWriteBuffer(ctx context.Context, wb utils.WriteBuffer) error {
	if err := wb.PushContext("S7AlarmTag"); err != nil {
		return err
	}
	if err := wb.WriteString("address", uint32(len(m.address)*8), m.address); err != nil {
		return err
	}
	if err := wb.PopContext("S7AlarmTag"); err != nil {
		return err
	}
	return nil
}

func (m *AlarmTag) String() string {
	return m.address
}
