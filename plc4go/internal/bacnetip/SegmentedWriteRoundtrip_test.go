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
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiTransports "github.com/apache/plc4x/plc4go/pkg/api/transports"
	model "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// fakeSegmentedDevice is a UDP BACnet/IP device that understands SEGMENTED
// confirmed requests (ASHRAE 135 clause 5.4 receive side): it collects the
// segments, acknowledges per its actual window size, reassembles the service
// request, and answers with a SimpleAck once the final segment arrived.
// Unsegmented confirmed requests get an immediate SimpleAck.
type fakeSegmentedDevice struct {
	conn *net.UDPConn
	wg   sync.WaitGroup
	log  zerolog.Logger
	ctx  context.Context

	// actualWindowSize is what the device answers in its SegmentAcks (and how
	// many segments it lets pass between acks).
	actualWindowSize uint8
	// nakFirstSegment, when true, answers segment 0 with a negative ack
	// (still confirming seq 0) — the sender must treat it as "received up to
	// 0" and continue with segment 1.
	nakFirstSegment bool

	mu               sync.Mutex
	buf              []byte
	sinceAck         uint8
	segmentsReceived int
	unsegmented      int
	reassembled      model.BACnetConfirmedServiceRequest
	nakSent          bool
}

func startFakeSegmentedDevice(t *testing.T, log zerolog.Logger, actualWindowSize uint8) *fakeSegmentedDevice {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	require.NoError(t, err)
	d := &fakeSegmentedDevice{conn: conn, log: log, ctx: t.Context(), actualWindowSize: actualWindowSize}
	d.wg.Add(1)
	go d.serve()
	return d
}

func (d *fakeSegmentedDevice) port() int { return d.conn.LocalAddr().(*net.UDPAddr).Port }

func (d *fakeSegmentedDevice) stop() {
	_ = d.conn.Close()
	d.wg.Wait()
}

func (d *fakeSegmentedDevice) stats() (segments, unsegmented int, reassembled model.BACnetConfirmedServiceRequest) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.segmentsReceived, d.unsegmented, d.reassembled
}

func (d *fakeSegmentedDevice) serve() {
	defer d.wg.Done()
	buf := make([]byte, 4096)
	for {
		n, src, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			return // socket closed
		}
		data := make([]byte, n)
		copy(data, buf[:n])

		bvlc, err := model.BVLCParse[model.BVLC](d.ctx, data)
		if err != nil {
			d.log.Error().Err(err).Msg("fake segmented device: parse BVLC")
			continue
		}
		npduRetriever, ok := bvlc.(interface{ GetNpdu() model.NPDU })
		if !ok {
			continue
		}
		cr, ok := npduRetriever.GetNpdu().GetApdu().(model.APDUConfirmedRequest)
		if !ok {
			continue
		}

		for _, reply := range d.handleConfirmedRequest(cr) {
			theBytes, err := reply.Serialize()
			if err != nil {
				d.log.Error().Err(err).Msg("fake segmented device: serialize reply")
				continue
			}
			if _, err := d.conn.WriteToUDP(theBytes, src); err != nil {
				d.log.Error().Err(err).Msg("fake segmented device: send reply")
			}
		}
	}
}

