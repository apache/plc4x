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

package _default

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// A SendRequest whose Send fails must not leave its expectation registered:
// the caller already receives the send error, so a later expectation timeout
// would fire the error handler a second time for the same request.
func Test_defaultCodec_SendRequest_failedSendRemovesExpectation(t *testing.T) {
	requirements := NewMockDefaultCodecRequirements(t)
	requirements.EXPECT().Send(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("nope"))
	m := &defaultCodec{
		DefaultCodecRequirements: requirements,
		notifyExpireWorker:       make(chan struct{}, 100),
		notifyReceiveWorker:      make(chan struct{}, 100),
		log:                      testutils.ProduceTestingLogger(t),
	}

	err := m.SendRequest(
		testutils.TestContext(t),
		t.Name(),
		nil,
		func(spi.Message) bool { return true },
		func(spi.Message) error { return nil },
		func(error) error { return nil },
	)
	assert.Error(t, err)

	m.expectationsChangeMutex.RLock()
	defer m.expectationsChangeMutex.RUnlock()
	assert.Empty(t, m.expectations, "failed send left its expectation registered")
}
