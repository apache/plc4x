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

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Example struct {
	*ExampleDelegate1
	ExampleDelegate2
	aField *uint16
	bField string
	cField int
	dField ExampleStringer `stringer:"true"`
	eField ExampleStringer `asPtr:"true"`
	fField *ExampleStringer
	gField *ExampleStringer
	hField time.Time
	iField time.Duration
	jField TestStruct1
	kField *TestStruct1
	lField TestStruct2
	mField *time.Time
	nField *time.Duration
}

type ExampleDelegate1 struct {
}

func (ExampleDelegate1) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	panic("no op")
}

type ExampleDelegate2 struct {
}

func (ExampleDelegate2) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	panic("no op")
}

type ExampleDelegate3 interface {
	utils.Serializable
	Example()
}

type ExampleStringer struct {
}

func (ExampleStringer) String() string {
	return ""
}

type TestStruct1 struct {
}

func (t *TestStruct1) Serialize() ([]byte, error) {
	panic("no use")
}

func (t *TestStruct1) SerializeWithWriteBuffer(_ context.Context, _ utils.WriteBuffer) error {
	panic("no use")
}

var _ utils.Serializable = (*TestStruct1)(nil)

type TestStruct2 struct {
}

func (t TestStruct2) String() string {
	panic("no use")
}

var _ fmt.Stringer = (*TestStruct2)(nil)
