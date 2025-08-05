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
package org.apache.plc4x.java.s7.readwrite;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class DatatypesTest {

    public static void main(String[] args) throws Exception {
        //try (PlcConnection connection = new DefaultPlcDriverManager().getConnection("s7://192.168.24.83")) {
        //
        String URL1 = "s7://192.168.23.30";
        
        //
        String URL2 = "s7://192.168.0.47?remote-rack=0&"
                + "remote-slot=3&"
                + "controller-type=S7_400&read-timeout=8&"
                + "ping=false&ping-time=2&retry-time=3";
        
        
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(URL1)) {
                final PlcReadRequest.Builder builder = connection.readRequestBuilder();
            builder.addTagAddress("bool-value-1", "%DB2:0.0:BOOL"); // true
            builder.addTagAddress("bool-value-2", "%DB2:2.1:BOOL"); // false
            // It seems S7 PLCs ignores the array notation for BOOL
            //builder.addField("bool-array", "%DB2:2.1:BOOL[4]");
            builder.addTagAddress("byte-value", "%DB2:2:BYTE");
            builder.addTagAddress("byte-array", "%DB2:2:BYTE[2]");
            builder.addTagAddress("word-value", "%DB2:2:WORD");
            builder.addTagAddress("word-array", "%DB2:2:WORD[2]");
            builder.addTagAddress("dword-value", "%DB2:2:DWORD");
            builder.addTagAddress("dword-array", "%DB2:2:DWORD[2]");
            builder.addTagAddress("sint-value", "%DB2:12:SINT"); // 7
            builder.addTagAddress("sint-array", "%DB2:14:SINT[2]"); // 1, -2
            builder.addTagAddress("int-value", "%DB2:18:INT"); // 23
            builder.addTagAddress("int-array", "%DB2:20:INT[2]"); // 123, -142
            builder.addTagAddress("dint-value", "%DB2:24:DINT"); // 24
            builder.addTagAddress("dint-array", "%DB2:28:DINT[2]"); // 1234, -2345
            builder.addTagAddress("usint-value", "%DB2:36:USINT"); // 42
            builder.addTagAddress("usint-array", "%DB2:38:USINT[2]"); // 3, 4
            builder.addTagAddress("uint-value", "%DB2:40:UINT"); // 3
            builder.addTagAddress("uint-array", "%DB2:42:UINT[2]"); // 242, 223
            builder.addTagAddress("udint-value", "%DB2:46:UDINT"); // 815
            builder.addTagAddress("udint-array", "%DB2:50:UDINT[2]"); // 12345, 23456
            builder.addTagAddress("real-value", "%DB2:58:REAL"); // 3.14159
            builder.addTagAddress("real-array", "%DB2:62:REAL[2]"); // 12.345, 12.345
            
//            builder.addTagAddress("lreal-value", "%DB2:70:LREAL"); // 3.14159265358979
//            builder.addTagAddress("lreal-array", "%DB2:78:LREAL[2]"); // 1.2345, -1.2345
            
            builder.addTagAddress("string-value", "%DB2:94:STRING(10)"); // "Hurz"
            
            // When reading a sized STRING array, this has to be translated into multiple items
            //builder.addField("string-array", "%DB2:350:STRING(10)[2]"); // "Wolf", "Lamm"
            builder.addTagAddress("time-value", "%DB2:862:TIME"); // 1234ms
            builder.addTagAddress("time-array", "%DB2:866:TIME[2]"); // 123ms, 234ms
            builder.addTagAddress("date-value", "%DB2:874:DATE"); // D#2020-08-20
            builder.addTagAddress("date-array", "%DB2:876:DATE[2]"); // D#1990-03-28, D#2020-10-25
            builder.addTagAddress("time-of-day-value", "%DB2:880:TIME_OF_DAY"); // TOD#12:34:56
            builder.addTagAddress("time-of-day-array", "%DB2:884:TIME_OF_DAY[2]"); // TOD#16:34:56, TOD#08:15:00
            builder.addTagAddress("date-and-time-value", "%DB2:892:DATE_AND_TIME"); // DTL#1978-03-28-12:34:56
//            builder.addTagAddress("date-and-time-array", "%DB2:904:DATE_AND_TIME[2]"); // DTL#1978-03-28-12:34:56, DTL#1978-03-28-12:34:56
            builder.addTagAddress("char-value", "%DB2:928:CHAR"); // "H"
            builder.addTagAddress("char-array", "%DB2:930:CHAR[4]"); // "H", "u", "r", "z"

            //builder.addTagAddress("huge-array", "%DB2:0:BYTE[900]");
            builder.addTagAddress("huge-array", "%DB2:0:UINT[4000]");
            //builder.addTagAddress("huge-array", "%DB2:0:UINT[4000]");
            final PlcReadRequest readRequest = builder.build();

            final PlcReadResponse readResponse = readRequest.execute().get();

            System.out.println(readResponse);
            byte[] responseBytes = readResponse.getPlcValue("huge-array").getRaw();
            for(int i = 0; i < 2000; i++) {
                byte firstByte = (byte) (i & 0x00FF);
                byte secondByte = (byte) ((i & 0xFF00) >> 8);
                if((responseBytes[i*2] != secondByte) || (responseBytes[i*2+1] != firstByte)) {
                    System.out.println("Difference");
                }
            }
        }
    }

}
