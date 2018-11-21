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

package org.apache.plc4x.java.isotp.protocol.model;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.isotp.protocol.model.params.Parameter;
import org.apache.plc4x.java.isotp.protocol.model.tpdus.ErrorTpdu;
import org.apache.plc4x.java.isotp.protocol.model.types.RejectCause;
import org.apache.plc4x.java.isotp.protocol.model.types.TpduCode;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.Collections;
import java.util.List;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class IsoTPMessageTests {

    @Test
    @Category(FastTests.class)
    public void isoTPMessage() {
        short destinationReference = 0x1;
        RejectCause rejectCause = RejectCause.REASON_NOT_SPECIFIED;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf errorData = Unpooled.buffer();
        ErrorTpdu tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, errorData);
        ByteBuf userData = Unpooled.buffer();
        IsoTPMessage isoTpMessage = new IsoTPMessage(tpdu, userData);

        errorData.writeByte(0x72);
        userData.writeByte(0x32);

        assertThat(isoTpMessage.getTpdu().getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        // Question: do we need two user data fields?
        assertThat(tpdu.getUserData().readByte(), equalTo((byte) 0x72));
        assertThat(isoTpMessage.getUserData().readByte(), equalTo((byte) 0x32));
    }


}