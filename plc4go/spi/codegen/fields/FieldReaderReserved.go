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

type FieldReaderReserved[T comparable] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderReserved[T comparable](logger zerolog.Logger) *FieldReaderReserved[T] {
	return &FieldReaderReserved[T]{log: logger}
}

func (f *FieldReaderReserved[T]) ReadReservedField(ctx context.Context, logicalName string, dataReader io.DataReader[T], referenceValue T, readerArgs ...utils.WithReaderArgs) (*T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	value, err := dataReader.Read(ctx, logicalName, readerArgs...)
	if value != referenceValue {
		return &value, utils.ParseAssertError{Message: fmt.Sprintf("Expected constant value %v but got %v for reserved field", referenceValue, value), Err: err}
	}
	return nil, nil
}
