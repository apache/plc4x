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
import org.apache.commons.io.IOUtils;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.mock.connection.TestChannelFactory;
import org.junit.jupiter.api.TestInfo;
import org.pcap4j.core.*;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.namednumber.DataLinkType;

import java.io.EOFException;
import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.Objects;
import java.util.concurrent.TimeoutException;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

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
        super.sendChannelCreatedEvent();

        ByteBuf writtenData = channel.readOutbound();
        if(writtenData == null) {
            throw new PlcRuntimeException("Error reading initial channel output");
        }
        byte[] connectionRequest = new byte[writtenData.readableBytes()];
        writtenData.readBytes(connectionRequest);
        // TODO: Check the content of the Iso TP connection request.

        // Send an Iso TP connection response back to the pipeline.
        byte[] connectionConfirm = readPcapFile("org/apache/plc4x/java/s7/connection/iso-tp-connect-response.pcap");
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

        // Read a S7 CPU Functions request.
        writtenData = channel.readOutbound();
        byte[] cpuFunctionsRequest = new byte[writtenData.readableBytes()];
        writtenData.readBytes(cpuFunctionsRequest);
        // TODO: Check the content of the S7 Setup Communication connection request.

        // Send an S7 CPU Functions response back to the pipeline.
        byte[] cpuFunctionsResponse = readPcapFile(
            "org/apache/plc4x/java/s7/connection/s7-cpu-functions-response.pcap");
        // Override the type of reported S7 device.
        switch (getParamControllerType()) {
            case S7_1200:
                cpuFunctionsResponse[48] = '2';
                break;
            case S7_1500:
                cpuFunctionsResponse[48] = '5';
                break;
            case S7_300:
                cpuFunctionsResponse[48] = '3';
                break;
            case S7_400:
                cpuFunctionsResponse[48] = '4';
                break;
            default:
                cpuFunctionsResponse[48] = '1';
                break;
        }
        channel.writeInbound(Unpooled.wrappedBuffer(cpuFunctionsResponse));
    }

    public static byte[] readPcapFile(String filename) {
        try {
            InputStream in = S7PlcTestConnection.class.getClassLoader().getResourceAsStream(filename);
            byte[] pcap = IOUtils.toByteArray(in);
            byte[] data = new byte[pcap.length - 94];
            System.arraycopy(pcap, 94, data, 0, pcap.length - 94);
            return data;
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    public void sendPcapFile(String filename) {
        try {
            ClassLoader classLoader = getClass().getClassLoader();
            File file = new File(Objects.requireNonNull(classLoader.getResource(filename)).getFile());
            PcapHandle responsePcap = Pcaps.openOffline(file.getAbsolutePath());
            Packet packet = responsePcap.getNextPacketEx();
            byte[] data = packet.getPayload().getPayload().getPayload().getRawData();

            EmbeddedChannel channel = (EmbeddedChannel) getChannel();
            channel.writeInbound(Unpooled.wrappedBuffer(data));
        } catch (PcapNativeException | EOFException | TimeoutException | NotOpenException e) {
            throw new RuntimeException("Error sending pacap file " + filename, e);
        }
    }

    public void verifyPcapFile(String filename, TestInfo testInfo) {
        try {
            ClassLoader classLoader = getClass().getClassLoader();
            File file = new File(Objects.requireNonNull(classLoader.getResource(filename)).getFile());
            PcapHandle responsePcap = Pcaps.openOffline(file.getAbsolutePath());
            Packet packet = responsePcap.getNextPacketEx();
            byte[] refData = packet.getPayload().getPayload().getPayload().getRawData();

            // Get the systems output.
            EmbeddedChannel channel = (EmbeddedChannel) getChannel();
            ByteBuf request = channel.readOutbound();

            // Check the sizes are equal.
            assertThat(refData.length, equalTo(request.readableBytes()));

            // Read the raw data sent to the output.
            byte[] actData = new byte[request.readableBytes()];
            request.readBytes(actData);

            // Compare the actual output to the reference output
            if(!Arrays.equals(actData, refData)) {
                String currentWorkingDir = System.getProperty("user.dir");
                Class<?> testClass = testInfo.getTestClass().orElse(Object.class);
                Method testMethod = testInfo.getTestMethod().orElse(null);
                String fileName = currentWorkingDir + "/target/failsafe-reports/failure-" + testClass.getSimpleName() + "-"  + testMethod.getName() + ".pcap";
                try (PcapHandle handle = Pcaps.openDead(DataLinkType.EN10MB, 65536)) {
                    PcapDumper dumper = handle.dumpOpen(fileName);
                    dumper.dumpRaw(actData);
                    dumper.flush();
                }
                for (int i = 0; i < actData.length; i++) {
                    if (actData[i] != refData[i]) {
                        fail("Mismatch at position " + i);
                    }
                }
            }
        } catch (Exception e) {
            throw new RuntimeException("Error sending pacap file " + filename, e);
        }
    }

}