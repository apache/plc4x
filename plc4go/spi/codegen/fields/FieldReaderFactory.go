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

	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func ReadAbstractField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderAbstract[T](log).ReadAbstractField(ctx, logicalName, dataReader, readerArgs...)
}

func ReadCountArrayField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], count uint64, readerArgs ...utils.WithReaderArgs) ([]T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderArray[T](log).ReadFieldCount(ctx, logicalName, dataReader, count, readerArgs...)
}

// ReadLengthArrayField
//
//	 In some protocols a long is used as length, but we simply can't address that many bytes,
//		so we simply cast it down to int as on java we couldn't even read more bytes as MAX-INT.
//
//		@param logicalName the logical name of this field
//		@param dataReader  the dataReader used to retrieve this field
//		@param length      the length of the array
//		@param readerArgs  optional read args
//		@param [T]         the type of the array elements
//		@return the read length array
//		@return err if something went wrong parsing
func ReadLengthArrayField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], length int, readerArgs ...utils.WithReaderArgs) ([]T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderArray[T](log).ReadFieldLength(ctx, logicalName, dataReader, length, readerArgs...)
}

func ReadTerminatedArrayField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], termination func() bool, readerArgs ...utils.WithReaderArgs) ([]T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderArray[T](log).ReadFieldTerminated(ctx, logicalName, dataReader, termination, readerArgs...)
}

func ReadAssertField[T comparable](ctx context.Context, logicalName string, dataReader io.DataReader[T], expectedValue T, readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderAssert[T](log).ReadAssertField(ctx, logicalName, dataReader, expectedValue, readerArgs...)
}

func ReadChecksumField[T comparable](ctx context.Context, logicalName string, dataReader io.DataReader[T], expectedValueProducer func() (expectedValue T, err error), readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	var zero T
	expectedValue, err := expectedValueProducer()
	if err != nil {
		return zero, err
	}
	return NewFieldReaderChecksum[T](log).ReadChecksumField(ctx, logicalName, dataReader, expectedValue, readerArgs...)
}

func ReadConstField[T comparable](ctx context.Context, logicalName string, dataReader io.DataReader[T], expectedValue T, readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderConst[T](log).ReadConstField(ctx, logicalName, dataReader, expectedValue, readerArgs...)
}

func ReadDiscriminatorField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderDiscriminator[T](log).ReadDiscriminatorField(ctx, logicalName, dataReader, readerArgs...)
}

func ReadDiscriminatorEnumField[T any](ctx context.Context, logicalName, innerName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderDiscriminatorEnum[T](log).ReadDiscriminatorEnumField(ctx, logicalName, innerName, dataReader, readerArgs...)
}

func ReadEnumField[T any](ctx context.Context, logicalName, innerName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderEnum[T](log).ReadEnumField(ctx, logicalName, innerName, dataReader, readerArgs...)
}

func ReadImplicitField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderImplicit[T](log).ReadImplicitField(ctx, logicalName, dataReader, readerArgs...)
}

func ReadOptionalField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], condition bool, readerArgs ...utils.WithReaderArgs) (*T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderOptional[T](log).ReadOptionalField(ctx, logicalName, dataReader, condition, readerArgs...)
}

func ReadManualByteArrayField(ctx context.Context, logicalName string, readBuffer utils.ReadBuffer, termination func([]byte) bool, parse func(context.Context) (byte, error), readerArgs ...utils.WithReaderArgs) ([]byte, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderManualArray[byte](log).ReadManualByteArrayField(ctx, logicalName, readBuffer, termination, parse, readerArgs...)
}

func ReadManualArrayField[T any](ctx context.Context, logicalName string, readBuffer utils.ReadBuffer, termination func([]T) bool, parse func(context.Context) (T, error), readerArgs ...utils.WithReaderArgs) ([]T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderManualArray[T](log).ReadManualArrayField(ctx, logicalName, readBuffer, termination, parse, readerArgs...)
}

func ReadManualField[T any](ctx context.Context, logicalName string, readBuffer utils.ReadBuffer, parseSupplier func(context.Context) (T, error), readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderManual[T](log).ReadManualField(ctx, logicalName, readBuffer, parseSupplier, readerArgs...)
}

func ReadPaddingField[T any](ctx context.Context, dataReader io.DataReader[T], timesPadding int, readerArgs ...utils.WithReaderArgs) error {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderPadding[T](log).ReadPaddingField(ctx, dataReader, timesPadding, readerArgs...)
}

func ReadReservedField[T comparable](ctx context.Context, logicalName string, dataReader io.DataReader[T], expectedValue T, readerArgs ...utils.WithReaderArgs) (*T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderReserved[T](log).ReadReservedField(ctx, logicalName, dataReader, expectedValue, readerArgs...)
}

func ReadSimpleField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderSimple[T](log).ReadSimpleField(ctx, logicalName, dataReader, readerArgs...)
}

func ReadUnknownField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderUnknown[T](log).ReadUnknownField(ctx, logicalName, dataReader, readerArgs...)
}

func ReadVirtualField[T any](ctx context.Context, logicalName string, klass any, valueExpression any, readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderVirtual[T](log).ReadVirtualField(ctx, logicalName, klass, valueExpression, readerArgs...)
}

func ReadPeekField[T any](ctx context.Context, logicalName string, dataReader io.DataReader[T], offset int, readerArgs ...utils.WithReaderArgs) (T, error) {
	log := *zerolog.Ctx(ctx)
	return NewFieldReaderPeek[T](log).ReadPeekField(ctx, logicalName, dataReader, offset, readerArgs...)
}
