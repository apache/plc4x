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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldReaderManualArray[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderManualArray[T any](log zerolog.Logger) *FieldReaderManualArray[T] {
	return &FieldReaderManualArray[T]{
		log: log,
	}
}

func (f *FieldReaderManualArray[T]) ReadManualByteArrayField(ctx context.Context, logicalName string, readBuffer utils.ReadBuffer, termination func([]byte) bool, parse func(context.Context) (byte, error), readerArgs ...utils.WithReaderArgs) ([]byte, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	// Ensure we have the render as list argument present
	readerArgs = append(readerArgs, utils.WithRenderAsList(true))
	if err := readBuffer.PullContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error pulling context for %s", logicalName)
	}
	var result = make([]byte, 0)
	for !termination(result) {
		read, err := parse(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading item")
		}
		a := any(read)
		elems, ok := a.(byte)
		if !ok {
			return nil, errors.Wrapf(err, "parse supplier didn't return bool. Was %T", a)
		}
		result = append(result, elems)
	}
	if err := readBuffer.CloseContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error closing context for %s", logicalName)
	}
	f.log.Debug().Str("logicalName", logicalName).Msg("done reading field")
	return result, nil
}

func (f *FieldReaderManualArray[T]) ReadManualArrayField(ctx context.Context, logicalName string, readBuffer utils.ReadBuffer, termination func([]T) bool, parse func(context.Context) (T, error), readerArgs ...utils.WithReaderArgs) ([]T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	// Ensure we have the render as list argument present
	readerArgs = append(readerArgs, utils.WithRenderAsList(true))
	if err := readBuffer.PullContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error pulling context for %s", logicalName)
	}
	var result = make([]T, 0)
	for !termination(result) {
		read, err := parse(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading item")
		}
		a := any(read)
		elems, ok := a.(T)
		if !ok {
			var t T
			return nil, errors.Wrapf(err, "parse supplier didn't return %T. Was %T", t, a)
		}
		result = append(result, elems)
	}
	if err := readBuffer.CloseContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error closing context for %s", logicalName)
	}
	f.log.Debug().Str("logicalName", logicalName).Msg("done reading field")
	return result, nil
}
