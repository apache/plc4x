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
package org.apache.plc4x.java.ads.protocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import java.util.Collections;
import java.util.Date;
import java.util.stream.Stream;

import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.assertThat;

public class AbstractProtocolTest {

    public static Stream<AmsPacket> amsPacketStream() {
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
                IndexOffset.NONE,
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
                Result.of(0),
                AdsState.of(0),
                DeviceState.of(0)
            ),
            AdsReadWriteRequest.of(
                targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
                IndexGroup.of(0),
                IndexOffset.NONE,
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
                IndexOffset.NONE,
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
        );
    }

    public void assertByteBufferEquals(ByteBuf expected, ByteBuf actual) {
        byte[] expectedBytes = new byte[expected.readableBytes()];
        expected.readBytes(expectedBytes);
        byte[] actualBytes = new byte[actual.readableBytes()];
        actual.readBytes(actualBytes);
        assertThat(expectedBytes, equalTo(actualBytes));
    }
}
