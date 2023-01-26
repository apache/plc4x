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
package org.apache.plc4x.protocol.ads;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.ads.readwrite.AdsDataType;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Scanner {

    private static Logger logger = LoggerFactory.getLogger(Scanner.class);

    public static void main(String[] args) throws Exception {
        if (args.length != 5) {
            System.out.println("Usage: java -jar ... [ip -address] [target ams] [target ams port] [source ams] [source ams port]");
            System.out.println("All tags are required. AMS route must be created beforehand.");
            System.out.println("Example parameter sequence:");
            System.out.println("  192.168.2.251  39.42.54.209.1.1  851  192.168.2.232.1.1  851");
            return;
        }

        String targetIp = args[0];
        String targetAmsId = args[1];
        String targetAmsPort = args[2];

        String sourceAmsId = args[3];
        String sourceAmsPort = args[4];

        String connectionString = "ads:tcp://" + targetIp + "?targetAmsNetId=" + targetAmsId + "&targetAmsPort=" + targetAmsPort + "&sourceAmsNetId=" + sourceAmsId + "&sourceAmsPort=" + sourceAmsPort;
        System.out.println("Launching connection " + connectionString);

        // Establish a connection to the plc using the url provided as first argument
        try (PlcConnection plcConnection = new DefaultPlcDriverManager().getConnection(connectionString)) {
            PlcConnectionMetadata metadata = plcConnection.getMetadata();
            System.out.println("read: " + metadata.canRead());
            System.out.println("write: " + metadata.canWrite());

            // read symbols
            System.out.println("Reading symbol info");
            PlcReadRequest.Builder readRequestBuilder = plcConnection.readRequestBuilder();
            PlcReadRequest request = readRequestBuilder.addTagAddress("SYM_UPLOADINFO2", "0xf00f/0x0:SINT[24]").build();
            PlcReadResponse rsp = request.execute().get();
            ByteBuffer buffer = toBuffer(rsp, "SYM_UPLOADINFO2");

            //System.err.println("first answer " + Hex.dump(buffer.array()));
            int symbolAnswerSize = buffer.getInt(4);
            System.out.println("Expecting symbol table containing " + symbolAnswerSize + " bytes");

            request = plcConnection.readRequestBuilder().addTagAddress("SYM_UPLOAD", "0xf00b/0x0:SINT[" + symbolAnswerSize + "]").build();
            PlcReadResponse readResponse = request.execute().get();
            buffer = toBuffer(readResponse, "SYM_UPLOAD");
            //System.err.println("second answer " + Hex.dump(buffer.array()));

            System.out.println("      |           PLC4X Tag          |          Hex          |                                                            ");
            System.out.println("##### |          Query Syntax        |   Index   |   Offset  |   Type   | Name                                   | Size (B) | Type | Flag | Comment");
            System.out.println("------+------------------------------+-----------+-----------+----------+----------------------------------------+----------+------+------+");

            int index = 0;
            int pos = 0;

            List<String> supportedTypes = Arrays.stream(AdsDataType.values())
                .map(AdsDataType::name)
                .map(String::toUpperCase)
                .map(s -> s + '\0')
                .collect(Collectors.toList());

            while (buffer.remaining() > 0) {
                int sectionLen = buffer.getInt();

                int group = buffer.getInt();
                int offset = buffer.getInt();
                int symbolSize = buffer.getInt();
                int symbolType = buffer.getInt();
                int symbolFlags = buffer.getInt();
                short nameLength = (short) (buffer.getShort() + 1);
                short typeLength = (short) (buffer.getShort() + 1);
                short commentLength = (short) (buffer.getShort() + 1);

                String name= slice(buffer, nameLength);
                String type = slice(buffer, typeLength);
                String comment = slice(buffer, commentLength);

                if (supportedTypes.contains(type.toUpperCase())) {
                    System.out.printf("%5s |", index++);
                    System.out.printf("%30s |", "0x" + Integer.toHexString(group) + "/0x" + Integer.toHexString(offset) + ":" + type);
                    System.out.printf("%10s |", "0x" + Integer.toHexString(group));
                    System.out.printf("%10s |", "0x" + Integer.toHexString(offset));
                    System.out.printf("%10s |", type);
                    System.out.printf("%-40s |", name);
                    System.out.printf("%9s |", symbolSize);
                    System.out.printf("%5s |", symbolType);
                    System.out.printf("%5s |", symbolFlags);
                    System.out.println(comment);
                }

                pos += sectionLen;
                buffer.position(pos);

            }
        }
    }

    private static String slice(ByteBuffer buffer, short length) {
        byte[] arr = new byte[length];
        buffer.get(arr);
        return new String(arr, StandardCharsets.UTF_8);
    }

    private static ByteBuffer toBuffer(PlcReadResponse rsp, String tagName) {
        System.out.println(rsp.getTagNames() + " " + rsp.getTag(tagName) + " " + rsp.getResponseCode(tagName));
        List<PlcValue> symbols = (List<PlcValue>) rsp.getObject(tagName);
        ByteBuffer buffer = ByteBuffer.allocate(symbols.size()).order(ByteOrder.LITTLE_ENDIAN);

        for (PlcValue byteVal : symbols) {
            byte byteValue = byteVal.getByte();
            //System.err.println("data " + Hex.encodeHexString(new byte[] {byteValue}));
            buffer.put(byteValue);
        }
        buffer.rewind();
        return buffer;
    }

}
