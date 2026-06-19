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
package org.apache.plc4x.java.s7.manual;

import org.apache.plc4x.java.api.authentication.PlcNullAuthentication;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.testutils.manual.BasicPlcTest;

import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;

public class ManualWallS7300DriverTest extends BasicPlcTest {

    /*
     * Test program code on the PLC with the test-data.
     *
     * Located in "main"
     *
     *
     */

    public ManualWallS7300DriverTest(String connectionString) {
        super(connectionString, new PlcNullAuthentication(), true, true, true, true, 100);
    }

    public static void main(String[] args) throws Exception {
        boolean testArrays = false;
        ManualWallS7300DriverTest test = new ManualWallS7300DriverTest("s7://192.168.24.60?remote-rack=0&remote-slot=1");//?log.audit-log-file=ManualWallS7300DriverTest-audit.log
        test.addTestCase(/*"g_b1",*/            "%DB42:0.0:BOOL",           new PlcBOOL(true));
        test.addTestCase(/*"g_b8",*/            "%DB42:1.0:BYTE",	        new PlcBYTE(0xAB));
        test.addTestCase(/*"g_b16",*/           "%DB42:2.0:WORD",	        new PlcWORD(0xBEEF));
        test.addTestCase(/*"g_s16",*/           "%DB42:4.0:INT",            new PlcINT(-1234));
        test.addTestCase(/*"g_b32",*/           "%DB42:6.0:DWORD",          new PlcDWORD(0xDEADBEEFL));
        test.addTestCase(/*"g_s32",*/           "%DB42:10.0:DINT",          new PlcDINT(-12345678));
        test.addTestCase(/*"g_r32",*/           "%DB42:14.0:REAL",          new PlcREAL(3.14159));
        test.addTestCase(/*"g_tim",*/           "%DB42:18.0:TIME",          new PlcTIME(2500));
        test.addTestCase(/*"g_dat",*/           "%DB42:22.0:DATE",          new PlcDATE(LocalDate.of(2025, 11, 12)));
        test.addTestCase(/*"g_timoday",*/       "%DB42:24.0:TIME_OF_DAY",	new PlcTIME_OF_DAY(LocalTime.of(14, 33, 21, 250000000)));
        test.addTestCase(/*"g_dattim",*/        "%DB42:28.0:DATE_AND_TIME", new PlcDATE_AND_TIME(LocalDateTime.of(2025, 11, 12, 14, 33, 21)));
        test.addTestCase(/*"g_str",*/           "%DB42:36.0:STRING(40)",    new PlcSTRING("Hello PLC4X"));
        if(testArrays) {
            test.addTestCase(/*"g_arrBool",*/       "%DB42:78.0:BOOL[8]", new PlcList(List.of(
                new PlcBOOL(true), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(false),
                new PlcBOOL(false), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(false))
            ));
            test.addTestCase(/*"g_arrByte",*/       "%DB42:80.0:BYTE[8]", new PlcList(List.of(
                new PlcBYTE(0xDE), new PlcBYTE(0xAD), new PlcBYTE(0xBE), new PlcBYTE(0xEF),
                new PlcBYTE(0x12), new PlcBYTE(0x34), new PlcBYTE(0x56), new PlcBYTE(0x78))
            ));
            test.addTestCase(/*"g_arrInt",*/        "%DB4:88.0:INT[5]", new PlcList(List.of(
                new PlcINT(-3), new PlcINT(-1), new PlcINT(0), new PlcINT(1), new PlcINT(3))
            ));
            test.addTestCase(/*"g_arrDInt",*/       "%DB42:98.0:DINT[4]", new PlcList(List.of(
                new PlcDINT(-1000), new PlcDINT(0), new PlcDINT(1000), new PlcDINT(2000000))
            ));
            test.addTestCase(/*"g_arrTime",*/       "%DB42:114.0:TIME[3]", new PlcList(List.of(
                new PlcTIME(Duration.ofMillis(10)), new PlcTIME(Duration.ofSeconds(1)), new PlcTIME(Duration.ofSeconds(10)))
            ));
            test.addTestCase(/*"g_arrString",*/     "%DB42:126.0:STRING(16)[3]", new PlcList(List.of(
                new PlcSTRING("alpha"), new PlcSTRING("beta"), new PlcSTRING("gamma"))
            ));
            test.addTestCase(/*"g_matI16_2x3",*/    "%DB42:180.0:INT[2][3]", new PlcList(List.of(
                new PlcList(List.of(
                    new PlcINT(10), new PlcINT(11), new PlcINT(12)
                )),
                new PlcList(List.of(
                    new PlcINT(-10), new PlcINT(-11), new PlcINT(-12)
                )))
            ));
            test.addTestCase(/*"g_matR32_3x2",*/    "%DB42:192.0:REAL[3][2]", new PlcList(List.of(
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
            test.addTestCase(/*"g_cubeU16_2x2x2",*/ "%DB42:216.0:INT[2][2][2]", new PlcList(List.of(
                new PlcList(List.of(
                    new PlcList(List.of(
                        new PlcINT(1), new PlcINT(2)
                    )),
                    new PlcList(List.of(
                        new PlcINT(3), new PlcINT(4)
                    ))
                )),
                new PlcList(List.of(
                    new PlcList(List.of(
                        new PlcINT(5), new PlcINT(6)
                    )),
                    new PlcList(List.of(
                        new PlcINT(7), new PlcINT(8)
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
