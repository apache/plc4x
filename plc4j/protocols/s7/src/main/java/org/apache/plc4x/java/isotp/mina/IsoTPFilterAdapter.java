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
package org.apache.plc4x.java.isotp.mina;

import org.apache.mina.core.buffer.IoBuffer;
import org.apache.mina.core.filterchain.IoFilterAdapter;
import org.apache.mina.core.session.IoSession;
import org.apache.mina.core.write.WriteRequest;
import org.apache.mina.core.write.WriteRequestWrapper;
import org.apache.plc4x.java.exceptions.PlcConnectionException;
import org.apache.plc4x.java.isoontcp.mina.IsoOnTcpFilterAdapter;
import org.apache.plc4x.java.isotp.mina.model.IsoTPMessage;
import org.apache.plc4x.java.isotp.mina.model.params.*;
import org.apache.plc4x.java.isotp.mina.model.tpdus.*;
import org.apache.plc4x.java.isotp.mina.model.types.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;

public class IsoTPFilterAdapter extends IoFilterAdapter {

    public final static String RACK_NO = "ISO_TP_RACK_NO";
    public final static String SLOT_NO = "ISO_TP_SLOT_NO";

    private final static Logger logger = LoggerFactory.getLogger(IsoOnTcpFilterAdapter.class);

    @Override
    public void sessionOpened(NextFilter nextFilter, IoSession session) throws Exception {
        // Get parameters from the session.
        byte rackNo;
        Object attr = session.getAttribute(RACK_NO);
        if ((attr != null) && (attr instanceof Byte)) {
            rackNo = (byte) attr;
        } else {
            throw new PlcConnectionException("Required session parameter " + RACK_NO + " not specified.");
        }
        byte slotNo;
        attr = session.getAttribute(SLOT_NO);
        if ((attr != null) && (attr instanceof Byte)) {
            slotNo = (byte) attr;
        } else {
            throw new PlcConnectionException("Required session parameter " + SLOT_NO + " not specified.");
        }

        // Open the session on ISO Transport Protocol first.
        ConnectionRequestTpdu connectionRequest = new ConnectionRequestTpdu(
            (short) 0x0000, (short) 0x000F, ProtocolClass.CLASS_0,
            Arrays.asList(
                new CalledTsapParameter(DeviceGroup.PG_OR_PC, (byte) 0, (byte) 0),
                new CallingTsapParameter(DeviceGroup.OTHERS, rackNo, slotNo),
                new TpduSizeParameter(TpduSize.SIZE_1024)),
            IoBuffer.allocate(0).setAutoExpand(true));
        session.write(connectionRequest);
    }

