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
package org.apache.plc4x.java.opcua.protocol;

import org.junit.jupiter.api.Test;

import static org.apache.plc4x.java.opcua.tag.OpcuaTag.ADDRESS_PATTERN;
import static org.assertj.core.api.Assertions.assertThat;

public class OpcuaTagTest {

    @Test
    public void testOpcuaAddressPattern() {

        //standard integer based param
        assertThat("ns=2;i=10846").matches(ADDRESS_PATTERN);
        //string based address values
        assertThat("ns=2;s=test.variable.name.inspect").matches(ADDRESS_PATTERN);
        assertThat("ns=2;s=::AsGlobalPV:ProductionOrder").matches(ADDRESS_PATTERN);
        assertThat("ns=2;s=::AsGlobalPV:ProductionOrder;BOOL").matches(ADDRESS_PATTERN);
        assertThat("ns=2;s=key param with some spaces").matches(ADDRESS_PATTERN);
        assertThat("ns=2;s=\"aweired\".\"siemens\".\"param\".\"submodule\".\"param").matches(ADDRESS_PATTERN);
        assertThat("ns=2;s=Weee314Waannaaa\\somenice=ext=a234a*#+1455!§$%&/()tttraaaaSymbols-:.,,").matches(ADDRESS_PATTERN);
        // GUID address tests
        assertThat("ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a").matches(ADDRESS_PATTERN);
        // binary encoded addresses
        assertThat("ns=2;b=asvaewavarahreb==").matches(ADDRESS_PATTERN);

    }

    @Test
    public void testOpcuaAddressDataTypePattern() {

        //standard integer based param
        assertThat("ns=2;i=10846;BOOL").matches(ADDRESS_PATTERN);
        //string based address values
        assertThat("ns=2;s=test.variable.name.inspect;DINT").matches(ADDRESS_PATTERN);
        assertThat("ns=2;s=key param with some spaces;ULINT").matches(ADDRESS_PATTERN);
        assertThat("ns=2;s=\"aweired\".\"siemens\".\"param\".\"submodule\".\"param;LREAL").matches(ADDRESS_PATTERN);
        //REGEX Valid, additional checks need to be done later
        assertThat("ns=2;s=Weee314Waannaaa\\somenice=ext=a234a*#+1455!§$%&/()tttraaaaSymbols-.,,;JIBBERISH").matches(ADDRESS_PATTERN);
        // GUID address tests
        assertThat("ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a;*&#%^*$(*)").doesNotMatch(ADDRESS_PATTERN);
        // binary encoded addresses
        assertThat("ns=2;b=asvae;wavarahreb==").doesNotMatch(ADDRESS_PATTERN);

    }
}
