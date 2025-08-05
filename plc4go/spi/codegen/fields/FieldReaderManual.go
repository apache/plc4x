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

type FieldReaderManual[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderManual[T any](logger zerolog.Logger) *FieldReaderManual[T] {
	return &FieldReaderManual[T]{log: logger}
}

func (f *FieldReaderManual[T]) ReadManualField(ctx context.Context, logicalName string, readBuffer utils.ReadBuffer, parseSupplier func(context.Context) (T, error), readerArgs ...utils.WithReaderArgs) (T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	var zero T
	if err := readBuffer.PullContext(logicalName, utils.WithRenderAsList(true)); err != nil {
		return zero, errors.Wrap(err, "error pulling context")
	}
	v, err := f.SwitchParseByteOrderIfNecessary(ctx, parseSupplier, readBuffer, f.ExtractByteOrder(utils.UpcastReaderArgs(readerArgs...)...))
	if err != nil {
		return zero, errors.Wrap(err, "error parsing field")
	}
	if err := readBuffer.CloseContext(logicalName, utils.WithRenderAsList(true)); err != nil {
		return zero, errors.Wrap(err, "error closing context")
	}
	return v, err
}
