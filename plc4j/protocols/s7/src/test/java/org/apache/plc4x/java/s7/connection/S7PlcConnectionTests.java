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

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.net.InetAddress;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.StringStartsWith.startsWith;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

public class S7PlcConnectionTests {

    private S7PlcConnection  s7PlcConnection;

    @Before
    public void setUp() throws Exception {
        InetAddress address = InetAddress.getByName("localhost");
        s7PlcConnection = new S7PlcConnection(address, 1, 2, "");
    }

    @After
    public void tearDown() {
        s7PlcConnection = null;
    }

    @Test
    public void initialState() {
        assertThat("Rack is incorrect", s7PlcConnection.getRack(), equalTo(1) );
        assertThat("Slot is incorrect", s7PlcConnection.getSlot(), equalTo(2) );
        assertThat("Pdu size is incorrect", s7PlcConnection.getParamPduSize(), equalTo(TpduSize.SIZE_1024));
        assertThat("Max AMQ Caller size is incorrect", s7PlcConnection.getParamMaxAmqCaller(), equalTo(8) );
        assertThat("Max AMQ Callee size is incorrect", s7PlcConnection.getParamMaxAmqCallee(), equalTo(8) );
    }

    @Test
    public void emptyParseAddress() {
        try {
            s7PlcConnection.parseAddress("");
        }
        catch (PlcException exception) {
            assertThat(exception.getMessage(), startsWith("Address string doesn't match") );
        }
    }

    @Test
    public void parseDatablockAddress() {
        try {
            S7DataBlockAddress address = (S7DataBlockAddress)
                s7PlcConnection.parseAddress("DATA_BLOCKS/20/100");

            assertThat("unexpected data block", address.getDataBlockNumber(), equalTo((short) 20) );
            assertThat("unexpected byte offset", address.getByteOffset(), equalTo((short) 100) );
        }
        catch (PlcException exception) {
            fail("valid data block address");
        }
    }

    @Test
    public void parseAddressAddress() {
        try {
            S7Address address = (S7Address) s7PlcConnection.parseAddress("TIMERS/10");

            assertThat("unexpected memory area", address.getMemoryArea(), equalTo(MemoryArea.TIMERS) );
            assertThat("unexpected byte offset", address.getByteOffset(), equalTo((short) 10) );
        }
        catch (PlcException exception) {
            fail("valid timer block address");
        }
    }

    @Test
    public void parseAddressBitAddress() {
        try {
            S7BitAddress address = (S7BitAddress) s7PlcConnection.parseAddress("TIMERS/10/4");

            assertThat("unexpected memory area", address.getMemoryArea(), equalTo(MemoryArea.TIMERS) );
            assertThat("unexpected byte offset", address.getByteOffset(), equalTo((short) 10) );
            assertThat("unexpected but offset", address.getBitOffset(), equalTo((byte) 4) );
        }
        catch (PlcException exception) {
            fail("valid timer block bit address");
        }
    }

}