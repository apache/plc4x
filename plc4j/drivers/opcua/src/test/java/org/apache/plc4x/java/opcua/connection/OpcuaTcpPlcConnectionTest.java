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
package org.apache.plc4x.java.opcua.connection;

import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.URI_PATTERN;
import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;

public class OpcuaTcpPlcConnectionTest {

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

    @Test
    public void testConectionStringPattern() {

        for (String address : validTCPOPC) {
            assertThat("opcua:tcp://127.0.0.1:555?discovery=true").matches(URI_PATTERN);
            assertThat("opcua:tcp://127.0.0.1:555?discovery=True").matches(URI_PATTERN);
            assertThat("opcua:tcp://127.0.0.1:555?discovery=TRUE").matches(URI_PATTERN);
            assertThat("opcua:tcp://127.0.0.1:555?Discovery=True").matches(URI_PATTERN);
            //No Port Specified
            assertThat("opcua:tcp://127.0.0.1?discovery=True").matches(URI_PATTERN);
            //No Transport Specified
            assertThat("opcua://127.0.0.1:647?discovery=True").matches(URI_PATTERN);
            //No Params Specified
            assertThat("opcua:tcp://127.0.0.1:111").matches(URI_PATTERN);
            //No Transport and Params Specified
            assertThat("opcua://127.0.0.1:754").matches(URI_PATTERN);
        }
    }
}
