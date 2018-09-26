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

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.base.util.HexUtil;

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
                PlcReader reader = plcConnection.getReader().orElseThrow(() -> new RuntimeException("No Reader found"));

                PlcReadResponse<?> readResponse = reader.read(builder -> builder.addItem("randomRegister", "register:7[3]")).get();
                System.out.println("Response " + readResponse);
                readResponse.getAllByteArrays("randomRegister").stream()
                    .map(HexUtil::toHex)
                    .map(hex -> "Register Value: " + hex)
                    .forEach(System.out::println);
            }

            {
                PlcReader reader = plcConnection.getReader().orElseThrow(() -> new RuntimeException("No Reader found"));

                PlcReadResponse<?> readResponse = reader.read(builder -> builder.addItem("randomCoil", "coil:1[9]")).get();
                System.out.println("Response " + readResponse);
                readResponse.getAllBooleans("randomCoil").stream()
                    .map(hex -> "Coil Value: " + hex)
                    .forEach(System.out::println);
            }

            {
                PlcWriter writer = plcConnection.getWriter().orElseThrow(() -> new RuntimeException("No Writer found"));

                PlcWriteResponse<?> writeResponse = writer.write(builder -> builder.addItem("randomCoilField", "coil:1", true)).get();
                System.out.println("Response " + writeResponse);
            }
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }
        System.exit(0);
    }
}
