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

import static org.assertj.core.api.Assertions.assertThat;

import org.apache.plc4x.java.isotp.netty.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.netty.model.types.ParameterCode;
import org.apache.plc4x.test.FastTests;
import org.junit.After;
import org.junit.Test;
import org.junit.experimental.categories.Category;

public class TsapParameterTests {

    private TsapParameter tsapParameter;

    @After
    public void tearDown() {
        tsapParameter = null;
    }

    @Test
    @Category(FastTests.class)
    public void calledPartameter() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0);
        tsapParameter = new CalledTsapParameter(deviceGroup, (byte)1, (byte)4);

        assertThat(tsapParameter.getDeviceGroup()).isEqualTo(DeviceGroup.valueOf((byte)0)).withFailMessage("Device group incorrect");
        assertThat(tsapParameter.getRackNumber()).isEqualTo((byte)1).withFailMessage("Rack number not correct");
        assertThat(tsapParameter.getSlotNumber()).isEqualTo((byte)4).withFailMessage("Slot number not coorect");
        assertThat(tsapParameter.getType()).isEqualTo(ParameterCode.CALLED_TSAP);
    }

    @Test
    @Category(FastTests.class)
    public void callingPartameter() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0);
        tsapParameter = new CallingTsapParameter(deviceGroup, (byte)2, (byte)5);

        assertThat(tsapParameter.getDeviceGroup()).isEqualTo(DeviceGroup.valueOf((byte)0)).withFailMessage("Device group incorrect");
        assertThat(tsapParameter.getRackNumber()).isEqualTo((byte)2).withFailMessage("Rack number not correct");
        assertThat(tsapParameter.getSlotNumber()).isEqualTo((byte)5).withFailMessage("Slot number not coorect");
        assertThat(tsapParameter.getType()).isEqualTo(ParameterCode.CALLING_TSAP);
    }

}