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
package org.apache.plc4x.java.cbus;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.cbus.readwrite.*;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.assertj.core.api.Assertions.assertThat;

public class RandomPackagesTest {

    static final String BACKSLASH = "5C";
    static final String CR = "0D";
    static final String LF = "0A";

    static final String TILDE = "7E";

    static final String AT = "40";

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Serial%20Interface%20User%20Guide.pdf
    @Nested
    class ReferenceDocumentationTest {
        // 4.2.9.1
        @Test
        void pointToPointCommandDirect() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0603002102D4" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.1
        @Test
        void pointToPointCommandBridged() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "06420903210289" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.2
        @Test
        void pointToMultiPointCommandDirect() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0538000108BA" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.2
        @Test
        void pointToMultiPointCommandBridged() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "05FF007A38004A" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.2.9.3
        @Test
        void pointToPointToMultiPointCommand2() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "03420938010871" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 4.3.3.1
        @Test
        void calRequest() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0605002102" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, false);
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

        // 7.3
        @Test
        void StandardFormatStatusReply1() throws Exception {
            byte[] bytes = Hex.decodeHex("D8380068AA0140550550001000000014000000000000000000CF" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void StandardFormatStatusReply2() throws Exception {
            byte[] bytes = Hex.decodeHex("D838580000000000000000000000000000000000000000000098" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void StandardFormatStatusReply3() throws Exception {
            byte[] bytes = Hex.decodeHex("D638B000000000FF00000000000000000000000000000043" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 7.4
        @Test
        void ExtendedFormatStatusReply1() throws Exception {
            byte[] bytes = Hex.decodeHex("F9073800AAAA000095990000000005555000000000005555555548" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void ExtendedFormatStatusReply2() throws Exception {
            byte[] bytes = Hex.decodeHex("F907380B0000000000005555000000000000000000000000000013" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void ExtendedFormatStatusReply3() throws Exception {
            byte[] bytes = Hex.decodeHex("f70738160000000000000000000000000000000000000000B4" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 9.1
        @Nested
        class PointToMultiPointCommandsIntoLocalCBusNetwork {
            @Test
            void LightningOff() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex(BACKSLASH + "0538000114AE" + CR);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void LightningStatus() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex(BACKSLASH + "05FF007A38004A" + CR);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply1() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex("D83800A8AA02000000000000000000000000000000000000009C" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply2() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex("D838580000000000000000000000000000000000000000000098" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply3() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex("D638B0000000000000000000000000000000000000000042" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply4() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex("86999900F8003800A8AA0200000000000000000000000000000000000000C4" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }


            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply5() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex("86999900F800385800000000000000000000000000000000000000000000C0" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply6() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex("86999900F60038B000000000000000000000000000000000000000008F" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }
        }

        // 9.2
        @Nested
        class PointToPointCommandsIntoLocalCBusNetwork {
            @Test
            void RecallCurrentValueOfParameter0x30onUnit0x04() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex(BACKSLASH + "0604001A3001AB" + CR);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void Reply() throws Exception {
                byte[] bytes = Hex.decodeHex("8604990082300328" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

        }

        // 9.3
        @Nested
        class PointToMultiPointCommandsIntoaRemoteCBusNetwork {
            @Test
            void IssueLightningOf() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = Hex.decodeHex(BACKSLASH + "03421B53643801149C" + CR);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("it is not clear if that is a request or reply... it fails in both variants")
            @Test
            void Reply() throws Exception {
                byte[] bytes = Hex.decodeHex("0565380354432101148E" + CR + LF);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

        }

        // 9.4
        @Disabled("no idea that is that here")
        @Test
        void SwitchMode() throws Exception {
            // TODO: the section describes that on non smart mode the message doesn't have the last CR
            byte[] bytes = Hex.decodeHex(/*TILDE +*/ AT + "A3300019" + CR);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 9.5
        @Disabled("no idea that is that here")
        @Test
        void MultipleCommands() throws Exception {
            // TODO: the section describes that on non smart mode the message doesn't have the last CR
            byte[] bytes = Hex.decodeHex(BACKSLASH + "05380001210122012301240A25010A2601D4" + CR);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }
    }

    @Nested
    class CBusQuickStartGuideTest {

        // 4.2.9.1
        @Test
        void pointToPointCommandDirect() throws Exception {
            byte[] bytes = Hex.decodeHex(BACKSLASH + "0538007902D4" + CR + LF);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

    }

}
