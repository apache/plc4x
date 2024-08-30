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

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func WriteSimpleTypeArrayField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value []T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterArray[T](log).WriteSimpleTypeArrayField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteComplexTypeArrayField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value []spi.Message, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterArray[spi.Message](log).WriteComplexTypeArrayField(ctx, logicalName, value, writeBuffer, writerArgs...)
}

func WriteByteArrayField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value []byte, dataWriter io.DataWriter[[]byte], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterArray[T](log).WriteByteArrayField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteChecksumField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterChecksum[T](log).WriteChecksumField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteConstField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterConst[T](log).WriteConstField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteEnumField[T any](ctx context.Context, log zerolog.Logger, logicalName, innerName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterEnum[T](log).WriteEnumField(ctx, logicalName, innerName, value, dataWriter, writerArgs...)
}

func WriteDiscriminatorField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterDiscriminator[T](log).WriteDiscriminatorField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteDiscriminatorEnumField[T any](ctx context.Context, log zerolog.Logger, logicalName, innerName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterDiscriminatorEnum[T](log).WriteDiscriminatorEnumField(ctx, logicalName, innerName, value, dataWriter, writerArgs...)
}

func WriteImplicitField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterImplicit[T](log).WriteImplicitField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteManualField[T any](ctx context.Context, log zerolog.Logger, logicalName string, runnable func(ctx context.Context) error, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterManual[T](log).WriteManualField(ctx, logicalName, runnable, writeBuffer, writerArgs...)
}

func WriteManualArrayField[T any](ctx context.Context, log zerolog.Logger, logicalName string, values []T, runnable func(ctx context.Context) error, writeBuffer utils.WriteBuffer, writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterManualArray[T](log).WriteManualArrayField(ctx, logicalName, values, runnable, writeBuffer, writerArgs...)
}

func WriteOptionalField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value T, dataWriter io.DataWriter[T], condition bool, writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterOptional[T](log).WriteOptionalField(ctx, logicalName, value, dataWriter, condition, writerArgs...)
}

func WriteOptionalEnumField[T any](ctx context.Context, log zerolog.Logger, logicalName, innerName string, value T, dataWriter io.DataWriter[T], condition bool, writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterOptionalEnum[T](log).WriteOptionalEnumField(ctx, logicalName, innerName, value, dataWriter, condition, writerArgs...)
}

func WritePaddingField[T any](ctx context.Context, log zerolog.Logger, logicalName string, timesPadding int, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterPadding[T](log).WritePaddingField(ctx, logicalName, timesPadding, value, dataWriter, writerArgs...)
}

func WriteReservedField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterReserved[T](log).WriteReservedField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteSimpleField[T any](ctx context.Context, log zerolog.Logger, logicalName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterSimple[T](log).WriteSimpleField(ctx, logicalName, value, dataWriter, writerArgs...)
}

func WriteSimpleEnumField[T any](ctx context.Context, log zerolog.Logger, logicalName, innerName string, value T, dataWriter io.DataWriter[T], writerArgs ...utils.WithWriterArgs) error {
	return NewFieldWriterSimpleEnum[T](log).WriteSimpleEnumField(ctx, logicalName, innerName, value, dataWriter, writerArgs...)
}
