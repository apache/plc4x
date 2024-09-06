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

package comp

import (
	"cmp"
	"container/heap"
	"fmt"
	"strings"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

type Arg any

type Args []any

var NoArgs = NewArgs()

func NewArgs(args ...any) Args {
	return args
}

func Get[T any](args Args, index int) T {
	return argsGetOrPanic[T](args, index)
}

func GetOptional[T any](args Args, index int, defaultValue T) T {
	return argsGetOrDefault(args, index, defaultValue)
}

func argsGetOrPanic[T any](args Args, index int) T {
	if index > len(args)-1 {
		panic(fmt.Sprintf("index out of bounds: %d(len %d of %s)", index, len(args), args))
	}
	aAtI := args[index]
	v, ok := aAtI.(T)
	if !ok {
		var _type T
		panic(fmt.Sprintf("argument #%d with type %T is not of type %T", index, aAtI, _type))
	}
	return v
}

func argsGetOrDefault[T any](args Args, index int, defaultValue T) T {
	if index > len(args)-1 {
		return defaultValue
	}
	return args[index].(T)
}

func (a Args) String() string {
	r := ""
	for i, ea := range a {
		eat := fmt.Sprintf("%T", ea)
		switch tea := ea.(type) {
		case []byte:
			ea = Btox(tea, ".")
		case fmt.Stringer:
			if tea != nil {
				teaString := tea.String()
				ea = teaString
				if strings.Contains(teaString, "\n") {
					ea = "\n" + teaString + "\n"
				}
			}
		}
		r += fmt.Sprintf("%d: %v (%s), ", i, ea, eat)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return "[" + r + "]"
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

// An PriorityItem is something we manage in a priority queue.
type PriorityItem[P cmp.Ordered, V any] struct {
	Value    V // The value of the item; arbitrary.
	Priority P // The priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The Index of the item in the heap.
}

func (p *PriorityItem[P, V]) String() string {
	v := fmt.Sprintf("%v", p.Value)
	if strings.Contains(v, "\n") {
		v = "\n" + v + "\n"
	}
	return fmt.Sprintf("[%v: prio %v - value %s], ", p.Index, p.Priority, v)
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[P cmp.Ordered, V any] []*PriorityItem[P, V]

//goland:noinspection GoMixedReceiverTypes
func (pq PriorityQueue[P, V]) Len() int { return len(pq) }

//goland:noinspection GoMixedReceiverTypes
func (pq PriorityQueue[P, V]) Less(i, j int) bool {
	return cmp.Less((pq)[i].Priority, (pq)[j].Priority)
}

//goland:noinspection GoMixedReceiverTypes
func (pq PriorityQueue[P, V]) Swap(i, j int) {
	(pq)[i], (pq)[j] = (pq)[j], (pq)[i]
	(pq)[i].Index = i
	(pq)[j].Index = j
}

//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) Push(x any) {
	n := len(*pq)
	item := x.(*PriorityItem[P, V])
	item.Index = n
	*pq = append(*pq, item)
}

//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) Clear() {
	if pq == nil {
		return
	}
	*pq = PriorityQueue[P, V]{}
}

// update modifies the priority and value of an Item in the queue.
//
//goland:noinspection GoMixedReceiverTypes
func (pq *PriorityQueue[P, V]) update(item *PriorityItem[P, V], value V, priority P) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
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
	if strings.Contains(s, "\n") {
		s = "\n" + s + "\n"
	}
	var p P
	var v V
	return fmt.Sprintf("PriorityQueue[%T,%T]{%s}", p, v, s[:len(s)-2])
}

// NillableKey is a key which can be used in maps
type NillableKey[T any] struct {
	Value T
	IsNil bool
}

func (n NillableKey[T]) String() string {
	if n.IsNil {
		return "nil"
	}
	return fmt.Sprintf("%v", n.Value)
}

// NK creates a new NillableKey of type K
func NK[T any, K NillableKey[T]](value *T) K {
	var _nk NillableKey[T]
	if value == nil {
		_nk.IsNil = true
		return K(_nk)
	}
	_nk.Value = *value
	return K(_nk)
}

type Comparable interface {
	Equals(other any) bool
}

type Copyable interface {
	DeepCopy() any
}

// DeepCopy copies things implementing Copyable
func DeepCopy[T Copyable](copyable Copyable) T {
	return copyable.DeepCopy().(T)
}

// CopyPtr copies things that are a pointer to something
func CopyPtr[T any](t *T) *T {
	if t == nil {
		return nil
	}
	tc := *t
	return &tc
}

type Updater interface {
	Update(Arg) error
}

type Encoder interface {
	Encode(Arg) error
}

type Decoder interface {
	Decode(Arg) error
}

// OptionalOption allows options to be applied that might be optional
func OptionalOption[V any, T any](value *V, opt func(V) func(*T)) func(*T) {
	if value != nil {
		return opt(*value)
	}
	return func(c *T) {}
}

// OptionalOptionDual allows options to be applied that might be optional
func OptionalOptionDual[V1 any, V2 any, T any](value1 *V1, value2 *V2, opt func(V1, V2) func(*T)) func(*T) {
	v1Set := value1 != nil
	v2Set := value2 != nil
	if (v1Set && !v2Set) || (!v1Set && v2Set) {
		panic("Dual options must be both set together")
	}
	if v1Set {
		return opt(*value1, *value2)
	}
	return func(c *T) {}
}

type MissingRequiredParameter struct {
	Message string
}

func (m MissingRequiredParameter) Error() string {
	return m.Message
}
