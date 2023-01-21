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
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type LocalDeviceObject struct {
	NumberOfAPDURetries       *uint
	APDUTimeout               *uint
	SegmentationSupported     *readWriteModel.BACnetSegmentation
	APDUSegmentTimeout        *uint
	MaxSegmentsAccepted       *readWriteModel.MaxSegmentsAccepted
	MaximumApduLengthAccepted *readWriteModel.MaxApduLengthAccepted
	App                       *Application
	ObjectName                string
	ObjectIdentifier          string
	VendorIdentifier          uint16
	ObjectList                []string
}

func (l *LocalDeviceObject) String() string {
	return fmt.Sprintf("LocalDeviceObject{NumberOfAPDURetries: %v, APDUTimeout: %v, SegmentationSupported: %v, APDUSegmentTimeout: %v, MaxSegmentsAccepted: %v, MaximumApduLengthAccepted: %v, ObjectName: %v, ObjectIdentifier: %v}", l.NumberOfAPDURetries, l.APDUTimeout, l.SegmentationSupported, l.APDUSegmentTimeout, l.MaxSegmentsAccepted, l.MaximumApduLengthAccepted, l.ObjectName, l.ObjectIdentifier)
}
