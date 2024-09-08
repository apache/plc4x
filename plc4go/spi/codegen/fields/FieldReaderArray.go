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
	"math"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type FieldReaderArray[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderArray[T any](log zerolog.Logger) *FieldReaderArray[T] {
	return &FieldReaderArray[T]{
		log: log,
	}
}

func (f *FieldReaderArray[T]) ReadFieldCount(ctx context.Context, logicalName string, dataReader io.DataReader[T], count uint64, readerArgs ...utils.WithReaderArgs) ([]T, error) {
	f.log.Debug().Str("logicalName", logicalName).Uint64("count", count).Msg("reading field")
	if count > math.MaxInt64 {
		return nil, errors.Errorf("Array count of %d exceeds the maximum allowed count of %d", count, math.MaxInt64)
	}
	// Ensure we have the render as list argument present
	readerArgs = append(readerArgs, utils.WithRenderAsList(true))
	if err := dataReader.PullContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error pulling context for %s", logicalName)
	}
	itemCount := int(max(0, count))
	var result = make([]T, itemCount)
	if itemCount == 0 {
		result = nil
	}
	for curItem := 0; curItem < itemCount; curItem++ {
		// Make some variables available that would be otherwise challenging to forward.
		ctx := codegen.NewContextCurItem(ctx, curItem)
		ctx = codegen.NewContextLastItem(ctx, curItem == itemCount-1)
		ctx = utils.CreateArrayContext(ctx, itemCount, curItem)
		read, err := dataReader.Read(ctx, "", readerArgs...)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading item %d", curItem)
		}
		result[curItem] = read
	}
	if err := dataReader.CloseContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error closing context for %s", logicalName)
	}
	f.log.Debug().Str("logicalName", logicalName).Msg("done reading field")
	return result, nil
}

func (f *FieldReaderArray[T]) ReadFieldLength(ctx context.Context, logicalName string, dataReader io.DataReader[T], length int, readerArgs ...utils.WithReaderArgs) ([]T, error) {
	f.log.Debug().Str("logicalName", logicalName).Int("length", length).Msg("reading field")
	// Ensure we have the render as list argument present
	readerArgs = append(readerArgs, utils.WithRenderAsList(true))
	if err := dataReader.PullContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error pulling context for %s", logicalName)
	}
	startPos := int(dataReader.GetPos())
	var result = make([]T, 0)
	numberOfElements := 0
	stopPosition := int(startPos) + length
	f.log.Debug().Int("startPos", startPos).Int("stopPosition", stopPosition).Msg("start reading at pos startPos while < stopPosition")
	for int(dataReader.GetPos()) < stopPosition {
		numberOfElements++
		read, err := dataReader.Read(ctx, "", readerArgs...)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading item %d", dataReader.GetPos())
		}
		result = append(result, read)
	}
	if err := dataReader.CloseContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error closing context for %s", logicalName)
	}
	f.log.Debug().Str("logicalName", logicalName).Msg("done reading field")
	return result, nil
}

func (f *FieldReaderArray[T]) ReadFieldTerminated(ctx context.Context, logicalName string, dataReader io.DataReader[T], termination func() bool, readerArgs ...utils.WithReaderArgs) ([]T, error) {
	f.log.Debug().Str("logicalName", logicalName).Msg("reading field")
	// Ensure we have the render as list argument present
	readerArgs = append(readerArgs, utils.WithRenderAsList(true))
	if err := dataReader.PullContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error pulling context for %s", logicalName)
	}
	var result = make([]T, 0)
	for !termination() {
		read, err := dataReader.Read(ctx, "", readerArgs...)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading item")
		}
		result = append(result, read)
	}
	if err := dataReader.CloseContext(logicalName, readerArgs...); err != nil {
		return nil, errors.Wrapf(err, "error closing context for %s", logicalName)
	}
	f.log.Debug().Str("logicalName", logicalName).Msg("done reading field")
	return result, nil
}
