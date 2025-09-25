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

import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcBYTE;
import org.apache.plc4x.test.manual.ManualTest;

public class ManualS71200DriverTest3 extends ManualTest {

    /*
     * Test program code on the PLC with the test-data.
     * Reading the states of the 5 input booleans. Used to test how multiple bits located in one byte are handled.
     */

    public ManualS71200DriverTest3(String connectionString) {
        super(connectionString, true, true, false, true, 100);
    }

    public static void main(String[] args) throws Exception {
        ManualS71200DriverTest3 test = new ManualS71200DriverTest3("s7-light://192.168.23.30");
        test.addTestCase("%I0:BYTE", new PlcBYTE(6));
        test.addTestCase("%I0.0:BOOL", new PlcBOOL(false));
        test.addTestCase("%I0.1:BOOL", new PlcBOOL(true));
        test.addTestCase("%I0.2:BOOL", new PlcBOOL(true));
        test.addTestCase("%I0.3:BOOL", new PlcBOOL(false));
        test.addTestCase("%I0.4:BOOL", new PlcBOOL(false));

        long start = System.currentTimeMillis();
        test.run();
        long end = System.currentTimeMillis();
        System.out.printf("Finished in %d ms", end - start);
    }

}
