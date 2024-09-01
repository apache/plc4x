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

type FieldWriterManualArray[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldWriterManualArray[T any](logger zerolog.Logger) *FieldWriterManualArray[T] {
	return &FieldWriterManualArray[T]{
		log: logger,
	}
}

func (f *FieldWriterManualArray[T]) WriteManualArrayField(ctx context.Context, logicalName string, values []T, consumer func(context.Context, utils.WriteBuffer, T) error, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	f.log.Debug().Str("logicalName", logicalName).Msg("write field")
	return f.SwitchParseByteOrderIfNecessarySerializeWrapped(ctx, func(ctx context.Context) error {
		if values != nil {
			if err := writeBuffer.PushContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
			for _, value := range values {
				if err := consumer(ctx, writeBuffer, value); err != nil {
					return errors.Wrapf(err, "error writing value %s for %v", logicalName, value)
				}
			}

			if err := writeBuffer.PopContext(logicalName, utils.WithRenderAsList(true)); err != nil {
				return errors.Wrap(err, "error pushing context for "+logicalName)
			}
		}
		return nil
	}, writeBuffer, f.ExtractByteOrder(utils.UpcastWriterArgs(writerArgs...)...))
}
