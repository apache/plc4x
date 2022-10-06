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

import org.apache.plc4x.java.profinet.readwrite.IpAddress;
import org.apache.plc4x.java.profinet.readwrite.PnDcp_FrameId;

public class StaticHelper {

    public static int stringLength(String str) {
        if (str == null) {
            return 0;
        }
        return str.length();
    }

    public static int arrayLength(byte[] arr) {
        return arr.length;
    }

    public static int calculateIPv4Checksum(int totalLength, int identification, int timeToLive, IpAddress sourceAddress, IpAddress destinationAddress) {
        // https://en.wikipedia.org/wiki/Ones%27_complement
        // https://www.thegeekstuff.com/2012/05/ip-header-checksum/
        int[] words = new int[10];
        // Version and header length
        words[0] = 0x4500;
        words[1] = totalLength;
        words[2] = identification;
        // Flags and fragment offset
        words[3] = 0x0000;
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

        return cur;
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

}
