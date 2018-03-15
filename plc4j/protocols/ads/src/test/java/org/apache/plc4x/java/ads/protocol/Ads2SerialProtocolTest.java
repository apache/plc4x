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
package org.apache.plc4x.java.ads.protocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.serial.AmsSerialFrame;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.Date;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThat;

@RunWith(Parameterized.class)
public class Ads2SerialProtocolTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(Ads2SerialProtocolTest.class);

    private Ads2SerialProtocol SUT;

    @Parameterized.Parameter
    public AmsSerialFrame amsSerialFrame;

    @Parameterized.Parameter(1)
    public String clazzName;

    @Parameterized.Parameter(2)
    public AmsPacket amsPacket;

    @Parameterized.Parameters(name = "{index} {1}")
    public static Collection<Object[]> data() {
        AmsNetId targetAmsNetId = AmsNetId.of("1.2.3.4.5.6");
        AmsPort targetAmsPort = AmsPort.of(7);
        AmsNetId sourceAmsNetId = AmsNetId.of("8.9.10.11.12.13");
        AmsPort sourceAmsPort = AmsPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Data data = Data.of("Hello World!".getBytes());
        return Stream.of(
            AdsAddDeviceNotificationRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(1), IndexOffset.of(1), Length.of(1), TransmissionMode.of(1), MaxDelay.of(1), CycleTime.of(1)),
            AdsAddDeviceNotificationResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                NotificationHandle.of(0)
            ),
            AdsDeleteDeviceNotificationRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                NotificationHandle.of(0)
            ),
            AdsDeleteDeviceNotificationResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            ),
            AdsDeviceNotificationRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Stamps.of(1),
                Collections.singletonList(
                    // Nano times need to be offset by (1.1.1970 - 1.1.1601) years in nanos
                    AdsStampHeader.of(TimeStamp.of(new Date()),
                        Collections.singletonList(
                            AdsNotificationSample.of(NotificationHandle.of(0), data))
                    )
                )
            ),
            AdsReadDeviceInfoRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId
            ),
            AdsReadDeviceInfoResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                MajorVersion.of((byte) 1),
                MinorVersion.of((byte) 2),
                Version.of(3),
                Device.of("Random DeviceId")
            ),
            AdsReadRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(0),
                IndexOffset.of(0),
                Length.of(1)
            ),
            AdsReadResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                data
            ),
            AdsReadStateRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId
            ),
            AdsReadStateResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            ),
            AdsReadWriteRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(0),
                IndexOffset.of(0),
                ReadLength.of(data.getCalculatedLength()),
                data
            ),
            AdsReadWriteResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0),
                data
            ),
            AdsWriteControlRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                AdsState.of(0xaffe),
                DeviceState.of(0xaffe),
                data
            ),
            AdsWriteControlResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            ),
            AdsWriteRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(0),
                IndexOffset.of(0),
                data
            ),
            AdsWriteResponse.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                Result.of(0)
            )
            /*,
            UnknownCommand.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, State.DEFAULT, invokeId,
                Unpooled.wrappedBuffer(new byte[]{42})
            )*/
        ).map(amsPacket -> new Object[]{amsPacket.toAmsSerialFrame((byte) 0), amsPacket.getClass().getSimpleName(), amsPacket}).collect(Collectors.toList());
    }

    @Before
    public void setUp() throws Exception {
        SUT = new Ads2SerialProtocol();
        byte[] bytes = amsSerialFrame.getBytes();
        LOGGER.info("amsPacket:\n{} has \n{}bytes\nHexDump:\n{}", amsSerialFrame, bytes.length, amsSerialFrame.dump());
    }

    @Test
    public void encode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.encode(null, amsSerialFrame.getAmsPacket(), out);
        assertEquals(1, out.size());
        assertThat(out, hasSize(1));
    }

    @Test
    public void decode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.decode(null, amsSerialFrame.getByteBuf(), out);
        assertThat(out, hasSize(1));
    }

    @Test
    public void roundTrip() throws Exception {
        ArrayList<Object> outbound = new ArrayList<>();
        SUT.encode(null, amsSerialFrame.getAmsPacket(), outbound);
        assertEquals(1, outbound.size());
        assertThat(outbound, hasSize(1));
        assertThat(outbound.get(0), instanceOf(ByteBuf.class));
        ByteBuf byteBuf = (ByteBuf) outbound.get(0);
        ArrayList<Object> inbound = new ArrayList<>();
        SUT.decode(null, byteBuf, inbound);
        assertEquals(1, inbound.size());
        assertThat(inbound, hasSize(1));
        assertThat(inbound.get(0), instanceOf(AmsPacket.class));
        AmsPacket inboundAmsPacket = (AmsPacket) inbound.get(0);
        assertThat("inbound divers from outbound", this.amsSerialFrame, equalTo(inboundAmsPacket.toAmsSerialFrame((byte) 0)));
    }
}