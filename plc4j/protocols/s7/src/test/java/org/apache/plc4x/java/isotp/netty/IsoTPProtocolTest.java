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
import org.apache.plc4x.java.isotp.netty.model.params.Parameter;
import org.apache.plc4x.java.isotp.netty.model.tpdus.*;
import org.apache.plc4x.java.isotp.netty.model.types.*;
import org.apache.plc4x.java.netty.NettyTestBase;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertTrue;

public class IsoTPProtocolTest extends NettyTestBase {

    private  IsoTPProtocol isoTPProtocol;

    @BeforeEach
    void setup() {
        byte rackNo = 1;
        byte slotNo = 4;
        TpduSize tpduSize = TpduSize.SIZE_512;

        isoTPProtocol = new IsoTPProtocol(rackNo, slotNo, tpduSize);
    }

    @AfterEach
    void terDown() {
        isoTPProtocol = null;
    }

    @Test
    @Tag("fast")
    public void encodeConnectionRequest() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ConnectionRequestTpdu tpdu = new ConnectionRequestTpdu((short)0x1, (short)(0x2), ProtocolClass.CLASS_0, Collections.emptyList(), buf);
        ArrayList<Object> out = new ArrayList<>();

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage)out.get(0)).getUserData();
        
        assertTrue(userData.writerIndex() == 7, "Incorrect message length");
        assertTrue(userData.readByte() == (byte)0x6, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.CONNECTION_REQUEST.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short)0x1, "Incorrect destination reference code");
        assertTrue(userData.readShort() == (short)0x2, "Incorrect source reference code");
        assertTrue(userData.readByte() == ProtocolClass.CLASS_0.getCode(), "Incorrect protocol class");
    }

    @Test
    @Tag("fast")
    public void decodeConnectionRequest() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ArrayList<Object> out = new ArrayList<>();

        buf.writeByte(0x6); // header length
        buf.writeByte(TpduCode.CONNECTION_REQUEST.getCode());
        buf.writeShort(0x01); // destination reference
        buf.writeShort(0x02); // source reference
        buf.writeByte(ProtocolClass.CLASS_0.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ConnectionRequestTpdu requestTpdu = (ConnectionRequestTpdu) ((IsoTPMessage)out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.CONNECTION_REQUEST, "Message code not correct");
        assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
        assertTrue(requestTpdu.getProtocolClass() == ProtocolClass.CLASS_0, "Message protocol class reference not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeDisconnectionRequest() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        DisconnectRequestTpdu tpdu = new DisconnectRequestTpdu((short)0x1, (short)(0x2), DisconnectReason.NORMAL, Collections.emptyList(), buf);
        ArrayList<Object> out = new ArrayList<>();

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");
        ByteBuf userData = ((IsoOnTcpMessage)out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 7, "Incorrect message length");
        assertTrue(userData.readByte() == (byte)0x6, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.DISCONNECT_REQUEST.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short)0x1, "Incorrect destination reference code");
        assertTrue(userData.readShort() == (short)0x2, "Incorrect source reference code");
        assertTrue(userData.readByte() == DisconnectReason.NORMAL.getCode(), "Incorrect disconnect reason");
    }

    @Test
    @Tag("fast")
    public void decodeDisconnectionRequest() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ArrayList<Object> out = new ArrayList<>();

        buf.writeByte(0x6); // header length
        buf.writeByte(TpduCode.DISCONNECT_REQUEST.getCode());
        buf.writeShort(0x01); // destination reference
        buf.writeShort(0x02); // source reference
        buf.writeByte(DisconnectReason.NORMAL.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        DisconnectRequestTpdu requestTpdu = (DisconnectRequestTpdu) ((IsoTPMessage)out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.DISCONNECT_REQUEST, "Message code not correct");
        assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
        assertTrue(requestTpdu.getDisconnectReason() == DisconnectReason.NORMAL, "Message disconnect reason not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeData() throws Exception {
      ChannelHandlerContext ctx = new MockChannelHandlerContext();
      ByteBuf buf = Unpooled.buffer();
      DataTpdu tpdu = new DataTpdu(true, (byte)0x7, Collections.emptyList(), buf);
      ArrayList<Object> out = new ArrayList<>();

      isoTPProtocol.encode(ctx, tpdu, out);

      assertTrue(out.size() == 1, "Message not decoded");

      ByteBuf userData = ((IsoOnTcpMessage)out.get(0)).getUserData();

      assertTrue(userData.writerIndex() == 3, "Incorrect message length");
      assertTrue(userData.readByte() == (byte)0x2, "Incorrect header length");
      assertTrue(userData.readByte() == TpduCode.DATA.getCode(), "Incorrect Tpdu code");
      assertTrue(userData.readByte() == (byte)0x87, "Incorrect Tpdu reference and EOT");
    }

    @Test
    @Tag("fast")
    public void decodeDataEOT() throws Exception {
      ChannelHandlerContext ctx = new MockChannelHandlerContext();
      ByteBuf buf = Unpooled.buffer();
      ArrayList<Object> out = new ArrayList<>();

      buf.writeByte(0x3); // header length
      buf.writeByte(TpduCode.DATA.getCode());
      buf.writeByte((byte) 0x81); // Tpdu code + EOT
      IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

      isoTPProtocol.decode(ctx, in, out);

      assertTrue(out.size() == 1, "Message not decoded");

      DataTpdu requestTpdu = (DataTpdu) ((IsoTPMessage)out.get(0)).getTpdu();

      assertTrue(requestTpdu.getTpduCode() == TpduCode.DATA, "Message code not correct");
      assertTrue(requestTpdu.getTpduRef() == (short) 0x1, "Message Tpdu reference not correct");
      assertTrue(requestTpdu.isEot(), "Message EOT not correct");
      assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void decodeData() throws Exception {
      ChannelHandlerContext ctx = new MockChannelHandlerContext();
      ByteBuf buf = Unpooled.buffer();
      ArrayList<Object> out = new ArrayList<>();

      buf.writeByte(0x3); // header length
      buf.writeByte(TpduCode.DATA.getCode());
      buf.writeByte((byte) 0x1); // Tpdu code
      IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

      isoTPProtocol.decode(ctx, in, out);

      assertTrue(out.size() == 1, "Message not decoded");

      DataTpdu requestTpdu = (DataTpdu) ((IsoTPMessage)out.get(0)).getTpdu();

      assertTrue(requestTpdu.getTpduCode() == TpduCode.DATA, "Message code not correct");
      assertTrue(requestTpdu.getTpduRef() == (short) 0x1, "Message Tpdu reference not correct");
      assertTrue(!requestTpdu.isEot(), "Message EOT not correct");
      assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeConnectionConfirm() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ConnectionConfirmTpdu tpdu = new ConnectionConfirmTpdu((short)0x1, (short)(0x2), ProtocolClass.CLASS_1, Collections.emptyList(), buf);
        ArrayList<Object> out = new ArrayList<>();

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage)out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 7, "Incorrect message length");
        assertTrue(userData.readByte() == (byte)0x6, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.CONNECTION_CONFIRM.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short)0x1, "Incorrect destination reference code");
        assertTrue(userData.readShort() == (short)0x2, "Incorrect source reference code");
        assertTrue(userData.readByte() == ProtocolClass.CLASS_1.getCode(), "Incorrect protocol class");
    }

    @Test
    @Tag("fast")
    public void decodeConnectionConfirm() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ArrayList<Object> out = new ArrayList<>();

        buf.writeByte(0x6); // header length
        buf.writeByte(TpduCode.CONNECTION_CONFIRM.getCode());
        buf.writeShort(0x01); // destination reference
        buf.writeShort(0x02); // source reference
        buf.writeByte(ProtocolClass.CLASS_0.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ConnectionConfirmTpdu requestTpdu = (ConnectionConfirmTpdu) ((IsoTPMessage)out.get(0)).getTpdu();

        assertTrue(requestTpdu.getTpduCode() == TpduCode.CONNECTION_CONFIRM, "Message code not correct");
        assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
        assertTrue(requestTpdu.getProtocolClass() == ProtocolClass.CLASS_0, "Message protocol class reference not correct");
        assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
      @Tag("fast")
      public void encodeDisconnectionConfirm() throws Exception {
          ChannelHandlerContext ctx = new MockChannelHandlerContext();
          ByteBuf buf = Unpooled.buffer();
          DisconnectConfirmTpdu tpdu = new DisconnectConfirmTpdu((short)0x1, (short)(0x2), Collections.emptyList(), buf);
          ArrayList<Object> out = new ArrayList<>();

          isoTPProtocol.encode(ctx, tpdu, out);

          assertTrue(out.size() == 1, "Message not decoded");

          ByteBuf userData = ((IsoOnTcpMessage)out.get(0)).getUserData();

          assertTrue(userData.writerIndex() == 6, "Incorrect message length");
          assertTrue(userData.readByte() == (byte)0x5, "Incorrect header length");
          assertTrue(userData.readByte() == TpduCode.DISCONNECT_CONFIRM.getCode(), "Incorrect Tpdu code");
          assertTrue(userData.readShort() == (short)0x1, "Incorrect destination reference code");
          assertTrue(userData.readShort() == (short)0x2, "Incorrect source reference code");
      }

      @Test
      @Tag("fast")
      public void decodeDisconnectionConfirm() throws Exception {
          ChannelHandlerContext ctx = new MockChannelHandlerContext();
          ByteBuf buf = Unpooled.buffer();
          ArrayList<Object> out = new ArrayList<>();

          buf.writeByte(0x5); // header length
          buf.writeByte(TpduCode.DISCONNECT_CONFIRM.getCode());
          buf.writeShort(0x01); // destination reference
          buf.writeShort(0x02); // source reference
          buf.writeByte(DisconnectReason.NORMAL.getCode());
          IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

          isoTPProtocol.decode(ctx, in, out);

          assertTrue(out.size() == 1, "Message not decoded");

          DisconnectConfirmTpdu requestTpdu = (DisconnectConfirmTpdu) ((IsoTPMessage)out.get(0)).getTpdu();

          assertTrue(requestTpdu.getTpduCode() == TpduCode.DISCONNECT_CONFIRM, "Message code not correct");
          assertTrue(requestTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
          assertTrue(requestTpdu.getSourceReference() == (short) 0x2, "Message source reference not correct");
          assertTrue(requestTpdu.getParameters().isEmpty(), "Message contains paramaters");
      }

    @Test
    @Tag("fast")
    public void encodeError() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ErrorTpdu tpdu = new ErrorTpdu((short)0x1, RejectCause.REASON_NOT_SPECIFIED, Collections.emptyList(), buf);
        ArrayList<Object> out = new ArrayList<>();

        isoTPProtocol.encode(ctx, tpdu, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ByteBuf userData = ((IsoOnTcpMessage)out.get(0)).getUserData();

        assertTrue(userData.writerIndex() == 5, "Incorrect message length");
        assertTrue(userData.readByte() == (byte)0x4, "Incorrect header length");
        assertTrue(userData.readByte() == TpduCode.TPDU_ERROR.getCode(), "Incorrect Tpdu code");
        assertTrue(userData.readShort() == (short)0x1, "Incorrect destination reference code");
        assertTrue(userData.readByte() == RejectCause.REASON_NOT_SPECIFIED.getCode(), "Incorrect reject cause code");
    }

    @Test
    @Tag("fast")
    public void decodeError() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ArrayList<Object> out = new ArrayList<>();

        buf.writeByte(0x6); // header length
        buf.writeByte(TpduCode.TPDU_ERROR.getCode());
        buf.writeShort(0x01); // destination reference
        buf.writeShort(0x02); // source reference
        buf.writeByte(ProtocolClass.CLASS_0.getCode());
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);

        assertTrue(out.size() == 1, "Message not decoded");

        ErrorTpdu errorTpdu = (ErrorTpdu) ((IsoTPMessage)out.get(0)).getTpdu();

        assertTrue(errorTpdu.getTpduCode() == TpduCode.TPDU_ERROR, "Message code not correct");
        assertTrue(errorTpdu.getDestinationReference() == (short) 0x1, "Message destination reference not correct");
        assertTrue(errorTpdu.getRejectCause() == RejectCause.REASON_NOT_SPECIFIED, "Message reject cause not correct");
        assertTrue(errorTpdu.getParameters().isEmpty(), "Message contains paramaters");
    }

    @Test
    @Tag("fast")
    public void encodeNullRequest() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ConnectionRequestTpdu tpdu =  null;
        ArrayList<Object> out = new ArrayList<>();

        isoTPProtocol.encode(ctx, tpdu, out);
        assertTrue(out.size() == 0, "Message decoded when null passed");

        isoTPProtocol.encode(ctx, null, out);
        assertTrue(out.size() == 0, "Message decoded when null passed");
    }


    @Test
    @Tag("fast")
    public void decodeNull() throws Exception {
        ChannelHandlerContext ctx = new MockChannelHandlerContext();
        ByteBuf buf = Unpooled.buffer();
        ArrayList<Object> out = new ArrayList<>();
        IsoOnTcpMessage in = new IsoOnTcpMessage(buf);

        isoTPProtocol.decode(ctx, in, out);
        assertTrue(out.size() == 0, "Message decoded when blank message passed");

        isoTPProtocol.decode(ctx, null, out);
        assertTrue(out.size() == 0, "Message decoded when blank message passed");
    }
}
