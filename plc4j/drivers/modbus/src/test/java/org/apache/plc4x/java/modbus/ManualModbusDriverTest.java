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

import org.apache.plc4x.test.manual.ManualTest;
import org.junit.jupiter.api.Disabled;

import java.util.Arrays;

@Disabled("Manual Test")
public class ManualModbusDriverTest extends ManualTest {

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

    public ManualModbusDriverTest(String connectionString) {
        super(connectionString);
    }

    public static void main(String[] args) throws Exception {
        ManualModbusDriverTest test = new ManualModbusDriverTest("modbus-tcp://192.168.23.30");
        test.addTestCase("holding-register:1:BOOL", true); // 0001
        test.addTestCase("holding-register:2:BYTE", Arrays.asList(false, false, true, false, true, false, true, false)); // 2A
        test.addTestCase("holding-register:3:WORD", Arrays.asList(true, false, true, false, false, true, false, true, true, false, true, true, true, false, false, false)); // A5B8
        test.addTestCase("holding-register:4:DWORD", Arrays.asList(true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, false, true, false, false, false, true, false, false, false, true, false, true, true, true, false, false, false)); // FCDE 88B8
        test.addTestCase("holding-register:6:LWORD", Arrays.asList(true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, false, true, false, false, false, true, false, false, false, true, false, true, true, true, false, false, false, true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, false, true, false, false, false, true, false, false, false, true, false, true, true, true, false, false, false)); // FCDE 88B8 FCDE 88B8
        test.addTestCase("holding-register:10:SINT", -42); // D6
        test.addTestCase("holding-register:11:USINT", 42); // 2A
        test.addTestCase("holding-register:12:INT", -2424); // F688
        test.addTestCase("holding-register:13:UINT", 42424); // A5B8
        test.addTestCase("holding-register:14:DINT", -242442424); // F18C 9F48
        test.addTestCase("holding-register:16:UDINT", 4242442424L);// FCDE 88B8
        test.addTestCase("holding-register:18:LINT", -4242442424242424242L);// C51F D117 B2FB B64E
        test.addTestCase("holding-register:22:ULINT", 4242442424242424242L);// 3AE0 2EE8 4D04 49B2
        test.addTestCase("holding-register:26:REAL", 3.141593F);// 4049 0FDC
        test.addTestCase("holding-register:28:LREAL", 2.71828182846D); // 4005 BF0A 8B14 5FCF
        //test.addTestCase("holding-register:32:TIME", "PT1.234S"); // 04D2
        //test.addTestCase("holding-register::LTIME", "PT24015H23M12.034002044S");
        //test.addTestCase("holding-register::DATE", "1998-03-28");
        //test.addTestCase("holding-register::TIME_OF_DAY", "15:36:30.123");
        //test.addTestCase("holding-register::TOD", "16:17:18.123");
        //test.addTestCase("holding-register::DATE_AND_TIME", "1996-05-06T15:36:30");
        //test.addTestCase("holding-register::DT", "1992-03-29T00:00");
        //test.addTestCase("holding-register::LDATE_AND_TIME", "1978-03-28T15:36:30");
        //test.addTestCase("holding-register::LDT", "1978-03-28T15:36:30");
        //test.addTestCase("holding-register::CHAR", "H");
        //test.addTestCase("holding-register::WCHAR", "w");
        //test.addTestCase("holding-register::STRING(10)", "hurz");
        //test.addTestCase("holding-register::WSTRING(10)", "wolf");
        test.run();
    }

}
