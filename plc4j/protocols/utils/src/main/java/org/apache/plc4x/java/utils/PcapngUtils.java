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
package org.apache.plc4x.java.utils;

import org.apache.commons.io.FileUtils;

import java.io.ByteArrayOutputStream;
import java.io.DataOutputStream;
import java.io.File;
import java.io.IOException;
import java.nio.ByteBuffer;

public class PcapngUtils {

    private PcapngUtils() {
      throw new IllegalStateException("Utility class!");
    }

    public static void dumpPacket(ByteBuffer data, int length, String name) throws IOException {
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        DataOutputStream dos = new DataOutputStream(out);
        byte[] pcapngHeader = {
            // PCAP header
            // Global Header:
            //  magic number
            (byte) 0xD4, (byte) 0xC3, (byte) 0xB2, (byte) 0xA1,
            //  Version(major / minor)
            (byte) 0x02, (byte) 0x00, (byte) 0x04, (byte) 0x00,
            //  Timezone
            (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            // 0
            (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            //  Snapshot length
            (byte) 0xFF, (byte) 0xFF, (byte) 0x00, (byte) 0x00,
            //  Network
            (byte) 0x01, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            // Packet header
            //  Timestamp (seconds)
            (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            //  Timestamp (microseconds)
            (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            // Packet length (in file)
            (byte) (length + 54), (byte) 0x00, (byte) 0x00, (byte) 0x00,
            // Packet length (real)
            (byte) (length + 54), (byte) 0x00, (byte) 0x00, (byte) 0x00,
            // Ethernet Frame Header
            // Dest MAC Address
            (byte) 0x00, (byte) 0x1B, (byte) 0x1B, (byte) 0x1A, (byte) 0xD6, (byte) 0xE0,
            // Source MAC Address
            (byte) 0xB8, (byte) 0x70, (byte) 0xF4, (byte) 0x6D, (byte) 0x2F, (byte) 0x58,
            // IPv4
            (byte) 0x08, (byte) 0x00,
            // IP Heaer
            // Version + Header Length
            (byte) 0x45,
            // Diferentiated Service Field
            (byte) 0x00,
            // Total Length (20 + 20 + length)
            (byte) 0x00, (byte) (length + 40),
            (byte) 0x6B, (byte) 0x6E, (byte) 0x40, (byte) 0x00,
            (byte) 0x40, (byte) 0x06, (byte) 0x44, (byte) 0xDF,
            (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            // TCP Header
            // Source Port
            (byte) 0xCC, (byte) 0xDE,
            // Dest Port
            (byte) 0x00, (byte) 0x66,
            // Sequence Number
            (byte) 0x8B, (byte) 0x9F, (byte) 0x0B, (byte) 0x76,
            // Ack Number
            (byte) 0x00, (byte) 0x1A, (byte) 0x4D, (byte) 0x96,
            // Header Length (first 6 bits, rest Flags)
            (byte) 0x50, (byte) 0x18,
            // Window Size
            (byte) 0x72, (byte) 0x10,
            // Checksum
            (byte) 0x1D, (byte) 0x43,
            // Urgent Pointer
            (byte) 0x00, (byte) 0x00
        };


        dos.write(pcapngHeader);
        dos.write(data.array(), 0, length);
        File output = new File(name);
        FileUtils.writeByteArrayToFile(output, out.toByteArray());
    }

}
