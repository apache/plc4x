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

	"github.com/stretchr/testify/assert"
)

// ActiveSubscriptionCount feeds the connection cache's subscription-aware
// idle reaping: it must reflect the live COV handle count across all
// subscribers of the connection.
func TestConnection_ActiveSubscriptionCount(t *testing.T) {
	conn := &Connection{}
	assert.Zero(t, conn.ActiveSubscriptionCount(), "fresh connection has no subscriptions")

	s1 := NewSubscriber(conn)
	s2 := NewSubscriber(conn)
	conn.addSubscriber(s1)
	conn.addSubscriber(s2)
	assert.Zero(t, conn.ActiveSubscriptionCount(), "subscribers without handles don't count")

	s1.storeHandle(&SubscriptionHandle{subscriberProcessId: 1})
	s1.storeHandle(&SubscriptionHandle{subscriberProcessId: 2})
	s2.storeHandle(&SubscriptionHandle{subscriberProcessId: 7})
	assert.Equal(t, 3, conn.ActiveSubscriptionCount(), "handles sum across subscribers")

	s1.removeHandle(2)
	assert.Equal(t, 2, conn.ActiveSubscriptionCount(), "removed handles drop out of the count")
}
