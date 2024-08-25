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

type APCI interface {
	PCI
}

type _APCI struct {
	*_PCI
	*DebugContents
}

type APDU interface {
	APCI
	PDUData

	GetApduInvokeID() *uint8 // TODO: check if we really need a pointer
}

type __APDU struct {
	*_APCI
}

type _APDU interface {
	APDU
}

type ___APDU struct {
	*__APDU
}

// TODO: implement it...
type ConfirmedRequestPDU struct {
	*__APDU
}

// TODO: implement it...
type UnconfirmedRequestPDU struct {
	*__APDU
}

// TODO: implement it...
type SimpleAckPDU struct {
	*__APDU
}

// TODO: implement it...
type ComplexAckPDU struct {
	*__APDU
}

// TODO: implement it...
type SegmentAckPDU struct {
	*__APDU
}

// TODO: implement it...
type ErrorPDU struct {
	*__APDU
}

// TODO: implement it...
type RejectReason struct {
	*Enumerated
}

// TODO: implement it...
type RejectPDU struct {
	*__APDU
}

// TODO: implement it...
type AbortReason struct {
	*Enumerated
}

// TODO: implement it...
type AbortPDU struct {
	*__APDU
}

// TODO: implement it...
type APCISequence struct {
	*_APCI
	*Sequence
}

// TODO: implement it...
type ConfirmedRequestSequence struct {
	*APCISequence
	*ConfirmedRequestPDU
}

// TODO: implement it...
type ComplexAckSequence struct {
	*APCISequence
	*ComplexAckPDU
}

// TODO: implement it...
type UnconfirmedRequestSequence struct {
	*APCISequence
	*UnconfirmedRequestPDU
}

// TODO: implement it...
type ErrorSequence struct {
	*APCISequence
	*ErrorPDU
}

// TODO: implement it...
type Error struct {
}

// TODO: implement it...
type ChangeListError struct {
}

// TODO: implement it...
type CreateObjectError struct {
}

// TODO: implement it...
type ConfirmedPrivateTransferError struct {
}

// TODO: implement it...
type WritePropertyMultipleError struct {
}

// TODO: implement it...
type VTCloseError struct {
}

// TODO: implement it...
type ReadPropertyRequest struct {
}

// TODO: implement it...
type ReadPropertyACK struct {
}

// TODO: implement it...
type ReadAccessSpecification struct {
}

// TODO: implement it...
type ReadPropertyMultipleRequest struct {
}

// TODO: implement it...
type ReadAccessResultElementChoice struct {
}

// TODO: implement it...
type ReadAccessResultElement struct {
}

// TODO: implement it...
type ReadAccessResult struct {
}

// TODO: implement it...
type ReadPropertyMultipleACK struct {
}

// TODO: implement it...
type RangeByPosition struct {
}

// TODO: implement it...
type RangeBySequenceNumber struct {
}

// TODO: implement it...
type RangeByTime struct {
}

// TODO: implement it...
type Range struct {
}

// TODO: implement it...
type ReadRangeRequest struct {
}

// TODO: implement it...
type ReadRangeACK struct {
}

// TODO: implement it...
type WritePropertyRequest struct {
}

// TODO: implement it...
type WriteAccessSpecification struct {
}

// TODO: implement it...
type WritePropertyMultipleRequest struct {
}

// TODO: implement it...
type GroupChannelValue struct {
}

// TODO: implement it...
type WriteGroupRequest struct {
}

// TODO: implement it...
type IAmRequest struct {
}

// TODO: implement it...
type IHaveRequest struct {
}

// TODO: implement it...
type WhoHasLimits struct {
}

// TODO: implement it...
type WhoHasObject struct {
}

// TODO: implement it...
type WhoHasRequest struct {
}

// TODO: implement it...
type WhoIsRequest struct {
}

// TODO: implement it...
type EventNotificationParameters struct {
}

// TODO: implement it...
type ConfirmedEventNotificationRequest struct {
}

// TODO: implement it...
type UnconfirmedEventNotificationRequest struct {
}

// TODO: implement it...
type COVNotificationParameters struct {
}

// TODO: implement it...
type ConfirmedCOVNotificationRequest struct {
}

// TODO: implement it...
type UnconfirmedCOVNotificationRequest struct {
}

// TODO: implement it...
type UnconfirmedPrivateTransferRequest struct {
}

// TODO: implement it...
type UnconfirmedTextMessageRequestMessageClass struct {
}

// TODO: implement it...
type UnconfirmedTextMessageRequestMessagePriority struct {
}

// TODO: implement it...
type UnconfirmedTextMessageRequest struct {
}

// TODO: implement it...
type TimeSynchronizationRequest struct {
}

