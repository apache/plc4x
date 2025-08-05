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

package basetypes

import (
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type ErrorCode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewErrorCode(arg Arg) (*ErrorCode, error) {
	s := &ErrorCode{
		enumerations: map[string]uint64{"abortApduTooLong": 123,
			"abortApplicationExceededReplyTime":  124,
			"abortBufferOverflow":                51,
			"abortInsufficientSecurity":          135,
			"abortInvalidApduInThisState":        52,
			"abortOther":                         56,
			"abortOutOfResources":                125,
			"abortPreemptedByHigherPriorityTask": 53,
			"abortProprietary":                   55,
			"abortSecurityError":                 136,
			"abortSegmentationNotSupported":      54,
			"abortTsmTimeout":                    126,
			"abortWindowSizeOutOfRange":          127,
			"accessDenied":                       85,
			"addressingError":                    115,
			"badDestinationAddress":              86,
			"badDestinationDeviceId":             87,
			"badSignature":                       88,
			"badSourceAddress":                   89,
			"badTimestamp":                       90, //Removed revision 22
			"busy":                               82,
			"bvlcFunctionUnknown":                143,
			"bvlcProprietaryFunctionUnknown":     144,
			"cannotUseKey":                       91, //Removed revision 22
			"cannotVerifyMessageId":              92, // Removed revision 22
			"characterSetNotSupported":           41,
			"communicationDisabled":              83,
			"configurationInProgress":            2,
			"correctKeyRevision":                 93, // Removed revision 22
			"covSubscriptionFailed":              43,
			"datatypeNotSupported":               47,
			"deleteFdtEntryFailed":               120,
			"deviceBusy":                         3,
			"destinationDeviceIdRequired":        94, // Removed revision 22
			"distributeBroadcastFailed":          121,
			"dnsError":                           192,
			"dnsNameResolutionFailed":            190,
			"dnsResolverFailure":                 191,
			"dnsUnavailable":                     189,
			"duplicateEntry":                     137,
			"duplicateMessage":                   95,
			"duplicateName":                      48,
			"duplicateObjectId":                  49,
			"dynamicCreationNotSupported":        4,
			"encryptionNotConfigured":            96,
			"encryptionRequired":                 97,
			"fileAccessDenied":                   5,
			"fileFull":                           128,
			"headerEncodingError":                145,
			"headerNotUnderstood":                146,
			"httpError":                          165,
			"httpNoUpgrade":                      153,
			"httpNotAServer":                     164,
			"httpResourceNotLocal":               154,
			"httpProxyAuthenticationFailed":      155,
			"httpResponseTimeout":                156,
			"httpResponseSyntaxError":            157,
			"httpResponseValueError":             158,
			"httpResponseMissingHeader":          159,
			"httpTemporaryUnavailable":           163,
			"httpUnexpectedResponseCode":         152,
			"httpUpgradeRequired":                161,
			"httpUpgradeError":                   162,
			"httpWebsocketHeaderError":           160,
			"inconsistentConfiguration":          129,
			"inconsistentObjectType":             130,
			"inconsistentParameters":             7,
			"inconsistentSelectionCriterion":     8,
			"incorrectKey":                       98, // Removed revision 22
			"internalError":                      131,
			"invalidArrayIndex":                  42,
			"invalidConfigurationData":           46,
			"invalidDataEncoding":                142,
			"invalidDataType":                    9,
			"invalidEventState":                  73,
			"invalidFileAccessMethod":            10,
			"invalidFileStartPosition":           11,
			"invalidKeyData":                     99, // Removed revision 22
			"invalidParameterDataType":           13,
			"invalidTag":                         57,
			"invalidTimeStamp":                   14,
			"invalidValueInThisState":            138,
			"ipAddressNotReachable":              198,
			"ipError":                            199,
			"keyUpdateInProgress":                100, // Removed revision 22
			"listElementNotFound":                81,
			"listItemNotNumbered":                140,
			"listItemNotTimestamped":             141,
			"logBufferFull":                      75,
			"loggedValuePurged":                  76,
			"malformedMessage":                   101,
			"messageTooLong":                     113,
			"messageIncomplete":                  147,
			"missingRequiredParameter":           16,
			"networkDown":                        58,
			"noAlarmConfigured":                  74,
			"noObjectsOfSpecifiedType":           17,
			"noPropertySpecified":                77,
			"noSpaceForObject":                   18,
			"noSpaceToAddListElement":            19,
			"noSpaceToWriteProperty":             20,
			"noVtSessionsAvailable":              21,
			"nodeDuplicateVmac":                  151,
			"notABacnetScHub":                    148,
			"notConfigured":                      132,
			"notConfiguredForTriggeredLogging":   78,
			"notCovProperty":                     44,
			"notKeyServer":                       102, // Removed revision 22
			"notRouterToDnet":                    110,
			"objectDeletionNotPermitted":         23,
			"objectIdentifierAlreadyExists":      24,
			"other":                              0,
			"operationalProblem":                 25,
			"optionalFunctionalityNotSupported":  45,
			"outOfMemory":                        133,
			"parameterOutOfRange":                80,
			"passwordFailure":                    26,
			"payloadExpected":                    149,
			"propertyIsNotAList":                 22,
			"propertyIsNotAnArray":               50,
			"readAccessDenied":                   27,
			"readBdtFailed":                      117,
			"readFdtFailed":                      119,
			"registerForeignDeviceFailed":        118,
			"rejectBufferOverflow":               59,
			"rejectInconsistentParameters":       60,
			"rejectInvalidParameterDataType":     61,
			"rejectInvalidTag":                   62,
			"rejectMissingRequiredParameter":     63,
			"rejectParameterOutOfRange":          64,
			"rejectTooManyArguments":             65,
			"rejectUndefinedEnumeration":         66,
			"rejectUnrecognizedService":          67,
			"rejectProprietary":                  68,
			"rejectOther":                        69,
			"routerBusy":                         111,
			"securityError":                      114,
			"securityNotConfigured":              103,
			"serviceRequestDenied":               29,
			"sourceSecurityRequired":             104,
			"success":                            84,
			"tcpClosedByLocal":                   195,
			"tcpClosedOther":                     196,
			"tcpConnectTimeout":                  193,
			"tcpConnectionRefused":               194,
			"tcpError":                           197,
			"timeout":                            30,
			"tlsClientAuthenticationFailed":      182,
			"tlsClientCertificateError":          180,
			"tlsClientCertificateExpired":        184,
			"tlsClientCertificateRevoked":        186,
			"tleError":                           188,
			"tlsServerAuthenticationFailed":      183,
			"tlsServerCertificateError":          181,
			"tlsServerCertificateExpired":        185,
			"tlsServerCertificateRevoked":        187,
			"tooManyKeys":                        105, // Removed revision 22
			"unexpectedData":                     150,
			"unknownAuthenticationType":          106,
			"unknownDevice":                      70,
			"unknownFileSize":                    122,
			"unknownKey":                         107, // Removed revision 22
			"unknownKeyRevision":                 108, // Removed revision 22
			"unknownNetworkMessage":              112,
			"unknownObject":                      31,
			"unknownProperty":                    32,
			"unknownSubscription":                79,
			"umknownRoute":                       71,
			"unknownSourceMessage":               109, // Removed revision 22
			"unknownVtClass":                     34,
			"unknownVtSession":                   35,
			"unsupportedObjectType":              36,
			"valueNotInitialized":                72,
			"valueOutOfRange":                    37,
			"valueTooLong":                       134,
			"vtSessionAlreadyClosed":             38,
			"vtSessionTerminationFailure":        39,
			"websocket-close-error":              168,
			"websocket-closed-abnormally":        173,
			"websocket-closed-by-peer":           169,
			"websocket-data-against-policy":      175,
			"websocket-data-inconsistent":        174,
			"websocket-data-not-accepted":        172,
			"websocket-endpoint-leaves":          170,
			"websocket-error":                    179,
			"websocket-extension-missing":        177,
			"websocket-frame-too-long":           176,
			"websocket-protocol-error":           171,
			"websocket-request-unavailable":      178,
			"websocket-scheme-not-supported":     166,
			"websocket-unknown-control-message":  167,
			"writeAccessDenied":                  40,
			"writeBdtFailed":                     116,
		},
	}
	var err error
	s.Enumerated, err = NewEnumerated(NoArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return s, nil
}
