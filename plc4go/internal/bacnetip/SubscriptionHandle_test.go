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

package bacnetip

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

func TestNewSubscriptionHandle_PopulatesIdentity(t *testing.T) {
	s := newTestSubscriber(t)
	tag := &plcTag{}
	h := NewSubscriptionHandle(s, "myTag", tag, apiModel.SubscriptionChangeOfState, 250*time.Millisecond)
	require.NotNil(t, h)
	assert.Equal(t, "myTag", h.tagName)
	assert.Equal(t, tag, h.tag)
	assert.Equal(t, apiModel.SubscriptionChangeOfState, h.subscriptionType)
	assert.Equal(t, 250*time.Millisecond, h.interval)
	// DefaultPlcSubscriptionHandle must be wired so callers can treat the
	// handle as a generic api/model.PlcSubscriptionHandle.
	assert.NotNil(t, h.DefaultPlcSubscriptionHandle)
}
