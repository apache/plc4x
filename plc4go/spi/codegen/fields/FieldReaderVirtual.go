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

package fields

import (
	"context"
	"fmt"

	"github.com/cstockton/go-conv"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldReaderVirtual[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderVirtual[T any](logger zerolog.Logger) *FieldReaderVirtual[T] {
	return &FieldReaderVirtual[T]{log: logger}
}

func (f *FieldReaderVirtual[T]) ReadVirtualField(ctx context.Context, logicalName string, klass any, valueExpression any, readerArgs ...utils.WithReaderArgs) (value T, err error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	switch klass.(type) {
	case *bool:
		return f.toT(conv.Bool(valueExpression))
	case *int8:
		return f.toT(conv.Int8(valueExpression))
	case *uint8:
		return f.toT(conv.Uint8(valueExpression))
	case *int16:
		return f.toT(conv.Int16(valueExpression))
	case *uint16:
		return f.toT(conv.Uint16(valueExpression))
	case *int32:
		return f.toT(conv.Int32(valueExpression))
	case *uint32:
		return f.toT(conv.Uint32(valueExpression))
	case *int64:
		return f.toT(conv.Int64(valueExpression))
	case *uint64:
		return f.toT(conv.Uint64(valueExpression))
	case *int:
		return f.toT(conv.Int(valueExpression))
	case *uint:
		return f.toT(conv.Uint(valueExpression))
	//case *uintptr:
	//	return f.toT(conv.Uintptr(valueExpression))
	case *float32:
		return f.toT(conv.Float32(valueExpression))
	case *float64:
		return f.toT(conv.Float64(valueExpression))
	//case *complex64:
	//	return f.toT(conv.Complex64(valueExpression))
	//case *complex128:
	//	return f.toT(conv.Complex128(valueExpression))
	case *string:
		return any(fmt.Sprintf("%v", valueExpression)).(T), nil
	}
	if valueExpression == nil {
		return f.toT(nil, nil)
	}
	return valueExpression.(T), nil
}

// EnsureType is used to convert from one (any,error) to (T, error)
func (f *FieldReaderVirtual[T]) toT(r any, err error) (T, error) {
	if err != nil {
		var zero T
		return zero, err
	}
	var zero T
	if r == nil {
		return zero, nil
	}
	return r.(T), nil
}
