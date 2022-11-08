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
package org.apache.plc4x.java.bacnetip;

import org.apache.commons.io.FileUtils;
import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.utils.ascii.AsciiBox;
import org.apache.plc4x.java.spi.utils.ascii.AsciiBoxWriter;
import org.apache.plc4x.java.spi.utils.hex.Hex;
import org.apache.plc4x.test.RequirePcapNg;
import org.apache.plc4x.test.hex.HexDiff;
import org.assertj.core.api.InstanceOfAssertFactories;
import org.assertj.core.presentation.HexadecimalRepresentation;
import org.junit.jupiter.api.*;
import org.opentest4j.TestAbortedException;
import org.pcap4j.core.*;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.UdpPacket;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;
import java.math.BigInteger;
import java.net.URL;
import java.nio.file.*;
import java.sql.Timestamp;
import java.util.*;
import java.util.concurrent.ConcurrentLinkedDeque;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

// Tests from http://kargs.net/captures
@RequirePcapNg
//@RequireAllTestsFlag
@Tag("require-all-tests")
@Tag("bacnet-regression")
public class RandomPackagesTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(RandomPackagesTest.class);

    public static final String BACNET_BPF_FILTER_UDP = "udp port 47808";

    Queue<Closeable> toBeClosed = new ConcurrentLinkedDeque<>();

    @AfterEach
    void closeStuff() {
        for (Closeable closeable = toBeClosed.poll(); closeable != null; closeable = toBeClosed.poll()) {
            LOGGER.info("Closing closeable " + closeable);
            IOUtils.closeQuietly(closeable);
        }
    }

    @TestFactory
    @DisplayName("BACnet-BBMD-on-same-subnet")
    Collection<DynamicNode> BACnet_BBMD_on_same_subnet() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnet-BBMD-on-same-subnet.cap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - BACnet Virtual Link Control BVLC Function Register-Foreign-Device",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    assertEquals(60000, ((BVLCRegisterForeignDevice) bvlc).getTtl());
                }),
            DynamicTest.dynamicTest("No. 2 - Unconfirmed whoIs",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCDistributeBroadcastToNetwork) bvlc).getNpdu().getApdu();
                    assertEquals(BACnetUnconfirmedServiceChoice.WHO_IS, apduUnconfirmedRequest.getServiceRequest().getServiceChoice());
                }),
            DynamicTest.dynamicTest("No. 3 - BACnet Virtual Link Control BVLC Function BVLC-Results",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    assertEquals(BVLCResultCode.SUCCESSFUL_COMPLETION, ((BVLCResult) bvlc).getCode());
                }),
            DynamicTest.dynamicTest("No. 4 - Unconfirmed-REQ who-Is",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCForwardedNPDU) bvlc).getNpdu().getApdu();
                    assertEquals(BACnetUnconfirmedServiceChoice.WHO_IS, apduUnconfirmedRequest.getServiceRequest().getServiceChoice());
                }),
            DynamicTest.dynamicTest("No. 5 - Unconfirmed-REQ i-Am device,123",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(123, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(1476, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.BA_CNET_STACKAT_SOURCE_FORGE, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 6 - Unconfirmed-REQ i-Am device,123",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCForwardedNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(123, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(1476, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.BA_CNET_STACKAT_SOURCE_FORGE, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 7 - Unconfirmed-REQ i-Am device,18",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(18, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(1476, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.ALERTON_HONEYWELL, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 8 - Unconfirmed-REQ i-Am device,18",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCForwardedNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(18, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(1476, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.ALERTON_HONEYWELL, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 9 - Unconfirmed-REQ i-Am device,2401",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(2401, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.AUTOMATED_LOGIC_CORPORATION, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 10 - Unconfirmed-REQ i-Am device,2401",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCForwardedNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(2401, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.AUTOMATED_LOGIC_CORPORATION, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 11 - Unconfirmed-REQ i-Am device,86114",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(86114, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(50, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.BA_CNET_STACKAT_SOURCE_FORGE, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 12 - Unconfirmed-REQ i-Am device,86114",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCForwardedNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(86114, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(50, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.BA_CNET_STACKAT_SOURCE_FORGE, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 13 - Unconfirmed-REQ i-Am device,884456",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(884456, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.THE_WATT_STOPPER_INC, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 14 - Unconfirmed-REQ i-Am device,884456",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) ((BVLCForwardedNPDU) bvlc).getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(884456, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.THE_WATT_STOPPER_INC, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                })
        );
    }

    @TestFactory
    @DisplayName("BACnet-MSTP-SNAP-Mixed")
    Collection<DynamicNode> BACnet_MSTP_SNAP_Mixed() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnet-MSTP-SNAP-Mixed.cap", BACNET_BPF_FILTER_UDP);
        return List.of(
            pcapEvaluator.parseEmAll(
            )
        );
    }

    @TestFactory
    @DisplayName("BACnetARRAY-element-0")
    Collection<DynamicNode> BACnetARRAY_element_0() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnetARRAY-element-0.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("BACnetARRAY-elements")
    Collection<DynamicNode> BACnetARRAY_elements() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnetARRAY-elements.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("BACnetDeviceObjectReference")
    Collection<DynamicNode> BACnetDeviceObjectReference() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnetDeviceObjectReference.pcap", BACNET_BPF_FILTER_UDP);
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   readProperty[  0] life-safety-zone,1 zone-members",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.LIFE_SAFETY_ZONE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.ZONE_MEMBERS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[  0] life-safety-zone,1 zone-members life-safety-zone,3 life-safety-zone,4 life-safety-zone,5 life-safety-zone,6 life-safety-zone,7 life-safety-zone,8 life-safety-zone,9 life-safety-zone,16 life-safety-zone,494 life-safety-zone,255 life-safety-zone,231 life-safety-zone,4193620 life-safety-zone,222 life-safety-zone,300 life-safety-zone,166",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.LIFE_SAFETY_ZONE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.ZONE_MEMBERS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetConstructedDataZoneMembers baCnetConstructedDataZoneMembers = (BACnetConstructedDataZoneMembers) baCnetServiceAckReadProperty.getValues();

                    List<BACnetDeviceObjectReference> members = baCnetConstructedDataZoneMembers.getMembers();
                    assertThat(members.get(0)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 3L);
                    assertThat(members.get(1)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 4L);
                    assertThat(members.get(2)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 5L);
                    assertThat(members.get(3)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 6L);
                    assertThat(members.get(4)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 7L);
                    assertThat(members.get(5)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 8L);
                    assertThat(members.get(6)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 9L);
                    assertThat(members.get(7)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 16L);
                    assertThat(members.get(8)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 494L);
                    assertThat(members.get(9)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 255L);
                    assertThat(members.get(10)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 231L);
                    assertThat(members.get(11)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 4193620L);
                    assertThat(members.get(12)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 222L);
                    assertThat(members.get(13)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 300L);
                    assertThat(members.get(14)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 166L);
                }),
            DynamicTest.dynamicTest("No. 3 - Confirmed-REQ   readProperty[  0] life-safety-zone,1 member-of",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.LIFE_SAFETY_ZONE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MEMBER_OF, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 4 - Complex-ACK     readProperty[  0] life-safety-zone,1 member-of life-safety-zone,1",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.LIFE_SAFETY_ZONE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MEMBER_OF, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetConstructedDataMemberOf baCnetConstructedDataMemberOf = (BACnetConstructedDataMemberOf) baCnetServiceAckReadProperty.getValues();

                    List<BACnetDeviceObjectReference> zones = baCnetConstructedDataMemberOf.getZones();
                    assertThat(zones.get(0)).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 1L);
                })
        );
    }

    @TestFactory
    @DisplayName("BACnetIP-MSTP-Mix")
    Collection<DynamicNode> BACnet_MSTP_Mix() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnetIP-MSTP-Mix.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("BBMD_Results")
    Collection<DynamicNode> BBMD_Results() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BBMD_Results.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("BBMD_readproperty")
    Collection<DynamicNode> BBMD_readProperty() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BBMD_readproperty.pcap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Unconfirmed-REQ who-Is 12345 12345",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCDistributeBroadcastToNetwork bvlcDistributeBroadcastToNetwork = (BVLCDistributeBroadcastToNetwork) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcDistributeBroadcastToNetwork.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) serviceRequest;
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 2 - Unconfirmed-REQ who-Is 12345 12345",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) serviceRequest;
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 3 - Unconfirmed-REQ i-Am device,12345",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(12345, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.BA_CNET_STACKAT_SOURCE_FORGE, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 4 - Unconfirmed-REQ who-Is 12345 12345",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCForwardedNPDU bvlcForwardedNPDU = (BVLCForwardedNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcForwardedNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) serviceRequest;
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 5 - Unconfirmed-REQ who-Is 12345 12345",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCForwardedNPDU bvlcForwardedNPDU = (BVLCForwardedNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcForwardedNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) serviceRequest;
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(12345L, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 6 - Unconfirmed-REQ i-Am device,12345",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCForwardedNPDU bvlcForwardedNPDU = (BVLCForwardedNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcForwardedNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(12345, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.BA_CNET_STACKAT_SOURCE_FORGE, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 7 - Confirmed-REQ readProperty[ 1] analog-output,0 priority-array",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PRIORITY_ARRAY, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 8 - Complex-ACK readProperty[ 1] analog-output,0 priority-array",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PRIORITY_ARRAY, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetConstructedDataPriorityArray priorityArray = (BACnetConstructedDataPriorityArray) baCnetServiceAckReadProperty.getValues();
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue01().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue02().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue03().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue04().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue05().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue06().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue07().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue08().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue09().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue10().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue11().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue12().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue13().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue14().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue15().getClass());
                    assertEquals(BACnetPriorityValueNull.class, priorityArray.getPriorityArray().getPriorityValue16().getClass());
                }),
            DynamicTest.dynamicTest("No. 9 - BACnet Virtual Link Control BVLC Function Register-Foreign-Device",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCRegisterForeignDevice bvlcRegisterForeignDevice = (BVLCRegisterForeignDevice) bvlc;
                    assertEquals(60000, bvlcRegisterForeignDevice.getTtl());
                }),
            DynamicTest.dynamicTest("No. 10 - Unconfirmed-REQ who-Is 12345 12345",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipPackages(1);
                }),
            DynamicTest.dynamicTest("No. 11 - BACnet Virtual Link Control BVLC Function BVLC-Result",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCResult bvlcResult = (BVLCResult) bvlc;
                    assertEquals(BVLCResultCode.SUCCESSFUL_COMPLETION, bvlcResult.getCode());
                }),
            DynamicTest.dynamicTest("No. 12-16 - Skip Unconfirmed-REQ who-Is/I-Am",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipPackages(5);
                }),
            DynamicTest.dynamicTest("No. 17 - Confirmed-REQ readProperty[ 1] analog-output,0 present-value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 18 - Complex-ACK readProperty[ 1] analog-output,0 present-value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagReal baCnetApplicationTagReal = ((BACnetConstructedDataAnalogOutputPresentValue) baCnetServiceAckReadProperty.getValues()).getPresentValue();
                    assertEquals(0, baCnetApplicationTagReal.getPayload().getValue());
                }),
            DynamicTest.dynamicTest("No. 19-26 - Skip Misc 8 packages",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipPackages(8);
                }),
            DynamicTest.dynamicTest("No. 27 - Confirmed-REQ readProperty[ 1] analog-output,0 relinquish-default",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.RELINQUISH_DEFAULT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 28 - Complex-ACK readProperty[ 1] analog-output,0 relinquish-default",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.RELINQUISH_DEFAULT, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagReal baCnetApplicationTagReal = ((BACnetConstructedDataAnalogOutputRelinquishDefault) baCnetServiceAckReadProperty.getValues()).getRelinquishDefault();
                    assertEquals(0.0f, baCnetApplicationTagReal.getActualValue());
                }),
            DynamicTest.dynamicTest("No. 29-76 - Skip Misc 48 packages",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipPackages(48);
                }),
            DynamicTest.dynamicTest("No. 77 - Confirmed-REQ writeProperty[ 1] analog-output,0 priority-array",
                () -> {
                    // This package is broken as from the spec it requires 16 values
                    pcapEvaluator.skipPackages(1);
                }),
            DynamicTest.dynamicTest("No. 78 - Error writeProperty[ 1]",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorGeneral baCnetErrorWriteProperty = (BACnetErrorGeneral) apduError.getError();
                    assertEquals(ErrorClass.PROPERTY, baCnetErrorWriteProperty.getError().getErrorClass().getValue());
                    assertEquals(ErrorCode.WRITE_ACCESS_DENIED, baCnetErrorWriteProperty.getError().getErrorCode().getValue());
                }),
            DynamicTest.dynamicTest("No. 79-86 - Skip Misc 8 packages",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipPackages(8);
                }),
            DynamicTest.dynamicTest("No. 87 - Confirmed-REQ writeProperty[ 1] analog-output,0 present-value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestWriteProperty baCnetConfirmedServiceRequestWriteProperty = (BACnetConfirmedServiceRequestWriteProperty) serviceRequest;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetConfirmedServiceRequestWriteProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetConfirmedServiceRequestWriteProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetConfirmedServiceRequestWriteProperty.getPropertyIdentifier().getValue());

                    BACnetApplicationTagReal baCnetApplicationTagReal = ((BACnetConstructedDataAnalogOutputPresentValue) baCnetConfirmedServiceRequestWriteProperty.getPropertyValue()).getPresentValue();
                    assertEquals(123.449997f, baCnetApplicationTagReal.getPayload().getValue());
                    BACnetContextTagUnsignedInteger priority = baCnetConfirmedServiceRequestWriteProperty.getPriority();
                    assertEquals(10, priority.getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 88 - Error writeProperty[ 1]",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorGeneral baCnetErrorWriteProperty = (BACnetErrorGeneral) apduError.getError();
                    assertEquals(ErrorClass.PROPERTY, baCnetErrorWriteProperty.getError().getErrorClass().getValue());
                    assertEquals(ErrorCode.VALUE_OUT_OF_RANGE, baCnetErrorWriteProperty.getError().getErrorCode().getValue());
                }),
            DynamicTest.dynamicTest("No. 89-142 - Skip to 142 Misc packages",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipTo(142);
                }),
            DynamicTest.dynamicTest("No. 143 - Confirmed-REQ writeProperty[ 1] analog-output,0 present-value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestWriteProperty baCnetConfirmedServiceRequestWriteProperty = (BACnetConfirmedServiceRequestWriteProperty) serviceRequest;
                    assertEquals(BACnetObjectType.ANALOG_OUTPUT, baCnetConfirmedServiceRequestWriteProperty.getObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetConfirmedServiceRequestWriteProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetConfirmedServiceRequestWriteProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagNull baCnetApplicationTagNull = (BACnetApplicationTagNull) ((BACnetConstructedDataUnspecified) baCnetConfirmedServiceRequestWriteProperty.getPropertyValue()).getData().get(0).getApplicationTag();
                    assertNotNull(baCnetApplicationTagNull);
                    BACnetContextTagUnsignedInteger priority = baCnetConfirmedServiceRequestWriteProperty.getPriority();
                    assertEquals(1, priority.getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 144 - Simple-ACK writeProperty[ 1]", () -> {
                BVLC bvlc = pcapEvaluator.nextBVLC();
                dump(bvlc);
                NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                APDUSimpleAck apduSimpleAck = (APDUSimpleAck) npdu.getApdu();
                assertEquals(15, apduSimpleAck.getServiceChoice());
            }),
            DynamicTest.dynamicTest("No. 145-200 - Skip to 200 Misc packages",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipTo(200);
                }),
            DynamicTest.dynamicTest("No. 201 - Confirmed-REQ readProperty[  1] device,12345 object-identifier",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(12345, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 202 - Complex-ACK   readProperty[  1] device,12345 object-identifier device,12345",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(12345, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagObjectIdentifier objectIdentifier = ((BACnetConstructedDataObjectIdentifier) baCnetServiceAckReadProperty.getValues()).getObjectIdentifier();
                    assertEquals(BACnetObjectType.DEVICE, objectIdentifier.getObjectType());
                    assertEquals(12345, objectIdentifier.getInstanceNumber());
                }),
            pcapEvaluator.parseFrom(203
            )
        );
    }

    @TestFactory
    @DisplayName("CEN_9_11")
    Collection<DynamicNode> CEN_9_11() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("CEN_9_11.pcap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   confirmedEventNotification[119] event-enrollment,11 analog-input,1",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestConfirmedEventNotification baCnetConfirmedServiceRequestConfirmedEventNotification = (BACnetConfirmedServiceRequestConfirmedEventNotification) serviceRequest;
                    assertEquals((short) 111, baCnetConfirmedServiceRequestConfirmedEventNotification.getProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.EVENT_ENROLLMENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getObjectType());
                    assertEquals(11, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.ANALOG_INPUT, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getInstanceNumber());
                    {
                        BACnetTimeStampSequence timestamp = (BACnetTimeStampSequence) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp().getTimestamp();
                        assertEquals(2, timestamp.getSequenceNumber().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.UNSIGNED_RANGE, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getValue());
                    }
                    {
                        assertEquals("My Message", baCnetConfirmedServiceRequestConfirmedEventNotification.getMessageText().getPayload().getValue());
                    }
                    {
                        assertEquals(BACnetNotifyType.EVENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsFalse());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getValue());
                    }
                    {
                        assertEquals(BACnetEventState.OFFNORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getValue());
                    }
                    {
                        BACnetNotificationParametersUnsignedRange baCnetNotificationParametersUnsignedRange = (BACnetNotificationParametersUnsignedRange) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertEquals(50, baCnetNotificationParametersUnsignedRange.getSequenceNumber().getPayload().getActualValue().longValue());
                        assertTrue(baCnetNotificationParametersUnsignedRange.getStatusFlags().getInAlarm());
                        assertFalse(baCnetNotificationParametersUnsignedRange.getStatusFlags().getFault());
                        assertFalse(baCnetNotificationParametersUnsignedRange.getStatusFlags().getOverridden());
                        assertFalse(baCnetNotificationParametersUnsignedRange.getStatusFlags().getOutOfService());
                        assertEquals(40, baCnetNotificationParametersUnsignedRange.getExceededLimit().getPayload().getActualValue().longValue());
                    }
                }),
            DynamicTest.dynamicTest("No. 2 - Simple-ACK      confirmedEventNotification[119]",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUSimpleAck apduSimpleAck = (APDUSimpleAck) npdu.getApdu();
                    assertEquals(119, apduSimpleAck.getOriginalInvokeId());
                    assertEquals(2, apduSimpleAck.getServiceChoice());
                }),
            DynamicTest.dynamicTest("No. 3 - Confirmed-REQ   confirmedEventNotification[120] event-enrollment,11 analog-input,1",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestConfirmedEventNotification baCnetConfirmedServiceRequestConfirmedEventNotification = (BACnetConfirmedServiceRequestConfirmedEventNotification) serviceRequest;
                    assertEquals((short) 111, baCnetConfirmedServiceRequestConfirmedEventNotification.getProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.EVENT_ENROLLMENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getObjectType());
                    assertEquals(11, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.ANALOG_INPUT, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getInstanceNumber());
                    {
                        BACnetTimeStampSequence timestamp = (BACnetTimeStampSequence) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp().getTimestamp();
                        assertEquals(2, timestamp.getSequenceNumber().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.EXTENDED, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getValue());
                    }
                    {
                        assertEquals("My Message", baCnetConfirmedServiceRequestConfirmedEventNotification.getMessageText().getPayload().getValue());
                    }
                    {
                        assertEquals(BACnetNotifyType.EVENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsFalse());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getValue());
                    }
                    {
                        assertEquals(BACnetEventState.OFFNORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getValue());
                    }
                    {
                        BACnetNotificationParametersExtended baCnetNotificationParametersExtended = (BACnetNotificationParametersExtended) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertEquals(BACnetVendorId.AUTOMATED_LOGIC_CORPORATION, baCnetNotificationParametersExtended.getVendorId().getValue());
                        assertEquals(33, baCnetNotificationParametersExtended.getExtendedEventType().getPayload().getActualValue().longValue());
                    }
                }),
            DynamicTest.dynamicTest("No. 4 - Simple-ACK      confirmedEventNotification[120]",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUSimpleAck apduSimpleAck = (APDUSimpleAck) npdu.getApdu();
                    assertEquals(120, apduSimpleAck.getOriginalInvokeId());
                    assertEquals(2, apduSimpleAck.getServiceChoice());
                }),
            DynamicTest.dynamicTest("No. 5 - Unconfirmed-REQ who-Is 140 140",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(140, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(140, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 6 - Unconfirmed-REQ who-Is 140 140",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(140, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(140, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 7 - Unconfirmed-REQ who-Is 871 871",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(871, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(871, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 8 - Unconfirmed-REQ who-Is 871 871",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(871, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(871, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getPayload().getActualValue().longValue());
                })
        );
    }

    @TestFactory
    @DisplayName("CEN_10")
    Collection<DynamicNode> CEN_10() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("CEN_10.pcap");
        return List.of(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   confirmedEventNotification[  7] device,151 trend-log,1 trend-log,1 log-buffer device,151",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestConfirmedEventNotification baCnetConfirmedServiceRequestConfirmedEventNotification = (BACnetConfirmedServiceRequestConfirmedEventNotification) serviceRequest;
                    assertEquals((short) 0, baCnetConfirmedServiceRequestConfirmedEventNotification.getProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getObjectType());
                    assertEquals(151, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.TREND_LOG, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getInstanceNumber());
                    {
                        BACnetTimeStampDateTime timestamp = (BACnetTimeStampDateTime) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp().getTimestamp();
                        BACnetTagPayloadDate payload = timestamp.getDateTimeValue().getDateTimeValue().getDateValue().getPayload();
                        assertEquals(2008, payload.getYear());
                        assertEquals(5, payload.getMonth());
                        assertEquals(2, payload.getDayOfMonth());
                        assertEquals(5, payload.getDayOfWeek());
                        BACnetTagPayloadTime baCnetTagPayloadTime = timestamp.getDateTimeValue().getDateTimeValue().getTimeValue().getPayload();
                        assertEquals(11, baCnetTagPayloadTime.getHour());
                        assertEquals(11, baCnetTagPayloadTime.getMinute());
                        assertEquals(30, baCnetTagPayloadTime.getSecond());
                        assertEquals(0, baCnetTagPayloadTime.getFractional());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(15, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.BUFFER_READY, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getValue());
                    }
                    {
                        assertEquals(BACnetNotifyType.EVENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsTrue());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getValue());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getValue());
                    }
                    {
                        BACnetNotificationParametersBufferReady baCnetNotificationParametersBufferReady = (BACnetNotificationParametersBufferReady) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertEquals(BACnetObjectType.TREND_LOG, baCnetNotificationParametersBufferReady.getBufferProperty().getValue().getObjectIdentifier().getObjectType());
                        assertEquals(BACnetPropertyIdentifier.LOG_BUFFER, baCnetNotificationParametersBufferReady.getBufferProperty().getValue().getPropertyIdentifier().getValue());
                        assertEquals(BACnetObjectType.DEVICE, baCnetNotificationParametersBufferReady.getBufferProperty().getValue().getDeviceIdentifier().getObjectType());
                        assertEquals(1640, baCnetNotificationParametersBufferReady.getPreviousNotification().getPayload().getActualValue().longValue());
                        assertEquals(1653, baCnetNotificationParametersBufferReady.getCurrentNotification().getPayload().getActualValue().longValue());
                    }
                })
        );
    }

    @TestFactory
    @DisplayName("COV_AWF_ARF")
    Collection<DynamicNode> COV_AWF_ARF() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("COV_AWF_ARF.cap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   subscribeCOV[ 10] binary-input,0",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestSubscribeCOV baCnetConfirmedServiceRequestSubscribeCOV = (BACnetConfirmedServiceRequestSubscribeCOV) serviceRequest;
                    assertEquals((short) 123, baCnetConfirmedServiceRequestSubscribeCOV.getSubscriberProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.BINARY_INPUT, baCnetConfirmedServiceRequestSubscribeCOV.getMonitoredObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetConfirmedServiceRequestSubscribeCOV.getMonitoredObjectIdentifier().getInstanceNumber());
                    assertTrue(baCnetConfirmedServiceRequestSubscribeCOV.getIssueConfirmed().getPayload().getIsFalse());
                    assertEquals(10, baCnetConfirmedServiceRequestSubscribeCOV.getLifetimeInSeconds().getPayload().getActualValue().longValue() / 60);
                }),
            DynamicTest.dynamicTest("No. 2 - Simple-ACK      subscribeCOV[ 10]",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUSimpleAck apduSimpleAck = (APDUSimpleAck) npdu.getApdu();
                    assertEquals((short) 5, apduSimpleAck.getServiceChoice());
                }),
            DynamicTest.dynamicTest("No. 3 - Unconfirmed-REQ unconfirmedCOVNotification device,12345 binary-input,0 present-value status-flags",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification = (BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) serviceRequest;
                    assertEquals((short) 123, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getSubscriberProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getInitiatingDeviceIdentifier().getObjectType());
                    assertEquals(12345, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.BINARY_INPUT, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getMonitoredObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getMonitoredObjectIdentifier().getInstanceNumber());
                    assertEquals(9, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getLifetimeInSeconds().getPayload().getActualValue().longValue() / 60);
                    {
                        assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(0).getPropertyIdentifier().getValue());
                    }
                    {
                        BACnetConstructedDataBinaryInputPresentValue constructedData = (BACnetConstructedDataBinaryInputPresentValue) baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(0).getPropertyValue().getConstructedData();
                        BACnetBinaryPV baCnetBinaryPV = constructedData.getPresentValue().getValue();
                        assertEquals(BACnetBinaryPV.INACTIVE, baCnetBinaryPV);
                    }
                    {
                        assertEquals(BACnetPropertyIdentifier.STATUS_FLAGS, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(1).getPropertyIdentifier().getValue());
                    }
                    {
                        assertThat((BACnetConstructedDataStatusFlags) baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(1).getPropertyValue().getConstructedData())
                            .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                            .satisfies(statusFlagsTagged ->
                                assertThat(statusFlagsTagged)
                                    .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                    .isEqualTo(List.of(false, false, false, false)));
                    }
                }),
            pcapEvaluator.parseTill(1347)
        );
    }

    @Disabled("BACnetAction is returned with 4 and that explodes so not sure how to handle that")
    @TestFactory
    @DisplayName("ContextTagAbove14Sample_1")
    Collection<DynamicNode> ContextTagAbove14Sample_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("ContextTagAbove14Sample_1.pcap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - skipLLC",
                () -> pcapEvaluator.skipPackages(1)),
            DynamicTest.dynamicTest("No. 2 - Confirmed-REQ   confirmedEventNotification[138] device,1 event-enrollment,1",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC(2);
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestConfirmedEventNotification baCnetConfirmedServiceRequestConfirmedEventNotification = (BACnetConfirmedServiceRequestConfirmedEventNotification) serviceRequest;
                    assertEquals((short) 1, baCnetConfirmedServiceRequestConfirmedEventNotification.getProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getObjectType());
                    assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.EVENT_ENROLLMENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getObjectType());
                    assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getInstanceNumber());
                    {
                        BACnetTimeStampSequence timestamp = (BACnetTimeStampSequence) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp().getTimestamp();
                        assertEquals(1, timestamp.getSequenceNumber().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(111, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.CHANGE_OF_STATE, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getValue());
                    }
                    {
                        assertEquals(BACnetNotifyType.EVENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsFalse());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getValue());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getValue());
                    }
                    {
                        BACnetNotificationParametersChangeOfState baCnetNotificationParametersChangeOfState = (BACnetNotificationParametersChangeOfState) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertTrue(baCnetNotificationParametersChangeOfState.getStatusFlags().getInAlarm());
                        assertFalse(baCnetNotificationParametersChangeOfState.getStatusFlags().getFault());
                        assertFalse(baCnetNotificationParametersChangeOfState.getStatusFlags().getOverridden());
                        assertFalse(baCnetNotificationParametersChangeOfState.getStatusFlags().getOutOfService());
                    }
                })
        );
    }

    @TestFactory
    @DisplayName("CriticalRoom55-1")
    Collection<DynamicNode> CriticalRoom55_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("CriticalRoom55-1.cap");
        return List.of(pcapEvaluator.parseEmAll(
            skip(8, "undefined reject... some glibberich"),
            skip(10, "undefined reject... some glibberich"),
            skip(12, "undefined reject... some glibberich"),
            skip(14, "undefined reject... some glibberich"),
            skip(16, "undefined reject... some glibberich"),
            skip(18, "undefined reject... some glibberich"),
            skip(22, "undefined reject... some glibberich"),
            skip(24, "undefined reject... some glibberich"),
            skip(26, "undefined reject... some glibberich")
        ));
    }

    @TestFactory
    @DisplayName("CriticalRoom55-2")
    Collection<DynamicNode> CriticalRoom55_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("CriticalRoom55-2.cap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ writeProperty[113] analog-value,1 present-value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestWriteProperty baCnetConfirmedServiceRequestWriteProperty = (BACnetConfirmedServiceRequestWriteProperty) serviceRequest;
                    BACnetApplicationTagReal baCnetApplicationTagReal = ((BACnetConstructedDataAnalogValuePresentValue) baCnetConfirmedServiceRequestWriteProperty.getPropertyValue()).getPresentValue();
                    assertEquals(123.0f, baCnetApplicationTagReal.getPayload().getValue());
                }),
            pcapEvaluator.parseFrom(2,
                IntStream.rangeClosed(2, 281)
                    .filter(i -> i % 2 == 0)
                    .mapToObj(i -> skip(i, SkipInstruction.SkipType.SKIP_COMPLETE, "The responses here are just garbage"))
                    .toArray(SkipInstruction[]::new)
            )
        );
    }

    @TestFactory
    @DisplayName("DRI%20CAVE%20log%20udp-0168-20081216-1117-03")
    Collection<DynamicNode> DRI_CAVE_log_udp_0168_20081216_1117_03() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("DRI%20CAVE%20log%20udp-0168-20081216-1117-03.cap");
        return List.of(
            pcapEvaluator.parseEmAll(
                skip(55, SkipInstruction.SkipType.SKIP_COMPARE, "This is a unknown service so we can't serialize because of dicriminator"),
                skip(60, SkipInstruction.SkipType.SKIP_COMPARE, "This is a unknown service so we can't serialize because of dicriminator"),
                skip(92, SkipInstruction.SkipType.SKIP_COMPARE, "This is a unknown service so we can't serialize because of dicriminator"),
                skip(99, SkipInstruction.SkipType.SKIP_COMPARE, "This is a unknown service so we can't serialize because of dicriminator"),
                skip(131, SkipInstruction.SkipType.SKIP_COMPARE, "This is a unknown service so we can't serialize because of dicriminator"),
                skip(134, SkipInstruction.SkipType.SKIP_COMPARE, "This is a unknown service so we can't serialize because of dicriminator"),
                skip(86, "incomplete captured package (size limit 100)"),
                skip(87, "broken package"),
                skip(94, "incomplete captured package (size limit 100)"),
                skip(95, "broken package"),
                skip(104, "incomplete captured package (size limit 100)"),
                skip(105, "broken package"),
                skip(110, "incomplete captured package (size limit 100)"),
                skip(111, "broken package"),
                skip(140, "incomplete captured package (size limit 100)"),
                skip(142, "broken package"),
                skip(143, "broken package"),
                skip(147, "broken package"),
                skip(151, "incomplete captured package (size limit 100)"),
                skip(152, "broken package")
            )
        );
    }

    @TestFactory
    @DisplayName("I-Am-Router-To-Network")
    Collection<DynamicNode> I_Am_Router_To_Network() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("I-Am-Router-To-Network.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("Ethereal-Misinterpreted-Packet")
    Collection<DynamicNode> Ethereal_Misinterpreted_Packet() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Ethereal-Misinterpreted-Packet.cap");
        return List.of(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   confirmedEventNotification[ 10] device,1041000 analog-input,3000016 (2200) Vendor Proprietary Value object-name (2201) Vendor Proprietary Value (2202) Vendor Proprietary Value reliability (661) VendorProprietary Value units (1659) Vendor Proprietary Value (2203) Vendor Proprietary Value vendor-identifier",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestConfirmedEventNotification baCnetConfirmedServiceRequestConfirmedEventNotification = (BACnetConfirmedServiceRequestConfirmedEventNotification) serviceRequest;
                    assertEquals((short) 0, baCnetConfirmedServiceRequestConfirmedEventNotification.getProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getObjectType());
                    assertEquals(1041000, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.ANALOG_INPUT, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getObjectType());
                    assertEquals(3000016, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getInstanceNumber());
                    {
                        BACnetTimeStampDateTime timestamp = (BACnetTimeStampDateTime) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp().getTimestamp();
                        assertEquals(2005, timestamp.getDateTimeValue().getDateTimeValue().getDateValue().getPayload().getYear());
                        assertEquals(12, timestamp.getDateTimeValue().getDateTimeValue().getDateValue().getPayload().getMonth());
                        assertEquals(8, timestamp.getDateTimeValue().getDateTimeValue().getDateValue().getPayload().getDayOfMonth());
                        assertEquals(4, timestamp.getDateTimeValue().getDateTimeValue().getDateValue().getPayload().getDayOfWeek());
                        assertEquals(14, timestamp.getDateTimeValue().getDateTimeValue().getTimeValue().getPayload().getHour());
                        assertEquals(12, timestamp.getDateTimeValue().getDateTimeValue().getTimeValue().getPayload().getMinute());
                        assertEquals(49, timestamp.getDateTimeValue().getDateTimeValue().getTimeValue().getPayload().getSecond());
                        assertEquals(0, timestamp.getDateTimeValue().getDateTimeValue().getTimeValue().getPayload().getFractional());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(200, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.OUT_OF_RANGE, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getValue());
                    }
                    {
                        assertEquals(BACnetNotifyType.ALARM, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsTrue());
                    }
                    {
                        assertEquals(BACnetEventState.HIGH_LIMIT, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getValue());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getValue());
                    }
                    {
                        BACnetNotificationParametersComplexEventType baCnetNotificationParametersComplexEventType = (BACnetNotificationParametersComplexEventType) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(0);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            assertEquals(2200, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagCharacterString baCnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals("StockingNAE", baCnetApplicationTagCharacterString.getValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(1);
                            assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            BACnetApplicationTagCharacterString baCnetApplicationTagCharacterString = ((BACnetConstructedDataObjectName) baCnetPropertyValue.getPropertyValue().getConstructedData()).getObjectName();
                            assertEquals("StockingNAE/N2-1.NAE4-N2A-DX1.OA-T", baCnetApplicationTagCharacterString.getValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(2);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            assertEquals(2201, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(85, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(3);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            assertEquals(2202, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagReal baCnetApplicationTagReal = (BACnetApplicationTagReal) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(35.093750, baCnetApplicationTagReal.getPayload().getValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(4);
                            assertEquals(BACnetPropertyIdentifier.RELIABILITY, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            BACnetReliability reliability = ((BACnetConstructedDataReliability) baCnetPropertyValue.getPropertyValue().getConstructedData()).getReliability().getValue();
                            assertEquals(BACnetReliability.NO_FAULT_DETECTED, reliability);
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(5);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            assertEquals(661, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(5, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(6);
                            assertEquals(BACnetPropertyIdentifier.UNITS, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            BACnetConstructedDataUnits baCnetConstructedDataUnits = (BACnetConstructedDataUnits) baCnetPropertyValue.getPropertyValue().getConstructedData();
                            assertEquals(BACnetEngineeringUnits.DEGREES_FAHRENHEIT, baCnetConstructedDataUnits.getUnits().getValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(7);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            assertEquals(1659, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(0, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(8);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            assertEquals(2203, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(0, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(9);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetPropertyValue.getPropertyIdentifier().getValue());
                            BACnetVendorIdTagged vendorId = ((BACnetConstructedDataVendorIdentifier) baCnetPropertyValue.getPropertyValue().getConstructedData()).getVendorIdentifier();
                            assertEquals(BACnetVendorId.JOHNSON_CONTROLS_INC, vendorId.getValue());
                        }
                    }
                })
        );
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("MSTP_Malformed_Packets")
    Collection<DynamicNode> MSTP_Malformed_Packets() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("MSTP_Malformed_Packets.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("NPDU")
    Collection<DynamicNode> NPDU() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("NPDU.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("PrivateTransferError-octetstring")
    Collection<DynamicNode> PrivateTransferError_octetstring() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("PrivateTransferError-octetstring.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("PrivateTransferError")
    Collection<DynamicNode> PrivateTransferError() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("PrivateTransferError.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("RPM_ALL_Allobjecttypes1")
    Collection<DynamicNode> RPM_ALL_Allobjecttypes1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("RPM_ALL_Allobjecttypes1.pcap");
        return List.of(pcapEvaluator.parseEmAll(
            skip(31, "set point reference looks completely broken")
        ));
    }

    @TestFactory
    @DisplayName("ReadPropertyMultiple")
    Collection<DynamicNode> ReadPropertyMultiple() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("ReadPropertyMultiple.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("ReadPropertyMultipleDeviceAll")
    Collection<DynamicNode> ReadPropertyMultipleDeviceAll() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("ReadPropertyMultipleDeviceAll.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("Subordinate List")
    Collection<DynamicNode> Subordinate_List() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Subordinate%20List.pcap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ readProperty[152] structured-view,1 subordinate-list",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apdu = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequestReadProperty serviceRequest = (BACnetConfirmedServiceRequestReadProperty) apdu.getServiceRequest();
                    assertThat(serviceRequest.getObjectIdentifier()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.STRUCTURED_VIEW, 1L);
                    assertThat(serviceRequest.getPropertyIdentifier()).extracting("value").isEqualTo(BACnetPropertyIdentifier.SUBORDINATE_LIST);
                }),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[152] structured-view,1 subordinate-list device,128 analog-input,1 device,128 analog-input,3 device,128 analog-output,1 device,128 analog-output,3",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAckReadProperty);
                    assertThat(baCnetServiceAckReadProperty.getObjectIdentifier()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.STRUCTURED_VIEW, 1L);
                    assertEquals(BACnetPropertyIdentifier.SUBORDINATE_LIST, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    List<BACnetDeviceObjectReference> data = ((BACnetConstructedDataSubordinateList) baCnetServiceAckReadProperty.getValues()).getSubordinateList();
                    assertThat(data)
                        .element(0)
                        .extracting(BACnetDeviceObjectReference::getDeviceIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(0)
                        .extracting(BACnetDeviceObjectReference::getObjectIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_INPUT, 1L);
                    assertThat(data)
                        .element(1)
                        .extracting(BACnetDeviceObjectReference::getDeviceIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(1)
                        .extracting(BACnetDeviceObjectReference::getObjectIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_INPUT, 3L);
                    assertThat(data)
                        .element(2)
                        .extracting(BACnetDeviceObjectReference::getDeviceIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(2)
                        .extracting(BACnetDeviceObjectReference::getObjectIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_OUTPUT, 1L);
                    assertThat(data)
                        .element(3)
                        .extracting(BACnetDeviceObjectReference::getDeviceIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(3)
                        .extracting(BACnetDeviceObjectReference::getObjectIdentifier)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_OUTPUT, 3L);
                })
        );
    }

    @TestFactory
    @DisplayName("Subordinate List2")
    Collection<DynamicNode> Subordinate_List2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Subordinate%20List2.pcap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   readProperty[143] structured-view,1 subordinate-list",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apdu = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequestReadProperty serviceRequest = (BACnetConfirmedServiceRequestReadProperty) apdu.getServiceRequest();
                    assertThat(serviceRequest.getObjectIdentifier()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.STRUCTURED_VIEW, 1L);
                    assertThat(serviceRequest.getPropertyIdentifier()).extracting("value").isEqualTo(BACnetPropertyIdentifier.SUBORDINATE_LIST);
                }),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[143] structured-view,1 subordinate-list device,4000 analog-input,1 analog-value,1 binary-input,1 binary-value,1",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAckReadProperty);
                    assertThat(baCnetServiceAckReadProperty.getObjectIdentifier()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.STRUCTURED_VIEW, 1L);
                    assertEquals(BACnetPropertyIdentifier.SUBORDINATE_LIST, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    List<BACnetDeviceObjectReference> data = ((BACnetConstructedDataSubordinateList) baCnetServiceAckReadProperty.getValues()).getSubordinateList();
                    assertThat(data).element(0).extracting(BACnetDeviceObjectReference::getDeviceIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 4000L);
                    assertThat(data).element(0).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_INPUT, 1L);
                    assertThat(data).element(1).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_VALUE, 1L);
                    assertThat(data).element(2).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.BINARY_INPUT, 1L);
                    assertThat(data).element(3).extracting(BACnetDeviceObjectReference::getObjectIdentifier).extracting("objectType", "instanceNumber").contains(BACnetObjectType.BINARY_VALUE, 1L);
                }),
            DynamicTest.dynamicTest("No. 3 - Confirmed-REQ   readProperty[144] structured-view,1 subordinate-annotations",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apdu = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequestReadProperty serviceRequest = (BACnetConfirmedServiceRequestReadProperty) apdu.getServiceRequest();
                    assertThat(serviceRequest.getObjectIdentifier()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.STRUCTURED_VIEW, 1L);
                    assertThat(serviceRequest.getPropertyIdentifier()).extracting("value").isEqualTo(BACnetPropertyIdentifier.SUBORDINATE_ANNOTATIONS);
                }),
            DynamicTest.dynamicTest("No. 4 - Complex-ACK     readProperty[144] structured-view,1 subordinate-annotations",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAckReadProperty);
                    assertThat(baCnetServiceAckReadProperty.getObjectIdentifier()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.STRUCTURED_VIEW, 1L);
                    assertEquals(BACnetPropertyIdentifier.SUBORDINATE_ANNOTATIONS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    List<BACnetApplicationTagCharacterString> data = ((BACnetConstructedDataSubordinateAnnotations) baCnetServiceAckReadProperty.getValues()).getSubordinateAnnotations();
                    assertThat(data.get(0)).extracting("value").isEqualTo("Subordinate 1");
                    assertThat(data.get(1)).extracting("value").isEqualTo("Subordinate 2");
                    assertThat(data.get(2)).extracting("value").isEqualTo("Subordinate 3");
                    assertThat(data.get(3)).extracting("value").isEqualTo("Subordinate 4");
                })
        );
    }

    @TestFactory
    @DisplayName("SubordinateListDecodeSample")
    Collection<DynamicNode> SubordinateListDecodeSample() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SubordinateListDecodeSample.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("SynergyBlinkWarn")
    Collection<DynamicNode> SynergyBlinkWarn() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SynergyBlinkWarn.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("SynergyReadProperties")
    Collection<DynamicNode> SynergyReadProperties() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SynergyReadProperties.cap");
        return List.of(pcapEvaluator.parseEmAll(
            skip(43, "seeems like a broken package"),
            skip(55, "TODO null as binaryPV looks wrong")
        ));
    }

    @TestFactory
    @DisplayName("SynergyWriteProperty")
    Collection<DynamicNode> SynergyWriteProperty() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SynergyWriteProperty.cap");
        return List.of(pcapEvaluator.parseEmAll(
            //3,
            //5
        ));
    }

    @TestFactory
    @DisplayName("Sysco-1")
    Collection<DynamicNode> Sysco_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Sysco-1.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("Sysco-2")
    Collection<DynamicNode> Sysco_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Sysco-2.cap");
        return List.of(pcapEvaluator.parseEmAll(
            // 25,
            // 26,
            // 27
        ));
    }

    @TestFactory
    @DisplayName("Sysco-3")
    Collection<DynamicNode> Sysco_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Sysco-3.cap");
        return List.of(pcapEvaluator.parseEmAll(
            //25,
            //26,
            //27
        ));
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("TA02 MST")
    Collection<DynamicNode> TA02_MST() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TA02%20MST.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("TAO2 TES.39A")
    Collection<DynamicNode> TAO2_TES_39_A() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TAO2%20TES.39A.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TC51103_BTL-9.21.1.X3_bool_ext_3")
    Collection<DynamicNode> TC51103_BTL_9_21_1_X3_bool_ext_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TC51103_BTL-9.21.1.X3_bool_ext_3.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TC51103_BTL-9.21.1.X3_int_ext_1")
    Collection<DynamicNode> TC51103_BTL_9_21_1_X3_int_ext_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TC51103_BTL-9.21.1.X3_int_ext_1.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TestRun4 - Internal side of Router")
    Collection<DynamicNode> TestRun4___Internal_side_of_Router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TestRun4%20-%20Internal%20side%20of%20Router.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TestRun4 - Outside of Router")
    Collection<DynamicNode> TestRun4___Outside_of_Router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TestRun4%20-%20Outside%20of%20Router.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TestRun5 - Internal side of Router")
    Collection<DynamicNode> TestRun5___Internal_side_of_Router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TestRun5%20-%20Internal%20side%20of%20Router.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TestRun5 - Outside of Router")
    Collection<DynamicNode> TestRun5___Outside_of_Router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TestRun5%20-%20Outside%20of%20Router.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TestRun8 - Internal side of Router")
    Collection<DynamicNode> TestRun8___Internal_side_of_Router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TestRun8%20-%20Internal%20side%20of%20Router.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TestRun8 - Outside of Router")
    Collection<DynamicNode> TestRun8___Outside_of_Router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TestRun8%20-%20Outside%20of%20Router.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TimeSync_Decode_Noon")
    Collection<DynamicNode> TimeSync_Decode_Noon() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TimeSync_Decode_Noon.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("Glorious Siemens implementation is a bit all over the place")
    @TestFactory
    @DisplayName("Tower333 lighting 5min IP")
    Collection<DynamicNode> Tower333_lighting_5min_IP() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Tower333%20lighting%205min%20IP.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            IntStream.of(
                    4, 7, 14, 15, 23, 28, 117, 118, 124, 126, 130, 131, 135, 166, 176, 177, 178, 180, 183, 185, 188, 194, 198,
                    199, 219, 223, 233, 254, 273, 284, 292, 294, 307, 311, 343, 345, 346, 347, 352, 354, 436, 437, 438,
                    447, 451, 454, 456, 459, 461, 468, 470, 517, 521, 602, 606, 624, 625, 650, 676, 700, 703, 704, 716,
                    718, 806, 815, 820, 851, 853, 856, 858, 860, 862, 885, 888, 891, 895, 921, 923, 924, 925, 926, 927, 928,
                    931, 933, 935, 936, 939, 941, 942, 944, 946, 948, 949, 982, 985, 986, 987, 1029, 1035, 1038, 1064, 1065, 1067, 1069,
                    1101, 1102, 1105, 1112, 1117, 1121, 1124, 1129, 1130, 1138, 1140, 1142, 1149, 1150, 1152, 1156, 1159,
                    1160, 1161, 1162, 1163, 1164, 1165, 1166, 1167, 1168, 1169, 1171, 1172, 1177, 1186, 1238, 1241, 1248, 1249,
                    1252, 1255, 1275, 1277, 1481, 1492, 1498, 1501, 1502, 1503, 1526, 1543, 1548, 1550, 1552, 1553, 1583, 1586, 1594,
                    1599, 1633, 1635, 1656, 1657, 1658, 1694, 1697, 1706, 1736, 1744, 1750, 1755, 1757, 1773, 1778, 1779,
                    1782, 1789, 1790, 1791, 1794, 1795, 1797, 1829, 1830, 1846, 1847, 1848, 1852, 1857, 1883, 1939, 1986,
                    1997, 1998, 1999, 2020, 2021, 2024, 2026, 2027, 2030, 2039, 2042, 2051, 2058, 2071, 2076, 2083, 2088,
                    2093, 2103, 2133, 2140, 2153, 2185, 2218, 2237, 2238, 2256, 2257, 2279, 2286, 2287, 2303, 2305, 2306, 2307,
                    2323, 2326, 2329, 2330, 2345, 2348, 2349, 2371, 2382, 2387, 2388, 2389, 2409, 2410, 2422, 2423, 2432,
                    2514, 2519, 2548, 2580, 2581, 2605, 2606, 2627, 2628, 2730, 2764, 2768, 2777, 2778, 2781, 2790, 2801,
                    2806, 2807, 2808, 2811, 2812, 2832, 2878, 2885, 2918, 2925, 2958, 2959, 2966, 2973, 2975, 3001, 3004, 3062,
                    3072, 3075, 3076, 3084, 3145, 3146, 3205, 3208, 3234, 3235, 3238, 3239, 3240, 3242, 3243, 3245, 3248, 3250, 3254,
                    3255, 3256, 3257, 3258, 3259, 3261, 3266, 3268, 3269, 3270, 3271, 3272, 3273, 3274, 3275, 3278, 3279,
                    3280, 3282, 3285, 3354, 3360, 3372, 3409, 3419, 3429, 3430, 3454, 3456, 3457, 3459, 3460, 3461, 3462, 3466,
                    3468, 3469, 3470, 3471, 3472, 3473, 3474, 3475, 3476, 3477, 3478, 3480, 3481, 3483, 3485, 3486, 3487, 3491,
                    3492, 3493, 3494, 3503, 3504, 3505, 3506, 3507, 3510, 3519, 3521, 3545, 3546, 3547, 3551, 3573, 3575,
                    3576, 3588, 3591, 3595, 3598, 3599, 3602, 3605, 3610, 3611, 3615, 3619, 3630, 3635, 3638, 3642, 3667,
                    3668, 3673, 3677, 3711, 3713, 3761, 3762, 3768, 3769, 3805, 3806, 3810, 3822, 3823, 3829, 3832, 3838, 3840,
                    3859, 3863, 3864, 3866, 3871, 3878, 3880, 3882, 3884, 3889, 3893, 3942, 3943, 3945, 3948, 3949, 3950,
                    3952, 3958, 3962, 3986, 3996, 4001, 4014, 4030, 4031, 4043, 4045, 4046, 4047, 4054, 4055, 4100, 4112, 4117,
                    4118, 4119, 4152, 4174, 4183, 4215, 4216, 4217, 4273, 4277, 4295, 4309, 4311, 4316, 4322, 4396, 4398,
                    4428, 4503, 4505, 4506, 4507, 4541, 4542, 4548, 4551, 4552, 4555, 4617, 4620, 4621, 4628, 4654, 4655,
                    4660, 4663, 4688, 4689, 4730, 4733, 4734, 4777, 4778, 4779, 4794, 4798, 5019, 5021, 5048, 5051, 5056, 5058, 5088,
                    5098, 5101, 5102, 5129, 5133, 5139, 5140, 5149, 5156, 5175, 5177, 5194, 5200, 5214, 5220, 5221, 5227, 5231,
                    5258, 5264, 5288, 5297, 5332, 5335, 5358, 5363, 5373, 5374, 5375, 5382, 5414, 5418, 5431, 5432, 5471,
                    5473, 5482, 5483, 5486, 5487, 5488, 5491, 5492, 5494, 5497, 5502, 5508, 5512, 5513, 5526, 5527, 5530,
                    5535, 5542, 5543, 5565, 5573, 5575, 5576, 5577, 5581, 5583, 5584, 5590, 5595, 5604, 5607, 5613, 5614,
                    5619, 5620, 5625, 5718, 5722, 5723, 5749, 5757, 5761, 5794, 5806, 5813, 5818, 5819, 5820, 5828, 5841,
                    5846, 5849, 5857, 5858, 5859, 5860, 5861, 5862, 5863, 5864, 5865, 5866, 5870, 5871, 5872, 5873, 5874, 5875,
                    5877, 5878, 5880, 5882, 5885, 5888, 5890, 5891, 5894, 5897, 5898, 5901, 5903, 5904, 5905, 5907, 5916,
                    5921, 5944, 5949, 5978, 5992, 5994, 6021, 6108, 6115, 6188, 6197, 6266, 6298, 6328, 6332, 6335, 6338,
                    6353, 6359, 6396, 6401, 6405, 6412, 6413, 6416, 6434, 6436, 6437, 6438, 6458, 6459, 6478, 6485, 6491, 6501,
                    6613, 6618, 6649, 6652, 6699, 6725, 6872, 6891, 7024, 7036, 7037, 7040, 7060, 7064, 7071, 7072, 7117,
                    7118, 7126, 7152, 7157, 7159, 7165, 7221, 7227, 7238, 7240, 7253, 7254, 7324, 7337, 7338, 7344, 7350, 7370,
                    7371, 7464, 7470, 7499, 7504, 7505, 7510, 7515, 7551, 7561, 7609, 7614, 7616, 7629, 7630, 7635, 7660, 7665,
                    7679, 7683, 7686, 7687, 7688, 7689, 7694, 7701, 7703, 7704, 7705, 7710, 7714, 7727, 7729, 7731, 7737, 7836,
                    7841, 7845, 7849, 7850, 7852, 7853, 7904, 7906, 7910, 7911, 7912, 7915, 7924, 7927, 7930, 7936, 7940,
                    7942, 7948, 7950, 7954, 7955, 7958, 7963, 7966, 7968, 7994, 7997, 8005, 8006, 8015, 8016, 8019, 8020, 8021, 8037,
                    8043, 8047, 8084, 8085, 8105, 8107, 8109, 8110, 8113, 8117, 8119, 8120, 8121, 8140, 8143, 8163, 8165,
                    8219, 8223, 8251, 8254, 8335, 8338, 8340, 8347, 8368, 8369, 8454, 8455, 8488, 8489, 8553, 8555, 8556, 8557,
                    8562, 8585, 8597, 8598, 8610, 8634, 8647, 8649, 8650, 8652, 8653, 8677, 8679, 8682, 8683, 8691, 8692,
                    8693, 8694, 8697, 8698, 8699, 8701, 8702, 8703, 8705, 8709, 8710, 8713, 8715, 8716, 8737, 8741, 8742,
                    8747, 8748, 8749, 8751, 8752, 8754, 8759, 8761, 8766, 8768, 8804, 8807, 8808, 8843, 8873, 8890, 8894,
                    8901, 8904, 8908, 8911, 8915, 8917, 8923, 8926, 8930, 8931, 8932, 8933, 8934, 8935, 8936, 8937, 8938,
                    8939, 8940, 8943, 8946, 8949, 8950, 8951, 8955, 8958, 8959, 8962, 8969, 8972, 8973, 8974, 8989, 8994,
                    9012, 9017, 9022, 9023, 9031, 9033, 9045)
                .mapToObj(i -> skip(i, SkipInstruction.SkipType.SKIP_COMPLETE, "Malformed Package. Siemens Implementation"))
                .toArray(SkipInstruction[]::new)
        ));
    }

    @TestFactory
    @DisplayName("TrendLogMultipleReadRange")
    Collection<DynamicNode> TrendLogMultipleReadRange() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TrendLogMultipleReadRange.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TrendLogMultipleReadRange2")
    Collection<DynamicNode> TrendLogMultipleReadRange2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TrendLogMultipleReadRange2.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TrendLogMultipleReadRange3")
    Collection<DynamicNode> TrendLogMultipleReadRange3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TrendLogMultipleReadRange3.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TrendLogMultipleReadRangeSimple")
    Collection<DynamicNode> TrendLogMultipleReadRangeSimple() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TrendLogMultipleReadRangeSimple.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("TrendLogMultipleUsage")
    Collection<DynamicNode> TrendLogMultipleUsage() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("TrendLogMultipleUsage.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("WhoIsRouterToNetwork-test")
    Collection<DynamicNode> WhoIsRouterToNetwork_test() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("WhoIsRouterToNetwork-test.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("WhoIsRouterToNetwork")
    Collection<DynamicNode> WhoIsRouterToNetwork() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("WhoIsRouterToNetwork.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("WhoIs_I-Am_Epics")
    Collection<DynamicNode> WhoIs_I_Am_Epics() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("WhoIs_I-Am_Epics.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("WireSharkError_ArrayIndex")
    Collection<DynamicNode> WireSharkError_ArrayIndex() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("WireSharkError_ArrayIndex.pcap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   readProperty[ 74] schedule,1 exception-schedule",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apdu = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequestReadProperty serviceRequest = (BACnetConfirmedServiceRequestReadProperty) apdu.getServiceRequest();
                    assertThat(serviceRequest.getObjectIdentifier())
                        .extracting(BACnetContextTagObjectIdentifier::getPayload)
                        .extracting(BACnetTagPayloadObjectIdentifier::getObjectType, BACnetTagPayloadObjectIdentifier::getInstanceNumber)
                        .contains(BACnetObjectType.SCHEDULE, 1L);
                    assertThat(serviceRequest.getPropertyIdentifier())
                        .extracting(BACnetPropertyIdentifierTagged::getValue)
                        .isEqualTo(BACnetPropertyIdentifier.EXCEPTION_SCHEDULE);
                    assertThat(serviceRequest.getArrayIndex()).extracting("payload").extracting("actualValue").isEqualTo(BigInteger.ZERO);
                }),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[ 74] schedule,1 exception-schedule",
                () -> {
                    if (true) {
                        throw new TestAbortedException("This is wrong. The exception schedule is a BACnetSpecialEvent and not a simple ApplicationTag");
                    }
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAckReadProperty);
                    assertThat(baCnetServiceAckReadProperty.getObjectIdentifier())
                        .extracting(BACnetContextTagObjectIdentifier::getPayload)
                        .extracting(BACnetTagPayloadObjectIdentifier::getObjectType, BACnetTagPayloadObjectIdentifier::getInstanceNumber)
                        .contains(BACnetObjectType.SCHEDULE, 1L);
                    assertEquals(BACnetPropertyIdentifier.EXCEPTION_SCHEDULE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    List<BACnetConstructedDataElement> data = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData();
                    assertThat(data.get(0).getApplicationTag()).extracting("payload").extracting("actualValue").isEqualTo(BigInteger.ZERO);
                })
        );
    }

    @TestFactory
    @DisplayName("WireSharkError_BufferReadyNotification")
    Collection<DynamicNode> WireSharkError_BufferReadyNotification() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("WireSharkError_BufferReadyNotification.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("WireSharkOfNewObjects")
    Collection<DynamicNode> WireSharkOfNewObjects() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("WireSharkOfNewObjects.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("WriteProperty_BinaryOutput")
    Collection<DynamicNode> WriteProperty_BinaryOutput() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("WriteProperty_BinaryOutput.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("action-list")
    Collection<DynamicNode> action_list() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("action-list.pcap");
        return List.of(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   readProperty[107] command,1 action",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUConfirmedRequest.class))
                    .extracting(APDUConfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetConfirmedServiceRequestReadProperty.class))
                    .satisfies(baCnetConfirmedServiceRequestReadProperty -> {
                        assertThat(baCnetConfirmedServiceRequestReadProperty)
                            .extracting(BACnetConfirmedServiceRequestReadProperty::getObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.COMMAND, 1L);
                        assertThat(baCnetConfirmedServiceRequestReadProperty)
                            .extracting(BACnetConfirmedServiceRequestReadProperty::getPropertyIdentifier)
                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                            .isEqualTo(BACnetPropertyIdentifier.ACTION);
                        assertThat(baCnetConfirmedServiceRequestReadProperty)
                            .extracting(BACnetConfirmedServiceRequestReadProperty::getArrayIndex)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ONE);
                    })),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[107] command,1 action binary-value,0 present-value device,0 analog-value,1 present-value binary-value,2 present-value device,0 analog-value,3 present-value binary-value,4 present-value device,0 analog-value,5 present-value binary-value,6 present-value device,0 analog-value,7 present-value binary-value,8 present-value device,0 analog-value,9 present-value",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getServiceAck)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetServiceAckReadProperty.class))
                    .satisfies(baCnetServiceAckReadProperty -> {
                        assertThat(baCnetServiceAckReadProperty.getObjectIdentifier())
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .contains(BACnetObjectType.COMMAND, 1L);
                        assertThat(baCnetServiceAckReadProperty.getPropertyIdentifier())
                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                            .isEqualTo(BACnetPropertyIdentifier.ACTION);
                        assertThat(baCnetServiceAckReadProperty.getValues())
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataCommandAction.class))
                            .extracting(BACnetConstructedDataCommandAction::getActionLists)
                            .satisfies(baCnetActionLists -> {
                                assertThat(baCnetActionLists).hasSize(1);
                                BACnetActionList baCnetActionList = baCnetActionLists.get(0);
                                assertThat(baCnetActionList.getAction()).isNotNull();
                                List<BACnetActionCommand> baCnetActionCommands = baCnetActionList.getAction();
                                assertThat(baCnetActionCommands)
                                    .element(0)
                                    .extracting(
                                        BACnetActionCommand::getObjectIdentifier,
                                        BACnetActionCommand::getPropertyIdentifier,
                                        BACnetActionCommand::getPropertyValue,
                                        BACnetActionCommand::getQuitOnFailure,
                                        BACnetActionCommand::getWriteSuccessful
                                    )
                                    .doesNotContainNull()
                                    .satisfiesExactly(
                                        objectIdentifier -> assertThat(objectIdentifier)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagObjectIdentifier.class))
                                            .extracting(BACnetContextTagObjectIdentifier::getObjectType)
                                            .isEqualTo(BACnetObjectType.BINARY_VALUE),
                                        propertyIdentifier -> assertThat(propertyIdentifier)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetPropertyIdentifierTagged.class))
                                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                                            .isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE),
                                        propertyValue -> assertThat(propertyValue)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataUnspecified.class))
                                            .extracting(BACnetConstructedDataUnspecified::getData)
                                            .asList().element(0)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataElement.class))
                                            .extracting(BACnetConstructedDataElement::getApplicationTag)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagReal.class))
                                            .extracting(BACnetApplicationTagReal::getPayload)
                                            .extracting(BACnetTagPayloadReal::getValue)
                                            .isEqualTo(0f),
                                        quitOnFailure -> assertThat(quitOnFailure)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagBoolean.class))
                                            .extracting(BACnetContextTagBoolean::getActualValue)
                                            .isEqualTo(false),
                                        writeSuccessful -> assertThat(writeSuccessful)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagBoolean.class))
                                            .extracting(BACnetContextTagBoolean::getActualValue)
                                            .isEqualTo(true)
                                    );
                                assertThat(baCnetActionCommands)
                                    .element(1)
                                    .extracting(
                                        BACnetActionCommand::getObjectIdentifier,
                                        BACnetActionCommand::getPropertyIdentifier,
                                        BACnetActionCommand::getPropertyValue,
                                        BACnetActionCommand::getQuitOnFailure,
                                        BACnetActionCommand::getWriteSuccessful
                                    )
                                    .doesNotContainNull()
                                    .satisfiesExactly(
                                        objectIdentifier -> assertThat(objectIdentifier)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagObjectIdentifier.class))
                                            .extracting(BACnetContextTagObjectIdentifier::getObjectType)
                                            .isEqualTo(BACnetObjectType.ANALOG_VALUE),
                                        propertyIdentifier -> assertThat(propertyIdentifier)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetPropertyIdentifierTagged.class))
                                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                                            .isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE),
                                        propertyValue -> assertThat(propertyValue)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataAnalogValuePresentValue.class))
                                            .extracting(BACnetConstructedDataAnalogValuePresentValue::getPresentValue)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagReal.class))
                                            .extracting(BACnetApplicationTagReal::getPayload)
                                            .extracting(BACnetTagPayloadReal::getValue)
                                            .isEqualTo(1.0f),
                                        quitOnFailure -> assertThat(quitOnFailure)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagBoolean.class))
                                            .extracting(BACnetContextTagBoolean::getActualValue)
                                            .isEqualTo(true),
                                        writeSuccessful -> assertThat(writeSuccessful)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagBoolean.class))
                                            .extracting(BACnetContextTagBoolean::getActualValue)
                                            .isEqualTo(false)
                                    );
                                assertThat(baCnetActionCommands)
                                    .element(2)
                                    .extracting(
                                        BACnetActionCommand::getObjectIdentifier,
                                        BACnetActionCommand::getPropertyIdentifier,
                                        BACnetActionCommand::getPropertyValue,
                                        BACnetActionCommand::getQuitOnFailure,
                                        BACnetActionCommand::getWriteSuccessful
                                    )
                                    .doesNotContainNull()
                                    .satisfiesExactly(
                                        objectIdentifier -> assertThat(objectIdentifier)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagObjectIdentifier.class))
                                            .extracting(BACnetContextTagObjectIdentifier::getObjectType)
                                            .isEqualTo(BACnetObjectType.BINARY_VALUE),
                                        propertyIdentifier -> assertThat(propertyIdentifier)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetPropertyIdentifierTagged.class))
                                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                                            .isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE),
                                        propertyValue -> assertThat(propertyValue)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataUnspecified.class))
                                            .extracting(BACnetConstructedDataUnspecified::getData)
                                            .asList().element(0)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataElement.class))
                                            .extracting(BACnetConstructedDataElement::getApplicationTag)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagReal.class))
                                            .extracting(BACnetApplicationTagReal::getPayload)
                                            .extracting(BACnetTagPayloadReal::getValue)
                                            .isEqualTo(2.0f),
                                        quitOnFailure -> assertThat(quitOnFailure)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagBoolean.class))
                                            .extracting(BACnetContextTagBoolean::getActualValue)
                                            .isEqualTo(false),
                                        writeSuccessful -> assertThat(writeSuccessful)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagBoolean.class))
                                            .extracting(BACnetContextTagBoolean::getActualValue)
                                            .isEqualTo(true)
                                    );
                            });
                    }))
        );
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("aha_220_to_20_lost_b")
    Collection<DynamicNode> aha_220_to_20_lost_b() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("aha_220_to_20_lost_b.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("alerton-plugfest-2")
    Collection<DynamicNode> alerton_plugfest_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("alerton-plugfest-2.cap", BACNET_BPF_FILTER_UDP);
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Complex-ACK readProperty[155] device,42222 protocol-version",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getServiceAck)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetServiceAckReadProperty.class))
                    .satisfies(baCnetServiceAckReadProperty -> {
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 42222L);
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getPropertyIdentifier)
                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                            .isEqualTo(BACnetPropertyIdentifier.PROTOCOL_VERSION);
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getArrayIndex)
                            .isNull();
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getValues)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataProtocolVersion.class))
                            .extracting(BACnetConstructedDataProtocolVersion::getProtocolVersion)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagUnsignedInteger.class))
                            .extracting(BACnetApplicationTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint8)
                            .isEqualTo((short) 1);
                    })),
            pcapEvaluator.parseFrom(2
            )
        );
    }

    @TestFactory
    @DisplayName("alerton-plugfest-3")
    Collection<DynamicNode> alerton_plugfest_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("alerton-plugfest-3.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @Disabled("too fat will crash ide (250k entries)")
    @TestFactory
    @DisplayName("alerton-plugfest")
    Collection<DynamicNode> alerton_plugfest() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("alerton-plugfest.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("arf-empty-file")
    Collection<DynamicNode> arf_empty_file() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("arf-empty-file.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("atomic-empty-file")
    Collection<DynamicNode> atomic_empty_file() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-empty-file.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("atomic-read-file-50")
    Collection<DynamicNode> atomic_read_file_50() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-read-file-50.cap");
        return List.of(pcapEvaluator.parseEmAll(
            // 5,
            // 8,
            // 11
        ));
    }

    @TestFactory
    @DisplayName("atomic-read-file-50x1500k")
    Collection<DynamicNode> atomic_read_file_50x1500k() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-read-file-50x1500k.cap");
        return List.of(pcapEvaluator.parseEmAll(
            //  35,
            //  60,
            //  75,
            //  173,
            //  201,
            //  216
        ));
    }

    @TestFactory
    @DisplayName("atomic-read-file-480")
    Collection<DynamicNode> atomic_read_file_480() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-read-file-480.cap");
        return List.of(pcapEvaluator.parseEmAll(
            //  5,
            //  8,
            //  9
        ));
    }

    @TestFactory
    @DisplayName("atomic-read-file-1470")
    Collection<DynamicNode> atomic_read_file_1470() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-read-file-1470.cap", BACNET_BPF_FILTER_UDP);
        // TODO: we have udps fragements here. Not sure how to handle it
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("atomic-read-file")
    Collection<DynamicNode> atomic_read_file() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-read-file.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("atomic-write-file-2")
    Collection<DynamicNode> atomic_write_file_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-write-file-2.cap");
        return List.of(pcapEvaluator.parseEmAll(
            //  1,
            //  3,
            //  6,
            //  1173,
            //  1179,
            //  1180
        ));
    }

    @TestFactory
    @DisplayName("atomic-write-file-3")
    Collection<DynamicNode> atomic_write_file_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-write-file-3.cap");
        return List.of(pcapEvaluator.parseEmAll(
            // 3,
            // 4,
            // 67,
            // 71
        ));
    }

    @TestFactory
    @DisplayName("atomic-write-file-50x1k")
    Collection<DynamicNode> atomic_write_file_50x1k() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-write-file-50x1k.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("atomic-write-file-480")
    Collection<DynamicNode> atomic_write_file_480() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-write-file-480.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("atomic-write-file-seg")
    Collection<DynamicNode> atomic_write_file_seg() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-write-file-seg.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("atomic-write-file")
    Collection<DynamicNode> atomic_write_file() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-write-file.cap");
        return List.of(pcapEvaluator.parseEmAll(
            //  1,
            //  4,
            //  3122
        ));
    }

    @TestFactory
    @DisplayName("atomic_write_file_bad_ack")
    Collection<DynamicNode> atomic_write_file_bad_ack() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic_write_file_bad_ack.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bacapp-malform")
    Collection<DynamicNode> bacapp_malform() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacapp-malform.cap", BACNET_BPF_FILTER_UDP);
        SkipInstruction[] undefinedExtraBytes = IntStream.of(
                1, 5, 13, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64,
                66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112,
                114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134, 136, 138, 140, 142, 144, 146, 148, 150, 152,
                154, 156, 158, 160, 162, 164, 166, 168, 170, 172, 174, 176, 178, 180, 182, 184, 186, 188, 190, 192,
                194, 196, 198, 200, 202, 204, 206, 208, 210, 212, 214, 216, 218, 220, 222, 224, 226, 228, 230, 232,
                234, 236, 238, 240, 242, 244, 246, 248, 250, 252, 254, 256, 258, 260, 262, 264, 266, 268, 270, 272,
                274, 276, 278, 280, 282, 284, 286, 288, 290, 292, 294, 296, 298, 300, 302, 304, 306, 308, 310, 312,
                314, 316, 318, 320, 322, 324, 326, 328, 330, 332, 334, 336, 338, 340, 342, 344, 346, 348, 350, 352,
                354, 356, 358, 360, 362, 364, 366, 368, 370, 372, 374, 376, 378, 380, 382, 384, 386, 388, 390, 392,
                394, 396, 398, 400, 402, 404, 406, 408, 410, 412, 414, 416, 418, 420, 422, 424, 426, 428, 430, 432,
                434, 436, 438, 440, 442, 444, 446, 448, 450, 452, 454, 456, 458, 460, 462, 464, 466, 468, 470, 472,
                474, 476, 478, 480, 482, 484, 486, 488, 490, 492, 494, 496, 498, 500, 502, 504, 506, 508, 510, 512,
                514, 516, 518, 520, 522, 524, 526, 528, 530, 532, 534, 536, 538, 540, 542, 544, 546, 548, 550, 552,
                554, 556, 558, 560, 562, 564, 566, 568, 570, 572, 574, 576, 578, 580, 582, 584, 586, 588, 590, 592,
                594, 596, 598, 600, 602, 604, 606, 608, 610, 612, 614, 616, 618, 620, 622, 624, 626, 628, 630, 632,
                634, 636, 638, 640, 642, 644, 646, 648, 650, 652, 654, 656, 658, 660, 662, 664, 666, 668, 670, 672,
                674, 676, 678, 680, 682, 684, 686, 688, 690, 692, 694, 696, 698, 700, 702, 704, 706, 708, 710, 712,
                714, 716, 718, 720, 722, 724, 726, 728, 730, 732, 734, 736, 738, 740, 742, 744, 746, 748, 750, 752,
                754, 756, 758, 760, 762, 764, 766, 768, 770, 772, 774, 776, 778, 780, 782, 784, 786, 788, 790, 792,
                794, 796, 798, 800, 802, 804, 806, 808, 810, 812, 814, 816, 818, 820, 822, 824, 826, 828, 830, 832,
                834, 836, 838, 840, 842, 844, 846, 848, 850, 852, 854, 856, 858, 860, 862, 864, 866, 868, 870, 872,
                874, 876, 878, 880, 882, 884, 886, 888, 890, 892, 894, 896, 898, 900, 902, 904, 906, 908, 910, 912,
                914, 916, 918, 920, 922, 924, 926, 928, 930, 932, 934, 936, 938, 940, 942, 944, 946, 948, 950, 952,
                954, 956, 958, 960, 962, 964, 966, 968, 970, 972, 974, 976, 978, 980, 982, 984, 986, 988, 990, 992,
                994, 996, 998, 1000, 1002, 1004, 1006, 1008, 1010, 1012, 1014, 1016, 1018, 1020, 1022, 1024, 1026,
                1028, 1030, 1032, 1034, 1036, 1038, 1040, 1042, 1044, 1046, 1048, 1050, 1052, 1054, 1056, 1058, 1060,
                1062, 1064, 1066, 1068, 1070, 1072, 1074, 1076, 1078, 1080, 1082, 1084, 1086, 1088, 1090, 1092, 1094,
                1096, 1098, 1100, 1102, 1104, 1106, 1108, 1110, 1112, 1114, 1116, 1118, 1120, 1122, 1124, 1126, 1128,
                1130, 1132, 1134, 1136, 1138, 1140, 1142, 1144, 1146, 1148, 1150, 1152, 1154, 1156, 1158, 1160, 1162,
                1164, 1166, 1168, 1170, 1172, 1174, 1176, 1178, 1180, 1182, 1184, 1186, 1188, 1190, 1192, 1194, 1196,
                1198, 1200, 1202, 1204, 1206, 1208, 1210, 1212, 1214, 1216, 1218, 1220, 1222, 1224, 1226, 1228, 1230,
                1232, 1234, 1236, 1238, 1240, 1242, 1244, 1246, 1248, 1250, 1252, 1254, 1256, 1258, 1260, 1262, 1264,
                1266, 1268, 1270, 1272, 1274, 1276, 1278, 1280, 1282, 1284, 1286, 1288, 1290, 1292, 1294, 1296, 1298,
                1300, 1302, 1304, 1306, 1308, 1310, 1312, 1314, 1316, 1318, 1320, 1322, 1324, 1326, 1328, 1330, 1332,
                1334, 1336, 1338, 1340, 1342, 1344, 1346, 1348, 1350, 1352, 1354, 1356, 1358, 1360, 1362, 1364, 1366,
                1368, 1370, 1372, 1374, 1376, 1378, 1380, 1382, 1384, 1386, 1388, 1390, 1392, 1394, 1396, 1398, 1400,
                1402, 1404, 1406, 1408, 1410, 1412, 1414, 1416, 1418, 1420, 1422, 1424, 1426, 1428, 1430, 1432, 1434,
                1436, 1438, 1440, 1442, 1444, 1446, 1448, 1450, 1452, 1454, 1456, 1458, 1460, 1462, 1464, 1466, 1468,
                1470, 1472, 1474, 1476, 1478, 1480, 1482, 1484, 1486, 1488, 1490, 1492, 1494, 1496, 1498, 1500, 1502,
                1504, 1506, 1508, 1510, 1512, 1514, 1516, 1518, 1520, 1522, 1524, 1526, 1528, 1530, 1532, 1534, 1536,
                1538, 1540, 1542, 1544, 1546, 1548, 1550, 1552, 1554, 1556, 1558, 1560, 1562, 1564, 1566, 1568, 1570,
                1572, 1574, 1576, 1578, 1580, 1582, 1584, 1586, 1588, 1590, 1592, 1594, 1596, 1598, 1600, 1602, 1604,
                1606, 1608, 1610, 1612, 1614, 1616, 1618, 1620, 1622, 1624, 1626, 1628, 1630, 1632, 1634, 1636, 1638,
                1640, 1642, 1644, 1646, 1648, 1650, 1652, 1654, 1656, 1658, 1660, 1662, 1664, 1666, 1668, 1670, 1672,
                1674, 1676, 1678, 1680, 1682
            )
            .mapToObj(i -> skip(i, SkipInstruction.SkipType.SKIP_COMPARE, "most of those packages contain extra undefined bytes"))
            .toArray(SkipInstruction[]::new);
        return List.of(pcapEvaluator.parseEmAll(
            ArrayUtils.addAll(undefinedExtraBytes
            )
        ));
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("bacnet-arcnet")
    Collection<DynamicNode> bacnet_arcnet() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-arcnet.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bacnet-ethernet-device")
    Collection<DynamicNode> bacnet_ethernet_device() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-ethernet-device.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bacnet-ethernet")
    Collection<DynamicNode> bacnet_ethernet() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-ethernet.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bacnet-ip")
    Collection<DynamicNode> bacnet_ip() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-ip.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("bacnet-properties")
    Collection<DynamicNode> bacnet_properties() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-properties.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bacnet-services")
    Collection<DynamicNode> bacnet_services() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-services.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("bacnet-stack-services")
    Collection<DynamicNode> bacnet_stack_services() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-stack-services.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(1, SkipInstruction.SkipType.SKIP_COMPARE, "contains extra bytes that we don't serialize"),
            skip(77, "Malformed Package"),
            skip(79, "Malformed Package"),
            skip(81, "Malformed Package"),
            skip(83, "Malformed Package")
        ));
    }

    @TestFactory
    @DisplayName("bacrpm-test")
    Collection<DynamicNode> bacrpm_test() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacrpm-test.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("bad_hub_restart")
    Collection<DynamicNode> bad_hub_restart() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bad_hub_restart.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bip-discover")
    Collection<DynamicNode> bip_discover() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bip-discover.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bip-readprop-2")
    Collection<DynamicNode> bip_readprop_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bip-readprop-2.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bip-readprop")
    Collection<DynamicNode> bip_readprop() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bip-readprop.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bip-readwrite-test")
    Collection<DynamicNode> bip_readwrite_test() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bip-readwrite-test.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("bo_command_failure")
    Collection<DynamicNode> bo_command_failure() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bo_command_failure.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bo_command_failure_original")
    Collection<DynamicNode> bo_command_failure_original() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bo_command_failure_original.pcap", BACNET_BPF_FILTER_UDP);
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 21 - Unconfirmed-REQ unconfirmedCOVNotification device,1 accumulator,21 present-value status-flags",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getSubscriberProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getMonitoredObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.ACCUMULATOR, 21L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getLifetimeInSeconds)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint32)
                            .isEqualTo(3108189395L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getListOfValues)
                            .extracting(BACnetPropertyValues::getData)
                            .satisfies(baCnetPropertyValues -> {
                                assertThat(baCnetPropertyValues).element(0).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataPresentValue.class))
                                        .extracting(BACnetConstructedDataPresentValue::getPresentValue)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagUnsignedInteger.class))
                                        .extracting(BACnetApplicationTagUnsignedInteger::getPayload)
                                        .extracting(BACnetTagPayloadUnsignedInteger::getValueUint16)
                                        .isEqualTo(1576);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                                assertThat(baCnetPropertyValues).element(1).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.STATUS_FLAGS);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataStatusFlags.class))
                                        .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                                        .satisfies(statusFlagsTagged ->
                                            assertThat(statusFlagsTagged)
                                                .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                                .isEqualTo(List.of(false, false, false, false)));
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                            });
                    })),
            DynamicTest.dynamicTest("No. 22 - Unconfirmed-REQ unconfirmedCOVNotification device,1 accumulator,22 present-value status-flags",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getSubscriberProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getMonitoredObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.ACCUMULATOR, 22L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getLifetimeInSeconds)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint32)
                            .isEqualTo(3108189395L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getListOfValues)
                            .extracting(BACnetPropertyValues::getData)
                            .satisfies(baCnetPropertyValues -> {
                                assertThat(baCnetPropertyValues).element(0).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataPresentValue.class))
                                        .extracting(BACnetConstructedDataPresentValue::getPresentValue)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagUnsignedInteger.class))
                                        .extracting(BACnetApplicationTagUnsignedInteger::getPayload)
                                        .extracting(BACnetTagPayloadUnsignedInteger::getValueUint16)
                                        .isEqualTo(1577);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                                assertThat(baCnetPropertyValues).element(1).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.STATUS_FLAGS);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataStatusFlags.class))
                                        .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                                        .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                        .isEqualTo(List.of(false, false, false, false));
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                            });
                    })),
            DynamicTest.dynamicTest("No. 23 - Unconfirmed-REQ unconfirmedCOVNotification device,1 binary-input,217 present-value status-flags",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getSubscriberProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getMonitoredObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.BINARY_INPUT, 217L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getLifetimeInSeconds)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint32)
                            .isEqualTo(3108189388L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getListOfValues)
                            .extracting(BACnetPropertyValues::getData)
                            .satisfies(baCnetPropertyValues -> {
                                assertThat(baCnetPropertyValues).element(0).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataBinaryInputPresentValue.class))
                                        .extracting(BACnetConstructedDataBinaryInputPresentValue::getPresentValue)
                                        .extracting(BACnetBinaryPVTagged::getValue)
                                        .isEqualTo(BACnetBinaryPV.ACTIVE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                                assertThat(baCnetPropertyValues).element(1).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.STATUS_FLAGS);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataStatusFlags.class))
                                        .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                                        .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                        .isEqualTo(List.of(false, false, false, false));
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                            });
                    })),
            DynamicTest.dynamicTest("No. 24 - Unconfirmed-REQ unconfirmedCOVNotification device,1 accumulator,21 present-value status-flags",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getSubscriberProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getMonitoredObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.ACCUMULATOR, 21L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getLifetimeInSeconds)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint32)
                            .isEqualTo(3108189388L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getListOfValues)
                            .extracting(BACnetPropertyValues::getData)
                            .satisfies(baCnetPropertyValues -> {
                                assertThat(baCnetPropertyValues).element(0).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataPresentValue.class))
                                        .extracting(BACnetConstructedDataPresentValue::getPresentValue)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagUnsignedInteger.class))
                                        .extracting(BACnetApplicationTagUnsignedInteger::getPayload)
                                        .extracting(BACnetTagPayloadUnsignedInteger::getValueUint16)
                                        .isEqualTo(1577);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                                assertThat(baCnetPropertyValues).element(1).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.STATUS_FLAGS);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataStatusFlags.class))
                                        .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                                        .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                        .isEqualTo(List.of(false, false, false, false));
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                            });
                    })),
            DynamicTest.dynamicTest("No. 25 - Unconfirmed-REQ unconfirmedCOVNotification device,1 binary-input,217 present-value status-flags",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getSubscriberProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getMonitoredObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.BINARY_INPUT, 217L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getLifetimeInSeconds)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint32)
                            .isEqualTo(3108189388L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getListOfValues)
                            .extracting(BACnetPropertyValues::getData)
                            .satisfies(baCnetPropertyValues -> {
                                assertThat(baCnetPropertyValues).element(0).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataBinaryInputPresentValue.class))
                                        .extracting(BACnetConstructedDataBinaryInputPresentValue::getPresentValue)
                                        .extracting(BACnetBinaryPVTagged::getValue)
                                        .isEqualTo(BACnetBinaryPV.ACTIVE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                                assertThat(baCnetPropertyValues).element(1).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.STATUS_FLAGS);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataStatusFlags.class))
                                        .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                                        .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                        .isEqualTo(List.of(true, false, false, false));
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                            });
                    })),
            DynamicTest.dynamicTest("No. 26 - Unconfirmed-REQ unconfirmedCOVNotification device,1 binary-output,1 present-value status-flags",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getSubscriberProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getMonitoredObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.BINARY_OUTPUT, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getLifetimeInSeconds)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint32)
                            .isEqualTo(3108189383L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getListOfValues)
                            .extracting(BACnetPropertyValues::getData)
                            .satisfies(baCnetPropertyValues -> {
                                assertThat(baCnetPropertyValues).element(0).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataBinaryOutputPresentValue.class))
                                        .extracting(BACnetConstructedDataBinaryOutputPresentValue::getPresentValue)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetBinaryPVTagged.class))
                                        .extracting(BACnetBinaryPVTagged::getValue)
                                        .isEqualTo(BACnetBinaryPV.INACTIVE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                                assertThat(baCnetPropertyValues).element(1).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.STATUS_FLAGS);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataStatusFlags.class))
                                        .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                                        .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                        .isEqualTo(List.of(true, false, false, false));
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                            });
                    })),
            DynamicTest.dynamicTest("No. 27 - Unconfirmed-REQ unconfirmedEventNotification device,1 binary-output,1",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getEventObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.BINARY_OUTPUT, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getTimestamp)
                            .extracting(BACnetTimeStampEnclosed::getTimestamp)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetTimeStampDateTime.class))
                            .extracting(BACnetTimeStampDateTime::getDateTimeValue)
                            .extracting(BACnetDateTimeEnclosed::getDateTimeValue)
                            .satisfies(baCnetDateTime -> {
                                assertThat(baCnetDateTime)
                                    .extracting(BACnetDateTime::getDateValue)
                                    .extracting(BACnetApplicationTagDate::getPayload)
                                    .extracting(BACnetTagPayloadDate::getYearMinus1900, BACnetTagPayloadDate::getMonth, BACnetTagPayloadDate::getDayOfMonth, BACnetTagPayloadDate::getDayOfWeek)
                                    .contains((short) 107, (short) 8, (short) 10, (short) 5);
                                assertThat(baCnetDateTime)
                                    .extracting(BACnetDateTime::getTimeValue)
                                    .extracting(BACnetApplicationTagTime::getPayload)
                                    .extracting(BACnetTagPayloadTime::getHour, BACnetTagPayloadTime::getMinute, BACnetTagPayloadTime::getSecond, BACnetTagPayloadTime::getFractional)
                                    .contains((short) 20, (short) 31, (short) 53, (short) 12);
                            });
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getNotificationClass)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint8)
                            .isEqualTo((short) 1);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getPriority)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint8)
                            .isEqualTo((short) 0);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getEventType)
                            .extracting(BACnetEventTypeTagged::getValue)
                            .isEqualTo(BACnetEventType.COMMAND_FAILURE);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getMessageText)
                            .extracting(BACnetContextTagCharacterString::getPayload)
                            .extracting(BACnetTagPayloadCharacterString::getValue)
                            .isEqualTo("BO_1");
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getNotifyType)
                            .extracting(BACnetNotifyTypeTagged::getValue)
                            .isEqualTo(BACnetNotifyType.ALARM);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getAckRequired)
                            .extracting(BACnetContextTagBoolean::getActualValue)
                            .isEqualTo(false);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getFromState)
                            .extracting(BACnetEventStateTagged::getValue)
                            .isEqualTo(BACnetEventState.NORMAL);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getToState)
                            .extracting(BACnetEventStateTagged::getValue)
                            .isEqualTo(BACnetEventState.OFFNORMAL);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification::getEventValues)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetNotificationParametersCommandFailure.class))
                            .satisfies(baCnetNotificationParametersCommandFailure -> {
                                assertThat(baCnetNotificationParametersCommandFailure)
                                    .extracting(BACnetNotificationParametersCommandFailure::getCommandValue)
                                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataUnspecified.class))
                                    .extracting(BACnetConstructedDataUnspecified::getData)
                                    .satisfies(baCnetConstructedDataElements ->
                                        assertThat(baCnetConstructedDataElements)
                                            .element(0)
                                            .extracting(BACnetConstructedDataElement::getApplicationTag)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagEnumerated.class))
                                            .extracting(BACnetApplicationTagEnumerated::getActualValue)
                                            .isEqualTo(0L)
                                    );
                                assertThat(baCnetNotificationParametersCommandFailure)
                                    .extracting(BACnetNotificationParametersCommandFailure::getStatusFlags)
                                    .extracting(BACnetStatusFlagsTagged::getInAlarm)
                                    .isEqualTo(true);
                                assertThat(baCnetNotificationParametersCommandFailure)
                                    .extracting(BACnetNotificationParametersCommandFailure::getFeedbackValue)
                                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataUnspecified.class))
                                    .extracting(BACnetConstructedDataUnspecified::getData)
                                    .satisfies(baCnetConstructedDataElements ->
                                        assertThat(baCnetConstructedDataElements)
                                            .element(0)
                                            .extracting(BACnetConstructedDataElement::getApplicationTag)
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagEnumerated.class))
                                            .extracting(BACnetApplicationTagEnumerated::getActualValue)
                                            .isEqualTo(1L)
                                    );
                            });
                    })),
            DynamicTest.dynamicTest("No. 28 - Unconfirmed-REQ unconfirmedCOVNotification device,1 accumulator,22 present-value status-flags",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalBroadcastNPDU.class))
                    .extracting(BVLCOriginalBroadcastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUUnconfirmedRequest.class))
                    .extracting(APDUUnconfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification.class))
                    .satisfies(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification -> {
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getSubscriberProcessIdentifier)
                            .extracting(BACnetContextTagUnsignedInteger::getActualValue)
                            .isEqualTo(BigInteger.ZERO);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getInitiatingDeviceIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.DEVICE, 1L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getMonitoredObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.ACCUMULATOR, 22L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getLifetimeInSeconds)
                            .extracting(BACnetContextTagUnsignedInteger::getPayload)
                            .extracting(BACnetTagPayloadUnsignedInteger::getValueUint32)
                            .isEqualTo(3108189381L);
                        assertThat(baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            .extracting(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification::getListOfValues)
                            .extracting(BACnetPropertyValues::getData)
                            .satisfies(baCnetPropertyValues -> {
                                assertThat(baCnetPropertyValues).element(0).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.PRESENT_VALUE);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataPresentValue.class))
                                        .extracting(BACnetConstructedDataPresentValue::getPresentValue)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagUnsignedInteger.class))
                                        .extracting(BACnetApplicationTagUnsignedInteger::getPayload)
                                        .extracting(BACnetTagPayloadUnsignedInteger::getValueUint16)
                                        .isEqualTo(1578);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                                assertThat(baCnetPropertyValues).element(1).satisfies(baCnetPropertyValue -> {
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyIdentifier).extracting(BACnetPropertyIdentifierTagged::getValue).isEqualTo(BACnetPropertyIdentifier.STATUS_FLAGS);
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyArrayIndex).isNull();
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPropertyValue)
                                        .extracting(BACnetConstructedDataElement::getConstructedData)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataStatusFlags.class))
                                        .extracting(BACnetConstructedDataStatusFlags::getStatusFlags)
                                        .extracting(BACnetStatusFlagsTagged::getInAlarm, BACnetStatusFlagsTagged::getFault, BACnetStatusFlagsTagged::getOverridden, BACnetStatusFlagsTagged::getOutOfService)
                                        .isEqualTo(List.of(false, false, false, false));
                                    assertThat(baCnetPropertyValue).extracting(BACnetPropertyValue::getPriority).isNull();
                                });
                            });
                    }))
        );
    }

    @TestFactory
    @DisplayName("btl-plugfest")
    Collection<DynamicNode> btl_plugfest() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("btl-plugfest.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("bvlc-bac4-rp")
    Collection<DynamicNode> bvlc_bac4_rp() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bvlc-bac4-rp.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bvlc-bac4")
    Collection<DynamicNode> bvlc_bac4() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bvlc-bac4.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("bvlc-fdreg-readprop-47809")
    Collection<DynamicNode> bvlc_fdreg_readprop_47809() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bvlc-fdreg-readprop-47809.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("bvlc-loop")
    Collection<DynamicNode> bvlc_loop() throws Exception {
        String filterBrokenUDPPackages = "udp[4:2] > 29";
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bvlc-loop.cap", BACNET_BPF_FILTER_UDP + " and " + filterBrokenUDPPackages);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no useful")
    @TestFactory
    @DisplayName("bvlc")
    Collection<DynamicNode>
    bvlc() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bvlc.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("cimetrics_mstp")
    Collection<DynamicNode> cimetrics_mstp() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("cimetrics_mstp.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("config-tool-discover")
    Collection<DynamicNode> config_tool_discover() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("config-tool-discover.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("confirmedEventNotification")
    Collection<DynamicNode> confirmedEventNotification() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("confirmedEventNotification.cap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ confirmedEventNotification[103] device,1041000 analog-input,3000016 present-value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestConfirmedEventNotification baCnetConfirmedServiceRequestConfirmedEventNotification = (BACnetConfirmedServiceRequestConfirmedEventNotification) serviceRequest;
                    assertEquals((short) 123, baCnetConfirmedServiceRequestConfirmedEventNotification.getProcessIdentifier().getPayload().getValueUint8());
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getObjectType());
                    assertEquals(1041000, baCnetConfirmedServiceRequestConfirmedEventNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.ANALOG_INPUT, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getObjectType());
                    assertEquals(3000016, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventObjectIdentifier().getInstanceNumber());
                    {
                        BACnetTimeStampTime timestamp = (BACnetTimeStampTime) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp().getTimestamp();
                        assertEquals(2, timestamp.getTimeValue().getPayload().getHour());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(200, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getValue());
                    }
                    {
                        assertEquals(BACnetNotifyType.ALARM, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertFalse(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsFalse());
                    }
                    {
                        assertEquals(BACnetEventState.HIGH_LIMIT, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getValue());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getValue());
                    }
                    {
                        BACnetNotificationParametersComplexEventType baCnetNotificationParametersComplexEventType = (BACnetNotificationParametersComplexEventType) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertNotNull(baCnetNotificationParametersComplexEventType);
                    }
                }),
            pcapEvaluator.parseFrom(2)
        );
    }

    @TestFactory
    @DisplayName("cov-testing-1")
    Collection<DynamicNode> cov_testing_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("cov-testing-1.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("cov-testing-2")
    Collection<DynamicNode> cov_testing_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("cov-testing-2.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("cov-testing-3")
    Collection<DynamicNode> cov_testing_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("cov-testing-3.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("device-address-binding")
    Collection<DynamicNode> device_address_binding() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("device-address-binding.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("epics-1")
    Collection<DynamicNode> epics_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("epics-1.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("epics-2")
    Collection<DynamicNode> epics_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("epics-2.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("eventLog_ReadRange")
    Collection<DynamicNode> eventLog_ReadRange() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("eventLog_ReadRange.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("eventLog_rpm")
    Collection<DynamicNode> eventLog_rpm() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("eventLog_rpm.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("eventTimeStamp_rp")
    Collection<DynamicNode> eventTimeStamp_rp() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("eventTimeStamp_rp.pcap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   readProperty[148] load-control,1 event-time-stamp",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUConfirmedRequest.class))
                    .extracting(APDUConfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetConfirmedServiceRequestReadProperty.class))
                    .satisfies(baCnetConfirmedServiceRequestReadProperty -> {
                        assertThat(baCnetConfirmedServiceRequestReadProperty)
                            .extracting(BACnetConfirmedServiceRequestReadProperty::getObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.LOAD_CONTROL, 1L);
                        assertThat(baCnetConfirmedServiceRequestReadProperty)
                            .extracting(BACnetConfirmedServiceRequestReadProperty::getPropertyIdentifier)
                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                            .isEqualTo(BACnetPropertyIdentifier.EVENT_TIME_STAMPS);
                        assertThat(baCnetConfirmedServiceRequestReadProperty)
                            .extracting(BACnetConfirmedServiceRequestReadProperty::getArrayIndex)
                            .isNull();
                    })),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[148] load-control,1 event-time-stamp",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getServiceAck)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetServiceAckReadProperty.class))
                    .satisfies(baCnetServiceAckReadProperty -> {
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getObjectIdentifier)
                            .extracting(BACnetContextTagObjectIdentifier::getObjectType, BACnetContextTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.LOAD_CONTROL, 1L);
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getPropertyIdentifier)
                            .extracting(BACnetPropertyIdentifierTagged::getValue)
                            .isEqualTo(BACnetPropertyIdentifier.EVENT_TIME_STAMPS);
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getArrayIndex)
                            .isNull();
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getValues)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataEventTimeStamps.class))
                            .satisfies(baCnetConstructedDataEventTimestamps ->
                                assertThat(baCnetConstructedDataEventTimestamps)
                                    .extracting(BACnetConstructedDataEventTimeStamps::getEventTimeStamps)
                                    .satisfies(baCnetTimeStamps -> assertThat(baCnetTimeStamps)
                                        .element(0)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetTimeStampTime.class))
                                        .extracting(BACnetTimeStampTime::getTimeValue)
                                        .extracting(BACnetContextTagTime::getPayload)
                                        .extracting(BACnetTagPayloadTime::getHourIsWildcard, BACnetTagPayloadTime::getMinuteIsWildcard, BACnetTagPayloadTime::getSecondIsWildcard, BACnetTagPayloadTime::getFractionalIsWildcard)
                                        .containsExactly(true, true, true, true)
                                    )
                            )
                            .satisfies(baCnetConstructedDataEventTimestamps ->
                                assertThat(baCnetConstructedDataEventTimestamps)
                                    .extracting(BACnetConstructedDataEventTimeStamps::getEventTimeStamps)
                                    .satisfies(baCnetTimeStamps -> assertThat(baCnetTimeStamps)
                                        .hasSize(3)
                                    )
                            )
                            .satisfies(baCnetConstructedDataEventTimestamps ->
                                assertThat(baCnetConstructedDataEventTimestamps)
                                    .extracting(BACnetConstructedDataEventTimeStamps::getEventTimeStamps)
                                    .satisfies(baCnetTimeStamps -> assertThat(baCnetTimeStamps)
                                        .element(2)
                                        .asInstanceOf(InstanceOfAssertFactories.type(BACnetTimeStampDateTime.class))
                                        .extracting(BACnetTimeStampDateTime::getDateTimeValue)
                                        .satisfies(baCnetDateTime ->
                                            assertThat(baCnetDateTime)
                                                .extracting(BACnetDateTimeEnclosed::getDateTimeValue)
                                                .extracting(BACnetDateTime::getDateValue)
                                                .extracting(BACnetApplicationTagDate::getPayload)
                                                .extracting(BACnetTagPayloadDate::getYearIsWildcard, BACnetTagPayloadDate::getMonthIsWildcard, BACnetTagPayloadDate::getDayOfMonthIsWildcard, BACnetTagPayloadDate::getDayOfWeekIsWildcard)
                                                .containsExactly(true, true, true, true)
                                        )
                                        .satisfies(baCnetDateTime ->
                                            assertThat(baCnetDateTime)
                                                .extracting(BACnetDateTimeEnclosed::getDateTimeValue)
                                                .extracting(BACnetDateTime::getTimeValue)
                                                .extracting(BACnetApplicationTagTime::getPayload)
                                                .extracting(BACnetTagPayloadTime::getFractionalIsWildcard, BACnetTagPayloadTime::getMinuteIsWildcard, BACnetTagPayloadTime::getSecondIsWildcard, BACnetTagPayloadTime::getFractionalIsWildcard)
                                                .containsExactly(true, true, true, true)
                                        )
                                    )
                            );
                    }))
        );
    }

    @TestFactory
    @DisplayName("eventTimeStamp_rpm")
    Collection<DynamicNode> eventTimeStamp_rpm() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("eventTimeStamp_rpm.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("foreign-device-npdu")
    Collection<DynamicNode> foreign_device_npdu() throws Exception {
        String filterBrokenUDPPackages = "udp[4:2] > 29";
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("foreign-device-npdu.cap", BACNET_BPF_FILTER_UDP + " and " + filterBrokenUDPPackages);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("getEventInformation")
    Collection<DynamicNode> getEventInformation() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("getEventInformation.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("i-am-vendor-id-over-255")
    Collection<DynamicNode> i_am_vendor_id_over_255() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("i-am-vendor-id-over-255.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("lmbc-300-bootload")
    Collection<DynamicNode> lmbc_300_bootload() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("lmbc-300-bootload.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("load-control-properties")
    Collection<DynamicNode> load_control_properties() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("load-control-properties.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("load-control")
    Collection<DynamicNode> load_control() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("load-control.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("log-buffer_readRange")
    Collection<DynamicNode> log_buffer_readRange() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("log-buffer_readRange.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("loop2")
    Collection<DynamicNode> loop2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("loop2.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp-cimetrics")
    Collection<DynamicNode> mstp_cimetrics() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp-cimetrics.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("mstp-test-4")
    Collection<DynamicNode> mstp_test_4() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp-test-4.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("mstp-whois-basrt-mix")
    Collection<DynamicNode> mstp_whois_basrt_mix() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp-whois-basrt-mix.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("mstp-whois-basrt-mix2")
    Collection<DynamicNode> mstp_whois_basrt_mix2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp-whois-basrt-mix2.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("mstp-whois-iam")
    Collection<DynamicNode> mstp_whois_iam() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp-whois-iam.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20090227094623")
    Collection<DynamicNode> mstp_20090227094623() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20090227094623.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20090304105820")
    Collection<DynamicNode> mstp_20090304105820() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20090304105820.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20090304110410")
    Collection<DynamicNode> mstp_20090304110410() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20090304110410.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20090807145500")
    Collection<DynamicNode> mstp_20090807145500() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20090807145500.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013121352")
    Collection<DynamicNode> mstp_20091013121352() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013121352.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013121410")
    Collection<DynamicNode> mstp_20091013121410() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013121410.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013122053")
    Collection<DynamicNode> mstp_20091013122053() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013122053.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013122451")
    Collection<DynamicNode> mstp_20091013122451() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013122451.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013123108")
    Collection<DynamicNode> mstp_20091013123108() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013123108.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013124218")
    Collection<DynamicNode> mstp_20091013124218() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013124218.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013130259")
    Collection<DynamicNode> mstp_20091013130259() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013130259.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091013162906")
    Collection<DynamicNode> mstp_20091013162906() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091013162906.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091014093519")
    Collection<DynamicNode> mstp_20091014093519() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091014093519.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20091014112427")
    Collection<DynamicNode> mstp_20091014112427() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20091014112427.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_20140225214217")
    Collection<DynamicNode> mstp_20140225214217() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_20140225214217.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("mstp_mix_basrt_V124")
    Collection<DynamicNode> mstp_mix_basrt_V124() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_mix_basrt_V124.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @Disabled("seems like this pcap breaks the parser")
    @TestFactory
    @DisplayName("mstp_mix_basrt_V124_bad")
    Collection<DynamicNode> mstp_mix_basrt_V124_bad() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_mix_basrt_V124_bad.cap");
        return List.of(pcapEvaluator.parseEmAll(

        ));
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("mstp_wtap")
    Collection<DynamicNode> mstp_wtap() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("mstp_wtap.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("nb-binary-output")
    Collection<DynamicNode> nb_binary_output() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("nb-binary-output.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("plugfest-2011-delta-1")
    Collection<DynamicNode> plugfest_2011_delta_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-delta-1.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("plugfest-2011-delta-2")
    Collection<DynamicNode> plugfest_2011_delta_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-delta-2.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("plugfest-2011-delta-3")
    Collection<DynamicNode> plugfest_2011_delta_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-delta-3.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("plugfest-2011-mstp-roundtable")
    Collection<DynamicNode> plugfest_2011_mstp_roundtable() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-mstp-roundtable.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("plugfest-2011-sauter-1")
    Collection<DynamicNode> plugfest_2011_sauter_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-sauter-1.pcap");
        return List.of(pcapEvaluator.parseEmAll(
            skip(70, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(72, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(74, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(76, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(228, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(306, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(309, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... no idea what is wrong here"),
            skip(313, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(321, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... no idea what is wrong here"),
            skip(323, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... no idea what is wrong here"),
            skip(335, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(337, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(341, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(345, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(347, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(351, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(357, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(367, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(371, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(373, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(377, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(381, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(385, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(391, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(405, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(416, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(423, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(425, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(429, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(433, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(435, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(1250, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... no idea what is wrong here"),
            skip(1252, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... no idea what is wrong here"),
            skip(1256, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... no idea what is wrong here"),
            skip(1278, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(1280, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... no idea what is wrong here")
        ));
    }

    @TestFactory
    @DisplayName("plugfest-2011-siemens-1")
    Collection<DynamicNode> plugfest_2011_siemens_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-siemens-1.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(225, "strange siemens package"),
            skip(278, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(291, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(363, SkipInstruction.SkipType.SKIP_COMPARE, "extra bytes at the end which are undefined"),
            skip(448, SkipInstruction.SkipType.SKIP_COMPARE, "extra bytes at the end which are undefined"),
            skip(827, "This request the results in a error response as it sends wrong data"),
            skip(828, "This request the results in a error response as it sends wrong data"),
            skip(865, "This request the results in a error response as it sends wrong data"),
            skip(866, "This request the results in a error response as it sends wrong data"),
            skip(875, "This request the results in a error response as it sends wrong data"),
            skip(876, "This request the results in a error response as it sends wrong data"),
            skip(1575, SkipInstruction.SkipType.SKIP_COMPARE, "extra bytes at the end which are undefined"),
            skip(2327, SkipInstruction.SkipType.SKIP_COMPARE, "Broken utf-8... 0xae should be 0xc2ae"),
            skip(2329, "strange siemens package"),
            skip(2345, "strange siemens package"),
            skip(2586, SkipInstruction.SkipType.SKIP_COMPARE, "extra bytes at the end which are undefined"),
            skip(2626, SkipInstruction.SkipType.SKIP_COMPARE, "extra bytes at the end which are undefined")
        ));
    }

    @TestFactory
    @DisplayName("plugfest-2011-trane-1")
    Collection<DynamicNode> plugfest_2011_trane_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-trane-1.pcap", BACNET_BPF_FILTER_UDP);
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Unconfirmed REQ who-Is",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    assertTrue(serviceRequest instanceof BACnetUnconfirmedServiceRequestWhoIs);
                }),
            pcapEvaluator.parseFrom(2
            )
        );
    }

    @TestFactory
    @DisplayName("plugfest-delta-2")
    Collection<DynamicNode> plugfest_delta_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-delta-2.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("plugfest-delta-2b")
    Collection<DynamicNode> plugfest_delta_2b() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-delta-2b.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("plugfest-tridium-1")
    Collection<DynamicNode> plugfest_tridium_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-tridium-1.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(13, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(15, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(391, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(393, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(423, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(425, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(438, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(440, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(489, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(495, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(828, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(830, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(945, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(947, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(1084, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(1086, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(3405, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(3407, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252")
        ));
    }

    @TestFactory
    @DisplayName("plugfest-tridium-2")
    Collection<DynamicNode> plugfest_tridium_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-tridium-2.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(60, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(62, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(82, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(84, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(103, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(105, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(156, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(158, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(170, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(172, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(184, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(186, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(205, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(207, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(227, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(229, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(245, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(251, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(270, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(272, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(302, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(304, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(319, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(325, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(366, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(368, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(386, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(388, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(410, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(412, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(430, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(434, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(440, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(444, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(460, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(462, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(487, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(489, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(511, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(513, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(527, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(533, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(551, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(553, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(574, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(576, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(597, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(599, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(620, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252"),
            skip(622, SkipInstruction.SkipType.SKIP_COMPARE, "no utf-8, the encoding is just Cp1252")

        ));
    }

    @TestFactory
    @DisplayName("polarsoft-free-range-router-init-routing-table")
    Collection<DynamicNode> polarsoft_free_range_router_init_routing_table() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("polarsoft-free-range-router-init-routing-table.cap");
        return List.of(pcapEvaluator.parseEmAll(
            skip(2, SkipInstruction.SkipType.SKIP_COMPARE, "we can't reserialize because of unknown service")
        ));
    }

    @TestFactory
    @DisplayName("polarsoft-free-range-router")
    Collection<DynamicNode> polarsoft_free_range_router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("polarsoft-free-range-router.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(6155, SkipInstruction.SkipType.SKIP_COMPARE, "we can't reserialize because of unknown service")
        ));
    }

    @TestFactory
    @DisplayName("properties")
    Collection<DynamicNode> properties() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("properties.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("read-file")
    Collection<DynamicNode> read_file() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("read-file.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("read-file-segments")
    Collection<DynamicNode> read_file_segments() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("read-file.cap");
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        // TODO: the tests below are not really independend as the last one uses the complete byte array
        return List.of(
            DynamicTest.dynamicTest("No. 3 - Confirmed-REQ   atomicReadFile[195] file,1",
                () -> assertThat(pcapEvaluator.nextBVLC(3))
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUConfirmedRequest.class))
                    .extracting(APDUConfirmedRequest::getServiceRequest)
                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetConfirmedServiceRequestAtomicReadFile.class))
                    .satisfies(baCnetConfirmedServiceRequestAtomicReadFile -> {
                        assertThat(baCnetConfirmedServiceRequestAtomicReadFile)
                            .extracting(BACnetConfirmedServiceRequestAtomicReadFile::getFileIdentifier)
                            .extracting(BACnetApplicationTagObjectIdentifier::getObjectType, BACnetApplicationTagObjectIdentifier::getInstanceNumber)
                            .containsExactly(BACnetObjectType.FILE, 1L);
                        assertThat(baCnetConfirmedServiceRequestAtomicReadFile)
                            .extracting(BACnetConfirmedServiceRequestAtomicReadFile::getAccessMethod)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConfirmedServiceRequestAtomicReadFileStream.class))
                            .satisfies(baCnetConfirmedServiceRequestAtomicReadFileStream -> {
                                assertThat(baCnetConfirmedServiceRequestAtomicReadFileStream)
                                    .extracting(BACnetConfirmedServiceRequestAtomicReadFileStream::getFileStartPosition)
                                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagSignedInteger.class))
                                    .extracting(BACnetApplicationTagSignedInteger::getPayload)
                                    .extracting(BACnetTagPayloadSignedInteger::getActualValue)
                                    .isEqualTo(BigInteger.valueOf(0));
                                assertThat(baCnetConfirmedServiceRequestAtomicReadFileStream)
                                    .extracting(BACnetConfirmedServiceRequestAtomicReadFileStream::getRequestOctetCount)
                                    .asInstanceOf(InstanceOfAssertFactories.type(BACnetApplicationTagUnsignedInteger.class))
                                    .extracting(BACnetApplicationTagUnsignedInteger::getPayload)
                                    .extracting(BACnetTagPayloadUnsignedInteger::getActualValue)
                                    .isEqualTo(BigInteger.valueOf(2048L));
                            });
                    })),
            DynamicTest.dynamicTest("No. 4 - Complex-ACK     atomicReadFile[195]  (Message fragment 0)",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getSegment)
                    .satisfies(baos::writeBytes)),
            DynamicTest.dynamicTest("No. 5 - Complex-ACK",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUSegmentAck.class))
                    .extracting(APDUSegmentAck::getNegativeAck, APDUSegmentAck::getServer, APDUSegmentAck::getOriginalInvokeId, APDUSegmentAck::getSequenceNumber, APDUSegmentAck::getProposedWindowSize)
                    .containsExactly(false, false, (short) 195, (short) 0, (short) 16)),
            DynamicTest.dynamicTest("No. 6 - Complex-ACK     atomicReadFile[195]  (Message fragment 1)",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getSegment)
                    .satisfies(bytes -> LOGGER.info("Segment 1:\n{}", Hex.dump(bytes)))
                    .satisfies(baos::writeBytes)),
            DynamicTest.dynamicTest("No. 7 - Complex-ACK     atomicReadFile[195]  (Message fragment 2)",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getSegment)
                    .satisfies(bytes -> LOGGER.info("Segment 2:\n{}", Hex.dump(bytes)))
                    .satisfies(baos::writeBytes)),
            DynamicTest.dynamicTest("No. 8 - Complex-ACK     atomicReadFile[195]  (Message fragment 3)",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getSegment)
                    .satisfies(bytes -> LOGGER.info("Segment 3:\n{}", Hex.dump(bytes)))
                    .satisfies(baos::writeBytes)),
            DynamicTest.dynamicTest("No. 9 - Complex-ACK     atomicReadFile[195]  (Message Reassembled)",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUComplexAck.class))
                    .extracting(APDUComplexAck::getSegment)
                    .satisfies(bytes -> LOGGER.info("Segment 4:\n{}", Hex.dump(bytes)))
                    .satisfies(baos::writeBytes)),
            DynamicTest.dynamicTest("No. 10 - Segment-ACK",
                () -> assertThat(pcapEvaluator.nextBVLC())
                    .asInstanceOf(InstanceOfAssertFactories.type(BVLCOriginalUnicastNPDU.class))
                    .extracting(BVLCOriginalUnicastNPDU::getNpdu)
                    .extracting(NPDU::getApdu)
                    .asInstanceOf(InstanceOfAssertFactories.type(APDUSegmentAck.class))
                    .extracting(APDUSegmentAck::getNegativeAck, APDUSegmentAck::getServer, APDUSegmentAck::getOriginalInvokeId, APDUSegmentAck::getSequenceNumber, APDUSegmentAck::getProposedWindowSize)
                    .containsExactly(false, false, (short) 195, (short) 4, (short) 16)),
            DynamicTest.dynamicTest("Manually put together payload",
                () -> assertThat(baos.toByteArray())
                    .satisfies(bytes -> {
                        LOGGER.info("Trying to parse\n{}", Hex.dump(bytes));
                        BACnetServiceAckAtomicReadFile baCnetServiceAck = (BACnetServiceAckAtomicReadFile) BACnetServiceAckAtomicReadFile.staticParse(new ReadBufferByteBased(bytes), (long) bytes.length);
                        assertThat(baCnetServiceAck)
                            .isNotNull()
                            .extracting(BACnetServiceAckAtomicReadFile::getAccessMethod)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetServiceAckAtomicReadFileStream.class))
                            .extracting(BACnetServiceAckAtomicReadFileStream::getFileData)
                            .extracting(BACnetApplicationTagOctetString::getPayload)
                            .isNotNull();
                    })
            )
        );
    }


    @TestFactory
    @DisplayName("read-properties")
    Collection<DynamicNode> read_properties() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("read-properties.cap", BACNET_BPF_FILTER_UDP);
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Unconfirmed-REQ i-Am device,111",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(111, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(50, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.ACUITY_BRANDS_LIGHTING_INC, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 2 - Unconfirmed-REQ who-Is",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) apduUnconfirmedRequest.getServiceRequest();
                    assertNotNull(baCnetUnconfirmedServiceRequestWhoIs);
                }),
            DynamicTest.dynamicTest("No. 3 - Unconfirmed-REQ i-Am device,111",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(111, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(50, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.NO_SEGMENTATION, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.ACUITY_BRANDS_LIGHTING_INC, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 8 - Unconfirmed-REQ i-Am device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(201, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(1476, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.SEGMENTED_BOTH, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.ALERTON_HONEYWELL, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 8 - Unconfirmed-REQ i-Am device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(61, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.SEGMENTED_BOTH, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.ACUITY_BRANDS_LIGHTING_INC, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicTest.dynamicTest("No. 11 - Confirmed-REQ   readProperty[ 29] device,201 object-identifier",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 12 - Complex-ACK     readProperty[ 29] device,201 object-identifier device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 13 - Confirmed-REQ   readProperty[ 30] device,201 object-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 14 - Complex-ACK     readProperty[ 30] device,201 object-name device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString baCnetApplicationTagCharacterString = ((BACnetConstructedDataObjectName) baCnetServiceAckReadProperty.getValues()).getObjectName();
                    assertEquals("Lithonia Router", baCnetApplicationTagCharacterString.getValue());
                }),
            DynamicTest.dynamicTest("No. 15 - Confirmed-REQ   readProperty[ 31] device,201 object-type",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 16 - Complex-ACK     readProperty[ 31] device,201 object-type device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 17 - Confirmed-REQ   readProperty[ 32] device,201 system-status",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 18 - Complex-ACK     readProperty[ 32] device,201 system-status device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetDeviceStatusTagged value = ((BACnetConstructedDataSystemStatus) baCnetServiceAckReadProperty.getValues()).getSystemStatus();
                    assertEquals(BACnetDeviceStatus.OPERATIONAL, value.getValue());
                }),
            DynamicTest.dynamicTest("No. 19 - Confirmed-REQ   readProperty[ 33] device,201 vendor-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 20 - Complex-ACK     readProperty[ 33] device,201 vendor-name device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataVendorName) baCnetServiceAckReadProperty.getValues()).getVendorName();
                    assertEquals("Alerton Technologies, Inc.", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 21 - Confirmed-REQ   readProperty[ 34] device,201 vendor-identifier",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 22 - Complex-ACK     readProperty[ 34] device,201 vendor-identifier device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 23 - Confirmed-REQ   readProperty[ 35] device,201 model-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 24 - Complex-ACK     readProperty[ 35] device,201 model-name device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataModelName) baCnetServiceAckReadProperty.getValues()).getModelName();
                    assertEquals("LSi Controller", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 25 - Confirmed-REQ   readProperty[ 36] device,201 model-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 26 - Complex-ACK     readProperty[ 36] device,201 model-name device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataFirmwareRevision) baCnetServiceAckReadProperty.getValues()).getFirmwareRevision();
                    assertEquals("BACtalk LSi   v3.10 A         ", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 27 - Confirmed-REQ   readProperty[ 37] device,201 application-software-version",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 28 - Complex-ACK     readProperty[ 37] device,201 application-software-version device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataApplicationSoftwareVersion) baCnetServiceAckReadProperty.getValues()).getApplicationSoftwareVersion();
                    assertEquals("LSi Controller v3.11 E", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 29 - Confirmed-REQ   readProperty[ 38] device,201 protocol-version",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 30 - Complex-ACK     readProperty[ 38] device,201 protocol-version device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataProtocolVersion) baCnetServiceAckReadProperty.getValues()).getProtocolVersion();
                    assertEquals(1, value.getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 31 - Confirmed-REQ   readProperty[ 39] device,201 protocol-conformance-class",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 32 - Complex-ACK     readProperty[ 39] device,201 protocol-conformance-class device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(3, BACnetApplicationTagUnsignedInteger.getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 33 - Confirmed-REQ   readProperty[ 40] device,201 protocol-services-supported",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 34 - Complex-ACK     readProperty[ 40] device,201 protocol-services-supported device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetServicesSupportedTagged value = ((BACnetConstructedDataProtocolServicesSupported) baCnetServiceAckReadProperty.getValues()).getProtocolServicesSupported();
                    assertEquals(Arrays.asList(true, false, true, true, false, true, true, true, false, false, true, true, true, false, true, true, true, true, true, false, true, false, false, false, false, false, true, true, false, false, true, false, true, true, true), value.getPayload().getData());
                }),
            DynamicTest.dynamicTest("No. 35 - Confirmed-REQ   readProperty[ 41] device,201 protocol-object-types-supported",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 36 - Complex-ACK     readProperty[ 41] device,201 protocol-object-types-supported device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetObjectTypesSupportedTagged value = ((BACnetConstructedDataProtocolObjectTypesSupported) baCnetServiceAckReadProperty.getValues()).getProtocolObjectTypesSupported();
                    assertEquals(Arrays.asList(false, false, true, false, false, true, true, false, true, true, true, false, false, false, false, true, true, true), value.getPayload().getData());
                }),
            DynamicTest.dynamicTest("No. 37 - Confirmed-REQ   readProperty[ 42] device,201 max-apdu-length-accepted",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 38 - Complex-ACK     readProperty[ 42] device,201 max-apdu-length-accepted device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataMaxAPDULengthAccepted) baCnetServiceAckReadProperty.getValues()).getMaxApduLengthAccepted();
                    assertEquals(1476, value.getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 39 - Confirmed-REQ   readProperty[ 43] device,201 segmentation-supported",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 40 - Complex-ACK     readProperty[ 43] device,201 segmentation-supported device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetSegmentationTagged value = ((BACnetConstructedDataSegmentationSupported) baCnetServiceAckReadProperty.getValues()).getSegmentationSupported();
                    assertEquals(BACnetSegmentation.SEGMENTED_BOTH, value.getValue());
                }),
            DynamicTest.dynamicTest("No. 41 - Confirmed-REQ   readProperty[ 44] device,201 local-time",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 42 - Complex-ACK     readProperty[ 44] device,201 local-time device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagTime value = ((BACnetConstructedDataLocalTime) baCnetServiceAckReadProperty.getValues()).getLocalTime();
                    assertEquals(15, value.getPayload().getHour());
                    assertEquals(28, value.getPayload().getMinute());
                    assertEquals(41, value.getPayload().getSecond());
                }),
            DynamicTest.dynamicTest("No. 43 - Confirmed-REQ   readProperty[ 45] device,201 local-date",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 44 - Complex-ACK     readProperty[ 45] device,201 local-date device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagDate value = ((BACnetConstructedDataLocalDate) baCnetServiceAckReadProperty.getValues()).getLocalDate();
                    assertEquals(2005, value.getPayload().getYear());
                    assertEquals(9, value.getPayload().getMonth());
                    assertEquals(1, value.getPayload().getDayOfMonth());
                    assertEquals(4, value.getPayload().getDayOfWeek());
                }),
            DynamicTest.dynamicTest("No. 45 - Confirmed-REQ   readProperty[ 46] device,201 utc-offset",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 46 - Complex-ACK     readProperty[ 46] device,201 utc-offset device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagSignedInteger value = ((BACnetConstructedDataUTCOffset) baCnetServiceAckReadProperty.getValues()).getUtcOffset();
                    assertEquals(0, value.getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 47 - Confirmed-REQ   readProperty[ 47] device,201 daylights-savings-status",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 48 - Complex-ACK     readProperty[ 47] device,201 daylights-savings-status device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagBoolean value = ((BACnetConstructedDataDaylightSavingsStatus) baCnetServiceAckReadProperty.getValues()).getDaylightSavingsStatus();
                    assertFalse(value.getActualValue());
                }),
            DynamicTest.dynamicTest("No. 49 - Confirmed-REQ   readProperty[ 48] device,201 apdu-segment-timeout",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 50 - Complex-ACK     readProperty[ 48] device,201 apdu-segment-timeout device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataAPDUSegmentTimeout) baCnetServiceAckReadProperty.getValues()).getApduSegmentTimeout();
                    assertEquals(6000, value.getPayload().getValueUint16());
                }),
            DynamicTest.dynamicTest("No. 51 - Confirmed-REQ   readProperty[ 49] device,201 apdu-timeout",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 52 - Complex-ACK     readProperty[ 49] device,201 apdu-timeout device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataAPDUTimeout) baCnetServiceAckReadProperty.getValues()).getApduTimeout();
                    assertEquals(6000, value.getPayload().getValueUint16());
                }),
            DynamicTest.dynamicTest("No. 53 - Confirmed-REQ   readProperty[ 50] device,201 number-of-APDU-retries",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 54 - Complex-ACK     readProperty[ 50] device,201 number-of-APDU-retries device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataNumberOfAPDURetries) baCnetServiceAckReadProperty.getValues()).getNumberOfApduRetries();
                    assertEquals((short) 3, value.getPayload().getValueUint8());
                }),
            DynamicTest.dynamicTest("No. 55 - Confirmed-REQ   readProperty[ 51] device,201 time-synchronization-recipients",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.TIME_SYNCHRONIZATION_RECIPIENTS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 56 - ERROR           readProperty[ 51] device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorGeneral baCnetErrorWriteProperty = (BACnetErrorGeneral) apduError.getError();
                    assertEquals(ErrorClass.PROPERTY, baCnetErrorWriteProperty.getError().getErrorClass().getValue());
                    assertEquals(ErrorCode.UNKNOWN_PROPERTY, baCnetErrorWriteProperty.getError().getErrorCode().getValue());
                }),
            DynamicTest.dynamicTest("No. 57 - Confirmed-REQ   readProperty[ 52] device,201 max-master",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 58 - Complex-ACK     readProperty[ 52] device,201 max-master device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataDeviceMaxMaster) baCnetServiceAckReadProperty.getValues()).getMaxMaster();
                    assertEquals((short) 127, value.getPayload().getValueUint8());
                }),
            DynamicTest.dynamicTest("No. 59 - Confirmed-REQ   readProperty[ 53] device,201 max-info-frames",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 60 - Complex-ACK     readProperty[ 53] device,201 max-info-frames device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataDeviceMaxInfoFrames) baCnetServiceAckReadProperty.getValues()).getMaxInfoFrames();
                    assertEquals((short) 40, value.getPayload().getValueUint8());
                }),
            DynamicTest.dynamicTest("No. 61 - Confirmed-REQ   readProperty[ 54] device,201 device-address-binding",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 62 - Complex-ACK     readProperty[ 54] device,201 device-address-binding device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 63 - Confirmed-REQ   readProperty[ 55] device,201 (514) Vendor Proprietary Value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 64 - Complex-ACK     readProperty[ 55] device,201 (514) Vendor Proprietary Value device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagBoolean BACnetApplicationTagBoolean = (BACnetApplicationTagBoolean) value;
                    assertTrue(BACnetApplicationTagBoolean.getPayload().getIsFalse());
                }),
            DynamicTest.dynamicTest("No. 65 - Confirmed-REQ   readProperty[ 56] device,201 (515) Vendor Proprietary Value device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(201, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 66 - Complex-ACK     readProperty[ 56] device,201 Error",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorGeneral baCnetErrorWriteProperty = (BACnetErrorGeneral) apduError.getError();
                    assertEquals(ErrorClass.PROPERTY, baCnetErrorWriteProperty.getError().getErrorClass().getValue());
                    assertEquals(ErrorCode.UNKNOWN_PROPERTY, baCnetErrorWriteProperty.getError().getErrorCode().getValue());
                }),
            DynamicTest.dynamicTest("No. 67 - Confirmed-REQ   readProperty[ 57] device,61 object-identifier",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 68 - Complex-ACK     readProperty[ 57] object-identifier device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagObjectIdentifier BACnetApplicationTagObjectIdentifier = ((BACnetConstructedDataObjectIdentifier) baCnetServiceAckReadProperty.getValues()).getObjectIdentifier();
                    assertEquals(BACnetObjectType.DEVICE, BACnetApplicationTagObjectIdentifier.getObjectType());
                    assertEquals(61, BACnetApplicationTagObjectIdentifier.getInstanceNumber());
                }),
            DynamicTest.dynamicTest("No. 69 - Confirmed-REQ   readProperty[ 58] device,61 object-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 70 - Complex-ACK     readProperty[ 58] object-name device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString BACnetApplicationTagObjectIdentifier = ((BACnetConstructedDataObjectName) baCnetServiceAckReadProperty.getValues()).getObjectName();
                    assertEquals("SYNERGY", BACnetApplicationTagObjectIdentifier.getValue());
                }),
            DynamicTest.dynamicTest("No. 71 - Confirmed-REQ   readProperty[ 59] device,61 object-type",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 72 - Complex-ACK     readProperty[ 59] object-type device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetObjectTypeTagged value = ((BACnetConstructedDataObjectType) baCnetServiceAckReadProperty.getValues()).getObjectType();
                    assertEquals(BACnetObjectType.DEVICE, value.getValue());
                }),
            DynamicTest.dynamicTest("No. 73 - Confirmed-REQ   readProperty[ 60] device,61 system-status",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 74 - Complex-ACK     readProperty[ 60] device,61 system-status device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetDeviceStatusTagged value = ((BACnetConstructedDataSystemStatus) baCnetServiceAckReadProperty.getValues()).getSystemStatus();
                    assertEquals(BACnetDeviceStatus.OPERATIONAL, value.getValue());
                }),
            DynamicTest.dynamicTest("No. 75 - Confirmed-REQ   readProperty[ 61] device,61 vendor-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 76 - Complex-ACK     readProperty[ 61] device,61 vendor-name device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataVendorName) baCnetServiceAckReadProperty.getValues()).getVendorName();
                    assertEquals("Lithonia Lighting, Inc.", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 77 - Confirmed-REQ   readProperty[ 62] device,61 vendor-identifier",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 78 - Complex-ACK     readProperty[ 62] device,61 vendor-identifier device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 79 - Confirmed-REQ   readProperty[ 63] device,61 model-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 80 - Complex-ACK     readProperty[ 63] device,61 model-name device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataModelName) baCnetServiceAckReadProperty.getValues()).getModelName();
                    assertEquals("SYSC MLX", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 81 - Confirmed-REQ   readProperty[ 64] device,61 model-name",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 82 - Complex-ACK     readProperty[ 64] device,61 model-name device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataFirmwareRevision) baCnetServiceAckReadProperty.getValues()).getFirmwareRevision();
                    assertEquals("2x62i", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 83 - Confirmed-REQ   readProperty[ 65] device,61 application-software-version",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 84 - Complex-ACK     readProperty[ 65] device,61 application-software-version device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagCharacterString value = ((BACnetConstructedDataApplicationSoftwareVersion) baCnetServiceAckReadProperty.getValues()).getApplicationSoftwareVersion();
                    assertEquals("10-Nov-2004", value.getValue());
                }),
            DynamicTest.dynamicTest("No. 85 - Confirmed-REQ   readProperty[ 66] device,61 protocol-version",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 86 - Complex-ACK     readProperty[ 66] device,61 protocol-version device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataProtocolVersion) baCnetServiceAckReadProperty.getValues()).getProtocolVersion();
                    assertEquals(1, value.getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 87 - Confirmed-REQ   readProperty[ 67] device,61 protocol-conformance-class",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 88 - Complex-ACK     readProperty[ 67] device,61 protocol-conformance-class device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(2, BACnetApplicationTagUnsignedInteger.getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 89 - Confirmed-REQ   readProperty[ 68] device,61 protocol-services-supported",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 90 - Complex-ACK     readProperty[ 68] device,61 protocol-services-supported device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetServicesSupportedTagged value = ((BACnetConstructedDataProtocolServicesSupported) baCnetServiceAckReadProperty.getValues()).getProtocolServicesSupported();
                    assertEquals(Arrays.asList(false, false, false, false, false, false, true, true, false, false, false, false, true, false, true, true, false, true, false, false, true, false, false, false, false, false, true, true, false, false, false, false, true, true, true), value.getPayload().getData());
                }),
            DynamicTest.dynamicTest("No. 91 - Confirmed-REQ   readProperty[ 69] device,61 protocol-object-types-supported",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 92 - Complex-ACK     readProperty[ 69] device,61 protocol-object-types-supported device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetObjectTypesSupportedTagged value = ((BACnetConstructedDataProtocolObjectTypesSupported) baCnetServiceAckReadProperty.getValues()).getProtocolObjectTypesSupported();
                    assertEquals(Arrays.asList(true, true, true, true, true, true, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false), value.getPayload().getData());
                }),
            DynamicTest.dynamicTest("No. 93 - Confirmed-REQ   readProperty[ 70] device,61 max-apdu-length-accepted",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 94 - Complex-ACK     readProperty[ 70] device,61 max-apdu-length-accepted device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataMaxAPDULengthAccepted) baCnetServiceAckReadProperty.getValues()).getMaxApduLengthAccepted();
                    assertEquals(480, value.getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 95 - Confirmed-REQ   readProperty[ 71] device,61 segmentation-supported",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 96 - Complex-ACK     readProperty[ 71] device,61 segmentation-supported device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetSegmentationTagged value = ((BACnetConstructedDataSegmentationSupported) baCnetServiceAckReadProperty.getValues()).getSegmentationSupported();
                    assertEquals(BACnetSegmentation.SEGMENTED_BOTH, value.getValue());
                }),
            DynamicTest.dynamicTest("No. 97 - Confirmed-REQ   readProperty[ 72] device,61 local-time",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 98 - Complex-ACK     readProperty[ 72] device,61 local-time device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagTime value = ((BACnetConstructedDataLocalTime) baCnetServiceAckReadProperty.getValues()).getLocalTime();
                    assertEquals(15, value.getPayload().getHour());
                    assertEquals(26, value.getPayload().getMinute());
                    assertEquals(40, value.getPayload().getSecond());
                }),
            DynamicTest.dynamicTest("No. 99 - Confirmed-REQ   readProperty[ 73] device,61 local-date",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 100 - Complex-ACK     readProperty[ 73] device,61 local-date device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagDate value = ((BACnetConstructedDataLocalDate) baCnetServiceAckReadProperty.getValues()).getLocalDate();
                    assertEquals(2005, value.getPayload().getYear());
                    assertEquals(9, value.getPayload().getMonth());
                    assertEquals(1, value.getPayload().getDayOfMonth());
                    assertEquals(255, value.getPayload().getDayOfWeek());
                }),
            DynamicTest.dynamicTest("No. 101 - Confirmed-REQ   readProperty[ 74] device,61 utc-offset",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 102 - Complex-ACK     readProperty[ 74] device,61 utc-offset device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagSignedInteger value = ((BACnetConstructedDataUTCOffset) baCnetServiceAckReadProperty.getValues()).getUtcOffset();
                    assertEquals(BigInteger.valueOf(-300), value.getPayload().getActualValue());
                }),
            DynamicTest.dynamicTest("No. 103 - Confirmed-REQ   readProperty[ 75] device,61 daylights-savings-status",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 104 - Complex-ACK     readProperty[ 75] device,61 daylights-savings-status device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagBoolean value = ((BACnetConstructedDataDaylightSavingsStatus) baCnetServiceAckReadProperty.getValues()).getDaylightSavingsStatus();
                    assertTrue(value.getPayload().getIsTrue());
                }),
            DynamicTest.dynamicTest("No. 105 - Confirmed-REQ   readProperty[ 76] device,61 apdu-segment-timeout",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 106 - Complex-ACK     readProperty[ 76] device,61 apdu-segment-timeout device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataAPDUSegmentTimeout) baCnetServiceAckReadProperty.getValues()).getApduSegmentTimeout();
                    assertEquals(8000, value.getPayload().getValueUint16());
                }),
            DynamicTest.dynamicTest("No. 107 - Confirmed-REQ   readProperty[ 77] device,61 apdu-timeout",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 108 - Complex-ACK     readProperty[ 77] device,61 apdu-timeout device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataAPDUTimeout) baCnetServiceAckReadProperty.getValues()).getApduTimeout();
                    assertEquals(8000, value.getPayload().getValueUint16());
                }),
            DynamicTest.dynamicTest("No. 109 - Confirmed-REQ   readProperty[ 78] device,61 number-of-APDU-retries",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 100 - Complex-ACK     readProperty[ 78] device,61 number-of-APDU-retries device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataNumberOfAPDURetries) baCnetServiceAckReadProperty.getValues()).getNumberOfApduRetries();
                    assertEquals((short) 3, value.getPayload().getValueUint8());
                }),
            DynamicTest.dynamicTest("No. 111 - Confirmed-REQ   readProperty[ 79] device,61 time-synchronization-recipients",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.TIME_SYNCHRONIZATION_RECIPIENTS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 112 - ERROR           readProperty[ 79] device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorGeneral baCnetErrorWriteProperty = (BACnetErrorGeneral) apduError.getError();
                    assertEquals(ErrorClass.PROPERTY, baCnetErrorWriteProperty.getError().getErrorClass().getValue());
                    assertEquals(ErrorCode.UNKNOWN_PROPERTY, baCnetErrorWriteProperty.getError().getErrorCode().getValue());
                }),
            DynamicTest.dynamicTest("No. 113 - Confirmed-REQ   readProperty[ 80] device,61 max-master",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 114 - Complex-ACK     readProperty[ 80] device,61 max-master device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataDeviceMaxMaster) baCnetServiceAckReadProperty.getValues()).getMaxMaster();
                    assertEquals((short) 127, value.getPayload().getValueUint8());
                }),
            DynamicTest.dynamicTest("No. 115 - Confirmed-REQ   readProperty[ 81] device,61 max-info-frames",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 116 - Complex-ACK     readProperty[ 81] device,61 max-info-frames device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTagUnsignedInteger value = ((BACnetConstructedDataDeviceMaxInfoFrames) baCnetServiceAckReadProperty.getValues()).getMaxInfoFrames();
                    assertEquals((short) 1, value.getPayload().getValueUint8());
                }),
            DynamicTest.dynamicTest("No. 117 - Confirmed-REQ   readProperty[ 82] device,61 device-address-binding",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 118 - Complex-ACK     readProperty[ 82] device,61 device-address-binding device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 119 - Confirmed-REQ   readProperty[ 83] device,61 (514) Vendor Proprietary Value",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 120 - Complex-ACK     readProperty[ 83] device,61 (514) Vendor Proprietary Value device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagTime baCnetApplicationTagTime = (BACnetApplicationTagTime) value;
                    assertEquals(7, baCnetApplicationTagTime.getPayload().getHour());
                    assertEquals(11, baCnetApplicationTagTime.getPayload().getMinute());
                    assertEquals(38, baCnetApplicationTagTime.getPayload().getSecond());
                }),
            DynamicTest.dynamicTest("No. 121 - Confirmed-REQ   readProperty[ 84] device,61 (515) Vendor Proprietary Value device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequest serviceRequest = apduConfirmedRequest.getServiceRequest();
                    assertNotNull(serviceRequest);
                    BACnetConfirmedServiceRequestReadProperty baCnetConfirmedServiceRequestReadProperty = (BACnetConfirmedServiceRequestReadProperty) serviceRequest;
                    assertEquals(BACnetObjectType.DEVICE, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetConfirmedServiceRequestReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getValue());
                }),
            DynamicTest.dynamicTest("No. 122 - Complex-ACK     readProperty[ 84] device,61 Error",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAck baCnetServiceAck = apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAck);
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) baCnetServiceAck;
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(61, baCnetServiceAckReadProperty.getObjectIdentifier().getInstanceNumber());
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getValue());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagTime baCnetApplicationTagTime = (BACnetApplicationTagTime) value;
                    assertEquals(20, baCnetApplicationTagTime.getPayload().getHour());
                    assertEquals(3, baCnetApplicationTagTime.getPayload().getMinute());
                    assertEquals(18, baCnetApplicationTagTime.getPayload().getSecond());
                })
        );
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("read-property-bad")
    Collection<DynamicNode> read_property_bad() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("read-property-bad.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("read-property-epics")
    Collection<DynamicNode> read_property_epics() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("read-property-epics.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("read-property-synergy")
    Collection<DynamicNode> read_property_synergy() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("read-property-synergy.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("read-property")
    Collection<DynamicNode> read_property() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("read-property.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("readfile")
    Collection<DynamicNode> readfile() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("readfile.cap");
        return List.of(pcapEvaluator.parseEmAll(
            // 3,
            // 4,
            // 5
        ));
    }

    @TestFactory
    @DisplayName("readrange_malformed")
    Collection<DynamicNode> readrange_malformed() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("readrange_malformed.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("reinit-device")
    Collection<DynamicNode> reinit_device() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("reinit-device.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(201, SkipInstruction.SkipType.SKIP_COMPLETE, "is using proprietary value which should not be allowed"),
            skip(207, SkipInstruction.SkipType.SKIP_COMPLETE, "is using proprietary value which should not be allowed"),
            skip(223, SkipInstruction.SkipType.SKIP_COMPLETE, "is using proprietary value which should not be allowed")
        ));
    }

    @TestFactory
    @DisplayName("router")
    Collection<DynamicNode> router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("router.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("routers")
    Collection<DynamicNode> routers() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("routers.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("rp-device")
    Collection<DynamicNode> rp_device() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("rp-device.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("rp-shed-level")
    Collection<DynamicNode> rp_shed_level() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("rp-shed-level.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("rp")
    Collection<DynamicNode> rp() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("rp.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("rpm-error")
    Collection<DynamicNode> rpm_error() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("rpm-error.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("rpm")
    Collection<DynamicNode> rpm() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("rpm.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("rpm_multiple_scheduler_bug")
    Collection<DynamicNode> rpm_multiple_scheduler_bug() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("rpm_multiple_scheduler_bug.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("schedule-rpm-foreign")
    Collection<DynamicNode> schedule_rpm_foreign() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("schedule-rpm-foreign.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("signed_value_negative")
    Collection<DynamicNode> signed_value_negative() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("signed_value_negative.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("single-RPM")
    Collection<DynamicNode> single_RPM() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("single-RPM.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("too many packages 293726... needs splitup")
    @TestFactory
    @DisplayName("softdel-BTL")
    Collection<DynamicNode> softdel_BTL() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("softdel-BTL.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("special-events")
    Collection<DynamicNode> special_events() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("special-events.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(13, "looks like broken package")
        ));
    }

    @TestFactory
    @DisplayName("startup-exchange")
    Collection<DynamicNode> startup_exchange() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("startup-exchange.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("state_text")
    Collection<DynamicNode> state_text() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("state_text.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("state_text_good")
    Collection<DynamicNode> state_text_good() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("state_text_good.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("subordinatelist_rpm")
    Collection<DynamicNode> subordinatelist_rpm() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("subordinatelist_rpm.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("synergy-binding-2x63y")
    Collection<DynamicNode> synergy_binding_2x63y() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("synergy-binding-2x63y.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("synergy-broken-rpm")
    Collection<DynamicNode> synergy_broken_rpm() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("synergy-broken-rpm.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("synergy-device")
    Collection<DynamicNode> synergy_device() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("synergy-device.cap");
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @TestFactory
    @DisplayName("time-sync")
    Collection<DynamicNode> time_sync() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("time-sync.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("tridium jace2")
    Collection<DynamicNode> tridium_jace2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("tridium%20jace2.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
        ));
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("u+4_MSTP")
    Collection<DynamicNode> u_4_MSTP() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("u+4_MSTP.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("weekend")
    Collection<DynamicNode> weekend() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("weekend.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("who-has-I-have")
    Collection<DynamicNode> who_has_I_have() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("who-has-I-have.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(1, SkipInstruction.SkipType.SKIP_COMPARE, "who is has some broken bytes")
        ));
    }

    @TestFactory
    @DisplayName("who-has")
    Collection<DynamicNode> who_has() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("who-has.cap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Unconfirmed-REQ who-Is 133 133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 2 - Unconfirmed-REQ who-Has device,133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestWhoHas baCnetUnconfirmedServiceRequestWhoHas = (BACnetUnconfirmedServiceRequestWhoHas) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHas.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHas.getDeviceInstanceRangeLowLimit().getPayload().getActualValue().longValue());
                    BACnetUnconfirmedServiceRequestWhoHasObjectIdentifier baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier = (BACnetUnconfirmedServiceRequestWhoHasObjectIdentifier) baCnetUnconfirmedServiceRequestWhoHas.getObject();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier.getObjectIdentifier().getInstanceNumber());
                }),
            DynamicTest.dynamicTest("No. 3-4 - skip 2 LLC packages",
                () -> pcapEvaluator.skipPackages(2)),
            DynamicTest.dynamicTest("No. 5 - Unconfirmed-REQ I-Am 133 133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalBroadcastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequestIAm baCnetUnconfirmedServiceRequestIAm = (BACnetUnconfirmedServiceRequestIAm) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestIAm.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(480, baCnetUnconfirmedServiceRequestIAm.getMaximumApduLengthAcceptedLength().getPayload().getActualValue().longValue());
                    assertEquals(BACnetSegmentation.SEGMENTED_BOTH, baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getValue());
                    assertEquals(BACnetVendorId.ACUITY_BRANDS_LIGHTING_INC, baCnetUnconfirmedServiceRequestIAm.getVendorId().getValue());
                }),
            DynamicContainer.dynamicContainer("Confirmed-REQ atomicWriteFile 1-30", () -> {
                Collection<DynamicNode> nodes = new LinkedList<>();
                IntStream.range(1, 31).forEach(i -> {
                    nodes.add(DynamicTest.dynamicTest("No. " + (4 + i * 2) + " - Confirmed-REQ atomicWriteFile [" + i + "] file,0", () -> {
                        BVLC bvlc = pcapEvaluator.nextBVLC();
                        dump(bvlc);
                        NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                        APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) npdu.getApdu();
                        BACnetConfirmedServiceRequestAtomicWriteFile baCnetConfirmedServiceRequestAtomicWriteFile = (BACnetConfirmedServiceRequestAtomicWriteFile) apduConfirmedRequest.getServiceRequest();
                        assertEquals(BACnetObjectType.FILE, baCnetConfirmedServiceRequestAtomicWriteFile.getDeviceIdentifier().getObjectType());
                        assertNotNull(baCnetConfirmedServiceRequestAtomicWriteFile.getFileStartPosition());
                        assertNotNull(baCnetConfirmedServiceRequestAtomicWriteFile.getFileData());
                    }));
                    nodes.add(DynamicTest.dynamicTest("No. " + (5 + i * 2) + " - Confirmed-Ack     atomicWriteFile [" + i + "]", () -> {
                        BVLC bvlc = pcapEvaluator.nextBVLC();
                        dump(bvlc);
                        NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                        APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                        BACnetServiceAckAtomicWriteFile baCnetServiceAckAtomicWriteFile = (BACnetServiceAckAtomicWriteFile) apduComplexAck.getServiceAck();
                        assertNotNull(baCnetServiceAckAtomicWriteFile.getFileStartPosition());
                    }));
                });
                return nodes.iterator();
            }),
            DynamicTest.dynamicTest("No. 66 - Unconfirmed-REQ who-Has device,133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestWhoHas baCnetUnconfirmedServiceRequestWhoHas = (BACnetUnconfirmedServiceRequestWhoHas) apduUnconfirmedRequest.getServiceRequest();
                    assertNull(baCnetUnconfirmedServiceRequestWhoHas.getDeviceInstanceRangeLowLimit());
                    assertNull(baCnetUnconfirmedServiceRequestWhoHas.getDeviceInstanceRangeLowLimit());
                    BACnetUnconfirmedServiceRequestWhoHasObjectIdentifier baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier = (BACnetUnconfirmedServiceRequestWhoHasObjectIdentifier) baCnetUnconfirmedServiceRequestWhoHas.getObject();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier.getObjectIdentifier().getInstanceNumber());
                }),
            DynamicTest.dynamicTest("No. 67 - skip 1 LLC packages",
                () -> pcapEvaluator.skipPackages(1)),
            DynamicTest.dynamicTest("No. 68 - Unconfirmed-REQ i-Have device,4194303 device,133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIHave baCnetUnconfirmedServiceRequestIHave = (BACnetUnconfirmedServiceRequestIHave) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getObjectType());
                    assertEquals(4194303, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getInstanceNumber());
                    assertEquals("Unknown", baCnetUnconfirmedServiceRequestIHave.getObjectName().getValue());
                }),
            DynamicTest.dynamicTest("No. 69 - Unconfirmed-REQ i-Have device,133 device,133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequestIHave baCnetUnconfirmedServiceRequestIHave = (BACnetUnconfirmedServiceRequestIHave) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getInstanceNumber());
                    assertEquals("SYNERGY", baCnetUnconfirmedServiceRequestIHave.getObjectName().getValue());
                }),
            DynamicTest.dynamicTest("No. 70 - Unconfirmed-REQ who-Has device,133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestWhoHas baCnetUnconfirmedServiceRequestWhoHas = (BACnetUnconfirmedServiceRequestWhoHas) apduUnconfirmedRequest.getServiceRequest();
                    assertNull(baCnetUnconfirmedServiceRequestWhoHas.getDeviceInstanceRangeLowLimit());
                    assertNull(baCnetUnconfirmedServiceRequestWhoHas.getDeviceInstanceRangeLowLimit());
                    BACnetUnconfirmedServiceRequestWhoHasObjectIdentifier baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier = (BACnetUnconfirmedServiceRequestWhoHasObjectIdentifier) baCnetUnconfirmedServiceRequestWhoHas.getObject();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHasObjectIdentifier.getObjectIdentifier().getInstanceNumber());
                }),
            DynamicTest.dynamicTest("No. 71 - skip 1 LLC packages",
                () -> pcapEvaluator.skipPackages(1)),
            DynamicTest.dynamicTest("No. 72 - Unconfirmed-REQ i-Have device,4194303 device,133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestIHave baCnetUnconfirmedServiceRequestIHave = (BACnetUnconfirmedServiceRequestIHave) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getObjectType());
                    assertEquals(4194303, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getInstanceNumber());
                    assertEquals("Unknown", baCnetUnconfirmedServiceRequestIHave.getObjectName().getValue());
                }),
            DynamicTest.dynamicTest("No. 73 - Unconfirmed-REQ i-Have device,133 device,133",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) npdu.getApdu();
                    BACnetUnconfirmedServiceRequestIHave baCnetUnconfirmedServiceRequestIHave = (BACnetUnconfirmedServiceRequestIHave) apduUnconfirmedRequest.getServiceRequest();
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestIHave.getDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestIHave.getObjectIdentifier().getInstanceNumber());
                    assertEquals("SYNERGY", baCnetUnconfirmedServiceRequestIHave.getObjectName().getValue());
                })
        );
    }

    @TestFactory
    @DisplayName("who-is-i-am")
    Collection<DynamicNode> who_is_i_am() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("who-is-i-am.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("whois-basrtp-b-1")
    Collection<DynamicNode> whois_basrtp_b_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("whois-basrtp-b-1.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @Disabled("no udp packages")
    @TestFactory
    @DisplayName("whois-basrtp-b-2")
    Collection<DynamicNode> whois_basrtp_b_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("whois-basrtp-b-2.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("whois-iam")
    Collection<DynamicNode> whois_iam() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("whois-iam.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("wireshark_BBMDError")
    Collection<DynamicNode> wireshark_BBMDError() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("wireshark_BBMDError.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("wireshark_CEN_9_11")
    Collection<DynamicNode> wireshark_CEN_9_11() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("wireshark_CEN_9_11.pcap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("wp-rp-index")
    Collection<DynamicNode> wp_rp_index() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("wp-rp-index.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(34, "Null doesn't seem to be a right value for object-list"),
            skip(42, "Null doesn't seem to be a right value for object-list"),
            skip(50, "Null doesn't seem to be a right value for object-list"),
            skip(62, "Null doesn't seem to be a right value for object-list"),
            skip(70, "Null doesn't seem to be a right value for object-list"),
            skip(78, "Null doesn't seem to be a right value for object-list"),
            skip(86, "Null doesn't seem to be a right value for object-list"),
            skip(94, "Null doesn't seem to be a right value for object-list")
        ));
    }

    @TestFactory
    @DisplayName("wp-weekly-schedule-index")
    Collection<DynamicNode> wp_weekly_schedule_index() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("wp-weekly-schedule-index.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("wp-weekly-schedule-test")
    Collection<DynamicNode> wp_weekly_schedule_test() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("wp-weekly-schedule-test.cap");
        return List.of(pcapEvaluator.parseEmAll(
            // TODO: fixme
            skip(17, SkipInstruction.SkipType.SKIP_COMPARE, "something is broken with the pcap reader as a ff is randomly missing")
        ));
    }

    @TestFactory
    @DisplayName("wp_weekly_schedule")
    Collection<DynamicNode> wp_weekly_schedule() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("wp_weekly_schedule.cap");
        return List.of(pcapEvaluator.parseEmAll(
            skip(1, "Incomplete/wrong request"),
            // TODO: there should be 7 values but apparently only 2 are transmitted.. could be right could be wrong
            skip(13, "Incomplete/wrong request"),
            skip(19, "Seems like wrong data")
        ));
    }

    @TestFactory
    @DisplayName("write-property-array")
    Collection<DynamicNode> write_property_array() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("write-property-array.cap");
        return List.of(pcapEvaluator.parseEmAll(
            skip(39, "strange priority array with two extra nulls at the end? apparently the device was confused too"),
            skip(44, "strange priority array with two extra nulls at the end? apparently the device was confused too")
        ));
    }

    @TestFactory
    @DisplayName("write-property-multiple")
    Collection<DynamicNode> write_property_multiple() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("write-property-multiple.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("write-property-wattstopper-panel")
    Collection<DynamicNode> write_property_wattstopper_panel() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("write-property-wattstopper-panel.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("write-property")
    Collection<DynamicNode> write_property() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("write-property.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll(
            skip(1594, "Malformed Package"),
            skip(1595, "Malformed Package"),
            skip(1596, "Malformed Package"),
            skip(1597, "Malformed Package"),
            skip(1598, "Malformed Package"),
            skip(1599, "Malformed Package"),
            skip(1600, "Malformed Package"),
            skip(1601, "Malformed Package"),
            skip(1602, "Malformed Package"),
            skip(1603, "Malformed Package"),
            skip(1604, "Malformed Package"),
            skip(1605, "Malformed Package"),
            skip(1606, "Malformed Package"),
            skip(1607, "Malformed Package"),
            skip(1608, "Malformed Package"),
            skip(1609, "Malformed Package")
        ));
    }

    @TestFactory
    @DisplayName("write-property2")
    Collection<DynamicNode> write_property2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("write-property2.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    private static void dump(Serializable serializable) throws SerializationException {
        if (!LOGGER.isInfoEnabled()) {
            return;
        }
        WriteBufferBoxBased writeBuffer = new WriteBufferBoxBased(true, true);
        serializable.serialize(writeBuffer);
        LOGGER.info("{}\n{}", serializable.getClass().getName(), writeBuffer.getBox());
    }

    private TestPcapEvaluator pcapEvaluator(String pcapFile) throws IOException, PcapNativeException, NotOpenException {
        TestPcapEvaluator pcapEvaluator = new TestPcapEvaluator(pcapFile);
        toBeClosed.offer(pcapEvaluator);
        return pcapEvaluator;
    }

    private TestPcapEvaluator pcapEvaluator(String pcapFile, String filter) throws IOException, PcapNativeException, NotOpenException {
        TestPcapEvaluator pcapEvaluator = new TestPcapEvaluator(pcapFile, filter);
        toBeClosed.offer(pcapEvaluator);
        return pcapEvaluator;
    }

    private static class TestPcapEvaluator extends PCAPEvaluator {

        public TestPcapEvaluator(String pcapFile) throws IOException, PcapNativeException, NotOpenException {
            super(pcapFile);
        }

        public TestPcapEvaluator(String pcapFile, String filter) throws IOException, PcapNativeException, NotOpenException {
            super(pcapFile, filter);
        }

        public DynamicContainer parseEmAll(SkipInstruction... skipInstructions) {
            return parseRange(1, maxPackages, skipInstructions);
        }

        public DynamicContainer parseFrom(int startPackageNumber, SkipInstruction... skipInstructions) {
            return parseRange(startPackageNumber, maxPackages, skipInstructions);
        }

        public DynamicContainer parseTill(int packageNumber, SkipInstruction... skipInstructions) {
            return parseRange(1, packageNumber, skipInstructions);
        }

        public DynamicContainer parseRange(int startPackageNumber, int endPackageNumber, SkipInstruction... skipInstructions) {
            Map<Integer, SkipInstruction> skipInstructionMap = Arrays.stream(skipInstructions).collect(Collectors.toMap(SkipInstruction::getPackageNumber, v -> v, (s1, s2) -> s1.skipType.ordinal() < s2.skipType.ordinal() ? s1 : s2));
            Set<Integer> numbersToSkip = Arrays.stream(skipInstructions).filter(SkipInstruction::shouldSkipAll).map(SkipInstruction::getPackageNumber).collect(Collectors.toSet());
            // This means we requested to skip no test number
            boolean hasNoSkipping = numbersToSkip.size() <= 0;
            boolean hasSkipping = !hasNoSkipping;
            // This function maps the test number to the package number if there is an offset. That happens when we apply a BPF filter
            Function<Integer, Integer> packageNumResolver = i -> packageNumbers != null ? packageNumbers.get(i - 1) : i;
            return DynamicContainer.dynamicContainer(
                "Parse em all (No. " + startPackageNumber + "-" + endPackageNumber + ")" + (hasNoSkipping ? "" : " [skipped unit-test-ids" + numbersToSkip + "]"),
                () -> IntStream.range(startPackageNumber, endPackageNumber + 1)
                    .mapToObj(
                        (i) -> DynamicTest.dynamicTest(
                            "No. " + packageNumResolver.apply(i) + " - Unit-Test nr." + i,
                            () -> {
                                Integer packageNumber = packageNumResolver.apply(i);
                                // TODO: maybe we can migrate that to the instruction below
                                if (hasSkipping && numbersToSkip.contains(packageNumber)) {
                                    LOGGER.info("Skipping unfiltered package {} with test nr. {}", packageNumber, i);
                                    skipPackages(1);
                                    throw new TestAbortedException("Package nr. " + packageNumber + " filtered (Unit-Test nr. " + i + ")");
                                }
                                SkipInstruction skipInstruction = skipInstructionMap.getOrDefault(packageNumber, noskip());
                                if (!skipInstruction.shouldParse()) {
                                    throw new TestAbortedException(skipInstruction.toString());
                                }
                                if (packageNumber <= currentPackageNumber) {
                                    throw new TestAbortedException("Apparently we read that package before");
                                }
                                BVLC bvlc = nextBVLC(packageNumber, skipInstruction);
                                LOGGER.info("Test number {} is package number {}", i, currentPackageNumber);
                                assumeTrue(bvlc != null, "No more package left");
                                dump(bvlc);
                            }
                        )
                    )
                    .map(DynamicNode.class::cast)
                    .iterator()
            );
        }
    }

    private static class PCAPEvaluator implements Closeable {
        protected int readPackages = 0;
        protected int currentPackageNumber = 0;
        protected boolean done = false;
        protected final String pcapFile;
        protected final PcapHandle pcapHandle;
        // maps timestamp to package number
        protected final Map<Timestamp, Integer> timestampToPackageNumberMap;
        // maps read package (index) to package number
        protected final List<Integer> packageNumbers;
        protected final int maxPackages;

        public PCAPEvaluator(String pcapFile) throws IOException, PcapNativeException, NotOpenException {
            this(pcapFile, null);
        }

        public PCAPEvaluator(String pcapFile, String filter) throws IOException, PcapNativeException, NotOpenException {
            this.pcapFile = pcapFile;
            String toParse = DownloadAndCache(pcapFile);
            LOGGER.info("Reading {}", toParse);
            PcapHandle intermediateHandle = getHandle(toParse);
            int packageNumber = 0;
            if (filter != null) {
                // In case of filtering we need to read all packages
                LOGGER.info("Building timestamp number map");
                timestampToPackageNumberMap = new HashMap<>();
                while (intermediateHandle.getNextPacket() != null) {
                    timestampToPackageNumberMap.put(intermediateHandle.getTimestamp(), ++packageNumber);
                }
                intermediateHandle.close();
                // Count package numbers now
                intermediateHandle = getHandle(toParse);
                intermediateHandle.setFilter(filter, BpfProgram.BpfCompileMode.OPTIMIZE);
                packageNumber = 0;
                packageNumbers = new LinkedList<>();
                while (intermediateHandle.getNextPacket() != null) {
                    packageNumber++;
                    packageNumbers.add(timestampToPackageNumberMap.get(intermediateHandle.getTimestamp()));
                }
                intermediateHandle.close();
                // We need a new handle as we consumed the old
                intermediateHandle = getHandle(toParse);
                intermediateHandle.setFilter(filter, BpfProgram.BpfCompileMode.OPTIMIZE);
            } else {
                timestampToPackageNumberMap = null;
                packageNumbers = null;
                while (intermediateHandle.getNextPacket() != null) {
                    packageNumber++;
                }
                intermediateHandle.close();
                intermediateHandle = getHandle(toParse);
            }
            pcapHandle = intermediateHandle;
            maxPackages = packageNumber;
        }

        public void skipTo(int targetPackageNumber) throws NotOpenException {
            if (targetPackageNumber <= currentPackageNumber) {
                throw new IllegalArgumentException("Package number must be bigger than " + currentPackageNumber);
            }
            LOGGER.info("Skipping to package number {} starting at package number {}. Current packages so far {}", targetPackageNumber, currentPackageNumber, readPackages);
            do {
                nextPacket();
            } while (currentPackageNumber < targetPackageNumber);
            LOGGER.info("Ended skipping at {} read packages with package number {}", readPackages, currentPackageNumber);
        }

        public void skipPackages(int numberOfReadPackages) {
            LOGGER.info("Skipping {} package reads", numberOfReadPackages);
            IntStream.rangeClosed(1, numberOfReadPackages).forEach(i -> {
                try {
                    nextPacket();
                    LOGGER.info("Skipping {}th package with package number {}", i, currentPackageNumber);
                } catch (NotOpenException e) {
                    throw new RuntimeException(e);
                }
            });
        }

        public BVLC nextBVLC() throws NotOpenException, ParseException, SerializationException {
            return nextBVLC(noskip());
        }

        public BVLC nextBVLC(SkipInstruction skipInstruction) throws NotOpenException, ParseException, SerializationException {
            return nextBVLC(null, skipInstruction);
        }

        public BVLC nextBVLC(Integer ensurePackageNumber) throws NotOpenException, ParseException, SerializationException {
            return nextBVLC(ensurePackageNumber, noskip());
        }

        public BVLC nextBVLC(Integer ensurePackageNumber, SkipInstruction skipInstruction) throws NotOpenException, ParseException, SerializationException {
            Packet packet = nextPacket();
            if (packet == null) {
                return null;
            }
            if (ensurePackageNumber != null) {
                if (ensurePackageNumber < 0) {
                    throw new IllegalArgumentException("Searched package number must be > 0");
                }
                // we ensure we find the searched package
                while (currentPackageNumber < ensurePackageNumber) {
                    packet = nextPacket();
                    if (packet == null) {
                        throw new IllegalArgumentException("Could not find package with package number " + ensurePackageNumber);
                    }
                }
                if (currentPackageNumber > ensurePackageNumber) {
                    throw new IllegalArgumentException("Could not find package with package number " + ensurePackageNumber);
                }
            }
            return getBvlc(packet, skipInstruction);
        }

        private BVLC getBvlc(Packet packet, SkipInstruction skipInstruction) throws ParseException, SerializationException {
            UdpPacket udpPacket = packet.get(UdpPacket.class);
            assumeTrue(udpPacket != null, "nextBVLC assumes a UDP Packet. If non is there it might by LLC");
            LOGGER.info("Handling UDP\n{}", udpPacket);
            byte[] rawData = udpPacket.getPayload().getRawData();
            LOGGER.info("Reading BVLC from:\n{}", Hex.dump(rawData));
            try {
                BVLC bvlc = BVLC.staticParse(new ReadBufferByteBased(rawData));
                if (skipInstruction.shouldSerialize()) {
                    WriteBufferByteBased writeBuffer = new WriteBufferByteBased(bvlc.getLengthInBytes());
                    bvlc.serialize(writeBuffer);
                    if (skipInstruction.shouldCompare()) {
                        @SuppressWarnings("redundant")
                        byte[] expectedBytes = rawData;
                        byte[] actualBytes = writeBuffer.getBytes();
                        if (!Arrays.equals(expectedBytes, actualBytes)) {
                            // This goes to std out on purpose to preserve coloring
                            System.out.println(HexDiff.diffHex(expectedBytes, actualBytes));
                        }
                        assertThat(actualBytes)
                            .withRepresentation(HexadecimalRepresentation.HEXA_REPRESENTATION)
                            .describedAs("re-serialized output doesn't match original bytes:%s\n", bvlc)
                            .isEqualTo(expectedBytes);
                    }
                } else {
                    LOGGER.debug("{}", skipInstruction);
                }
                return bvlc;
            } catch (ParseException e) {
                throw new ParseException(String.format("Caught at current package number: %d. Packages read so far %d", currentPackageNumber, readPackages), e);
            } catch (SerializationException e) {
                throw new SerializationException(String.format("Caught at current package number: %d. Packages read so far %d", currentPackageNumber, readPackages), e);
            }
        }

        private Packet nextPacket() throws NotOpenException {
            Packet packet = pcapHandle.getNextPacket();
            if (packet == null) {
                done = true;
                LOGGER.warn("No more package. Suggestion: You can read up to {} packages.", readPackages);
                return null;
            }
            readPackages++;
            if (timestampToPackageNumberMap == null) {
                currentPackageNumber++;
            } else {
                currentPackageNumber = timestampToPackageNumberMap.get(pcapHandle.getTimestamp());
            }
            LOGGER.info("({}) Next packet:\n{}", currentPackageNumber, packet);
            return packet;
        }

        private PcapHandle getHandle(String file) throws PcapNativeException {
            return Pcaps.openOffline(file, PcapHandle.TimestampPrecision.NANO);
        }

        private String DownloadAndCache(String file) throws IOException {
            String tempDirectory = FileUtils.getTempDirectoryPath();
            File pcapFile = FileSystems.getDefault().getPath(tempDirectory, RandomPackagesTest.class.getSimpleName(), file).toFile();
            FileUtils.createParentDirectories(pcapFile);
            if (!pcapFile.exists()) {
                URL source = new URL("http://kargs.net/captures/" + file);
                LOGGER.info("Downloading {}", source);
                FileUtils.copyURLToFile(source, pcapFile);
            }
            return pcapFile.getAbsolutePath();
        }

        @Override
        public void close() {
            pcapHandle.close();
        }

        @Override
        public String toString() {
            return "PCAPEvaluator{" +
                "pcapFile='" + pcapFile + '\'' +
                ", pcapHandle=" + pcapHandle +
                '}';
        }
    }

    static SkipInstruction noskip() {
        return skip(-1, SkipInstruction.SkipType.NO_SKIPPING, "there should be no reason to skip this package");
    }

    static SkipInstruction skip(int packageNumber, String reason) {
        return skip(packageNumber, SkipInstruction.SkipType.SKIP_COMPLETE, reason);
    }

    static SkipInstruction skip(int packageNumber, SkipInstruction.SkipType skipType, String reason) {
        return new SkipInstruction(packageNumber, skipType, reason);
    }

    static class SkipInstruction {
        int packageNumber;

        SkipType skipType;

        String reason;

        SkipInstruction(int packageNumber, SkipType skipType, String reason) {
            this.packageNumber = packageNumber;
            this.skipType = Objects.requireNonNull(skipType);
            this.reason = Objects.requireNonNull(reason);
        }

        enum SkipType {
            /**
             * don't do anything with that package
             */
            SKIP_COMPLETE,
            /**
             * only parse don't serialize
             */
            SKIP_SERIALIZE,
            /**
             * ony parse, serialize and don't compare
             */
            SKIP_COMPARE,
            /**
             * don't skip at all
             */
            NO_SKIPPING
        }

        boolean shouldSkipAll() {
            return skipType == SkipType.SKIP_COMPLETE;
        }

        boolean shouldParse() {
            return skipType.ordinal() > SkipType.SKIP_COMPLETE.ordinal();
        }

        boolean shouldSerialize() {
            return skipType.ordinal() > SkipType.SKIP_SERIALIZE.ordinal();
        }

        boolean shouldCompare() {
            return skipType.ordinal() > SkipType.SKIP_COMPARE.ordinal();
        }

        int getPackageNumber() {
            return packageNumber;
        }

        @Override
        public String toString() {
            return "Package " + packageNumber + " skipped with skipType=" + skipType + ". Reason: " + reason;
        }
    }

    static void appendPackageNumberToFile(int packageNumber) {
        try {
            OpenOption openOption = StandardOpenOption.CREATE_NEW;
            Path path = Paths.get("target", "collectedPackageNumbers.txt");
            if (path.toFile().exists()) {
                openOption = StandardOpenOption.APPEND;
            }
            Files.write(path, (" " + packageNumber + ",").getBytes(), openOption);
        } catch (IOException ignore) {
        }
    }
}
