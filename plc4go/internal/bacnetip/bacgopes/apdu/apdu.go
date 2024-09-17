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

package apdu

import (
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

var _debug = CreateDebugPrinter()

var APDUTypes map[readWriteModel.ApduType]func() Decoder

func init() {
	APDUTypes = map[readWriteModel.ApduType]func() Decoder{
		readWriteModel.ApduType_CONFIRMED_REQUEST_PDU: func() Decoder {
			pdu, _ := NewConfirmedRequestPDU(Nothing())
			return pdu
		},
		readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU: func() Decoder {
			pdu, _ := NewUnconfirmedRequestPDU(Nothing())
			return pdu
		},
		readWriteModel.ApduType_SIMPLE_ACK_PDU: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_COMPLEX_ACK_PDU: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_SEGMENT_ACK_PDU: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_ERROR_PDU: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_REJECT_PDU: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_ABORT_PDU: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_8: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_9: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_A: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_B: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_C: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_D: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_E: func() Decoder {
			panic("Implement Me")
		},
		readWriteModel.ApduType_APDU_UNKNOWN_F: func() Decoder {
			panic("Implement Me")
		},
	}
}

var UnconfirmedRequestTypes map[readWriteModel.BACnetUnconfirmedServiceChoice]func() Decoder

func init() {
	UnconfirmedRequestTypes = map[readWriteModel.BACnetUnconfirmedServiceChoice]func() Decoder{
		readWriteModel.BACnetUnconfirmedServiceChoice_I_AM: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_I_HAVE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_COV_NOTIFICATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_EVENT_NOTIFICATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_PRIVATE_TRANSFER: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_TEXT_MESSAGE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_TIME_SYNCHRONIZATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_WHO_HAS: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_WHO_IS: func() Decoder {
			request, _ := NewWhoIsRequest(Nothing())
			return request
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UTC_TIME_SYNCHRONIZATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_WRITE_GROUP: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetUnconfirmedServiceChoice_UNCONFIRMED_COV_NOTIFICATION_MULTIPLE: func() Decoder {
			panic("implement me")
		},
	}
}

var ConfirmedRequestTypes map[readWriteModel.BACnetConfirmedServiceChoice]func() Decoder

func init() {
	ConfirmedRequestTypes = map[readWriteModel.BACnetConfirmedServiceChoice]func() Decoder{
		readWriteModel.BACnetConfirmedServiceChoice_ACKNOWLEDGE_ALARM: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_CONFIRMED_COV_NOTIFICATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_CONFIRMED_COV_NOTIFICATION_MULTIPLE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_CONFIRMED_EVENT_NOTIFICATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_GET_ALARM_SUMMARY: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_GET_ENROLLMENT_SUMMARY: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_GET_EVENT_INFORMATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_LIFE_SAFETY_OPERATION: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_SUBSCRIBE_COV: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_SUBSCRIBE_COV_PROPERTY: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_SUBSCRIBE_COV_PROPERTY_MULTIPLE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_ATOMIC_READ_FILE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_ATOMIC_WRITE_FILE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_ADD_LIST_ELEMENT: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_REMOVE_LIST_ELEMENT: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_CREATE_OBJECT: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_DELETE_OBJECT: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_READ_PROPERTY: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_READ_PROPERTY_MULTIPLE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_READ_RANGE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_WRITE_PROPERTY: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_WRITE_PROPERTY_MULTIPLE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_DEVICE_COMMUNICATION_CONTROL: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_CONFIRMED_PRIVATE_TRANSFER: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_CONFIRMED_TEXT_MESSAGE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_REINITIALIZE_DEVICE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_VT_OPEN: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_VT_CLOSE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_VT_DATA: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_AUTHENTICATE: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_REQUEST_KEY: func() Decoder {
			panic("implement me")
		},
		readWriteModel.BACnetConfirmedServiceChoice_READ_PROPERTY_CONDITIONAL: func() Decoder {
			panic("implement me")
		},
	}
}
