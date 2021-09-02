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
package org.apache.plc4x.java.s7.readwrite;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class DatatypesTest {

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = new PlcDriverManager().getConnection("s7://192.168.23.30")) {
            final PlcReadRequest.Builder builder = connection.readRequestBuilder();
            builder.addItem("bool-value-1", "%DB2:0.0:BOOL"); // true
            builder.addItem("bool-value-2", "%DB2:2.1:BOOL"); // false
            // It seems S7 PLCs ignores the array notation for BOOL
            //builder.addItem("bool-array", "%DB2:2.1:BOOL[4]");
            builder.addItem("byte-value", "%DB2:2:BYTE");
            builder.addItem("byte-array", "%DB2:2:BYTE[2]");
            builder.addItem("word-value", "%DB2:2:WORD");
            builder.addItem("word-array", "%DB2:2:WORD[2]");
            builder.addItem("dword-value", "%DB2:2:DWORD");
            builder.addItem("dword-array", "%DB2:2:DWORD[2]");
            builder.addItem("sint-value", "%DB2:12:SINT"); // 7
            builder.addItem("sint-array", "%DB2:14:SINT[2]"); // 1, -2
            builder.addItem("int-value", "%DB2:18:INT"); // 23
            builder.addItem("int-array", "%DB2:20:INT[2]"); // 123, -142
            builder.addItem("dint-value", "%DB2:24:DINT"); // 24
            builder.addItem("dint-array", "%DB2:28:DINT[2]"); // 1234, -2345
            builder.addItem("usint-value", "%DB2:36:USINT"); // 42
            builder.addItem("usint-array", "%DB2:38:USINT[2]"); // 3, 4
            builder.addItem("uint-value", "%DB2:40:UINT"); // 3
            builder.addItem("uint-array", "%DB2:42:UINT[2]"); // 242, 223
            builder.addItem("udint-value", "%DB2:46:UDINT"); // 815
            builder.addItem("udint-array", "%DB2:50:UDINT[2]"); // 12345, 23456
            builder.addItem("real-value", "%DB2:58:REAL"); // 3.14159
            builder.addItem("real-array", "%DB2:62:REAL[2]"); // 12.345, 12.345
            builder.addItem("lreal-value", "%DB2:70:LREAL"); // 3.14159265358979
            builder.addItem("lreal-array", "%DB2:78:LREAL[2]"); // 1.2345, -1.2345
            builder.addItem("string-value", "%DB2:94:STRING(10)"); // "Hurz"
            // When reading a sized STRING string array, this has to be translated into multiple items
            //builder.addItem("string-array", "%DB2:350:STRING(10)[2]"); // "Wolf", "Lamm"
            builder.addItem("time-value", "%DB2:862:TIME"); // 1234ms
            builder.addItem("time-array", "%DB2:866:TIME[2]"); // 123ms, 234ms
            builder.addItem("date-value", "%DB2:874:DATE"); // D#2020-08-20
            builder.addItem("date-array", "%DB2:876:DATE[2]"); // D#1990-03-28, D#2020-10-25
            builder.addItem("time-of-day-value", "%DB2:880:TIME_OF_DAY"); // TOD#12:34:56
            builder.addItem("time-of-day-array", "%DB2:884:TIME_OF_DAY[2]"); // TOD#16:34:56, TOD#08:15:00
            builder.addItem("date-and-time-value", "%DB2:892:DATE_AND_TIME"); // DTL#1978-03-28-12:34:56
            builder.addItem("date-and-time-array", "%DB2:904:DATE_AND_TIME[2]"); // DTL#1978-03-28-12:34:56, DTL#1978-03-28-12:34:56
            builder.addItem("char-value", "%DB2:928:CHAR"); // "H"
            builder.addItem("char-array", "%DB2:930:CHAR[4]"); // "H", "u", "r", "z"
            final PlcReadRequest readRequest = builder.build();

            final PlcReadResponse readResponse = readRequest.execute().get();

            System.out.println(readResponse);

        }
    }

}
