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
package org.apache.plc4x.java.modbus;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.base.util.HexUtil;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.Arrays;
import java.util.Collection;
import java.util.stream.IntStream;

public class ManualPlc4XModbusTest {

    public static void main(String... args) {
        String connectionUrl;
        if (args.length > 0 && "serial".equalsIgnoreCase(args[0])) {
            System.out.println("Using serial");
            connectionUrl = "modbus:serial:///dev/ttys003";
        } else {
            System.out.println("Using tcp");
            connectionUrl = "modbus:tcp://localhost:5440";
        }
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionUrl)) {
            System.out.println("PlcConnection " + plcConnection);

            {
                PlcReadRequest readRequest = plcConnection.readRequestBuilder().orElseThrow(() -> new RuntimeException("No Reader found"))
                .addItem("randomRegister", "register:7[3]").build();
                PlcReadResponse readResponse = readRequest.execute().get();
                System.out.println("Response " + readResponse);
                readResponse.getAllByteArrays("randomRegister").stream()
                    .map(HexUtil::toHex)
                    .map(hex -> "Register Value: " + hex)
                    .forEach(System.out::println);
            }

            {
                // Read an int from 2 registers

                // Just dump the actual values
                PlcReadRequest readRequest = plcConnection.readRequestBuilder().orElseThrow(() -> new RuntimeException("No Reader found"))
                    .addItem("randomRegister", "register:3[2]").build();
                PlcReadResponse readResponse = readRequest.execute().get();
                System.out.println("Response " + readResponse);
                Collection<Byte[]> randomRegisters = readResponse.getAllByteArrays("randomRegister");
                randomRegisters.stream()
                    .map(HexUtil::toHex)
                    .map(hex -> "Register Value: " + hex)
                    .forEach(System.out::println);

                // Read an actual int
                Byte[] registerBytes = randomRegisters.stream()
                    .flatMap(Arrays::stream)
                    .toArray(Byte[]::new);
                int readInt = ByteBuffer.wrap(ArrayUtils.toPrimitive(registerBytes))
                    .order(ByteOrder.BIG_ENDIAN)
                    .getInt();
                System.out.println("Read int " + readInt + " from register");
            }

            {
                // Read an int from 2 registers and multiple requests

                // Just dump the actual values
                PlcReadRequest readRequest = plcConnection.readRequestBuilder().orElseThrow(() -> new RuntimeException("No Reader found"))
                    .addItem("randomRegister1", "register:1[2]")
                    .addItem("randomRegister2", "register:10[3]")
                    .addItem("randomRegister3", "register:20[4]")
                    .addItem("randomRegister4", "register:30[5]")
                    .addItem("randomRegister5", "register:40[6]")
                    .build();
                PlcReadResponse readResponse = readRequest.execute().get();
                System.out.println("Response " + readResponse);
                IntStream.range(1, 6).forEach(i -> {
                    Collection<Byte[]> randomRegisters = readResponse.getAllByteArrays("randomRegister" + i);
                    randomRegisters.stream()
                        .map(HexUtil::toHex)
                        .map(hex -> "Register " + i + " Value: " + hex)
                        .forEach(System.out::println);

                    // Read an actual int
                    Byte[] registerBytes = randomRegisters.stream()
                        .flatMap(Arrays::stream)
                        .toArray(Byte[]::new);
                    int readInt = ByteBuffer.wrap(ArrayUtils.toPrimitive(registerBytes))
                        .order(ByteOrder.BIG_ENDIAN)
                        .getInt();
                    System.out.println("Read int " + i + " " + readInt + " from register");
                });
            }

            {
                PlcReadRequest readRequest = plcConnection.readRequestBuilder().orElseThrow(() -> new RuntimeException("No Reader found"))
                    .addItem("randomCoil", "coil:1[9]").build();
                PlcReadResponse readResponse = readRequest.execute().get();
                System.out.println("Response " + readResponse);
                readResponse.getAllBooleans("randomCoil").stream()
                    .map(hex -> "Coil Value: " + hex)
                    .forEach(System.out::println);
            }

            {
                PlcWriteRequest writeRequest = plcConnection.writeRequestBuilder().orElseThrow(() -> new RuntimeException("No Writer found"))
                    .addItem("randomCoilField", "coil:1", true).build();
                PlcWriteResponse writeResponse = writeRequest.execute().get();
                System.out.println("Response " + writeResponse);
            }
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }
        System.exit(0);
    }
}
