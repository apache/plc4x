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
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// DriverContext carries the resolved Configuration plus values derived from it
// (enum mappings, lifecycle flags) into the Connection layer. It is passed by
// value because all fields are small.
//
//go:generate go tool plc4xGenerator -type=DriverContext
type DriverContext struct {
	configuration Configuration

	// maxApduLengthAccepted is the Configuration.MaxApduLengthAccepted byte count
	// translated to its BACnet enum value (used in IAm responses and confirmed
	// requests). E.g. 1476 → MaxApduLengthAccepted_NUM_OCTETS_1476.
	maxApduLengthAccepted model.MaxApduLengthAccepted `stringer:"true"`

	// segmentation is Configuration.SegmentationSupported translated to its
	// BACnet enum. Defaults to SEGMENTED_BOTH for unknown strings.
	segmentation model.BACnetSegmentation `stringer:"true"`

	// maxSegmentsAccepted is Configuration.MaxSegmentsAccepted translated to its
	// BACnet enum value (e.g. 16 → NUM_SEGMENTS_16).
	maxSegmentsAccepted model.MaxSegmentsAccepted `stringer:"true"`

	awaitSetupComplete      bool
	awaitDisconnectComplete bool
}

func NewDriverContext(configuration Configuration) DriverContext {
	return DriverContext{
		configuration:           configuration,
		maxApduLengthAccepted:   bytesToMaxApduLength(configuration.MaxApduLengthAccepted),
		segmentation:            stringToSegmentation(configuration.SegmentationSupported),
		maxSegmentsAccepted:     numToMaxSegments(configuration.MaxSegmentsAccepted),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,
	}
}

// bytesToMaxApduLength picks the largest BACnet APDU-length enum whose byte count
// fits within the user-specified ceiling. A user value below the minimum (50) is
// clamped to MINIMUM_MESSAGE_SIZE.
func bytesToMaxApduLength(bytes uint16) model.MaxApduLengthAccepted {
	switch {
	case bytes >= 1476:
		return model.MaxApduLengthAccepted_NUM_OCTETS_1476
	case bytes >= 1024:
		return model.MaxApduLengthAccepted_NUM_OCTETS_1024
	case bytes >= 480:
		return model.MaxApduLengthAccepted_NUM_OCTETS_480
	case bytes >= 206:
		return model.MaxApduLengthAccepted_NUM_OCTETS_206
	case bytes >= 128:
		return model.MaxApduLengthAccepted_NUM_OCTETS_128
	default:
		return model.MaxApduLengthAccepted_MINIMUM_MESSAGE_SIZE
	}
}

func stringToSegmentation(s string) model.BACnetSegmentation {
	switch s {
	case "segmented-transmit":
		return model.BACnetSegmentation_SEGMENTED_TRANSMIT
	case "segmented-receive":
		return model.BACnetSegmentation_SEGMENTED_RECEIVE
	case "no-segmentation":
		return model.BACnetSegmentation_NO_SEGMENTATION
	default:
		return model.BACnetSegmentation_SEGMENTED_BOTH
	}
}

// numToMaxSegments picks the smallest BACnet MaxSegmentsAccepted enum whose count
// is greater than or equal to n. 0 is treated as "unspecified".
func numToMaxSegments(n uint8) model.MaxSegmentsAccepted {
	switch {
	case n == 0:
		return model.MaxSegmentsAccepted_UNSPECIFIED
	case n <= 2:
		return model.MaxSegmentsAccepted_NUM_SEGMENTS_02
	case n <= 4:
		return model.MaxSegmentsAccepted_NUM_SEGMENTS_04
	case n <= 8:
		return model.MaxSegmentsAccepted_NUM_SEGMENTS_08
	case n <= 16:
		return model.MaxSegmentsAccepted_NUM_SEGMENTS_16
	case n <= 32:
		return model.MaxSegmentsAccepted_NUM_SEGMENTS_32
	case n <= 64:
		return model.MaxSegmentsAccepted_NUM_SEGMENTS_64
	default:
		return model.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS
	}
}
