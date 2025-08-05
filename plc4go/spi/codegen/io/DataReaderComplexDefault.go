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

package io

import (
	"context"
	"encoding/binary"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DataReaderComplexDefault[T any] struct {
	complexTypeSupplier func(context.Context, utils.ReadBuffer) (T, error)
	readBuffer          utils.ReadBuffer
}

var _ DataReaderComplex[string] = (*DataReaderComplexDefault[string])(nil)

func NewDataReaderComplexDefault[T any](complexTypeSupplier func(context.Context, utils.ReadBuffer) (T, error), readBuffer utils.ReadBuffer) *DataReaderComplexDefault[T] {
	return &DataReaderComplexDefault[T]{
		complexTypeSupplier: complexTypeSupplier,
		readBuffer:          readBuffer,
	}
}

func (d *DataReaderComplexDefault[T]) GetPos() uint16 {
	return d.readBuffer.GetPos()
}

func (d *DataReaderComplexDefault[T]) SetPos(pos uint16) {
	d.readBuffer.Reset(pos)
}

func (d *DataReaderComplexDefault[T]) GetByteOrder() binary.ByteOrder {
	return d.readBuffer.GetByteOrder()
}

func (d *DataReaderComplexDefault[T]) SetByteOrder(byteOrder binary.ByteOrder) {
	d.readBuffer.SetByteOrder(byteOrder)
}

func (d *DataReaderComplexDefault[T]) Read(ctx context.Context, logicalName string, readerArgs ...utils.WithReaderArgs) (T, error) {
	return d.ReadComplex(ctx, logicalName, d.complexTypeSupplier, readerArgs...)
}

func (d *DataReaderComplexDefault[T]) ReadComplex(ctx context.Context, logicalName string, supplier func(context.Context, utils.ReadBuffer) (T, error), readerArgs ...utils.WithReaderArgs) (T, error) {
	var zero T
	// TODO: it might be even better if we default to value like in other places... on the other hand a complex type has always a proper logical name so this might be fine like that
	if logicalName != "" {
		err := d.readBuffer.PullContext(logicalName, readerArgs...)
		if err != nil {
			return zero, errors.Wrap(err, "error pulling context")
		}
	}
	t, err := supplier(ctx, d.readBuffer)
	if err != nil {
		return zero, errors.Wrap(err, "error getting value")
	}
	if logicalName != "" {
		err := d.readBuffer.CloseContext(logicalName, readerArgs...)
		if err != nil {
			return zero, errors.Wrap(err, "error closing context")
		}
	}
	return t, err
}

func (d *DataReaderComplexDefault[T]) PullContext(logicalName string, readerArgs ...utils.WithReaderArgs) error {
	return d.readBuffer.PullContext(logicalName, readerArgs...)
}

func (d *DataReaderComplexDefault[T]) CloseContext(logicalName string, readerArgs ...utils.WithReaderArgs) error {
	return d.readBuffer.CloseContext(logicalName, readerArgs...)
}

func (d *DataReaderComplexDefault[T]) GetReadBuffer() utils.ReadBuffer {
	return d.readBuffer
}
