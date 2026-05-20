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

package bacnetip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func TestTagHandler_ParseSimpleAddress(t *testing.T) {
	h := NewTagHandler()
	tag, err := h.ParseTag("ANALOG_INPUT,1/PRESENT_VALUE")
	require.NoError(t, err)
	bt, ok := tag.(BacNetPlcTag)
	require.True(t, ok)
	require.NotNil(t, bt.GetObjectId().ObjectIdType)
	assert.Equal(t, readWriteModel.BACnetObjectType_ANALOG_INPUT, *bt.GetObjectId().ObjectIdType)
	assert.Equal(t, uint32(1), bt.GetObjectId().ObjectIdInstance)
	require.Len(t, bt.GetProperties(), 1)
	require.NotNil(t, bt.GetProperties()[0].PropertyIdentifier)
	assert.Equal(t, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, *bt.GetProperties()[0].PropertyIdentifier)
}

func TestTagHandler_ParseNumericObjectAndProperty(t *testing.T) {
	// Numeric object type maps to ObjectIdTypeProprietary; numeric property
	// maps to PropertyIdentifierProprietary.
	h := NewTagHandler()
	tag, err := h.ParseTag("12345,7/4242")
	require.NoError(t, err)
	bt := tag.(BacNetPlcTag)
	require.Nil(t, bt.GetObjectId().ObjectIdType)
	require.NotNil(t, bt.GetObjectId().ObjectIdTypeProprietary)
	assert.Equal(t, uint16(12345), *bt.GetObjectId().ObjectIdTypeProprietary)
	assert.Equal(t, uint32(7), bt.GetObjectId().ObjectIdInstance)
	require.Len(t, bt.GetProperties(), 1)
	require.Nil(t, bt.GetProperties()[0].PropertyIdentifier)
	require.NotNil(t, bt.GetProperties()[0].PropertyIdentifierProprietary)
	assert.Equal(t, uint32(4242), *bt.GetProperties()[0].PropertyIdentifierProprietary)
}

func TestTagHandler_ParseWithArrayIndex(t *testing.T) {
	h := NewTagHandler()
	tag, err := h.ParseTag("ANALOG_VALUE,2/PRESENT_VALUE[3]")
	require.NoError(t, err)
	props := tag.(BacNetPlcTag).GetProperties()
	require.Len(t, props, 1)
	require.NotNil(t, props[0].ArrayIndex)
	assert.Equal(t, uint(3), *props[0].ArrayIndex)
}

func TestTagHandler_ParseMultipleProperties(t *testing.T) {
	h := NewTagHandler()
	tag, err := h.ParseTag("ANALOG_INPUT,1/PRESENT_VALUE&OBJECT_NAME&UNITS")
	require.NoError(t, err)
	props := tag.(BacNetPlcTag).GetProperties()
	require.Len(t, props, 3)
	assert.Equal(t, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, *props[0].PropertyIdentifier)
	assert.Equal(t, readWriteModel.BACnetPropertyIdentifier_OBJECT_NAME, *props[1].PropertyIdentifier)
	assert.Equal(t, readWriteModel.BACnetPropertyIdentifier_UNITS, *props[2].PropertyIdentifier)
}

func TestTagHandler_ParseMultiplePropertiesWithArrayIndex(t *testing.T) {
	h := NewTagHandler()
	tag, err := h.ParseTag("ANALOG_VALUE,5/PRESENT_VALUE[1]&PRIORITY_ARRAY[16]")
	require.NoError(t, err)
	props := tag.(BacNetPlcTag).GetProperties()
	require.Len(t, props, 2)
	require.NotNil(t, props[0].ArrayIndex)
	assert.Equal(t, uint(1), *props[0].ArrayIndex)
	require.NotNil(t, props[1].ArrayIndex)
	assert.Equal(t, uint(16), *props[1].ArrayIndex)
}

func TestTagHandler_ParseInvalidString(t *testing.T) {
	h := NewTagHandler()
	cases := []string{
		"",
		"garbage",
		"ANALOG_INPUT,1",       // missing property
		"ANALOG_INPUT/X",       // missing instance
		"ANALOG_INPUT,abc/X",   // non-numeric instance
		",1/X",                 // missing object type
	}
	for _, c := range cases {
		_, err := h.ParseTag(c)
		assert.Error(t, err, "input %q should fail", c)
	}
}

func TestTagHandler_ParseUnknownObjectType(t *testing.T) {
	h := NewTagHandler()
	_, err := h.ParseTag("NOT_AN_OBJECT,1/PRESENT_VALUE")
	assert.Error(t, err)
}

func TestTagHandler_ParseUnknownPropertyType(t *testing.T) {
	h := NewTagHandler()
	_, err := h.ParseTag("ANALOG_INPUT,1/NOT_A_PROPERTY")
	assert.Error(t, err)
}

func TestTagHandler_ParseQueryNotSupported(t *testing.T) {
	h := NewTagHandler()
	_, err := h.ParseQuery("anything")
	assert.Error(t, err)
}
