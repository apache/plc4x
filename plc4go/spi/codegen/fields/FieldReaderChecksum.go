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

type FieldReaderChecksum[T comparable] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderChecksum[T comparable](logger zerolog.Logger) *FieldReaderChecksum[T] {
	return &FieldReaderChecksum[T]{log: logger}
}

func (f *FieldReaderChecksum[T]) ReadChecksumField(ctx context.Context, logicalName string, dataReader io.DataReader[T], referenceValue T, readerArgs ...utils.WithReaderArgs) (T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	checksumValue, err := f.SwitchParseByteOrderIfNecessary(ctx, func(ctx context.Context) (T, error) {
		return dataReader.Read(ctx, logicalName, readerArgs...)
	}, dataReader, f.ExtractByteOrder(utils.UpcastReaderArgs(readerArgs...)...))
	if checksumValue != referenceValue {
		return checksumValue, utils.ParseAssertError{Message: fmt.Sprintf("Checksum value %v: doesn't match expected %v", checksumValue, referenceValue), Err: err}
	}
	return checksumValue, nil
}