// TODO: implement it...
type UTCTimeSynchronizationRequest struct {
}

// TODO: implement it...
type AcknowledgeAlarmRequest struct {
}

// TODO: implement it...
type GetAlarmSummaryRequest struct {
}

// TODO: implement it...
type GetAlarmSummaryAlarmSummary struct {
}

// TODO: implement it...
type GetAlarmSummaryACK struct {
}

// TODO: implement it...
type GetEnrollmentSummaryRequestAcknowledgmentFilterType struct {
}

// TODO: implement it...
type GetEnrollmentSummaryRequestEventStateFilterType struct {
}

// TODO: implement it...
type GetEnrollmentSummaryRequestPriorityFilterType struct {
}

// TODO: implement it...
type GetEnrollmentSummaryRequest struct {
}

// TODO: implement it...
type GetEnrollmentSummaryEnrollmentSummary struct {
}

// TODO: implement it...
type GetEnrollmentSummaryACK struct {
}

// TODO: implement it...
type GetEventInformationRequest struct {
}

// TODO: implement it...
type GetEventInformationEventSummary struct {
}

// TODO: implement it...
type GetEventInformationACK struct {
}

// TODO: implement it...
type LifeSafetyOperationRequest struct {
}

// TODO: implement it...
type SubscribeCOVRequest struct {
}

// TODO: implement it...
type SubscribeCOVPropertyRequest struct {
}

// TODO: implement it...
type AtomicReadFileRequestAccessMethodChoiceStreamAccess struct {
}

// TODO: implement it...
type AtomicReadFileRequestAccessMethodChoiceRecordAccess struct {
}

// TODO: implement it...
type AtomicReadFileRequestAccessMethodChoice struct {
}

// TODO: implement it...
type AtomicReadFileRequest struct {
}

// TODO: implement it...
type AtomicReadFileACKAccessMethodStreamAccess struct {
}

// TODO: implement it...
type AtomicReadFileACKAccessMethodRecordAccess struct {
}

// TODO: implement it...
type AtomicReadFileACKAccessMethodChoice struct {
}

// TODO: implement it...
type AtomicReadFileACK struct {
}

// TODO: implement it...
type AtomicWriteFileRequestAccessMethodChoiceStreamAccess struct {
}

// TODO: implement it...
type AtomicWriteFileRequestAccessMethodChoiceRecordAccess struct {
}

// TODO: implement it...
type AtomicWriteFileRequestAccessMethodChoice struct {
}

// TODO: implement it...
type AtomicWriteFileRequest struct {
}

// TODO: implement it...
type AtomicWriteFileACK struct {
}

// TODO: implement it...
type AddListElementRequest struct {
}

// TODO: implement it...
type CreateObjectRequestObjectSpecifier struct {
}

// TODO: implement it...
type CreateObjectRequest struct {
}

// TODO: implement it...
type CreateObjectACK struct {
}

// TODO: implement it...
type DeleteObjectRequest struct {
}

// TODO: implement it...
type RemoveListElementRequest struct {
}

// TODO: implement it...
type DeviceCommunicationControlRequestEnableDisable struct {
}

// TODO: implement it...
type DeviceCommunicationControlRequest struct {
}

// TODO: implement it...
type ConfirmedPrivateTransferRequest struct {
	*ConfirmedRequestSequence
}

func NewConfirmedPrivateTransferRequest() *ConfirmedPrivateTransferRequest {
	c := &ConfirmedPrivateTransferRequest{}
	panic("implement me")
	return c
}

// TODO: implement it...
type ConfirmedPrivateTransferACK struct {
}

// TODO: implement it...
type ConfirmedTextMessageRequestMessageClass struct {
}

// TODO: implement it...
type ConfirmedTextMessageRequestMessagePriority struct {
}

// TODO: implement it...
type ConfirmedTextMessageRequest struct {
}

// TODO: implement it...
type ReinitializeDeviceRequestReinitializedStateOfDevice struct {
}

// TODO: implement it...
type ReinitializeDeviceRequest struct {
}

// TODO: implement it...
type VTOpenRequest struct {
}

// TODO: implement it...
type VTOpenACK struct {
}

// TODO: implement it...
type VTCloseRequest struct {
}

// TODO: implement it...
type VTDataRequest struct {
}

// TODO: implement it...
type VTDataACK struct {
}

// TODO: implement it...
type AuthenticateRequest struct {
}

// TODO: implement it...
type AuthenticateACK struct {
}

// TODO: implement it...
type RequestKeyRequest struct {
}

// TODO: implement it...
type ConfirmedServiceChoice struct {
}

// TODO: implement it...
type UnconfirmedServiceChoice struct {
}
