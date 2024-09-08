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
	io2 "io"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldReaderOptional[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderOptional[T any](logger zerolog.Logger) *FieldReaderOptional[T] {
	return &FieldReaderOptional[T]{log: logger}
}

func (f *FieldReaderOptional[T]) ReadOptionalField(ctx context.Context, logicalName string, dataReader io.DataReader[T], condition bool, readerArgs ...utils.WithReaderArgs) (*T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	if !condition {
		f.log.Debug().Str("logicalName", logicalName).Msg("Condition doesn't match for field")
		return nil, nil
	}

	// TODO: add hex support

	currentPos := dataReader.GetPos()
	optionalValue, err := dataReader.Read(ctx, logicalName, readerArgs...)
	switch {
	case errors.Is(err, utils.ParseAssertError{}):
		f.log.Debug().Err(err).Str("logicalName", logicalName).Uint16("oldPos", currentPos).Msg("Assertion doesn't match for field. Resetting read position to oldPos")
		dataReader.SetPos(currentPos)
		return nil, nil
	case errors.Is(err, io2.EOF):
		f.log.Debug().Err(err).Str("logicalName", logicalName).Uint16("oldPos", currentPos).Msg("Not enough bytes. Resetting read position to oldPos")
		dataReader.SetPos(currentPos)
		return nil, nil
	case err != nil:
		return nil, errors.Wrapf(err, "Error parsing '%s' field", logicalName)
	default:
		// All good
		f.log.Debug().Str("logicalName", logicalName).Msg("done reading field")
		return &optionalValue, nil
	}
}
