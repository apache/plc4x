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

import readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

type WhoIsIAmServices struct {
}

func NewWhoIsIAmServices() (*WhoIsIAmServices, error) {
	// TODO: implement me
	return nil, nil
}

var defaultMaxApduLength = readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1024
var defaultMaxSegmentsAccepted = readWriteModel.MaxSegmentsAccepted_NUM_SEGMENTS_16

// _LocalDeviceObjectDefault is a device entry with default entries
var _LocalDeviceObjectDefault = LocalDeviceObject{
	MaximumApduLengthAccepted: &defaultMaxApduLength,
	SegmentationSupported:     readWriteModel.BACnetSegmentation_SEGMENTED_BOTH,
	MaxSegmentsAccepted:       &defaultMaxSegmentsAccepted,
	APDUSegmentTimeout:        5000,
	APDUTimeout:               3000,
	NumberOfAPDURetries:       3,
}

type LocalDeviceObject struct {
	NumberOfAPDURetries       uint
	APDUTimeout               uint
	SegmentationSupported     readWriteModel.BACnetSegmentation
	APDUSegmentTimeout        uint
	MaxSegmentsAccepted       *readWriteModel.MaxSegmentsAccepted
	MaximumApduLengthAccepted *readWriteModel.MaxApduLengthAccepted
}

func NewLocalDeviceObject() *LocalDeviceObject {
	return &LocalDeviceObject{
		NumberOfAPDURetries:       _LocalDeviceObjectDefault.NumberOfAPDURetries,
		APDUTimeout:               _LocalDeviceObjectDefault.APDUTimeout,
		SegmentationSupported:     _LocalDeviceObjectDefault.SegmentationSupported,
		APDUSegmentTimeout:        _LocalDeviceObjectDefault.APDUSegmentTimeout,
		MaxSegmentsAccepted:       _LocalDeviceObjectDefault.MaxSegmentsAccepted,
		MaximumApduLengthAccepted: _LocalDeviceObjectDefault.MaximumApduLengthAccepted,
	}
}
