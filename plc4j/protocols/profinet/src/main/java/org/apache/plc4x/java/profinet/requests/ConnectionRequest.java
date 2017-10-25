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
package org.apache.plc4x.java.profinet.requests;

import org.apache.plc4x.java.profinet.types.DeviceGroup;

public class ConnectionRequest {

    /*
     * Template for a connection request packet.
     * Parts with "(parameter)" will be customized when using the template.
     * Related Links:
     * - S7 Protocol (http://gmiru.com/article/s7comm/)
     * - ISO Transport Protocol (Class 0) (https://tools.ietf.org/html/rfc905)
     * - ISO on TCP (https://tools.ietf.org/html/rfc1006)
     * - Reference to calculating the TSAP ids: https://www.tanindustrie.de/fr/Help/ConfigClient/tsap_s7.htm
     * - Structure and some constants of a variable read/write request:
     *      https://support.industry.siemens.com/tf/ww/en/posts/classic-style-any-pounter-to-variant-type/126024/?page=0&pageSize=10
     */
    private static final byte[] TEMPLATE = {
        ////////////////////////////////////////////////////
        // RFC 1006 (ISO on TCP)
        (byte) 0x03,                // Version (is always constant 0x03)
        (byte) 0x00,                // Reserved (is always constant 0x00)
        (byte) 0x00, (byte) 0x16,   // Packet length (including ISOonTCP header)

        ////////////////////////////////////////////////////
        // RFC 905 (ISO Transport Protocol)
        (byte) 0x11,                // Length indicator field
        //                             (Length of the header excluding the length indicator itself)
        //  Fixed part:
        (byte) 0xe0,                // TPDU Code (First 4 bits, 1110 = Connection Request)
        //                             (Second 4 bits: Initial Credit Allocation (constantly set to 0)
        (byte) 0x00, (byte) 0x00,   // Destination Reference (constantly set to 0)
        (byte) 0x00, (byte) 0x00,   // (parameter) Source Reference (Connection ID set by the client to identify the connection)
        (byte) 0x00,                // Preferred Transport Protocol Class (First 4 bits, 0000 = Class 0)
        //                             Options (All options set to 0 for Class 0)
        //  Variable part:
        //      - TPDU Size
        (byte) 0xc0,                // Parameter Code (1100 0000 = TPDU size)
        (byte) 0x01,                // Length (TPDU size constantly 1)
        (byte) 0x0a,                // Size of the TPDU packets sent in this connection
        //                             0x0b = 2048 bytes
        //                             0x0a = 1024 bytes (*)
        //                             0x09 = 512 bytes
        //                             0x08 = 256 bytes
        //                             0x07 = 128 bytes
        //      - Calling Transport Service Access Point Identifier
        (byte) 0xc1,                // Parameter Code (1100 0001 = Calling Transport Service Access Point Identifier)
        (byte) 0x02,                // Length (constantly 2)
        DeviceGroup.PG_OR_PC.getCode(),// Device Group:
        //                             0x01 = PG or PC (*)
        //                             0x02 = OS (operating or monitoring device)
        //                             0x03 = Others, such as OPC server, S7 PLC, ...
        (byte) 0x00,                // 3 bits (5..7) Rack Number
        //                             5 bits (4..0) CPU Slot Number
        //      - Called Transport Service Access Point Identifier
        (byte) 0xc2,                // Parameter Code (1100 0010 = Called Transport Service Access Point Identifier)
        (byte) 0x02,                // Length (constantly 2)
        DeviceGroup.OTHERS.getCode(),// Device Group:
        //                             0x01 = PG or PC (*)
        //                             0x02 = OS (operating or monitoring device)
        //                             0x03 = Others, such as OPC server, S7 PLC, ...
        (byte) 0x00                 // 3 bits (Bits 6-8) Rack Number
        //                             5 bits (Bity 1-5) CPU Slot Number

        ////////////////////////////////////////////////////
        // S7 User-Data
        // None: The connection request is completely handled by the lower layers.
    };

}
