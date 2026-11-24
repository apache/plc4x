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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// A count from a negative signed field (OPC UA's -1 for a null array) reads as an empty array.
func TestReadCountArrayField_negativeCountIsEmpty(t *testing.T) {
	readBuffer := utils.NewReadBufferByteBased([]byte{1, 2, 3})
	count := int32(-1)
	result, err := ReadCountArrayField[byte](context.Background(), "array", io.ReadUnsignedByte(readBuffer, 8), uint64(count))
	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, uint32(0), readBuffer.GetPos(), "nothing is read")

	result, err = ReadCountArrayField[byte](context.Background(), "array", io.ReadUnsignedByte(readBuffer, 8), 2)
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2}, result)
}
