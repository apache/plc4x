/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.profinet.readwrite.utils;

import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;

import java.util.Arrays;
import java.util.List;

public class StaticHelper {

    public static int arrayLength(byte[] arr) {
        return arr.length;
    }

    public static short calculateIPv4Checksum(int totalLength, int identification, int timeToLive, IpAddress sourceAddress, IpAddress destinationAddress) {
        // https://en.wikipedia.org/wiki/Ones%27_complement
        // https://www.thegeekstuff.com/2012/05/ip-header-checksum/
        int[] words = new int[10];
        // Version and header length
        words[0] = 0x4500;
        words[1] = totalLength;
        words[2] = identification;
        // Flags and fragment offset
        words[3] = 0x4000;
        // Time to live and protocol
        words[4] = (timeToLive & 0xFF) << 8 | 0x11;
        // Checksum set to 0 for calculation
        words[5] = 0x0000;
        // Source address
        byte[] data = sourceAddress.getData();
        words[6] = ((((int) data[0]) & 0xFF) << 8) | ((int) data[1] & 0xFF);
        words[7] = ((((int) data[2]) & 0xFF) << 8) | ((int) data[3] & 0xFF);
        // Target address
        data = destinationAddress.getData();
        words[8] = ((((int) data[0]) & 0xFF) << 8) | ((int) data[1] & 0xFF);
        words[9] = ((((int) data[2]) & 0xFF) << 8) | ((int) data[3] & 0xFF);

        int cur = 0;
        for(int i = 0; i < 10; i++) {
            cur = cur + words[i];
            // The sum can result in max one bit above 0xFFFF.
            // Not sure if it could cascade in a second round, let's stay on the safe side.
            while(cur > 0xFFFF) {
                cur = cur & 0xFFFF;
                cur += 1;
            }
        }

        return (short) ~((short) cur);
    }

    public static short calculateUdpChecksum(IpAddress sourceAddress, IpAddress destinationAddress, int sourcePort, int destPort, int packetLength, DceRpc_Packet payload) {
        // https://en.wikipedia.org/wiki/Ones%27_complement
        // https://www.thegeekstuff.com/2012/05/ip-header-checksum/
        int[] words = new int[10];
        byte[] data = sourceAddress.getData();
        words[0] = ((((int) data[0]) & 0xFF) << 8) | ((int) data[1] & 0xFF);
        words[1] = ((((int) data[2]) & 0xFF) << 8) | ((int) data[3] & 0xFF);
        // Target address
        data = destinationAddress.getData();
        words[2] = ((((int) data[0]) & 0xFF) << 8) | ((int) data[1] & 0xFF);
        words[3] = ((((int) data[2]) & 0xFF) << 8) | ((int) data[3] & 0xFF);
        words[4] = 0x0011;
        words[5] = packetLength;
        words[6] = sourcePort;
        words[7] = destPort;
        words[8] = packetLength;
        words[9] = 0x0000;


        int cur = 0;
        for(int i = 0; i < 10; i++) {
            cur = cur + words[i];
            // The sum can result in max one bit above 0xFFFF.
            // Not sure if it could cascade in a second round, let's stay on the safe side.
            while(cur > 0xFFFF) {
                cur = cur & 0xFFFF;
                cur += 1;
            }
        }
        WriteBufferByteBased buffer;
        boolean evenSize = (payload.getLengthInBytes() % 2) == 0;
        int byteLen = evenSize ? payload.getLengthInBytes() : payload.getLengthInBytes() + 1;
        buffer = new WriteBufferByteBased(new byte[byteLen],
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));

        try {
            payload.serialize(buffer);
            if (!evenSize) {
                buffer.writeSignedByte(8, (byte) 0x00);
            }
            byte[] byteBuffer = buffer.getBytes();

            for(int i = 0; i < byteBuffer.length - 1; i += 2) {
                int w =  ((((int) byteBuffer[i]) & 0xFF) << 8) | ((int) byteBuffer[i+1] & 0xFF);
                cur = cur + w;
                // The sum can result in max one bit above 0xFFFF.
                // Not sure if it could cascade in a second round, let's stay on the safe side.

                while (cur > 0xFFFF) {
                    cur = cur & 0xFFFF;
                    cur += 1;
                }
            }
        } catch (BufferException e) {
            return 0x0000;
        }

