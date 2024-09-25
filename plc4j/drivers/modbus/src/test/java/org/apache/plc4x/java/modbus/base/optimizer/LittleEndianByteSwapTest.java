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

package org.apache.plc4x.java.modbus.base.optimizer;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagCoil;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHandler;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.tcp.context.ModbusTcpContext;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public class LittleEndianByteSwapTest {

    @Test
    public void testLittleEndianByteSwap() throws Exception {
        Map<String, String[]> input = new LinkedHashMap<>();
        input.put("variable0",     new String[]{"holding-register:1:REAL",   "22.009401"});
        input.put("variable1",     new String[]{"holding-register:5:REAL",   "25084.002"});
        input.put("variable2",     new String[]{"holding-register:9:REAL",   "-159.22"});
        input.put("variable3",     new String[]{"holding-register:13:REAL",  "21.7556"});
        input.put("variable4",     new String[]{"holding-register:17:REAL",  "28798.0"});
        input.put("variable5",     new String[]{"holding-register:21:REAL",  "-159.01"});
        input.put("variable6",     new String[]{"holding-register:25:REAL",  "21.471199"});
        input.put("variable7",     new String[]{"holding-register:29:REAL",  "28330.0"});
        input.put("variable8",     new String[]{"holding-register:33:REAL",  "-158.75"});
        input.put("variable9",     new String[]{"holding-register:37:REAL",  "0.69699997"});
        input.put("variable10",    new String[]{"holding-register:41:REAL",  "32.899994"});
        input.put("variable11",    new String[]{"holding-register:45:REAL",  "0.019"});
        input.put("variable12",    new String[]{"holding-register:49:REAL",  "31.600006"});
        input.put("variable13",    new String[]{"holding-register:53:REAL",  "0.683625"});
        input.put("variable14",    new String[]{"holding-register:57:REAL",  "4.9614224"});
        input.put("variable15",    new String[]{"holding-register:61:REAL",  "0.72375"});
        input.put("variable16",    new String[]{"holding-register:65:REAL",  "0.729"});
        input.put("variable17",    new String[]{"holding-register:69:REAL",  "1.9530001"});
        input.put("variable18",    new String[]{"holding-register:73:REAL",  "0.00125"});
        input.put("variable19",    new String[]{"holding-register:77:REAL",  "0.175"});
        input.put("variable20",    new String[]{"holding-register:81:REAL",  "10.285"});
        input.put("variable21",    new String[]{"holding-register:85:REAL",  "425.439"});
        input.put("variable22",    new String[]{"holding-register:89:REAL",  "0.84125"});
        input.put("variable23",    new String[]{"holding-register:93:REAL",  "0.5225"});
        input.put("variable24",    new String[]{"holding-register:97:REAL",  "0.73700005"});
        input.put("variable25",    new String[]{"holding-register:101:REAL", "11.6865"});
        input.put("variable26",    new String[]{"holding-register:105:REAL", "0.0"});
        input.put("variable27",    new String[]{"holding-register:109:REAL", "0.0"});
        input.put("variable28",    new String[]{"holding-register:113:REAL", "0.7339"});
        input.put("variable29",    new String[]{"holding-register:117:REAL", "0.7232"});
        input.put("variable30",    new String[]{"holding-register:121:REAL", "40.16"});
        input.put("variable31",    new String[]{"holding-register:125:REAL", "-0.008889999"});
        input.put("variable32",    new String[]{"holding-register:129:REAL", "-0.010159999"});
        input.put("variable33",    new String[]{"holding-register:133:REAL", "0.7351"});
        input.put("variable34",    new String[]{"holding-register:137:REAL", "0.7208"});
        input.put("variable35",    new String[]{"holding-register:141:REAL", "0.56"});
        input.put("variable36",    new String[]{"holding-register:145:REAL", "-0.012700001"});
        input.put("variable37",    new String[]{"holding-register:149:REAL", "-0.01651"});
        input.put("variable38",    new String[]{"holding-register:153:REAL", "0.7338"});
        input.put("variable39",    new String[]{"holding-register:157:REAL", "0.7296"});
        input.put("variable40",    new String[]{"holding-register:161:REAL", "0.72"});
        input.put("variable41",    new String[]{"holding-register:165:REAL", "-0.008889999"});
        input.put("variable42",    new String[]{"holding-register:169:REAL", "-0.020319998"});
        input.put("variable43",    new String[]{"holding-register:173:REAL", "0.424"});
        input.put("variable44",    new String[]{"holding-register:177:REAL", "0.3768"});
        input.put("variable45",    new String[]{"holding-register:181:REAL", "0.0"});
        input.put("variable46",    new String[]{"holding-register:185:REAL", "-0.01143"});
        input.put("variable47",    new String[]{"coil:1:BOOL",               "false"});
        input.put("variable48",    new String[]{"coil:3:BOOL",               "false"});
        input.put("variable49",    new String[]{"coil:5:BOOL",               "true"});
        input.put("variable50",    new String[]{"coil:7:BOOL",               "false"});
        input.put("variable51",    new String[]{"coil:9:BOOL",               "false"});
        input.put("variable52",    new String[]{"coil:11:BOOL",              "false"});
        input.put("variable53",    new String[]{"coil:13:BOOL",              "false"});
        input.put("variable54",    new String[]{"coil:15:BOOL",              "true"});
        input.put("variable55",    new String[]{"coil:17:BOOL",              "false"});
        input.put("variable56",    new String[]{"coil:19:BOOL",              "false"});
        input.put("variable57",    new String[]{"coil:21:BOOL",              "false"});
        input.put("variable58",    new String[]{"coil:23:BOOL",              "true"});
        input.put("variable59",    new String[]{"coil:25:BOOL",              "false"});
        input.put("variable60",    new String[]{"coil:27:BOOL",              "true"});
        input.put("variable61",    new String[]{"coil:29:BOOL",              "false"});
        input.put("variable62",    new String[]{"coil:31:BOOL",              "false"});
        input.put("variable63",    new String[]{"coil:33:BOOL",              "false"});
        input.put("variable64",    new String[]{"coil:35:BOOL",              "false"});
        input.put("variable65",    new String[]{"coil:37:BOOL",              "false"});
        input.put("variable66",    new String[]{"coil:39:BOOL",              "false"});
        input.put("variable67",    new String[]{"coil:41:BOOL",              "false"});
        input.put("variable68",    new String[]{"coil:43:BOOL",              "false"});
        input.put("variable69",    new String[]{"coil:45:BOOL",              "false"});
        input.put("variable70",    new String[]{"coil:47:BOOL",              "true"});
        input.put("variable71",    new String[]{"coil:49:BOOL",              "false"});
        input.put("variable72",    new String[]{"coil:51:BOOL",              "true"});
        input.put("variable73",    new String[]{"coil:53:BOOL",              "false"});
        input.put("variable74",    new String[]{"coil:55:BOOL",              "false"});
        input.put("variable75",    new String[]{"coil:57:BOOL",              "false"});
        input.put("variable76",    new String[]{"coil:59:BOOL",              "false"});
        input.put("variable77",    new String[]{"coil:61:BOOL",              "false"});
        input.put("variable78",    new String[]{"coil:63:BOOL",              "false"});
        input.put("variable79",    new String[]{"coil:65:BOOL",              "false"});
        input.put("variable80",    new String[]{"coil:67:BOOL",              "false"});
        input.put("variable81",    new String[]{"coil:69:BOOL",              "false"});
        input.put("variable82",    new String[]{"coil:71:BOOL",              "false"});
        input.put("variable83",    new String[]{"coil:73:BOOL",              "false"});
        input.put("variable84",    new String[]{"coil:75:BOOL",              "false"});
        input.put("variable85",    new String[]{"coil:77:BOOL",              "false"});
        input.put("variable86",    new String[]{"coil:79:BOOL",              "false"});
        input.put("variable87",    new String[]{"coil:81:BOOL",              "false"});
        input.put("variable88",    new String[]{"coil:83:BOOL",              "false"});
        input.put("variable89",    new String[]{"coil:85:BOOL",              "false"});
        input.put("variable90",    new String[]{"coil:87:BOOL",              "false"});
        input.put("variable91",    new String[]{"coil:89:BOOL",              "false"});
        input.put("variable92",    new String[]{"coil:91:BOOL",              "false"});
        input.put("variable93",    new String[]{"coil:93:BOOL",              "false"});
        input.put("variable94",    new String[]{"coil:95:BOOL",              "false"});
        input.put("variable95",    new String[]{"coil:97:BOOL",              "false"});
        input.put("variable96",    new String[]{"coil:99:BOOL",              "false"});
        input.put("variable97",    new String[]{"coil:101:BOOL",             "false"});
        input.put("variable98",    new String[]{"coil:103:BOOL",             "false"});
        input.put("variable99",    new String[]{"coil:105:BOOL",             "false"});
        input.put("variable100",   new String[]{"coil:107:BOOL",             "false"});
        input.put("variable101",   new String[]{"coil:109:BOOL",             "false"});
        input.put("variable102",   new String[]{"coil:111:BOOL",             "false"});
        input.put("variable103",   new String[]{"coil:113:BOOL",             "false"});
        input.put("variable104",   new String[]{"coil:115:BOOL",             "false"});
        input.put("variable105",   new String[]{"coil:117:BOOL",             "false"});
        input.put("variable106",   new String[]{"coil:119:BOOL",             "false"});
        input.put("variable107",   new String[]{"coil:121:BOOL",             "false"});
        input.put("variable108",   new String[]{"coil:123:BOOL",             "false"});
        input.put("variable109",   new String[]{"coil:125:BOOL",             "false"});
        input.put("variable110",   new String[]{"coil:127:BOOL",             "false"});
        input.put("variable111",   new String[]{"coil:129:BOOL",             "false"});
        input.put("variable112",   new String[]{"coil:131:BOOL",             "false"});
        input.put("variable113",   new String[]{"coil:133:BOOL",             "false"});
        input.put("variable114",   new String[]{"coil:135:BOOL",             "false"});
        input.put("variable115",   new String[]{"coil:137:BOOL",             "false"});
        input.put("variable116",   new String[]{"coil:139:BOOL",             "false"});
        input.put("variable117",   new String[]{"coil:141:BOOL",             "false"});
        input.put("variable118",   new String[]{"coil:143:BOOL",             "false"});
        input.put("variable119",   new String[]{"coil:145:BOOL",             "false"});
        input.put("variable120",   new String[]{"coil:147:BOOL",             "false"});
        input.put("variable121",   new String[]{"coil:149:BOOL",             "false"});
        input.put("variable122",   new String[]{"coil:151:BOOL",             "false"});
        input.put("variable123",   new String[]{"coil:153:BOOL",             "false"});
        input.put("variable124",   new String[]{"coil:155:BOOL",             "false"});
        input.put("variable125",   new String[]{"coil:157:BOOL",             "false"});
        input.put("variable126",   new String[]{"coil:159:BOOL",             "false"});
        input.put("variable127",   new String[]{"coil:161:BOOL",             "false"});
        input.put("variable128",   new String[]{"coil:163:BOOL",             "false"});
        input.put("variable129",   new String[]{"coil:165:BOOL",             "false"});
        input.put("variable130",   new String[]{"coil:167:BOOL",             "false"});
        input.put("variable131",   new String[]{"coil:169:BOOL",             "false"});
        input.put("variable132",   new String[]{"coil:171:BOOL",             "false"});
        input.put("variable133",   new String[]{"coil:173:BOOL",             "false"});
        input.put("variable134",   new String[]{"coil:175:BOOL",             "false"});
        PlcReader mockPlcReader = Mockito.mock(PlcReader.class);
        PlcTagHandler modbusTagHandler = new ModbusTagHandler();
        PlcReadRequest.Builder builder = new DefaultPlcReadRequest.Builder(mockPlcReader, modbusTagHandler);
        for (String name : input.keySet()) {
            String[] data = input.get(name);
            builder.addTagAddress(name, data[0]);
        }
        PlcReadRequest readRequest = builder.build();

        ModbusOptimizer sut = new ModbusOptimizer();

        ModbusTcpContext mockContext = Mockito.mock(ModbusTcpContext.class);
        Mockito.when(mockContext.getByteOrder()).thenReturn(ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP);
        Mockito.when(mockContext.getMaxCoilsPerRequest()).thenReturn(2000);
        Mockito.when(mockContext.getMaxRegistersPerRequest()).thenReturn(125);

        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Do the first part of the optimizer ... split up into multiple requests ...
        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////

        List<PlcReadRequest> optimizedReadRequests = sut.processReadRequest(readRequest, mockContext);

        // Validate the results of this first step.

        Assertions.assertNotNull(optimizedReadRequests);
        Assertions.assertEquals(3, optimizedReadRequests.size());

        // Check the expected first request (Coil)
        PlcReadRequest firstRequest = optimizedReadRequests.get(0);
        Assertions.assertEquals(1, firstRequest.getNumberOfTags());
        Assertions.assertEquals("coils0", firstRequest.getTagNames().stream().findFirst().orElseThrow());
        PlcTag coilsTag = firstRequest.getTag("coils0");
        Assertions.assertInstanceOf(ModbusTagCoil.class, coilsTag);
        ModbusTagCoil coil = (ModbusTagCoil) coilsTag;
        Assertions.assertEquals(0, coil.getAddress());
        Assertions.assertEquals(175, coil.getNumberOfElements());

        // Check the expected second request (Register)
        PlcReadRequest secondRequest = optimizedReadRequests.get(1);
        Assertions.assertEquals(1, secondRequest.getNumberOfTags());
        Assertions.assertEquals("registers0", secondRequest.getTagNames().stream().findFirst().orElseThrow());
        PlcTag registers0Tag = secondRequest.getTag("registers0");
        Assertions.assertInstanceOf(ModbusTagHoldingRegister.class, registers0Tag);
        ModbusTagHoldingRegister registers0 = (ModbusTagHoldingRegister) registers0Tag;
        Assertions.assertEquals(0, registers0.getAddress());
        Assertions.assertEquals(122, registers0.getNumberOfElements());

        // Check the expected third request (Register)
        PlcReadRequest thirdRequest = optimizedReadRequests.get(2);
        Assertions.assertEquals(1, thirdRequest.getNumberOfTags());
        Assertions.assertEquals("registers1", thirdRequest.getTagNames().stream().findFirst().orElseThrow());
        PlcTag registers1Tag = thirdRequest.getTag("registers1");
        Assertions.assertInstanceOf(ModbusTagHoldingRegister.class, registers1Tag);
        ModbusTagHoldingRegister registers1 = (ModbusTagHoldingRegister) registers1Tag;
        Assertions.assertEquals(124, registers1.getAddress());
        Assertions.assertEquals(62, registers1.getNumberOfElements());

        // Prepare the results as we got them on the wire.

        Map<PlcReadRequest, BaseOptimizer.SubResponse<PlcReadResponse>> readResponses = new HashMap<>();
        readResponses.put(firstRequest, new BaseOptimizer.SubResponse<>(
            new DefaultPlcReadResponse(firstRequest, Map.of(
                "coils0", new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcRawByteArray(Hex.decodeHex("3060480c00c084000000000000000000000000000000")))))));
        readResponses.put(secondRequest, new BaseOptimizer.SubResponse<>(
            new DefaultPlcReadResponse(secondRequest, Map.of(
                "registers0", new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcRawByteArray(Hex.decodeHex("134141b000000000f80146c3000000003852c31f000000000b7841ae00000000fc0046e000000000028fc31f00000000c50441ab00000000540046dd00000000c000c31e000000006e973f32000000009998420300000000a5e33c9b00000000ccd041fc00000000020c3f2f00000000c3f9409e0000000047ae3f39000000009fbe3f3a00000000fbe83ff900000000d70a3aa30000000033333e33000000008f5c412400000000b83143d4000000005c293f5700000000c28f3f0500000000ac093f3c00000000fbe7413a0000000000000000000000000000000000000000e0df3f3b0000000023a33f3900000000a3d74220")))))));
        readResponses.put(thirdRequest, new BaseOptimizer.SubResponse<>(
            new DefaultPlcReadResponse(thirdRequest, Map.of(
                "registers1", new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcRawByteArray(Hex.decodeHex("a75cbc11000000007620bc26000000002f833f3c0000000086593f38000000005c293f0f0000000013aabc50000000003ffbbc8700000000da513f3b00000000c7113f3a0000000051ec3f3800000000a75cbc11000000007620bca60000000016873ed900000000ebee3ec000000000000000000000000044e5bc3b")))))));

        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Process the responses
        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////

        PlcReadResponse readResponse = sut.processReadResponses(readRequest, readResponses, mockContext);

        // Check if there were no invalid items
        List<String> failedTags = readResponse.getTagNames().stream().filter(tagName -> readResponse.getResponseCode(tagName) != PlcResponseCode.OK).collect(Collectors.toList());
        failedTags.forEach(failedTag -> Assertions.fail("Field " + failedTag + "failed."));

        // Check if the returned values match the expected ones
        for (String name : input.keySet()) {
            String[] data = input.get(name);
            String readValue = readResponse.getString(name);
            Assertions.assertEquals(data[1], readValue);
        }
    }
    
}
