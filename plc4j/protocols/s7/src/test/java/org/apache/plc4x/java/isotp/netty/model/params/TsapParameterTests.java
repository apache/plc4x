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

package org.apache.plc4x.java.isotp.netty.model.params;

import org.apache.plc4x.java.isotp.netty.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.netty.model.types.ParameterCode;
import org.junit.jupiter.api.*;

import static org.junit.jupiter.api.Assertions.*;

class TsapParameterTests {
    private TsapParameter tsapParameter;

    @BeforeEach
    void setUp() {
    }

    @AfterEach
    void tearDown() {
        tsapParameter = null;
    }

    @Test
    @Tag("fast")
    void calledPartameter() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0);
        tsapParameter = new CalledTsapParameter(deviceGroup, (byte)1, (byte)4);

        assertTrue(tsapParameter.getDeviceGroup() == DeviceGroup.valueOf((byte)0), "Device group incorrect");
        assertTrue(tsapParameter.getRackNumber() == (byte)1, "Rack number not correct");
        assertTrue(tsapParameter.getSlotNumber() == (byte)4, "Slot number not coorect");
        assertTrue(tsapParameter.getType() == ParameterCode.CALLED_TSAP);
    }

    @Test
    @Tag("fast")
    void callingPartameter() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0);
        tsapParameter = new CallingTsapParameter(deviceGroup, (byte)2, (byte)5);

        assertTrue(tsapParameter.getDeviceGroup() == DeviceGroup.valueOf((byte)0), "Device group incorrect");
        assertTrue(tsapParameter.getRackNumber() == (byte)2, "Rack number not correct");
        assertTrue(tsapParameter.getSlotNumber() == (byte)5, "Slot number not coorect");
        assertTrue(tsapParameter.getType() == ParameterCode.CALLING_TSAP);
    }

}