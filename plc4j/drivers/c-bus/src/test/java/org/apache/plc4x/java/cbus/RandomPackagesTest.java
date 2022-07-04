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

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Serial%20Interface%20User%20Guide.pdf
    @Nested
    class ReferenceDocumentationTest {

        // 3.4
        @Nested
        class Header {
            @Test
            void Point_point_multipoint_lowest_priority_class() throws Exception {
                byte[] bytes = new byte[]{0x03};
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusHeader msg = CBusHeader.staticParse(readBufferByteBased);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void Point_multipoint_lowest_priority_class() throws Exception {
                byte[] bytes = new byte[]{0x05};
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusHeader msg = CBusHeader.staticParse(readBufferByteBased);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void Point_point_lowest_priority_class() throws Exception {
                byte[] bytes = new byte[]{0x06};
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusHeader msg = CBusHeader.staticParse(readBufferByteBased);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

        }


        // 4.2.9.1
        @Nested
        class PointToPointCommands {
            @Test
            void pointToPointCommandDirect() throws Exception {
                byte[] bytes = "\\0603002102D4\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void pointToPointCommandBridged() throws Exception {
                byte[] bytes = "\\06420903210289\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }
        }


        // 4.2.9.2
        @Nested
        class PointToMultiPointCommands {
            @Test
            void pointToMultiPointCommandDirect() throws Exception {
                byte[] bytes = "\\0538000108BA\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void pointToMultiPointCommandBridged() throws Exception {
                byte[] bytes = "\\05FF007A38004A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }
        }

        // 4.2.9.3
        @Nested
        class PointToPointToMultoPointCommands {
            @Test
            void pointToPointToMultiPointCommand2() throws Exception {
                byte[] bytes = "\\03420938010871\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }
        }

        // 4.3.3.1
        @Nested
        class _CALReply {
            @Test
            void calRequest() throws Exception {
                byte[] bytes = "\\0605002102\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, false, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void calReplyNormal() throws Exception {
                byte[] bytes = Hex.decodeHex("8902312E322E363620200A");
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CALReply msg = CALReplyShort.staticParse(readBufferByteBased, false);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void calReplySmart() throws Exception {
                byte[] bytes = Hex.decodeHex("860593008902312E322E363620207F");
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CALReply msg = CALReplyLong.staticParse(readBufferByteBased, false);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }
        }

        // 4.3.3.2
        @Nested
        class MonitoredSAL {
            @Test
            @Disabled("Not yet implemented")
            void monitoredSal() throws Exception {
                byte[] bytes = "0503380079083F\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                Reply msg = Reply.staticParse(readBufferByteBased, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }
        }

        // 4.3.3.3
        @Nested
        class Confirmation {
            @Test
            void successful() throws Exception {
                byte[] bytes = "g.".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                Reply msg = Reply.staticParse(readBufferByteBased, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void toManyRetransmissions() throws Exception {
                byte[] bytes = "g#".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                Reply msg = Reply.staticParse(readBufferByteBased, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void corruption() throws Exception {
                byte[] bytes = "g$".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                Reply msg = Reply.staticParse(readBufferByteBased, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void desync() throws Exception {
                byte[] bytes = "g%".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                Reply msg = Reply.staticParse(readBufferByteBased, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void tooLong() throws Exception {
                byte[] bytes = "g'".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                Reply msg = Reply.staticParse(readBufferByteBased, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }
        }

        // 7.3
        @Test
        void StandardFormatStatusReply1() throws Exception {
            byte[] bytes = Hex.decodeHex("D8380068AA0140550550001000000014000000000000000000CF");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void StandardFormatStatusReply2() throws Exception {
            byte[] bytes = Hex.decodeHex("D838580000000000000000000000000000000000000000000098");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void StandardFormatStatusReply3() throws Exception {
            byte[] bytes = Hex.decodeHex("D638B000000000FF00000000000000000000000000000043");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 7.4
        @Test
        void ExtendedFormatStatusReply1() throws Exception {
            byte[] bytes = Hex.decodeHex("F9073800AAAA000095990000000005555000000000005555555548");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void ExtendedFormatStatusReply2() throws Exception {
            byte[] bytes = Hex.decodeHex("F907380B0000000000005555000000000000000000000000000013");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void ExtendedFormatStatusReply3() throws Exception {
            byte[] bytes = Hex.decodeHex("f70738160000000000000000000000000000000000000000B4");
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
                byte[] bytes = "\\0538000114AE\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void LightningStatus() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = "\\05FF007A38004A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here, the command requires to many bytes")
            @Test
            void LightningStatusReply1() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = "D83800A8AA02000000000000000000000000000000000000009C\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply2() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = "D838580000000000000000000000000000000000000000000098\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply3() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = "D638B0000000000000000000000000000000000000000042\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply4() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = "86999900F8003800A8AA0200000000000000000000000000000000000000C4\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }


            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply5() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = "86999900F800385800000000000000000000000000000000000000000000C0\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("something is wrong here")
            @Test
            void LightningStatusReply6() throws Exception {
                // TODO: the section describes that on non smart mode the message doesn't have the last CR
                byte[] bytes = "86999900F60038B000000000000000000000000000000000000000008F\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true, bytes.length);
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
                byte[] bytes = "\\0604001A3001AB\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Test
            void Reply() throws Exception {
                byte[] bytes = "8604990082300328\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, true, bytes.length);
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
                byte[] bytes = "\\03421B53643801149C\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

            @Disabled("it is not clear if that is a request or reply... it fails in both variants")
            @Test
            void Reply() throws Exception {
                byte[] bytes = "0565380354432101148E\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
                assertThat(msg)
                    .isNotNull();
                System.out.println(msg);
            }

        }

        // 9.4
        @Test
        void SwitchMode() throws Exception {
            // TODO: the section describes that on non smart mode the message doesn't have the last CR
            byte[] bytes = "@A3300019\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        // 9.5
        @Test
        void MultipleCommands() throws Exception {
            byte[] bytes = "\\05380001210122012301240A25010A2601D4\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void testParameterSet() throws Exception {
            byte[] bytes = "@A3470011\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
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
            byte[] bytes = "\\0538007902D4\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

    }

    @Nested
    class OwnCaptures {

        @Disabled
        @Test
        void whatEverThisIs() throws Exception {
            byte[] bytes = "\\3436303230303231303167\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }

        @Test
        void deviceManagementInstruction() throws Exception {
            byte[] bytes = "@1A2001\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, true, bytes.length);
            assertThat(msg)
                .isNotNull();
            System.out.println(msg);
        }
    }
}