    @Override
    public void filterWrite(NextFilter nextFilter, IoSession session, WriteRequest writeRequest) throws Exception {
        if (writeRequest.getMessage() instanceof Tpdu) {
            logger.debug("ISO TP Message sent: {}", new Object[]{session});

            Tpdu tpdu = (Tpdu) writeRequest.getMessage();

            // Allocate a buffer that's big enough for the entire message.
            final IoBuffer buffer = IoBuffer.allocate(getHeaderLength(tpdu) + tpdu.getUserData().remaining());
            // Header length indicator field (The length byte doesn't count)
            buffer.put((byte) (getHeaderLength(tpdu) - 1));
            // TPDU Code (First 4 bits), Initial Credit Allocation (Second 4 bits)
            buffer.put(tpdu.getTpduCode().getCode());
            // The fixed header of a TCP TP Packet depends highly on the selected type.
            switch (tpdu.getTpduCode()) {
                case CONNECTION_REQUEST:
                case CONNECTION_CONFIRM: {
                    ConnectionTpdu connectionTpdu = (ConnectionTpdu) tpdu;
                    buffer.putShort(connectionTpdu.getDestinationReference());
                    buffer.putShort(connectionTpdu.getSourceReference());
                    buffer.put(connectionTpdu.getProtocolClass().getCode());
                    outputParameters(buffer, tpdu.getParameters());
                    break;
                }
                case DATA: {
                    DataTpdu dataTpdu = (DataTpdu) tpdu;
                    // EOT (Bit 8 = 1) / TPDU (All other bits 0)
                    buffer.put((byte) (dataTpdu.getTpduRef() | (dataTpdu.isEot() ? 0x80 : 0x00)));
                    // Note: A Data TPDU in Class 0 doesn't have parameters
                    break;
                }
                case DISCONNECT_REQUEST:
                case DISCONNECT_CONFIRM: {
                    DisconnectTpdu disconnectTpdu = (DisconnectTpdu) tpdu;
                    buffer.putShort(disconnectTpdu.getDestinationReference());
                    buffer.putShort(disconnectTpdu.getSourceReference());
                    if (disconnectTpdu instanceof DisconnectRequestTpdu) {
                        DisconnectRequestTpdu disconnectRequestTpdu = (DisconnectRequestTpdu) disconnectTpdu;
                        buffer.put(disconnectRequestTpdu.getDisconnectReason().getCode());
                    }
                    outputParameters(buffer, tpdu.getParameters());
                    break;
                }
                case TPDU_ERROR: {
                    ErrorTpdu errorTpdu = (ErrorTpdu) tpdu;
                    buffer.putShort(errorTpdu.getDestinationReference());
                    buffer.put(errorTpdu.getRejectCause().getCode());
                    outputParameters(buffer, tpdu.getParameters());
                }
                default: {
                    logger.error("TDPU Value {} not implemented yet", new Object[]{tpdu.getTpduCode().name()});
                    return;
                }
            }
            // Add the user-data itself.
            buffer.put(tpdu.getUserData());

            // Prepare the buffer for sending.
            buffer.flip();

            // The lower protocol filters only care about sending data, so we replace
            // the payload with the data output by this level.
            writeRequest = new WriteRequestWrapper(writeRequest) {
                @Override
                public Object getMessage() {
                    return buffer;
                }
            };
        }
        nextFilter.filterWrite(session, writeRequest);
    }

    @Override
    public void messageReceived(NextFilter nextFilter, IoSession session, Object message) throws Exception {
        if (message instanceof IoBuffer) {
            logger.debug("ISO TP Message received: {}", new Object[]{session});

            final IoBuffer buffer = (IoBuffer) message;
            int packetStart = buffer.position();
            byte headerLength = buffer.get();
            int headerEnd = packetStart + headerLength;
            TpduCode tpduCode = TpduCode.valueOf(buffer.get());

            // Read fixed header part.
            Tpdu tpdu = null;
            List<Parameter> parameters = new LinkedList<>();
            switch (tpduCode) {
                case CONNECTION_REQUEST:
                case CONNECTION_CONFIRM: {
                    short destinationReference = buffer.getShort();
                    short sourceReference = buffer.getShort();
                    ProtocolClass protocolClass = ProtocolClass.valueOf(buffer.get());
                    switch (tpduCode) {
                        case CONNECTION_REQUEST:
                            tpdu = new ConnectionRequestTpdu(destinationReference, sourceReference, protocolClass, parameters, buffer);
                            break;
                        case CONNECTION_CONFIRM:
                            tpdu = new ConnectionConfirmTpdu(destinationReference, sourceReference, protocolClass, parameters, buffer);

                            break;
                    }
                    break;
                }
                case DATA: {
                    byte tmp = buffer.get();
                    boolean eot = (tmp & 0x80) == 0x80;
                    byte tpduRef = (byte) (tmp & 0x7F);
                    tpdu = new DataTpdu(eot, tpduRef, parameters, buffer);
                    break;
                }
                case DISCONNECT_REQUEST:
                case DISCONNECT_CONFIRM: {
                    short destinationReference = buffer.getShort();
                    short sourceReference = buffer.getShort();
                    if (tpduCode == TpduCode.DISCONNECT_REQUEST) {
                        DisconnectReason disconnectReason = DisconnectReason.valueOf(buffer.get());
                        tpdu = new DisconnectRequestTpdu(destinationReference, sourceReference, disconnectReason, parameters, buffer);
                    } else {
                        tpdu = new DisconnectConfirmTpdu(destinationReference, sourceReference, parameters, buffer);
                    }
                    break;
                }
                case TPDU_ERROR: {
                    short destinationReference = buffer.getShort();
                    RejectCause rejectCause = RejectCause.valueOf(buffer.get());
                    tpdu = new ErrorTpdu(destinationReference, rejectCause, parameters, buffer);
                    break;
                }
            }

            // Read variable header parameters
            while (buffer.position() < headerEnd) {
                Parameter parameter = parseParameter(buffer);
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
                    CalledTsapParameter calledTsapParameter = tpdu.getParameter(CalledTsapParameter.class);
                    if (calledTsapParameter != null) {
                        session.setAttribute("ISO_TP_CALLED_TSAP", calledTsapParameter);
                    }
                    TpduSizeParameter tpduSizeParameter = tpdu.getParameter(TpduSizeParameter.class);
                    if (tpduSizeParameter != null) {
                        session.setAttribute("ISO_TP_TPDU_SIZE", tpduSizeParameter.getTpduSize().name());
                    }
                    nextFilter.sessionOpened(session);
                }
                message = new IsoTPMessage(tpdu, buffer);
            }
        }

