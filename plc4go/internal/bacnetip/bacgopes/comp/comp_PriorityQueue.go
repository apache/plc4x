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

// PriorityQueue implements heap.Interface and holds Items.
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
