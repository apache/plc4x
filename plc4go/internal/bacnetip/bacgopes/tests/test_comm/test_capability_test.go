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

package test_comm

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/capability"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// TODO: big WIP

type BaseCollector struct {
	Collector
}

func NewBaseCollector(localLog zerolog.Logger, subs ...CollectorOrCapability) (*BaseCollector, func()) {
	b := &BaseCollector{}
	var init func()
	b.Collector, init = NewCollector(localLog, append([]CollectorOrCapability{CoCCo(b)}, subs...)...)
	return b, init
}

func (b *BaseCollector) Transform(value any) any {
	for fn := range b.CapabilityFunctions("transform") {
		value = fn(NA(b, value), NoKWArgs())
	}
	return value
}

func (b *BaseCollector) SearchCapability(_ ...CollectorOrCapability) []Capability {
	// No-op
	return []Capability{}
}

type PlusOne struct {
	Capability
}

func NewPlusOne() *PlusOne {
	p := &PlusOne{}
	p.Capability = NewCapability()
	return p
}

func (p *PlusOne) transform(value any) any {
	return value.(int) + 1
}

type TimesTen struct {
	Capability
}

func NewTimesTen() *TimesTen {
	t := &TimesTen{}
	t.Capability = NewCapability()
	return t
}

func (p *TimesTen) transform(value any) any {
	return value.(int) * 10
}

type MakeList struct {
	Capability
}

//####################################
//####################################
//##
//## Example classes

type Example1 struct {
	*BaseCollector
}

func NewExample1(localLog zerolog.Logger) *Example1 {
	e := &Example1{}
	var init func()
	e.BaseCollector, init = NewBaseCollector(localLog, CoCCo(e))
	init()
	return e
}

type Example2 struct {
	*BaseCollector
	*PlusOne
}

func NewExample2(localLog zerolog.Logger) *Example2 {
	e := &Example2{}
	var init func()
	e.BaseCollector, init = NewBaseCollector(localLog, CoCCo(e))
	e.PlusOne = NewPlusOne()
	init()
	return e
}

func (e *Example2) Serialize() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *Example2) SerializeWithWriteBuffer(_ context.Context, _ utils.WriteBuffer) error {
	// NO-OP
	return nil
}

func (e *Example2) String() string {
	return e.BaseCollector.String() + e.PlusOne.String()
}

type Example3 struct {
	*BaseCollector
	*TimesTen
	*PlusOne
}

func NewExample3(localLog zerolog.Logger) *Example3 {
	e := &Example3{}
	var init func()
	e.BaseCollector, init = NewBaseCollector(localLog, CoCCo(e))
	init()
	e.TimesTen = NewTimesTen()
	e.PlusOne = NewPlusOne()
	return e
}

func (e *Example3) Serialize() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *Example3) SerializeWithWriteBuffer(_ context.Context, _ utils.WriteBuffer) error {
	// NO-OP
	return nil
}

func (e *Example3) String() string {
	return e.BaseCollector.String() + e.TimesTen.String() + e.PlusOne.String()
}

type Example4 struct {
	*BaseCollector
	*MakeList
	*TimesTen
}

func (e *Example4) String() string {
	return e.BaseCollector.String() + e.MakeList.String() + e.TimesTen.String()
}

func (e *Example4) Serialize() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *Example4) SerializeWithWriteBuffer(_ context.Context, _ utils.WriteBuffer) error {
	// NO-OP
	return nil
}

func TestExamples(t *testing.T) {
	t.Skip("needs more looking at when actually being needed") // TODO: ignore for now
	t.Run("test_example_1", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 1, NewExample1(testingLogger).Transform(1))
	})
	t.Run("test_example_2", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 3, NewExample2(testingLogger).Transform(2))
	})
	t.Run("test_example_3", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 31, NewExample3(testingLogger).Transform(3))
	})
	t.Run("test_example_4", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, NewExample3(testingLogger).Transform(4))
	})
	t.Run("test_example_5", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 6, NewExample3(testingLogger).Transform(5))
	})
}