        // Let the higher level protocol filters parse their part.
        super.messageReceived(nextFilter, session, message);
    }

    private void outputParameters(IoBuffer buffer, List<Parameter> parameters) {
        if (parameters != null) {
            for (Parameter parameter : parameters) {
                buffer.put(parameter.getType().getCode());
                buffer.put((byte) (getParameterLength(parameter) - 2));
                switch (parameter.getType()) {
                    case CALLED_TSAP:
                    case CALLING_TSAP: {
                        TsapParameter tsap = (TsapParameter) parameter;
                        buffer.put(tsap.getDeviceGroup().getCode());
                        buffer.put((byte)
                            ((tsap.getRackNumber() << 4) | (tsap.getSlotNumber())));
                        break;
                    }
                    case CHECKSUM: {
                        ChecksumParameter checksum = (ChecksumParameter) parameter;
                        buffer.put(checksum.getChecksum());
                        break;
                    }
                    case DISCONNECT_ADDITIONAL_INFORMATION: {
                        DisconnectAdditionalInformationParameter disconnectAdditionalInformation = (DisconnectAdditionalInformationParameter) parameter;
                        buffer.put(disconnectAdditionalInformation.getData());
                        break;
                    }
                    case TPDU_SIZE: {
                        TpduSizeParameter tpduSize = (TpduSizeParameter) parameter;
                        buffer.put(tpduSize.getTpduSize().getCode());
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

    private Parameter parseParameter(IoBuffer buffer) {
        ParameterCode parameterCode = ParameterCode.valueOf(buffer.get());
        if (parameterCode == null) {
            logger.error("Could not find parameter code");
            return null;
        }
        byte length = buffer.get();
        switch (parameterCode) {
            case CALLING_TSAP:
            case CALLED_TSAP: {
                DeviceGroup deviceGroup = DeviceGroup.valueOf(buffer.get());
                byte tmp = buffer.get();
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
                byte checksum = buffer.get();
                return new ChecksumParameter(checksum);
            }
            case DISCONNECT_ADDITIONAL_INFORMATION: {
                byte[] data = new byte[length];
                buffer.get(data);
                return new DisconnectAdditionalInformationParameter(data);
            }
            case TPDU_SIZE: {
                TpduSize tpduSize = TpduSize.valueOf(buffer.get());
                return new TpduSizeParameter(tpduSize);
            }
            default: {
                logger.error("Parameter not implemented yet " + parameterCode.name());
                return null;
            }
        }
    }

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
                    headerLength = 6;
                    break;
                case DISCONNECT_CONFIRM:
                    headerLength = 7;
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

}
