/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.protocol.ads;

import org.apache.plc4x.test.manual.ManualTest;

import java.util.Arrays;

public class ManualAdsDriverTest extends ManualTest {

    public ManualAdsDriverTest(String connectionString) {
        super(connectionString);
    }

    public static void main(String[] args) throws Exception {
        ManualAdsDriverTest test = new ManualAdsDriverTest("ads:tcp://192.168.23.20?sourceAmsNetId=192.168.23.200.1.1&sourceAmsPort=65534&targetAmsNetId=192.168.23.20.1.1&targetAmsPort=851");
        test.addTestCase("main.hurz_BOOL:BOOL", true);
        test.addTestCase("main.hurz_BYTE:BYTE", Arrays.asList(false, false, true, false, true, false, true, false));
        test.addTestCase("main.hurz_WORD:WORD", Arrays.asList(true, false, true, false, false, true, false, true, true, false, true, true, true, false, false, false));
        test.addTestCase("main.hurz_DWORD:DWORD", Arrays.asList(true, true, true, true, true, true, false, false, true, true, false, true, true, true, true, false, true, false, false, false, true, false, false, false, true, false, true, true, true, false, false, false));
        test.addTestCase("main.hurz_SINT:SINT", -42);
        test.addTestCase("main.hurz_USINT:USINT", 42);
        test.addTestCase("main.hurz_INT:INT", -2424);
        test.addTestCase("main.hurz_UINT:UINT", 42424);
        test.addTestCase("main.hurz_DINT:DINT", -242442424);
        test.addTestCase("main.hurz_UDINT:UDINT", 4242442424L);
        test.addTestCase("main.hurz_LINT:LINT", -4242442424242424242L);
        test.addTestCase("main.hurz_ULINT:ULINT", 4242442424242424242L);
        test.addTestCase("main.hurz_REAL:REAL", 3.14159265359F);
        test.addTestCase("main.hurz_LREAL:LREAL", 2.71828182846D);
        test.addTestCase("main.hurz_STRING:STRING", "hurz");
        test.addTestCase("main.hurz_WSTRING:WSTRING", "wolf");
        test.addTestCase("main.hurz_TIME:TIME", "PT1.234S");
        test.addTestCase("main.hurz_LTIME:LTIME", "PT24015H23M12.034002044S");
        test.addTestCase("main.hurz_DATE:DATE", "1978-03-28");
        test.addTestCase("main.hurz_TIME_OF_DAY:TIME_OF_DAY", "15:36:30.123");
        test.addTestCase("main.hurz_TOD:TOD", "16:17:18.123");
        test.addTestCase("main.hurz_DATE_AND_TIME:DATE_AND_TIME", "1996-05-06T15:36:30");
        test.addTestCase("main.hurz_DT:DT", "1972-03-29T00:00");
        test.run();
    }

}
