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
package org.apache.plc4x.java.isotp.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.isoontcp.netty.model.IsoOnTcpMessage;
import org.apache.plc4x.java.isotp.netty.model.IsoTPMessage;
import org.apache.plc4x.java.isotp.netty.model.params.*;
import org.apache.plc4x.java.isotp.netty.model.tpdus.*;
import org.apache.plc4x.java.isotp.netty.model.types.*;
import org.apache.plc4x.test.FastTests;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.ArrayList;
import java.util.Collections;

import static org.hamcrest.collection.IsCollectionWithSize.hasSize;
import static org.hamcrest.collection.IsEmptyCollection.empty;
import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;

public class IsoTPProtocolTest {

    private IsoTPProtocol isoTPProtocol;
    private ChannelHandlerContext ctx;
    private ByteBuf buf;
    private ArrayList<Object> out;

    @Before
    public void setup() {
        byte rackNo = 1;
        byte slotNo = 4;
        TpduSize tpduSize = TpduSize.SIZE_512;

        isoTPProtocol = new IsoTPProtocol(rackNo, slotNo, tpduSize);
        ctx = mock(ChannelHandlerContext.class, RETURNS_DEEP_STUBS);
        buf = Unpooled.buffer();
        out = new ArrayList<>();
    }

    @After
    public void terDown() {
        isoTPProtocol = null;
    }

