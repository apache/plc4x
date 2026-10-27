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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// recordCodec is a captureCodec that also hands over every request it sends.
type recordCodec struct {
	*captureCodec
	sent chan spi.Message
}

func (c *recordCodec) SendRequest(ctx context.Context, name string, message spi.Message, accepts spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	c.sent <- message
	return c.captureCodec.SendRequest(ctx, name, message, accepts, handleMessage, handleError)
}

// sentSubscribeCOV unwraps the SubscribeCOV service request of a sent frame.
func sentSubscribeCOV(t *testing.T, codec *recordCodec) readWriteModel.BACnetConfirmedServiceRequestSubscribeCOV {
	t.Helper()
	var message spi.Message
	select {
	case message = <-codec.sent:
	case <-time.After(time.Second):
		t.Fatal("no SubscribeCOV was sent")
	}
	_, err := message.(interface{ Serialize() ([]byte, error) }).Serialize()
	require.NoError(t, err, "the frame must encode")
	apdu := message.(readWriteModel.BVLC).(interface{ GetNpdu() readWriteModel.NPDU }).GetNpdu().GetApdu()
	request, ok := apdu.(readWriteModel.APDUConfirmedRequest).GetServiceRequest().(readWriteModel.BACnetConfirmedServiceRequestSubscribeCOV)
	require.True(t, ok, "not a SubscribeCOV")
	return request
}

// ASHRAE 135 13.14.1: a SubscribeCOV without issueConfirmedNotifications and
// lifetime cancels the subscription, while lifetime 0 subscribes for an
// indefinite time. Unsubscribe must send the cancellation, or it leaves the
// device with a subscription that never expires.
func TestSubscriber_UnsubscribeSendsACancellationNotAnIndefiniteSubscription(t *testing.T) {
	codec := &recordCodec{captureCodec: newCaptureCodec(nil), sent: make(chan spi.Message, 4)}
	conn := &Connection{configuration: createDefaultConfiguration(), messageCodec: codec}
	conn.configuration.CovLifetimeSeconds = 0 // indefinite, so the subscription itself carries lifetime 0
	conn.driverContext = NewDriverContext(conn.configuration)
	s := NewSubscriber(conn)

	handle := subscribeOne(t, s, codec.captureCodec, "ai1", 1)
	subscribe := sentSubscribeCOV(t, codec)
	require.NotNil(t, subscribe.GetLifetimeInSeconds(), "an indefinite subscription carries lifetime 0")
	assert.EqualValues(t, 0, subscribe.GetLifetimeInSeconds().GetPayload().GetActualValue())
	assert.NotNil(t, subscribe.GetIssueConfirmed())

	results := s.Unsubscribe(testutils.TestContext(t), spiModel.NewDefaultPlcUnsubscriptionRequest(
		[]apiModel.PlcSubscriptionHandle{handle.DefaultPlcSubscriptionHandle}))
	cancel := sentSubscribeCOV(t, codec)
	assert.Nil(t, cancel.GetIssueConfirmed(), "a cancellation leaves out issueConfirmedNotifications")
	assert.Nil(t, cancel.GetLifetimeInSeconds(), "a cancellation leaves out the lifetime")
	assert.EqualValues(t, handle.subscriberProcessId, cancel.GetSubscriberProcessIdentifier().GetPayload().GetActualValue())

	select {
	case h := <-codec.handlers:
		require.NoError(t, h.handleMessage(wrapAPDU(readWriteModel.NewAPDUSimpleAck(0, readWriteModel.BACnetConfirmedServiceChoice_SUBSCRIBE_COV), false, nil)))
	case <-time.After(time.Second):
		t.Fatal("Unsubscribe never sent its request")
	}
	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
	case <-time.After(time.Second):
		t.Fatal("Unsubscribe did not deliver a result")
	}
	assert.Nil(t, s.lookupHandle(handle.subscriberProcessId), "the handle is gone after Unsubscribe")
}
