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
package org.apache.plc4x.java.opcua;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.INET_ADDRESS_PATTERN;
import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.OPCUA_URI_PATTERN;
import static org.apache.plc4x.java.opcua.UtilsTest.assertMatching;
/**
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
public class OpcuaPlcDriverTest {
    @BeforeEach
    public void before() {
    }

    @AfterEach
    public void after() {

    }

    @Test
    public void testOpcuaAddressPattern() {

        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost:3131");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de:443");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://127.0.0.1");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://127.0.0.1:251");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://254.254.254.254:1337");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://254.254.254.254");


        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://localhost");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://localhost:3131");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://www.google.de");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://www.google.de:443");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://127.0.0.1");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://127.0.0.1:251");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://254.254.254.254:1337");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://254.254.254.254");


    }

}
