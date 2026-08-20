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

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldReaderReserved[T comparable] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderReserved[T comparable](logger zerolog.Logger) *FieldReaderReserved[T] {
	return &FieldReaderReserved[T]{log: logger}
}

// ReadReservedField reads a field the protocol specification reserves for a fixed value.
//
// A value other than the expected one is reported and kept, not treated as a parse failure. Devices
// do put unexpected bytes in reserved fields, and refusing the frame over one of them loses
// everything that follows - for a record stream, every later record too. This is also what plc4j
// does (FieldReaderReserved.java logs and returns).
//
// The value is handed back so the caller can record what arrived, but note it does not survive a
// round trip: the generated serializer writes the specified constant, not the stored value, so
// re-serializing a message normalizes the field. That is the same in plc4j.
//
// A genuine read error is still an error. It used to be swallowed whenever the failed read left the
// zero value and the reserved field expected zero, which is the common case - so a truncated buffer
// silently parsed on with the remaining fields reading nonsense.
func (f *FieldReaderReserved[T]) ReadReservedField(ctx context.Context, logicalName string, dataReader io.DataReader[T], referenceValue T, readerArgs ...utils.WithReaderArgs) (*T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	value, err := dataReader.Read(ctx, logicalName, readerArgs...)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading reserved field %s", logicalName)
	}
	if value != referenceValue {
		f.log.Info().
			Str("logicalName", logicalName).
			Any("expected", referenceValue).
			Any("actual", value).
			Msg("Unexpected value in reserved field, keeping what the device sent")
		return &value, nil
	}
	return nil, nil
}
