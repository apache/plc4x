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

package comp

const (
	////
	// General keys

	KWAddActor   = KnownKey("addActor")
	KWDelActor   = KnownKey("delActor")
	KWActorError = KnownKey("actorError")
	KWError      = KnownKey("error")

	////
	// comm.PCI related keys

	KWCPCIUserData    = KnownKey("user_data")
	KWCPCISource      = KnownKey("source")
	KWCPCIDestination = KnownKey("destination")

	////
	// PCI related keys

	KWPCIExpectingReply  = KnownKey("expecting_reply")
	KWPCINetworkPriority = KnownKey("network_priority")

	////
	// NPDU related keys

	KWWirtnNetwork           = KnownKey("wirtnNetwork")
	KWIartnNetworkList       = KnownKey("iartnNetworkList")
	KWIcbrtnNetwork          = KnownKey("icbrtnNetwork")
	KWIcbrtnPerformanceIndex = KnownKey("icbrtnPerformanceIndex")
	KWRmtnRejectionReason    = KnownKey("rmtnRejectionReason")
	KWRmtnDNET               = KnownKey("rmtnDNET")
	KWRbtnNetworkList        = KnownKey("rbtnNetworkList")
	KWRatnNetworkList        = KnownKey("ratnNetworkList")
	KWIrtTable               = KnownKey("irtTable")
	KWIrtaTable              = KnownKey("irtaTable")
	KWEctnDNET               = KnownKey("ectnDNET")
	KWEctnTerminationTime    = KnownKey("ectnTerminationTime")
	KWDctnDNET               = KnownKey("dctnDNET")
	KWNniNet                 = KnownKey("nniNet")
	KWNniFlag                = KnownKey("nniFlag")

	////
	// BVLL related keys

	KWBvlciResultCode = KnownKey("bvlciResultCode")
	KWBvlciBDT        = KnownKey("bvlciBDT")
	KWBvlciAddress    = KnownKey("bvlciAddress")
	KWFdAddress       = KnownKey("fdAddress")
	KWFdTTL           = KnownKey("fdTTL")
	KWFdRemain        = KnownKey("fdRemain")
	KWBvlciTimeToLive = KnownKey("bvlciTimeToLive")
	KWBvlciFDT        = KnownKey("bvlciFDT")

	////
	// APDU keys

	KWConfirmedServiceChoice   = KnownKey("choice")
	KWUnconfirmedServiceChoice = KnownKey("choice")
	KWErrorClass               = KnownKey("errorClass")
	KWErrorCode                = KnownKey("errorCode")
	KWContext                  = KnownKey("context")
	KWInvokedID                = KnownKey("invokeID")

	////
	// Object keys

	KWObjectName                = KnownKey("objectName")
	KWObjectIdentifier          = KnownKey("objectIdentifier")
	KWMaximumApduLengthAccepted = KnownKey("maximumApduLengthAccepted")
	KWSegmentationSupported     = KnownKey("segmentationSupported")
	KWMaxSegmentsAccepted       = KnownKey("segmentsAccepted")
	KWNumberOfAPDURetries       = KnownKey("numberOfAPDURetries")
	KWVendorIdentifier          = KnownKey("vendorIdentifier")
	KWCls                       = KnownKey("cls")
)
