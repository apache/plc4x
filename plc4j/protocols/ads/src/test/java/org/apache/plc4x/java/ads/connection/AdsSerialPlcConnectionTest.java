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
package org.apache.plc4x.java.ads.connection;

import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.model.AdsAddress;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.*;

public class AdsSerialPlcConnectionTest {

    private AdsSerialPlcConnection SUT;

    @Before
    public void setUp() throws Exception {
        SUT = AdsSerialPlcConnection.of("/dev/tty0", AmsNetId.of("0.0.0.0.0.0"), AmsPort.of(13));
    }

    @After
    public void tearDown() {
        SUT = null;
    }

    @Test
    public void initialState() {
        assertEquals(SUT.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
        assertEquals(SUT.getTargetAmsPort().toString(), "13");
    }

    @Test
    public void emptyParseAddress() throws Exception {
        try {
            SUT.parseAddress("");
        } catch (IllegalArgumentException exception) {
            assertTrue("Unexpected exception", exception.getMessage().startsWith("address  doesn't match "));
        }
    }

    @Test
    public void parseAddress() throws Exception {
        try {
            AdsAddress address = (AdsAddress) SUT.parseAddress("1/1");
            assertEquals(address.getIndexGroup(), 1);
            assertEquals(address.getIndexOffset(), 1);
        } catch (IllegalArgumentException exception) {
            fail("valid data block address");
        }
    }
}