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
package org.apache.plc4x.java.cbus;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.cbus.readwrite.*;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

public class RandomPackagesTest {

    static final String BACKSLASH = "5C";
    static final String CR = "0D";
    static final String LF = "0A";

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Serial%20Interface%20User%20Guide.pdf
    @Nested
    class ReferenceDocumentationTest {
        // 4.2.9.1
        @Test
        void pointToPointCommandDirect() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0603002102D4" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand msg = CBusCommand.staticParse(readBufferByteBased, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.1
        @Test
        void pointToPointCommandBridged() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "06420903210289" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand msg = CBusCommand.staticParse(readBufferByteBased, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.2
        @Test
        void pointToMultiPointCommandDirect() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0538000108BA" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand msg = CBusCommand.staticParse(readBufferByteBased, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.2
        @Test
        void pointToMultiPointCommandBridged() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "05FF007A38004A" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand msg = CBusCommand.staticParse(readBufferByteBased, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.3
        @Test
        void pointToPointToMultiPointCommand2() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "03420938010871" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand msg = CBusCommand.staticParse(readBufferByteBased, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.3.3.1
        @Test
        void calRequest() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0605002102" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand msg = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.3.3.1
        @Test
        void calReplyNormal() throws Exception {
            byte[] bytes = Hex.decodeHex("8902312E322E363620200A" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CALReply msg = CALReplyShort.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void calReplySmart() throws Exception {
            byte[] bytes = Hex.decodeHex("860593008902312E322E363620207F" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CALReply msg = CALReplyLong.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.3.3.2
        @Test
        @Disabled("Not yet implemented")
        void monitoredSal() throws Exception {
            byte[] bytes = Hex.decodeHex("0503380079083F");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            Reply msg = MonitoredSALReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }
    }

}
