/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.netty;

import org.apache.commons.io.HexDump;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.junit.Assert.assertEquals;

@RunWith(Parameterized.class)
public class ADSProtocolTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(ADSProtocolTest.class);

    private ADSProtocol SUT;

    @Parameterized.Parameter
    public AMSTCPPaket amstcpPaket;

    @Parameterized.Parameters(name = "{0} {index}")
    public static Collection<Object[]> data() {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Data data = Data.of("Hello World!".getBytes());
        return Stream.of(
            new ADSAddDeviceNotificationRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(1), IndexOffset.of(1), Length.of(1), TransmissionMode.of(1), MaxDelay.of(1), CycleTime.of(1)),
            new ADSAddDeviceNotificationResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                NotificationHandle.of(0)
            ),
            new ADSDeleteDeviceNotificationRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                NotificationHandle.of(0)
            ),
            new ADSDeleteDeviceNotificationResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            ),
            new ADSDeviceNotificationRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                // TODO: must be calculated when using this constructor
                Length.of(32),
                Stamps.of(1),
                Collections.singletonList(
                    // Nano times need to be offset by (1.1.1970 - 1.1.1601) years in nanos
                    AdsStampHeader.of(TimeStamp.of(System.nanoTime()),
                        Samples.of(1),
                        Collections.singletonList(
                            AdsNotificationSample.of(NotificationHandle.of(0), SampleSize.of(data.getCalculatedLength()), data))
                    )
                )
            ),
            new ADSReadDeviceInfoRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId
            ),
            new ADSReadDeviceInfoResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                MajorVersion.of((byte) 1),
                MinorVersion.of((byte) 2),
                Version.of(3),
                Device.of(
                    (byte) 1, (byte) 2, (byte) 3, (byte) 4, (byte) 5, (byte) 6, (byte) 7, (byte) 8,
                    (byte) 9, (byte) 10, (byte) 11, (byte) 12, (byte) 13, (byte) 14, (byte) 15, (byte) 16
                )
            ),
            new ADSReadRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(0),
                IndexOffset.of(0),
                Length.of(1)
            ),
            new ADSReadResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                Length.of(data.getCalculatedLength()),
                data
            ),
            new ADSReadStateRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId
            ),
            new ADSReadStateResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            ),
            new ADSReadWriteRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(0),
                IndexOffset.of(0),
                ReadLength.of(data.getCalculatedLength()),
                WriteLength.of(data.getCalculatedLength()),
                data
            ),
            new ADSReadWriteResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                Length.of(data.getCalculatedLength()),
                data
            ),
            new ADSWriteControlRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                ADSState.of(0xaffe),
                DeviceState.of(0xaffe),
                Length.of(data.getCalculatedLength()),
                data
            ),
            new ADSWriteControlResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            ),
            new ADSWriteRequest(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(0),
                IndexOffset.of(0),
                Length.of(data.getCalculatedLength()),
                data
            ),
            new ADSWriteResponse(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            )/*,
            new UnknownCommand(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Unpooled.wrappedBuffer(new byte[]{42})
            )*/
        ).map(amstcpPaket -> new Object[]{amstcpPaket}).collect(Collectors.toList());
    }

    @Before
    public void setUp() throws Exception {
        SUT = new ADSProtocol();
        byte[] bytes = amstcpPaket.getBytes();
        LOGGER.info("amstcpPaket:\n{} has \n{}bytes", amstcpPaket, bytes.length);
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            HexDump.dump(bytes, 0, byteArrayOutputStream, 0);
            byteArrayOutputStream.flush();
            LOGGER.info("HexDump:\n{}", byteArrayOutputStream);
        }
    }

    @Test
    public void encode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.encode(null, amstcpPaket, out);
        assertEquals(1, out.size());
        // TODO: replace with hamcrest
        //assertThat(out, hasSize(1));
    }

    @Test
    public void decode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.decode(null, amstcpPaket.getByteBuf(), out);
        assertEquals(1, out.size());
        // TODO: replace with hamcrest
        //assertThat(out, hasSize(1));
    }

}