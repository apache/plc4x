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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.test.manual.ManualTest;
import org.junit.jupiter.api.Disabled;

@Disabled("Manual Test")
public class ManualModbusDriverTestModbusPal extends ManualTest {

    /*
     * Test program code on the PLC with the test-data.
     *
     * Located in "main"
     *

    hurz_BOOL  := TRUE;
	hurz_BYTE  := 42;
	hurz_WORD  := 42424;
	hurz_DWORD := 4242442424;
	hurz_LWORD := 4242442424242424242;
	hurz_SINT  := -42;
	hurz_USINT := 42;
	hurz_INT   := -2424;
	hurz_UINT  := 42424;
	hurz_DINT  := -242442424;
	hurz_UDINT := 4242442424;
	hurz_LINT  := -4242442424242424242;
	hurz_ULINT := 4242442424242424242;
	hurz_REAL  := 3.14159265359;
	hurz_LREAL := 2.71828182846;
	hurz_TIME  := T#1S234MS;
	hurz_LTIME := LTIME#1000D15H23M12S34MS2US44NS;
	hurz_DATE  := D#1998-03-28;
	//hurz_LDATE:LDATE;
	hurz_TIME_OF_DAY 	:= TIME_OF_DAY#15:36:30.123;
	hurz_TOD         	:= TOD#16:17:18.123;
	//hurz_LTIME_OF_DAY:LTIME_OF_DAY;
	//hurz_LTOD:LTOD;
	hurz_DATE_AND_TIME 	:= DATE_AND_TIME#1996-05-06-15:36:30;
	hurz_DT				:= DT#1972-03-29-00:00:00;
	//hurz_LDATE_AND_TIME:LDATE_AND_TIME;
	//hurz_LDT:LDT;
	hurz_STRING			:= 'hurz';
	hurz_WSTRING		:= "wolf";

     *
     */

    public ManualModbusDriverTestModbusPal(String connectionString) {
        super(connectionString);
    }

    public static void main(String[] args) throws Exception {
        // ! See modbus-pal-project.xmpp for a config made to be used by this test
        // Tested with ModbusPal
        ManualModbusDriverTestModbusPal test = new ManualModbusDriverTestModbusPal("modbus-tcp://127.0.0.1");
        test.addTestCase("holding-register:1000:BOOL", new PlcBOOL(true)); // 0001 # 1
        test.addTestCase("holding-register:1001:BYTE", new PlcBYTE(42)); // 2A # 42
        test.addTestCase("holding-register:1002:WORD", new PlcWORD(42424)); // A5B8 # 42424
        test.addTestCase("holding-register:1003:DWORD", new PlcDWORD(4242442424L)); // FCDE 88B8 # 64734 35000
        test.addTestCase("holding-register:1005:LWORD", new PlcLWORD(4242442424242424242L)); // 3AE0 2EE8 4D04 49B2 # 15072 12008 19716 18866
        test.addTestCase("holding-register:1009:SINT", new PlcSINT(-42)); // D6 # 214
        test.addTestCase("holding-register:1010:USINT", new PlcUSINT(42)); // 2A # 42
        test.addTestCase("holding-register:1011:INT", new PlcINT(-2424)); // F688 # 63112
        test.addTestCase("holding-register:1012:UINT", new PlcUINT(42424)); // A5B8 # 42424
        test.addTestCase("holding-register:1013:DINT", new PlcDINT(-242442424)); // F18C 9F48 # 61836 40776
        test.addTestCase("holding-register:1015:UDINT", new PlcUDINT(4242442424L));// FCDE 88B8 # 64734 35000
        test.addTestCase("holding-register:1017:LINT", new PlcLINT(-4242442424242424242L));// C51F D117 B2FB B64E # 50463 53527 45819 46670
        test.addTestCase("holding-register:1021:ULINT", new PlcULINT(4242442424242424242L));// 3AE0 2EE8 4D04 49B2 # 15072 12008 19716 18866
        test.addTestCase("holding-register:1025:REAL", new PlcREAL(3.141593F));// 4049 0FDC # 16457 4060
        test.addTestCase("holding-register:1027:LREAL", new PlcLREAL(2.71828182846D)); // 4005 BF0A 8B14 5FCF # 16389 48906 35604 24527
// TODO: These datatypes are not yet fully implemented
//        test.addTestCase("holding-register:1031:STRING", new PlcSTRING("hurz")); // 6875 727A # 26741 29306
//        test.addTestCase("holding-register:1033:WSTRING", new PlcWSTRING("wolf")); // 0068 0075 0072 007A # 104 117 114 122
//        test.addTestCase("holding-register:1037:TIME", new PlcTIME(Duration.parse("PT1.234S"))); // 04D2 # 1234
//        test.addTestCase("holding-register::LTIME", new PlcLTIME(Duration.parse("PT24015H23M12.034002044S")));
//        test.addTestCase("holding-register::DATE", new PlcDATE(LocalDate.parse("1978-03-28")));
//        test.addTestCase("holding-register::TIME_OF_DAY", new PlcTIME_OF_DAY(LocalTime.parse("15:36:30.123")));
//        test.addTestCase("holding-register::DATE_AND_TIME", new PlcDATE_AND_TIME(LocalDateTime.parse("1996-05-06T15:36:30")));
        test.run();
    }

}
