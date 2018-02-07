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

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class ADSProtocol extends MessageToMessageCodec<ByteBuf, AMSTCPPaket> {

    private static final Logger LOGGER = LoggerFactory.getLogger(ADSProtocol.class);

    private Map<Invoke, AMSTCPPaket> requests;

    public ADSProtocol() {
        this.requests = new HashMap<>();
    }

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, AMSTCPPaket amstcpPaket, List<Object> out) throws Exception {
        Invoke invokeId = amstcpPaket.getAmsHeader().getInvokeId();
        if (invokeId != Invoke.NONE) {
            requests.put(invokeId, amstcpPaket);
        }
        out.add(amstcpPaket.getByteBuf());
    }


    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) throws Exception {
        // Reserved
        byteBuf.skipBytes(2);
        long packetLength = byteBuf.readUnsignedIntLE();
        AMSNetId targetAmsNetId = AMSNetId.of(byteBuf);
        AMSPort targetAmsPort = AMSPort.of(byteBuf);
        AMSNetId sourceAmsNetId = AMSNetId.of(byteBuf);
        AMSPort sourceAmsPort = AMSPort.of(byteBuf);
        Command commandId = Command.of(byteBuf);
        State stateId = State.of(byteBuf);
        DataLength dataLength = DataLength.of(byteBuf);
        AMSError errorCode = AMSError.of(byteBuf);
        Invoke invoke = Invoke.of(byteBuf);
        AMSTCPPaket correlatedAmstcpPacket = requests.get(invoke);
        LOGGER.debug("Correlated packet received {}", correlatedAmstcpPacket);
        if (dataLength.getAsLong() > Integer.MAX_VALUE) {
            throw new IllegalStateException("Overflow in datalength: " + dataLength.getAsLong());
        }
        ByteBuf commandBuffer = byteBuf.readBytes((int) dataLength.getAsLong());
        AMSTCPHeader amstcpHeader = AMSTCPHeader.of(packetLength);
        AMSHeader amsHeader = AMSHeader.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, commandId, stateId, dataLength, errorCode, invoke);
        switch (commandId) {
            case Invalid:
                out.add(new UnknownCommand(amstcpHeader, amsHeader, commandBuffer));
                break;
            case ADS_Read_Device_Info:
                if (stateId.isRequest()) {
                    out.add(new ADSReadDeviceInfoRequest(amstcpHeader, amsHeader));
                } else {
                    Result result = Result.of(commandBuffer);
                    MajorVersion majorVersion = MajorVersion.of(commandBuffer);
                    MinorVersion minorVersion = MinorVersion.of(commandBuffer);
                    Version version = Version.of(commandBuffer);
                    Device device = Device.of(commandBuffer);
                    out.add(new ADSReadDeviceInfoResponse(amstcpHeader, amsHeader, result, majorVersion, minorVersion, version, device));
                }
                break;
            case ADS_Read:
                if (stateId.isRequest()) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer);
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer);
                    Length length = Length.of(commandBuffer);
                    out.add(new ADSReadRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, length));
                } else {
                    Result result = Result.of(commandBuffer);
                    Length length = Length.of(byteBuf);
                    if (length.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) length.getAsLong()).array());
                    out.add(new ADSReadResponse(amstcpHeader, amsHeader, result, length, data));
                }
                break;
            case ADS_Write:
                if (stateId.isRequest()) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer);
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer);
                    Length length = Length.of(byteBuf);
                    if (length.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) length.getAsLong()).array());
                    out.add(new ADSWriteRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, length, data));
                } else {
                    Result result = Result.of(commandBuffer);
                    out.add(new ADSWriteResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Read_State:
                if (stateId.isRequest()) {
                    out.add(new ADSReadStateRequest(amstcpHeader, amsHeader));
                } else {
                    Result result = Result.of(commandBuffer);
                    out.add(new ADSReadStateResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Write_Control:
                if (stateId.isRequest()) {
                    ADSState adsState = ADSState.of(commandBuffer);
                    DeviceState deviceState = DeviceState.of(commandBuffer);
                    Length length = Length.of(byteBuf);
                    if (length.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) length.getAsLong()).array());
                    out.add(new ADSWriteControlRequest(amstcpHeader, amsHeader, adsState, deviceState, length, data));
                } else {
                    Result result = Result.of(commandBuffer);
                    out.add(new ADSWriteControlResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Add_Device_Notification:
                if (stateId.isRequest()) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer);
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer);
                    Length length = Length.of(commandBuffer);
                    TransmissionMode transmissionMode = TransmissionMode.of(commandBuffer);
                    MaxDelay maxDelay = MaxDelay.of(commandBuffer);
                    CycleTime cycleTime = CycleTime.of(commandBuffer);
                    out.add(new ADSAddDeviceNotificationRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime));
                } else {
                    Result result = Result.of(commandBuffer);
                    NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer);
                    out.add(new ADSAddDeviceNotificationResponse(amstcpHeader, amsHeader, result, notificationHandle));
                }
                break;
            case ADS_Delete_Device_Notification:
                if (stateId.isRequest()) {
                    NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer);
                    out.add(new ADSDeleteDeviceNotificationRequest(amstcpHeader, amsHeader, notificationHandle));
                } else {
                    Result result = Result.of(commandBuffer);
                    out.add(new ADSDeleteDeviceNotificationResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Device_Notification:
                if (stateId.isRequest()) {
                    Length length = Length.of(commandBuffer);
                    if (length.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
                    }
                    Stamps stamps = Stamps.of(commandBuffer);
                    if (stamps.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + stamps.getAsLong());
                    }
                    ByteBuf adsDeviceNotificationBuffer = commandBuffer.readBytes((int) length.getAsLong());
                    List<AdsStampHeader> adsStampHeaders = new ArrayList<>((int) stamps.getAsLong());
                    for (int i = 1; i <= stamps.getAsLong(); i++) {
                        TimeStamp timeStamp = TimeStamp.of(adsDeviceNotificationBuffer);
                        Samples samples = Samples.of(adsDeviceNotificationBuffer);

                        List<AdsNotificationSample> adsNotificationSamples = new LinkedList<>();
                        for (int j = 1; j <= samples.getAsLong(); j++) {
                            NotificationHandle notificationHandle = NotificationHandle.of(adsDeviceNotificationBuffer);
                            SampleSize sampleSize = SampleSize.of(adsDeviceNotificationBuffer);
                            if (sampleSize.getAsLong() > Integer.MAX_VALUE) {
                                throw new IllegalStateException("Overflow in datalength: " + sampleSize.getAsLong());
                            }
                            // TODO: do we need a special marker class for: Notice: If your handle becomes invalid, one notification without data will be send once as advice.
                            Data data = Data.of(adsDeviceNotificationBuffer.readBytes((int) sampleSize.getAsLong()).array());
                            AdsNotificationSample adsNotificationSample = AdsNotificationSample.of(notificationHandle, sampleSize, data);
                            adsNotificationSamples.add(adsNotificationSample);

                        }
                        AdsStampHeader adsStampHeader = AdsStampHeader.of(timeStamp, samples, adsNotificationSamples);
                        adsStampHeaders.add(adsStampHeader);
                    }
                    out.add(new ADSDeviceNotificationRequest(amstcpHeader, amsHeader, length, stamps, adsStampHeaders));
                } else {
                    out.add(new UnknownCommand(amstcpHeader, amsHeader, commandBuffer));
                }
                break;
            case ADS_Read_Write:
                if (stateId.isRequest()) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer);
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer);
                    ReadLength readLength = ReadLength.of(commandBuffer);
                    if (readLength.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + readLength.getAsLong());
                    }
                    WriteLength writeLength = WriteLength.of(commandBuffer);
                    if (writeLength.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + writeLength.getAsLong());
                    }
                    if (readLength.getAsLong() + writeLength.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + readLength.getAsLong() + writeLength.getAsLong());
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) (readLength.getAsLong() + writeLength.getAsLong())).array());
                    out.add(new ADSReadWriteRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, readLength, writeLength, data));
                } else {
                    Result result = Result.of(commandBuffer);
                    Length length = Length.of(commandBuffer);
                    if (length.getAsLong() > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) length.getAsLong()).array());
                    out.add(new ADSReadWriteResponse(amstcpHeader, amsHeader, result, length, data));
                }
                break;
            case UNKNOWN:
                out.add(new UnknownCommand(amstcpHeader, amsHeader, commandBuffer));
        }
        if (commandBuffer.readableBytes() > 0) {
            throw new IllegalStateException("Unread bytes left: " + commandBuffer.readableBytes());
        }
    }

}
