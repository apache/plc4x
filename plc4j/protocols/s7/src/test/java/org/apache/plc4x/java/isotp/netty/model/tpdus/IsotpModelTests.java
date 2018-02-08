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

package org.apache.plc4x.java.isotp.netty.model.tpdus;

import static org.assertj.core.api.Assertions.assertThat;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.isotp.netty.model.params.CallingTsapParameter;
import org.apache.plc4x.java.isotp.netty.model.params.ChecksumParameter;
import org.apache.plc4x.java.isotp.netty.model.params.Parameter;
import org.apache.plc4x.java.isotp.netty.model.params.TpduSizeParameter;
import org.apache.plc4x.java.isotp.netty.model.types.*;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class IsotpModelTests {

    @Test
    @Category(FastTests.class)
    public void errorTpdu() {
        short destinationReference = 0x1;
        RejectCause rejectCause = RejectCause.REASON_NOT_SPECIFIED;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x7F);

        ErrorTpdu tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, userData);

        assertThat(tpdu.getTpduCode()).isEqualTo(TpduCode.TPDU_ERROR);
        assertThat(tpdu.getDestinationReference()).isEqualTo((short) 0x1).withFailMessage("Unexpected destination reference");
        assertThat(tpdu.getRejectCause()).isEqualTo(RejectCause.REASON_NOT_SPECIFIED);
        assertThat(tpdu.getParameters()).isEmpty();
        assertThat(tpdu.getUserData().readByte()).isEqualTo((byte) 0x7F).withFailMessage("Unexpected user data");
    }

    @Test
    @Category(FastTests.class)
    public void errorTpduParameter() {
        short destinationReference = 0x1;
        RejectCause rejectCause = RejectCause.REASON_NOT_SPECIFIED;
        ArrayList<Parameter> parameters = new ArrayList<>();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x7F);

        ErrorTpdu tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, userData);

        parameters.add(new TpduSizeParameter(TpduSize.SIZE_1024));
        parameters.add(new ChecksumParameter((byte) 0xFF));

        assertThat(tpdu.getParameters()).hasSize(2).withFailMessage("Unexpected number of parameters");
        assertThat(tpdu.getParameters()).containsAll(parameters).withFailMessage("Unexpected parameter");
        assertThat(tpdu.getParameter(ChecksumParameter.class).isPresent()).isTrue().withFailMessage("Checksum parameter should exist");
        assertThat(!tpdu.getParameter(CallingTsapParameter.class).isPresent()).isTrue().withFailMessage("CallingTsapParameter parameter should not exist");
    }

    @Test
    @Category(FastTests.class)
    public void dataTpdu() {
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x66);

        DataTpdu tpdu = new DataTpdu(true, (byte) 0x7F, parameters, userData);

        assertThat(tpdu.getTpduCode()).isEqualTo(TpduCode.DATA);
        assertThat(tpdu.isEot()).isTrue().withFailMessage("Unexpected eot reference");
        assertThat(tpdu.getTpduRef()).isEqualTo((byte) 0x7F);
        assertThat(tpdu.getParameters().isEmpty()).isTrue().withFailMessage("Unexpected parameters");
        assertThat(tpdu.getUserData().readByte()).isEqualTo((byte) 0x66).withFailMessage("Unexpected user data");
    }

    @Test
    @Category(FastTests.class)
    public void connectionRequestTpdu() {
        short destinationReference = 0x1;
        short sourceReference = 0x2;
        ProtocolClass protocolClass = ProtocolClass.CLASS_0;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x33);

        ConnectionRequestTpdu tpdu = new ConnectionRequestTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);

        assertThat(tpdu.getTpduCode() == TpduCode.CONNECTION_REQUEST);
        assertThat(tpdu.getDestinationReference()).isEqualTo((short) 0x1).withFailMessage("Unexpected destination reference");
        assertThat(tpdu.getSourceReference()).isEqualTo((short) 0x2).withFailMessage("Unexpected source reference");
        assertThat(tpdu.getProtocolClass()).isEqualTo(ProtocolClass.CLASS_0);
        assertThat(tpdu.getParameters().isEmpty()).isTrue().withFailMessage("Unexpected parameters");
        assertThat(tpdu.getUserData().readByte()).isEqualTo((byte) 0x33).withFailMessage("Unexpected user data");
    }

    @Test
    @Category(FastTests.class)
    public void connectionConfirmTpdu() {
        short destinationReference = 0x3;
        short sourceReference = 0x4;
        ProtocolClass protocolClass = ProtocolClass.CLASS_1;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x44);

        ConnectionConfirmTpdu tpdu = new ConnectionConfirmTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);

        assertThat(tpdu.getTpduCode() == TpduCode.CONNECTION_CONFIRM);
        assertThat(tpdu.getDestinationReference()).isEqualTo((short) 0x3).withFailMessage("Unexpected destination reference");
        assertThat(tpdu.getSourceReference()).isEqualTo((short) 0x4).withFailMessage("Unexpected source reference");
        assertThat(tpdu.getProtocolClass()).isEqualTo(ProtocolClass.CLASS_1);
        assertThat(tpdu.getParameters().isEmpty()).isTrue().withFailMessage("Unexpected parameters");
        assertThat(tpdu.getUserData().readByte()).isEqualTo((byte) 0x44).withFailMessage("Unexpected user data");
    }

    @Test
    @Category(FastTests.class)
    public void disconnectionRequestTpdu() {
        short destinationReference = 0x1;
        short sourceReference = 0x2;
        DisconnectReason disconnectReason = DisconnectReason.ADDRESS_UNKNOWN;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x22);

        DisconnectRequestTpdu tpdu = new DisconnectRequestTpdu(destinationReference, sourceReference, disconnectReason, parameters, userData);

        assertThat(tpdu.getTpduCode() == TpduCode.DISCONNECT_REQUEST);
        assertThat(tpdu.getDestinationReference()).isEqualTo((short) 0x1).withFailMessage("Unexpected destination reference");
        assertThat(tpdu.getSourceReference()).isEqualTo((short) 0x2).withFailMessage("Unexpected source reference");
        assertThat(tpdu.getDisconnectReason()).isEqualTo(DisconnectReason.ADDRESS_UNKNOWN);
        assertThat(tpdu.getParameters().isEmpty()).isTrue().withFailMessage("Unexpected parameters");
        assertThat(tpdu.getUserData().readByte()).isEqualTo((byte) 0x22).withFailMessage("Unexpected user data");
    }

    @Test
    @Category(FastTests.class)
    public void disconnectionConfirmTpdu() {
        short destinationReference = 0x3;
        short sourceReference = 0x4;
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x11);

        DisconnectConfirmTpdu tpdu = new DisconnectConfirmTpdu(destinationReference, sourceReference, parameters, userData);

        assertThat(tpdu.getTpduCode() == TpduCode.DISCONNECT_CONFIRM);
        assertThat(tpdu.getDestinationReference()).isEqualTo((short) 0x3).withFailMessage("Unexpected destination reference");
        assertThat(tpdu.getSourceReference()).isEqualTo((short) 0x4).withFailMessage("Unexpected source reference");
        assertThat(tpdu.getParameters().isEmpty()).isTrue().withFailMessage("Unexpected parameters");
        assertThat(tpdu.getUserData().readByte()).isEqualTo((byte) 0x11).withFailMessage("Unexpected user data");
    }

}