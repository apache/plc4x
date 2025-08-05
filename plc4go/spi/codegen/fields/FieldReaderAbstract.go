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

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// FieldReaderAbstract should actually never be used as an abstract field never gets read, it only makes sure
// abstract accessor methods are generated in the abstract parent type.
type FieldReaderAbstract[T any] struct {
	codegen.FieldCommons[T]

	log zerolog.Logger
}

func NewFieldReaderAbstract[T any](logger zerolog.Logger) *FieldReaderAbstract[T] {
	return &FieldReaderAbstract[T]{log: logger}
}

func (FieldReaderAbstract[T]) ReadAbstractField(context.Context, string, io.DataReader[T], ...utils.WithReaderArgs) (T, error) {
	var zero T
	return zero, nil
}
