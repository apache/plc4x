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

import static org.hamcrest.collection.IsCollectionWithSize.hasSize;
import static org.hamcrest.collection.IsEmptyCollection.empty;
import static org.hamcrest.collection.IsIterableContainingInOrder.contains;
import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

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

        assertThat(tpdu.getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        assertThat("Unexpected destination reference", tpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(tpdu.getRejectCause(), equalTo(RejectCause.REASON_NOT_SPECIFIED));
        assertThat(tpdu.getParameters(), empty() );
        assertThat("Unexpected user data", tpdu.getUserData().readByte(), equalTo((byte) 0x7F));
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

        assertThat("Unexpected number of parameters", tpdu.getParameters(), hasSize(2));
        assertThat("Unexpected parameter", tpdu.getParameters(), contains(parameters.toArray()));
        assertThat("Checksum parameter should exist", tpdu.getParameter(ChecksumParameter.class).isPresent(), is(true));
        assertThat("CallingTsapParameter parameter should not exist", !tpdu.getParameter(CallingTsapParameter.class).isPresent(), is(true));
    }

    @Test
    @Category(FastTests.class)
    public void dataTpdu() {
        List<Parameter> parameters = Collections.emptyList();
        ByteBuf userData = Unpooled.buffer();

        userData.writeByte(0x66);

        DataTpdu tpdu = new DataTpdu(true, (byte) 0x7F, parameters, userData);

        assertThat(tpdu.getTpduCode(), equalTo(TpduCode.DATA));
        assertThat("Unexpected eot reference", tpdu.isEot(), is(true));
        assertThat(tpdu.getTpduRef(), equalTo((byte) 0x7F));
        assertThat("Unexpected parameters", tpdu.getParameters().isEmpty(), is(true));
        assertThat("Unexpected user data", tpdu.getUserData().readByte(), equalTo((byte) 0x66));
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

        assertThat(tpdu.getTpduCode(), equalTo(TpduCode.CONNECTION_REQUEST));
        assertThat("Unexpected destination reference", tpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat("Unexpected source reference", tpdu.getSourceReference(), equalTo((short) 0x2));
        assertThat(tpdu.getProtocolClass(), equalTo(ProtocolClass.CLASS_0));
        assertThat("Unexpected parameters", tpdu.getParameters().isEmpty(), is(true));
        assertThat("Unexpected user data", tpdu.getUserData().readByte(), equalTo((byte) 0x33));
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

        assertThat(tpdu.getTpduCode(), equalTo(TpduCode.CONNECTION_CONFIRM));
        assertThat("Unexpected destination reference", tpdu.getDestinationReference(), equalTo((short) 0x3));
        assertThat("Unexpected source reference", tpdu.getSourceReference(), equalTo((short) 0x4));
        assertThat(tpdu.getProtocolClass(), equalTo(ProtocolClass.CLASS_1));
        assertThat("Unexpected parameters", tpdu.getParameters().isEmpty(), is(true));
        assertThat("Unexpected user data", tpdu.getUserData().readByte(), equalTo((byte) 0x44));
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

        assertThat(tpdu.getTpduCode(), equalTo(TpduCode.DISCONNECT_REQUEST));
        assertThat("Unexpected destination reference", tpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat("Unexpected source reference", tpdu.getSourceReference(), equalTo((short) 0x2));
        assertThat(tpdu.getDisconnectReason(), equalTo(DisconnectReason.ADDRESS_UNKNOWN));
        assertThat("Unexpected parameters", tpdu.getParameters(), empty());
        assertThat("Unexpected user data", tpdu.getUserData().readByte(), equalTo((byte) 0x22));
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

        assertThat(tpdu.getTpduCode(), equalTo(TpduCode.DISCONNECT_CONFIRM));
        assertThat("Unexpected destination reference", tpdu.getDestinationReference(), equalTo((short) 0x3));
        assertThat("Unexpected source reference", tpdu.getSourceReference(), equalTo((short) 0x4));
        assertThat("Unexpected parameters", tpdu.getParameters().isEmpty(), is(true));
        assertThat("Unexpected user data", tpdu.getUserData().readByte(), equalTo((byte) 0x11));
    }

}