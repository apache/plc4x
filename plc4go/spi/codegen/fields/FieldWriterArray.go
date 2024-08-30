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

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldWriterArray[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldWriterArray[T any](logger zerolog.Logger) *FieldWriterArray[T] {
	return &FieldWriterArray[T]{
		log: logger,
	}
}

func (f *FieldWriterArray[T]) WriteByteArrayField(ctx context.Context, logicalName string, values []byte, dataWriter io.DataWriter[[]byte], writerArgs ...utils.WithWriterArgs) error {
	f.log.Debug().Str("logicalName", logicalName).Msg("write field")
	if values != nil {
		return dataWriter.Write(ctx, logicalName, values, writerArgs...)
	}
	return nil
}

func (f *FieldWriterArray[T]) WriteSimpleTypeArrayField(ctx context.Context, logicalName string, values []T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	f.log.Debug().Str("logicalName", logicalName).Msg("write field")
	return f.SwitchParseByteOrderIfNecessarySerializeWrapped(ctx, func(ctx context.Context) error {
		if values != nil {
			if err := dataWriter.PushContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
			for curItem := 0; curItem < len(values); curItem++ {
				value := values[curItem]
				ctx := codegen.NewContextLastItem(ctx, curItem == len(values)-1)
				if err := dataWriter.Write(ctx, "value", value, writerArgs...); err != nil {
					return errors.Wrapf(err, "error writing value %s for %v", logicalName, value)
				}
			}
			if err := dataWriter.PopContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
		}
		return nil
	}, dataWriter, f.ExtractByteOrder(utils.UpcastWriterArgs(writerArgs...)...))
}

func (f *FieldWriterArray[T]) WriteComplexTypeArrayField(ctx context.Context, logicalName string, values []spi.Message, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	f.log.Debug().Str("logicalName", logicalName).Msg("write field")
	return f.SwitchParseByteOrderIfNecessarySerializeWrapped(ctx, func(ctx context.Context) error {
		if values != nil {
			if err := writeBuffer.PushContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
			for curItem := 0; curItem < len(values); curItem++ {
				value := values[curItem]
				ctx := codegen.NewContextLastItem(ctx, curItem == len(values)-1)
				if err := value.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
					return errors.Wrapf(err, "error writing value %v for %s", logicalName, value)
				}
			}
			if err := writeBuffer.PopContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
		}
		return nil
	}, writeBuffer, f.ExtractByteOrder(utils.UpcastWriterArgs(writerArgs...)...))
}
