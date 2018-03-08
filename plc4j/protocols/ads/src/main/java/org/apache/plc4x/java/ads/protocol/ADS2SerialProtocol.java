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
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.ads.api.generic.AMSPacket;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.serial.AMSSerialAcknowledgeFrame;
import org.apache.plc4x.java.ads.api.serial.AMSSerialFrame;
import org.apache.plc4x.java.ads.api.serial.AMSSerialResetFrame;
import org.apache.plc4x.java.ads.api.serial.types.*;
import org.apache.plc4x.java.ads.api.tcp.AMSTCPHeader;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

public class ADS2SerialProtocol extends MessageToMessageCodec<ByteBuf, AMSPacket> {

    private static final Logger LOGGER = LoggerFactory.getLogger(ADS2TcpProtocol.class);

    private final ConcurrentMap<Invoke, AMSPacket> requests;

    private final ADS2TcpProtocol ads2TcpProtocol;

    public ADS2SerialProtocol() {
        this.requests = new ConcurrentHashMap<>();
        this.ads2TcpProtocol = new ADS2TcpProtocol(true);
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
        byte asLong = (byte) (invokeId.getAsLong() % 255);
        out.add(amsPacket.toAmsSerialFrame(asLong).getByteBuf());
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) throws Exception {
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
            case AMSSerialFrame.ID:
                // This is a lazy implementation. we just reuse the tcp implementation
                ByteBuf fakeTcpHeader = AMSTCPHeader.of(0).getByteBuf();
                ads2TcpProtocol.decode(channelHandlerContext, Unpooled.wrappedBuffer(fakeTcpHeader, userData.getByteBuf()), out);
                AMSPacket amsPacket = (AMSPacket) out.get(0);
                AMSSerialFrame amsSerialFrame = amsPacket.toAmsSerialFrame(fragmentNumber.getBytes()[0]);
                LOGGER.debug("Ams Serial Frame received {}", amsSerialFrame);
                break;
            case AMSSerialAcknowledgeFrame.ID:
                AMSSerialAcknowledgeFrame amsSerialAcknowledgeFrame = AMSSerialAcknowledgeFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial ACK Frame received {}", amsSerialAcknowledgeFrame);
                break;
            case AMSSerialResetFrame.ID:
                AMSSerialResetFrame amsSerialResetFrame = AMSSerialResetFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial Reset Frame received {}", amsSerialResetFrame);
                break;
        }

        MessageDigest messageDigest;
        try {
            messageDigest = MessageDigest.getInstance("CRC-16");
        } catch (NoSuchAlgorithmException e) {
            throw new PlcRuntimeException(e);
        }
        messageDigest.update(magicCookie.getBytes());
        messageDigest.update(transmitterAddress.getBytes());
        messageDigest.update(receiverAddress.getBytes());
        messageDigest.update(fragmentNumber.getBytes());
        messageDigest.update(userDataLength.getBytes());
        byte[] digest = messageDigest.digest(userData.getBytes());
        if (digest.length > 2) {
            throw new PlcRuntimeException("Digest length too great " + digest.length);
        }
        if (!Arrays.equals(digest, crc.getBytes())) {
            throw new PlcProtocolException("CRC checksum wrong");
        }

        byteBuf.release();
    }
}
