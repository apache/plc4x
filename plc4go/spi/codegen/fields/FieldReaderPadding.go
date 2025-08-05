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

type FieldReaderPadding[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderPadding[T any](logger zerolog.Logger) *FieldReaderPadding[T] {
	return &FieldReaderPadding[T]{log: logger}
}

func (f *FieldReaderPadding[T]) ReadPaddingField(ctx context.Context, dataReader io.DataReader[T], timesPadding int, readerArgs ...utils.WithReaderArgs) error {
	f.log.Debug().Msg("reading padding field")
	if err := dataReader.PullContext("padding", utils.WithRenderAsList(true)); err != nil {
		return errors.Wrap(err, "error pulling context")
	}
	for timesPadding > 0 {
		// Just read the padding data and ignore it
		_, err := dataReader.Read(ctx, "value", readerArgs...)
		_ = err
		// Ignore ...
		// This could simply be that we're out of data to read for padding.
		// In protocols like the S7 protocol, this can happen if this is the
		// last field item, then the packet might end here.
		timesPadding--
	}
	if err := dataReader.CloseContext("padding", utils.WithRenderAsList(true)); err != nil {
		return errors.Wrap(err, "error closing context")
	}
	f.log.Debug().Msg("done reading padding field")
	return nil
}
