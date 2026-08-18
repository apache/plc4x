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

package eip

import (
	"context"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestListIdentityRequestRoundTrip(t *testing.T) {
	request := buildListIdentityRequest()
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	require.NoError(t, request.SerializeWithWriteBuffer(context.Background(), wb))
	raw := wb.GetBytes()
	// 24-byte encapsulation header, ListIdentity command 0x63, empty payload
	require.Len(t, raw, 24)
	assert.Equal(t, uint16(0x63), binary.LittleEndian.Uint16(raw[0:2]))
	assert.Equal(t, uint16(0), binary.LittleEndian.Uint16(raw[2:4]))

	rb := utils.NewReadBufferByteBased(raw, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	parsed, err := readWriteModel.EipPacketParseWithBuffer[readWriteModel.EipPacket](context.Background(), rb, false)
	require.NoError(t, err)
	_, ok := parsed.(readWriteModel.EipListIdentityRequest)
	assert.True(t, ok)
}
