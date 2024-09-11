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
	"iter"
	"os"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

type GenericFunction = func(args Args, kwargs KWArgs) error

type Arg any

type Args []any

var NoArgs = NewArgs()

func NewArgs(args ...any) Args {
	return args
}

// NA is a shortcut for NewArgs
var NA = NewArgs

// GetFromArgs gets a value fromArgs and if not present panics
func GetFromArgs[T any](args Args, index int) T {
	if index > len(args)-1 {
		panic(fmt.Sprintf("index out of bounds: %d(len %d of %s)", index, len(args), args))
	}
	aAtI := args[index]
	v, ok := aAtI.(T)
	if !ok {
		panic(fmt.Sprintf("argument #%d with type %T is not of type %T", index, aAtI, *new(T)))
	}
	return v
}

// GA is a shortcut for GetFromArgs
func GA[T any](args Args, index int) T {
	return GetFromArgs[T](args, index)
}

// GetFromArgsOptional gets a value from Args or return default if not present
func GetFromArgsOptional[T any](args Args, index int, defaultValue T) T {
	if index > len(args)-1 {
		return defaultValue
	}
	return args[index].(T)
}

// GAO is a shortcut for GetFromArgsOptional
func GAO[T any](args Args, index int, defaultValue T) T {
	return GetFromArgsOptional(args, index, defaultValue)
}

func (a Args) Format(f fmt.State, verb rune) {
	switch verb {
	case 'r':
		s := a.String()[1 : len(a.String())-1]
		_, _ = f.Write([]byte("(" + s + ")"))
	}
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

// NKW is a shortcut for NewKWArgs
var NKW = NewKWArgs

func (k KWArgs) Format(f fmt.State, verb rune) {
	switch verb {
	case 'r':
		_, _ = f.Write([]byte(k.String()))
	}
}

func (k KWArgs) String() string {
	r := ""
	for kk, ea := range k {
		switch kk {
		case KWCompRootMessage, KWCompBVLCIRequirements:
			// TODO: figure out if we want to control that for the %r above and do something different here
			continue
		}
		switch tea := ea.(type) {
		case []byte:
			ea = Btox(tea, ".")
		}
		r += fmt.Sprintf("'%s'=%v, ", kk, ea)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return "{" + r + "}"
}

// KW gets a value from KWArgs and if not present panics
func KW[T any](kwargs KWArgs, key KnownKey) T {
	r, ok := kwargs[key]
	if !ok {
		panic(fmt.Sprintf("key %v not found in kwargs", key))
	}
	return r.(T)
}

// KWO gets a value from KWArgs and if not present returns the supplied default value
func KWO[T any](kwargs KWArgs, key KnownKey, defaultValue T) T {
	r, ok := kwargs[key]
	if !ok {
		return defaultValue
	}
	return r.(T)
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
	// comm.PCI related keys

	KWCPCISource      = KnownKey("source")
	KWCPCIDestination = KnownKey("destination")
	KWCPCIData        = KnownKey("data")
	KWCPCIUserData    = KnownKey("user_data")

	////
	// PCI related keys

	KWPCIExpectingReply  = KnownKey("expecting_reply")
	KWPCINetworkPriority = KnownKey("network_priority")

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

	////
	// Compability layer keys

	KWCompRootMessage       = KnownKey("compRootMessage")
	KWCompBVLCIRequirements = KnownKey("compBVLCIRequirements")
)

// Nothing give NoArgs and NoKWArgs
func Nothing() (Args, KWArgs) {
	return NoArgs, NoKWArgs
}

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

// GA PriorityQueue implements heap.Interface and holds Items.
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

// OptionalOption2 allows options to be applied that might be optional
func OptionalOption2[V1 any, V2 any, T any](value1 *V1, value2 *V2, opt func(V1, V2) func(*T)) func(*T) {
	v1Set := value1 != nil
	v2Set := value2 != nil
	if (v1Set && !v2Set) || (!v1Set && v2Set) {
		return func(c *T) {}
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

// SortedMapIterator lets you iterate over an array in a deterministic way
func SortedMapIterator[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

type DebugPrinter = func(format string, a ...any)

func CreateDebugPrinter() DebugPrinter {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, "bacgopes")
	dir = dir[rootIndex:]
	qualifier := strings.ReplaceAll(dir, "/", ".")
	dirPrefix := path.Base(dir) + "_"

	bacgopesDebug := os.Getenv("BACGOPES_DEBUG")
	if strings.Contains(bacgopesDebug, qualifier) {
		return func(format string, a ...any) {
			_, file, _, ok := runtime.Caller(1)
			if !ok {
				return
			}
			base := path.Base(file)
			prefix := strings.TrimSuffix(base, ".go")
			prefix = strings.TrimPrefix(prefix, dirPrefix)
			_, _ = fmt.Fprintf(os.Stderr, "DEBUG:"+qualifier+"."+prefix+":"+format+"\n", a...)
		}
	}
	return nil
}

// OR returns a or b
func OR[T comparable](a T, b T) T {
	if reflect.ValueOf(a).IsNil() || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil()) { // TODO: check if there is another way than using reflect
		return b
	} else {
		return a
	}
}

// ToPtr gives a Ptr
func ToPtr[T any](value T) *T {
	return &value
}

type DefaultRFormatter struct {
	output []byte
}

func NewDefaultRFormatter() DefaultRFormatter {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		return DefaultRFormatter{}
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, "bacgopes")
	dir = dir[rootIndex:]
	dirPrefix := path.Base(dir) + "_"
	base := path.Base(file)
	prefix := strings.TrimSuffix(base, ".go")
	prefix = strings.TrimPrefix(prefix, dirPrefix)
	return DefaultRFormatter{
		output: []byte(fmt.Sprintf("<%s at 0x%x>", prefix, pc)),
	}
}

func (d DefaultRFormatter) Format(s fmt.State, v rune) {
	switch v {
	case 'r':
		_, _ = s.Write(d.output)
	}
}

func StructName() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, "bacgopes")
	dir = dir[rootIndex:]
	dirPrefix := path.Base(dir) + "_"
	base := path.Base(file)
	prefix := strings.TrimSuffix(base, ".go")
	return strings.TrimPrefix(prefix, dirPrefix)
}
