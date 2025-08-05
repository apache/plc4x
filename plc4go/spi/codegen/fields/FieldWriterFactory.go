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
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func WriteSimpleTypeArrayField[T any](ctx context.Context, logicalName string, value []T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterArray[T, spi.Message](log).WriteSimpleTypeArrayField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteComplexTypeArrayField[T spi.Message](ctx context.Context, logicalName string, value []T, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterArray[any, T](log).WriteComplexTypeArrayField(ctx, logicalName, value, writeBuffer, writerArgs...)
}

func WriteByteArrayField(ctx context.Context, logicalName string, value []byte, dataWriter io.DataWriter[[]byte], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterArray[any, spi.Message](log).WriteByteArrayField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteChecksumField[T any](ctx context.Context, logicalName string, valueProducer func() (T, error), dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	value, err := valueProducer()
	if err != nil {
		return errors.Wrap(err, "error producing value")
	}
	return NewFieldWriterChecksum[T](log).WriteChecksumField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteConstField[T any](ctx context.Context, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterConst[T](log).WriteConstField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteEnumField[T any](ctx context.Context, logicalName, innerName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterEnum[T](log).WriteEnumField(ctx, logicalName, innerName, value, dataWriter, writerArgs...)
}

func WriteDiscriminatorField[T any](ctx context.Context, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterDiscriminator[T](log).WriteDiscriminatorField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteDiscriminatorEnumField[T any](ctx context.Context, logicalName, innerName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterDiscriminatorEnum[T](log).WriteDiscriminatorEnumField(ctx, logicalName, innerName, value, dataWriter, writerArgs...)
}

func WriteImplicitField[T any](ctx context.Context, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterImplicit[T](log).WriteImplicitField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteManualField[T any](ctx context.Context, logicalName string, runnable func(context.Context) error, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterManual[T](log).WriteManualField(ctx, logicalName, runnable, writeBuffer, writerArgs...)
}

func WriteManualArrayField[T any](ctx context.Context, logicalName string, values []T, runnable func(context.Context, utils.WriteBuffer, T) error, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterManualArray[T](log).WriteManualArrayField(ctx, logicalName, values, runnable, writeBuffer, writerArgs...)
}

func WriteOptionalField[T any](ctx context.Context, logicalName string, value *T, dataWriter io.DataWriter[T], condition bool, writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterOptional[T](log).WriteOptionalField(ctx, logicalName, value, dataWriter, condition, writerArgs...)
}

func WriteOptionalEnumField[T any](ctx context.Context, logicalName, innerName string, value *T, dataWriter io.DataWriter[T], condition bool, writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterOptionalEnum[T](log).WriteOptionalEnumField(ctx, logicalName, innerName, value, dataWriter, condition, writerArgs...)
}

func WritePaddingField[T any](ctx context.Context, logicalName string, timesPadding int, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterPadding[T](log).WritePaddingField(ctx, logicalName, timesPadding, value, dataWriter, writerArgs...)
}

func WriteReservedField[T any](ctx context.Context, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterReserved[T](log).WriteReservedField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteSimpleField[T any](ctx context.Context, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterSimple[T](log).WriteSimpleField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteSimpleEnumField[T any](ctx context.Context, logicalName, innerName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldWriterSimpleEnum[T](log).WriteSimpleEnumField(ctx, logicalName, innerName, value, dataWriter, writerArgs...)
}
