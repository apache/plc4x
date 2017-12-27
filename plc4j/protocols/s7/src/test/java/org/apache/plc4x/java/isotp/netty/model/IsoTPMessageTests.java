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

package org.apache.plc4x.java.isotp.netty.model;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.isotp.netty.model.params.Parameter;
import org.apache.plc4x.java.isotp.netty.model.tpdus.ErrorTpdu;
import org.apache.plc4x.java.isotp.netty.model.types.RejectCause;
import org.apache.plc4x.java.isotp.netty.model.types.TpduCode;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class IsoTPMessageTests {

    @Test
    @Tag("fast")
    void isoTPMessage() {
        short destinationReference = 0x1;
        RejectCause rejectCause = RejectCause.REASON_NOT_SPECIFIED;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf errorData = Unpooled.buffer();
        ErrorTpdu tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, errorData);
        ByteBuf userData = Unpooled.buffer();
        IsoTPMessage isoTpMessage = new IsoTPMessage(tpdu, userData);

        errorData.writeByte(0x72);
        userData.writeByte(0x32);

        assertTrue(isoTpMessage.getTpdu().getTpduCode() == TpduCode.TPDU_ERROR, "Unexpected Tpdu");
        // Question: do we need two user data fields?
        assertTrue(tpdu.getUserData().readByte() == (byte) 0x72, "Unexpected user data");
        assertTrue(isoTpMessage.getUserData().readByte() == (byte) 0x32, "Unexpected user data");
    }


}