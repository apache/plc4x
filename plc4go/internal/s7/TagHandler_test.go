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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
)

func TestTagHandlerParseTag(t *testing.T) {
	handler := NewTagHandler()

	t.Run("S5TIME address", func(t *testing.T) {
		tag, err := handler.ParseTag("%M0:S5TIME")
		require.NoError(t, err)
		s7Tag := tag.(PlcTag)
		assert.Equal(t, readWriteModel.TransportSize_S5TIME, s7Tag.GetDataType())
		assert.Equal(t, readWriteModel.MemoryArea_FLAGS_MARKERS, s7Tag.GetMemoryArea())
	})
	t.Run("S5TIME in data block", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB10:4:S5TIME")
		require.NoError(t, err)
		s7Tag := tag.(PlcTag)
		assert.Equal(t, readWriteModel.TransportSize_S5TIME, s7Tag.GetDataType())
		assert.Equal(t, uint16(10), s7Tag.GetBlockNumber())
	})
	t.Run("fixed length string", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB69:68:STRING(20)")
		require.NoError(t, err)
		stringTag := tag.(PlcStringTag)
		assert.Equal(t, uint16(20), stringTag.stringLength)
	})
	t.Run("var length string short form", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB69:68:STRING")
		require.NoError(t, err)
		stringTag := tag.(PlcStringTag)
		assert.Equal(t, uint16(254), stringTag.stringLength)
		assert.Equal(t, readWriteModel.TransportSize_STRING, stringTag.GetDataType())
	})
	t.Run("var length string long form with array", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB69.DBX68:WSTRING[3]")
		require.NoError(t, err)
		stringTag := tag.(PlcStringTag)
		assert.Equal(t, uint16(254), stringTag.stringLength)
		assert.Equal(t, uint16(3), stringTag.GetNumElements())
		assert.Equal(t, readWriteModel.TransportSize_WSTRING, stringTag.GetDataType())
	})
	t.Run("alarm subscription tag", func(t *testing.T) {
		tag, err := handler.ParseTag("ALM")
		require.NoError(t, err)
		alarmTag := tag.(*AlarmTag)
		assert.Equal(t, AlarmTagPush, alarmTag.GetKind())
	})
	t.Run("alarm query tags", func(t *testing.T) {
		tag, err := handler.ParseTag("QUERY:ALARM_S")
		require.NoError(t, err)
		alarmTag := tag.(*AlarmTag)
		assert.Equal(t, AlarmTagQuery, alarmTag.GetKind())
		assert.Equal(t, readWriteModel.QueryType_ALARM_S, alarmTag.GetQueryType())

		tag, err = handler.ParseTag("query:alarm_8")
		require.NoError(t, err)
		alarmTag = tag.(*AlarmTag)
		assert.Equal(t, readWriteModel.QueryType_ALARM_8, alarmTag.GetQueryType())
	})
	t.Run("existing forms still parse", func(t *testing.T) {
		for _, address := range []string{"%Q0.0:BOOL", "%M100:INT[10]", "%DB1.DBX0.0:BOOL", "%DB1:0.0:BOOL", "%I0:BYTE"} {
			_, err := handler.ParseTag(address)
			assert.NoError(t, err, address)
		}
	})
	t.Run("oversized byte offset is rejected", func(t *testing.T) {
		_, err := handler.ParseTag("%DB1:70000:INT")
		assert.Error(t, err)
	})
	t.Run("bogus string type suffix is rejected", func(t *testing.T) {
		_, err := handler.ParseTag("%DB1:0:STRINGX")
		assert.Error(t, err)
	})
}
