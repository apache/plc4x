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
import org.apache.plc4x.java.ads.protocol.util.DigestUtil;

/**
 * An AMS packet can be transferred via RS232 with the help of an AMS serial frame.
 * The actual AMS packet is in the user data field of the frame.
 * The max. length of the AMS packet is limited to 255 bytes.
 * Therefore the max. size of an AMS serial frame is 263 bytes.
 * The fragment number is compared with an internal counter by the receiver.
 * The frame number is simply accepted and not checked when receiving the first AMS frame or in case a timeout is
 * exceeded. The CRC16 algorithm is used for calculating the checksum.
 *
 * @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec_amsframe.htm?id=8115637053270715044">TwinCAT AMS via RS232 Specification</a>
 */
public class AmsSerialFrame implements ByteReadable {

    public static final int ID = 0xA501;

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

    /**
     * The AMS packet to be sent.
     */
    private final UserData userData;

    private final CRC crc;

    private AmsSerialFrame(MagicCookie magicCookie, TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber, UserDataLength userDataLength, UserData userData, CRC crc) {
        this.magicCookie = magicCookie;
        this.transmitterAddress = transmitterAddress;
        this.receiverAddress = receiverAddress;
        this.fragmentNumber = fragmentNumber;
        this.userDataLength = userDataLength;
        this.userData = userData;
        this.crc = crc;
    }

    private AmsSerialFrame(FragmentNumber fragmentNumber, UserData userData) {
        this.magicCookie = MagicCookie.of(ID);
        this.transmitterAddress = TransmitterAddress.RS232_COMM_ADDRESS;
        this.receiverAddress = ReceiverAddress.RS232_COMM_ADDRESS;
        this.fragmentNumber = fragmentNumber;
        long calculatedLength = userData.getCalculatedLength();
        if (calculatedLength > 255) {
            throw new IllegalArgumentException("Paket length must not exceed 255");
        }
        this.userDataLength = UserDataLength.of((byte) calculatedLength);
        byte[] amsPacketBytes = userData.getBytes();
        this.userData = UserData.of(amsPacketBytes);
        this.crc = CRC.of(DigestUtil.calculateCrc16(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData));
    }

    public static AmsSerialFrame of(MagicCookie magicCookie, TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber, UserDataLength userDataLength, UserData userData, CRC crc) {
        return new AmsSerialFrame(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData, crc);
    }

    public static AmsSerialFrame of(FragmentNumber fragmentNumber, UserData userData) {
        return new AmsSerialFrame(fragmentNumber, userData);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, userData, crc);
    }

    public MagicCookie getMagicCookie() {
        return magicCookie;
    }

    public TransmitterAddress getTransmitterAddress() {
        return transmitterAddress;
    }

    public ReceiverAddress getReceiverAddress() {
        return receiverAddress;
    }

    public FragmentNumber getFragmentNumber() {
        return fragmentNumber;
    }

    public UserDataLength getUserDataLength() {
        return userDataLength;
    }

    public UserData getUserData() {
        return userData;
    }

    public CRC getCrc() {
        return crc;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof AmsSerialFrame)) {
            return false;
        }

        AmsSerialFrame that = (AmsSerialFrame) o;

        if (!magicCookie.equals(that.magicCookie)) {
            return false;
        }
        if (!transmitterAddress.equals(that.transmitterAddress)) {
            return false;
        }
        if (!receiverAddress.equals(that.receiverAddress)) {
            return false;
        }
        if (!fragmentNumber.equals(that.fragmentNumber)) {
            return false;
        }
        if (!userDataLength.equals(that.userDataLength)) {
            return false;
        }
        if (!userData.equals(that.userData)) {
            return false;
        }
        return crc.equals(that.crc);
    }

    @Override
    public int hashCode() {
        int result = magicCookie.hashCode();
        result = 31 * result + transmitterAddress.hashCode();
        result = 31 * result + receiverAddress.hashCode();
        result = 31 * result + fragmentNumber.hashCode();
        result = 31 * result + userDataLength.hashCode();
        result = 31 * result + userData.hashCode();
        result = 31 * result + crc.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AmsSerialFrame{" +
            "magicCookie=" + magicCookie +
            ", transmitterAddress=" + transmitterAddress +
            ", receiverAddress=" + receiverAddress +
            ", fragmentNumber=" + fragmentNumber +
            ", userDataLength=" + userDataLength +
            ", userData=" + userData +
            ", crc=" + crc +
            '}';
    }
}
