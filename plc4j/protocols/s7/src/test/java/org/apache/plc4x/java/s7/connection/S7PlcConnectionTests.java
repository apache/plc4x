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

import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.apache.plc4x.java.s7.types.S7ControllerType;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.apache.plc4x.java.base.util.Junit5Backport.assertThrows;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;

public class S7PlcConnectionTests {

    private S7PlcTestConnection SUT;

    @Before
    public void setUp() {
        SUT = new S7PlcTestConnection(1, 2,
            "pdu-size=1&max-amq-caller=2&max-amq-callee=3&unknown=parameter&unknown-flag", S7ControllerType.S7_1200);
    }

    @After
    public void tearDown() {
        SUT = null;
    }

    @Test
    public void initialState() {
        assertThat("Rack is incorrect", SUT.getRack(), equalTo(1) );
        assertThat("Slot is incorrect", SUT.getSlot(), equalTo(2) );
        assertThat("Pdu size is incorrect", SUT.getParamPduSize(), equalTo(TpduSize.SIZE_128));
        assertThat("Max AMQ Caller size is incorrect", SUT.getParamMaxAmqCaller(), equalTo(2) );
        assertThat("Max AMQ Callee size is incorrect", SUT.getParamMaxAmqCallee(), equalTo(3) );
    }

    @Test
    public void capabilities() {
        assertThat(SUT.canRead(), equalTo(true));
        assertThat(SUT.readRequestBuilder(), notNullValue());

        assertThat(SUT.canWrite(), equalTo(true));
        assertThat(SUT.writeRequestBuilder(), notNullValue());

        assertThat(SUT.canSubscribe(), equalTo(false));
        assertThrows(PlcUnsupportedOperationException.class, () -> SUT.subscriptionRequestBuilder());
        assertThrows(PlcUnsupportedOperationException.class, () -> SUT.unsubscriptionRequestBuilder());
    }

}