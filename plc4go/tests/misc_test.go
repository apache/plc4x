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

package tests

import (
	"encoding/binary"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByteOrderUpcast(t *testing.T) {
	t.Run("little endian", func(t *testing.T) {
		rw := utils.UpcastReaderArgs(codegen.WithByteOrder(binary.LittleEndian))
		fc := codegen.FieldCommons[any]{}
		byteOrder := fc.ExtractByteOrder(rw...)
		require.NotNil(t, byteOrder)
		assert.Equal(t, binary.LittleEndian, *byteOrder)
	})
	t.Run("big endian", func(t *testing.T) {
		rw := utils.UpcastReaderArgs(codegen.WithByteOrder(binary.BigEndian))
		fc := codegen.FieldCommons[any]{}
		byteOrder := fc.ExtractByteOrder(rw...)
		require.NotNil(t, byteOrder)
		assert.Equal(t, binary.BigEndian, *byteOrder)
	})
}
