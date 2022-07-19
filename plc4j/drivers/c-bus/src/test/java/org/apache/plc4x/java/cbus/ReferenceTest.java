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
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.apache.plc4x.java.cbus.Util.assertMessageMatches;
import static org.assertj.core.api.Assertions.assertThat;

public class ReferenceTest {

    RequestContext requestContext;
    CBusOptions cBusOptions;

    @BeforeEach
    void setUp() {
        requestContext = new RequestContext(false, false, false);
        cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, false);
    }

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
                assertThat(msg).isNotNull();
                System.out.println(msg);
            }

            @Test
            void Point_multipoint_lowest_priority_class() throws Exception {
                byte[] bytes = new byte[]{0x05};
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusHeader msg = CBusHeader.staticParse(readBufferByteBased);
                assertThat(msg).isNotNull();
                System.out.println(msg);
            }

            @Test
            void Point_point_lowest_priority_class() throws Exception {
                byte[] bytes = new byte[]{0x06};
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusHeader msg = CBusHeader.staticParse(readBufferByteBased);
                assertThat(msg).isNotNull();
                System.out.println(msg);
            }

        }

        // 4
        @Nested
        class SerialInterface {

            // 4.2.3
            @Test
            void reset() throws Exception {
                byte[] bytes = "~\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
            }

            @Disabled("not implemented yet")
            // 4.2.4
            @Test
            void cancel() throws Exception {
                byte[] bytes = "AB0123?9876\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
            }

            // 4.2.5
            @Test
            void smartConnectShortcut() throws Exception {
                byte[] bytes = "\r|\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
            }

            @Disabled("not implemented yet")
            // 4.2.4
            @Test
            void confirmation() throws Exception {
                // If you follow the spec a confirmation can occur at any place... seems strange
                byte[] bytes = "AB0123n9876\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
            }

            // 4.2.7
            @Test
            void directCommandAccess1() throws Exception {
                byte[] bytes = "@2102\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CALData calData = ((CALDataOrSetParameterValue) ((RequestDirectCommandAccess) ((CBusMessageToServer) msg).getRequest()).getCalDataOrSetParameter()).getCalData();
                System.out.println(calData);
                assertMessageMatches(bytes, msg);
            }

            // 4.2.7
            @Test
            void directCommandAccess2() throws Exception {
                // TODO: this should be the same as the @above but that is not yet implemented
                byte[] bytes = "~2102\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CALData calData = ((CALDataOrSetParameterValue) ((RequestObsolete) ((CBusMessageToServer) msg).getRequest()).getCalDataOrSetParameter()).getCalData();
                System.out.println(calData);
                assertMessageMatches(bytes, msg);
            }

        }


        // 4.2.9.1
        @Nested
        class PointToPointCommands {
            @Test
            void pointToPointCommandDirect() throws Exception {
                byte[] bytes = "\\0603002102D4\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CBusCommand cbusCommand = ((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand();
                System.out.println(cbusCommand);
            }

            @Test
            void pointToPointCommandBridged() throws Exception {
                byte[] bytes = "\\06420903210289\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CBusCommand cbusCommand = ((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand();
                System.out.println(cbusCommand);
            }
        }


        // 4.2.9.2
        @Nested
        class PointToMultiPointCommands {
            @Test
            void pointToMultiPointCommandDirect() throws Exception {
                byte[] bytes = "\\0538000108BA\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CBusCommand cbusCommand = ((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand();
                System.out.println(cbusCommand);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void pointToMultiPointCommandBridged() throws Exception {
                byte[] bytes = "\\05FF007A38004A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CBusCommand cbusCommand = ((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand();
                System.out.println(cbusCommand);
                assertMessageMatches(bytes, msg);
            }
        }

        // 4.2.9.3
        @Nested
        class PointToPointToMultoPointCommands {
            @Test
            void pointToPointToMultiPointCommand2() throws Exception {
                byte[] bytes = "\\03420938010871\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CBusCommand cbusCommand = ((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand();
                System.out.println(cbusCommand);
                assertMessageMatches(bytes, msg);
            }
        }

        // 4.3.3.1
        @Nested
        class _CALReply {
            @Test
            void calRequest() throws Exception {
                byte[] bytes = "\\0605002102\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void calReplyNormal() throws Exception {
                byte[] bytes = Hex.decodeHex("8902312E322E363620200A");
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CALReply msg = CALReplyShort.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void calReplySmart() throws Exception {
                byte[] bytes = Hex.decodeHex("860593008902312E322E363620207F");
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CALReply msg = CALReplyLong.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }
        }

        // 4.3.3.2
        @Nested
        class _MonitoredSAL {
            @Test
            void monitoredSal() throws Exception {
                byte[] bytes = "0503380079083F\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, bytes.length, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                EncodedReply encodedReply = ((ReplyEncodedReply) ((ReplyOrConfirmationReply) msg).getReply()).getEncodedReply();
                MonitoredSAL monitoredSAL = ((MonitoredSALReply) encodedReply).getMonitoredSAL();
                System.out.println(monitoredSAL);
                assertMessageMatches(bytes, msg);
            }
        }

        // 4.3.3.3
        @Nested
        class Confirmation {
            @Test
            void successful() throws Exception {
                byte[] bytes = "g.".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, bytes.length, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void toManyRetransmissions() throws Exception {
                byte[] bytes = "g#".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, bytes.length, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void corruption() throws Exception {
                byte[] bytes = "g$".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, bytes.length, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void desync() throws Exception {
                byte[] bytes = "g%".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, bytes.length, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void tooLong() throws Exception {
                byte[] bytes = "g'".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, bytes.length, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }
        }

        // 7.3
        @Test
        void StandardFormatStatusReply1() throws Exception {
            byte[] bytes = Hex.decodeHex("D8380068AA0140550550001000000014000000000000000000CF");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        @Test
        void StandardFormatStatusReply2() throws Exception {
            byte[] bytes = Hex.decodeHex("D838580000000000000000000000000000000000000000000098");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        @Test
        void StandardFormatStatusReply3() throws Exception {
            byte[] bytes = Hex.decodeHex("D638B000000000FF00000000000000000000000000000043");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            StandardFormatStatusReply msg = StandardFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        // 7.4
        @Test
        void ExtendedFormatStatusReply1() throws Exception {
            byte[] bytes = Hex.decodeHex("F9073800AAAA000095990000000005555000000000005555555548");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        @Test
        void ExtendedFormatStatusReply2() throws Exception {
            byte[] bytes = Hex.decodeHex("F907380B0000000000005555000000000000000000000000000013");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        @Test
        void ExtendedFormatStatusReply3() throws Exception {
            byte[] bytes = Hex.decodeHex("f70738160000000000000000000000000000000000000000B4");
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            ExtendedFormatStatusReply msg = ExtendedFormatStatusReply.staticParse(readBufferByteBased, false);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        // 9.1
        @Nested
        class PointToMultiPointCommandsIntoLocalCBusNetwork {
            @Test
            void LightningOff() throws Exception {
                byte[] bytes = "\\0538000114AE\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatus() throws Exception {
                byte[] bytes = "\\05FF007A38004A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply1() throws Exception {
                byte[] bytes = "D83800A8AA02000000000000000000000000000000000000009C\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                requestContext = new RequestContext(false, true, false);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                StandardFormatStatusReply reply = ((EncodedReplyStandardFormatStatusReply) ((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply()).getReply();
                System.out.println(reply);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply2() throws Exception {
                byte[] bytes = "D838580000000000000000000000000000000000000000000098\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((EncodedReplyStandardFormatStatusReply) ((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply()).getReply());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply3() throws Exception {
                byte[] bytes = "D638B0000000000000000000000000000000000000000042\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((EncodedReplyStandardFormatStatusReply) ((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply()).getReply());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply4() throws Exception {
                byte[] bytes = "86999900F8003800A8AA0200000000000000000000000000000000000000C4\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((EncodedReplyCALReply) ((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply()).getCalReply());
                assertMessageMatches(bytes, msg);
            }


            @Test
            void LightningStatusReply5() throws Exception {
                byte[] bytes = "86999900F800385800000000000000000000000000000000000000000000C0\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((EncodedReplyCALReply) ((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply()).getCalReply());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply6() throws Exception {
                byte[] bytes = "86999900F60038B000000000000000000000000000000000000000008F\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((EncodedReplyCALReply) ((ReplyEncodedReply) ((ReplyOrConfirmationReply) ((CBusMessageToClient) msg).getReply()).getReply()).getEncodedReply()).getCalReply());
                assertMessageMatches(bytes, msg);
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
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

            @Test
            void Reply() throws Exception {
                byte[] bytes = "8604990082300328\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
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
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                CBusCommand cbusCommand = ((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand();
                System.out.println(cbusCommand);
                assertMessageMatches(bytes, msg);
            }

            @Disabled("the transformation from point to point to multipoint message is not yet implemented as described in 8.4... not sure if we will ever implement that")
            @Test
            void Reply() throws Exception {
                byte[] bytes = Hex.decodeHex("0565380354432101148E");
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusCommand msg = CBusCommand.staticParse(readBufferByteBased, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                assertMessageMatches(bytes, msg);
            }

        }

        // 9.4
        @Test
        void SwitchMode() throws Exception {
            byte[] bytes = "~@A3300019\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        // 9.5
        @Test
        void MultipleCommands() throws Exception {
            byte[] bytes = "\\05380001210122012301240A25010A2601D4\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        // 10.2.1
        @Test
        void testParameterSet() throws Exception {
            byte[] bytes = "@A3470011\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

        // 10.2.1
        @Test
        void testParameterSetObsolete() throws Exception {
            byte[] bytes = "A3470011\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            CALDataOrSetParameterSetParameter calDataOrSetParameter = (CALDataOrSetParameterSetParameter) ((RequestObsolete) ((CBusMessageToServer) msg).getRequest()).getCalDataOrSetParameter();
            System.out.println(calDataOrSetParameter);
            assertMessageMatches(bytes, msg);
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/C-Bus%20Quick%20Start%20Guide.pdf
    @Nested
    class CBusQuickStartGuideTest {

        // 4.2.9.1
        @Test
        void pointToPointCommandDirect() throws Exception {
            byte[] bytes = "\\0538007902D4\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            assertMessageMatches(bytes, msg);
        }

    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2002%20-%20C-Bus%20Lighting%20Application.pdf
    @Nested
    class LightningApplicationsTest {

        // 2.9.6.4.4 Command Sequence
        @Nested
        class CommandSquence {

            @Test
            void StartDynamicIcon() throws Exception {
                byte[] bytes = "\\053800A412080020\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void IconBitmap() throws Exception {
                byte[] bytes = "\\053800A412080021\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void CompleteDynamicIcon() throws Exception {
                byte[] bytes = "\\053800A412080022\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Nested
            class ChineseTable {

                @Test
                void StartDynamicIcon() throws Exception {
                    byte[] bytes = "\\053800A401080020\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);
                    System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                    assertMessageMatches(bytes, msg);
                }

                @Test
                void IconHeader() throws Exception {
                    byte[] bytes = "\\053800A80104CA00130C0600\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);
                    System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                    assertMessageMatches(bytes, msg);
                }

                @Test
                void AppendDynamicIcon() throws Exception {
                    byte[] bytes = "\\053800A401080021\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);
                    System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                    assertMessageMatches(bytes, msg);
                }

                @Test
                void WriteIconBitmapData1() throws Exception {
                    byte[] bytes = "\\053800A80104AAF05500FF50\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);
                    System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                    assertMessageMatches(bytes, msg);
                }

                @Test
                void WriteIconBitmapData2() throws Exception {
                    byte[] bytes = "\\053800A801040000F0F00F00\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);
                    System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                    assertMessageMatches(bytes, msg);
                }

                @Test
                void useDynamicIcon() throws Exception {
                    byte[] bytes = "\\053800A401080022\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);
                    System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                    assertMessageMatches(bytes, msg);
                }

                @Test
                void displayIcon() throws Exception {
                    byte[] bytes = "\\053800A60102CA010013\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);
                    System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                    assertMessageMatches(bytes, msg);
                }

            }
        }

        // 2.11 Examples
        @Nested
        class Examples {
            @Test
            void switchElectricalLoads() throws Exception {
                byte[] bytes = "\\0538007993B7\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void switchElectricalLoadsBridged() throws Exception {
                byte[] bytes = "\\0356093879935A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2005%20-%20C-Bus%20Security%20Application.pdf
    @Nested
    class SecurityApplicationsTest {

        //5.11.1
        @Nested
        class SecuritySystemEmitsAlarmOn{
            @Test
            void AlarmOnWrongPrio() throws Exception {
                byte[] bytes = "\\05D00079832F\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void AlarmOn() throws Exception {
                byte[] bytes = "\\85D0007983AF\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

        }

        //5.11.2
        @Test
        void Zone3Unsealed() throws Exception {
            byte[] bytes = "\\05D0000A860398\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
            assertMessageMatches(bytes, msg);
        }

        //5.11.3
        @Test
        void ZoneName() throws Exception {
            byte[] bytes = "\\05D000AD8D034B49544348452E2020202088\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
            assertThat(msg).isNotNull();
            System.out.println(msg);
            System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
            assertMessageMatches(bytes, msg);
        }

        // 5.11.4
        @Nested
        class DeviceRequestsSecuritySystemtoArm {

            @Test
            void ArmSecurity() throws Exception {
                byte[] bytes = "\\05D0000AA2FF80\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
            @Test
            void ArmSecurityRemote() throws Exception {
                byte[] bytes = "\\039209D00AA2FFE7\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
        }
    }


    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2006%20-%20C-Bus%20Metering%20Application.pdf
    @Nested
    class MeteringApplicationsTest{

        //6.11.1
        @Nested
        class DeviceRequestsMeteringApplicationtoMeasureElectricity{
            @Test
            void LocalMeasurement() throws Exception {
                byte[] bytes = "\\05D100090120\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteMeasurement() throws Exception {
                byte[] bytes = "\\035609D10901C3\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

        }

        //6.11.2
        @Nested
        class MeterMeasurementDevicesendsElectricityUse{
            @Test
            void LocalMeasurement() throws Exception {
                byte[] bytes = "\\05D1000D810000DBF8C9\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteMeasurement() throws Exception {
                byte[] bytes = "\\033709D10D810000DBF88B\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2007%20-%20C-Bus%20Trigger%20Control%20Application.pdf
    @Nested
    class TriggerControlApplicationsTest{

        //7.12
        @Nested
        class Examples{
            @Test
            void LocalTrigger() throws Exception {
                byte[] bytes = "\\05CA0002250109\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteTrigger() throws Exception {
                byte[] bytes = "\\035609CA022501AC\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2008%20-%20C-Bus%20Enable%20Control%20Application.pdf
    @Nested
    class EnableControlApplicationsTest{

        //8.11
        @Nested
        class Examples{
            @Test
            void LocalTrigger() throws Exception {
                byte[] bytes = "\\05CB0002378275\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteTrigger() throws Exception {
                byte[] bytes = "\\035609CB02378216\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2009%20-%20C-Bus%20Temperature%20Broadcast%20Application.pdf
    @Nested
    class TemperatureBroadcastApplicationsTest{

        //9.11
        @Nested
        class Examples{

            @Test
            void temperatureBroadcast() throws Exception {
                byte[] bytes = "\\051900020564\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2010%20-%20C-Bus%20Ventilation%20Application.pdf
    @Nested
    class VentilationApplicationTest{
        // TODO: no tests described here but it should work by adjusting the values from Lightning...
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2009%20-%20C-Bus%20Temperature%20Control%20Application.pdf
    @Nested
    class AccessControlApplicationsTest{

        //9.11
        @Nested
        class Examples{

            @Test
            void validAccessRequest() throws Exception {
                byte[] bytes = "\\05D500A4010300017D\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void closeAccessPoint() throws Exception {
                byte[] bytes = "\\05D5000201FF24\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void lockAccessPoint() throws Exception {
                byte[] bytes = "\\05D5000AFFFF1E\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void lockAccessPointRemote() throws Exception {
                byte[] bytes = "\\039209D50AFFFF85\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2021%20-%20C-Bus%20Media%20Transport%20Control%20Application.pdf
    @Nested
    class MediaTransportControlApplicationsTest{
        // TODO: no tests described here
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2023%20-%20C-Bus%20Clock%20and%20Timekeeping%20Application.pdf
    @Nested
    class ClockAndTimekeeping{

        //23.13
        @Nested
        class Examples{

            @Test
            void outputATimeCommand() throws Exception {
                byte[] bytes = "\\05DF000D010A2B1700C2\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void outputADateCommand() throws Exception {
                byte[] bytes = "\\05DF000E0207D502190411\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Disabled("This example just seems plain wrong... First of all there is no command for 0x10 defined and the argument which should be it 0x11 requires a argument of 0x03... So either documenation wrong or the example")
            @Test
            void outputARequestRefreshCommand() throws Exception {
                byte[] bytes = "\\05DF00100C\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }

            @Test
            void outputARequestRefreshCommandFixedQuestionMark() throws Exception {
                byte[] bytes = "\\05DF001103\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions, bytes.length);
                assertThat(msg).isNotNull();
                System.out.println(msg);
                System.out.println(((RequestCommand) ((CBusMessageToServer) msg).getRequest()).getCbusCommand());
                assertMessageMatches(bytes, msg);
            }
        }
    }
}
