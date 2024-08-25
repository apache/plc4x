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

package test_apdu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func TestMaxApduLengthAcceptedEncode(t *testing.T) {
	t.Skip("Plc4x doesn't normalise at model level")
	apdu := model.NewAPDU(50)
	assert.Equal(t, 0, 50, apdu.ApduLength)
}

func TestMaxApduLengthAcceptedDecode(t *testing.T) {
	t.Skip("Plc4x doesn't normalise at model level")
	apdu := model.NewAPDU(0)
	serialize, err := apdu.Serialize()
	require.NoError(t, err)
	apduParse, err := model.APDUParse(context.Background(), serialize, 0)
	require.NoError(t, err)
	// TODO: no way to access the length
	_ = apduParse
}
