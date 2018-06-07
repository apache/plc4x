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
import org.apache.plc4x.java.ads.api.serial.AmsSerialAcknowledgeFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialResetFrame;
import org.apache.plc4x.java.ads.api.serial.types.*;
import org.apache.plc4x.java.ads.protocol.util.DigestUtil;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class Payload2SerialProtocol extends MessageToMessageCodec<ByteBuf, ByteBuf> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Payload2SerialProtocol.class);

    private final AtomicInteger fragmentCounter = new AtomicInteger(0);

    private final AtomicBoolean frameOnTheWay = new AtomicBoolean(false);

    private volatile ScheduledFuture<ChannelFuture> retryHandler;

    private final Lock lock = new ReentrantLock(true);

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, ByteBuf amsPacket, List<Object> out) throws Exception {
        LOGGER.trace("(<--OUT): {}, {}, {}", channelHandlerContext, amsPacket, out);
        while (frameOnTheWay.get() || !lock.tryLock()) {
            // In this case we might not send it yet.
            TimeUnit.MILLISECONDS.sleep(10);
        }
        int fragmentNumber = fragmentCounter.getAndIncrement();
        if (fragmentNumber > 255) {
            fragmentNumber = 0;
            fragmentCounter.set(fragmentNumber);
        }
        try {
            // TODO: we need to remember the fragment and maybe even need to spilt up the package
            // TODO: if we exceed 255 byte
            AmsSerialFrame amsSerialFrame = AmsSerialFrame.of(FragmentNumber.of((byte) fragmentNumber), UserData.of(amsPacket));
            out.add(amsSerialFrame.getByteBuf());
            retryHandler = channelHandlerContext.executor().schedule(() -> {
                try {
                    TimeUnit.SECONDS.sleep(2);
                    LOGGER.info("Retransmitting {}", amsSerialFrame);
                    channelHandlerContext.channel().writeAndFlush(amsSerialFrame.getByteBuf());
                } catch (InterruptedException e) {
                    LOGGER.debug("Interrupted", e);
                    Thread.currentThread().interrupt();
                }
                return channelHandlerContext.voidPromise();
            }, 0, TimeUnit.MILLISECONDS);
            frameOnTheWay.set(true);
        } finally {
            lock.unlock();
        }
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) throws Exception {
        LOGGER.trace("(-->IN): {}, {}, {}", channelHandlerContext, byteBuf, out);
        MagicCookie magicCookie = MagicCookie.of(byteBuf);
        TransmitterAddress transmitterAddress = TransmitterAddress.of(byteBuf);
        ReceiverAddress receiverAddress = ReceiverAddress.of(byteBuf);
        FragmentNumber fragmentNumber = FragmentNumber.of(byteBuf);
        UserDataLength userDataLength = UserDataLength.of(byteBuf);
        UserData userData;
        byte userDataLengthAsByte = userDataLength.getAsByte();
        if (userDataLengthAsByte > 0) {
            byte[] userDataByteArray = new byte[userDataLengthAsByte];
            byteBuf.readBytes(userDataByteArray);
            userData = UserData.of(userDataByteArray);
        } else {
            userData = UserData.EMPTY;
        }
        CRC crc = CRC.of(byteBuf);

        switch (magicCookie.getAsInt()) {
            case AmsSerialFrame.ID:
                AmsSerialFrame amsSerialFrame = AmsSerialFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData, crc);
                LOGGER.debug("Ams Serial Frame received {}", amsSerialFrame);
                // TODO: check if this is the right way to ack a package.
                ChannelFuture channelFuture = channelHandlerContext.writeAndFlush(AmsSerialAcknowledgeFrame.of(transmitterAddress, receiverAddress, fragmentNumber));
                // waiting for the ack-frame to be transmitted before we forward the package
                channelFuture.await();
                frameOnTheWay.set(false);
                out.add(userData.getByteBuf());
                break;
            case AmsSerialAcknowledgeFrame.ID:
                AmsSerialAcknowledgeFrame amsSerialAcknowledgeFrame = AmsSerialAcknowledgeFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial ACK Frame received {}", amsSerialAcknowledgeFrame);
                retryHandler.cancel(true);
                ReferenceCountUtil.release(byteBuf);
                break;
            case AmsSerialResetFrame.ID:
                // TODO: how to react to a reset
                AmsSerialResetFrame amsSerialResetFrame = AmsSerialResetFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial Reset Frame received {}", amsSerialResetFrame);
                ReferenceCountUtil.release(byteBuf);
                break;
        }
        CRC calculatedCrc = CRC.of(DigestUtil.calculateCrc16(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData));
        if (!crc.equals(calculatedCrc)) {
            throw new PlcProtocolException("CRC checksum wrong. Got " + crc + " expected " + calculatedCrc);
        }

        if (byteBuf.readableBytes() > 0) {
            throw new IllegalStateException("Unread bytes left: " + byteBuf.readableBytes());
        }

    }
}
