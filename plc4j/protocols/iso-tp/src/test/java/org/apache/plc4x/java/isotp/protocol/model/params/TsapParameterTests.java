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

package org.apache.plc4x.java.isotp.protocol.model.params;

import org.apache.plc4x.java.isotp.protocol.model.types.ParameterCode;
import org.apache.plc4x.test.FastTests;
import org.junit.After;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class TsapParameterTests {

    private TsapParameter tsapParameter;

    @After
    public void tearDown() {
        tsapParameter = null;
    }

    @Test
    @Category(FastTests.class)
    public void calledParameter() {
        short calledTsapId = 0x1234;
        tsapParameter = new CalledTsapParameter(calledTsapId);

        assertThat("TSAP Id incorrect", tsapParameter.getTsapId(), equalTo(calledTsapId));
        assertThat(tsapParameter.getType(), equalTo(ParameterCode.CALLED_TSAP));
    }

    @Test
    @Category(FastTests.class)
    public void callingParameter() {
        short callingTsapId = 0x4321;
        tsapParameter = new CallingTsapParameter(callingTsapId);

        assertThat("TSAP Id incorrect", tsapParameter.getTsapId(), equalTo(callingTsapId));
        assertThat(tsapParameter.getType(), equalTo(ParameterCode.CALLING_TSAP));
    }

}