// handleConfirmedRequest implements the device-side segment collection and
// returns the BVLC replies to put on the wire (segment acks and/or the final
// SimpleAck).
func (d *fakeSegmentedDevice) handleConfirmedRequest(cr model.APDUConfirmedRequest) []model.BVLC {
	d.mu.Lock()
	defer d.mu.Unlock()

	invokeId := cr.GetInvokeId()

	if !cr.GetSegmentedMessage() {
		d.unsegmented++
		return []model.BVLC{wrapAPDU(model.NewAPDUSimpleAck(invokeId, cr.GetServiceRequest().GetServiceChoice()), false, nil)}
	}

	seqPtr := cr.GetSequenceNumber()
	if seqPtr == nil {
		d.log.Error().Msg("fake segmented device: segmented request without sequence number")
		return nil
	}
	seq := *seqPtr
	d.segmentsReceived++

	if seq == 0 {
		// First segment: raw bytes start with the service-choice octet.
		d.buf = append([]byte{}, cr.GetSegment()...)
	} else {
		if cr.GetSegmentServiceChoice() == nil {
			d.log.Error().Uint8("seq", seq).Msg("fake segmented device: follow-up segment without service choice")
			return nil
		}
		d.buf = append(d.buf, cr.GetSegment()...)
	}

	var replies []model.BVLC

	if d.nakFirstSegment && seq == 0 && !d.nakSent {
		// Negative ack still confirming segment 0 — sender must continue at 1.
		d.nakSent = true
		d.sinceAck = 0
		replies = append(replies, wrapAPDU(model.NewAPDUSegmentAck(true, true, invokeId, 0, d.actualWindowSize), false, nil))
		return replies
	}

	d.sinceAck++
	final := !cr.GetMoreFollows()
	// Clause 5.4.5: segment 0 is always acknowledged on its own; afterwards the
	// device acks every actual-window-size segments and the final segment.
	if seq == 0 || d.sinceAck >= d.actualWindowSize || final {
		d.sinceAck = 0
		replies = append(replies, wrapAPDU(model.NewAPDUSegmentAck(false, true, invokeId, seq, d.actualWindowSize), false, nil))
	}

	if final {
		reassembled, err := model.BACnetConfirmedServiceRequestParse[model.BACnetConfirmedServiceRequest](d.ctx, d.buf, uint32(len(d.buf)))
		if err != nil {
			d.log.Error().Err(err).Int("bytes", len(d.buf)).Msg("fake segmented device: reassembled payload does not parse")
			return replies
		}
		d.reassembled = reassembled
		replies = append(replies, wrapAPDU(model.NewAPDUSimpleAck(invokeId, reassembled.GetServiceChoice()), false, nil))
	}
	return replies
}

// newSegmentedTestConnection opens a driver connection declaring the peer's
// APDU ceiling and segmentation capability via the new connection options.
func newSegmentedTestConnection(t *testing.T, log zerolog.Logger, port int, peerOptions string) plc4go.PlcConnection {
	t.Helper()
	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=3000&%s", port, peerOptions)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	conn, err := dm.GetConnection(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, conn.Close()) })
	return conn
}

// executeLargeWrite issues a WritePropertyMultiple across enough tags that the
// serialized request exceeds a 206-byte peer APDU ceiling.
func executeLargeWrite(t *testing.T, conn plc4go.PlcConnection, tags int) (apiModel.PlcWriteResponse, error) {
	t.Helper()
	builder := conn.WriteRequestBuilder()
	for i := range tags {
		builder.AddTagAddress(fmt.Sprintf("pv%d", i), fmt.Sprintf("ANALOG_VALUE,%d/PRESENT_VALUE", i+1), float32(i)+0.5)
	}
	wr, err := builder.Build()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	t.Cleanup(cancel)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-wr.Execute(ctx):
		return res.GetResponse(), res.GetErr()
	}
}

// TestSegmentedWrite_RoundtripWindow1 pins the core clause 5.4 transmit flow
// against a device acking every segment (actual window size 1): the oversized
// WritePropertyMultiple must arrive as multiple segments whose reassembly
// parses back to the original request, and the write must resolve OK.
func TestSegmentedWrite_RoundtripWindow1(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeSegmentedDevice(t, log, 1)
	t.Cleanup(device.stop)

	conn := newSegmentedTestConnection(t, log, device.port(),
		"PeerMaxApduLengthAccepted=206&PeerSegmentationSupported=segmented-both")

	response, err := executeLargeWrite(t, conn, 24)
	require.NoError(t, err)

	segments, unsegmented, reassembled := device.stats()
	t.Logf("device saw %d segments (%d unsegmented requests)", segments, unsegmented)
	require.GreaterOrEqual(t, segments, 2, "request must have been segmented")
	assert.Zero(t, unsegmented, "the oversized request must not go out unsegmented")
	require.NotNil(t, reassembled, "device must have reassembled the request")
	wpm, ok := reassembled.(model.BACnetConfirmedServiceRequestWritePropertyMultiple)
	require.True(t, ok, "reassembled request must be a WritePropertyMultiple, got %T", reassembled)
	assert.Len(t, wpm.GetData(), 24, "all write-access specs must survive reassembly")

	for i := range 12 {
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode(fmt.Sprintf("pv%d", i)))
	}
}

