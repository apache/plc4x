/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class ManualDriverTest {

    /**
     * Test programm made to work with the Modbus Simulator from https://www.modbustools.com/modbus_slave.html
     * In contrast to others, this allows simulating the various types of modbus byte orders.
     *
     * @param args
     * @throws Exception
     */
    public static void main(String[] args) throws Exception {
        //final PlcConnection connection = new DefaultPlcDriverManager().getConnection("modbus-tcp://10.211.55.3?default-payload-byte-order=BIG_ENDIAN");
        final PlcConnection connection = new DefaultPlcDriverManager().getConnection("modbus-tcp://10.211.55.3?default-payload-byte-order=LITTLE_ENDIAN");
        //final PlcConnection connection = new DefaultPlcDriverManager().getConnection("modbus-tcp://10.211.55.3?default-payload-byte-order=BIG_ENDIAN_BYTE_SWAP");
        //final PlcConnection connection = new DefaultPlcDriverManager().getConnection("modbus-tcp://10.211.55.3?default-payload-byte-order=LITTLE_ENDIAN_BYTE_SWAP");
        final PlcReadRequest readRequest = connection.readRequestBuilder()
            .addTagAddress("16 bit BigEndian", "holding-register:1:WORD{byte-order:'BIG_ENDIAN'}")
            .addTagAddress("16 bit LittleEndian", "holding-register:2:WORD{byte-order:'LITTLE_ENDIAN'}")
            .addTagAddress("32 bit BigEndian", "holding-register:3:DWORD{byte-order:'BIG_ENDIAN'}")
            .addTagAddress("32 bit LittleEndian", "holding-register:5:DWORD{byte-order:'LITTLE_ENDIAN'}")
            .addTagAddress("32 bit BigEndianByteSwap", "holding-register:7:DWORD{byte-order:'BIG_ENDIAN_BYTE_SWAP'}")
            .addTagAddress("32 bit LittleEndianByteSwap", "holding-register:9:DWORD{byte-order:'LITTLE_ENDIAN_BYTE_SWAP'}")
            .addTagAddress("64 bit BigEndian", "holding-register:11:LWORD{byte-order:'BIG_ENDIAN'}")
            .addTagAddress("64 bit LittleEndian", "holding-register:15:LWORD{byte-order:'LITTLE_ENDIAN'}")
            .addTagAddress("64 bit BigEndianByteSwap", "holding-register:19:LWORD{byte-order:'BIG_ENDIAN_BYTE_SWAP'}")
            .addTagAddress("64 bit LittleEndianByteSwap", "holding-register:23:LWORD{byte-order:'LITTLE_ENDIAN_BYTE_SWAP'}")
            .build();
        final PlcReadResponse plcReadResponse = readRequest.execute().get();
        connection.close();
        System.out.println("16 bit BigEndian:            " + String.format("0x%04X", plcReadResponse.getInteger("16 bit BigEndian")));
        System.out.println("16 bit LittleEndian:         " + String.format("0x%04X", plcReadResponse.getInteger("16 bit LittleEndian")));
        System.out.println("32 bit BigEndian:            " + String.format("0x%08X", plcReadResponse.getInteger("32 bit BigEndian")));
        System.out.println("32 bit LittleEndian:         " + String.format("0x%08X", plcReadResponse.getInteger("32 bit LittleEndian")));
        System.out.println("32 bit BigEndianByteSwap:    " + String.format("0x%08X", plcReadResponse.getInteger("32 bit BigEndianByteSwap")));
        System.out.println("32 bit LittleEndianByteSwap: " + String.format("0x%08X", plcReadResponse.getInteger("32 bit LittleEndianByteSwap")));
        System.out.println("64 bit BigEndian:            " + String.format("0x%016X", plcReadResponse.getLong("64 bit BigEndian")));
        System.out.println("64 bit LittleEndian:         " + String.format("0x%016X", plcReadResponse.getLong("64 bit LittleEndian")));
        System.out.println("64 bit BigEndianByteSwap:    " + String.format("0x%016X", plcReadResponse.getLong("64 bit BigEndianByteSwap")));
        System.out.println("64 bit LittleEndianByteSwap: " + String.format("0x%016X", plcReadResponse.getLong("64 bit LittleEndianByteSwap")));
    }

}
