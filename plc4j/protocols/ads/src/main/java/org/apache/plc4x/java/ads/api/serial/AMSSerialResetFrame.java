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
package org.apache.plc4x.java.ads.api.serial;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.serial.types.*;
import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 * In case the transmitter does not receive a valid acknowledgement after multiple transmission, then a reset frame is
 * sent. In this way the receiver is informed that a new communication is running and the receiver then accepts the
 * fragment number during the next AMS-Frame, without carrying out a check.
 */
public class AMSSerialResetFrame implements ByteReadable {

    public static final int ID = 0xA503;

    /**
     * Id for detecting an AMS serial frame.
     */
    private final MagicCookie magicCookie;

    /**
     * Address of the sending participant. This value can always be set to 0 for an RS232 communication,
     * since it is a 1 to 1 connection and hence the participants are unique.
     */
    private final TransmitterAddress transmitterAddress;

    /**
     * Receiver’s address. This value can always be set to 0 for an RS232 communication, since it is a 1 to 1
     * connection and hence the participants are unique.
     */
    private final ReceiverAddress receiverAddress;

    /**
     * Number of the frame sent. Once the number 255 has been sent, it starts again from 0. The receiver checks this
     * number with an internal counter.
     */
    private final FragmentNumber fragmentNumber;

    /**
     * The max. length of the AMS packet to be sent is 255. If larger AMS packets are to be sent then they have to be
     * fragmented (not published at the moment).
     */
    private final UserDataLength userDataLength;

    private final CRC crc;

    private AMSSerialResetFrame(MagicCookie magicCookie, TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber, UserDataLength userDataLength, CRC crc) {
        this.magicCookie = magicCookie;
        this.transmitterAddress = transmitterAddress;
        this.receiverAddress = receiverAddress;
        this.fragmentNumber = fragmentNumber;
        this.userDataLength = userDataLength;
        this.crc = crc;
    }

    private AMSSerialResetFrame(TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber) {
        this.magicCookie = MagicCookie.of(ID);
        this.transmitterAddress = transmitterAddress;
        this.receiverAddress = receiverAddress;
        this.fragmentNumber = FragmentNumber.of((byte) 0);
        this.userDataLength = UserDataLength.of((byte) 0);
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
        byte[] digest = messageDigest.digest(userDataLength.getBytes());
        if (digest.length > 2) {
            throw new PlcRuntimeException("Digest length too great " + digest.length);
        }
        this.crc = CRC.of(digest[0], digest[1]);
    }

    public static AMSSerialResetFrame of(MagicCookie magicCookie, TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber, UserDataLength userDataLength, CRC crc) {
        return new AMSSerialResetFrame(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
    }

    public static AMSSerialResetFrame of(TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber) {
        return new AMSSerialResetFrame(transmitterAddress, receiverAddress, fragmentNumber);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof AMSSerialResetFrame)) return false;

        AMSSerialResetFrame that = (AMSSerialResetFrame) o;

        if (!magicCookie.equals(that.magicCookie)) return false;
        if (!transmitterAddress.equals(that.transmitterAddress)) return false;
        if (!receiverAddress.equals(that.receiverAddress)) return false;
        if (!fragmentNumber.equals(that.fragmentNumber)) return false;
        if (!userDataLength.equals(that.userDataLength)) return false;
        return crc.equals(that.crc);
    }

    @Override
    public int hashCode() {
        int result = magicCookie.hashCode();
        result = 31 * result + transmitterAddress.hashCode();
        result = 31 * result + receiverAddress.hashCode();
        result = 31 * result + fragmentNumber.hashCode();
        result = 31 * result + userDataLength.hashCode();
        result = 31 * result + crc.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AMSSerialResetFrame{" +
            "magicCookie=" + magicCookie +
            ", transmitterAddress=" + transmitterAddress +
            ", receiverAddress=" + receiverAddress +
            ", fragmentNumber=" + fragmentNumber +
            ", userDataLength=" + userDataLength +
            ", crc=" + crc +
            '}';
    }
}
