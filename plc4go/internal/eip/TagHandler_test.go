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

package eip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
)

func TestParseTagSimpleNoType(t *testing.T) {
	// %rate without type → defaults to DINT, elementNb 1
	tag, err := NewTagHandler().ParseTag("%rate")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "rate", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_DINT, plcTag.GetType())
	assert.Equal(t, uint16(1), plcTag.GetElementNb())
}

func TestParseTagWithDINT(t *testing.T) {
	// %rate:DINT → DINT, elementNb 1
	tag, err := NewTagHandler().ParseTag("%rate:DINT")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "rate", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_DINT, plcTag.GetType())
	assert.Equal(t, uint16(1), plcTag.GetElementNb())
}

func TestParseTagArray(t *testing.T) {
	// %arr[0]:DINT:4 → DINT, elementNb 4, tag preserved as arr[0]
	tag, err := NewTagHandler().ParseTag("%arr[0]:DINT:4")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "arr[0]", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_DINT, plcTag.GetType())
	assert.Equal(t, uint16(4), plcTag.GetElementNb())
}

func TestParseTagWithDottedMember(t *testing.T) {
	// %struct.member:INT → INT, elementNb 1
	tag, err := NewTagHandler().ParseTag("%struct.member:INT")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "struct.member", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_INT, plcTag.GetType())
	assert.Equal(t, uint16(1), plcTag.GetElementNb())
}

func TestParseTagWithREAL(t *testing.T) {
	// %f:REAL → REAL, elementNb 1
	tag, err := NewTagHandler().ParseTag("%f:REAL")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "f", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_REAL, plcTag.GetType())
	assert.Equal(t, uint16(1), plcTag.GetElementNb())
}

func TestParseTagWithBOOL(t *testing.T) {
	// %b:BOOL → BOOL, elementNb 1
	tag, err := NewTagHandler().ParseTag("%b:BOOL")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "b", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_BOOL, plcTag.GetType())
	assert.Equal(t, uint16(1), plcTag.GetElementNb())
}

func TestParseTagWithSTRING(t *testing.T) {
	// %s:STRING → STRING, elementNb 1
	tag, err := NewTagHandler().ParseTag("%s:STRING")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "s", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_STRING, plcTag.GetType())
	assert.Equal(t, uint16(1), plcTag.GetElementNb())
}

func TestParseTagUnknownType(t *testing.T) {
	// %x:UNKNOWNTYPE → error
	_, err := NewTagHandler().ParseTag("%x:UNKNOWNTYPE")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown data type")
}

func TestParseTagGarbageInput(t *testing.T) {
	// Invalid prefix (not %): error
	_, err := NewTagHandler().ParseTag("rate:DINT")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tag address")
}

func TestParseTagEmptyString(t *testing.T) {
	// Empty string: error
	_, err := NewTagHandler().ParseTag("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tag address")
}

func TestParseTagWithTypeButNoElements(t *testing.T) {
	// %rate:LINT with no element count → LINT, elementNb defaults to 1
	tag, err := NewTagHandler().ParseTag("%rate:LINT")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, "rate", plcTag.GetTag())
	assert.Equal(t, readWriteModel.CIPDataTypeCode_LINT, plcTag.GetType())
	assert.Equal(t, uint16(1), plcTag.GetElementNb())
}

func TestParseTagLargeElementCount(t *testing.T) {
	// %arr[0]:DINT:1000 → elementNb 1000
	tag, err := NewTagHandler().ParseTag("%arr[0]:DINT:1000")
	require.NoError(t, err)
	plcTag := tag.(PlcTag)
	assert.Equal(t, readWriteModel.CIPDataTypeCode_DINT, plcTag.GetType())
	assert.Equal(t, uint16(1000), plcTag.GetElementNb())
}
