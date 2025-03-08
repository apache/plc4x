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

package appservice

import (
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// NOTE: needs to reside here otherwise there is a circular dependency

//go:generate go tool plc4xGenerator -type=DeviceInfo -prefix=app_
type DeviceInfo struct {
	DeviceIdentifier readWriteModel.BACnetTagPayloadObjectIdentifier
	Address          Address

	MaximumApduLengthAccepted *readWriteModel.MaxApduLengthAccepted `stringer:"true"`
	SegmentationSupported     *readWriteModel.BACnetSegmentation    `stringer:"true"`
	MaxSegmentsAccepted       *readWriteModel.MaxSegmentsAccepted   `stringer:"true"`
	VendorId                  *readWriteModel.BACnetVendorId        `stringer:"true"`
	MaximumNpduLength         *uint

	_refCount int
	_cacheKey DeviceInfoCacheKey
}

func NewDeviceInfo(deviceIdentifier readWriteModel.BACnetTagPayloadObjectIdentifier, address Address) DeviceInfo {
	return DeviceInfo{
		DeviceIdentifier: deviceIdentifier,
		Address:          address,

		MaximumApduLengthAccepted: func() *readWriteModel.MaxApduLengthAccepted {
			octets1024 := readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1024
			return &octets1024
		}(),
		SegmentationSupported: func() *readWriteModel.BACnetSegmentation {
			noSegmentation := readWriteModel.BACnetSegmentation_NO_SEGMENTATION
			return &noSegmentation
		}(),
	}
}