    @Test
    @Category(FastTests.class)
    public void encodeConnectionRequest() throws Exception {
        ConnectionRequestTpdu tpdu = new ConnectionRequestTpdu((short) 0x1, (short) (0x2), ProtocolClass.CLASS_0, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(7));
        assertThat(userData.readByte(), equalTo((byte) 0x6));
        assertThat(userData.readByte(), equalTo(TpduCode.CONNECTION_REQUEST.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readShort(), equalTo((short) 0x2));
        assertThat(userData.readByte(), equalTo(ProtocolClass.CLASS_0.getCode()));
    }

    @Test
    @Category(FastTests.class)
    public void decodeConnectionRequest() throws Exception {
        buf.writeByte(0x6) // header length
            .writeByte(TpduCode.CONNECTION_REQUEST.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(ProtocolClass.CLASS_0.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ConnectionRequestTpdu requestTpdu = (ConnectionRequestTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(requestTpdu.getTpduCode(), equalTo(TpduCode.CONNECTION_REQUEST));
        assertThat(requestTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(requestTpdu.getSourceReference(), equalTo((short) 0x2));
        assertThat(requestTpdu.getProtocolClass(), equalTo(ProtocolClass.CLASS_0));
        assertThat(requestTpdu.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void encodeDisconnectionRequest() throws Exception {
        DisconnectRequestTpdu tpdu = new DisconnectRequestTpdu((short) 0x1, (short) (0x2), DisconnectReason.NORMAL, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));
        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(7));
        assertThat(userData.readByte(), equalTo((byte) 0x6));
        assertThat(userData.readByte(), equalTo(TpduCode.DISCONNECT_REQUEST.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readShort(), equalTo((short) 0x2));
        assertThat(userData.readByte(), equalTo(DisconnectReason.NORMAL.getCode()));
    }

    @Test
    @Category(FastTests.class)
    public void decodeDisconnectionRequest() throws Exception {
        buf.writeByte(0x6) // header length
            .writeByte(TpduCode.DISCONNECT_REQUEST.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(DisconnectReason.NORMAL.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        DisconnectRequestTpdu requestTpdu = (DisconnectRequestTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(requestTpdu.getTpduCode(), equalTo(TpduCode.DISCONNECT_REQUEST));
        assertThat(requestTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(requestTpdu.getSourceReference(), equalTo((short) 0x2));
        assertThat(requestTpdu.getDisconnectReason(), equalTo(DisconnectReason.NORMAL));
        assertThat(requestTpdu.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void encodeData() throws Exception {
        DataTpdu tpdu = new DataTpdu(true, (byte) 0x7, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(3));
        assertThat(userData.readByte(), equalTo((byte) 0x2));
        assertThat(userData.readByte(), equalTo(TpduCode.DATA.getCode()));
        assertThat(userData.readByte(), equalTo((byte) 0x87));
    }

    @Test
    @Category(FastTests.class)
    public void decodeDataEOT() throws Exception {
        buf.writeByte(0x3) // header length
            .writeByte(TpduCode.DATA.getCode())
            .writeByte((byte) 0x81); // Tpdu code + EOT
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        DataTpdu requestTpdu = (DataTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(requestTpdu.getTpduCode(), equalTo(TpduCode.DATA));
        assertThat(requestTpdu.getTpduRef(), equalTo((byte) 0x1));
        assertThat(requestTpdu.isEot(), is(true));
        assertThat(requestTpdu.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void decodeData() throws Exception {
        buf.writeByte(0x3) // header length
            .writeByte(TpduCode.DATA.getCode())
            .writeByte((byte) 0x1); // Tpdu code
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        DataTpdu requestTpdu = (DataTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(requestTpdu.getTpduCode(), equalTo(TpduCode.DATA));
        assertThat(requestTpdu.getTpduRef(), equalTo((byte) 0x1));
        assertThat(!requestTpdu.isEot(), is(true));
        assertThat(requestTpdu.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void encodeConnectionConfirm() throws Exception {
        ConnectionConfirmTpdu tpdu = new ConnectionConfirmTpdu((short) 0x1, (short) (0x2), ProtocolClass.CLASS_1, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(7));
        assertThat(userData.readByte(), equalTo((byte) 0x6));
        assertThat(userData.readByte(), equalTo(TpduCode.CONNECTION_CONFIRM.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readShort(), equalTo((short) 0x2));
        assertThat(userData.readByte(), equalTo(ProtocolClass.CLASS_1.getCode()));
    }

    @Test
    @Category(FastTests.class)
    public void decodeConnectionConfirm() throws Exception {
        buf.writeByte(0x6) // header length
            .writeByte(TpduCode.CONNECTION_CONFIRM.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(ProtocolClass.CLASS_0.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ConnectionConfirmTpdu requestTpdu = (ConnectionConfirmTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(requestTpdu.getTpduCode(), equalTo(TpduCode.CONNECTION_CONFIRM));
        assertThat(requestTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(requestTpdu.getSourceReference(), equalTo((short) 0x2));
        assertThat(requestTpdu.getProtocolClass(), equalTo(ProtocolClass.CLASS_0));
        assertThat(requestTpdu.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void encodeDisconnectionConfirm() throws Exception {
        DisconnectConfirmTpdu tpdu = new DisconnectConfirmTpdu((short) 0x1, (short) (0x2), Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(6));
        assertThat(userData.readByte(), equalTo((byte) 0x5));
        assertThat(userData.readByte(), equalTo(TpduCode.DISCONNECT_CONFIRM.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readShort(), equalTo((short) 0x2));
    }

    @Test
    @Category(FastTests.class)
    public void decodeDisconnectionConfirm() throws Exception {
        buf.writeByte(0x5) // header length
            .writeByte(TpduCode.DISCONNECT_CONFIRM.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(DisconnectReason.NORMAL.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        DisconnectConfirmTpdu requestTpdu = (DisconnectConfirmTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(requestTpdu.getTpduCode(), equalTo(TpduCode.DISCONNECT_CONFIRM));
        assertThat(requestTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(requestTpdu.getSourceReference(), equalTo((short) 0x2));
        assertThat(requestTpdu.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void encodeError() throws Exception {
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(5));
        assertThat(userData.readByte(), equalTo((byte) 0x4));
        assertThat(userData.readByte(), equalTo(TpduCode.TPDU_ERROR.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readByte(), equalTo(RejectCause.REASON_NOT_SPECIFIED.getCode()));
    }

    @Test
    @Category(FastTests.class)
    public void encodeCallingParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        CallingTsapParameter callingParameter = new CallingTsapParameter(DeviceGroup.PG_OR_PC, (byte) 0x7, (byte) 0xe1); // slot number too big and overflows into rack
        parmameters.add(callingParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(9));
        assertThat(userData.readByte(), equalTo((byte) 0x8));
        assertThat(userData.readByte(), equalTo(TpduCode.TPDU_ERROR.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readByte(), equalTo(RejectCause.REASON_NOT_SPECIFIED.getCode()));
        assertThat(userData.readByte(), equalTo(ParameterCode.CALLING_TSAP.getCode()));
        assertThat(userData.readByte(), equalTo((byte) 0x2));
        assertThat(userData.readByte(), equalTo(DeviceGroup.PG_OR_PC.getCode()));
        byte rackAndSlot = userData.readByte();
        assertThat((rackAndSlot & 0xf0) >> 4, equalTo(0x7));
        assertThat((rackAndSlot & 0x0f), equalTo(0xe1 & 0x0f));
    }

    @Test
    @Category(FastTests.class)
    public void encodeChecksumParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        ChecksumParameter checksumParameter = new ChecksumParameter((byte) 0x77);
        parmameters.add(checksumParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(8));
        assertThat(userData.readByte(), equalTo((byte) 0x7));
        assertThat(userData.readByte(), equalTo(TpduCode.TPDU_ERROR.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readByte(), equalTo(RejectCause.REASON_NOT_SPECIFIED.getCode()));
        assertThat(userData.readByte(), equalTo(ParameterCode.CHECKSUM.getCode()));
        assertThat(userData.readByte(), equalTo((byte) 0x1));
        assertThat(userData.readByte(), equalTo((byte) 0x77));
    }

    @Test
    @Category(FastTests.class)
    public void encodeAditionalInformationParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        byte[] data = {'O', 'p', 'p', 's'};
        DisconnectAdditionalInformationParameter informationParameter = new DisconnectAdditionalInformationParameter(data);
        parmameters.add(informationParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(11));
        assertThat(userData.readByte(), equalTo((byte) 0xA));
        assertThat(userData.readByte(), equalTo(TpduCode.TPDU_ERROR.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readByte(), equalTo(RejectCause.REASON_NOT_SPECIFIED.getCode()));
        assertThat(userData.readByte(), equalTo(ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION.getCode()));
        assertThat(userData.readByte(), equalTo((byte) 0x4));
        assertThat(userData.readByte(), equalTo((byte) 'O'));
        assertThat(userData.readByte(), equalTo((byte) 'p'));
        assertThat(userData.readByte(), equalTo((byte) 'p'));
        assertThat(userData.readByte(), equalTo((byte) 's'));
    }

    @Test
    @Category(FastTests.class)
    public void encodeSizeParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        TpduSizeParameter sizeParameter = new TpduSizeParameter(TpduSize.SIZE_512);
        parmameters.add(sizeParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertThat("Message not decoded", out, hasSize(1));

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertThat(userData.writerIndex(), equalTo(8));
        assertThat(userData.readByte(), equalTo((byte) 0x7));
        assertThat(userData.readByte(), equalTo(TpduCode.TPDU_ERROR.getCode()));
        assertThat(userData.readShort(), equalTo((short) 0x1));
        assertThat(userData.readByte(), equalTo(RejectCause.REASON_NOT_SPECIFIED.getCode()));
        assertThat(userData.readByte(), equalTo(ParameterCode.TPDU_SIZE.getCode()));
        assertThat(userData.readByte(), equalTo((byte) 0x1));
        assertThat(userData.readByte(), equalTo(TpduSize.SIZE_512.getCode()));
    }

    @Test
    @Category(FastTests.class)
    public void decodeError() throws Exception {
        buf.writeByte(0x4) // header length
            .writeByte(TpduCode.TPDU_ERROR.getCode())
            .writeShort(0x01) // destination reference
            .writeByte(RejectCause.REASON_NOT_SPECIFIED.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(errorTpdu.getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        assertThat(errorTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(errorTpdu.getRejectCause(), equalTo(RejectCause.REASON_NOT_SPECIFIED));
        assertThat(errorTpdu.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void encodeNullRequest() throws Exception {
        ConnectionRequestTpdu tpdu = null;

        isoTPProtocol.encode(ctx, tpdu, out);
        assertThat("Message not decoded", out, empty());

        isoTPProtocol.encode(ctx, null, out);
        assertThat("Message not decoded", out, empty());
    }


    @Test
    @Category(FastTests.class)
    public void decodeNull() throws Exception {
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);
        assertThat("Message not decoded", out, empty());

        isoTPProtocol.decode(ctx, null, out);
        assertThat("Message not decoded", out, empty());
    }

    @Test
    @Category(FastTests.class)
    public void encodeUnsupported() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        CustomTpdu tpdu = new CustomTpdu((byte) 0x7F, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);
        assertThat("Message not decoded", out, empty());
    }


    @Test
    @Category(FastTests.class)
    public void decodeUnsupported() throws Exception {
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);
        buf.writeByte(0x3) // header length
            .writeByte(0x7F)
            .writeShort(0x01); // destination reference
        isoTPProtocol.decode(ctx, in, out);
        assertThat("Message not decoded", out, empty());
    }

    @Test
    @Category(FastTests.class)
    public void decodeCallingParameter() throws Exception {
        buf.writeByte(0x8) // header length
            .writeByte(TpduCode.TPDU_ERROR.getCode())
            .writeShort(0x01) // destination reference
            .writeByte(RejectCause.REASON_NOT_SPECIFIED.getCode()) // reject clause
            .writeByte(ParameterCode.CALLING_TSAP.getCode())
            .writeByte(0x2) // parameter length
            .writeByte(DeviceGroup.PG_OR_PC.getCode())
            .writeByte((byte) ((0x1 << 4) | 0x7));
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(errorTpdu.getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        assertThat(errorTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(errorTpdu.getRejectCause(), equalTo(RejectCause.REASON_NOT_SPECIFIED));
        assertThat(errorTpdu.getParameters(), hasSize(1));
        CallingTsapParameter parameter = (CallingTsapParameter) errorTpdu.getParameters().get(0);
        assertThat(parameter.getType(), equalTo(ParameterCode.CALLING_TSAP));
        assertThat(parameter.getDeviceGroup(), equalTo(DeviceGroup.PG_OR_PC));
        assertThat(parameter.getRackNumber(), equalTo((byte) 0x1));
        assertThat(parameter.getSlotNumber(), equalTo((byte) 0x7));
    }

    @Test
    @Category(FastTests.class)
    public void decodeCalledParameter() throws Exception {
        buf.writeByte(0x8) // header length
            .writeByte(TpduCode.TPDU_ERROR.getCode())
            .writeShort(0x01) // destination reference
            .writeByte(RejectCause.REASON_NOT_SPECIFIED.getCode()) // reject clause
            .writeByte(ParameterCode.CALLED_TSAP.getCode())
            .writeByte(0x2) // parameter length
            .writeByte(DeviceGroup.PG_OR_PC.getCode())
            .writeByte((byte) ((0x2 << 4) | 0x3));
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(errorTpdu.getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        assertThat(errorTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(errorTpdu.getRejectCause(), equalTo(RejectCause.REASON_NOT_SPECIFIED));
        assertThat(errorTpdu.getParameters(), hasSize(1));
        CalledTsapParameter parameter = (CalledTsapParameter) errorTpdu.getParameters().get(0);
        assertThat(parameter.getType(), equalTo(ParameterCode.CALLED_TSAP));
        assertThat(parameter.getDeviceGroup(), equalTo(DeviceGroup.PG_OR_PC));
        assertThat(parameter.getRackNumber(), equalTo((byte) 0x2));
        assertThat(parameter.getSlotNumber(), equalTo((byte) 0x3));
    }

    @Test
    @Category(FastTests.class)
    public void decodeChecksumParameter() throws Exception {
        buf.writeByte(0x8) // header length
            .writeByte(TpduCode.TPDU_ERROR.getCode())
            .writeShort(0x01) // destination reference
            .writeByte(RejectCause.REASON_NOT_SPECIFIED.getCode()) // reject clause
            .writeByte(ParameterCode.CHECKSUM.getCode())
            .writeByte(0x1) // parameter length
            .writeByte(0x33); // checksum
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(errorTpdu.getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        assertThat(errorTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(errorTpdu.getRejectCause(), equalTo(RejectCause.REASON_NOT_SPECIFIED));
        assertThat(errorTpdu.getParameters(), hasSize(1));
        ChecksumParameter parameter = (ChecksumParameter) errorTpdu.getParameters().get(0);
        assertThat(parameter.getType(), equalTo(ParameterCode.CHECKSUM));
        assertThat(parameter.getChecksum(), equalTo((byte) 0x33));
    }

    @Test
    @Category(FastTests.class)
    public void decodeSizeParameter() throws Exception {
        buf.writeByte(0x8) // header length
            .writeByte(TpduCode.TPDU_ERROR.getCode())
            .writeShort(0x01) // destination reference
            .writeByte(RejectCause.REASON_NOT_SPECIFIED.getCode()) // reject clause
            .writeByte(ParameterCode.TPDU_SIZE.getCode())
            .writeByte(0x1) // parameter length
            .writeByte(TpduSize.SIZE_256.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(errorTpdu.getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        assertThat(errorTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(errorTpdu.getRejectCause(), equalTo(RejectCause.REASON_NOT_SPECIFIED));
        assertThat(errorTpdu.getParameters(), hasSize(1));
        TpduSizeParameter parameter = (TpduSizeParameter) errorTpdu.getParameters().get(0);
        assertThat(parameter.getType(), equalTo(ParameterCode.TPDU_SIZE));
        assertThat(parameter.getTpduSize(), equalTo(TpduSize.SIZE_256));
    }

    @Test
    @Category(FastTests.class)
    public void decodeAdditionalInformationParameter() throws Exception {
        buf.writeByte(0x8) // header length
            .writeByte(TpduCode.TPDU_ERROR.getCode())
            .writeShort(0x01) // destination reference
            .writeByte(RejectCause.REASON_NOT_SPECIFIED.getCode()) // reject clause
            .writeByte(ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION.getCode())
            .writeByte(0x5) // parameter length
            .writeByte('E')
            .writeByte('r')
            .writeByte('r')
            .writeByte('o')
            .writeByte('r');
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertThat("Message not decoded", out, hasSize(1));

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertThat(errorTpdu.getTpduCode(), equalTo(TpduCode.TPDU_ERROR));
        assertThat(errorTpdu.getDestinationReference(), equalTo((short) 0x1));
        assertThat(errorTpdu.getRejectCause(), equalTo(RejectCause.REASON_NOT_SPECIFIED));
        assertThat(errorTpdu.getParameters(), hasSize(1));
        DisconnectAdditionalInformationParameter parameter = (DisconnectAdditionalInformationParameter) errorTpdu.getParameters().get(0);
        assertThat(parameter.getType(), equalTo(ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION));
        byte[] data = parameter.getData();
        assertThat(data[0], equalTo((byte) 'E'));
        assertThat(data[1], equalTo((byte) 'r'));
        assertThat(data[2], equalTo((byte) 'r'));
        assertThat(data[3], equalTo((byte) 'o'));
        assertThat(data[4], equalTo((byte) 'r'));
    }

}
