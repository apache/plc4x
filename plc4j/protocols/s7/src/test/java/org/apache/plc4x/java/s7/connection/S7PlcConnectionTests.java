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

package org.apache.plc4x.java.s7.connection;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.fail;

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

public class S7PlcConnectionTests {

    private S7PlcConnection  s7PlcConnection;

    @Before
    public void setUp() {
        s7PlcConnection = new S7PlcConnection("localhost", 1, 2, "");
    }

    @After
    public void tearDown() {
        s7PlcConnection = null;
    }

    @Test
    public void initialState() {
        assertThat(s7PlcConnection.getHostName()).isEqualToIgnoringCase("localhost")
            .withFailMessage("Hostname is incorrect");
        assertThat(s7PlcConnection.getRack()).isEqualTo(1)
            .withFailMessage("Rack is incorrect");
        assertThat(s7PlcConnection.getSlot()).isEqualTo(2)
            .withFailMessage("Slot is incorrect");
        assertThat(s7PlcConnection.getParamPduSize()).isEqualTo(TpduSize.SIZE_1024)
            .withFailMessage("Pdu size is incorrect");
        assertThat(s7PlcConnection.getParamMaxAmqCaller()).isEqualTo(8)
            .withFailMessage("Max AMQ Caller size is incorrect");
        assertThat(s7PlcConnection.getParamMaxAmqCallee()).isEqualTo(8)
            .withFailMessage("Max AMQ Callee size is incorrect");
    }

    @Test
    public void emptyParseAddress() {
        try {
            s7PlcConnection.parseAddress("");
        }
        catch (PlcException exception) {
            assertThat(exception.getMessage()).startsWith("Address string doesn't match");
        }
    }

    @Test
    public void parseDatablockAddress() {
        try {
            S7DataBlockAddress address = (S7DataBlockAddress)
                s7PlcConnection.parseAddress("DATA_BLOCKS/20/100");

            assertThat(address.getDataBlockNumber()).isEqualTo((short) 20)
                .withFailMessage("unexpected data block");
            assertThat(address.getByteOffset()).isEqualTo((short) 100)
                .withFailMessage("unexpected byte offset");
        }
        catch (PlcException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseAddressAddress() {
        try {
            S7Address address = (S7Address) s7PlcConnection.parseAddress("TIMERS/10");

            assertThat(address.getMemoryArea()).isEqualTo(MemoryArea.TIMERS)
                .withFailMessage("unexpected memory area");
            assertThat(address.getByteOffset()).isEqualTo((short) 10)
                .withFailMessage("unexpected byte offset");
        }
        catch (PlcException exception) {
            fail("valid timer block address");
        }
    }

    @Test
    public void parseAddressBitAddress() {
        try {
            S7BitAddress address = (S7BitAddress) s7PlcConnection.parseAddress("TIMERS/10/4");

            assertThat(address.getMemoryArea()).isEqualTo(MemoryArea.TIMERS)
                .withFailMessage("unexpected memory area");
            assertThat(address.getByteOffset()).isEqualTo((short) 10)
                .withFailMessage("unexpected byte offset");
            assertThat(address.getBitOffset()).isEqualTo((byte) 4)
                .withFailMessage("unexpected but offset");
        }
        catch (PlcException exception) {
            fail("valid timer block bit address");
        }
    }

}