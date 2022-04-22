/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferBoxBased;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.utils.hex.Hex;
import org.apache.plc4x.test.RequirePcapNg;
import org.assertj.core.api.InstanceOfAssertFactories;
import org.junit.jupiter.api.*;
import org.pcap4j.core.*;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.UdpPacket;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.Closeable;
import java.io.File;
import java.io.IOException;
import java.math.BigInteger;
import java.net.URL;
import java.nio.file.FileSystems;
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
                    assertEquals((short) 0x8, apduUnconfirmedRequest.getServiceRequest().getServiceChoice());
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
                    assertEquals((short) 0x8, apduUnconfirmedRequest.getServiceRequest().getServiceChoice());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(260, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(260, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(18, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(18, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(24, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(24, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(260, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(260, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(86, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(86, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
                })
        );
    }

    @TestFactory
    @DisplayName("BACnet-MSTP-SNAP-Mixed")
    Collection<DynamicNode> BACnet_MSTP_SNAP_Mixed() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnet-MSTP-SNAP-Mixed.cap", BACNET_BPF_FILTER_UDP);
        return List.of(
            pcapEvaluator.parseEmAll()
        );
    }

    @TestFactory
    @DisplayName("BACnetARRAY-element-0")
    Collection<DynamicNode> BACnetARRAY_element_0() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnetARRAY-element-0.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("BACnetARRAY-elements")
    Collection<DynamicNode> BACnetARRAY_elements() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnetARRAY-elements.cap");
        return List.of(pcapEvaluator.parseEmAll());
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
                    assertEquals(BACnetPropertyIdentifier.ZONE_MEMBERS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.ZONE_MEMBERS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetConstructedDataLifeSafetyZone baCnetConstructedDataLifeSafetyZone = (BACnetConstructedDataLifeSafetyZone) baCnetServiceAckReadProperty.getValues();

                    List<BACnetContextTagObjectIdentifier> zones = baCnetConstructedDataLifeSafetyZone.getZones();
                    assertThat(zones.get(0)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 3L);
                    assertThat(zones.get(1)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 4L);
                    assertThat(zones.get(2)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 5L);
                    assertThat(zones.get(3)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 6L);
                    assertThat(zones.get(4)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 7L);
                    assertThat(zones.get(5)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 8L);
                    assertThat(zones.get(6)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 9L);
                    assertThat(zones.get(7)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 16L);
                    assertThat(zones.get(8)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 494L);
                    assertThat(zones.get(9)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 255L);
                    assertThat(zones.get(10)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 231L);
                    assertThat(zones.get(11)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 4193620L);
                    assertThat(zones.get(12)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 222L);
                    assertThat(zones.get(13)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 300L);
                    assertThat(zones.get(14)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 166L);
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
                    assertEquals(BACnetPropertyIdentifier.MEMBER_OF, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MEMBER_OF, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetConstructedDataLifeSafetyZone baCnetConstructedDataLifeSafetyZone = (BACnetConstructedDataLifeSafetyZone) baCnetServiceAckReadProperty.getValues();

                    List<BACnetContextTagObjectIdentifier> zones = baCnetConstructedDataLifeSafetyZone.getZones();
                    assertThat(zones.get(0)).extracting("objectType", "instanceNumber").contains(BACnetObjectType.LIFE_SAFETY_ZONE, 1L);
                })
        );
    }

    @TestFactory
    @DisplayName("BACnetIP-MSTP-Mix")
    Collection<DynamicNode> BACnet_MSTP_Mix() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("BACnetIP-MSTP-Mix.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
                    // TODO: change to enum
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(260L, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    // TODO: change to enum
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(260L, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetPropertyIdentifier.PRIORITY_ARRAY, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PRIORITY_ARRAY, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    /* FIXME: we get now a bunch of tags here
                    BACnetPropertyValuePriorityValue baCnetPropertyValuePriorityValue = (BACnetPropertyValuePriorityValue) ((BACnetConstructedDataUnspecified)baCnetServiceAckReadProperty.getValues()).getData();
                    assertArrayEquals(new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, baCnetPropertyValuePriorityValue.getValues());
                     */
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
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTagReal baCnetApplicationTagReal = (BACnetApplicationTagReal) ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
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
                    assertEquals(BACnetPropertyIdentifier.RELINQUISH_DEFAULT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.RELINQUISH_DEFAULT, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    /* FIXME: wrong data here too
                    BACnetApplicationTagReal baCnetApplicationTagReal = (BACnetApplicationTagReal) ((BACnetConstructedDataUnspecified)baCnetServiceAckReadProperty.getValues()).getData().get(0);
                    assertEquals(0f, baCnetApplicationTagReal);
                     */
                }),
            DynamicTest.dynamicTest("No. 29-76 - Skip Misc 48 packages",
                () -> {
                    // this is a repeat from the package above
                    pcapEvaluator.skipPackages(48);
                }),
            DynamicTest.dynamicTest("No. 77 - Confirmed-REQ writeProperty[ 1] analog-output,0 priority-array",
                () -> {
                    // This package is broken as from the spec it requires 16 values // TODO: validate that
                    pcapEvaluator.skipPackages(1);
                }),
            DynamicTest.dynamicTest("No. 78 - Error writeProperty[ 1]",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorWriteProperty baCnetErrorWriteProperty = (BACnetErrorWriteProperty) apduError.getError();
                    // TODO: change to enum
                    assertEquals(0x02, baCnetErrorWriteProperty.getErrorClass().getActualValue());
                    // TODO: change to enum
                    assertEquals(0x28, baCnetErrorWriteProperty.getErrorCode().getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetConfirmedServiceRequestWriteProperty.getPropertyIdentifier().getPropertyIdentifier());

                    BACnetApplicationTagReal baCnetApplicationTagReal = (BACnetApplicationTagReal) ((BACnetConstructedDataUnspecified) baCnetConfirmedServiceRequestWriteProperty.getPropertyValue()).getData().get(0).getApplicationTag();
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
                    BACnetErrorWriteProperty baCnetErrorWriteProperty = (BACnetErrorWriteProperty) apduError.getError();
                    // TODO: change to enum
                    assertEquals(0x02, baCnetErrorWriteProperty.getErrorClass().getActualValue());
                    // TODO: change to enum
                    assertEquals(0x25, baCnetErrorWriteProperty.getErrorCode().getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetConfirmedServiceRequestWriteProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTagObjectIdentifier objectIdentifier = (BACnetApplicationTagObjectIdentifier) ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    assertEquals(BACnetObjectType.DEVICE, objectIdentifier.getObjectType());
                    assertEquals(12345, objectIdentifier.getInstanceNumber());
                }),
            pcapEvaluator.parseFrom(203)
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
                        BACnetTimeStampSequence timestamp = (BACnetTimeStampSequence) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp();
                        assertEquals(2, timestamp.getSequenceNumber().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.UNSIGNED_RANGE, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getEventType());
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
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getEventState());
                    }
                    {
                        assertEquals(BACnetEventState.OFFNORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getEventState());
                    }
                    {
                        BACnetNotificationParametersUnsignedRange baCnetNotificationParametersUnsignedRange = (BACnetNotificationParametersUnsignedRange) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertEquals(50, baCnetNotificationParametersUnsignedRange.getSequenceNumber().getPayload().getActualValue().longValue());
                        assertTrue(baCnetNotificationParametersUnsignedRange.getStatusFlags().getInAlarm());
                        assertFalse(baCnetNotificationParametersUnsignedRange.getStatusFlags().getFault());
                        assertFalse(baCnetNotificationParametersUnsignedRange.getStatusFlags().getOverriden());
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
                        BACnetTimeStampSequence timestamp = (BACnetTimeStampSequence) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp();
                        assertEquals(2, timestamp.getSequenceNumber().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.EXTENDED, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getEventType());
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
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getEventState());
                    }
                    {
                        assertEquals(BACnetEventState.OFFNORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getEventState());
                    }
                    {
                        BACnetNotificationParametersExtended baCnetNotificationParametersExtended = (BACnetNotificationParametersExtended) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertEquals(24, baCnetNotificationParametersExtended.getVendorId().getPayload().getActualValue().longValue());
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
                        BACnetTimeStampDateTime timestamp = (BACnetTimeStampDateTime) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp();
                        BACnetTagPayloadDate payload = timestamp.getDateTimeValue().getDateValue().getPayload();
                        assertEquals(2008, payload.getYear());
                        assertEquals(5, payload.getMonth());
                        assertEquals(2, payload.getDayOfMonth());
                        assertEquals(5, payload.getDayOfWeek());
                        BACnetTagPayloadTime baCnetTagPayloadTime = timestamp.getDateTimeValue().getTimeValue().getPayload();
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
                        assertEquals(BACnetEventType.BUFFER_READY, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getEventType());
                    }
                    {
                        assertEquals(BACnetNotifyType.EVENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsTrue());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getEventState());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getEventState());
                    }
                    {
                        BACnetNotificationParametersBufferReady baCnetNotificationParametersBufferReady = (BACnetNotificationParametersBufferReady) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertEquals(BACnetObjectType.TREND_LOG, baCnetNotificationParametersBufferReady.getBufferProperty().getObjectIdentifier().getObjectType());
                        assertEquals(BACnetPropertyIdentifier.LOG_BUFFER, baCnetNotificationParametersBufferReady.getBufferProperty().getPropertyIdentifier().getPropertyIdentifier());
                        assertEquals(BACnetObjectType.DEVICE, baCnetNotificationParametersBufferReady.getBufferProperty().getDeviceIdentifier().getObjectType());
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
                    assertEquals((long) 12345, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getInitiatingDeviceIdentifier().getInstanceNumber());
                    assertEquals(BACnetObjectType.BINARY_INPUT, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getMonitoredObjectIdentifier().getObjectType());
                    assertEquals(0, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getMonitoredObjectIdentifier().getInstanceNumber());
                    assertEquals(9, baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getLifetimeInSeconds().getPayload().getActualValue().longValue() / 60);
                    {
                        BACnetContextTagPropertyIdentifier baCnetContextTagPropertyIdentifier = baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(0).getPropertyIdentifier();
                        assertEquals(BACnetPropertyIdentifier.PRESENT_VALUE, baCnetContextTagPropertyIdentifier.getPropertyIdentifier());
                    }
                    {
                        BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(0).getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                        assertEquals(0, baCnetApplicationTagEnumerated.getActualValue());
                    }
                    {
                        BACnetContextTagPropertyIdentifier baCnetContextTagPropertyIdentifier = baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(1).getPropertyIdentifier();
                        assertEquals(BACnetPropertyIdentifier.STATUS_FLAGS, baCnetContextTagPropertyIdentifier.getPropertyIdentifier());
                    }
                    {
                        BACnetApplicationTagBitString baCnetApplicationTagBitString = (BACnetApplicationTagBitString) ((BACnetConstructedDataUnspecified) baCnetUnconfirmedServiceRequestUnconfirmedCOVNotification.getListOfValues().getData().get(1).getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                        assertEquals(Arrays.asList(false, false, false, false), baCnetApplicationTagBitString.getPayload().getData());
                    }
                }),
            pcapEvaluator.parseTill(1347)
        );
    }

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
                        BACnetTimeStampSequence timestamp = (BACnetTimeStampSequence) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp();
                        assertEquals(1, timestamp.getSequenceNumber().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(111, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.CHANGE_OF_STATE, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getEventType());
                    }
                    {
                        assertEquals(BACnetNotifyType.EVENT, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsFalse());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getEventState());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getEventState());
                    }
                    {
                        BACnetNotificationParametersChangeOfState baCnetNotificationParametersChangeOfState = (BACnetNotificationParametersChangeOfState) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        assertTrue(baCnetNotificationParametersChangeOfState.getStatusFlags().getInAlarm());
                        assertFalse(baCnetNotificationParametersChangeOfState.getStatusFlags().getFault());
                        assertFalse(baCnetNotificationParametersChangeOfState.getStatusFlags().getOverriden());
                        assertFalse(baCnetNotificationParametersChangeOfState.getStatusFlags().getOutOfService());
                    }
                })
        );
    }

    @TestFactory
    @DisplayName("CriticalRoom55-1")
    Collection<DynamicNode> CriticalRoom55_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("CriticalRoom55-1.cap");
        return List.of(pcapEvaluator.parseEmAll());
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
                    BACnetApplicationTagReal baCnetApplicationTagReal = (BACnetApplicationTagReal) ((BACnetConstructedDataUnspecified) baCnetConfirmedServiceRequestWriteProperty.getPropertyValue()).getData().get(0).getApplicationTag();
                    assertEquals(123.0f, baCnetApplicationTagReal.getPayload().getValue());
                }),
            pcapEvaluator.parseFrom(2)
        );
    }

    @TestFactory
    @DisplayName("DRI%20CAVE%20log%20udp-0168-20081216-1117-03")
    Collection<DynamicNode> DRI_CAVE_log_udp_0168_20081216_1117_03() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("DRI%20CAVE%20log%20udp-0168-20081216-1117-03.cap");
        return List.of(pcapEvaluator.parseEmAll());
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
                        BACnetTimeStampDateTime timestamp = (BACnetTimeStampDateTime) baCnetConfirmedServiceRequestConfirmedEventNotification.getTimestamp();
                        assertEquals(2005, timestamp.getDateTimeValue().getDateValue().getPayload().getYear());
                        assertEquals(12, timestamp.getDateTimeValue().getDateValue().getPayload().getMonth());
                        assertEquals(8, timestamp.getDateTimeValue().getDateValue().getPayload().getDayOfMonth());
                        assertEquals(4, timestamp.getDateTimeValue().getDateValue().getPayload().getDayOfWeek());
                        assertEquals(14, timestamp.getDateTimeValue().getTimeValue().getPayload().getHour());
                        assertEquals(12, timestamp.getDateTimeValue().getTimeValue().getPayload().getMinute());
                        assertEquals(49, timestamp.getDateTimeValue().getTimeValue().getPayload().getSecond());
                        assertEquals(0, timestamp.getDateTimeValue().getTimeValue().getPayload().getFractional());
                    }
                    {
                        assertEquals(1, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotificationClass().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(200, baCnetConfirmedServiceRequestConfirmedEventNotification.getPriority().getPayload().getActualValue().longValue());
                    }
                    {
                        assertEquals(BACnetEventType.OUT_OF_RANGE, baCnetConfirmedServiceRequestConfirmedEventNotification.getEventType().getEventType());
                    }
                    {
                        assertEquals(BACnetNotifyType.ALARM, baCnetConfirmedServiceRequestConfirmedEventNotification.getNotifyType().getValue());
                    }
                    {
                        assertTrue(baCnetConfirmedServiceRequestConfirmedEventNotification.getAckRequired().getPayload().getIsTrue());
                    }
                    {
                        assertEquals(BACnetEventState.HIGH_LIMIT, baCnetConfirmedServiceRequestConfirmedEventNotification.getFromState().getEventState());
                    }
                    {
                        assertEquals(BACnetEventState.NORMAL, baCnetConfirmedServiceRequestConfirmedEventNotification.getToState().getEventState());
                    }
                    {
                        BACnetNotificationParametersComplexEventType baCnetNotificationParametersComplexEventType = (BACnetNotificationParametersComplexEventType) baCnetConfirmedServiceRequestConfirmedEventNotification.getEventValues();
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(0);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            assertEquals(2200, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagCharacterString baCnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals("StockingNAE", baCnetApplicationTagCharacterString.getValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(1);
                            assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            BACnetApplicationTagCharacterString baCnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals("StockingNAE/N2-1.NAE4-N2A-DX1.OA-T", baCnetApplicationTagCharacterString.getValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(2);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            assertEquals(2201, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(85, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(3);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            assertEquals(2202, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagReal baCnetApplicationTagReal = (BACnetApplicationTagReal) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(35.093750, baCnetApplicationTagReal.getPayload().getValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(4);
                            assertEquals(BACnetPropertyIdentifier.RELIABILITY, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(0, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(5);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            assertEquals(661, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(5, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(6);
                            assertEquals(BACnetPropertyIdentifier.UNITS, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(64, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(7);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            assertEquals(1659, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(0, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(8);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            assertEquals(2203, baCnetPropertyValue.getPropertyIdentifier().getProprietaryValue());
                            BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(0, baCnetApplicationTagEnumerated.getActualValue());
                        }
                        {
                            BACnetPropertyValue baCnetPropertyValue = baCnetNotificationParametersComplexEventType.getListOfValues().getData().get(9);
                            assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetPropertyValue.getPropertyIdentifier().getPropertyIdentifier());
                            BACnetApplicationTagUnsignedInteger baCnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) ((BACnetConstructedDataUnspecified) baCnetPropertyValue.getPropertyValue().getConstructedData()).getData().get(0).getApplicationTag();
                            assertEquals(5, baCnetApplicationTagUnsignedInteger.getPayload().getActualValue().longValue());
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
        return List.of(pcapEvaluator.parseEmAll(
            3,
            4
        ));
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
        return List.of(pcapEvaluator.parseEmAll());
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
                    assertThat(serviceRequest.getPropertyIdentifier()).extracting("propertyIdentifier").isEqualTo(BACnetPropertyIdentifier.SUBORDINATE_LIST);
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
                    assertEquals(BACnetPropertyIdentifier.SUBORDINATE_LIST, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    List<BACnetConstructedDataElement> data = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData();
                    assertThat(data)
                        .element(0)
                        .extracting(BACnetConstructedDataElement::getContextTag)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(1)
                        .extracting(BACnetConstructedDataElement::getContextTag)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_INPUT, 1L);
                    assertThat(data)
                        .element(2)
                        .extracting(BACnetConstructedDataElement::getContextTag)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(3)
                        .extracting(BACnetConstructedDataElement::getContextTag)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_INPUT, 3L);
                    assertThat(data)
                        .element(4)
                        .extracting(BACnetConstructedDataElement::getContextTag)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(5)
                        .extracting(BACnetConstructedDataElement::getContextTag)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_OUTPUT, 1L);
                    assertThat(data)
                        .element(6)
                        .extracting(BACnetConstructedDataElement::getContextTag)
                        .extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 128L);
                    assertThat(data)
                        .element(7)
                        .extracting(BACnetConstructedDataElement::getContextTag)
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
                    assertThat(serviceRequest.getPropertyIdentifier()).extracting("propertyIdentifier").isEqualTo(BACnetPropertyIdentifier.SUBORDINATE_LIST);
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
                    assertEquals(BACnetPropertyIdentifier.SUBORDINATE_LIST, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    List<BACnetConstructedDataElement> data = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData();
                    assertThat(data.get(0).getContextTag()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.DEVICE, 4000L);
                    assertThat(data.get(1).getContextTag()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_INPUT, 1L);
                    assertThat(data.get(2).getContextTag()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.ANALOG_VALUE, 1L);
                    assertThat(data.get(3).getContextTag()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.BINARY_INPUT, 1L);
                    assertThat(data.get(4).getContextTag()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.BINARY_VALUE, 1L);
                }),
            DynamicTest.dynamicTest("No. 3 - Confirmed-REQ   readProperty[144] structured-view,1 subordinate-annotations",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUConfirmedRequest apdu = (APDUConfirmedRequest) npdu.getApdu();
                    BACnetConfirmedServiceRequestReadProperty serviceRequest = (BACnetConfirmedServiceRequestReadProperty) apdu.getServiceRequest();
                    assertThat(serviceRequest.getObjectIdentifier()).extracting("objectType", "instanceNumber").contains(BACnetObjectType.STRUCTURED_VIEW, 1L);
                    assertThat(serviceRequest.getPropertyIdentifier()).extracting("propertyIdentifier").isEqualTo(BACnetPropertyIdentifier.SUBORDINATE_ANNOTATIONS);
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
                    assertEquals(BACnetPropertyIdentifier.SUBORDINATE_ANNOTATIONS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    List<BACnetConstructedDataElement> data = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData();
                    assertThat(data.get(0).getApplicationTag()).extracting("value").isEqualTo("Subordinate 1");
                    assertThat(data.get(1).getApplicationTag()).extracting("value").isEqualTo("Subordinate 2");
                    assertThat(data.get(2).getApplicationTag()).extracting("value").isEqualTo("Subordinate 3");
                    assertThat(data.get(3).getApplicationTag()).extracting("value").isEqualTo("Subordinate 4");
                })
        );
    }

    @TestFactory
    @DisplayName("SubordinateListDecodeSample")
    Collection<DynamicNode> SubordinateListDecodeSample() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SubordinateListDecodeSample.pcap");
        return List.of(pcapEvaluator.parseEmAll(
            1,
            4
        ));
    }

    @TestFactory
    @DisplayName("SynergyBlinkWarn")
    Collection<DynamicNode> SynergyBlinkWarn() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SynergyBlinkWarn.cap");
        return List.of(pcapEvaluator.parseEmAll(
            3,
            5
        ));
    }

    @TestFactory
    @DisplayName("SynergyReadProperties")
    Collection<DynamicNode> SynergyReadProperties() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SynergyReadProperties.cap");
        return List.of(pcapEvaluator.parseEmAll(
            37,
            38
        ));
    }

    @TestFactory
    @DisplayName("SynergyWriteProperty")
    Collection<DynamicNode> SynergyWriteProperty() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("SynergyWriteProperty.cap");
        return List.of(pcapEvaluator.parseEmAll(
            3,
            5
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
            25,
            26,
            27
        ));
    }

    @TestFactory
    @DisplayName("Sysco-3")
    Collection<DynamicNode> Sysco_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Sysco-3.cap");
        return List.of(pcapEvaluator.parseEmAll(
            25,
            26,
            27
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
        return List.of(pcapEvaluator.parseEmAll(
            3,
            4,
            51,
            52,
            117
        ));
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
        return List.of(pcapEvaluator.parseEmAll(
            1,
            4,
            5,
            13,
            14,
            15,
            20,
            25,
            205,
            206,
            350,
            351,
            451
        ));
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

    @TestFactory
    @DisplayName("Tower333 lighting 5min IP")
    Collection<DynamicNode> Tower333_lighting_5min_IP() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("Tower333%20lighting%205min%20IP.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
                        .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
                        .isEqualTo(BACnetPropertyIdentifier.EXCEPTION_SCHEDULE);
                    assertThat(serviceRequest.getArrayIndex()).extracting("payload").extracting("actualValue").isEqualTo(BigInteger.ZERO);
                }),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[ 74] schedule,1 exception-schedule",
                () -> {
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
                    assertEquals(BACnetPropertyIdentifier.EXCEPTION_SCHEDULE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                            .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
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
                            .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
                            .isEqualTo(BACnetPropertyIdentifier.ACTION);
                        assertThat(baCnetServiceAckReadProperty.getValues())
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataCommand.class))
                            .extracting(BACnetConstructedDataCommand::getAction)
                            .satisfies(baCnetActionCommands -> {
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
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagPropertyIdentifier.class))
                                            .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
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
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagPropertyIdentifier.class))
                                            .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
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
                                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetContextTagPropertyIdentifier.class))
                                            .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
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
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);

                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUComplexAck apduComplexAck = (APDUComplexAck) npdu.getApdu();
                    BACnetServiceAckReadProperty baCnetServiceAckReadProperty = (BACnetServiceAckReadProperty) apduComplexAck.getServiceAck();
                    assertNotNull(baCnetServiceAckReadProperty);
                    assertEquals(BACnetObjectType.DEVICE, baCnetServiceAckReadProperty.getObjectIdentifier().getObjectType());
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK readProperty[155] device,42222 protocol-conformance-class",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            pcapEvaluator.parseFrom(3)
        );
    }

    @TestFactory
    @DisplayName("alerton-plugfest-3")
    Collection<DynamicNode> alerton_plugfest_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("alerton-plugfest-3.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
            5,
            8,
            11
        ));
    }

    @TestFactory
    @DisplayName("atomic-read-file-50x1500k")
    Collection<DynamicNode> atomic_read_file_50x1500k() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-read-file-50x1500k.cap");
        return List.of(pcapEvaluator.parseEmAll(
            35,
            60,
            75,
            173,
            201,
            216
        ));
    }

    @TestFactory
    @DisplayName("atomic-read-file-480")
    Collection<DynamicNode> atomic_read_file_480() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-read-file-480.cap");
        return List.of(pcapEvaluator.parseEmAll(
            5,
            8,
            9
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
            1,
            3,
            6,
            1173,
            1179,
            1180
        ));
    }

    @TestFactory
    @DisplayName("atomic-write-file-3")
    Collection<DynamicNode> atomic_write_file_3() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("atomic-write-file-3.cap");
        return List.of(pcapEvaluator.parseEmAll(
            3,
            4,
            67,
            71
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
            1,
            4,
            3122
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
        return List.of(pcapEvaluator.parseEmAll());
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
            2
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
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bacnet-stack-services")
    Collection<DynamicNode> bacnet_stack_services() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bacnet-stack-services.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 22 - Unconfirmed-REQ unconfirmedCOVNotification device,1 accumulator,22 present-value status-flags",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 23 - Unconfirmed-REQ unconfirmedCOVNotification device,1 binary-input,217 present-value status-flags",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 24 - Unconfirmed-REQ unconfirmedCOVNotification device,1 accumulator,21 present-value status-flags",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 25 - Unconfirmed-REQ unconfirmedCOVNotification device,1 binary-input,217 present-value status-flags",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 26 - Unconfirmed-REQ unconfirmedCOVNotification device,1 binary-output,1 present-value status-flags",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 26 - Unconfirmed-REQ unconfirmedEventNotification device,1 binary-output,1",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 28 - Unconfirmed-REQ unconfirmedCOVNotification device,1 accumulator,22 present-value status-flags",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                })
        );
    }

    @TestFactory
    @DisplayName("btl-plugfest")
    Collection<DynamicNode> btl_plugfest() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("btl-plugfest.cap");
        return List.of(pcapEvaluator.parseEmAll());
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
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("bvlc-fdreg-readprop-47809")
    Collection<DynamicNode> bvlc_fdreg_readprop_47809() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("bvlc-fdreg-readprop-47809.cap");
        return List.of(pcapEvaluator.parseEmAll());
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
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
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
            2
        ));
    }

    @TestFactory
    @DisplayName("epics-2")
    Collection<DynamicNode> epics_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("epics-2.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
                            .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
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
                            .extracting(BACnetContextTagPropertyIdentifier::getPropertyIdentifier)
                            .isEqualTo(BACnetPropertyIdentifier.EVENT_TIME_STAMPS);
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getArrayIndex)
                            .isNull();
                        assertThat(baCnetServiceAckReadProperty)
                            .extracting(BACnetServiceAckReadProperty::getValues)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetConstructedDataEventTimestamps.class))
                            .satisfies(baCnetConstructedDataEventTimestamps ->
                                assertThat(baCnetConstructedDataEventTimestamps)
                                    .extracting(BACnetConstructedDataEventTimestamps::getToOffnormal)
                                    .extracting(BACnetContextTagTime::getPayload)
                                    .extracting(BACnetTagPayloadTime::getHourIsWildcard, BACnetTagPayloadTime::getMinuteIsWildcard, BACnetTagPayloadTime::getSecondIsWildcard, BACnetTagPayloadTime::getFractionalIsWildcard)
                                    .containsExactly(true, true, true, true)
                            )
                            .satisfies(baCnetConstructedDataEventTimestamps ->
                                assertThat(baCnetConstructedDataEventTimestamps)
                                    .extracting(BACnetConstructedDataEventTimestamps::getToFault)
                                    .extracting(BACnetContextTagUnsignedInteger::getPayload)
                                    .extracting(BACnetTagPayloadUnsignedInteger::getActualValue)
                                    .isEqualTo(BigInteger.ZERO)
                            )
                            .satisfies(baCnetConstructedDataEventTimestamps ->
                                assertThat(baCnetConstructedDataEventTimestamps)
                                    .extracting(BACnetConstructedDataEventTimestamps::getToNormal)
                                    .satisfies(baCnetDateTime ->
                                        assertThat(baCnetDateTime)
                                            .extracting(BACnetDateTime::getDateValue)
                                            .extracting(BACnetApplicationTagDate::getPayload)
                                            .extracting(BACnetTagPayloadDate::getYearIsWildcard, BACnetTagPayloadDate::getMonthIsWildcard, BACnetTagPayloadDate::getDayOfMonthIsWildcard, BACnetTagPayloadDate::getDayOfWeekIsWildcard)
                                            .containsExactly(true, true, true, true)
                                    )
                                    .satisfies(baCnetDateTime ->
                                        assertThat(baCnetDateTime)
                                            .extracting(BACnetDateTime::getTimeValue)
                                            .extracting(BACnetApplicationTagTime::getPayload)
                                            .extracting(BACnetTagPayloadTime::getFractionalIsWildcard, BACnetTagPayloadTime::getMinuteIsWildcard, BACnetTagPayloadTime::getSecondIsWildcard, BACnetTagPayloadTime::getFractionalIsWildcard)
                                            .containsExactly(true, true, true, true)
                                    )
                            )
                        ;
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
        return List.of(pcapEvaluator.parseEmAll());
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
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("plugfest-2011-siemens-1")
    Collection<DynamicNode> plugfest_2011_siemens_1() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-2011-siemens-1.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
            pcapEvaluator.parseFrom(2)
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
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("plugfest-tridium-2")
    Collection<DynamicNode> plugfest_tridium_2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("plugfest-tridium-2.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("polarsoft-free-range-router-init-routing-table")
    Collection<DynamicNode> polarsoft_free_range_router_init_routing_table() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("polarsoft-free-range-router-init-routing-table.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("polarsoft-free-range-router")
    Collection<DynamicNode> polarsoft_free_range_router() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("polarsoft-free-range-router.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("properties")
    Collection<DynamicNode> properties() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("properties.cap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
                        BACnetServiceAckAtomicReadFile baCnetServiceAck = (BACnetServiceAckAtomicReadFile) BACnetServiceAckAtomicReadFile.staticParse(new ReadBufferByteBased(bytes),bytes.length);
                        assertThat(baCnetServiceAck)
                            .isNotNull()
                            .extracting(BACnetServiceAckAtomicReadFile::getAccessMethod)
                            .asInstanceOf(InstanceOfAssertFactories.type(BACnetServiceAckAtomicReadFileStream.class))
                            .extracting(BACnetServiceAckAtomicReadFileStream::getFileData)
                            .extracting(BACnetApplicationTagOctetString::getValue)
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedReceive());
                    assertEquals(42, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
                }),
            DynamicTest.dynamicTest("No. 2 - Unconfirmed-REQ who-Is",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    BVLCOriginalBroadcastNPDU bvlcOriginalBroadcastNPDU = (BVLCOriginalBroadcastNPDU) bvlc;
                    APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) bvlcOriginalBroadcastNPDU.getNpdu().getApdu();
                    BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) apduUnconfirmedRequest.getServiceRequest();
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
                    assertFalse(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedBoth());
                    assertEquals(42, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedBoth());
                    assertEquals(18, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedBoth());
                    assertEquals(42, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString baCnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) value;
                    assertEquals(0x0, baCnetApplicationTagEnumerated.getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("Alerton Technologies, Inc.", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("LSi Controller", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("BACtalk LSi   v3.10 A         ", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("LSi Controller v3.11 E", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(1, BACnetApplicationTagUnsignedInteger.getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagBitString BACnetApplicationTagBitString = (BACnetApplicationTagBitString) value;
                    assertEquals(Arrays.asList(true, false, true, true, false, true, true, true, false, false, true, true, true, false, true, true, true, true, true, false, true, false, false, false, false, false, true, true, false, false, true, false, true, true, true), BACnetApplicationTagBitString.getPayload().getData());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagBitString BACnetApplicationTagBitString = (BACnetApplicationTagBitString) value;
                    assertEquals(Arrays.asList(false, false, true, false, false, true, true, false, true, true, true, false, false, false, false, true, true, true), BACnetApplicationTagBitString.getPayload().getData());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(1476, BACnetApplicationTagUnsignedInteger.getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagEnumerated BACnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) value;
                    assertEquals(0, BACnetApplicationTagEnumerated.getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagTime BACnetApplicationTagTime = (BACnetApplicationTagTime) value;
                    assertEquals(15, BACnetApplicationTagTime.getPayload().getHour());
                    assertEquals(28, BACnetApplicationTagTime.getPayload().getMinute());
                    assertEquals(41, BACnetApplicationTagTime.getPayload().getSecond());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagDate BACnetApplicationTagDate = (BACnetApplicationTagDate) value;
                    assertEquals(2005, BACnetApplicationTagDate.getPayload().getYear());
                    assertEquals(9, BACnetApplicationTagDate.getPayload().getMonth());
                    assertEquals(1, BACnetApplicationTagDate.getPayload().getDayOfMonth());
                    assertEquals(4, BACnetApplicationTagDate.getPayload().getDayOfWeek());
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
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagSignedInteger BACnetApplicationTagSignedInteger = (BACnetApplicationTagSignedInteger) value;
                    assertEquals(0, BACnetApplicationTagSignedInteger.getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagBoolean BACnetApplicationTagBoolean = (BACnetApplicationTagBoolean) value;
                    assertFalse(BACnetApplicationTagBoolean.getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(6000, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint16());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(6000, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint16());
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
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals((short) 3, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint8());
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
                    assertEquals(BACnetPropertyIdentifier.TIME_SYNCHRONIZATION_RECIPIENTS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                }),
            DynamicTest.dynamicTest("No. 56 - ERROR           readProperty[ 51] device,201",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorReadProperty baCnetErrorReadProperty = (BACnetErrorReadProperty) apduError.getError();
                    // TODO: use proper enums
                    assertEquals(32, baCnetErrorReadProperty.getErrorCode().getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals((short) 127, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint8());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals((short) 40, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint8());
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
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                }),
            DynamicTest.dynamicTest("No. 66 - Complex-ACK     readProperty[ 56] device,201 Error",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorReadProperty baCnetErrorReadProperty = (BACnetErrorReadProperty) apduError.getError();
                    // TODO: use proper enums
                    assertEquals(32, baCnetErrorReadProperty.getErrorCode().getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagObjectIdentifier BACnetApplicationTagObjectIdentifier = (BACnetApplicationTagObjectIdentifier) value;
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagObjectIdentifier = (BACnetApplicationTagCharacterString) value;
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.OBJECT_TYPE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) value;
                    assertEquals(8, baCnetApplicationTagEnumerated.getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.SYSTEM_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagEnumerated baCnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) value;
                    assertEquals(0, baCnetApplicationTagEnumerated.getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("Lithonia Lighting, Inc.", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_IDENTIFIER, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MODEL_NAME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("SYSC MLX", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.FIRMWARE_REVISION, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("2x62i", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.APPLICATION_SOFTWARE_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagCharacterString BACnetApplicationTagCharacterString = (BACnetApplicationTagCharacterString) value;
                    assertEquals("10-Nov-2004", BACnetApplicationTagCharacterString.getValue());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_VERSION, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(1, BACnetApplicationTagUnsignedInteger.getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_CONFORMANCE_CLASS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_SERVICES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagBitString BACnetApplicationTagBitString = (BACnetApplicationTagBitString) value;
                    assertEquals(Arrays.asList(false, false, false, false, false, false, true, true, false, false, false, false, true, false, true, true, false, true, false, false, true, false, false, false, false, false, true, true, false, false, false, false, true, true, true), BACnetApplicationTagBitString.getPayload().getData());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.PROTOCOL_OBJECT_TYPES_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagBitString BACnetApplicationTagBitString = (BACnetApplicationTagBitString) value;
                    assertEquals(Arrays.asList(true, true, true, true, true, true, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false), BACnetApplicationTagBitString.getPayload().getData());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_APDU_LENGTH_ACCEPTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(480, BACnetApplicationTagUnsignedInteger.getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.SEGMENTATION_SUPPORTED, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagEnumerated BACnetApplicationTagEnumerated = (BACnetApplicationTagEnumerated) value;
                    assertEquals(0, BACnetApplicationTagEnumerated.getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_TIME, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagTime BACnetApplicationTagTime = (BACnetApplicationTagTime) value;
                    assertEquals(15, BACnetApplicationTagTime.getPayload().getHour());
                    assertEquals(26, BACnetApplicationTagTime.getPayload().getMinute());
                    assertEquals(40, BACnetApplicationTagTime.getPayload().getSecond());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.LOCAL_DATE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagDate BACnetApplicationTagDate = (BACnetApplicationTagDate) value;
                    assertEquals(2005, BACnetApplicationTagDate.getPayload().getYear());
                    assertEquals(9, BACnetApplicationTagDate.getPayload().getMonth());
                    assertEquals(1, BACnetApplicationTagDate.getPayload().getDayOfMonth());
                    assertEquals(255, BACnetApplicationTagDate.getPayload().getDayOfWeek());
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
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.UTC_OFFSET, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagSignedInteger BACnetApplicationTagSignedInteger = (BACnetApplicationTagSignedInteger) value;
                    assertEquals(BigInteger.valueOf(-300), BACnetApplicationTagSignedInteger.getPayload().getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.DAYLIGHT_SAVINGS_STATUS, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagBoolean BACnetApplicationTagBoolean = (BACnetApplicationTagBoolean) value;
                    assertTrue(BACnetApplicationTagBoolean.getPayload().getIsTrue());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_SEGMENT_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(8000, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint16());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.APDU_TIMEOUT, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals(8000, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint16());
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
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.NUMBER_OF_APDU_RETRIES, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals((short) 3, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint8());
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
                    assertEquals(BACnetPropertyIdentifier.TIME_SYNCHRONIZATION_RECIPIENTS, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                }),
            DynamicTest.dynamicTest("No. 112 - ERROR           readProperty[ 79] device,61",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    NPDU npdu = ((BVLCOriginalUnicastNPDU) bvlc).getNpdu();
                    APDUError apduError = (APDUError) npdu.getApdu();
                    BACnetErrorReadProperty baCnetErrorReadProperty = (BACnetErrorReadProperty) apduError.getError();
                    // TODO: use proper enums
                    assertEquals(32, baCnetErrorReadProperty.getErrorCode().getActualValue());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_MASTER, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals((short) 127, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint8());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.MAX_INFO_FRAMES, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
                    BACnetApplicationTag value = ((BACnetConstructedDataUnspecified) baCnetServiceAckReadProperty.getValues()).getData().get(0).getApplicationTag();
                    BACnetApplicationTagUnsignedInteger BACnetApplicationTagUnsignedInteger = (BACnetApplicationTagUnsignedInteger) value;
                    assertEquals((short) 1, BACnetApplicationTagUnsignedInteger.getPayload().getValueUint8());
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
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.DEVICE_ADDRESS_BINDING, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetConfirmedServiceRequestReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
                    assertEquals(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, baCnetServiceAckReadProperty.getPropertyIdentifier().getPropertyIdentifier());
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
        return List.of(pcapEvaluator.parseEmAll());
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
            3,
            4,
            5
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
        return List.of(pcapEvaluator.parseEmAll());
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
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("rp-shed-level")
    Collection<DynamicNode> rp_shed_level() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("rp-shed-level.cap");
        return Arrays.asList(
            DynamicTest.dynamicTest("No. 1 - Confirmed-REQ   readProperty[  1] load-control,0 expected-shed-level",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            DynamicTest.dynamicTest("No. 2 - Complex-ACK     readProperty[  1] load-control,0 expected-shed-level",
                () -> {
                    BVLC bvlc = pcapEvaluator.nextBVLC();
                    dump(bvlc);
                    // TODO:
                    assumeTrue(false, "not properly implemented. Check manually and add asserts");
                }),
            pcapEvaluator.parseFrom(3)
        );
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
        return List.of(pcapEvaluator.parseEmAll());
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
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("state_text_good")
    Collection<DynamicNode> state_text_good() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("state_text_good.cap");
        return List.of(pcapEvaluator.parseEmAll());
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
        return List.of(pcapEvaluator.parseEmAll(
            5,
            8,
            9,
            10,
            13,
            55,
            56,
            57
        ));
    }

    @TestFactory
    @DisplayName("synergy-device")
    Collection<DynamicNode> synergy_device() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("synergy-device.cap");
        return List.of(pcapEvaluator.parseEmAll(
            1,
            4
        ));
    }

    @TestFactory
    @DisplayName("time-sync")
    Collection<DynamicNode> time_sync() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("time-sync.cap");
        return List.of(pcapEvaluator.parseEmAll(
            2
        ));
    }

    @TestFactory
    @DisplayName("tridium jace2")
    Collection<DynamicNode> tridium_jace2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("tridium%20jace2.pcap", BACNET_BPF_FILTER_UDP);
        return List.of(pcapEvaluator.parseEmAll());
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
        return List.of(pcapEvaluator.parseEmAll());
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
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestWhoHas.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHas.getObjectIdentifier().getInstanceNumber());
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
                    // TODO: we should use a enum here
                    assertTrue(baCnetUnconfirmedServiceRequestIAm.getSegmentationSupported().getIsSegmentedBoth());
                    assertEquals(42, baCnetUnconfirmedServiceRequestIAm.getVendorId().getPayload().getActualValue().longValue());
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
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestWhoHas.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHas.getObjectIdentifier().getInstanceNumber());
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
                    assertEquals(BACnetObjectType.DEVICE, baCnetUnconfirmedServiceRequestWhoHas.getObjectIdentifier().getObjectType());
                    assertEquals(133, baCnetUnconfirmedServiceRequestWhoHas.getObjectIdentifier().getInstanceNumber());
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
        return List.of(pcapEvaluator.parseEmAll());
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
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("wp_weekly_schedule")
    Collection<DynamicNode> wp_weekly_schedule() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("wp_weekly_schedule.cap");
        return List.of(pcapEvaluator.parseEmAll());
    }

    @TestFactory
    @DisplayName("write-property-array")
    Collection<DynamicNode> write_property_array() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("write-property-array.cap");
        return List.of(pcapEvaluator.parseEmAll());
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
            1594, // Malformed Package
            1595,   // Malformed Package
            1596,   // Malformed Package
            1597,   // Malformed Package
            1598,   // Malformed Package
            1599,   // Malformed Package
            1600,   // Malformed Package
            1601,   // Malformed Package
            1602,   // Malformed Package
            1603,   // Malformed Package
            1604,   // Malformed Package
            1605,   // Malformed Package
            1606,   // Malformed Package
            1607,   // Malformed Package
            1608,   // Malformed Package
            1609    // Malformed Package
        ));
    }

    @TestFactory
    @DisplayName("write-property2")
    Collection<DynamicNode> write_property2() throws Exception {
        TestPcapEvaluator pcapEvaluator = pcapEvaluator("write-property2.cap");
        return List.of(pcapEvaluator.parseEmAll(
            10,
            14
        ));
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

        public DynamicContainer parseEmAll(int... skippedNumbers) {
            return parseRange(1, maxPackages, skippedNumbers);
        }

        public DynamicContainer parseFrom(int startPackageNumber, int... skippedNumbers) {
            return parseRange(startPackageNumber, maxPackages, skippedNumbers);
        }

        public DynamicContainer parseTill(int packageNumber, int... skippedNumbers) {
            return parseRange(1, packageNumber, skippedNumbers);
        }

        public DynamicContainer parseRange(int startPackageNumber, int packageNumber, int... skippedNumbers) {
            Set<Integer> numbersToSkip = Arrays.stream(skippedNumbers).boxed().collect(Collectors.toSet());
            boolean hasNoSkipping = skippedNumbers.length <= 0;
            boolean hasNoPackageNumberMapping = packageNumbers == null;
            Function<Integer, Integer> packageNumResolver = i -> hasNoPackageNumberMapping ? i : packageNumbers.get(i - 1);
            return DynamicContainer.dynamicContainer("Parse em all (No. " + startPackageNumber + "-" + packageNumber + ")", () -> IntStream.range(startPackageNumber, packageNumber + 1).filter(i -> hasNoSkipping || hasNoPackageNumberMapping || !numbersToSkip.contains(packageNumResolver.apply(i))).mapToObj((i) -> DynamicTest.dynamicTest("No. " + packageNumResolver.apply(i) + " - Test nr." + i, () -> {
                if (!hasNoSkipping && hasNoPackageNumberMapping && numbersToSkip.contains(i)) {
                    LOGGER.info("Skipping unfiltered package {}", i);
                    skipPackages(1);
                    return;
                }
                BVLC bvlc = nextBVLC();
                LOGGER.info("Test number {} is package number {}", i, currentPackageNumber);
                assumeTrue(bvlc != null, "No more package left");
                dump(bvlc);
            })).map(DynamicNode.class::cast).iterator());
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

        public BVLC nextBVLC() throws NotOpenException, ParseException {
            return nextBVLC(null);
        }

        public BVLC nextBVLC(Integer ensurePackageNumber) throws NotOpenException, ParseException {
            Packet packet = nextPacket();
            if (packet == null) {
                return null;
            }
            if (ensurePackageNumber != null) {
                if (ensurePackageNumber < 0) {
                    throw new IllegalArgumentException("Seached package number must be > 0");
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
            UdpPacket udpPacket = packet.get(UdpPacket.class);
            assumeTrue(udpPacket != null, "nextBVLC assumes a UDP Packet. If non is there it might by LLC");
            LOGGER.info("Handling UDP\n{}", udpPacket);
            byte[] rawData = udpPacket.getPayload().getRawData();
            LOGGER.info("Reading BVLC from:\n{}", Hex.dump(rawData));
            try {
                return BVLC.staticParse(new ReadBufferByteBased(rawData));
            } catch (ParseException e) {
                throw new ParseException(String.format("Caught at current package number: %d. Packages read so far %d", currentPackageNumber, readPackages), e);
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
}
