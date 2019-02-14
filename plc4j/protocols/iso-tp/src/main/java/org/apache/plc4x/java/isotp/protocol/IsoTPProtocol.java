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
package org.apache.plc4x.java.isotp.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.api.exceptions.PlcProtocolPayloadTooBigException;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.isoontcp.protocol.IsoOnTcpProtocol;
import org.apache.plc4x.java.isoontcp.protocol.model.IsoOnTcpMessage;
import org.apache.plc4x.java.isotp.protocol.events.IsoTPConnectedEvent;
import org.apache.plc4x.java.isotp.protocol.model.IsoTPMessage;
import org.apache.plc4x.java.isotp.protocol.model.params.*;
import org.apache.plc4x.java.isotp.protocol.model.tpdus.*;
import org.apache.plc4x.java.isotp.protocol.model.types.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;

public class IsoTPProtocol extends PlcMessageToMessageCodec<IsoOnTcpMessage, Tpdu> {

    private static final Logger logger = LoggerFactory.getLogger(IsoTPProtocol.class);

    private short callingTsapId;
    private short calledTsapId;
    private TpduSize tpduSize;

    public IsoTPProtocol(short callingTsapId, short calledTsapId, TpduSize tpduSize) {
        this.callingTsapId = callingTsapId;
        this.calledTsapId = calledTsapId;
        this.tpduSize = tpduSize;
    }

