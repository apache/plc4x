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

package org.apache.plc4x.java.profinet;

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.profinet.gsdml.ProfinetInterfaceSubmoduleItem;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.readwrite.utils.StaticHelper;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;

import static org.junit.jupiter.api.Assertions.assertEquals;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetCheckSumTests {


    @Test
    public void calculateChecksumTest()  {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST,
            true,
            false,
            false,
            IntegerEncoding.BIG_ENDIAN,
            CharacterEncoding.ASCII,
            FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0002, 0x0493),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            new DceRpc_ActivityUuid(0x0aa499a5L, 0x1df0, 0x11b2, new byte[]{(byte) 0x80, (byte) 0x00, (byte) 0xa9, (byte) 0xa2, (byte) 0x59, (byte) 0x4b, (byte) 0x74, (byte) 0x00}),
            0,
            0,
            DceRpc_Operation.CONNECT,
            new PnIoCm_Packet_Req(
                16696,
                16696,
                0,
                0,
                new ArrayList<PnIoCm_Block>())
        );
        try {
            assertEquals(
                (short) 0x4411,
                StaticHelper.calculateUdpChecksum(
                    new IpAddress(Hex.decodeHex("c0a85a01")),
                    new IpAddress(Hex.decodeHex("c0a85a80")),
                    50000,
                    34964,
                    packet.getLengthInBytes() + 8,
                    packet)
                );
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }

    }

    @Test
    public void calculateIpChecksumTest() throws DecoderException {
        int checksum = StaticHelper.calculateIPv4Checksum(
            425,
            0x4dc2,
            64,
            new IpAddress(Hex.decodeHex("c0a85a01")),
            new IpAddress(Hex.decodeHex("c0a85a80"))
        );

        assertEquals(
            checksum
            , (short) 0xb5af);
    }

    @Test
    public void calculateIpChecksumTest2() throws DecoderException {
        int checksum = StaticHelper.calculateIPv4Checksum(
            425,
            0x1e85,
            64,
            new IpAddress(Hex.decodeHex("c0a85a01")),
            new IpAddress(Hex.decodeHex("c0a85a80"))
        );

        assertEquals(
            checksum
            , (short) 0xe4ec);
    }



}
