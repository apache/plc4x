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
package org.apache.plc4x.java.umas.manual;

import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.test.manual.ManualTest;

import java.time.*;

public class ManualUmasDriverTest extends ManualTest {

    public ManualUmasDriverTest(String connectionString) {
        super(connectionString, true, true, true, true, 100);
    }

    public static void main(String[] args) throws Exception {
        String spsIp = "192.168.42.99";
        String connectionString = String.format("umas://%s", spsIp);
        ManualUmasDriverTest test = new ManualUmasDriverTest(connectionString);

        // =========================================================
        // 1) Scalars — all primitive PlcValueTypes that are supported on the UMAS device
        // =========================================================
        test.addTestCase("g_b1",      new PlcBOOL(true));
        test.addTestCase("g_b8",      new PlcBYTE(0xAB));
        test.addTestCase("g_b16",     new PlcWORD(0xBEEF));
        test.addTestCase("g_s16",     new PlcINT(-1234));
        test.addTestCase("g_u16",     new PlcUINT(54321));
        test.addTestCase("g_b32",     new PlcDWORD(0xDEADBEEFL));
        test.addTestCase("g_s32",     new PlcDINT(-12345678));
        test.addTestCase("g_u32",     new PlcUDINT(305_419_896L)); // 0x12345678
        test.addTestCase("g_r32",     new PlcREAL(3.14159f));
        test.addTestCase("g_tim",     new PlcTIME(java.time.Duration.parse("PT2.5S")));
        test.addTestCase("g_dat",     new PlcDATE(java.time.LocalDate.parse("2025-11-12")));
        test.addTestCase("g_timoday", new PlcTIME_OF_DAY(java.time.LocalTime.parse("14:33:21")));
        test.addTestCase("g_dattim",  new PlcDATE_AND_TIME(java.time.LocalDateTime.parse("2025-11-12T14:33:21")));
        test.addTestCase("g_str",     new PlcSTRING("Hello PLC4X"));

        long startMillis = System.currentTimeMillis();
        test.run();
        long endMillis = System.currentTimeMillis();
        System.out.println("Test executed in " + (endMillis - startMillis) + "ms");
    }

}