    /**
     * If the IsoTP protocol is used on top of the ISO on TCP protocol, then as soon as the pipeline receives the
     * request to connect, an IsoTP connection request TPDU must be sent in order to initialize the connection.
     *
     * @param ctx the current protocol layers context
     * @param evt the event
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        ChannelHandler prevHandler = getPrevChannelHandler(ctx);

        // If the connection has just been established, start setting up the connection
        // by sending a connection request to the plc.
        if ((prevHandler instanceof IsoOnTcpProtocol) && (evt instanceof ConnectEvent)) {
            logger.debug("ISO Transport Protocol Sending Connection Request");
            // Open the session on ISO Transport Protocol first.
            ConnectionRequestTpdu connectionRequest = new ConnectionRequestTpdu(
                (short) 0x0000, (short) 0x000F, ProtocolClass.CLASS_0,
                Arrays.asList(
                    new CalledTsapParameter(calledTsapId),
                    new CallingTsapParameter(callingTsapId),
                    new TpduSizeParameter(tpduSize)),
                Unpooled.buffer());
            ctx.channel().writeAndFlush(connectionRequest);
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void encode(ChannelHandlerContext ctx, Tpdu in, List<Object> out) {
        logger.trace("ISO Transport Protocol Message sent");

        if (in == null) {
            return;
        }

        ByteBuf buf = Unpooled.buffer();

        // Header length indicator field (The length byte doesn't count)
        buf.writeByte((byte) (getHeaderLength(in) - 1));
        // TPDU Code (First 4 bits), Initial Credit Allocation (Second 4 bits)
        buf.writeByte(in.getTpduCode().getCode());
        // The fixed header of a TCP TP Packet depends highly on the selected type.
        switch (in.getTpduCode()) {
            case CONNECTION_REQUEST:
            case CONNECTION_CONFIRM:
                encodeConnectionTpdu(in, buf);
                break;
            case DATA:
                encodeDataTpdu((DataTpdu) in, buf);
                break;
            case DISCONNECT_REQUEST:
            case DISCONNECT_CONFIRM:
                encodeDisconnectTpdu(in, buf);
                break;
            case TPDU_ERROR:
                encodeErrorTpdu(in, buf);
                break;
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("TDPU Value {} not implemented yet", in.getTpduCode().name());
                }
                return;
        }
        // Add the user-data itself.
        buf.writeBytes(in.getUserData());

        // Check if the message doesn't exceed the negotiated maximum size.
        if(buf.writerIndex() > tpduSize.getValue()) {
            ctx.fireExceptionCaught(new PlcProtocolPayloadTooBigException(
                "iso-tp", tpduSize.getValue(), buf.writerIndex(), in));
        } else {
            out.add(new IsoOnTcpMessage(buf));
        }
    }

    private void encodeErrorTpdu(Tpdu in, ByteBuf buf) {
        ErrorTpdu errorTpdu = (ErrorTpdu) in;
        buf.writeShort(errorTpdu.getDestinationReference());
        buf.writeByte(errorTpdu.getRejectCause().getCode());
        encodeParameters(buf, in.getParameters());
    }

    private void encodeDisconnectTpdu(Tpdu in, ByteBuf buf) {
        DisconnectTpdu disconnectTpdu = (DisconnectTpdu) in;
        buf.writeShort(disconnectTpdu.getDestinationReference());
        buf.writeShort(disconnectTpdu.getSourceReference());
        if (disconnectTpdu instanceof DisconnectRequestTpdu) {
            DisconnectRequestTpdu disconnectRequestTpdu = (DisconnectRequestTpdu) disconnectTpdu;
            buf.writeByte(disconnectRequestTpdu.getDisconnectReason().getCode());
        }
        encodeParameters(buf, in.getParameters());
    }

    private void encodeDataTpdu(DataTpdu in, ByteBuf buf) {
        // EOT (Bit 8 = 1) / TPDU (All other bits 0)
        buf.writeByte((byte) (in.getTpduRef() | (in.isEot() ? 0x80 : 0x00)));
        // Note: A Data TPDU in Class 0 doesn't have parameters
    }

    private void encodeConnectionTpdu(Tpdu in, ByteBuf buf) {
        ConnectionTpdu connectionTpdu = (ConnectionTpdu) in;
        buf.writeShort(connectionTpdu.getDestinationReference());
        buf.writeShort(connectionTpdu.getSourceReference());
        buf.writeByte(connectionTpdu.getProtocolClass().getCode());
        encodeParameters(buf, in.getParameters());
    }

    private void encodeParameters(ByteBuf out, List<Parameter> parameters) {
        if (parameters == null) {
            return;
        }

        for (Parameter parameter : parameters) {
            out.writeByte(parameter.getType().getCode());
            out.writeByte((byte) (getParameterLength(parameter) - 2));
            switch (parameter.getType()) {
                case CALLED_TSAP:
                case CALLING_TSAP:
                    TsapParameter tsap = (TsapParameter) parameter;
                    out.writeShort(tsap.getTsapId());
                    break;
                case CHECKSUM:
                    ChecksumParameter checksum = (ChecksumParameter) parameter;
                    out.writeByte(checksum.getChecksum());
                    break;
                case DISCONNECT_ADDITIONAL_INFORMATION:
                    DisconnectAdditionalInformationParameter disconnectAdditionalInformation = (DisconnectAdditionalInformationParameter) parameter;
                    out.writeBytes(disconnectAdditionalInformation.getData());
                    break;
                case TPDU_SIZE:
                    TpduSizeParameter sizeParameter = (TpduSizeParameter) parameter;
                    out.writeByte(sizeParameter.getTpduSize().getCode());
                    break;
                default:
                    if (logger.isErrorEnabled()) {
                        logger.error("TDPU tarameter type {} not implemented yet", parameter.getType().name());
                    }
                    return;
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void decode(ChannelHandlerContext ctx, IsoOnTcpMessage in, List<Object> out) {
        logger.trace("ISO TP Message received");
        if (in == null) {
            return;
        }

        if (logger.isDebugEnabled()) {
            logger.debug("Got Data: {}", ByteBufUtil.hexDump(in.getUserData()));
        }

        ByteBuf userData = in.getUserData();
        if (userData.writerIndex() < 1) {
            return;
        }

        int packetStart = userData.readerIndex();
        byte headerLength = userData.readByte();
        int headerEnd = packetStart + headerLength;
        TpduCode tpduCode = TpduCode.valueOf(userData.readByte());
        // Read fixed header part.
        Tpdu tpdu = null;
        List<Parameter> parameters = new LinkedList<>();
        switch (tpduCode) {
            case CONNECTION_REQUEST:
            case CONNECTION_CONFIRM:
                tpdu = decodeConnectTpdu(ctx, userData, tpduCode, parameters);
                break;
            case DATA:
                tpdu = decodeDataTpdu(userData, parameters);
                break;
            case DISCONNECT_REQUEST:
            case DISCONNECT_CONFIRM:
                tpdu = decodeDisconnectTpdu(userData, tpduCode, parameters);
                break;
            case TPDU_ERROR:
                tpdu = decodeErrorTpdu(userData, parameters);
                break;
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("Tpdu Code {} not implemented", tpduCode.name());
                }
                break;
        }

        // Read variable header parameters
        while (userData.readerIndex() < headerEnd) {
            Parameter parameter = decodeParameter(userData);
            if (parameter != null) {
                parameters.add(parameter);
            }
        }

        if (tpdu != null) {
            // If we got a ConnectionConfirmTpdu response we are currently
            // in the process of establishing a connection with the PLC, so
            // Save some of the information in the session and tell the next
            // layer to negotiate the connection parameters.
            if (tpdu instanceof ConnectionConfirmTpdu) {
                tpdu.getParameter(CalledTsapParameter.class).ifPresent(
                    calledTsapParameter -> calledTsapId = calledTsapParameter.getTsapId());
                tpdu.getParameter(TpduSizeParameter.class).ifPresent(
                    tpduSizeParameter -> tpduSize = tpduSizeParameter.getTpduSize());
            }
            out.add(new IsoTPMessage(tpdu, userData));
        }
    }

    private Tpdu decodeErrorTpdu(ByteBuf userData, List<Parameter> parameters) {
        Tpdu tpdu;
        short destinationReference = userData.readShort();
        RejectCause rejectCause = RejectCause.valueOf(userData.readByte());
        tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, userData);
        return tpdu;
    }

    private Tpdu decodeDisconnectTpdu(ByteBuf userData, TpduCode tpduCode, List<Parameter> parameters) {
        Tpdu tpdu;
        short destinationReference = userData.readShort();
        short sourceReference = userData.readShort();
        if (tpduCode == TpduCode.DISCONNECT_REQUEST) {
            DisconnectReason disconnectReason = DisconnectReason.valueOf(userData.readByte());
            tpdu = new DisconnectRequestTpdu(
                destinationReference, sourceReference, disconnectReason, parameters, userData);
        } else {  // TpduCode.DISCONNECT_CONFIRM
            tpdu = new DisconnectConfirmTpdu(
                destinationReference, sourceReference, parameters, userData);
        }
        return tpdu;
    }

    private Tpdu decodeDataTpdu(ByteBuf userData, List<Parameter> parameters) {
        Tpdu tpdu;
        byte tmp = userData.readByte();
        // Bit 8 is the EOT indicator (1 = last TPDU)
        boolean eot = (tmp & 0x80) == 0x80;
        // The rest is simply a 7 bit number identifying the current request.
        byte tpduRef = (byte) (tmp & 0x7F);
        tpdu = new DataTpdu(eot, tpduRef, parameters, userData);
        return tpdu;
    }

    private Tpdu decodeConnectTpdu(ChannelHandlerContext ctx, ByteBuf userData, TpduCode tpduCode, List<Parameter> parameters) {
        Tpdu tpdu;
        short destinationReference = userData.readShort();
        short sourceReference = userData.readShort();
        ProtocolClass protocolClass = ProtocolClass.valueOf(userData.readByte());

        if (tpduCode == TpduCode.CONNECTION_REQUEST) {
            tpdu = new ConnectionRequestTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);

        } else { // TpduCode.CONNECTION_CONFIRM
            tpdu = new ConnectionConfirmTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);
            ctx.channel().pipeline().fireUserEventTriggered(new IsoTPConnectedEvent());
        }
        return tpdu;
    }

    private Parameter decodeParameter(ByteBuf out) {
        ParameterCode parameterCode = ParameterCode.valueOf(out.readByte());
        if (parameterCode == null) {
            logger.error("Could not find parameter code");
            return null;
        }
        byte length = out.readByte();
        switch (parameterCode) {
            case CALLING_TSAP:
            case CALLED_TSAP:
                return decodeTsapParameter(out, parameterCode);
            case CHECKSUM:
                byte checksum = out.readByte();
                return new ChecksumParameter(checksum);
            case DISCONNECT_ADDITIONAL_INFORMATION:
                byte[] data = new byte[length];
                out.readBytes(data);
                return new DisconnectAdditionalInformationParameter(data);
            case TPDU_SIZE:
                TpduSize size = TpduSize.valueOf(out.readByte());
                return new TpduSizeParameter(size);
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("Parameter not implemented yet {}", parameterCode.name());
                }
                return null;
        }
    }

    private Parameter decodeTsapParameter(ByteBuf out, ParameterCode parameterCode) {
        short tsapId = out.readShort();
        switch (parameterCode) {
            case CALLING_TSAP:
                return new CallingTsapParameter(tsapId);
            case CALLED_TSAP:
                return new CalledTsapParameter(tsapId);
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("Parameter not implemented yet {}", parameterCode.name());
                }
                return null;
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * Return the length of the entire header in bytes (including the size field itself)
     * This is a sum of the fixed size header defined for the given tpdu type and the
     * lengths of all parameters.
     *
     * @param tpdu Tpdu to get the header length for
     * @return length of the iso tp header
     */
    private short getHeaderLength(Tpdu tpdu) {
        if (tpdu == null) {
            return 0;
        }

        short headerLength;
        switch (tpdu.getTpduCode()) {
            case CONNECTION_REQUEST:
            case CONNECTION_CONFIRM:
                headerLength = 7;
                break;
            case DATA:
                headerLength = 3;
                break;
            case DISCONNECT_REQUEST:
                headerLength = 7;
                break;
            case DISCONNECT_CONFIRM:
                headerLength = 6;
                break;
            case TPDU_ERROR:
                headerLength = 5;
                break;
            default:
                headerLength = 0;
                break;
        }
        return (short) (headerLength + getParametersLength(tpdu.getParameters()));
    }

    private short getParametersLength(List<Parameter> parameters) {
        short length = 0;
        if (parameters != null) {
            for (Parameter parameter : parameters) {
                length += getParameterLength(parameter);
            }
        }
        return length;
    }

    private short getParameterLength(Parameter parameter) {
        if (parameter == null) {
            return 0;
        }
        switch (parameter.getType()) {
            case CALLED_TSAP:
            case CALLING_TSAP:
                return 4;
            case CHECKSUM:
                return 3;
            case DISCONNECT_ADDITIONAL_INFORMATION:
                DisconnectAdditionalInformationParameter disconnectAdditionalInformationParameter =
                    (DisconnectAdditionalInformationParameter) parameter;
                return (short) (2 + ((disconnectAdditionalInformationParameter.getData() != null) ?
                    disconnectAdditionalInformationParameter.getData().length : 0));
            case TPDU_SIZE:
                return 3;
            default:
                return 0;
        }
    }

}
