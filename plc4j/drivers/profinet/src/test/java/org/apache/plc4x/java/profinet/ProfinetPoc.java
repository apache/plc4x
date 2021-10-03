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
package org.apache.plc4x.java.profinet;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.readwrite.types.*;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;

public class ProfinetPoc {

    public static void main(String[] args) throws Exception {

        // Send Profinet IO Context Manager (PNIO-CM) Connection Request (UDP)
        DceRpc_Packet connectionRequest = new DceRpc_Packet_Req(
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE, 1, 0, 0, 0,0,0, DceRpc_Operation.CONNECT,
            new PnIoCm_Packet_Req(404, 404, 0, 404,
                new PnIoCm_Block[]{
                    new PnIoCm_Block_ArReq((short) 1, (short) 0, PnIoCm_ArType.IO_CONTROLLER,
                        new Uuid(Hex.decodeHex("654519352df3b6428f874371217c2b51")), 2,
                        new MacAddress(Hex.decodeHex("606D3C3DA9A3")),
                        new Uuid(Hex.decodeHex("dea000006c9711d1827100640008002a")),
                        false, false, false,
                        false, PnIoCm_CompanionArType.SINGLE_AR, false,
                        true, false, PnIoCm_State.ACTIVE,
                        600, 0x8892,
                        /*"plc4x-pn-master"*/"profinetxadriver4933"),
                    new PnIoCm_Block_IoCrReq((short) 1, (short) 0, PnIoCm_IoCrType.INPUT_CR,
                        0x0001, 0x8892, false, false,
                        false, false, PnIoCm_RtClass.RT_CLASS_2, 40,
                        0xBBF0, 128, 8, 1, 0, 0xffffffff,
                        3, 3, 0xC000,
                        new MacAddress(Hex.decodeHex("000000000000")),
                        new PnIoCm_IoCrBlockReqApi[] {
                            new PnIoCm_IoCrBlockReqApi(
                                new PnIoCm_IoDataObject[] {
                                    new PnIoCm_IoDataObject(0, 0x0001, 0),
                                    new PnIoCm_IoDataObject(0, 0x8000, 1),
                                    new PnIoCm_IoDataObject(0, 0x8001, 2),
                                    new PnIoCm_IoDataObject(0, 0x8002, 3),
                                    new PnIoCm_IoDataObject(1, 0x0001, 4)
                                },
                                    new PnIoCm_IoCs[] {
                                    new PnIoCm_IoCs(0x0001, 0x0001, 0x0019)
                                })
                        }),
                    new PnIoCm_Block_IoCrReq((short) 1, (short) 0, PnIoCm_IoCrType.OUTPUT_CR,
                        0x0002, 0x8892,  false, false,
                        false, false, PnIoCm_RtClass.RT_CLASS_2, 40,
                        0x8000, 128, 8, 1, 0, 0xffffffff,
                        3, 3, 0xC000,
                        new MacAddress(Hex.decodeHex("000000000000")),
                        new PnIoCm_IoCrBlockReqApi[]{
                            new PnIoCm_IoCrBlockReqApi(
                                new PnIoCm_IoDataObject[] {
                                    new PnIoCm_IoDataObject(0x0001, 0x0001, 0x0005)
                                },
                                new PnIoCm_IoCs[] {
                                    new PnIoCm_IoCs(0, 0x0001, 0),
                                    new PnIoCm_IoCs(0, 0x8000, 1),
                                    new PnIoCm_IoCs(0, 0x8001, 2),
                                    new PnIoCm_IoCs(0, 0x8002, 3),
                                    new PnIoCm_IoCs(1, 0x0001, 4)
                                })
                        }
                    ),
                    new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                        new PnIoCm_ExpectedSubmoduleBlockReqApi[] {
                            new PnIoCm_ExpectedSubmoduleBlockReqApi(0,
                                0x00000010, 0x00000000, new PnIoCm_Submodule[] {
                                    new PnIoCm_Submodule_NoInputNoOutputData(0x0001,
                                        0x00000001, false, false,
                                        false, false),
                                    new PnIoCm_Submodule_NoInputNoOutputData(0x8000,
                                        0x00000002, false, false,
                                        false, false),
                                    new PnIoCm_Submodule_NoInputNoOutputData(0x8001,
                                        0x00000003, false, false,
                                        false, false),
                                    new PnIoCm_Submodule_NoInputNoOutputData(0x8002,
                                        0x00000003, false, false,
                                        false, false)
                            })
                        }
                    ),
                    new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                        new PnIoCm_ExpectedSubmoduleBlockReqApi[] {
                            new PnIoCm_ExpectedSubmoduleBlockReqApi(1,
                                0x00000022, 0x00000000, new PnIoCm_Submodule[] {
                                new PnIoCm_Submodule_InputAndOutputData(0x0001, 0x00000010,
                                    false, false, false,
                                    false, 20, (short) 1, (short) 1,
                                    6, (short) 1, (short) 1)
                            })
                        }
                    ),
                    new PnIoCm_Block_AlarmCrReq((short) 1, (short) 0,
                        PnIoCm_AlarmCrType.ALARM_CR, 0x8892, false, false, 1, 3,
                        0x0000, 200, 0xC000, 0xA000)
                },404)
            );

        // Serialize the message
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(connectionRequest.getLengthInBytes());
        connectionRequest.getMessageIO().serialize(writeBuffer, connectionRequest);

        // Send the message as UDP telegram.
        InetAddress[] deviceAddress = InetAddress.getAllByName("192.168.24.31");
        DatagramPacket packet = new DatagramPacket(
            writeBuffer.getData(), writeBuffer.getData().length, deviceAddress[0], 34964);
        DatagramSocket datagramSocket = new DatagramSocket();
        datagramSocket.send(packet);

        // The PNIO_PS message seems to come in earlier than the connection response ...
        // Receive PNIO_PS (ProfiSafe Cyclic Data Unit)
        // Receive PNIO-CM Connection Response (UDP)


    }

}
