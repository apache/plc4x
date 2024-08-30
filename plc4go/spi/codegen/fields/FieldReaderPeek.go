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
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldReaderPeek[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderPeek[T any](logger zerolog.Logger) *FieldReaderPeek[T] {
	return &FieldReaderPeek[T]{log: logger}
}

func (f *FieldReaderPeek[T]) ReadPeekField(ctx context.Context, logicalName string, dataReader io.DataReader[T], offset int, readerArgs ...utils.WithReaderArgs) (T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	var zero T

	currentPos := dataReader.GetPos()
	// TODO: implement offset. We either need to pass the readBuffer or add a instruction to the dataReader

	field, err := f.SwitchParseByteOrderIfNecessary(ctx, func(ctx context.Context) (T, error) {
		return dataReader.Read(ctx, logicalName, readerArgs...)
	}, dataReader, f.ExtractByteOrder(utils.UpcastReaderArgs(readerArgs...)...))
	switch {
	case errors.Is(err, utils.ParseAssertError{}):
		f.log.Debug().Err(err).Str("logicalName", logicalName).Uint16("oldPos", currentPos).Msg("Peek failed for field. Resetting read position to oldPos")
	case err != nil:
		return zero, errors.Wrapf(err, "Error parsing '%s' field", logicalName)
	default:
		// All good
	}
	dataReader.SetPos(currentPos)
	f.log.Debug().Str("logicalName", logicalName).Msg("done reading field")
	return field, nil
}
