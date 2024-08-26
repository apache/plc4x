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

	"github.com/pkg/errors"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type APCI interface {
	PCI

	GetApduInvokeID() *uint8

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setAPDU(readWriteModel.APDU)
	getAPDU() readWriteModel.APDU
}

type _APCI struct {
	*_PCI
	*DebugContents

	apdu readWriteModel.APDU // TODO: check if this is part of the _APCI or _APDU
}

var _ APCI = (*_APCI)(nil)

func NewAPCI(pduUserData spi.Message, apdu readWriteModel.APDU) APCI {
	a := &_APCI{
		apdu: apdu,
	}
	a._PCI = newPCI(pduUserData, nil, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
	switch apdu := pduUserData.(type) {
	case readWriteModel.APDUExactly:
		a.apdu = apdu
	}
	return a
}

func (n *_APCI) GetApduInvokeID() *uint8 {
	if n.apdu == nil {
		return nil
	}
	switch apdu := n.apdu.(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		invokeId := apdu.GetInvokeId()
		return &invokeId
	default:
		return nil
	}
}

// Deprecated: check if needed as we do it in update
func (n *_APCI) setAPDU(apdu readWriteModel.APDU) {
	n.apdu = apdu
}

func (n *_APCI) getAPDU() readWriteModel.APDU {
	return n.apdu
}

func (n *_APCI) Update(apci Arg) error {
	if err := n._PCI.Update(apci); err != nil {
		return errors.Wrap(err, "error updating _PCI")
	}
	switch apci := apci.(type) {
	case APDU:
		n.apdu = apci.getAPDU()
		// TODO: update coordinates...
		return nil
	default:
		return errors.Errorf("invalid APCI type %T", apci)
	}
}

func (n *_APCI) Encode(pdu Arg) error {
	if err := pdu.(interface{ Update(Arg) error }).Update(n); err != nil { // TODO: better validate that arg is really PDUData... use switch similar to Update
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (n *_APCI) Decode(pdu Arg) error {
	if err := n._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (n *_APCI) deepCopy() *_APCI {
	return &_APCI{_PCI: n._PCI.deepCopy(), apdu: n.apdu}
}

type APDU interface {
	readWriteModel.APDU
	APCI
	PDUData
}

type __APDU struct {
	*_APCI
	*_PDUData
}

var _ APDU = (*__APDU)(nil)

func NewAPDU() (APDU, error) {
	a := &__APDU{}

	a._APCI = NewAPCI(nil, nil).(*_APCI)
	a._PDUData = NewPDUData(NoArgs).(*_PDUData)
	return a, nil
}

func (a *__APDU) Encode(pdu Arg) error {
	panic("implement me")
}

func (a *__APDU) Decode(pdu Arg) error {
	panic("implement me")
}

func (a *__APDU) GetApduType() readWriteModel.ApduType {
	if a.apdu == nil {
		return 0xf
	}
	return a.apdu.GetApduType()
}

func (a *__APDU) deepCopy() *__APDU {
	return &__APDU{_APCI: a._APCI.deepCopy(), _PDUData: a._PDUData.deepCopy()}
}

func (a *__APDU) DeepCopy() PDU {
	return a.deepCopy()
}

func (a *__APDU) String() string {
	return fmt.Sprintf("__APDU{%s}", a._PCI)
}

type _APDU interface {
	APDU
}

type ___APDU struct {
	*__APDU
}

var _ _APDU = (*___APDU)(nil)

func new_APDU() (_APDU, error) {
	i := &___APDU{}
	var err error
	apdu, err := NewAPDU()
	if err != nil {
		return nil, errors.Wrap(err, "error creating APDU")
	}
	i.__APDU = apdu.(*__APDU)
	return i, nil
}

// TODO: implement it...
type ConfirmedRequestPDU struct {
	*___APDU
}

// TODO: implement it...
type UnconfirmedRequestPDU struct {
	*___APDU
}

func NewUnconfirmedRequestPDU() (*UnconfirmedRequestPDU, error) {
	u := &UnconfirmedRequestPDU{}
	apdu, err := new_APDU()
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APDU")
	}
	u.__APDU = apdu.(*__APDU)
	return u, nil
}

// TODO: implement it...
type SimpleAckPDU struct {
	*___APDU
}

// TODO: implement it...
type ComplexAckPDU struct {
	*___APDU
}

// TODO: implement it...
type SegmentAckPDU struct {
	*___APDU
}

// TODO: implement it...
type ErrorPDU struct {
	*___APDU
}

// TODO: implement it...
type RejectReason struct {
	*Enumerated
}

// TODO: implement it...
type RejectPDU struct {
	*___APDU
}

// TODO: implement it...
type AbortReason struct {
	*Enumerated
}

// TODO: implement it...
type AbortPDU struct {
	*___APDU
}

// TODO: implement it...
type APCISequence struct {
	*_APCI
	*Sequence

	tagList *TagList
}

func NewAPCISequence() *APCISequence {
	a := &APCISequence{}
	a._APCI = NewAPCI(nil, nil).(*_APCI) // TODO: what to pass up?
	a.Sequence = NewSequence()

	// start with an empty tag list
	a.tagList = NewTagList(nil)
	return a
}

func (a *APCISequence) Encode(apdu Arg) error {
	switch apdu := apdu.(type) {
	case APDU:
		if err := apdu.Update(a); err != nil {
			return errors.Wrap(err, "error updating APDU")
		}

		// create a tag list
		a.tagList = NewTagList(nil)
		if err := a.Sequence.Encode(a.tagList); err != nil {
			return errors.Wrap(err, "error encoding TagList")
		}

		// encode the tag list
		a.tagList.Encode(apdu)

		apdu.setAPDU(a.apdu)
		return nil
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
}

func (a *APCISequence) Decode(apdu Arg) error {
	switch apdu := apdu.(type) {
	case APDU:
		if err := a.Update(apdu); err != nil {
			return errors.Wrap(err, "error updating APDU")
		}
		switch pduUserData := apdu.GetRootMessage().(type) {
		case readWriteModel.APDUExactly:
			a.tagList = NewTagList(nil)
			if err := a.tagList.Decode(apdu); err != nil {
				return errors.Wrap(err, "error decoding TagList")
			}
			// pass the taglist to the Sequence for additional decoding
			if err := a.Sequence.Decode(a.tagList); err != nil {
				return errors.Wrap(err, "error encoding TagList")
			}

			_ = pduUserData
		}
		return nil
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
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

func NewUnconfirmedRequestSequence() (*UnconfirmedRequestSequence, error) {
	u := &UnconfirmedRequestSequence{}
	u.APCISequence = NewAPCISequence()
	var err error
	u.UnconfirmedRequestPDU, err = NewUnconfirmedRequestPDU()
	if err != nil {
		return nil, errors.Wrap(err, "error creating UnconfirmedRequestPDU")
	}
	return u, nil
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
	*UnconfirmedRequestSequence
}

func NewWhoIsRequest() (*WhoIsRequest, error) {
	w := &WhoIsRequest{}
	var err error
	w.UnconfirmedRequestSequence, err = NewUnconfirmedRequestSequence()
	if err != nil {
		return nil, errors.Wrap(err, "error creating UnconfirmedRequestSequence")
	}
	return w, nil
}

func (r *WhoIsRequest) String() string {
	return fmt.Sprintf("WhoIsRequest{%s}", r.UnconfirmedRequestSequence)
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

func (r *ConfirmedPrivateTransferRequest) String() string {
	return fmt.Sprintf("ConfirmedPrivateTransferRequest{%s}", r.ConfirmedRequestSequence)
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
