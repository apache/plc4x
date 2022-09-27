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
package org.apache.plc4x.protocol.ads;

import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.test.manual.ManualTest;

import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

public class ManualAdsDriverTest extends ManualTest {

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
	hurz_DATE  := D#1978-03-28;
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

    public ManualAdsDriverTest(String connectionString) {
        super(connectionString, true);
    }

    public static void main(String[] args) throws Exception {
        String spsIp = "192.168.23.20";
        String clientIp = "192.168.23.200";
        String sourceAmsNetId = clientIp + ".1.1";
        int sourceAmsPort = 65534;
        String targetAmsNetId = spsIp + ".1.1";
        int targetAmsPort = 851;
        String connectionString = String.format("ads:tcp://%s?sourceAmsNetId=%s&sourceAmsPort=%d&targetAmsNetId=%s&targetAmsPort=%d", spsIp, sourceAmsNetId, sourceAmsPort, targetAmsNetId, targetAmsPort);
        ManualAdsDriverTest test = new ManualAdsDriverTest(connectionString);
        test.addTestCase("MAIN.hurz_BOOL", new PlcBOOL(true));
        test.addTestCase("MAIN.hurz_BYTE", new PlcBitString(new boolean[]{false, false, true, false, true, false, true, false}));
        test.addTestCase("MAIN.hurz_WORD", new PlcBitString(new boolean[]{true, false, true, false, false, true, false, true, true, false, true, true, true, false, false, false}));
        test.addTestCase("MAIN.hurz_DWORD", new PlcBitString(new boolean[]{true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, false, true, false, false, false, true, false, false, false, true, false, true, true, true, false, false, false}));
        // These are the values, if we decide to go with numeric values instead of bit-strings.
        //test.addTestCase("MAIN.hurz_BYTE", new PlcBYTE(42));
        //test.addTestCase("MAIN.hurz_WORD", new PlcWORD(42424));
        //test.addTestCase("MAIN.hurz_DWORD", new PlcDWORD(4242442424L));
        test.addTestCase("MAIN.hurz_SINT", new PlcSINT(-42));
        test.addTestCase("MAIN.hurz_USINT", new PlcUSINT(42));
        test.addTestCase("MAIN.hurz_INT", new PlcINT(-2424));
        test.addTestCase("MAIN.hurz_UINT", new PlcUINT(42424));
        test.addTestCase("MAIN.hurz_DINT", new PlcDINT(-242442424));
        test.addTestCase("MAIN.hurz_UDINT", new PlcUDINT(4242442424L));
        test.addTestCase("MAIN.hurz_LINT", new PlcLINT(-4242442424242424242L));
        test.addTestCase("MAIN.hurz_ULINT", new PlcULINT(4242442424242424242L));
        test.addTestCase("MAIN.hurz_REAL", new PlcREAL(3.14159265359F));
        test.addTestCase("MAIN.hurz_LREAL", new PlcLREAL(2.71828182846D));
        test.addTestCase("MAIN.hurz_STRING", new PlcSTRING("hurz"));
        test.addTestCase("MAIN.hurz_WSTRING", new PlcWSTRING("wolf"));
        test.addTestCase("MAIN.hurz_TIME", new PlcTIME(Duration.parse("PT1.234S")));
        test.addTestCase("MAIN.hurz_LTIME", new PlcLTIME(Duration.parse("PT24015H23M12.034002044S")));
        test.addTestCase("MAIN.hurz_DATE", new PlcDATE(LocalDate.parse("1978-03-28")));
        test.addTestCase("MAIN.hurz_TIME_OF_DAY", new PlcTIME_OF_DAY(LocalTime.parse("15:36:30.123")));
        test.addTestCase("MAIN.hurz_DATE_AND_TIME", new PlcDATE_AND_TIME(LocalDateTime.parse("1996-05-06T15:36:30")));
        //test.addTestCase("MAIN.hurz_DT", new PlcDT("1972-03-29T00:00"));
        test.run();
    }

}
