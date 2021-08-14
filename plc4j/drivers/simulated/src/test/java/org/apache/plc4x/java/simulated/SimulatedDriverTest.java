/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.simulated;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.simulated.connection.SimulatedConnection;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

class SimulatedDriverTest implements WithAssertions {

    SimulatedDriver SUT = new SimulatedDriver();

    @Test
    void getProtocolCode() {
        assertThat(SUT.getProtocolCode()).isEqualTo("simulated");
    }

    @Test
    void getProtocolName() {
        assertThat(SUT.getProtocolName()).isEqualTo("Simulated PLC4X Datasource");
    }

    @Test
    void connect() throws Exception {
        assertThat(SUT.getConnection("simulated:foobar")).isInstanceOf(SimulatedConnection.class);
    }

    @Test
    void connect_secure() {
        assertThatThrownBy(() -> SUT.getConnection(null, null)).isInstanceOf(PlcConnectionException.class);
    }

    @Test
    void wrongUrl() {
        assertThatThrownBy(() -> SUT.getConnection("simulated:"))
            .isInstanceOf(PlcConnectionException.class)
            .hasMessage("Invalid URL: no device name given.");
    }

}