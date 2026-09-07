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

// fakeWriteDevice is a minimal UDP BACnet/IP device that answers every
// confirmed WriteProperty / WritePropertyMultiple request with the spec-correct
// success response: a SimpleAck (clause 15.9) echoing the request's invoke id
// and service choice.
type fakeWriteDevice struct {
	conn    *net.UDPConn
	wg      sync.WaitGroup
	log     zerolog.Logger
	gotReqs int
	ctx     context.Context
	mu      sync.Mutex

	// wrongInvokeId, when true, makes the device reply with an off-by-one
	// invoke id so we can assert the client does NOT accept the mismatch.
	wrongInvokeId bool
	// replyError, when true, makes the device reply with a BACnet Error-PDU
	// (WRITE_ACCESS_DENIED) instead of a SimpleAck.
	replyError bool
}

func startFakeWriteDevice(t *testing.T, log zerolog.Logger) *fakeWriteDevice {
	t.Helper()
	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
	conn, err := net.ListenUDP("udp4", addr)
	require.NoError(t, err)
	d := &fakeWriteDevice{conn: conn, log: log, ctx: t.Context()}
	d.wg.Add(1)
	go d.serve()
	return d
}

func (d *fakeWriteDevice) port() int { return d.conn.LocalAddr().(*net.UDPAddr).Port }

func (d *fakeWriteDevice) requestCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.gotReqs
}

func (d *fakeWriteDevice) serve() {
	defer d.wg.Done()
	buf := make([]byte, 4096)
	for {
		n, src, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			return // socket closed
		}
		data := make([]byte, n)
		copy(data, buf[:n])

		invokeId, choice, ok := d.extractInvokeIdAndChoice(data)
		if !ok {
			d.log.Warn().Msg("fake write device: could not extract invoke id; ignoring")
			continue
		}
		d.mu.Lock()
		d.gotReqs++
		d.mu.Unlock()

		replyInvokeId := invokeId
		if d.wrongInvokeId {
			replyInvokeId = invokeId + 1
		}

		var resp model.BVLC
		switch {
		case choice == model.BACnetConfirmedServiceChoice_READ_PROPERTY:
			resp = buildReadAckFor(replyInvokeId)
		case d.replyError:
			resp = d.buildErrorReply(replyInvokeId, choice)
		default:
			resp = d.buildSimpleAck(replyInvokeId, choice)
		}
		theBytes, err := resp.Serialize()
		if err != nil {
			d.log.Error().Err(err).Msg("fake write device: serialize response")
			continue
		}
		if _, err := d.conn.WriteToUDP(theBytes, src); err != nil {
			d.log.Error().Err(err).Msg("fake write device: write response")
			continue
		}
		d.log.Info().Uint8("invokeId", replyInvokeId).Int("bytes", len(theBytes)).Stringer("dst", src).Msg("fake write device sent reply")
	}
}

func (d *fakeWriteDevice) extractInvokeIdAndChoice(data []byte) (uint8, model.BACnetConfirmedServiceChoice, bool) {
	bvlc, err := model.BVLCParse[model.BVLC](d.ctx, data)
	if err != nil {
		d.log.Error().Err(err).Msg("fake write device: parse BVLC")
		return 0, 0, false
	}
	npduRetriever, ok := bvlc.(interface{ GetNpdu() model.NPDU })
	if !ok {
		return 0, 0, false
	}
	apdu := npduRetriever.GetNpdu().GetApdu()
	cr, ok := apdu.(model.APDUConfirmedRequest)
	if !ok {
		d.log.Warn().Msgf("fake write device: not a confirmed request: %T", apdu)
		return 0, 0, false
	}
	return cr.GetInvokeId(), cr.GetServiceRequest().GetServiceChoice(), true
}

func (d *fakeWriteDevice) buildSimpleAck(invokeId uint8, choice model.BACnetConfirmedServiceChoice) model.BVLC {
	apdu := model.NewAPDUSimpleAck(invokeId, choice)
	return wrapAPDU(apdu, false, nil)
}

// buildReadAckFor mirrors fakeBacnetDevice.buildReadPropertyAck: a ReadProperty
// ComplexAck for ANALOG_INPUT,1/PRESENT_VALUE echoing the invoke id.
func buildReadAckFor(invokeId uint8) model.BVLC {
	serviceAck := model.NewBACnetServiceAckReadProperty(
		0,
		model.CreateBACnetContextTagObjectIdentifier(0, uint16(model.BACnetObjectType_ANALOG_INPUT), 1),
		model.CreateBACnetPropertyIdentifierTagged(1, uint32(model.BACnetPropertyIdentifier_PRESENT_VALUE)),
		nil,
		constructedDataFromTag(model.CreateBACnetApplicationTagReal(23.5)),
	)
	apdu := model.NewAPDUComplexAck(false, false, invokeId, nil, nil, serviceAck, nil, nil)
	return wrapAPDU(apdu, false, nil)
}

