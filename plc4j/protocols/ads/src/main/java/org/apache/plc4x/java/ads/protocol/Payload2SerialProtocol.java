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
import org.apache.plc4x.java.ads.api.serial.AmsSerialAcknowledgeFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialResetFrame;
import org.apache.plc4x.java.ads.api.serial.types.*;
import org.apache.plc4x.java.ads.protocol.util.DigestUtil;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

import static java.lang.Integer.toHexString;
import static org.apache.commons.lang3.StringUtils.leftPad;

public class Payload2SerialProtocol extends MessageToMessageCodec<ByteBuf, ByteBuf> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Payload2TcpProtocol.class);

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, ByteBuf amsPacket, List<Object> out) throws Exception {
        out.add(AmsSerialFrame.of(FragmentNumber.of((byte) 0), UserData.of(amsPacket)).getByteBuf());
        // TODO: we need to remember the fragment and maybe even need to spilt up the package
        // TODO: if we exceed 255 byte
        amsPacket.release();
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
            case AmsSerialFrame.ID:
                AmsSerialFrame amsSerialFrame = AmsSerialFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData, crc);
                LOGGER.debug("Ams Serial Frame received {}", amsSerialFrame);
                out.add(userData.getByteBuf());
                break;
            case AmsSerialAcknowledgeFrame.ID:
                AmsSerialAcknowledgeFrame amsSerialAcknowledgeFrame = AmsSerialAcknowledgeFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial ACK Frame received {}", amsSerialAcknowledgeFrame);
                break;
            case AmsSerialResetFrame.ID:
                AmsSerialResetFrame amsSerialResetFrame = AmsSerialResetFrame.of(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
                LOGGER.debug("Ams Serial Reset Frame received {}", amsSerialResetFrame);
                break;
        }
        int calculatedCrc16 = DigestUtil.calculateCrc16(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData);
        if (!crc.equals(CRC.of(calculatedCrc16))) {
            throw new PlcProtocolException("CRC checksum wrong. Got "
                + "0x" + leftPad(toHexString(crc.getAsInt()), 4, "0")
                + " expected "
                + "0x" + leftPad(toHexString(calculatedCrc16), 4, "0")
            );
        }

        if (byteBuf.readableBytes() > 0) {
            byteBuf.release();
            throw new IllegalStateException("Unread bytes left: " + byteBuf.readableBytes());
        }

        byteBuf.release();
    }
}
