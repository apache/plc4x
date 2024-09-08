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

package codegen

import (
	"context"
	"encoding/binary"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldCommons[T any] struct {
}

func (FieldCommons[T]) ExtractByteOrder(readerWriterArgs ...utils.WithReaderWriterArgs) *binary.ByteOrder {
	for _, arg := range readerWriterArgs {
		switch rwArg := arg.(type) {
		case withOptionByteOrder:
			return &rwArg.byteOrder
		}
	}
	return nil
}

func (FieldCommons[T]) SwitchParseByteOrderIfNecessary(ctx context.Context, runnable func(context.Context) (T, error), byteOrderAware utils.ByteOrderAware, wantedByteOrder *binary.ByteOrder) (T, error) {
	var zero T
	if runnable == nil {
		return zero, errors.New("runnable is nil")
	}
	if byteOrderAware == nil {
		return zero, errors.New("byteOrderAware is nil")
	}
	currentByteOrder := byteOrderAware.GetByteOrder()
	if wantedByteOrder == nil || currentByteOrder == *wantedByteOrder {
		return runnable(ctx)
	}
	byteOrderAware.SetByteOrder(*wantedByteOrder)
	run, err := runnable(ctx)
	byteOrderAware.SetByteOrder(currentByteOrder)
	return run, err
}

func (FieldCommons[T]) SwitchParseByteOrderIfNecessarySerializeWrapped(ctx context.Context, runnable func(context.Context) error, byteOrderAware utils.ByteOrderAware, wantedByteOrder *binary.ByteOrder) error {
	if runnable == nil {
		return errors.New("runnable is nil")
	}
	if byteOrderAware == nil {
		return errors.New("byteOrderAware is nil")
	}
	currentByteOrder := byteOrderAware.GetByteOrder()
	if wantedByteOrder == nil || currentByteOrder == *wantedByteOrder {
		return runnable(ctx)
	}
	byteOrderAware.SetByteOrder(*wantedByteOrder)
	err := runnable(ctx)
	byteOrderAware.SetByteOrder(currentByteOrder)
	return err
}
