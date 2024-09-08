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
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/capability"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// TODO: big WIP

type BaseCollector struct {
	*Collector
}

func NewBaseCollector(localLog zerolog.Logger) *BaseCollector {
	b := &BaseCollector{}
	b.Collector = NewCollector(localLog)
	return b
}

func (b BaseCollector) transform(value any) any {
	panic("not implemented") // TODO: implement me
	return value
}

type PlusOne struct {
	*Capability
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
	*Capability
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
	*Capability
}

//####################################
//####################################
//##
//## Example classes

type Example1 struct {
	*BaseCollector
}

func NewExample1(localLog zerolog.Logger) *Example1 {
	b := &Example1{}
	b.BaseCollector = NewBaseCollector(localLog)
	return b
}

func (e *Example1) transform(value any) any {
	return e.BaseCollector.transform(value)
}

type Example2 struct {
	*BaseCollector
	*PlusOne
}

func NewExample2(localLog zerolog.Logger) *Example2 {
	b := &Example2{}
	b.BaseCollector = NewBaseCollector(localLog)
	b.PlusOne = NewPlusOne()
	return b
}

func (e *Example2) transform(value any) any {
	return e.BaseCollector.transform(value)
}

type Example3 struct {
	*BaseCollector
	*TimesTen
	*PlusOne
}

func (e *Example3) transform(value any) any {
	return e.BaseCollector.transform(value)
}

func NewExample3(localLog zerolog.Logger) *Example3 {
	b := &Example3{}
	b.BaseCollector = NewBaseCollector(localLog)
	b.TimesTen = NewTimesTen()
	b.PlusOne = NewPlusOne()
	return b
}

type Example4 struct {
	*BaseCollector
	*MakeList
	*TimesTen
}

func TestExamples(t *testing.T) {
	t.Skip("big WIP...") // TODO: big WIP
	t.Run("test_example_1", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 1, NewExample1(testingLogger).transform(1))
	})
	t.Run("test_example_2", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 3, NewExample2(testingLogger).transform(2))
	})
	t.Run("test_example_3", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 31, NewExample3(testingLogger).transform(3))
	})
	t.Run("test_example_4", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, NewExample3(testingLogger).transform(4))
	})
	t.Run("test_example_5", func(t *testing.T) {
		testingLogger := testutils.ProduceTestingLogger(t)
		assert.Equal(t, 6, NewExample3(testingLogger).transform(5))
	})
}
