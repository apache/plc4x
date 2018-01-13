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
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;

public class IsoTPProtocolTest {

    private IsoTPProtocol isoTPProtocol;
    private ChannelHandlerContext ctx;
    private ByteBuf buf;
    private ArrayList<Object> out;

    @BeforeEach
    void setup() {
        byte rackNo = 1;
        byte slotNo = 4;
        TpduSize tpduSize = TpduSize.SIZE_512;

        isoTPProtocol = new IsoTPProtocol(rackNo, slotNo, tpduSize);
        ctx = mock(ChannelHandlerContext.class, RETURNS_DEEP_STUBS);
        buf = Unpooled.buffer();
        out = new ArrayList<>();
    }

    @AfterEach
    void terDown() {
        isoTPProtocol = null;
    }

    @Test
    @Tag("fast")
    public void encodeConnectionRequest() throws Exception {
        ConnectionRequestTpdu tpdu = new ConnectionRequestTpdu((short) 0x1, (short) (0x2), ProtocolClass.CLASS_0, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 7, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x6, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.CONNECTION_REQUEST.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readShort() == (short) 0x2, "Incorrect source reference code");
        assertTrue(userData.readByte() == ProtocolClass.CLASS_0.getCode(), "Incorrect protocol class");
    }

    @Test
    @Tag("fast")
    public void decodeConnectionRequest() throws Exception {
        buf.writeByte(0x6) // header length
            .writeByte(TpduCode.CONNECTION_REQUEST.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(ProtocolClass.CLASS_0.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ConnectionRequestTpdu requestTpdu = (ConnectionRequestTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.CONNECTION_REQUEST, "Message code not correct");
        assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
        assertTrue(requestTpdu.getProtocolClass() == ProtocolClass.CLASS_0, "Message protocol class reference not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeDisconnectionRequest() throws Exception {
        DisconnectRequestTpdu tpdu = new DisconnectRequestTpdu((short) 0x1, (short) (0x2), DisconnectReason.NORMAL, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");
        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 7, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x6, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.DISCONNECT_REQUEST.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readShort() == (short) 0x2, "Incorrect source reference code");
        assertTrue(userData.readByte() == DisconnectReason.NORMAL.getCode(), "Incorrect disconnect reason");
    }

    @Test
    @Tag("fast")
    public void decodeDisconnectionRequest() throws Exception {
        buf.writeByte(0x6) // header length
            .writeByte(TpduCode.DISCONNECT_REQUEST.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(DisconnectReason.NORMAL.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        DisconnectRequestTpdu requestTpdu = (DisconnectRequestTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.DISCONNECT_REQUEST, "Message code not correct");
        assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
        assertTrue(requestTpdu.getDisconnectReason() == DisconnectReason.NORMAL, "Message disconnect reason not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeData() throws Exception {
        DataTpdu tpdu = new DataTpdu(true, (byte) 0x7, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 3, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x2, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.DATA.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readByte() == (byte) 0x87, "Incorrect Tpdu reference and EOT");
    }

    @Test
    @Tag("fast")
    public void decodeDataEOT() throws Exception {
        buf.writeByte(0x3) // header length
            .writeByte(TpduCode.DATA.getCode())
            .writeByte((byte) 0x81); // Tpdu code + EOT
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        DataTpdu requestTpdu = (DataTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.DATA, "Message code not correct");
        assertTrue(requestTpdu.getTpduRef() == (short) 0x1, "Message Tpdu reference not correct");
        assertTrue(requestTpdu.isEot(), "Message EOT not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void decodeData() throws Exception {
        buf.writeByte(0x3) // header length
            .writeByte(TpduCode.DATA.getCode())
            .writeByte((byte) 0x1); // Tpdu code
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        DataTpdu requestTpdu = (DataTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.DATA, "Message code not correct");
        assertTrue(requestTpdu.getTpduRef() == (short) 0x1, "Message Tpdu reference not correct");
        assertTrue(!requestTpdu.isEot(), "Message EOT not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeConnectionConfirm() throws Exception {
        ConnectionConfirmTpdu tpdu = new ConnectionConfirmTpdu((short) 0x1, (short) (0x2), ProtocolClass.CLASS_1, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 7, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x6, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.CONNECTION_CONFIRM.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readShort() == (short) 0x2, "Incorrect source reference code");
        assertTrue(userData.readByte() == ProtocolClass.CLASS_1.getCode(), "Incorrect protocol class");
    }

    @Test
    @Tag("fast")
    public void decodeConnectionConfirm() throws Exception {
        buf.writeByte(0x6) // header length
            .writeByte(TpduCode.CONNECTION_CONFIRM.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(ProtocolClass.CLASS_0.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ConnectionConfirmTpdu requestTpdu = (ConnectionConfirmTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.CONNECTION_CONFIRM, "Message code not correct");
        assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
        assertTrue(requestTpdu.getProtocolClass() == ProtocolClass.CLASS_0, "Message protocol class reference not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeDisconnectionConfirm() throws Exception {
        DisconnectConfirmTpdu tpdu = new DisconnectConfirmTpdu((short) 0x1, (short) (0x2), Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 6, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x5, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.DISCONNECT_CONFIRM.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readShort() == (short) 0x2, "Incorrect source reference code");
    }

    @Test
    @Tag("fast")
    public void decodeDisconnectionConfirm() throws Exception {
        buf.writeByte(0x5) // header length
            .writeByte(TpduCode.DISCONNECT_CONFIRM.getCode())
            .writeShort(0x01) // destination reference
            .writeShort(0x02) // source reference
            .writeByte(DisconnectReason.NORMAL.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        DisconnectConfirmTpdu requestTpdu = (DisconnectConfirmTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.DISCONNECT_CONFIRM, "Message code not correct");
        assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeError() throws Exception {
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, Collections.emptyList(), buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 5, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x4, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.TPDU_ERROR.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readByte() == RejectCause.REASON_NOT_SPECIFIED.getCode(), "Incorrect reject cause code");
    }

    @Test
    @Tag("fast")
    public void encodeCallingParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        CallingTsapParameter callingParameter = new CallingTsapParameter(DeviceGroup.PG_OR_PC, (byte) 0x7, (byte) 0xe1); // slot number too big and overflows into rack
        parmameters.add(callingParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 9, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x8, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.TPDU_ERROR.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readByte() == RejectCause.REASON_NOT_SPECIFIED.getCode(), "Incorrect reject cause code");
        assertTrue(userData.readByte() == ParameterCode.CALLING_TSAP.getCode(), "Incorrect parameter code");
        assertTrue(userData.readByte() == (byte) 0x2, "Incorrect parameter length");
        assertTrue(userData.readByte() == DeviceGroup.PG_OR_PC.getCode(), "Incorrect device group code");
        byte rackAndSlot = userData.readByte();
        assertTrue((rackAndSlot & 0xf0) >> 4 == 0x7, "Incorrect rack number");
        assertTrue((rackAndSlot & 0x0f) == (0xe1 & 0x0f), "Incorrect slot number");
    }

    @Test
    @Tag("fast")
    public void encodeChecksumParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        ChecksumParameter checksumParameter = new ChecksumParameter((byte) 0x77);
        parmameters.add(checksumParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 8, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x7, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.TPDU_ERROR.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readByte() == RejectCause.REASON_NOT_SPECIFIED.getCode(), "Incorrect reject cause code");
        assertTrue(userData.readByte() == ParameterCode.CHECKSUM.getCode(), "Incorrect parameter code");
        assertTrue(userData.readByte() == (byte) 0x1, "Incorrect parameter length");
        assertTrue(userData.readByte() == 0x77, "Incorrect checksum");
    }

    @Test
    @Tag("fast")
    public void encodeAditionalInformationParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        byte[] data = {'O', 'p', 'p', 's'};
        DisconnectAdditionalInformationParameter informationParameter = new DisconnectAdditionalInformationParameter(data);
        parmameters.add(informationParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 11, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0xA, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.TPDU_ERROR.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readByte() == RejectCause.REASON_NOT_SPECIFIED.getCode(), "Incorrect reject cause code");
        assertTrue(userData.readByte() == ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION.getCode(), "Incorrect parameter code");
        assertTrue(userData.readByte() == (byte) 0x4, "Incorrect parameter length");
        assertTrue(userData.readByte() == 'O', "Incorrect data");
        assertTrue(userData.readByte() == 'p', "Incorrect data");
        assertTrue(userData.readByte() == 'p', "Incorrect data");
        assertTrue(userData.readByte() == 's', "Incorrect data");
    }

    @Test
    @Tag("fast")
    public void encodeSizeParameter() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        TpduSizeParameter sizeParameter = new TpduSizeParameter(TpduSize.SIZE_512);
        parmameters.add(sizeParameter);
        ErrorTpdu tpdu = new ErrorTpdu((short) 0x1, RejectCause.REASON_NOT_SPECIFIED, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage) out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 8, "Incorrect message length");
        assertTrue(userData.readByte() == (byte) 0x7, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.TPDU_ERROR.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short) 0x1, "Incorrect destination reference code");
        assertTrue(userData.readByte() == RejectCause.REASON_NOT_SPECIFIED.getCode(), "Incorrect reject cause code");
        assertTrue(userData.readByte() == ParameterCode.TPDU_SIZE.getCode(), "Incorrect parameter code");
        assertTrue(userData.readByte() == (byte) 0x1, "Incorrect parameter length");
        assertTrue(userData.readByte() == TpduSize.SIZE_512.getCode(), "Incorrect tdpu size");
    }

    @Test
    @Tag("fast")
    public void decodeError() throws Exception {
        buf.writeByte(0x4) // header length
            .writeByte(TpduCode.TPDU_ERROR.getCode())
            .writeShort(0x01) // destination reference
            .writeByte(RejectCause.REASON_NOT_SPECIFIED.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(errorTpdu.getTpduCode() == TpduCode.TPDU_ERROR, "Message code not correct");
        assertTrue(errorTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(errorTpdu.getRejectCause() == RejectCause.REASON_NOT_SPECIFIED, "Message reject cause not correct");
        assertTrue(errorTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeNullRequest() throws Exception {
        ConnectionRequestTpdu tpdu = null;

        isoTPProtocol.encode(ctx, tpdu, out);
        assertTrue(out.size() == 0, "Message decoded when null passed");

        isoTPProtocol.encode(ctx, null, out);
        assertTrue(out.size() == 0, "Message decoded when null passed");
    }


    @Test
    @Tag("fast")
    public void decodeNull() throws Exception {
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);
        assertTrue(out.size() == 0, "Message decoded when blank message passed");

        isoTPProtocol.decode(ctx, null, out);
        assertTrue(out.size() == 0, "Message decoded when blank message passed");
    }

    @Test
    @Tag("fast")
    public void encodeUnsupported() throws Exception {
        ArrayList<Parameter> parmameters = new ArrayList<>();
        CustomTpdu tpdu = new CustomTpdu((byte)0x7F, parmameters, buf);

        isoTPProtocol.encode(ctx, tpdu, out);
        assertTrue(out.size() == 0, "Message encoded when unsupported Tpdu code passed");
   }


    @Test
    @Tag("fast")
    public void decodeUnsupported() throws Exception {
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);
        buf.writeByte(0x3) // header length
            .writeByte(0x7F)
            .writeShort(0x01); // destination reference
        isoTPProtocol.decode(ctx, in, out);
        assertTrue(out.size() == 0, "Message decoded when unsupported Tpdu code passed");
    }

    @Test
    @Tag("fast")
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

        assertTrue(out.size() == 1, "Message not decoded");

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(errorTpdu.getTpduCode() == TpduCode.TPDU_ERROR, "Message code not correct");
        assertTrue(errorTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(errorTpdu.getRejectCause() == RejectCause.REASON_NOT_SPECIFIED, "Message reject cause not correct");
        assertTrue(errorTpdu.getParameters().size() == 1, "Incorrect number of parameters");
        CallingTsapParameter parameter = (CallingTsapParameter) errorTpdu.getParameters().get(0);
        assertTrue(parameter.getType() == ParameterCode.CALLING_TSAP, "Parameter type incorrect");
        assertTrue(parameter.getDeviceGroup() == DeviceGroup.PG_OR_PC, "Device group incorrect");
        assertTrue(parameter.getRackNumber() == 0x1, "Rack number incorrect");
        assertTrue(parameter.getSlotNumber() == 0x7, "Slot number incorrect");
    }

    @Test
    @Tag("fast")
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

        assertTrue(out.size() == 1, "Message not decoded");

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(errorTpdu.getTpduCode() == TpduCode.TPDU_ERROR, "Message code not correct");
        assertTrue(errorTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(errorTpdu.getRejectCause() == RejectCause.REASON_NOT_SPECIFIED, "Message reject cause not correct");
        assertTrue(errorTpdu.getParameters().size() == 1, "Incorrect number of parameters");
        CalledTsapParameter parameter = (CalledTsapParameter) errorTpdu.getParameters().get(0);
        assertTrue(parameter.getType() == ParameterCode.CALLED_TSAP, "Parameter type incorrect");
        assertTrue(parameter.getDeviceGroup() == DeviceGroup.PG_OR_PC, "Device group incorrect");
        assertTrue(parameter.getRackNumber() == 0x2, "Rack number incorrect");
        assertTrue(parameter.getSlotNumber() == 0x3, "Slot number incorrect");
    }

    @Test
    @Tag("fast")
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

        assertTrue(out.size() == 1, "Message not decoded");

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(errorTpdu.getTpduCode() == TpduCode.TPDU_ERROR, "Message code not correct");
        assertTrue(errorTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(errorTpdu.getRejectCause() == RejectCause.REASON_NOT_SPECIFIED, "Message reject cause not correct");
        assertTrue(errorTpdu.getParameters().size() == 1, "Incorrect number of parameters");
        ChecksumParameter parameter = (ChecksumParameter) errorTpdu.getParameters().get(0);
        assertTrue(parameter.getType() == ParameterCode.CHECKSUM, "Parameter type incorrect");
        assertTrue(parameter.getChecksum() == 0x33, "Checksum incorrect");
    }

    @Test
    @Tag("fast")
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

        assertTrue(out.size() == 1, "Message not decoded");

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(errorTpdu.getTpduCode() == TpduCode.TPDU_ERROR, "Message code not correct");
        assertTrue(errorTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(errorTpdu.getRejectCause() == RejectCause.REASON_NOT_SPECIFIED, "Message reject cause not correct");
        assertTrue(errorTpdu.getParameters().size() == 1, "Incorrect number of parameters");
        TpduSizeParameter parameter = (TpduSizeParameter) errorTpdu.getParameters().get(0);
        assertTrue(parameter.getType() == ParameterCode.TPDU_SIZE, "Parameter type incorrect");
        assertTrue(parameter.getTpduSize() == TpduSize.SIZE_256, "Size incorrect");
    }

    @Test
    @Tag("fast")
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

        assertTrue(out.size() == 1, "Message not decoded");

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage) out.get(0)).getTpdu();

        assertTrue(errorTpdu.getTpduCode() == TpduCode.TPDU_ERROR, "Message code not correct");
        assertTrue(errorTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(errorTpdu.getRejectCause() == RejectCause.REASON_NOT_SPECIFIED, "Message reject cause not correct");
        assertTrue(errorTpdu.getParameters().size() == 1, "Incorrect number of parameters");
        DisconnectAdditionalInformationParameter parameter = (DisconnectAdditionalInformationParameter) errorTpdu.getParameters().get(0);
        assertTrue(parameter.getType() == ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION, "Parameter type incorrect");
        byte[] data = parameter.getData();
        assertTrue(data[0] == 'E', "Data incorrect");
        assertTrue(data[1] == 'r', "Data incorrect");
        assertTrue(data[2] == 'r', "Data incorrect");
        assertTrue(data[3] == 'o', "Data incorrect");
        assertTrue(data[4] == 'r', "Data incorrect");
    }
}
