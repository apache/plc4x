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

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldReaderConst[T comparable] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderConst[T comparable](logger zerolog.Logger) *FieldReaderConst[T] {
	return &FieldReaderConst[T]{log: logger}
}

func (f *FieldReaderConst[T]) ReadConstField(ctx context.Context, logicalName string, dataReader io.DataReader[T], expectedValue T, readerArgs ...utils.WithReaderArgs) (T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	value, err := f.SwitchParseByteOrderIfNecessary(ctx, func(ctx context.Context) (T, error) {
		return dataReader.Read(ctx, logicalName, readerArgs...)
	}, dataReader, f.ExtractByteOrder(utils.UpcastReaderArgs(readerArgs...)...))
	if value != expectedValue {
		return value, utils.ParseAssertError{Message: fmt.Sprintf("Checksum value %v: doesn't match expected %v", value, expectedValue), Err: err}
	}
	return value, nil
}
