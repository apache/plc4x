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

    public static final CBusOptions C_BUS_OPTIONS_WITH_SRCHK = new CBusOptions(false, false, false, false, false, false, false, false, true);

    RequestContext requestContext;
    CBusOptions cBusOptions;

    @BeforeEach
    void setUp() {
        requestContext = new RequestContext(false);
        cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, false);
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/C-Bus%20Interface%20Requirements.pdf
    @Nested
    class InterfaceRequirementsTest {

        // 8.2
        @Nested
        class Level1InterfaceImplementationRequirements {

            // 8.2.4
            @Nested
            class SerialInterfaceInitialisation {

                @Test
                void Step_1_Reset() throws Exception {
                    byte[] bytes = "~~~\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_2_SetInterfaceOptions3() throws Exception {
                    byte[] bytes = "@A3420002\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_3_SetInterfaceOptions1_PUN() throws Exception {
                    byte[] bytes = "@A3410058\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_4_SetInterfaceOptions1() throws Exception {
                    byte[] bytes = "@A3300058\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }

            // 8.2.5
            @Nested
            class ConfirmationOfTransmission {
                @Test
                void SomeCommand() throws Exception {
                    byte[] bytes = "\\0538000121A1g\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void success() throws Exception {
                    byte[] bytes = "g.".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void checksumFailure() throws Exception {
                    byte[] bytes = "g!".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void tooManyRetransmissions() throws Exception {
                    byte[] bytes = "g#".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void corruptionInTransmission() throws Exception {
                    byte[] bytes = "g$".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void noSystemClock() throws Exception {
                    byte[] bytes = "g%".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }
        }

        // 8.3
        @Nested
        class Level2InterfaceImplementationRequirements {

            // 8.3.4
            @Nested
            class SerialInterfaceInitialisation {
                @Test
                void Step_1_Reset() throws Exception {
                    byte[] bytes = "~~~\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_2_AnyApplicationFilter() throws Exception {
                    byte[] bytes = "@A3210038\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_3_SetInterfaceOptions3() throws Exception {
                    byte[] bytes = "@A3420002\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_4_SetInterfaceOptions1_PUN() throws Exception {
                    byte[] bytes = "@A3410059\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_5_SetInterfaceOptions1() throws Exception {
                    byte[] bytes = "@A3300059\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }

            // 8.3.5
            @Nested
            class ProgrammingTheSerialInterfaceToFilterSALMessageTraffic {

                @Test
                void Step_1_SelectOnlyLighting() throws Exception {
                    byte[] bytes = "@A3210038\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_2_SelectHeatingAsSecondApplication() throws Exception {
                    byte[] bytes = "@A3220088\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }
        }

        // 8.4
        @Nested
        class Level3InterfaceImplementationRequirements {
            // No specific tests
        }

        // 8.5
        @Nested
        class Level4InterfaceImplementationRequirements {

            // 8.5.4
            @Nested
            class SerialInterfaceInitialisation {
                @Test
                void Step_1_Reset() throws Exception {
                    byte[] bytes = "~~~\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_2_AnyApplicationFilter() throws Exception {
                    byte[] bytes = "@A3210038\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_3_SetInterfaceOptions3() throws Exception {
                    byte[] bytes = "@A342000A\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_4_SetInterfaceOptions1_PUN() throws Exception {
                    byte[] bytes = "@A3410079\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Step_5_SetInterfaceOptions1() throws Exception {
                    byte[] bytes = "@A3300079\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }
        }

        // 8.6
        @Nested
        class Level5InterfaceImplementationRequirements {
            // No specific tests
        }

        // 8.7
        @Nested
        class Level6InterfaceImplementationRequirements {
            // No specific tests
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Serial%20Interface%20User%20Guide.pdf
    @Nested
    class SerialInterfaceGuideTest {

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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();

                System.out.println(msg);
            }

            @Disabled("not implemented yet")
            // 4.2.4
            @Test
            void cancel() throws Exception {
                byte[] bytes = "AB0123?9876\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();

                System.out.println(msg);
            }

            // 4.2.5
            @Test
            void smartConnectShortcut() throws Exception {
                byte[] bytes = "\r|\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();

                System.out.println(msg);
            }

            // 4.2.7
            @Test
            void directCommandAccess1() throws Exception {
                byte[] bytes = "@2102\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 4.2.7
            @Test
            void directCommandAccess2() throws Exception {
                // TODO: this should be the same as the @above but that is not yet implemented
                byte[] bytes = "~2102\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

        }


        // 4.2.9.1
        @Nested
        class PointToPointCommands {
            @Test
            void pointToPointCommandDirect() throws Exception {
                byte[] bytes = "\\0603002102D4\r".getBytes(StandardCharsets.UTF_8);
                cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void pointToPointCommandBridged() throws Exception {
                byte[] bytes = "\\06420903210289\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void pointToMultiPointCommandBridged() throws Exception {
                byte[] bytes = "\\05FF007A38004A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void calReplyNormal() throws Exception {
                byte[] bytes = "8902312E322E363620200A\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                requestContext = new RequestContext(false);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void calReplySmart() throws Exception {
                // TODO: seems like the checksum is wrong here???
                //byte[] bytes = "860593008902312E322E363620207F\r\n".getBytes(StandardCharsets.UTF_8);
                byte[] bytes = "860593008902312E322E36362020EC\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                requestContext = new RequestContext(false);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
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
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);

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
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void toManyRetransmissions() throws Exception {
                byte[] bytes = "g#".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void corruption() throws Exception {
                byte[] bytes = "g$".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void desync() throws Exception {
                byte[] bytes = "g%".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void tooLong() throws Exception {
                byte[] bytes = "g'".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                ReplyOrConfirmation msg = ReplyOrConfirmation.staticParse(readBufferByteBased, cBusOptions, requestContext);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }

        // 7.3
        @Test
        void StandardFormatStatusReply1() throws Exception {
            byte[] bytes = "D8380068AA0140550550001000000014000000000000000000CF\r\n".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            requestContext = new RequestContext(true);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        @Test
        void StandardFormatStatusReply2() throws Exception {
            byte[] bytes = "D838580000000000000000000000000000000000000000000098\r\n".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            requestContext = new RequestContext(true);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        @Test
        void StandardFormatStatusReply3() throws Exception {
            byte[] bytes = "D638B000000000FF00000000000000000000000000000043\r\n".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            requestContext = new RequestContext(true);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        // 7.4
        @Test
        void ExtendedFormatStatusReply1() throws Exception {
            byte[] bytes = "F9073800AAAA000095990000000055550000000000005555555548\r\n".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            requestContext = new RequestContext(true);
            cBusOptions = new CBusOptions(false, false, false, true, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        @Test
        void ExtendedFormatStatusReply2() throws Exception {
            byte[] bytes = "F907380B0000000000005555000000000000000000000000000013\r\n".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            requestContext = new RequestContext(true);
            cBusOptions = new CBusOptions(false, false, false, true, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        @Test
        void ExtendedFormatStatusReply3() throws Exception {
            byte[] bytes = "F70738160000000000000000000000000000000000000000B4\r\n".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            requestContext = new RequestContext(true);
            cBusOptions = new CBusOptions(false, false, false, true, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatus() throws Exception {
                byte[] bytes = "\\05FF007A38004A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply1() throws Exception {
                byte[] bytes = "D83800A8AA02000000000000000000000000000000000000009C\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                requestContext = new RequestContext(true);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply2() throws Exception {
                byte[] bytes = "D838580000000000000000000000000000000000000000000098\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply3() throws Exception {
                byte[] bytes = "D638B0000000000000000000000000000000000000000042\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply4() throws Exception {
                // TODO: the command header seems wrong as it is missing a byte
                //byte[] bytes = "86999900F8003800A8AA0200000000000000000000000000000000000000C4\r\n".getBytes(StandardCharsets.UTF_8);
                byte[] bytes = "86999900F9003800A8AA0200000000000000000000000000000000000000C3\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, true, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }


            @Test
            void LightningStatusReply5() throws Exception {
                // TODO: the command header seems wrong as it is missing a byte
                // byte[] bytes = "86999900F800385800000000000000000000000000000000000000000000C0\r\n".getBytes(StandardCharsets.UTF_8);
                byte[] bytes = "86999900F900385800000000000000000000000000000000000000000000BF\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, true, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void LightningStatusReply6() throws Exception {
                // TODO: wrong checksum in this example???
                // TODO: the command header seems wrong as it is missing a byte
                //byte[] bytes = "86999900F60038B000000000000000000000000000000000000000008F\r\n".getBytes(StandardCharsets.UTF_8);
                byte[] bytes = "86999900F70038B0000000000000000000000000000000000000000069\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, true, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // TODO: due the usage of reserved we lose bits here so we need to fix that
            @Disabled("TODO: due the usage of reserved we lose bits here so we need to fix that")
            @Test
            void Reply() throws Exception {
                byte[] bytes = "8604990082300328\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

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
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
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
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        // 10.2.1
        @Test
        void testParameterSet() throws Exception {
            byte[] bytes = "@A3470011\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        // 10.2.1
        @Test
        void testParameterSetObsolete() throws Exception {
            byte[] bytes = "A3470011\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/C-Bus%20Quick%20Start%20Guide.pdf
    @Nested
    class CBusQuickStartGuideTest {

        // 4.3
        @Test
        void checksums() throws Exception {
            byte[] bytes = "\\0538007988C2g\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        //5.2
        @Nested
        class PCI_Setup {

            // 5.2.1
            @Nested
            class MMIMessagesNotRequired {

                @Test
                void init() throws Exception {
                    byte[] bytes = "~~~\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void writeSomething() throws Exception {
                    byte[] bytes = "A3210038g\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void writeSomethingResponse() throws Exception {
                    byte[] bytes = "g.322100AD\r\n".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
                    requestContext = new RequestContext(false);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void writeSomething2() throws Exception {
                    byte[] bytes = "A3420002g\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void writeSomethingResponse2() throws Exception {
                    byte[] bytes = "g.3242008C\r\n".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    requestContext = new RequestContext(false);
                    cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void writeSomething3() throws Exception {
                    byte[] bytes = "A3300059g\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void writeSomethingResponse3() throws Exception {
                    byte[] bytes = "g.8600000032300018\r\n".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = C_BUS_OPTIONS_WITH_SRCHK;
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }

            // 5.2.2
            @Test
            void MMIMessagesRequired() throws Exception {
                byte[] bytes = "A3300079g\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }

        // 6
        @Nested
        class TransmittingCBusLightingControlCommands {
            // 6.1
            @Test
            void TransmitAnONCommand() throws Exception {
                byte[] bytes = "\\053800790842u\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 6.2
            @Test
            void TransmitAnOFFCommand() throws Exception {
                byte[] bytes = "\\0538000108BAu\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 6.3
            @Test
            void TransmitAnRampToLevelCommand() throws Exception {
                byte[] bytes = "\\0538005A08550Cu\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }

        // 7
        @Nested
        class ReceivingCBusLightingControlCommands {
            // 7.1
            @Test
            void ReceiveAnONCommand() throws Exception {
                byte[] bytes = "05003800790842\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 7.1
            @Test
            void ReceiveAnONCommandAlternative() throws Exception {
                byte[] bytes = "0500380100790841\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 7.2
            @Test
            void ReceiveAnOFFCommand() throws Exception {
                byte[] bytes = "050038000108BA\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 7.2
            @Test
            void ReceiveAnOFFCommandAlternative() throws Exception {
                byte[] bytes = "05003801000108B9\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 7.3
            @Test
            void ReceiveAnRampToLevelCommand() throws Exception {
                byte[] bytes = "050038005A08550C\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 7.3
            @Test
            void ReceiveAnRampToLevelCommandAlternative() throws Exception {
                byte[] bytes = "05003801005A08550B\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Disabled("Needs to be implemented")
            // 7.4
            @Nested
            class ReceivingOtherCommands {
                @Test
                void Case1() throws Exception {
                    // Test with nn not 00 or 01... they should be discarded
                    byte[] bytes = "05ss38nn....zz\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Case2() throws Exception {
                    // Test with nn not 00 or 01... they should be discarded
                    byte[] bytes = "05ss3800cc....zz\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Case2Alternative() throws Exception {
                    // Test with nn not 00 or 01... they should be discarded
                    byte[] bytes = "05ss380100cc....zz\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }

        }

        // 8
        @Nested
        class InterpretingTheMmi {
            @Test
            void BigMMI1() throws Exception {
                byte[] bytes = "D8380068AA0140550550001000000014000000000000000000CF\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void BigMMI2() throws Exception {
                byte[] bytes = "D838580000000000000000000000000000000000000000000098\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void BigMMI3() throws Exception {
                byte[] bytes = "D638B000000000FF00000000000000000000000000000043\r\n".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }

        // 9
        @Nested
        class Example {
            @Nested
            class ControlExamples {
                @Test
                void turnOnLight() throws Exception {
                    byte[] bytes = "\\053800792129g\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void turnOffLight() throws Exception {
                    byte[] bytes = "\\0538000121A1g\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void rampLight() throws Exception {
                    byte[] bytes = "\\0538000A217F19g\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }

            @Nested
            class MontiorExamples {
                @Test
                void onCommand() throws Exception {
                    byte[] bytes = "050B380079201F\r\n".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void offCommand() throws Exception {
                    byte[] bytes = "050B3800012097\r\n".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Ramp() throws Exception {
                    byte[] bytes = "050B38000220484E\r\n".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, true, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }
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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void IconBitmap() throws Exception {
                byte[] bytes = "\\053800A412080021\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void CompleteDynamicIcon() throws Exception {
                byte[] bytes = "\\053800A412080022\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Nested
            class ChineseTable {

                @Test
                void StartDynamicIcon() throws Exception {
                    byte[] bytes = "\\053800A401080020\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void IconHeader() throws Exception {
                    byte[] bytes = "\\053800A80104CA00130C0600\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void AppendDynamicIcon() throws Exception {
                    byte[] bytes = "\\053800A401080021\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void WriteIconBitmapData1() throws Exception {
                    byte[] bytes = "\\053800A80104AAF05500FF50\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void WriteIconBitmapData2() throws Exception {
                    byte[] bytes = "\\053800A801040000F0F00F00\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void useDynamicIcon() throws Exception {
                    byte[] bytes = "\\053800A401080022\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void displayIcon() throws Exception {
                    byte[] bytes = "\\053800A60102CA010013\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void switchElectricalLoadsBridged() throws Exception {
                byte[] bytes = "\\0356093879935A\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2005%20-%20C-Bus%20Security%20Application.pdf
    @Nested
    class SecurityApplicationsTest {

        //5.11.1
        @Nested
        class SecuritySystemEmitsAlarmOn {
            @Test
            void AlarmOnWrongPrio() throws Exception {
                byte[] bytes = "\\05D00079832F\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void AlarmOn() throws Exception {
                byte[] bytes = "\\85D0007983AF\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

        }

        //5.11.2
        @Test
        void Zone3Unsealed() throws Exception {
            byte[] bytes = "\\05D0000A860398\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

            assertMessageMatches(bytes, msg);
        }

        //5.11.3
        @Test
        void ZoneName() throws Exception {
            byte[] bytes = "\\05D000AD8D034B49544348452E2020202088\r".getBytes(StandardCharsets.UTF_8);
            ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
            cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
            CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
            assertThat(msg).isNotNull();
            System.out.println(msg);

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
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void ArmSecurityRemote() throws Exception {
                byte[] bytes = "\\039209D00AA2FFE7\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }
    }


    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2006%20-%20C-Bus%20Metering%20Application.pdf
    @Nested
    class MeteringApplicationsTest {

        //6.11.1
        @Nested
        class DeviceRequestsMeteringApplicationtoMeasureElectricity {
            @Test
            void LocalMeasurement() throws Exception {
                byte[] bytes = "\\05D100090120\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteMeasurement() throws Exception {
                byte[] bytes = "\\035609D10901C3\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

        }

        //6.11.2
        @Nested
        class MeterMeasurementDevicesendsElectricityUse {
            @Test
            void LocalMeasurement() throws Exception {
                byte[] bytes = "\\05D1000D810000DBF8C9\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteMeasurement() throws Exception {
                byte[] bytes = "\\033709D10D810000DBF88B\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2007%20-%20C-Bus%20Trigger%20Control%20Application.pdf
    @Nested
    class TriggerControlApplicationsTest {

        //7.12
        @Nested
        class Examples {
            @Test
            void LocalTrigger() throws Exception {
                byte[] bytes = "\\05CA0002250109\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteTrigger() throws Exception {
                byte[] bytes = "\\035609CA022501AC\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2008%20-%20C-Bus%20Enable%20Control%20Application.pdf
    @Nested
    class EnableControlApplicationsTest {

        //8.11
        @Nested
        class Examples {
            @Test
            void LocalTrigger() throws Exception {
                byte[] bytes = "\\05CB0002378275\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void RemoteTrigger() throws Exception {
                // TODO: seems like the checksum is wrong here again...
                //byte[] bytes = "\\035609CB02378216\r".getBytes(StandardCharsets.UTF_8);
                byte[] bytes = "\\035609CB02378218\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2009%20-%20C-Bus%20Temperature%20Broadcast%20Application.pdf
    @Nested
    class TemperatureBroadcastApplicationsTest {

        //9.11
        @Nested
        class Examples {

            @Test
            void temperatureBroadcast() throws Exception {
                byte[] bytes = "\\051900020564\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2010%20-%20C-Bus%20Ventilation%20Application.pdf
    @Nested
    class VentilationApplicationTest {
        // TODO: no tests described here but it should work by adjusting the values from Lightning...
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2009%20-%20C-Bus%20Temperature%20Control%20Application.pdf
    @Nested
    class AccessControlApplicationsTest {

        //9.11
        @Nested
        class Examples {

            @Test
            void validAccessRequest() throws Exception {
                byte[] bytes = "\\05D500A4010300017D\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void closeAccessPoint() throws Exception {
                byte[] bytes = "\\05D5000201FF24\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void lockAccessPoint() throws Exception {
                byte[] bytes = "\\05D5000AFFFF1E\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void lockAccessPointRemote() throws Exception {
                byte[] bytes = "\\039209D50AFFFF85\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2021%20-%20C-Bus%20Media%20Transport%20Control%20Application.pdf
    @Nested
    class MediaTransportControlApplicationsTest {
        // TODO: no tests described here
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2023%20-%20C-Bus%20Clock%20and%20Timekeeping%20Application.pdf
    @Nested
    class ClockAndTimekeeping {

        //23.13
        @Nested
        class Examples {

            @Test
            void outputATimeCommand() throws Exception {
                byte[] bytes = "\\05DF000D010A2B1700C2\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void outputADateCommand() throws Exception {
                byte[] bytes = "\\05DF000E0207D502190411\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Disabled("This example just seems plain wrong... First of all there is no command for 0x10 defined and the argument which should be it 0x11 requires a argument of 0x03... So either documenation wrong or the example")
            @Test
            void outputARequestRefreshCommand() throws Exception {
                byte[] bytes = "\\05DF00100C\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void outputARequestRefreshCommandFixedQuestionMark() throws Exception {
                byte[] bytes = "\\05DF001103\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }
        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2024%20-%20C-Bus%20Telephony%20Application.pdf
    @Nested
    class Telephony {

        //24.11
        @Nested
        class Examples {

            @Test
            void LineOnHook() throws Exception {
                byte[] bytes = "\\05E000090111\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Disabled("Again it seems as this command is just wrong... there is no command definition for 2C")
            @Test
            void LineOffHook() throws Exception {
                byte[] bytes = "\\05E0002C020230333935323734333231FD\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            @Test
            void LineOffHookFixedQuestionMark() throws Exception {
                byte[] bytes = "\\05E000AC02013033393532373433323168\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, true);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

        }
    }

    // from: https://updates.clipsal.com/ClipsalSoftwareDownload/DL/downloads/OpenCBus/Chapter%2034%20-%20C-Bus%20Error%20Reporting%20Application.pdf
    @Nested
    class ErrorReporting {

        //34.13
        @Nested
        class Examples {

            // 34.13.1
            @Test
            void AllOk() throws Exception {
                byte[] bytes = "\\05CE0015FF20DE0000\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }

            // 34.13.2
            @Test
            void MinorFailure() throws Exception {
                byte[] bytes = "\\05CE0015882A6721B4\r".getBytes(StandardCharsets.UTF_8);
                ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                assertThat(msg).isNotNull();
                System.out.println(msg);

                assertMessageMatches(bytes, msg);
            }


            // 34.13.3
            @Nested
            class GeneralFailureWhichGetsAcknowledged {
                @Test
                void Reporting() throws Exception {
                    byte[] bytes = "\\05CE00159023426633\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void Acknowledge() throws Exception {
                    byte[] bytes = "\\05CE00259033426633\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }

            @Nested
            class LatchedExtremeFailureWhichGetsCleared {
                @Test
                void mostRecent() throws Exception {
                    byte[] bytes = "\\05CE001569E1FE0100\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void mostSevere() throws Exception {
                    byte[] bytes = "\\05CE001569CCFE0102\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void clearMostSevere() throws Exception {
                    byte[] bytes = "\\05CE003569C9FE0102\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }

                @Test
                void newError() throws Exception {
                    byte[] bytes = "\\05CE001569E9FE0100\r".getBytes(StandardCharsets.UTF_8);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
                    CBusMessage msg = CBusMessage.staticParse(readBufferByteBased, false, requestContext, cBusOptions);
                    assertThat(msg).isNotNull();
                    System.out.println(msg);

                    assertMessageMatches(bytes, msg);
                }
            }
        }
    }
}
