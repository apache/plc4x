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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// subscribeOne runs Subscribe for one tag on s, acks its SubscribeCOV and
// returns the handle it stored.
func subscribeOne(t *testing.T, s *Subscriber, codec *captureCodec, tagName string, instance uint32) *SubscriptionHandle {
	t.Helper()
	objType := readWriteModel.BACnetObjectType_ANALOG_INPUT
	tag := &plcTag{ObjectId: objectId{ObjectIdType: &objType, ObjectIdInstance: instance}}
	req := spiModel.NewDefaultPlcSubscriptionRequest(s, []string{tagName},
		map[string]apiModel.PlcSubscriptionTag{tagName: tag},
		map[string]apiModel.PlcSubscriptionType{tagName: apiModel.SubscriptionChangeOfState},
		map[string]time.Duration{tagName: 0}, nil)
	results := s.Subscribe(testutils.TestContext(t), req)
	select {
	case h := <-codec.handlers:
		require.NoError(t, h.handleMessage(wrapAPDU(readWriteModel.NewAPDUSimpleAck(0, readWriteModel.BACnetConfirmedServiceChoice_SUBSCRIBE_COV), false, nil)))
	case <-time.After(time.Second):
		t.Fatal("Subscribe never sent a SubscribeCOV")
	}
	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		require.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode(tagName))
	case <-time.After(time.Second):
		t.Fatal("Subscribe did not deliver a result")
	}
	for _, h := range s.handles {
		if h.tagName == tagName {
			return h
		}
	}
	t.Fatalf("no handle stored for %s", tagName)
	return nil
}

// Every request builder makes its own Subscriber, and a notification reaches
// all subscribers of the connection, each routing by process id alone. Two
// subscribers on one connection must therefore never hand out the same id, or
// one object's notifications update the other object's tag too.
func TestSubscriber_ProcessIdsAreUniquePerConnection(t *testing.T) {
	codec := newCaptureCodec(nil)
	conn := &Connection{configuration: createDefaultConfiguration(), messageCodec: codec}
	conn.driverContext = NewDriverContext(conn.configuration)

	first, second := NewSubscriber(conn), NewSubscriber(conn)
	a := subscribeOne(t, first, codec, "a", 1)
	b := subscribeOne(t, second, codec, "b", 2)
	require.NotEqual(t, a.subscriberProcessId, b.subscriberProcessId, "process ids collide across subscribers of one connection")

	gotA, gotB := &captureConsumer{}, &captureConsumer{}
	first.Register(gotA.consume, []apiModel.PlcSubscriptionHandle{a.DefaultPlcSubscriptionHandle})
	second.Register(gotB.consume, []apiModel.PlcSubscriptionHandle{b.DefaultPlcSubscriptionHandle})

	notification := buildCOVNotification(b.subscriberProcessId, 1234, readWriteModel.BACnetObjectType_ANALOG_INPUT, 2, 21.5)
	for _, s := range conn.subscribers {
		s.HandleUnconfirmedCOVNotification(notification)
	}
	assert.Empty(t, gotA.events, "b's notification must not reach a's consumer")
	assert.Len(t, gotB.events, 1)
}
