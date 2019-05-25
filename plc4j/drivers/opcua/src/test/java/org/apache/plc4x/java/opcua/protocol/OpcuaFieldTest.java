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
package org.apache.plc4x.java.opcua.protocol;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.apache.plc4x.java.opcua.UtilsTest.assertMatching;
import static org.apache.plc4x.java.opcua.protocol.OpcuaField.ADDRESS_PATTERN;
/**
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
public class OpcuaFieldTest {

    @BeforeEach
    public void before() {
    }

    @AfterEach
    public void after() {

    }

    @Test
    public void testOpcuaAddressPattern() {

        assertMatching(ADDRESS_PATTERN, "ns=2;i=10846");
        assertMatching(ADDRESS_PATTERN, "ns=2;s=test.variable.name.inspect");
        assertMatching(ADDRESS_PATTERN, "ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a");
        assertMatching(ADDRESS_PATTERN, "ns=2;b=asvaewavarahreb==");

    }
}
