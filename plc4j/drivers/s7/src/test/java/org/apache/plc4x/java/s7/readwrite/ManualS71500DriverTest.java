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
package org.apache.plc4x.java.s7.readwrite;

import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.test.manual.ManualTest;

import java.time.*;

public class ManualS71500DriverTest extends ManualTest {

    public ManualS71500DriverTest(String connectionString) {
        super(connectionString, true, true, true, true, 100);
    }

    public static void main(String[] args) throws Exception {
        ManualS71500DriverTest test = new ManualS71500DriverTest("s7://192.168.24.83");
        test.addTestCase("%DB1:0.0:BOOL", new PlcBOOL(true));
        test.addTestCase("%DB1:1:BYTE", new PlcBYTE(42));
        test.addTestCase("%DB1:2:WORD", new PlcWORD(42424));
        test.addTestCase("%DB1:4:DWORD", new PlcDWORD(4242442424L));
        test.addTestCase("%DB1:8:LWORD", new PlcLWORD(4242442424424244242L));
        test.addTestCase("%DB1:16:SINT", new PlcSINT(-42));
        test.addTestCase("%DB1:17:USINT", new PlcUSINT(42));
        test.addTestCase("%DB1:18:INT", new PlcINT(-2424));
        test.addTestCase("%DB1:20:UINT", new PlcUINT(42424));
        test.addTestCase("%DB1:22:DINT", new PlcDINT(-242442424));
        test.addTestCase("%DB1:26:UDINT", new PlcUDINT(4242442424L));
        test.addTestCase("%DB1:30:LINT", new PlcLINT(-4242442424424244242L));
        test.addTestCase("%DB1:38:ULINT", new PlcULINT(4242442424424244242L));
        test.addTestCase("%DB1:46:REAL", new PlcREAL(3.141593F));
        test.addTestCase("%DB1:50:LREAL", new PlcLREAL(2.71828182846D));
        test.addTestCase("%DB1:58:CHAR", new PlcCHAR("H"));
        test.addTestCase("%DB1:60:WCHAR", new PlcWCHAR("w"));
        test.addTestCase("%DB1:62:STRING(10)", new PlcSTRING("hurz"));
        test.addTestCase("%DB1:318:WSTRING(10)", new PlcWSTRING("wolf"));
        test.addTestCase("%DB1:830:STRING", new PlcSTRING("hurz"));
        test.addTestCase("%DB1:1086:WSTRING", new PlcWSTRING("wolf"));
        test.addTestCase("%DB1:1598:TIME", new PlcTIME(Duration.parse("PT1.234S")));
        test.addTestCase("%DB1:1640:S5TIME", new PlcTIME(Duration.parse("PT10S")));
        test.addTestCase("%DB1:1602:LTIME", new PlcLTIME(Duration.parse("PT589H23M12.034002044S")));
        test.addTestCase("%DB1:1610:DATE", new PlcDATE(LocalDate.parse("1998-03-28")));
        test.addTestCase("%DB1:1612:TIME_OF_DAY", new PlcTIME_OF_DAY(LocalTime.parse("15:36:30.123")));
        test.addTestCase("%DB1:1616:LTIME_OF_DAY", new PlcLTIME_OF_DAY(LocalTime.parse("15:36:30")));
        test.addTestCase("%DB1:1624:DATE_AND_TIME", new PlcDATE_AND_TIME(LocalDateTime.parse("1996-05-06T15:36:30")));
        test.addTestCase("%DB1:1632:DATE_AND_LTIME", new PlcDATE_AND_LTIME(LocalDateTime.parse("1978-03-28T15:36:30")));
        test.addTestCase("%DB1:1642:DTL", new PlcDATE_AND_LTIME(LocalDateTime.parse("1978-03-28T15:36:30").withNano(34002044)));
        long start = System.currentTimeMillis();
        test.run();
        long end = System.currentTimeMillis();
        System.out.printf("Finished in %d ms", end - start);
    }

}
