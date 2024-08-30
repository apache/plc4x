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

type DataReaderEnumDefault[T any, I any] struct {
	enumResolver func(I) (T, bool)
	dataReader   DataReader[I]
}

var _ DataReaderEnum[string] = (*DataReaderEnumDefault[string, string])(nil)

func NewDataReaderEnumDefault[T any, I any](enumResolver func(I) (T, bool), dataReader DataReader[I]) *DataReaderEnumDefault[T, I] {
	return &DataReaderEnumDefault[T, I]{
		enumResolver: enumResolver,
		dataReader:   dataReader,
	}
}

func (d *DataReaderEnumDefault[T, I]) GetPos() uint16 {
	return d.dataReader.GetPos()
}

func (d *DataReaderEnumDefault[T, I]) SetPos(pos uint16) {
	d.dataReader.SetPos(pos)
}

func (d *DataReaderEnumDefault[T, I]) GetByteOrder() binary.ByteOrder {
	return d.dataReader.GetByteOrder()
}

func (d *DataReaderEnumDefault[T, I]) SetByteOrder(byteOrder binary.ByteOrder) {
	d.dataReader.SetByteOrder(byteOrder)
}

func (d *DataReaderEnumDefault[T, I]) Read(ctx context.Context, logicalName string, readerArgs ...utils.WithReaderArgs) (T, error) {
	var zero T
	rawValue, err := d.dataReader.Read(ctx, logicalName, readerArgs...)
	if err != nil {
		return zero, errors.Wrap(err, "error reading raw data")
	}
	enumValue, ok := d.enumResolver(rawValue)
	if !ok { // TODO: decide if we want to log... maybe we always pass the context so we can extract that
		//d.log.Debug().Str("logicalName", logicalName).Interface("rawValue", rawValue).Msg("no enum value found")
	}
	return enumValue, nil
}

func (d *DataReaderEnumDefault[T, I]) PullContext(logicalName string, readerArgs ...utils.WithReaderArgs) error {
	return d.dataReader.PullContext(logicalName, readerArgs...)
}

func (d *DataReaderEnumDefault[T, I]) CloseContext(logicalName string, readerArgs ...utils.WithReaderArgs) error {
	return d.dataReader.CloseContext(logicalName, readerArgs...)
}

func (d *DataReaderEnumDefault[T, I]) GetReadBuffer() utils.ReadBuffer {
	return d.dataReader.GetReadBuffer()
}
