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

package org.apache.plc4x.java.simulated;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.simulated.connection.SimulatedPlcConnection;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

class SimulatedPlcDriverTest implements WithAssertions {

    SimulatedPlcDriver SUT = new SimulatedPlcDriver();

    @Test
    void getProtocolCode() {
        assertThat(SUT.getProtocolCode()).isEqualTo("test");
    }

    @Test
    void getProtocolName() {
        assertThat(SUT.getProtocolName()).isEqualTo("PLC4X Test Protocol");
    }

    @Test
    void connect() throws Exception {
        assertThat(SUT.connect("test:foobar")).isInstanceOf(SimulatedPlcConnection.class);
    }

    @Test
    void connect_secure() {
        assertThatThrownBy(() -> SUT.connect(null, null)).isInstanceOf(PlcConnectionException.class);
    }

    @Test
    void wrongUrl() {
        assertThatThrownBy(() -> SUT.connect("test:"))
            .isInstanceOf(PlcConnectionException.class)
            .hasMessage("Invalid URL: no device name given.");
    }
}