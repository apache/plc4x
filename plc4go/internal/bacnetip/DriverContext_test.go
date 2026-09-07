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

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func TestNewDriverContext_Defaults(t *testing.T) {
	dc := NewDriverContext(createDefaultConfiguration())
	assert.Equal(t, model.MaxApduLengthAccepted_NUM_OCTETS_1476, dc.maxApduLengthAccepted)
	assert.Equal(t, model.BACnetSegmentation_SEGMENTED_BOTH, dc.segmentation)
	assert.Equal(t, model.MaxSegmentsAccepted_NUM_SEGMENTS_16, dc.maxSegmentsAccepted)
	assert.True(t, dc.awaitSetupComplete)
	assert.True(t, dc.awaitDisconnectComplete)
}

func TestBytesToMaxApduLength(t *testing.T) {
	cases := []struct {
		in   uint16
		want model.MaxApduLengthAccepted
	}{
		{0, model.MaxApduLengthAccepted_MINIMUM_MESSAGE_SIZE},
		{127, model.MaxApduLengthAccepted_MINIMUM_MESSAGE_SIZE},
		{128, model.MaxApduLengthAccepted_NUM_OCTETS_128},
		{205, model.MaxApduLengthAccepted_NUM_OCTETS_128},
		{206, model.MaxApduLengthAccepted_NUM_OCTETS_206},
		{479, model.MaxApduLengthAccepted_NUM_OCTETS_206},
		{480, model.MaxApduLengthAccepted_NUM_OCTETS_480},
		{1023, model.MaxApduLengthAccepted_NUM_OCTETS_480},
		{1024, model.MaxApduLengthAccepted_NUM_OCTETS_1024},
		{1475, model.MaxApduLengthAccepted_NUM_OCTETS_1024},
		{1476, model.MaxApduLengthAccepted_NUM_OCTETS_1476},
		{65535, model.MaxApduLengthAccepted_NUM_OCTETS_1476},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, bytesToMaxApduLength(c.in), "bytesToMaxApduLength(%d)", c.in)
	}
}

func TestStringToSegmentation(t *testing.T) {
	cases := []struct {
		in   string
		want model.BACnetSegmentation
	}{
		{"segmented-both", model.BACnetSegmentation_SEGMENTED_BOTH},
		{"segmented-transmit", model.BACnetSegmentation_SEGMENTED_TRANSMIT},
		{"segmented-receive", model.BACnetSegmentation_SEGMENTED_RECEIVE},
		{"no-segmentation", model.BACnetSegmentation_NO_SEGMENTATION},
		{"", model.BACnetSegmentation_SEGMENTED_BOTH},
		{"bogus", model.BACnetSegmentation_SEGMENTED_BOTH},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, stringToSegmentation(c.in), "stringToSegmentation(%q)", c.in)
	}
}

func TestNumToMaxSegments(t *testing.T) {
	cases := []struct {
		in   uint8
		want model.MaxSegmentsAccepted
	}{
		{0, model.MaxSegmentsAccepted_UNSPECIFIED},
		{1, model.MaxSegmentsAccepted_NUM_SEGMENTS_02},
		{2, model.MaxSegmentsAccepted_NUM_SEGMENTS_02},
		{3, model.MaxSegmentsAccepted_NUM_SEGMENTS_04},
		{4, model.MaxSegmentsAccepted_NUM_SEGMENTS_04},
		{5, model.MaxSegmentsAccepted_NUM_SEGMENTS_08},
		{8, model.MaxSegmentsAccepted_NUM_SEGMENTS_08},
		{9, model.MaxSegmentsAccepted_NUM_SEGMENTS_16},
		{16, model.MaxSegmentsAccepted_NUM_SEGMENTS_16},
		{17, model.MaxSegmentsAccepted_NUM_SEGMENTS_32},
		{32, model.MaxSegmentsAccepted_NUM_SEGMENTS_32},
		{33, model.MaxSegmentsAccepted_NUM_SEGMENTS_64},
		{64, model.MaxSegmentsAccepted_NUM_SEGMENTS_64},
		{65, model.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS},
		{255, model.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, numToMaxSegments(c.in), "numToMaxSegments(%d)", c.in)
	}
}
