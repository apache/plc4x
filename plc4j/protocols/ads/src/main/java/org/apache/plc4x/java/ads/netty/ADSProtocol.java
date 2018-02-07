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

import java.util.*;

public class ADSProtocol extends MessageToMessageCodec<ByteBuf, AMSTCPPaket> {

    // TODO: better track this a layer above as this here might be useless
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
        AMSNetId targetAmsNetId = AMSNetId.of(byteBuf.readBytes(6).array());
        AMSPort targetAmsPort = AMSPort.of(byteBuf.readUnsignedShortLE());
        AMSNetId sourceAmsNetId = AMSNetId.of(byteBuf.readBytes(6).array());
        AMSPort sourceAmsPort = AMSPort.of(byteBuf.readUnsignedShortLE());
        // TODO: could be transformed to readUnsignedShortLE someday
        Command commandId = Command.of(byteBuf.readBytes(2).array());
        // TODO: could be transformed to readUnsignedShortLE someday
        State stateId = State.of(byteBuf.readBytes(2).array());
        long dataLengthLong = byteBuf.readUnsignedIntLE();
        DataLength dataLength = DataLength.of(dataLengthLong);
        AMSError errorCode = AMSError.of(byteBuf.readBytes(4).array());
        Invoke invoke = Invoke.of(byteBuf.readBytes(4).array());
        AMSTCPPaket correlatedAmstcpPaket = requests.get(invoke);
        if (dataLengthLong > Integer.MAX_VALUE) {
            throw new IllegalStateException("Overflow in datalength: " + dataLengthLong);
        }
        ByteBuf commandBuffer = byteBuf.readBytes((int) dataLengthLong);
        boolean request = false;
        AMSTCPHeader amstcpHeader = AMSTCPHeader.of(packetLength);
        AMSHeader amsHeader = AMSHeader.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, commandId, stateId, dataLength, errorCode, invoke);
        switch (stateId) {
            case ADS_REQUEST_TCP:
                request = true;
                break;
            case ADS_RESPONSE_TCP:
                request = false;
                break;
            case ADS_REQUEST_UDP:
                request = true;
                break;
            case ADS_RESPONSE_UDP:
                request = false;
                break;
            case UNKNOWN:
                out.add(new UnknownCommand(amstcpHeader, amsHeader, commandBuffer));
                return;
        }
        switch (commandId) {
            case Invalid:
                out.add(new UnknownCommand(amstcpHeader, amsHeader, commandBuffer));
                break;
            case ADS_Read_Device_Info:
                if (request) {
                    out.add(new ADSReadDeviceInfoRequest(amstcpHeader, amsHeader));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    MajorVersion majorVersion = MajorVersion.of(commandBuffer.readByte());
                    MinorVersion minorVersion = MinorVersion.of(commandBuffer.readByte());
                    // TODO: could be transformed to readUnsignedShortLE someday
                    Version version = Version.of(commandBuffer.readBytes(2).array());
                    Device device = Device.of(commandBuffer.readBytes(16).array());
                    out.add(new ADSReadDeviceInfoResponse(amstcpHeader, amsHeader, result, majorVersion, minorVersion, version, device));
                }
                break;
            case ADS_Read:
                if (request) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer.readBytes(4).array());
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer.readBytes(4).array());
                    Length length = Length.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSReadRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, length));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    long adsReadLength = byteBuf.readUnsignedIntLE();
                    Length length = Length.of(adsReadLength);
                    if (adsReadLength > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + adsReadLength);
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) adsReadLength).array());
                    out.add(new ADSReadResponse(amstcpHeader, amsHeader, result, length, data));
                }
                break;
            case ADS_Write:
                if (request) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer.readBytes(4).array());
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer.readBytes(4).array());
                    long adsWriteLength = byteBuf.readUnsignedIntLE();
                    Length length = Length.of(adsWriteLength);
                    if (adsWriteLength > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + adsWriteLength);
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) adsWriteLength).array());
                    out.add(new ADSWriteRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, length, data));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSWriteResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Read_State:
                if (request) {
                    out.add(new ADSReadStateRequest(amstcpHeader, amsHeader));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSReadStateResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Write_Control:
                if (request) {
                    // TODO: could be transformed to readUnsignedShortLE someday
                    ADSState adsState = ADSState.of(commandBuffer.readBytes(2).array());
                    // TODO: could be transformed to readUnsignedShortLE someday
                    DeviceState deviceState = DeviceState.of(commandBuffer.readBytes(2).array());
                    long adsWriteControlLength = byteBuf.readUnsignedIntLE();
                    Length length = Length.of(adsWriteControlLength);
                    if (adsWriteControlLength > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + adsWriteControlLength);
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) adsWriteControlLength).array());
                    out.add(new ADSWriteControlRequest(amstcpHeader, amsHeader, adsState, deviceState, length, data));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSWriteControlResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Add_Device_Notification:
                if (request) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer.readBytes(4).array());
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer.readBytes(4).array());
                    Length length = Length.of(commandBuffer.readUnsignedIntLE());
                    TransmissionMode transmissionMode = TransmissionMode.of(commandBuffer.readBytes(4).array());
                    MaxDelay maxDelay = MaxDelay.of(commandBuffer.readBytes(4).array());
                    CycleTime cycleTime = CycleTime.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSAddDeviceNotificationRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSAddDeviceNotificationResponse(amstcpHeader, amsHeader, result, notificationHandle));
                }
                break;
            case ADS_Delete_Device_Notification:
                if (request) {
                    NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSDeleteDeviceNotificationRequest(amstcpHeader, amsHeader, notificationHandle));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    out.add(new ADSDeleteDeviceNotificationResponse(amstcpHeader, amsHeader, result));
                }
                break;
            case ADS_Device_Notification:
                if (request) {
                    long adsDeviceNotificationLength = commandBuffer.readUnsignedIntLE();
                    Length length = Length.of(adsDeviceNotificationLength);
                    if (adsDeviceNotificationLength > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + adsDeviceNotificationLength);
                    }
                    long numberOfStamps = commandBuffer.readUnsignedIntLE();
                    Stamps stamps = Stamps.of(numberOfStamps);
                    if (numberOfStamps > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + numberOfStamps);
                    }
                    ByteBuf adsDeviceNotificationBuffer = commandBuffer.readBytes((int) adsDeviceNotificationLength);
                    List<AdsStampHeader> adsStampHeaders = new ArrayList<>((int) numberOfStamps);
                    for (int i = 1; i <= numberOfStamps; i++) {
                        TimeStamp timeStamp = TimeStamp.of(adsDeviceNotificationBuffer.readBytes(8).array());
                        long numberOfSamples = adsDeviceNotificationBuffer.readUnsignedIntLE();
                        Samples samples = Samples.of(numberOfSamples);

                        List<AdsNotificationSample> adsNotificationSamples = new LinkedList<>();
                        for (int j = 1; j <= numberOfSamples; j++) {
                            NotificationHandle notificationHandle = NotificationHandle.of(adsDeviceNotificationBuffer.readBytes(4).array());
                            long sampleSizeLong = adsDeviceNotificationBuffer.readUnsignedIntLE();
                            SampleSize sampleSize = SampleSize.of(sampleSizeLong);
                            if (sampleSizeLong > Integer.MAX_VALUE) {
                                throw new IllegalStateException("Overflow in datalength: " + sampleSizeLong);
                            }
                            Data data = Data.of(adsDeviceNotificationBuffer.readBytes((int) sampleSizeLong).array());
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
                if (request) {
                    IndexGroup indexGroup = IndexGroup.of(commandBuffer.readBytes(4).array());
                    IndexOffset indexOffset = IndexOffset.of(commandBuffer.readBytes(4).array());
                    long readLengthLong = commandBuffer.readUnsignedIntLE();
                    ReadLength readLength = ReadLength.of(readLengthLong);
                    if (readLengthLong > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + readLengthLong);
                    }
                    long writeLengthLong = commandBuffer.readUnsignedIntLE();
                    WriteLength writeLength = WriteLength.of(writeLengthLong);
                    if (writeLengthLong > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + writeLengthLong);
                    }
                    if (readLengthLong + writeLengthLong > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + readLengthLong + writeLengthLong);
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) (readLengthLong + writeLengthLong)).array());
                    out.add(new ADSReadWriteRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, readLength, writeLength, data));
                } else {
                    Result result = Result.of(commandBuffer.readBytes(4).array());
                    long adsReadLength = byteBuf.readUnsignedIntLE();
                    Length length = Length.of(adsReadLength);
                    if (adsReadLength > Integer.MAX_VALUE) {
                        throw new IllegalStateException("Overflow in datalength: " + adsReadLength);
                    }
                    Data data = Data.of(commandBuffer.readBytes((int) adsReadLength).array());
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
