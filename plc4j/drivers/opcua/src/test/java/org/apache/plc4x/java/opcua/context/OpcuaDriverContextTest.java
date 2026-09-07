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
import java.lang.reflect.Field;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;

import static org.assertj.core.api.Assertions.assertThat;

class OpcuaDriverContextTest {

    @Test
    void initializeBuildsEndpointFromUrlComponents() {
        OpcuaDriverContext ctx = new OpcuaDriverContext();
        OpcuaConfiguration cfg = new OpcuaConfiguration();
        try {
            java.lang.reflect.Field f = cfg.getClass().getDeclaredField("securityPolicy");
            f.setAccessible(true);
            f.set(cfg, SecurityPolicy.NONE);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        ctx.initialize("tcp", "192.168.0.1", "4840", "/UA/SampleServer", cfg);
        assertThat(ctx.getHost()).isEqualTo("192.168.0.1");
        assertThat(ctx.getPort()).isEqualTo("4840");
        assertThat(ctx.getTransportEndpoint()).isEqualTo("/UA/SampleServer");
        assertThat(ctx.getEndpoint()).isEqualTo("opc.tcp://192.168.0.1:4840/UA/SampleServer");
    }

    @Test
    void initializeAcceptsNullPortAndNullTransportEndpoint() {
        OpcuaDriverContext ctx = new OpcuaDriverContext();
        OpcuaConfiguration cfg = new OpcuaConfiguration();
        try {
            java.lang.reflect.Field f = cfg.getClass().getDeclaredField("securityPolicy");
            f.setAccessible(true);
            f.set(cfg, SecurityPolicy.NONE);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        ctx.initialize("tcp", "host", null, null, cfg);
        assertThat(ctx.getEndpoint()).isEqualTo("opc.tcp://host");
        assertThat(ctx.getTransportEndpoint()).isEmpty();
    }

    @Test
    void setConfigurationParsesUriIntoHostPortAndEndpoint() throws Exception {
        OpcuaConfiguration cfg = new OpcuaConfiguration();
        setField(cfg, "protocolCode", "opcua");
        setField(cfg, "transportCode", "tcp");
        setField(cfg, "transportConfig", "localhost:4840/UA/SampleServer");
        // SecurityPolicy.NONE skips keystore loading, so no files are needed.
        setField(cfg, "securityPolicy", SecurityPolicy.NONE);

        OpcuaDriverContext ctx = new OpcuaDriverContext();
        ctx.setConfiguration(cfg);

        assertThat(ctx.getHost()).isEqualTo("localhost");
        assertThat(ctx.getPort()).isEqualTo("4840");
        assertThat(ctx.getTransportEndpoint()).isEqualTo("/UA/SampleServer");
        assertThat(ctx.getEndpoint()).isEqualTo("opc.tcp://localhost:4840/UA/SampleServer");
        // No keystore was loaded → no certificate / thumbprint / verifier.
        assertThat(ctx.getCertificateKeyPair()).isNull();
        assertThat(ctx.getServerCertificate()).isNull();
        assertThat(ctx.getApplicationUri()).isEmpty();
    }

    private static void setField(Object target, String name, Object value) throws Exception {
        Field f = target.getClass().getDeclaredField(name);
        f.setAccessible(true);
        f.set(target, value);
    }


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