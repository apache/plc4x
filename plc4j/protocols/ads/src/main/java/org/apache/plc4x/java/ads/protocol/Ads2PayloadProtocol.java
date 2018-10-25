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
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.plc4x.java.ads.api.commands.*;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.protocol.exception.AdsProtocolOverflowException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

public class Ads2PayloadProtocol extends MessageToMessageCodec<ByteBuf, AmsPacket> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Ads2PayloadProtocol.class);

    private static final Configuration CONF = new SystemConfiguration();
    private static final long MAX_NUM_STAMPS = CONF.getLong("plc4x.ads2payloadprotocol.max_num_stamps", 512L);
    private static final long MAX_NUM_SAMPLES = CONF.getLong("plc4x.ads2payloadprotocol.max_num_samples", 1024L);
    private static final long ADS_READ_COMMAND_MAX_BYTES = CONF.getLong("plc4x.ads2payloadprotocol.ads_read_command_max_bytes", 134217728L);
    private static final long ADS_WRITE_COMMAND_MAX_BYTES = CONF.getLong("plc4x.ads2payloadprotocol.ads_write_command_max_bytes", 134217728L);
    private static final long ADS_WRITE_CONTROL_COMMAND_MAX_BYTES = CONF.getLong("plc4x.ads2payloadprotocol.ads_write_control_command_max_bytes", 134217728L);
    private static final long ADS_NOTIFICATION_SAMPLE_MAX_BYTES = CONF.getLong("plc4x.ads2payloadprotocol.ads_notification_sample_max_bytes", 134217728L);
    private static final long ADS_READ_WRITE_COMMAND_REQUEST_MAX_BYTES = CONF.getLong("plc4x.ads2payloadprotocol.ads_read_write_command_request_max_bytes", 134217728L);
    private static final long ADS_READ_WRITE_COMMAND_RESPONSE_MAX_BYTES = CONF.getLong("plc4x.ads2payloadprotocol.ads_read_write_command_response_max_bytes", 134217728L);

    private final ConcurrentMap<Invoke, AmsPacket> requests;

    public Ads2PayloadProtocol() {
        this.requests = new ConcurrentHashMap<>();
    }

    /**
     * Resets this protocol and discard all send requests.
     */
    public void reset() {
        requests.clear();
    }

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, AmsPacket amsPacket, List<Object> out) {
        LOGGER.trace("(<--OUT): {}, {}, {}", channelHandlerContext, amsPacket, out);
        Invoke invokeId = amsPacket.getAmsHeader().getInvokeId();
        if (invokeId != Invoke.NONE) {
            requests.put(invokeId, amsPacket);
        }
        out.add(amsPacket.getByteBuf());
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) {
        LOGGER.trace("(-->IN): {}, {}, {}", channelHandlerContext, byteBuf, out);
        AmsNetId targetAmsNetId = AmsNetId.of(byteBuf);
        AmsPort targetAmsPort = AmsPort.of(byteBuf);
        AmsNetId sourceAmsNetId = AmsNetId.of(byteBuf);
        AmsPort sourceAmsPort = AmsPort.of(byteBuf);
        Command commandId = Command.of(byteBuf);
        State stateId = State.of(byteBuf);
        DataLength dataLength = DataLength.of(byteBuf);
        AmsError errorCode = AmsError.of(byteBuf);
        Invoke invoke = Invoke.of(byteBuf);
        AmsPacket correlatedAmsPacket = requests.remove(invoke);
        if (correlatedAmsPacket != null) {
            LOGGER.debug("Correlated packet received {}", correlatedAmsPacket);
        }
        if (dataLength.getAsLong() > Integer.MAX_VALUE) {
            throw new AdsProtocolOverflowException(Integer.class, dataLength.getAsLong());
        }
        ByteBuf commandBuffer = byteBuf.readBytes((int) dataLength.getAsLong());
        AmsHeader amsHeader = AmsHeader.of(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, commandId, stateId, dataLength, errorCode, invoke);
        final AmsPacket amsPacket;
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
        LOGGER.debug("Received amsPacket {}", amsPacket);
        out.add(amsPacket);
        if (commandBuffer.readableBytes() > 0) {
            throw new IllegalStateException("Unread bytes left: " + commandBuffer.readableBytes());
        }
    }


    private AmsPacket handleInvalidCommand(ByteBuf commandBuffer, AmsHeader amsHeader) {
        return UnknownCommand.of(amsHeader, commandBuffer);
    }

    private AmsPacket handleADSReadDeviceInfoCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            amsPacket = AdsReadDeviceInfoRequest.of(amsHeader);
        } else {
            Result result = Result.of(commandBuffer);
            MajorVersion majorVersion = MajorVersion.of(commandBuffer);
            MinorVersion minorVersion = MinorVersion.of(commandBuffer);
            Version version = Version.of(commandBuffer);
            Device device = Device.of(commandBuffer);
            amsPacket = AdsReadDeviceInfoResponse.of(amsHeader, result, majorVersion, minorVersion, version, device);
        }
        return amsPacket;
    }

    private AmsPacket handleADSReadCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            IndexGroup indexGroup = IndexGroup.of(commandBuffer);
            IndexOffset indexOffset = IndexOffset.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            amsPacket = AdsReadRequest.of(amsHeader, indexGroup, indexOffset, length);
        } else {
            Result result = Result.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, length.getAsLong());
            }
            if (length.getAsLong() > ADS_READ_COMMAND_MAX_BYTES) {
                throw new AdsProtocolOverflowException("ADS_READ_COMMAND_MAX_BYTES", ADS_READ_COMMAND_MAX_BYTES, length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = AdsReadResponse.of(amsHeader, result, length, data);
        }
        return amsPacket;
    }

    private AmsPacket handleADSWriteCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            IndexGroup indexGroup = IndexGroup.of(commandBuffer);
            IndexOffset indexOffset = IndexOffset.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, length.getAsLong());
            }
            if (length.getAsLong() > ADS_WRITE_COMMAND_MAX_BYTES) {
                throw new AdsProtocolOverflowException("ADS_WRITE_COMMAND_MAX_BYTES", ADS_WRITE_COMMAND_MAX_BYTES, length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = AdsWriteRequest.of(amsHeader, indexGroup, indexOffset, length, data);
        } else {
            Result result = Result.of(commandBuffer);
            amsPacket = AdsWriteResponse.of(amsHeader, result);
        }
        return amsPacket;
    }


    private AmsPacket handleADSReadStateCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            amsPacket = AdsReadStateRequest.of(amsHeader);
        } else {
            Result result = Result.of(commandBuffer);
            AdsState adsState = AdsState.of(commandBuffer);
            DeviceState deviceState = DeviceState.of(commandBuffer);
            amsPacket = AdsReadStateResponse.of(amsHeader, result, adsState, deviceState);
        }
        return amsPacket;
    }

    private AmsPacket handleADSWriteControlCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            AdsState adsState = AdsState.of(commandBuffer);
            DeviceState deviceState = DeviceState.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, length.getAsLong());
            }
            if (length.getAsLong() > ADS_WRITE_CONTROL_COMMAND_MAX_BYTES) {
                throw new AdsProtocolOverflowException("ADS_WRITE_CONTROL_COMMAND_MAX_BYTES", ADS_WRITE_CONTROL_COMMAND_MAX_BYTES, length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = AdsWriteControlRequest.of(amsHeader, adsState, deviceState, length, data);
        } else {
            Result result = Result.of(commandBuffer);
            amsPacket = AdsWriteControlResponse.of(amsHeader, result);
        }
        return amsPacket;
    }

    private AmsPacket handleADSAddDeviceNotificationCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            IndexGroup indexGroup = IndexGroup.of(commandBuffer);
            IndexOffset indexOffset = IndexOffset.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            TransmissionMode transmissionMode = TransmissionMode.of(commandBuffer);
            MaxDelay maxDelay = MaxDelay.of(commandBuffer);
            CycleTime cycleTime = CycleTime.of(commandBuffer);
            commandBuffer.skipBytes(AdsAddDeviceNotificationRequest.Reserved.NUM_BYTES);
            amsPacket = AdsAddDeviceNotificationRequest.of(amsHeader, indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime);
        } else {
            Result result = Result.of(commandBuffer);
            NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer);
            amsPacket = AdsAddDeviceNotificationResponse.of(amsHeader, result, notificationHandle);
        }
        return amsPacket;
    }

    private AmsPacket handADSDeleteDeviceNotificationCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            NotificationHandle notificationHandle = NotificationHandle.of(commandBuffer);
            amsPacket = AdsDeleteDeviceNotificationRequest.of(amsHeader, notificationHandle);
        } else {
            Result result = Result.of(commandBuffer);
            amsPacket = AdsDeleteDeviceNotificationResponse.of(amsHeader, result);
        }
        return amsPacket;
    }

    private AmsPacket handleADSDeviceNotificationCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, length.getAsLong());
            }
            Stamps stamps = Stamps.of(commandBuffer);
            if (stamps.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, stamps.getAsLong());
            }
            // Note: the length includes the 4 Bytes of stamps which we read already so we substract.
            ByteBuf adsDeviceNotificationBuffer = commandBuffer.readBytes((int) length.getAsLong() - Stamps.NUM_BYTES);
            List<AdsStampHeader> adsStampHeaders = new ArrayList<>((int) stamps.getAsLong());
            if (stamps.getAsLong() > MAX_NUM_STAMPS) {
                throw new AdsProtocolOverflowException("MAX_NUM_STAMPS", MAX_NUM_STAMPS, length.getAsLong());
            }
            for (int i = 1; i <= stamps.getAsLong(); i++) {
                AdsStampHeader adsStampHeader = handleStampHeader(adsDeviceNotificationBuffer);
                adsStampHeaders.add(adsStampHeader);
            }
            amsPacket = AdsDeviceNotificationRequest.of(amsHeader, length, stamps, adsStampHeaders);
        } else {
            amsPacket = UnknownCommand.of(amsHeader, commandBuffer);
        }
        return amsPacket;
    }

    private AdsStampHeader handleStampHeader(ByteBuf adsDeviceNotificationBuffer) {
        TimeStamp timeStamp = TimeStamp.of(adsDeviceNotificationBuffer);
        Samples samples = Samples.of(adsDeviceNotificationBuffer);

        List<AdsNotificationSample> adsNotificationSamples = new LinkedList<>();
        if (samples.getAsLong() > MAX_NUM_SAMPLES) {
            throw new AdsProtocolOverflowException("MAX_NUM_SAMPLES", MAX_NUM_SAMPLES, samples.getAsLong());
        }
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
            throw new AdsProtocolOverflowException(Integer.class, sampleSize.getAsLong());
        }
        if (sampleSize.getAsLong() > ADS_NOTIFICATION_SAMPLE_MAX_BYTES) {
            throw new AdsProtocolOverflowException("ADS_NOTIFICATION_SAMPLE_MAX_BYTES", ADS_NOTIFICATION_SAMPLE_MAX_BYTES, sampleSize.getAsLong());
        }
        // TODO: do we need a special marker class for: Notice: If your handle becomes invalid, one notification without data will be send once as advice.
        byte[] dataToRead = new byte[(int) sampleSize.getAsLong()];
        adsDeviceNotificationBuffer.readBytes(dataToRead);
        Data data = Data.of(dataToRead);
        return AdsNotificationSample.of(notificationHandle, sampleSize, data);
    }

    private AmsPacket handleADSReadWriteCommand(State stateId, ByteBuf commandBuffer, AmsHeader amsHeader) {
        AmsPacket amsPacket;
        if (stateId.isRequest()) {
            IndexGroup indexGroup = IndexGroup.of(commandBuffer);
            IndexOffset indexOffset = IndexOffset.of(commandBuffer);
            ReadLength readLength = ReadLength.of(commandBuffer);
            if (readLength.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, readLength.getAsLong());
            }
            WriteLength writeLength = WriteLength.of(commandBuffer);
            if (writeLength.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, writeLength.getAsLong());
            }
            if (readLength.getAsLong() + writeLength.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, readLength.getAsLong() + writeLength.getAsLong());
            }
            if (readLength.getAsLong() > ADS_READ_WRITE_COMMAND_REQUEST_MAX_BYTES) {
                throw new AdsProtocolOverflowException("ADS_READ_WRITE_COMMAND_REQUEST_MAX_BYTES", ADS_READ_WRITE_COMMAND_REQUEST_MAX_BYTES, readLength.getAsLong());
            }
            byte[] dataToRead = new byte[(int) readLength.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = AdsReadWriteRequest.of(amsHeader, indexGroup, indexOffset, readLength, writeLength, data);
        } else {
            Result result = Result.of(commandBuffer);
            Length length = Length.of(commandBuffer);
            if (length.getAsLong() > Integer.MAX_VALUE) {
                throw new AdsProtocolOverflowException(Integer.class, length.getAsLong());
            }
            if (length.getAsLong() > ADS_READ_WRITE_COMMAND_RESPONSE_MAX_BYTES) {
                throw new AdsProtocolOverflowException("ADS_READ_WRITE_COMMAND_RESPONSE_MAX_BYTES", ADS_READ_WRITE_COMMAND_RESPONSE_MAX_BYTES, length.getAsLong());
            }
            byte[] dataToRead = new byte[(int) length.getAsLong()];
            commandBuffer.readBytes(dataToRead);
            Data data = Data.of(dataToRead);
            amsPacket = AdsReadWriteResponse.of(amsHeader, result, length, data);
        }
        return amsPacket;
    }

    private AmsPacket handleUnknownCommand(ByteBuf commandBuffer, AmsHeader amsHeader) {
        return UnknownCommand.of(amsHeader, commandBuffer);
    }
}
