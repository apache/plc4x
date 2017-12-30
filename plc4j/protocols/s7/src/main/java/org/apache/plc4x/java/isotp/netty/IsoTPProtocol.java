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

    private final static Logger logger = LoggerFactory.getLogger(IsoTPProtocol.class);

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
        if(evt instanceof S7ConnectionEvent && ((S7ConnectionEvent) evt).getState() == S7ConnectionState.INITIAL) {
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
    protected void encode(ChannelHandlerContext ctx, Tpdu in, List<Object> out) throws Exception {
        logger.debug("ISO Transport Protocol Message sent");

        ByteBuf buf = Unpooled.buffer();

        // Header length indicator field (The length byte doesn't count)
        buf.writeByte((byte) (getHeaderLength(in) - 1));
        // TPDU Code (First 4 bits), Initial Credit Allocation (Second 4 bits)
        buf.writeByte(in.getTpduCode().getCode());
        // The fixed header of a TCP TP Packet depends highly on the selected type.
        switch (in.getTpduCode()) {
            case CONNECTION_REQUEST:
            case CONNECTION_CONFIRM: {
                ConnectionTpdu connectionTpdu = (ConnectionTpdu) in;
                buf.writeShort(connectionTpdu.getDestinationReference());
                buf.writeShort(connectionTpdu.getSourceReference());
                buf.writeByte(connectionTpdu.getProtocolClass().getCode());
                outputParameters(buf, in.getParameters());
                break;
            }
            case DATA: {
                DataTpdu dataTpdu = (DataTpdu) in;
                // EOT (Bit 8 = 1) / TPDU (All other bits 0)
                buf.writeByte((byte) (dataTpdu.getTpduRef() | (dataTpdu.isEot() ? 0x80 : 0x00)));
                // Note: A Data TPDU in Class 0 doesn't have parameters
                break;
            }
            case DISCONNECT_REQUEST:
            case DISCONNECT_CONFIRM: {
                DisconnectTpdu disconnectTpdu = (DisconnectTpdu) in;
                buf.writeShort(disconnectTpdu.getDestinationReference());
                buf.writeShort(disconnectTpdu.getSourceReference());
                if (disconnectTpdu instanceof DisconnectRequestTpdu) {
                    DisconnectRequestTpdu disconnectRequestTpdu = (DisconnectRequestTpdu) disconnectTpdu;
                    buf.writeByte(disconnectRequestTpdu.getDisconnectReason().getCode());
                }
                outputParameters(buf, in.getParameters());
                break;
            }
            case TPDU_ERROR: {
                ErrorTpdu errorTpdu = (ErrorTpdu) in;
                buf.writeShort(errorTpdu.getDestinationReference());
                buf.writeByte(errorTpdu.getRejectCause().getCode());
                outputParameters(buf, in.getParameters());
                break;
            }
            default: {
                logger.error("TDPU Value {} not implemented yet", new Object[]{in.getTpduCode().name()});
                return;
            }
        }
        // Add the user-data itself.
        buf.writeBytes(in.getUserData());

        out.add(new IsoOnTcpMessage(buf));
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, IsoOnTcpMessage in, List<Object> out) throws Exception {
        if(logger.isTraceEnabled()) {
            logger.trace("Got Data: {}", ByteBufUtil.hexDump(in.getUserData()));
        }
        logger.debug("ISO TP Message received");

        ByteBuf userData = in.getUserData();
        int packetStart = userData.readerIndex();
        byte headerLength = userData.readByte();
        int headerEnd = packetStart + headerLength;
        TpduCode tpduCode = TpduCode.valueOf(userData.readByte());
        // Read fixed header part.
        Tpdu tpdu = null;
        List<Parameter> parameters = new LinkedList<>();
        switch (tpduCode) {
            case CONNECTION_REQUEST:
            case CONNECTION_CONFIRM: {
                short destinationReference = userData.readShort();
                short sourceReference = userData.readShort();
                ProtocolClass protocolClass = ProtocolClass.valueOf(userData.readByte());
                switch (tpduCode) {
                    case CONNECTION_REQUEST:
                        tpdu = new ConnectionRequestTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);
                        break;
                    case CONNECTION_CONFIRM:
                        tpdu = new ConnectionConfirmTpdu(destinationReference, sourceReference, protocolClass, parameters, userData);
                        ctx.channel().pipeline().fireUserEventTriggered(
                            new S7ConnectionEvent(S7ConnectionState.ISO_TP_CONNECTION_RESPONSE_RECEIVED));
                        break;
                }
                break;
            }
            case DATA: {
                byte tmp = userData.readByte();
                boolean eot = (tmp & 0x80) == 0x80;
                byte tpduRef = (byte) (tmp & 0x7F);
                tpdu = new DataTpdu(eot, tpduRef, parameters, userData);
                break;
            }
            case DISCONNECT_REQUEST:
            case DISCONNECT_CONFIRM: {
                short destinationReference = userData.readShort();
                short sourceReference = userData.readShort();
                if (tpduCode == TpduCode.DISCONNECT_REQUEST) {
                    DisconnectReason disconnectReason = DisconnectReason.valueOf(userData.readByte());
                    tpdu = new DisconnectRequestTpdu(
                        destinationReference, sourceReference, disconnectReason, parameters, userData);
                } else {
                    tpdu = new DisconnectConfirmTpdu(
                        destinationReference, sourceReference, parameters, userData);
                }
                break;
            }
            case TPDU_ERROR: {
                short destinationReference = userData.readShort();
                RejectCause rejectCause = RejectCause.valueOf(userData.readByte());
                tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, userData);
                break;
            }
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
                calledTsapParameter = tpdu.getParameter(CalledTsapParameter.class);
                tpduSizeParameter = tpdu.getParameter(TpduSizeParameter.class);
            }
            out.add(new IsoTPMessage(tpdu, userData));
        }
    }

    private void outputParameters(ByteBuf out, List<Parameter> parameters) {
        if (parameters != null) {
            for (Parameter parameter : parameters) {
                out.writeByte(parameter.getType().getCode());
                out.writeByte((byte) (getParameterLength(parameter) - 2));
                switch (parameter.getType()) {
                    case CALLED_TSAP:
                    case CALLING_TSAP: {
                        TsapParameter tsap = (TsapParameter) parameter;
                        out.writeByte(tsap.getDeviceGroup().getCode());
                        out.writeByte((byte)
                            ((tsap.getRackNumber() << 4) | (tsap.getSlotNumber())));
                        break;
                    }
                    case CHECKSUM: {
                        ChecksumParameter checksum = (ChecksumParameter) parameter;
                        out.writeByte(checksum.getChecksum());
                        break;
                    }
                    case DISCONNECT_ADDITIONAL_INFORMATION: {
                        DisconnectAdditionalInformationParameter disconnectAdditionalInformation = (DisconnectAdditionalInformationParameter) parameter;
                        out.writeBytes(disconnectAdditionalInformation.getData());
                        break;
                    }
                    case TPDU_SIZE: {
                        TpduSizeParameter tpduSize = (TpduSizeParameter) parameter;
                        out.writeByte(tpduSize.getTpduSize().getCode());
                        break;
                    }
                    default: {
                        logger.error("TDPU tarameter type {} not implemented yet",
                            new Object[]{parameter.getType().name()});
                        return;
                    }
                }
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
        if (tpdu != null) {
            short headerLength = 0;
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
            }
            return (short) (headerLength + getParametersLength(tpdu.getParameters()));
        }
        return 0;
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
        if (parameter != null) {
            switch (parameter.getType()) {
                case CALLED_TSAP:
                case CALLING_TSAP:
                    return 4;
                case CHECKSUM:
                    return 3;
                case DISCONNECT_ADDITIONAL_INFORMATION: {
                    DisconnectAdditionalInformationParameter disconnectAdditionalInformationParameter =
                        (DisconnectAdditionalInformationParameter) parameter;
                    return (short) (2 + ((disconnectAdditionalInformationParameter.getData() != null) ?
                        disconnectAdditionalInformationParameter.getData().length : 0));
                }
                case TPDU_SIZE:
                    return 3;
            }
        }
        return 0;
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
            case CALLED_TSAP: {
                DeviceGroup deviceGroup = DeviceGroup.valueOf(out.readByte());
                byte tmp = out.readByte();
                byte rackId = (byte) ((tmp & 0xF0) >> 4);
                byte slotId = (byte) (tmp & 0x0F);
                switch (parameterCode) {
                    case CALLING_TSAP:
                        return new CallingTsapParameter(deviceGroup, rackId, slotId);
                    case CALLED_TSAP:
                        return new CalledTsapParameter(deviceGroup, rackId, slotId);
                    default:
                        logger.error("Parameter not implemented yet " + parameterCode.name());
                        return null;
                }
            }
            case CHECKSUM: {
                byte checksum = out.readByte();
                return new ChecksumParameter(checksum);
            }
            case DISCONNECT_ADDITIONAL_INFORMATION: {
                byte[] data = new byte[length];
                out.readBytes(data);
                return new DisconnectAdditionalInformationParameter(data);
            }
            case TPDU_SIZE: {
                TpduSize tpduSize = TpduSize.valueOf(out.readByte());
                return new TpduSizeParameter(tpduSize);
            }
            default: {
                logger.error("Parameter not implemented yet " + parameterCode.name());
                return null;
            }
        }
    }


}
