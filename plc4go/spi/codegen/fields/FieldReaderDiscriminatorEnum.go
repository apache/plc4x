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

type FieldReaderDiscriminatorEnum[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderDiscriminatorEnum[T any](logger zerolog.Logger) *FieldReaderDiscriminatorEnum[T] {
	return &FieldReaderDiscriminatorEnum[T]{log: logger}
}

func (f *FieldReaderDiscriminatorEnum[T]) ReadDiscriminatorEnumField(ctx context.Context, logicalName, innerName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	var zero T
	if err := dataReader.PullContext(logicalName, utils.WithRenderAsList(true)); err != nil {
		return zero, errors.Wrap(err, "error pulling context")
	}
	v, err := f.SwitchParseByteOrderIfNecessary(ctx, func(ctx context.Context) (T, error) {
		return dataReader.Read(ctx, innerName, readerArgs...)
	}, dataReader, f.ExtractByteOrder(utils.UpcastReaderArgs(readerArgs...)...))
	if err != nil {
		return zero, errors.Wrap(err, "error parsing field")
	}
	if err := dataReader.CloseContext(logicalName, utils.WithRenderAsList(true)); err != nil {
		return zero, errors.Wrap(err, "error closing context")
	}
	return v, err
}
