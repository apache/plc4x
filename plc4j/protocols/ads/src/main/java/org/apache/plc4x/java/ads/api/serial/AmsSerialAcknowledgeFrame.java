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
 * If an AMS serial frame has been received and the frame is OK (magic cookie OK, CRC OK, correct fragment number etc.),
 * then the receiver has to send an acknowledge frame, to inform the transmitter that the frame has arrived.
 *
 * @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec_amsframe.htm?id=8115637053270715044">TwinCAT AMS via RS232 Specification</a>
 */
public class AmsSerialAcknowledgeFrame implements ByteReadable {

    public static final int ID = 0x5A01;

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

    private AmsSerialAcknowledgeFrame(MagicCookie magicCookie, TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber, UserDataLength userDataLength, CRC crc) {
        this.magicCookie = magicCookie;
        this.transmitterAddress = transmitterAddress;
        this.receiverAddress = receiverAddress;
        this.fragmentNumber = fragmentNumber;
        this.userDataLength = userDataLength;
        this.crc = crc;
    }

    private AmsSerialAcknowledgeFrame(TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber) {
        this.magicCookie = MagicCookie.of(ID);
        this.transmitterAddress = transmitterAddress;
        this.receiverAddress = receiverAddress;
        this.fragmentNumber = fragmentNumber;
        this.userDataLength = UserDataLength.of((byte) 0);
        this.crc = CRC.of(DigestUtil.calculateCrc16(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength));
    }

    public static AmsSerialAcknowledgeFrame of(MagicCookie magicCookie, TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber, UserDataLength userDataLength, CRC crc) {
        return new AmsSerialAcknowledgeFrame(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
    }

    public static AmsSerialAcknowledgeFrame of(TransmitterAddress transmitterAddress, ReceiverAddress receiverAddress, FragmentNumber fragmentNumber) {
        return new AmsSerialAcknowledgeFrame(transmitterAddress, receiverAddress, fragmentNumber);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(magicCookie, transmitterAddress, receiverAddress, fragmentNumber, userDataLength, crc);
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

    public CRC getCrc() {
        return crc;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof AmsSerialAcknowledgeFrame)) {
            return false;
        }

        AmsSerialAcknowledgeFrame that = (AmsSerialAcknowledgeFrame) o;

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
        return "AmsSerialAcknowledgeFrame{" +
            "magicCookie=" + magicCookie +
            ", transmitterAddress=" + transmitterAddress +
            ", receiverAddress=" + receiverAddress +
            ", fragmentNumber=" + fragmentNumber +
            ", userDataLength=" + userDataLength +
            ", crc=" + crc +
            '}';
    }
}
