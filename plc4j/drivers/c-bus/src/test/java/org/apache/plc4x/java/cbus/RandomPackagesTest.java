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
import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

//@Disabled("non of those work yet")
public class RandomPackagesTest {

    static final String BACKSLASH = "5C";
    static final String CR = "0D";

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Serial%20Interface%20User%20Guide.pdf
    @Nested
    class ReferenceDocumentationTest {
        // 4.2.9.1
        @Test
        void pointToPointCommand1() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0603002102D4");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
        }

        // 4.2.9.1
        @Test
        void pointToPointCommand2() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "06420903210289");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
        }

        // 4.2.9.2
        @Test
        void pointToMultiPointCommand1() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0538000108BA");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
        }

        // 4.2.9.2
        @Test
        void pointToMultiPointCommand2() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "05FF007A38004A");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
        }

        // 4.2.9.3
        @Test
        void pointToPointToMultiPointCommand2() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "03420938010871");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
        }

        // 4.3.3.1
        @Test
        void calReply1() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0605002102" + CR);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
            System.out.println(cBusCommand);
        }

        // 4.3.3.1
        @Test
        void calReply2() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "860593008902312E322E363620207F" + CR);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
        }

        // 4.3.3.2
        @Test
        void monitoredSal() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0503380079083F");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusCommand cBusCommand = CBusCommand.staticParse(readBufferByteBased, false);
            assertThat(cBusCommand)
                .isNotNull();
        }
    }

}
