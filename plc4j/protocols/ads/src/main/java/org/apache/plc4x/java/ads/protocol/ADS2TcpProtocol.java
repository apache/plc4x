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
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSPacket;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.api.tcp.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.tcp.types.TcpLength;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

public class ADS2TcpProtocol extends MessageToMessageCodec<ByteBuf, AMSPacket> {

    private static final Logger LOGGER = LoggerFactory.getLogger(ADS2TcpProtocol.class);

    private final ConcurrentMap<Invoke, AMSPacket> requests;

    private final boolean ignoreBrokenPackages;

    public ADS2TcpProtocol() {
        this(false);
    }

    public ADS2TcpProtocol(boolean ignoreBrokenPackages) {
        this.requests = new ConcurrentHashMap<>();
        this.ignoreBrokenPackages = ignoreBrokenPackages;
    }

    /**
     * Resets this protocol and discard all send requests.
     */
    public void reset() {
        requests.clear();
    }

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, AMSPacket amsPacket, List<Object> out) throws Exception {
        Invoke invokeId = amsPacket.getAmsHeader().getInvokeId();
        if (invokeId != Invoke.NONE) {
            requests.put(invokeId, amsPacket);
        }
        out.add(amsPacket.toAmstcpPacket().getByteBuf());
    }

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) throws Exception {
        // Tcp decoding
        // Reserved
        byteBuf.skipBytes(AMSTCPHeader.Reserved.NUM_BYTES);
        TcpLength packetLength = TcpLength.of(byteBuf);
        AMSTCPHeader amstcpHeader = AMSTCPHeader.of(packetLength);
        LOGGER.debug("AMS TCP Header {}", amstcpHeader);

        // Ams decoding
        AMSNetId targetAmsNetId = AMSNetId.of(byteBuf);
        AMSPort targetAmsPort = AMSPort.of(byteBuf);
        AMSNetId sourceAmsNetId = AMSNetId.of(byteBuf);
        AMSPort sourceAmsPort = AMSPort.of(byteBuf);
        Command commandId = Command.of(byteBuf);
        State stateId = State.of(byteBuf);
        DataLength dataLength = DataLength.of(byteBuf);
        AMSError errorCode = AMSError.of(byteBuf);
        Invoke invoke = Invoke.of(byteBuf);
        AMSPacket correlatedamsPacket = requests.remove(invoke);
        if (correlatedamsPacket != null) {
            LOGGER.debug("Correlated packet received {}", correlatedamsPacket);
        }
        if (dataLength.getAsLong() > Integer.MAX_VALUE) {
            byteBuf.release();
            throw new IllegalStateException("Overflow in datalength: " + dataLength.getAsLong());
        }
        ByteBuf commandBuffer = byteBuf.readBytes((int) dataLength.getAsLong());
        AMSHeader amsHeader = AMSHeader.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, commandId, stateId, dataLength, errorCode, invoke);
        final AMSPacket amsPacket;
        switch (commandId) {
            case INVALID:
                amsPacket = handleInvalidCommand(commandBuffer, amsHeader);
                break;
            case ADS_READ_DEVICE_INFO:
                amsPacket = handleADSReadDeviceInfoCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_READ:
                amsPacket = handleADSReadCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_WRITE:
                amsPacket = handleADSWriteCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_READ_STATE:
                amsPacket = handleADSReadStateCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_WRITE_CONTROL:
                amsPacket = handleADSWriteControlCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_ADD_DEVICE_NOTIFICATION:
                amsPacket = handleADSAddDeviceNotificationCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_DELETE_DEVICE_NOTIFICATION:
                amsPacket = handADSDeleteDeviceNotificationCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_DEVICE_NOTIFICATION:
                amsPacket = handleADSDeviceNotificationCommand(stateId, commandBuffer, amsHeader);
                break;
            case ADS_READ_WRITE:
                amsPacket = handleADSReadWriteCommand(stateId, commandBuffer, amsHeader);
                break;
            case UNKNOWN:
            default:
                amsPacket = handleUnknownCommand(commandBuffer, amsHeader);
        }
        out.add(amsPacket);
        LOGGER.trace("Set amsPacket {} to out", amsPacket);
        if (!ignoreBrokenPackages && commandBuffer.readableBytes() > 0) {
            commandBuffer.release();
            byteBuf.release();
            throw new IllegalStateException("Unread bytes left: " + commandBuffer.readableBytes());
        }
        commandBuffer.release();
        byteBuf.release();
    }

    private AMSPacket handleInvalidCommand(ByteBuf commandBuffer, AMSHeader amsHeader) {
        return UnknownCommand.of(amsHeader, commandBuffer);
    }

    private AMSPacket handleADSReadDeviceInfoCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
        if (stateId.isRequest()) {
            amsPacket = ADSReadDeviceInfoRequest.of(amsHeader);
        } else {
            Result result = Result.of(commandBuffer);
            MajorVersion majorVersion = MajorVersion.of(commandBuffer);
            MinorVersion minorVersion = MinorVersion.of(commandBuffer);
            Version version = Version.of(commandBuffer);
            Device device = Device.of(commandBuffer);
            amsPacket = ADSReadDeviceInfoResponse.of(amsHeader, result, majorVersion, minorVersion, version, device);
        }
        return amsPacket;
    }

    private AMSPacket handleADSReadCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
        if (stateId.isRequest()) {
            IndexGroup indexGroup = IndexGroup.of(commandBuffer);
            IndexOffset indexOffset = IndexOffset.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            amsPacket = ADSReadRequest.of(amsHeader, indexGroup, indexOffset, length);
        } else {
            Result result = Result.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = ADSReadResponse.of(amsHeader, result, length, data);
        }
        return amsPacket;
    }

    private AMSPacket handleADSWriteCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
        if (stateId.isRequest()) {
            IndexGroup indexGroup = IndexGroup.of(commandBuffer);
            IndexOffset indexOffset = IndexOffset.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = ADSWriteRequest.of(amsHeader, indexGroup, indexOffset, length, data);
        } else {
            Result result = Result.of(commandBuffer);
            amsPacket = ADSWriteResponse.of(amsHeader, result);
        }
        return amsPacket;
    }


    private AMSPacket handleADSReadStateCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
        if (stateId.isRequest()) {
            amsPacket = ADSReadStateRequest.of(amsHeader);
        } else {
            Result result = Result.of(commandBuffer);
            amsPacket = ADSReadStateResponse.of(amsHeader, result);
        }
        return amsPacket;
    }

    private AMSPacket handleADSWriteControlCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
        if (stateId.isRequest()) {
            ADSState adsState = ADSState.of(commandBuffer);
            DeviceState deviceState = DeviceState.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = ADSWriteControlRequest.of(amsHeader, adsState, deviceState, length, data);
        } else {
            Result result = Result.of(commandBuffer);
            amsPacket = ADSWriteControlResponse.of(amsHeader, result);
        }
        return amsPacket;
    }

    private AMSPacket handleADSAddDeviceNotificationCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
        if (stateId.isRequest()) {
            IndexGroup indexGroup = IndexGroup.of(commandBuffer);
            IndexOffset indexOffset = IndexOffset.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            TransmissionMode transmissionMode = TransmissionMode.of(commandBuffer);
            MaxDelay maxDelay = MaxDelay.of(commandBuffer);
            CycleTime cycleTime = CycleTime.of(commandBuffer);
            commandBuffer.skipBytes(ADSAddDeviceNotificationRequest.Reserved.NUM_BYTES);
            amsPacket = ADSAddDeviceNotificationRequest.of(amsHeader, indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime);
        } else {
            Result result = Result.of(commandBuffer);
            NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer);
            amsPacket = ADSAddDeviceNotificationResponse.of(amsHeader, result, notificationHandle);
        }
        return amsPacket;
    }

    private AMSPacket handADSDeleteDeviceNotificationCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
        if (stateId.isRequest()) {
            NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer);
            amsPacket = ADSDeleteDeviceNotificationRequest.of(amsHeader, notificationHandle);
        } else {
            Result result = Result.of(commandBuffer);
            amsPacket = ADSDeleteDeviceNotificationResponse.of(amsHeader, result);
        }
        return amsPacket;
    }

    private AMSPacket handleADSDeviceNotificationCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
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
                AdsStampHeader adsStampHeader = handleStampHeader(adsDeviceNotificationBuffer);
                adsStampHeaders.add(adsStampHeader);
            }
            adsDeviceNotificationBuffer.release();
            amsPacket = ADSDeviceNotificationRequest.of(amsHeader, length, stamps, adsStampHeaders);
        } else {
            amsPacket = UnknownCommand.of(amsHeader, commandBuffer);
        }
        return amsPacket;
    }

    private AdsStampHeader handleStampHeader(ByteBuf adsDeviceNotificationBuffer) {
        TimeStamp timeStamp = TimeStamp.of(adsDeviceNotificationBuffer);
        Samples samples = Samples.of(adsDeviceNotificationBuffer);

        List<AdsNotificationSample> adsNotificationSamples = new LinkedList<>();
        for (int i = 1; i <= samples.getAsLong(); i++) {
            AdsNotificationSample adsNotificationSample = handleAdsNotificartionSample(adsDeviceNotificationBuffer);
            adsNotificationSamples.add(adsNotificationSample);

        }
        return AdsStampHeader.of(timeStamp, samples, adsNotificationSamples);
    }

    private AdsNotificationSample handleAdsNotificartionSample(ByteBuf adsDeviceNotificationBuffer) {
        NotificationHandle notificationHandle = NotificationHandle.of(adsDeviceNotificationBuffer);
        SampleSize sampleSize = SampleSize.of(adsDeviceNotificationBuffer);
        if (sampleSize.getAsLong() > Integer.MAX_VALUE) {
            throw new IllegalStateException("Overflow in datalength: " + sampleSize.getAsLong());
        }
        // TODO: do we need a special marker class for: Notice: If your handle becomes invalid, one notification without data will be send once as advice.
        byte[] dataToRead = new byte[(int) sampleSize.getAsLong()];
        adsDeviceNotificationBuffer.readBytes(dataToRead);
        Data data = Data.of(dataToRead);
        return AdsNotificationSample.of(notificationHandle, sampleSize, data);
    }

    private AMSPacket handleADSReadWriteCommand(State stateId, ByteBuf commandBuffer, AMSHeader amsHeader) {
        AMSPacket amsPacket;
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
            byte[] dataToRead = new byte[(int) readLength.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = ADSReadWriteRequest.of(amsHeader, indexGroup, indexOffset, readLength, writeLength, data);
        } else {
            Result result = Result.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new IllegalStateException("Overflow in datalength: " + length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = ADSReadWriteResponse.of(amsHeader, result, length, data);
        }
        return amsPacket;
    }

    private AMSPacket handleUnknownCommand(ByteBuf commandBuffer, AMSHeader amsHeader) {
        return UnknownCommand.of(amsHeader, commandBuffer);
    }
}
