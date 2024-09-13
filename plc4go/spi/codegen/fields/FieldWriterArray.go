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

type FieldWriterArray[T any, C spi.Message] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldWriterArray[T any, C spi.Message](logger zerolog.Logger) *FieldWriterArray[T, C] {
	return &FieldWriterArray[T, C]{
		log: logger,
	}
}

func (f *FieldWriterArray[_, _]) WriteByteArrayField(ctx context.Context, logicalName string, values []byte, dataWriter io.DataWriter[[]byte], writerArgs ...utils.WithWriterArgs) error {
	f.log.Debug().Str("logicalName", logicalName).Msg("write field")
	if values != nil && len(values) > 0 {
		return dataWriter.Write(ctx, logicalName, values, writerArgs...)
	}
	return nil
}

func (f *FieldWriterArray[T, _]) WriteSimpleTypeArrayField(ctx context.Context, logicalName string, values []T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	f.log.Debug().Str("logicalName", logicalName).Msg("write field")
	return f.SwitchParseByteOrderIfNecessarySerializeWrapped(ctx, func(ctx context.Context) error {
		if values != nil && len(values) > 0 {
			if err := dataWriter.PushContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
			for curItem := 0; curItem < len(values); curItem++ {
				value := values[curItem]
				ctx := codegen.NewContextLastItem(ctx, curItem == len(values)-1)
				ctx = utils.CreateArrayContext(ctx, len(values), curItem)
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

func (f *FieldWriterArray[_, C]) WriteComplexTypeArrayField(ctx context.Context, logicalName string, values []C, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	f.log.Debug().Str("logicalName", logicalName).Msg("write field")
	return f.SwitchParseByteOrderIfNecessarySerializeWrapped(ctx, func(ctx context.Context) error {
		if values != nil && len(values) > 0 {
			if err := writeBuffer.PushContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
			for curItem := 0; curItem < len(values); curItem++ {
				value := values[curItem]
				ctx := codegen.NewContextLastItem(ctx, curItem == len(values)-1)
				ctx = utils.CreateArrayContext(ctx, len(values), curItem)
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
