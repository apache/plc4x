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

package org.apache.plc4x.java.opcua.context;

import org.junit.jupiter.api.*;

import java.util.Arrays;
import java.util.stream.Stream;

import static org.apache.plc4x.java.opcua.context.OpcuaDriverContext.INET_ADDRESS_PATTERN;
import static org.apache.plc4x.java.opcua.context.OpcuaDriverContext.URI_PATTERN;
import static org.assertj.core.api.Assertions.assertThat;

class OpcuaDriverContextTest {

    @Test
    public void testOpcuaAddressPattern() {
        assertThat(":tcp://localhost").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://localhost:3131").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://www.google.de").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://www.google.de:443").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://127.0.0.1").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://127.0.0.1:251").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://254.254.254.254:1337").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://254.254.254.254").matches(INET_ADDRESS_PATTERN);


        assertThat("opcua:tcp://localhost").matches(URI_PATTERN);
        assertThat("opcua:tcp://localhost:3131").matches(URI_PATTERN);
        assertThat("opcua:tcp://www.google.de").matches(URI_PATTERN);
        assertThat("opcua:tcp://www.google.de:443").matches(URI_PATTERN);
        assertThat("opcua:tcp://127.0.0.1").matches(URI_PATTERN);
        assertThat("opcua:tcp://127.0.0.1:251").matches(URI_PATTERN);
        assertThat("opcua:tcp://254.254.254.254:1337").matches(URI_PATTERN);
        assertThat("opcua:tcp://254.254.254.254").matches(URI_PATTERN);

        assertThat("opcua:tcp://127.0.0.1?discovery=false").matches(URI_PATTERN);
        assertThat("opcua:tcp://opcua.demo-this.com:51210/UA/SampleServer?discovery=false").matches(URI_PATTERN);
    }


    private final String[] validTCPOPC = {
        "localhost",
        "127.0.0.1",
        "254.254.254.254"
    };
    private final int[] validPorts = {
        1337,
        42,
        1,
        24152
    };
    private final String[] nDiscoveryParams = {
        "discovery=false"
    };

    @TestFactory
    Stream<DynamicNode> testConnectionStringPattern() throws Exception {
        return Arrays.stream(validTCPOPC)
            .map(address -> DynamicContainer.dynamicContainer("Address: " + address, () -> Arrays.stream(validPorts)
                    .mapToObj(port -> DynamicTest.dynamicTest("Port: " + port, () -> {
                            assertThat("opcua:tcp://" + address + ":555?discovery=true").matches(URI_PATTERN);
                            assertThat("opcua:tcp://" + address + ":555?discovery=True").matches(URI_PATTERN);
                            assertThat("opcua:tcp://" + address + ":555?discovery=TRUE").matches(URI_PATTERN);
                            assertThat("opcua:tcp://" + address + ":555?Discovery=True").matches(URI_PATTERN);
                            //No Port Specified
                            assertThat("opcua:tcp://" + address + "?discovery=True").matches(URI_PATTERN);
                            //No Transport Specified
                            assertThat("opcua://" + address + ":647?discovery=True").matches(URI_PATTERN);
                            //No Params Specified
                            assertThat("opcua:tcp://" + address + ":111").matches(URI_PATTERN);
                            //No Transport and Params Specified
                            assertThat("opcua://" + address + ":754").matches(URI_PATTERN);
                        })
                    )
                    .map(DynamicNode.class::cast)
                    .iterator()
                )
            );
    }
}