        return (short) ~((short) cur);
    }

    public static void main(String[] args) {
        System.out.println(calculateIPv4Checksum(532, 0x44DF, 64,
            new IpAddress(new byte[]{(byte) 0xC0, (byte) 0xA8, (byte) 0x18, (byte) 0xC8}),
            new IpAddress(new byte[]{(byte) 0xC0, (byte) 0xA8, (byte) 0x18, (byte) 0x1F})));
        System.out.println(calculateIPv4Checksum(198, 0x0048, 30,
            new IpAddress(new byte[]{(byte) 0xC0, (byte) 0xA8, (byte) 0x18, (byte) 0x1F}),
            new IpAddress(new byte[]{(byte) 0xC0, (byte) 0xA8, (byte) 0x18, (byte) 0xC8})));
    }

    public static PnDcp_FrameId getFrameId(int frameIdValue) {
        // Range 1
        if(frameIdValue == 0x0020) {
            return PnDcp_FrameId.PTCP_RTSyncPDUWithFollowUp;
        }
        // Range 2
        if(frameIdValue == 0x0080) {
            return PnDcp_FrameId.PTCP_RTSyncPDU;
        }
        // Range 3
        if(frameIdValue >= 0x0100 && frameIdValue <= 0x0FFF) {
            return PnDcp_FrameId.RT_CLASS_3;
        }
        // (We do not support RT Class 3 (No need to implement these))
        // Range 4
        // (Not used)
        // Range 5
        // (Not used)
        // Range 6
        if(frameIdValue >= 8000 && frameIdValue <= 0xBFFF) {
            return PnDcp_FrameId.RT_CLASS_1;
        }
        // Range 7
        if(frameIdValue >= 0xC000 && frameIdValue <= 0xFBFF) {
            return PnDcp_FrameId.RT_CLASS_UDP;
        }
        // Range 8
        if(frameIdValue == 0xFC01) {
            return PnDcp_FrameId.Alarm_High;
        }
        if(frameIdValue == 0xFE01) {
            return PnDcp_FrameId.Alarm_Low;
        }
        if(frameIdValue == 0xFEFC) {
            return PnDcp_FrameId.DCP_Hello_ReqPDU;
        }
        if(frameIdValue == 0xFEFD) {
            return PnDcp_FrameId.DCP_GetSet_PDU;
        }
        if(frameIdValue == 0xFEFE) {
            return PnDcp_FrameId.DCP_Identify_ReqPDU;
        }
        if(frameIdValue == 0xFEFF) {
            return PnDcp_FrameId.DCP_Identify_ResPDU;
        }
        // Range 9
        if(frameIdValue == 0xFF00) {
            return PnDcp_FrameId.PTCP_AnnouncePDU;
        }
        if(frameIdValue == 0xFF20) {
            return PnDcp_FrameId.PTCP_FollowUpPDU;
        }
        if(frameIdValue == 0xFF40) {
            return PnDcp_FrameId.PTCP_DelayReqPDU;
        }
        if(frameIdValue == 0xFF41) {
            return PnDcp_FrameId.PTCP_DelayResPDUWithFollowUp;
        }
        if(frameIdValue == 0xFF42) {
            return PnDcp_FrameId.PTCP_DelayFuResPDUWithFollowUp;
        }
        if(frameIdValue == 0xFF43) {
            return PnDcp_FrameId.PTCP_DelayResPDUWithoutFollowUp;
        }
        // Range 12
        // 0xFF80 - 0xFF8F FragmentationFrameId
        if(frameIdValue >= 0xFF80 && frameIdValue <= 0xFF8F) {
            return PnDcp_FrameId.FragmentationFrameId;
        }

        return PnDcp_FrameId.RESERVED;
    }

    public static boolean isSysexEnd(ReadBuffer io) {
        // Peek two bytes; the LLDP sysex terminator is a 0x00 in the second
        // position. New SPI buffers don't expose absolute byte access so we
        // snapshot/rewind to peek non-destructively.
        if (!(io instanceof ReadBufferByteBased rb)) {
            return false;
        }
        int pos = rb.getPositionInBits();
        try {
            byte[] bytes = rb.readBits(16);
            return bytes.length >= 2 && bytes[1] == (byte) 0x00;
        } catch (BufferException e) {
            return false;
        } finally {
            rb.setPositionInBits(pos);
        }
    }

    public static LldpUnit parseSysexString(ReadBuffer io) {
        try {
            return LldpUnit.staticParse(io);
        } catch (BufferException e) {
            return null;
        }
    }

    public static void serializeSysexString(WriteBuffer io, LldpUnit unit) {
        try {
            unit.serialize(io);
        } catch (BufferException ignored) {
            // best-effort
        }
    }

    public static int lengthSysexString(List<LldpUnit> data) {
        int lengthInBytes = 0;
        for (LldpUnit unit : data) {
            lengthInBytes += unit.getLengthInBytes();
        }
        return lengthInBytes;
    }

    public static void writeDataUnit(WriteBuffer writeBuffer, PnIo_CyclicServiceDataUnit dataUnit) throws BufferException {
        dataUnit.serialize(writeBuffer);
    }

    public static PnIo_CyclicServiceDataUnit readDataUnit(ReadBuffer readBuffer) throws BufferException {
        // Subtract 4 bytes for the trailing footer (cycle counter + flags) we
        // don't want to consume as part of the data unit.
        int NO_TRAILING_BYTES = 4;
        if (!(readBuffer instanceof ReadBufferByteBased rb)) {
            return PnIo_CyclicServiceDataUnit.staticParse(readBuffer, (short) 0);
        }
        int remainingBytes = rb.getRemainingBits() / 8;
        int dataUnitLength = Math.max(remainingBytes - NO_TRAILING_BYTES, 0);
        return PnIo_CyclicServiceDataUnit.staticParse(readBuffer, (short) dataUnitLength);
    }

    private static final byte[] nullUuid = new byte[] {0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00};
    public static boolean isNullUuid(Uuid uuid) {
        return Arrays.equals(uuid.getData(), nullUuid);
    }

}
