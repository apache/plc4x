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

var APDUTypes map[readWriteModel.ApduType]func() interface{ Decode(Arg) error }

func init() {
	APDUTypes = map[readWriteModel.ApduType]func() interface{ Decode(Arg) error }{
		readWriteModel.ApduType_CONFIRMED_REQUEST_PDU: func() interface{ Decode(Arg) error } {
			pdu, _ := NewConfirmedRequestPDU(nil)
			return pdu
		},
		readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU: func() interface{ Decode(Arg) error } {
			pdu, _ := NewUnconfirmedRequestPDU(nil)
			return pdu
		},
		readWriteModel.ApduType_SIMPLE_ACK_PDU: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_COMPLEX_ACK_PDU: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_SEGMENT_ACK_PDU: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_ERROR_PDU: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_REJECT_PDU: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_ABORT_PDU: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_8: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_9: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_A: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_B: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_C: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_D: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_E: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_F: func() interface{ Decode(Arg) error } {
			panic("Implement Me")
		},
	}
}

var UnconfirmedRequestTypes map[readWriteModel.BACnetUnconfirmedServiceChoice]func() interface{ Decode(Arg) error }

func init() {
	UnconfirmedRequestTypes = map[readWriteModel.BACnetUnconfirmedServiceChoice]func() interface{ Decode(Arg) error }{
		readWriteModel.BACnetUnconfirmedServiceChoice_I_AM: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_I_HAVE: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_COV_NOTIFICATION: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_EVENT_NOTIFICATION: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_PRIVATE_TRANSFER: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_TEXT_MESSAGE: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_TIME_SYNCHRONIZATION: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_WHO_HAS: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_WHO_IS: func() interface{ Decode(Arg) error } {
			request, _ := NewWhoIsRequest()
			return request
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UTC_TIME_SYNCHRONIZATION: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_WRITE_GROUP: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_COV_NOTIFICATION_MULTIPLE: func() interface{ Decode(Arg) error } {
			panic("implement me")
		},
	}
}
