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

package object

import (
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/basetypes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type NetworkPortObject struct {
	Object
}

func NewNetworkPortObject(options ...Option) (*NetworkPortObject, error) {
	n := new(NetworkPortObject)
	objectType := "networkPort" //56,
	properties := []Property{
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),                      //111
		NewReadableProperty("reliability", V2P(NewReliability)),                      //103
		NewReadableProperty("outOfService", V2P(NewBoolean)),                         //81
		NewReadableProperty("networkType", V2P(NewNetworkType)),                      //427
		NewReadableProperty("protocolLevel", V2P(NewProtocolLevel)),                  //482
		NewOptionalProperty("referencePort", V2P(NewUnsigned)),                       //483
		NewReadableProperty("networkNumber", V2P(NewUnsigned16)),                     //425
		NewReadableProperty("networkNumberQuality", V2P(NewNetworkNumberQuality)),    //426
		NewReadableProperty("changesPending", V2P(NewBoolean)),                       //416
		NewOptionalProperty("command", V2P(NewNetworkPortCommand)),                   //417
		NewOptionalProperty("macAddress", V2P(NewOctetString)),                       //423
		NewReadableProperty("apduLength", V2P(NewUnsigned)),                          //399
		NewReadableProperty("linkSpeed", V2P(NewReal)),                               //420
		NewOptionalProperty("linkSpeeds", ArrayOfP(NewReal, 0, 0)),                   //421
		NewOptionalProperty("linkSpeedAutonegotiate", V2P(NewBoolean)),               //422
		NewOptionalProperty("networkInterfaceName", V2P(NewCharacterString)),         //424
		NewOptionalProperty("bacnetIPMode", V2P(NewIPMode)),                          //408
		NewOptionalProperty("ipAddress", V2P(NewOctetString)),                        //400
		NewOptionalProperty("bacnetIPUDPPort", V2P(NewUnsigned16)),                   //412
		NewOptionalProperty("ipSubnetMask", V2P(NewOctetString)),                     //411
		NewOptionalProperty("ipDefaultGateway", V2P(NewOctetString)),                 //401
		NewOptionalProperty("bacnetIPMulticastAddress", V2P(NewOctetString)),         //409
		NewOptionalProperty("ipDNSServer", ArrayOfP(NewOctetString, 0, 0)),           //406
		NewOptionalProperty("ipDHCPEnable", V2P(NewBoolean)),                         //402
		NewOptionalProperty("ipDHCPLeaseTime", V2P(NewUnsigned)),                     //403
		NewOptionalProperty("ipDHCPLeaseTimeRemaining", V2P(NewUnsigned)),            //404
		NewOptionalProperty("ipDHCPServer", V2P(NewOctetString)),                     //405
		NewOptionalProperty("bacnetIPNATTraversal", V2P(NewBoolean)),                 //410
		NewOptionalProperty("bacnetIPGlobalAddress", V2P(NewHostNPort)),              //407
		NewOptionalProperty("bbmdBroadcastDistributionTable", ListOfP(NewBDTEntry)),  //414
		NewOptionalProperty("bbmdAcceptFDRegistrations", V2P(NewBoolean)),            //413
		NewOptionalProperty("bbmdForeignDeviceTable", ListOfP(NewFDTEntry)),          //415
		NewOptionalProperty("fdBBMDAddress", V2P(NewHostNPort)),                      //418
		NewOptionalProperty("fdSubscriptionLifetime", V2P(NewUnsigned16)),            //419
		NewOptionalProperty("bacnetIPv6Mode", V2P(NewIPMode)),                        //435
		NewOptionalProperty("ipv6Address", V2P(NewOctetString)),                      //436
		NewOptionalProperty("ipv6PrefixLength", V2P(NewUnsigned8)),                   //437
		NewOptionalProperty("bacnetIPv6UDPPort", V2P(NewUnsigned16)),                 //438
		NewOptionalProperty("ipv6DefaultGateway", V2P(NewOctetString)),               //439
		NewOptionalProperty("bacnetIPv6MulticastAddress", V2P(NewOctetString)),       //440
		NewOptionalProperty("ipv6DNSServer", V2P(NewOctetString)),                    //441
		NewOptionalProperty("ipv6AutoAddressingEnable", V2P(NewBoolean)),             //442
		NewOptionalProperty("ipv6DHCPLeaseTime", V2P(NewUnsigned)),                   //443
		NewOptionalProperty("ipv6DHCPLeaseTimeRemaining", V2P(NewUnsigned)),          //444
		NewOptionalProperty("ipv6DHCPServer", V2P(NewOctetString)),                   //445
		NewOptionalProperty("ipv6ZoneIndex", V2P(NewCharacterString)),                //446
		NewOptionalProperty("maxMaster", V2P(NewUnsigned8)),                          //64
		NewOptionalProperty("maxInfoFrames", V2P(NewUnsigned8)),                      //63
		NewOptionalProperty("slaveProxyEnable", V2P(NewBoolean)),                     //172
		NewOptionalProperty("manualSlaveAddressBinding", ListOfP(NewAddressBinding)), //170
		NewOptionalProperty("autoSlaveDiscovery", V2P(NewBoolean)),                   //169
		NewOptionalProperty("slaveAddressBinding", ListOfP(NewAddressBinding)),       //171
		NewOptionalProperty("virtualMACAddressTable", ListOfP(NewVMACEntry)),         //429
		NewOptionalProperty("routingTable", ListOfP(NewRouterEntry)),                 //428
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),                 //353
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),                   //17
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),              //35
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),         //0
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),                        //72
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventState", V2P(NewEventState)),                //36,
		NewReadableProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)), //357,
	}
	var err error
	n.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, n)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return n, nil
}
