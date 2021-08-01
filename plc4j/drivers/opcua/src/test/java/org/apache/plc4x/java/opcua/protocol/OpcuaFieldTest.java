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
import static org.apache.plc4x.java.opcua.UtilsTest.assertNoMatching;
import static org.apache.plc4x.java.opcua.field.OpcuaField.ADDRESS_PATTERN;

/**
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

        //standard integer based param
        assertMatching(ADDRESS_PATTERN, "ns=2;i=10846");
        //string based address values
        assertMatching(ADDRESS_PATTERN, "ns=2;s=test.variable.name.inspect");
        assertMatching(ADDRESS_PATTERN, "ns=2;s=key param with some spaces");
        assertMatching(ADDRESS_PATTERN, "ns=2;s=\"aweired\".\"siemens\".\"param\".\"submodule\".\"param");
        assertMatching(ADDRESS_PATTERN, "ns=2;s=Weee314Waannaaa\\somenice=ext=a234a*#+1455!§$%&/()tttraaaaSymbols-.,,");
        // GUID address tests
        assertMatching(ADDRESS_PATTERN, "ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a");
        // binary encoded addresses
        assertMatching(ADDRESS_PATTERN, "ns=2;b=asvaewavarahreb==");

    }

    @Test
    public void testOpcuaAddressDataTypePattern() {

        //standard integer based param
        assertMatching(ADDRESS_PATTERN, "ns=2;i=10846:BOOL");
        //string based address values
        assertMatching(ADDRESS_PATTERN, "ns=2;s=test.variable.name.inspect:DINT");
        assertMatching(ADDRESS_PATTERN, "ns=2;s=key param with some spaces:ULINT");
        assertMatching(ADDRESS_PATTERN, "ns=2;s=\"aweired\".\"siemens\".\"param\".\"submodule\".\"param:LREAL");
        //REGEX Valid, additional checks need to be done later
        assertMatching(ADDRESS_PATTERN, "ns=2;s=Weee314Waannaaa\\somenice=ext=a234a*#+1455!§$%&/()tttraaaaSymbols-.,,:JIBBERISH");
        // GUID address tests
        assertNoMatching(ADDRESS_PATTERN, "ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a:*&#%^*$(*)");
        // binary encoded addresses
        assertNoMatching(ADDRESS_PATTERN, "ns=2;b=asvae:wavarahreb==");

    }
}
