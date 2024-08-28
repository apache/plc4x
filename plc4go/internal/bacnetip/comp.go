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
	"cmp"
	"container/heap"
	"context"
	"fmt"
	"strings"

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

func (a Args) Get1NPDU() NPDU {
	return a[1].(NPDU)
}

func (a Args) Get0APDU() APDU {
	return a[0].(APDU)
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
		r += fmt.Sprintf("%d: %v (%T), ", i, ea, ea)
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
	PDUData
}

type messageBridge struct {
	*_PDUData
}

func NewMessageBridge(bytes ...byte) MessageBridge {
	return &messageBridge{&_PDUData{data: bytes}}
}

var _ MessageBridge = (*messageBridge)(nil)

func (m *messageBridge) Serialize() ([]byte, error) {
	return m.data, nil
}

func (m *messageBridge) SerializeWithWriteBuffer(_ context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteByteArray("Bytes", m.data)
}

func (m *messageBridge) GetLengthInBytes(_ context.Context) uint16 {
	return uint16(len(m.data))
}

func (m *messageBridge) GetLengthInBits(ctx context.Context) uint16 {
	return m.GetLengthInBytes(ctx) * 8
}

// An PriorityItem is something we manage in a priority queue.
type PriorityItem[P cmp.Ordered, V any] struct {
	value    V // The value of the item; arbitrary.
	priority P // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func (p *PriorityItem[P, V]) String() string {
	return fmt.Sprintf("[%v: %v-%v], ", p.index, p.priority, p.value)
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[P cmp.Ordered, V any] []*PriorityItem[P, V]

//goland:noinspection GoMixedReceiverTypes
func (pq PriorityQueue[P, V]) Len() int { return len(pq) }

//goland:noinspection GoMixedReceiverTypes
func (pq PriorityQueue[P, V]) Less(i, j int) bool {
	return cmp.Less((pq)[i].priority, (pq)[j].priority)
}

//goland:noinspection GoMixedReceiverTypes
func (pq PriorityQueue[P, V]) Swap(i, j int) {
	(pq)[i], (pq)[j] = (pq)[j], (pq)[i]
	(pq)[i].index = i
	(pq)[j].index = j
}

//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) Push(x any) {
	n := len(*pq)
	item := x.(*PriorityItem[P, V])
	item.index = n
	*pq = append(*pq, item)
}

//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) clear() {
	if pq == nil {
		return
	}
	*pq = PriorityQueue[P, V]{}
}

// update modifies the priority and value of an Item in the queue.
//
//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) update(item *PriorityItem[P, V], value V, priority P) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

//goland:noinspection GoMixedReceiverTypes
func (pq PriorityQueue[P, V]) String() string {
	var buf strings.Builder
	for _, p := range pq {
		buf.WriteString(p.String())
	}
	s := buf.String()
	if s == "" {
		return ""
	}
	var p P
	var v V
	return fmt.Sprintf("PriorityQueue[%T,%T]{%s}", p, v, s[:len(s)-2])
}

// NillableKey is a key which can be used in maps
type NillableKey[T any] struct {
	value T
	isNil bool
}

// NK creates a new NillableKey of type K
func NK[T any, K NillableKey[T]](value *T) K {
	var _nk NillableKey[T]
	if value == nil {
		_nk.isNil = true
		return K(_nk)
	}
	_nk.value = *value
	return K(_nk)
}
