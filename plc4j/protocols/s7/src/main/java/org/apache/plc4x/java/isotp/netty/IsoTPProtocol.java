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
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.isoontcp.netty.model.IsoOnTcpMessage;
import org.apache.plc4x.java.isotp.netty.model.IsoTPMessage;
import org.apache.plc4x.java.isotp.netty.model.params.*;
import org.apache.plc4x.java.isotp.netty.model.tpdus.*;
import org.apache.plc4x.java.isotp.netty.model.types.*;
import org.apache.plc4x.java.netty.events.S7ConnectionEvent;
import org.apache.plc4x.java.netty.events.S7ConnectionState;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;

public class IsoTPProtocol extends MessageToMessageCodec<IsoOnTcpMessage, Tpdu> {

    private static final Logger logger = LoggerFactory.getLogger(IsoTPProtocol.class);

    private final byte rackNo;
    private final byte slotNo;
    private final TpduSize tpduSize;

    private CalledTsapParameter calledTsapParameter;
    private TpduSizeParameter tpduSizeParameter;

    public IsoTPProtocol(byte rackNo, byte slotNo, TpduSize tpduSize) {
        this.rackNo = rackNo;
        this.slotNo = slotNo;
        this.tpduSize = tpduSize;
    }

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        // If the connection has just been established, start setting up the connection
        // by sending a connection request to the plc.
        if (evt instanceof S7ConnectionEvent && ((S7ConnectionEvent) evt).getState() == S7ConnectionState.INITIAL) {
            logger.debug("ISO Transport Protocol Sending Connection Request");
            // Open the session on ISO Transport Protocol first.
            ConnectionRequestTpdu connectionRequest = new ConnectionRequestTpdu(
                (short) 0x0000, (short) 0x000F, ProtocolClass.CLASS_0,
                Arrays.asList(
                    new CalledTsapParameter(DeviceGroup.PG_OR_PC, (byte) 0, (byte) 0),
                    new CallingTsapParameter(DeviceGroup.OTHERS, rackNo, slotNo),
                    new TpduSizeParameter(tpduSize)),
                Unpooled.buffer());
            ctx.channel().writeAndFlush(connectionRequest);
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, Tpdu in, List<Object> out) {
        logger.debug("ISO Transport Protocol Message sent");

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
                encodeConnecton(in, buf);
                break;
            case DATA:
                encodeData((DataTpdu) in, buf);
                break;
            case DISCONNECT_REQUEST:
            case DISCONNECT_CONFIRM:
                encodeDisconnect(in, buf);
                break;
            case TPDU_ERROR:
                enocdeError(in, buf);
                break;
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("TDPU Value {} not implemented yet", in.getTpduCode().name());
                }
                return;
        }
        // Add the user-data itself.
        buf.writeBytes(in.getUserData());

        out.add(new IsoOnTcpMessage(buf));
    }

    private void enocdeError(Tpdu in, ByteBuf buf) {
        ErrorTpdu errorTpdu = (ErrorTpdu) in;
        buf.writeShort(errorTpdu.getDestinationReference());
        buf.writeByte(errorTpdu.getRejectCause().getCode());
        outputParameters(buf, in.getParameters());
    }

    private void encodeDisconnect(Tpdu in, ByteBuf buf) {
        DisconnectTpdu disconnectTpdu = (DisconnectTpdu) in;
        buf.writeShort(disconnectTpdu.getDestinationReference());
        buf.writeShort(disconnectTpdu.getSourceReference());
        if (disconnectTpdu instanceof DisconnectRequestTpdu) {
            DisconnectRequestTpdu disconnectRequestTpdu = (DisconnectRequestTpdu) disconnectTpdu;
            buf.writeByte(disconnectRequestTpdu.getDisconnectReason().getCode());
        }
        outputParameters(buf, in.getParameters());
    }

    private void encodeData(DataTpdu in, ByteBuf buf) {
        // EOT (Bit 8 = 1) / TPDU (All other bits 0)
        buf.writeByte((byte) (in.getTpduRef() | (in.isEot() ? 0x80 : 0x00)));
        // Note: A Data TPDU in Class 0 doesn't have parameters
    }

    private void encodeConnecton(Tpdu in, ByteBuf buf) {
        ConnectionTpdu connectionTpdu = (ConnectionTpdu) in;
        buf.writeShort(connectionTpdu.getDestinationReference());
        buf.writeShort(connectionTpdu.getSourceReference());
        buf.writeByte(connectionTpdu.getProtocolClass().getCode());
        outputParameters(buf, in.getParameters());
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, IsoOnTcpMessage in, List<Object> out) {
        if (logger.isTraceEnabled()) {
            logger.trace("Got Data: {}", ByteBufUtil.hexDump(in.getUserData()));
        }
        logger.debug("ISO TP Message received");

        if (in == null) {
            return;
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
                tpdu = decodeConnection(ctx, userData, tpduCode, parameters);
                break;
            case DATA:
                tpdu = decodeData(userData, parameters);
                break;
            case DISCONNECT_REQUEST:
            case DISCONNECT_CONFIRM:
                tpdu = decodeDisconnect(userData, tpduCode, parameters);
                break;
            case TPDU_ERROR:
                tpdu = decodeError(userData, parameters);
                break;
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("Tpdu Code {} not implemented", tpduCode.name());
                }
                break;
        }

        // Read variable header parameters
        while (userData.readerIndex() < headerEnd) {
            Parameter parameter = parseParameter(userData);
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
                // TODO: check if null is a valid value (fails test: org.apache.plc4x.java.isotp.netty.IsoTPProtocolTest.decodeConnectionConfirm)
                calledTsapParameter = tpdu.getParameter(CalledTsapParameter.class).orElse(null);
                // TODO: check if null is a valid value (fails test: org.apache.plc4x.java.isotp.netty.IsoTPProtocolTest.decodeConnectionConfirm)
                tpduSizeParameter = tpdu.getParameter(TpduSizeParameter.class).orElse(null);
            }
            out.add(new IsoTPMessage(tpdu, userData));
        }
    }

    private Tpdu decodeError(ByteBuf userData, List<Parameter> parameters) {
        Tpdu tpdu;
        short destinationReference = userData.readShort();
        RejectCause rejectCause = RejectCause.valueOf(userData.readByte());
        tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, userData);
        return tpdu;
    }

    private Tpdu decodeDisconnect(ByteBuf userData, TpduCode tpduCode, List<Parameter> parameters) {
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

    private Tpdu decodeData(ByteBuf userData, List<Parameter> parameters) {
        Tpdu tpdu;
        byte tmp = userData.readByte();
        boolean eot = (tmp & 0x80) == 0x80;
        byte tpduRef = (byte) (tmp & 0x7F);
        tpdu = new DataTpdu(eot, tpduRef, parameters, userData);
        return tpdu;
    }

    private Tpdu decodeConnection(ChannelHandlerContext ctx, ByteBuf userData, TpduCode tpduCode, List<Parameter> parameters) {
        Tpdu tpdu;
        short destinationReference = userData.readShort();
        short sourceReference = userData.readShort();
        ProtocolClass protocolClass = ProtocolClass.valueOf(userData.readByte());

        if (tpduCode == TpduCode.CONNECTION_REQUEST) {
            tpdu = new ConnectionRequestTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);

        } else { // TpduCode.CONNECTION_CONFIRM
            tpdu = new ConnectionConfirmTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);
            ctx.channel().pipeline().fireUserEventTriggered(
                new S7ConnectionEvent(S7ConnectionState.ISO_TP_CONNECTION_RESPONSE_RECEIVED));

        }
        return tpdu;
    }

    private void outputParameters(ByteBuf out, List<Parameter> parameters) {
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
                    out.writeByte(tsap.getDeviceGroup().getCode());
                    out.writeByte((byte) ((tsap.getRackNumber() << 4) | (tsap.getSlotNumber() & 0x0F)));
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

    private Parameter parseParameter(ByteBuf out) {
        ParameterCode parameterCode = ParameterCode.valueOf(out.readByte());
        if (parameterCode == null) {
            logger.error("Could not find parameter code");
            return null;
        }
        byte length = out.readByte();
        switch (parameterCode) {
            case CALLING_TSAP:
            case CALLED_TSAP:
                return parseCallParameter(out, parameterCode);
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

    private Parameter parseCallParameter(ByteBuf out, ParameterCode parameterCode) {
        DeviceGroup deviceGroup = DeviceGroup.valueOf(out.readByte());
        byte rackAndSlot = out.readByte();
        byte rackId = (byte) ((rackAndSlot & 0xF0) >> 4);
        byte slotId = (byte) (rackAndSlot & 0x0F);
        switch (parameterCode) {
            case CALLING_TSAP:
                return new CallingTsapParameter(deviceGroup, rackId, slotId);
            case CALLED_TSAP:
                return new CalledTsapParameter(deviceGroup, rackId, slotId);
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("Parameter not implemented yet {}", parameterCode.name());
                }
                return null;
        }
    }


}
