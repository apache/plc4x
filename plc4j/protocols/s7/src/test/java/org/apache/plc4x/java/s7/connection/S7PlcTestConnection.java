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
package org.apache.plc4x.java.s7.connection;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.embedded.EmbeddedChannel;
import org.apache.commons.io.FileUtils;
import org.apache.commons.io.IOUtils;
import org.apache.plc4x.java.base.connection.TestChannelFactory;
import org.apache.plc4x.java.netty.events.S7ConnectionEvent;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;

public class S7PlcTestConnection extends S7PlcConnection {

    public S7PlcTestConnection(int rack, int slot, String params) {
        super(new TestChannelFactory(), rack, slot, params);
    }

    /*
            byte[] setupCommunicationResponse = toByteArray(
            new int[] {
                // ISO on TCP packet
                0x03, 0x00,
                0x00, 0x1B,
                // ISO TP packet
                0x02,
                TpduCode.DATA.getCode(),
                0x80,
                S7Protocol.S7_PROTOCOL_MAGIC_NUMBER,
                MessageType.ACK_DATA.getCode(), 0x00, 0x00,
                0x00, 0x01,
                // Parameter Length
                0x00, 0x08,
                // Data Length
                0x00, 0x00,
                // Error codes
                0x00, 0x00,
                // Parameters:
                ParameterType.SETUP_COMMUNICATION.getCode(), 0x00, 0x00, 0x08, 0x00, 0x08, 0x01, 0x00
            });
        dumpArrayToPcapFile("s7-setup-communication-response.pcap", setupCommunicationResponse);
     */

    @Override
    protected void sendChannelCreatedEvent() {
        EmbeddedChannel channel = (EmbeddedChannel) getChannel();

        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new S7ConnectionEvent());

        ByteBuf writtenData = channel.readOutbound();
        byte[] connectionRequest = new byte[writtenData.readableBytes()];
        writtenData.readBytes(connectionRequest);
        // TODO: Check the content of the Iso TP connection request.

        // Send an Iso TP connection response back to the pipeline.
        byte[] connectionConfirm = readPcapFile(
            "org/apache/plc4x/java/s7/connection/iso-tp-connect-response.pcap");
        channel.writeInbound(Unpooled.wrappedBuffer(connectionConfirm));

        // Read a S7 Setup Communication request.
        writtenData = channel.readOutbound();
        byte[] setupCommunicationRequest = new byte[writtenData.readableBytes()];
        writtenData.readBytes(setupCommunicationRequest);
        // TODO: Check the content of the S7 Setup Communication connection request.

        // Send an S7 Setup Communication response back to the pipeline.
        byte[] setupCommunicationResponse = readPcapFile(
            "org/apache/plc4x/java/s7/connection/s7-setup-communication-response.pcap");
        channel.writeInbound(Unpooled.wrappedBuffer(setupCommunicationResponse));
    }

    public static byte[] toByteArray(int[] in) {
        byte[] out = new byte[in.length];
        for(int i = 0; i < in.length; i++) {
            out[i] = (byte) in[i];
        }
        return out;
    }

    public static void dumpArrayToPcapFile(String filename, byte[] data) {
        byte[] pcapHeader = toByteArray(
            new int[]{
                // Magic Bytes (Big Endian)
                0xa1, 0xb2, 0xc3, 0xd4,
                // Version (2.4)
                0x00, 0x02, 0x00, 0x04,
                // Timestamp
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
                // Max Length
                0x00, 0x00, 0xFF, 0xFF,
                // Link Layer: Ethernet
                0x00, 0x00, 0x00, 0x01,
                // Timestamp (Seconds)
                0x00, 0x00, 0x00, 0x00,
                // Timestamp (Microseconds)
                0x00, 0x00, 0x00, 0x00,
                // Packet Length
                0x00, 0x00, 0x00, 0x00,
                // Packet Length (duplicate)
                0x00, 0x00, 0x00, 0x00, // 40
                /////////////////////////////////////
                // Ethernet Header
                // Destination Address (Mac Address)
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
                // Source Address (Mac Address)
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
                // IPv4
                0x08, 0x00,
                /////////////////////////////////////
                // IPv4 Header (type + 20 byte length)
                0x45, 0x00,
                // Length
                0x00, 0x00, // 40 + payload
                0x00, 0x00,
                0x00, 0x00,
                0x1e,
                // Protocol (TCP)
                0x06,
                // Header Checksum
                0x00, 0x00,
                // Source IP
                0x00, 0x00, 0x00, 0x00,
                // Dest IP
                0x00, 0x00, 0x00, 0x00, // 34
                /////////////////////////////////////
                // TCP Header
                // Source Port
                0x00, 0x66,
                // Dest Port
                0x10, 0xd1,
                0x00, 0x00, 0x00, 0x00,
                0x00, 0x00, 0x00, 0x00,
                // Header Length
                0x50, 0x18,
                0xFA, 0xDA,
                0x00, 0x00,
                0x00, 0x00 // 54
            });

        byte[] total = new byte[pcapHeader.length + data.length];
        System.arraycopy(pcapHeader, 0, total, 0, pcapHeader.length);
        System.arraycopy(data, 0, total, pcapHeader.length, data.length);
        total[35] = (byte) (data.length + 54);
        total[39] = (byte) (data.length + 54);
        total[39+18] = (byte) (data.length + 40);

        try {
            File file = new File(filename);
            System.out.println("Dumping to file: "  + file.getAbsolutePath());
            FileUtils.writeByteArrayToFile(file, total);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public byte[] readPcapFile(String filename) {
        try {
            InputStream in = getClass().getClassLoader().getResourceAsStream(filename);
            byte[] pcap = IOUtils.toByteArray(in);
            byte[] data = new byte[pcap.length - 94];
            System.arraycopy(pcap, 94, data, 0, pcap.length - 94);
            return data;
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

}
