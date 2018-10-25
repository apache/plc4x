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
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.util.ReferenceCountUtil;
import io.netty.util.concurrent.ScheduledFuture;
import org.apache.commons.lang3.mutable.MutableInt;
import org.apache.plc4x.java.ads.api.serial.AmsSerialAcknowledgeFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialResetFrame;
import org.apache.plc4x.java.ads.api.serial.types.*;
import org.apache.plc4x.java.ads.protocol.util.DigestUtil;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolPayloadTooBigException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

public class Payload2SerialProtocol extends MessageToMessageCodec<ByteBuf, ByteBuf> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Payload2SerialProtocol.class);

    private final AtomicInteger fragmentCounter = new AtomicInteger(0);

    private AtomicReference<ScheduledFuture<?>> currentRetryer = new AtomicReference<>();

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, ByteBuf amsPacket, List<Object> out) throws PlcProtocolPayloadTooBigException {
        LOGGER.trace("(<--OUT): {}, {}, {}", channelHandlerContext, amsPacket, out);
        int fragmentNumber = fragmentCounter.getAndUpdate(value -> value > 255 ? 0 : ++value);
        LOGGER.debug("Using fragmentNumber {} for {}", fragmentNumber, amsPacket);
        UserData userData = UserData.of(amsPacket);
        if (userData.getCalculatedLength() > 255) {
            throw new PlcProtocolPayloadTooBigException("ADS/AMS", 255, (int) userData.getCalculatedLength(), amsPacket);
        }
        AmsSerialFrame amsSerialFrame = AmsSerialFrame.of(FragmentNumber.of((byte) fragmentNumber), userData);

        MutableInt retryCount = new MutableInt(0);
        ScheduledFuture<?> oldRetryer = currentRetryer.get();
        if (oldRetryer != null) {
            oldRetryer.cancel(false);
        }
        currentRetryer.set(channelHandlerContext.executor().scheduleAtFixedRate(() -> {
            LOGGER.trace("Retrying {} the {} time", amsSerialFrame, retryCount);
            int currentTry = retryCount.incrementAndGet();
            if (currentTry > 10) {
                // TODO: we might need to throw an exception to potentially cancel upstream waiting
                channelHandlerContext.writeAndFlush(AmsSerialResetFrame.of(FragmentNumber.of((byte) fragmentNumber)));
                PlcRuntimeException plcRuntimeException = new PlcRuntimeException("Retry exhausted after " + retryCount + " times");
                channelHandlerContext.fireExceptionCaught(plcRuntimeException);
                throw plcRuntimeException;
            } else {
                channelHandlerContext.writeAndFlush(amsSerialFrame);
            }
        }, 100, 100, TimeUnit.MILLISECONDS));
        out.add(amsSerialFrame.getByteBuf());
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) throws Exception {
        LOGGER.trace("(-->IN): {}, {}, {}", channelHandlerContext, byteBuf, out);
        if (byteBuf.readableBytes() < MagicCookie.NUM_BYTES + TransmitterAddress.NUM_BYTES + ReceiverAddress.NUM_BYTES + FragmentNumber.NUM_BYTES) {
            return;
        }
        MagicCookie magicCookie = MagicCookie.of(byteBuf);
        TransmitterAddress transmitterAddress = TransmitterAddress.of(byteBuf);
        ReceiverAddress receiverAddress = ReceiverAddress.of(byteBuf);
        FragmentNumber fragmentNumber = FragmentNumber.of(byteBuf);
        int expectedFrameNumber = fragmentCounter.get() - 1;
        if (expectedFrameNumber < 0) {
            expectedFrameNumber = 255;
        }
        if (fragmentNumber.getAsByte() != expectedFrameNumber) {
            LOGGER.warn("Unexpected fragment {} received. Expected {}", fragmentNumber, expectedFrameNumber);
        }
        UserDataLength userDataLength = UserDataLength.of(byteBuf);
        UserData userData;
        byte userDataLengthAsByte = userDataLength.getAsByte();
        if (byteBuf.readableBytes() < userDataLengthAsByte) {
            return;
        }
        if (userDataLengthAsByte > 0) {
            byte[] userDataByteArray = new byte[userDataLengthAsByte];
            byteBuf.readBytes(userDataByteArray);
            userData = UserData.of(userDataByteArray);
        } else {
            userData = UserData.EMPTY;
        }
        CRC crc = CRC.of(byteBuf);

        // we don't need to retransmit
        ScheduledFuture<?> scheduledFuture = currentRetryer.get();
        if (scheduledFuture != null) {
            scheduledFuture.cancel(false);
        }

        Runnable postAction = null;
        switch (magicCookie.getAsInt()) {
            case AmsSerialFrame.ID:
                AmsSerialFrame amsSerialFrame = AmsSerialFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData, crc);
                LOGGER.debug("Ams Serial Frame received {}", amsSerialFrame);
                postAction = () -> {
                    // TODO: check if this is the right way to ack a package.
                    ChannelFuture channelFuture = channelHandlerContext.writeAndFlush(AmsSerialAcknowledgeFrame.of(transmitterAddress, receiverAddress, fragmentNumber).getByteBuf());
                    // waiting for the ack-frame to be transmitted before we forward the package
                    try {
                        channelFuture.await();
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        throw new PlcRuntimeException(e);
                    }
                    out.add(userData.getByteBuf());
                };
                break;
            case AmsSerialAcknowledgeFrame.ID:
                AmsSerialAcknowledgeFrame amsSerialAcknowledgeFrame = AmsSerialAcknowledgeFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial ACK Frame received {}", amsSerialAcknowledgeFrame);
                ReferenceCountUtil.release(byteBuf);
                break;
            case AmsSerialResetFrame.ID:
                // TODO: how to react to a reset
                AmsSerialResetFrame amsSerialResetFrame = AmsSerialResetFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial Reset Frame received {}", amsSerialResetFrame);
                ReferenceCountUtil.release(byteBuf);
                break;
            default:
                throw new PlcProtocolException("Unknown type: " + magicCookie);
        }
        CRC calculatedCrc = CRC.of(DigestUtil.calculateCrc16(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData));
        if (!crc.equals(calculatedCrc)) {
            throw new PlcProtocolException("CRC checksum wrong. Got " + crc + " expected " + calculatedCrc);
        }

        if (postAction != null) {
            postAction.run();
        }

        if (byteBuf.readableBytes() > 0) {
            throw new IllegalStateException("Unread bytes left: " + byteBuf.readableBytes());
        }
    }
}
