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
package org.apache.plc4x.java.opcua;

import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.test.manual.ManualTest;

import java.math.BigDecimal;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;

public class ManualS71500NewFWDriverTest extends ManualTest {

    public ManualS71500NewFWDriverTest(String connectionString) {
        super(connectionString, true, false, true, true, 100);
    }

    public static void main(String[] args) throws Exception {
        boolean testArrays = true;
        ManualS71500NewFWDriverTest test = new ManualS71500NewFWDriverTest("opcua://192.168.23.28:4840");
        test.addTestCase(/*"g_b1",*/            "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_b1\"",                   new PlcBOOL(true));
        test.addTestCase(/*"g_b8",*/            "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_b8\"",	                new PlcBYTE(0xAB));
        test.addTestCase(/*"g_s8",*/            "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_s8\"",	                new PlcSINT(-12));
        test.addTestCase(/*"g_u8",*/            "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_u8\"",	                new PlcUSINT(250));
        test.addTestCase(/*"g_b16",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_b16\"",	                new PlcWORD(0xBEEF));
        test.addTestCase(/*"g_s16",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_s16\"",                  new PlcINT(-1234));
        test.addTestCase(/*"g_u16",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_u16\"",                  new PlcUINT(54321));
        test.addTestCase(/*"g_b32",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_b32\"",                  new PlcDWORD(0xDEADBEEFL));
        test.addTestCase(/*"g_s32",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_s32\"",                  new PlcDINT(-12345678));
        test.addTestCase(/*"g_u32",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_u32\"",                  new PlcUDINT(305419896));
        test.addTestCase(/*"g_b64",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_b64\"",                  new PlcLWORD(0x0123_4567_89AB_CDEFL));
        test.addTestCase(/*"g_s64",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_s64\"",                  new PlcLINT(-9223372036854770000L));
        test.addTestCase(/*"g_u64",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_u64\"",                  new PlcULINT(new BigDecimal("18446744073709551000")));
        test.addTestCase(/*"g_r32",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_r32\"",                  new PlcREAL(3.14159));
        test.addTestCase(/*"g_r64",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_r64\"",                  new PlcLREAL(2.71828182845905));
        test.addTestCase(/*"g_tim",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_tim\";TIME",             new PlcTIME(2500)); // Is returned as Int32
        test.addTestCase(/*"g_dat",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_dat\";DATE",             new PlcDATE(LocalDate.of(2025, 11, 12))); // Is returned as UInt16
        test.addTestCase(/*"g_timoday",*/       "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_timoday\";TIME_OF_DAY", 	new PlcTIME_OF_DAY(LocalTime.of(14, 33, 21, 250000000))); // Is returned as UInt32
//        test.addTestCase(/*"g_dattim",*/        "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_dattim\";LTIME",         new PlcDATE_AND_LTIME(LocalDateTime.of(2025, 11, 12, 14, 33, 21, 500_000_000))); // TODO: Getting a class cast error, because OpcuaMessageResponse cannot be cast to OpcuaAPU
        test.addTestCase(/*"g_str",*/           "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_str\"",                  new PlcSTRING("Hello PLC4X"));
        test.addTestCase(/*"g_wstr",*/          "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_wstr\"",                 new PlcWSTRING("Grüße von PLC4X"));
        if(testArrays) {
            test.addTestCase(/*"g_arrBool",*/       "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrBool\"", new PlcList(List.of(
                new PlcBOOL(true), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(true),
                new PlcBOOL(false), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(false))
            ));
            test.addTestCase(/*"g_arrByte",*/       "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrByte\"", new PlcList(List.of(
                new PlcBYTE(0xDE), new PlcBYTE(0xAD), new PlcBYTE(0xBE), new PlcBYTE(0xEF),
                new PlcBYTE(0x12), new PlcBYTE(0x34), new PlcBYTE(0x56), new PlcBYTE(0x78))
            ));
            test.addTestCase(/*"g_arrInt",*/        "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrInt\"", new PlcList(List.of(
                new PlcINT(-3), new PlcINT(-1), new PlcINT(0), new PlcINT(1), new PlcINT(3))
            ));
            test.addTestCase(/*"g_arrUInt",*/       "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrUInt\"", new PlcList(List.of(
                new PlcUINT(1), new PlcUINT(10), new PlcUINT(100), new PlcUINT(1000), new PlcUINT(10000))
            ));
            test.addTestCase(/*"g_arrDInt",*/       "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrDInt\"", new PlcList(List.of(
                new PlcDINT(-1000), new PlcDINT(0), new PlcDINT(1000), new PlcDINT(2000000))
            ));
            test.addTestCase(/*"g_arrUDInt",*/       "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrUDInt\"", new PlcList(List.of(
                new PlcUDINT(0), new PlcUDINT(1), new PlcUDINT(0xFFFF), new PlcUDINT(0x12345678))
            ));
            test.addTestCase(/*"g_arrLReal",*/      "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrLReal\"", new PlcList(List.of(
                new PlcLREAL(1.5), new PlcLREAL(-2.0), new PlcLREAL(0.125))
            ));
            test.addTestCase(/*"g_arrTime",*/       "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrTime\";TIME", new PlcList(List.of(
                new PlcTIME(Duration.ofMillis(10)), new PlcTIME(Duration.ofSeconds(1)), new PlcTIME(Duration.ofSeconds(10)))
            ));
            test.addTestCase(/*"g_arrString",*/     "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrString\"", new PlcList(List.of(
                new PlcSTRING("alpha"), new PlcSTRING("beta"), new PlcSTRING("gamma"))
            ));
            test.addTestCase(/*"g_arrWString",*/     "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_arrWString\"", new PlcList(List.of(
                new PlcWSTRING("Äpfel"), new PlcWSTRING("Öl"))
            ));
            test.addTestCase(/*"g_matI16_2x3",*/    "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_matI16_2x3\"", new PlcList(List.of(
                new PlcList(List.of(
                    new PlcINT(10), new PlcINT(11), new PlcINT(12)
                )),
                new PlcList(List.of(
                    new PlcINT(-10), new PlcINT(-11), new PlcINT(-12)
                )))
            ));
            test.addTestCase(/*"g_matR32_3x2",*/    "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_matR32_3x2\"", new PlcList(List.of(
                new PlcList(List.of(
                    new PlcREAL(1.0), new PlcREAL(1.5)
                )),
                new PlcList(List.of(
                    new PlcREAL(2.0), new PlcREAL(2.5)
                )),
                new PlcList(List.of(
                    new PlcREAL(3.0), new PlcREAL(3.5)
                )))
            ));
            test.addTestCase(/*"g_cubeU16_2x2x2",*/ "ns=3;s=\"OPC_UA_DB\".\"OPC Data\".\"g_cubeU16_2x2x2\"", new PlcList(List.of(
                new PlcList(List.of(
                    new PlcList(List.of(
                        new PlcUINT(1), new PlcUINT(2)
                    )),
                    new PlcList(List.of(
                        new PlcUINT(3), new PlcUINT(4)
                    ))
                )),
                new PlcList(List.of(
                    new PlcList(List.of(
                        new PlcUINT(5), new PlcUINT(6)
                    )),
                    new PlcList(List.of(
                        new PlcUINT(7), new PlcUINT(8)
                    ))
                )))
            ));
        }

        long start = System.currentTimeMillis();
        test.run();
        long end = System.currentTimeMillis();
        System.out.printf("Finished in %d ms", end - start);
    }

}