// TestSegmentedWrite_RoundtripWindow3 exercises the windowed burst path: the
// device only acks every 3rd segment, so the sender must honor the actual
// window size from the first ack and keep multiple segments in flight.
func TestSegmentedWrite_RoundtripWindow3(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeSegmentedDevice(t, log, 3)
	t.Cleanup(device.stop)

	conn := newSegmentedTestConnection(t, log, device.port(),
		"PeerMaxApduLengthAccepted=206&PeerSegmentationSupported=segmented-both")

	response, err := executeLargeWrite(t, conn, 50)
	require.NoError(t, err)

	segments, _, reassembled := device.stats()
	t.Logf("device saw %d segments", segments)
	require.GreaterOrEqual(t, segments, 5, "a 50-tag request at a 206-byte ceiling must span several segments")
	require.NotNil(t, reassembled)
	wpm, ok := reassembled.(model.BACnetConfirmedServiceRequestWritePropertyMultiple)
	require.True(t, ok)
	assert.Len(t, wpm.GetData(), 50)

	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("pv0"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("pv49"))
}

// TestSegmentedWrite_NegativeAckContinues pins the NAK handling: a negative
// SegmentAck confirming segment 0 must not abort the transfer — the sender
// resumes with segment 1 and the write still completes.
func TestSegmentedWrite_NegativeAckContinues(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeSegmentedDevice(t, log, 1)
	device.nakFirstSegment = true
	t.Cleanup(device.stop)

	conn := newSegmentedTestConnection(t, log, device.port(),
		"PeerMaxApduLengthAccepted=206&PeerSegmentationSupported=segmented-both")

	response, err := executeLargeWrite(t, conn, 24)
	require.NoError(t, err)
	_, _, reassembled := device.stats()
	require.NotNil(t, reassembled, "transfer must complete despite the NAK on segment 0")
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("pv0"))
}

// TestSegmentedWrite_PeerWithoutSegmentationFailsFast: with a declared APDU
// ceiling but no segmentation support, an oversized request must fail fast
// with an actionable error instead of provoking a device-side abort.
func TestSegmentedWrite_PeerWithoutSegmentationFailsFast(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeSegmentedDevice(t, log, 1)
	t.Cleanup(device.stop)

	conn := newSegmentedTestConnection(t, log, device.port(),
		"PeerMaxApduLengthAccepted=206&PeerSegmentationSupported=no-segmentation")

	_, err := executeLargeWrite(t, conn, 24)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support segmented requests")
	segments, unsegmented, _ := device.stats()
	assert.Zero(t, segments, "nothing may reach the wire")
	assert.Zero(t, unsegmented, "nothing may reach the wire")
}

// TestSegmentedWrite_SmallRequestStaysUnsegmented guards the common path: with
// peer capabilities declared, a request that FITS the ceiling must go out as a
// plain unsegmented confirmed request.
func TestSegmentedWrite_SmallRequestStaysUnsegmented(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeSegmentedDevice(t, log, 1)
	t.Cleanup(device.stop)

	conn := newSegmentedTestConnection(t, log, device.port(),
		"PeerMaxApduLengthAccepted=480&PeerSegmentationSupported=segmented-both")

	response, err := executeLargeWrite(t, conn, 2)
	require.NoError(t, err)
	segments, unsegmented, _ := device.stats()
	assert.Zero(t, segments, "a fitting request must not be segmented")
	assert.Equal(t, 1, unsegmented)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("pv0"))
}

// TestSegmentationOptionParsing_Unknown guards the conservative default: an
// unknown or absent peer segmentation string must NOT enable segmented sends.
func TestSegmentationOptionParsing_Unknown(t *testing.T) {
	for _, s := range []string{"", "bogus", "SEGMENTED-BOTH", strings.ToUpper("segmented-receive")} {
		assert.False(t, segmentationAcceptsSegmentedRequests(s), "%q must not enable segmented requests", s)
	}
	assert.True(t, segmentationAcceptsSegmentedRequests("segmented-both"))
	assert.True(t, segmentationAcceptsSegmentedRequests("segmented-receive"))
	assert.False(t, segmentationAcceptsSegmentedRequests("segmented-transmit"),
		"a peer that can only TRANSMIT segments cannot receive ours")
}
