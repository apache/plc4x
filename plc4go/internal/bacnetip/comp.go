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

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Arg any

type Args []any

var NoArgs = NewArgs()

func NewArgs(args ...any) Args {
	return args
}

func (a Args) Get0PDU() PDU {
	return a[0].(PDU)
}

func (a Args) Get1PDU() PDU {
	return a[1].(PDU)
}

func (a Args) Get0NPDU() NPDU {
	return a[0].(NPDU)
}

func (a Args) Get0NetworkAdapter() *NetworkAdapter {
	return a[0].(*NetworkAdapter)
}

func (a Args) Get0MultiplexClient() *_MultiplexClient {
	return a[0].(*_MultiplexClient)
}

func (a Args) Get0MultiplexServer() *_MultiplexServer {
	return a[0].(*_MultiplexServer)
}

func (a Args) String() string {
	r := ""
	for i, ea := range a {
		switch tea := ea.(type) {
		case []byte:
			ea = Btox(tea, ".")
		}
		r += fmt.Sprintf("%d: %v, ", i, ea)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return r
}

type KWArgs map[KnownKey]any

var NoKWArgs = NewKWArgs()

func NewKWArgs(kw ...any) KWArgs {
	if len(kw)%2 != 0 {
		panic("KWArgs must have an even number of arguments")
	}
	r := make(KWArgs)
	for i := 0; i < len(kw)-1; i += 2 {
		key, ok := kw[i].(KnownKey)
		if !ok {
			panic("keys must be of type KnownKey")
		}
		r[key] = kw[i+1]
	}
	return r
}
func (k KWArgs) String() string {
	r := ""
	for kk, ea := range k {
		switch tea := ea.(type) {
		case []byte:
			ea = Btox(tea, ".")
		}
		r += fmt.Sprintf("%s=%v, ", kk, ea)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return "{" + r + "}"
}

type KnownKey string

const (
	////
	// General keys

	KWAddActor   = KnownKey("addActor")
	KWDelActor   = KnownKey("delActor")
	KWActorError = KnownKey("actorError")
	KWError      = KnownKey("error")

	////
	// PDU related Keys

	KWPPDUSource     = KnownKey("pduSource")
	KWPDUDestination = KnownKey("pduDestination")
	KWPDUData        = KnownKey("pduData")

	////
	// NPDU related keys

	KWWirtnNetwork           = KnownKey("wirtnNetwork")
	KWIartnNetworkList       = KnownKey("iartnNetworkList")
	KWIcbrtnNetwork          = KnownKey("icbrtnNetwork")
	KWIcbrtnPerformanceIndex = KnownKey("icbrtnPerformanceIndex")
	KWRmtnRejectionReason    = KnownKey("rmtnRejectionReason")
	KWRmtnDNET               = KnownKey("rmtnDNET")
	KWRbtnNetworkList        = KnownKey("rbtnNetworkList")
	KWRatnNetworkList        = KnownKey("ratnNetworkList")
	KWIrtTable               = KnownKey("irtTable")
	KWIrtaTable              = KnownKey("irtaTable")
	KWEctnDNET               = KnownKey("ectnDNET")
	KWEctnTerminationTime    = KnownKey("ectnTerminationTime")
	KWDctnDNET               = KnownKey("dctnDNET")
	KWNniNet                 = KnownKey("nniNet")
	KWNniFlag                = KnownKey("nniFlag")

	////
	// BVLL related keys

	KWBvlciResultCode = KnownKey("bvlciResultCode")
	KWBvlciBDT        = KnownKey("bvlciBDT")
	KWBvlciAddress    = KnownKey("bvlciAddress")
	KWFdAddress       = KnownKey("fdAddress")
	KWFdTTL           = KnownKey("fdTTL")
	KWFdRemain        = KnownKey("fdRemain")
	KWBvlciTimeToLive = KnownKey("bvlciTimeToLive")
	KWBvlciFDT        = KnownKey("bvlciFDT")
)

type MessageBridge interface {
	spi.Message
	_PDUDataRequirements
}

type messageBridge struct {
	Bytes []byte
}

func NewMessageBridge(bytes ...byte) MessageBridge {
	m := &messageBridge{Bytes: make([]byte, len(bytes))}
	copy(m.Bytes, bytes)
	if len(m.Bytes) == 0 {
		m.Bytes = nil
	}
	return m
}

var _ MessageBridge = (*messageBridge)(nil)

func (m *messageBridge) String() string {
	return Btox(m.Bytes, "")
}

func (m *messageBridge) Serialize() ([]byte, error) {
	return m.Bytes, nil
}

func (m *messageBridge) SerializeWithWriteBuffer(_ context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteByteArray("Bytes", m.Bytes)
}

func (m *messageBridge) GetLengthInBytes(_ context.Context) uint16 {
	return uint16(len(m.Bytes))
}

func (m *messageBridge) GetLengthInBits(ctx context.Context) uint16 {
	return m.GetLengthInBytes(ctx) * 8
}

func (m *messageBridge) getPDUData() []byte {
	return m.Bytes
}
