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

package knxnetip

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// newSilentReader builds a Reader whose gateway never answers a group-address read, so
// the only way out of a read is the context.
func newSilentReader(t *testing.T) *Reader {
	t.Helper()
	release := make(chan struct{})
	var once sync.Once
	t.Cleanup(func() { once.Do(func() { close(release) }) })
	connection, _ := newStubbedConnection(t, nil, func(_ spi.Message, _ spi.HandleMessage, _ spi.HandleError) error {
		// The gateway stays silent until the test is over.
		<-release
		return errors.New("codec closed")
	})
	connection.valueCache = map[uint16][]byte{}
	connection.tagHandler = NewTagHandler()
	connection.ClientKnxAddress = driverModel.NewKnxAddress(1, 1, 1)
	return NewReader(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
}

// readGroupAddressAsync runs a read in the background so a regression shows up as a test
// failure instead of a hanging test-run.
func readGroupAddressAsync(t *testing.T, reader *Reader, ctx context.Context, tag GroupAddressTag) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	t.Helper()
	type readResult struct {
		responseCode apiModel.PlcResponseCode
		plcValue     apiValues.PlcValue
	}
	results := make(chan readResult, 1)
	go func() {
		responseCode, plcValue := reader.readGroupAddress(ctx, tag)
		results <- readResult{responseCode, plcValue}
	}()
	select {
	case result := <-results:
		return result.responseCode, result.plcValue
	case <-time.After(10 * time.Second):
		require.FailNow(t, "readGroupAddress didn't honor the context")
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
}

func Test_Reader_readGroupAddress_canceledContext(t *testing.T) {
	reader := newSilentReader(t)
	tagType := driverModel.KnxDatapointType_DPT_Switch
	tag := NewGroupAddress3LevelPlcTag("1", "2", "3", &tagType)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	responseCode, plcValue := readGroupAddressAsync(t, reader, ctx, tag)
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, responseCode)
	assert.Nil(t, plcValue)
}

func Test_Reader_readGroupAddress_expiringContext(t *testing.T) {
	reader := newSilentReader(t)
	tagType := driverModel.KnxDatapointType_DPT_Switch
	tag := NewGroupAddress3LevelPlcTag("1", "2", "3", &tagType)

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()

	responseCode, plcValue := readGroupAddressAsync(t, reader, ctx, tag)
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, responseCode)
	assert.Nil(t, plcValue)
}

// Test_Reader_Read_canceledContext drives the same case through the public read api.
func Test_Reader_Read_canceledContext(t *testing.T) {
	reader := newSilentReader(t)
	readRequestBuilder := reader.connection.ReadRequestBuilder()
	readRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch")
	readRequest, err := readRequestBuilder.Build()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	select {
	case readResult := <-readRequest.Execute(ctx):
		require.NoError(t, readResult.GetErr())
		assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, readResult.GetResponse().GetResponseCode("switch"))
	case <-time.After(10 * time.Second):
		require.FailNow(t, "read didn't honor the context")
	}
}

// Test_Reader_readGroupAddress_fromCache makes sure the happy path is untouched: a value
// which is already in the local cache never reaches the codec.
//
// It doubles as the regression test for the "skip the first byte" guard which used to be
// driven by KnxDatapointType.GetLengthInBits (hardcoded to 32 for every datapoint-type by
// the generated model, so it always fired): the cached payload is the raw group-value
// payload from the bus, and the datapoint parser consumes the reserved bits/byte itself.
func Test_Reader_readGroupAddress_fromCache(t *testing.T) {
	tests := []struct {
		name      string
		tagType   driverModel.KnxDatapointType
		payload   []byte
		assertion func(t *testing.T, value apiValues.PlcValue)
	}{
		{
			name:    "a datapoint-type which fits into the embedded data bits",
			tagType: driverModel.KnxDatapointType_DPT_Switch,
			payload: []byte{0x01},
			assertion: func(t *testing.T, value apiValues.PlcValue) {
				assert.True(t, value.GetBool())
			},
		},
		{
			name:    "a datapoint-type which needs its own data byte",
			tagType: driverModel.KnxDatapointType_DPT_Scaling,
			payload: []byte{0x00, 0x80},
			assertion: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint8(0x80), value.GetUint8())
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := newSilentReader(t)
			// 1/2/3 -> 5 bits main, 3 bits middle, 8 bits sub
			reader.connection.valueCache[1<<11|2<<8|3] = test.payload

			tagType := test.tagType
			tag := NewGroupAddress3LevelPlcTag("1", "2", "3", &tagType)

			responseCode, plcValue := readGroupAddressAsync(t, reader, t.Context(), tag)
			assert.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
			require.NotNil(t, plcValue)
			test.assertion(t, plcValue)
		})
	}
}
