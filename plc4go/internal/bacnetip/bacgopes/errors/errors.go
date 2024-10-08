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

package errors

import (
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// ConfigurationError This error is raised when there is a configuration problem such as bindings between layers or
//
//	required parameters that are missing.
type ConfigurationError struct {
	ValueError
}

// EncodingError this error is raised if there is a problem during encoding.
type EncodingError struct {
	ValueError
}

// DecodingError this error is raised if there is a problem during decoding.
type DecodingError struct {
	ValueError
}

// ExecutionError This error is raised for if there is an error during the execution of
//
//	a service or function at the application layer of stack and the error
//	translated into an ErrorPDU.
type ExecutionError struct {
	ErrorClass readWriteModel.ErrorClass
	ErrorCode  readWriteModel.ErrorCode
	Message    string
}

var _ error = ExecutionError{}

func (e ExecutionError) Error() string {
	// TODO: check if we want to return class and code
	return e.Message
}

// RejectException Exceptions in this family correspond to reject reasons.  If the
//
//	application raises one of these errors while processing a confirmed
//	service request, the stack will form an appropriate RejectPDU and
//	send it to the client.
type RejectException struct {
	Exception
	// TODO: add reject reason check
	// rejectReason = nil
	args Args
}

// RejectOther Generated in response to a confirmed request APDU that contains a
//
//	syntax error for which an error code has not been explicitly defined.
type RejectOther struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'other'
}

// RejectBufferOverflow A buffer capacity has been exceeded.
type RejectBufferOverflow struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'bufferOverflow'
}

// InconsistentParameters Generated in response to a confirmed request APDU that omits a
//
//	conditional service argument that should be present or contains a
//	conditional service argument that should not be present. This condition
//	could also elicit a Reject PDU with a Reject Reason of INVALID_TAG.
type InconsistentParameters struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'inconsistentParameters'
}

// InvalidParameterDatatype Generated in response to a confirmed request APDU in which the encoding
//
//	of one or more of the service parameters does not follow the correct type
//	specification. This condition could also elicit a Reject PDU with a Reject
//	Reason of INVALID_TAG.
type InvalidParameterDatatype struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'invalidParameterDatatype'
}

// InvalidTag While parsing a message, an invalid tag was encountered. Since an
//
//	invalid tag could confuse the parsing logic, any of the following Reject
//	Reasons may also be generated in response to a confirmed request
//	containing an invalid tag: INCONSISTENT_PARAMETERS,
//	INVALID_PARAMETER_DATA_TYPE, MISSING_REQUIRED_PARAMETER, and
//	TOO_MANY_ARGUMENTS.
type InvalidTag struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'invalidTag'
}

// MissingRequiredParameter Generated in response to a confirmed request APDU that is missing at
//
//	least one mandatory service argument. This condition could also elicit a
//	Reject PDU with a Reject Reason of INVALID_TAG.
type MissingRequiredParameter struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'missingRequiredParameter'
}

// ParameterOutOfRange Generated in response to a confirmed request APDU that conveys a
//
//	parameter whose value is outside the range defined for this service.
type ParameterOutOfRange struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'parameterOutOfRange'
}

// TooManyArguments Generated in response to a confirmed request APDU in which the total
//
//	number of service arguments is greater than specified for the service.
//	This condition could also elicit a Reject PDU with a Reject Reason of
//	INVALID_TAG.
type TooManyArguments struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'tooManyArguments'
}

// UndefinedEnumeration Generated in response to a confirmed request APDU in which one or
//
//	more of the service parameters is decoded as an enumeration that is not
//	defined by the type specification of this parameter.
type UndefinedEnumeration struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'undefinedEnumeration'
}

// UnrecognizedService Generated in response to a confirmed request APDU in which the Service
//
//	Choice field specifies an unknown or unsupported service.
type UnrecognizedService struct {
	RejectException
	// TODO: add reject reason
	// rejectReason = 'unrecognizedService'
}

// AbortException Exceptions in this family correspond to abort reasons.  If the
//
//	application raises one of these errors while processing a confirmed
//	service request, the stack will form an appropriate AbortPDU and
//	send it to the client.
type AbortException struct {
	Exception
	// TODO: add reject reason check
	// abortReason = nil
	args Args
}

// AbortOther This abort reason is returned for a reason other than any of those
//
//	for which an error code has not been explicitly defined.
type AbortOther struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'other'
}

// AbortBufferOverflow A buffer capacity has been exceeded.
type AbortBufferOverflow struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'bufferOverflow'
}

// InvalidAPDUInThisState Generated in response to an APDU that is not expected in the present
//
//	state of the Transaction State Machine.
type InvalidAPDUInThisState struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'invalidApduInThisState'
}

// PreemptedByHigherPriorityTask The transaction shall be aborted to permit higher priority processing.
type PreemptedByHigherPriorityTask struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'preemptedByHigherPriorityTask'
}

// SegmentationNotSupported Generated in response to an APDU that has its segmentation bit set to
//
//	TRUE when the receiving device does not support segmentation. It is also
//	generated when a BACnet-ComplexACK-PDU is large enough to require
//	segmentation but it cannot be transmitted because either the transmitting
//	device or the receiving device does not support segmentation.
type SegmentationNotSupported struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'segmentationNotSupported'
}

// SecurityError The Transaction is aborted due to receipt of a security error.
type SecurityError struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'securityError'
}

// InsufficientSecurity The transaction is aborted due to receipt of a PDU secured differently
//
//	than the original PDU of the transaction.
type InsufficientSecurity struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'insufficientSecurity'
}

// WindowSizeOutOfRange A device receives a request that is segmented, or receives any segment
//
//	of a segmented request, where the Proposed Window Size field of the PDU
//	header is either zero or greater than 127.
type WindowSizeOutOfRange struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'windowSizeOutOfRange'
}

// ApplicationExceededReplyTime A device receives a confirmed request but its application layer has not
//
//	responded within the published APDU Timeout period.
type ApplicationExceededReplyTime struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'applicationExceededReplyTime'
}

// OutOfResources A device receives a request but cannot start processing because it has
//
//	run out of some internal resource.
type OutOfResources struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'outOfResources'
}

// TSMTimeout A transaction state machine timer exceeded the timeout applicable for
//
//	the current state, causing the transaction machine to abort the
//	transaction.
type TSMTimeout struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'tsmTimeout'
}

// APDUTooLong An APDU was received from the local application program whose overall
//
//	size exceeds the maximum transmittable length or exceeds the maximum number
//	of segments accepted by the server.
type APDUTooLong struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'apduTooLong'
}

// ServerTimeout BACgopes specific.
type ServerTimeout struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'serverTimeout'
}

// NoResponse BACgopes specific.
type NoResponse struct {
	AbortException
	// TODO: add reject reason check
	// abortReason = 'noResponse'
}