func (d *fakeWriteDevice) buildErrorReply(invokeId uint8, choice model.BACnetConfirmedServiceChoice) model.BVLC {
	base := buildErrorAPDU(model.ErrorClass_PROPERTY, model.ErrorCode_WRITE_ACCESS_DENIED)
	apdu := model.NewAPDUError(invokeId, choice, base.GetError())
	return wrapAPDU(apdu, false, nil)
}

func (d *fakeWriteDevice) stop() {
	_ = d.conn.Close()
	d.wg.Wait()
}

// writePresentValue issues a single WriteRequest for ANALOG_VALUE,1/PRESENT_VALUE
// against the given connection and returns the response code. The bool reports
// whether the Execute channel resolved at all (false == timed out / hung).
func writePresentValue(t *testing.T, conn plc4go.PlcConnection) (apiModel.PlcResponseCode, bool) {
	t.Helper()
	wr, err := conn.WriteRequestBuilder().
		AddTagAddress("pv", "ANALOG_VALUE,1/PRESENT_VALUE", float32(75.0)).
		Build()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 6*time.Second)
	t.Cleanup(cancel)
	select {
	case <-ctx.Done():
		return 0, false
	case res := <-wr.Execute(ctx):
		if res.GetErr() != nil {
			t.Logf("write error: %v", res.GetErr())
			return 0, false
		}
		return res.GetResponse().GetResponseCode("pv"), true
	}
}

func newWriteTestConnection(t *testing.T, log zerolog.Logger, port int) plc4go.PlcConnection {
	t.Helper()
	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=3000", port)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	conn, err := dm.GetConnection(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, conn.Close()) })
	return conn
}

// TestNativeBacnetWrite_SimpleAckResolves is the core regression: a WriteProperty
// answered by a spec-correct SimpleAck must resolve the Execute channel with OK.
func TestNativeBacnetWrite_SimpleAckResolves(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeWriteDevice(t, log)
	t.Cleanup(device.stop)

	conn := newWriteTestConnection(t, log, device.port())
	code, ok := writePresentValue(t, conn)
	t.Logf("device received %d request(s)", device.requestCount())
	require.True(t, ok, "write did not complete (Execute never resolved — hung on SimpleAck)")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
}

// TestNativeBacnetWrite_WrongInvokeIdDoesNotResolve confirms a SimpleAck bearing
// the wrong invoke id is NOT accepted (the request must hang / time out rather
// than falsely report success).
func TestNativeBacnetWrite_WrongInvokeIdDoesNotResolve(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeWriteDevice(t, log)
	device.wrongInvokeId = true
	t.Cleanup(device.stop)

	conn := newWriteTestConnection(t, log, device.port())
	_, ok := writePresentValue(t, conn)
	assert.False(t, ok, "a SimpleAck with the wrong invoke id must NOT resolve the write")
}

// TestNativeBacnetWrite_AfterReads reproduces the real E2E ordering: the same
// connection is used for several ReadProperty requests first (advancing the
// shared invoke-id generator so the write lands on a non-zero invoke id, and
// exercising the idle period between requests), then a WriteProperty is issued
// on that same connection. The write must resolve with OK.
func TestNativeBacnetWrite_AfterReads(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeWriteDevice(t, log)
	t.Cleanup(device.stop)

	conn := newWriteTestConnection(t, log, device.port())

	// Several reads first (mirrors discovery) so the write uses a later invoke id.
	for i := range 6 {
		code, val, ok := readPresentValue(t, conn)
		require.True(t, ok, "read %d hung", i)
		assert.Equal(t, apiModel.PlcResponseCode_OK, code)
		assert.InDelta(t, 23.5, val, 0.001)
	}

	code, ok := writePresentValue(t, conn)
	t.Logf("device received %d request(s)", device.requestCount())
	require.True(t, ok, "write after reads did not complete (Execute never resolved)")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
}

// TestNativeBacnetWrite_ErrorPduResolvesWithCode confirms a BACnet Error-PDU for
// the write resolves the request with the mapped error code.
func TestNativeBacnetWrite_ErrorPduResolvesWithCode(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeWriteDevice(t, log)
	device.replyError = true
	t.Cleanup(device.stop)

	conn := newWriteTestConnection(t, log, device.port())
	code, ok := writePresentValue(t, conn)
	require.True(t, ok, "an Error-PDU must resolve the write (not hang)")
	assert.Equal(t, apiModel.PlcResponseCode_ACCESS_DENIED, code)
}
