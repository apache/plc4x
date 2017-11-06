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
package org.apache.plc4x.java.s7.messages.s7;

import org.apache.plc4x.java.exception.PlcException;
import org.apache.plc4x.java.exception.PlcIoException;
import org.apache.plc4x.java.exception.PlcProtocolException;
import org.apache.plc4x.java.s7.messages.*;
import org.apache.plc4x.java.s7.messages.isotp.messages.AbstractConnectionMessage;
import org.apache.plc4x.java.s7.messages.isotp.params.CalledTsapIsoTpParameter;
import org.apache.plc4x.java.s7.messages.isotp.params.IsoTpParameter;
import org.apache.plc4x.java.s7.messages.isotp.params.IsoTpParameterFactory;
import org.apache.plc4x.java.s7.messages.isotp.types.ProtocolClass;
import org.apache.plc4x.java.s7.messages.isotp.types.TpduCode;
import org.apache.plc4x.java.s7.messages.s7.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.messages.s7.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.messages.s7.params.S7Parameter;
import org.apache.plc4x.java.s7.messages.s7.params.S7ParameterFactory;
import org.apache.plc4x.java.s7.messages.s7.payload.ReadVarPayload;
import org.apache.plc4x.java.s7.messages.s7.payload.S7Payload;
import org.apache.plc4x.java.s7.messages.s7.types.MessageType;
import org.apache.plc4x.java.s7.messages.s7.types.DataTransportSize;

import java.io.DataInputStream;
import java.io.IOException;
import java.util.LinkedList;
import java.util.List;

public class S7Parser implements Parser {

    @Override
    public Message parse(DataInputStream dis) throws PlcException {
        try {
            ////////////////////////////////////////////////////
            // RFC 1006 (ISO on TCP)
            if (dis.readByte() != AbstractConnectionMessage.ISO_ON_TCP_MAGIC_NUMBER) {
                throw new PlcProtocolException("Error parsing message. Expecting ISO on TCP magic number");
            }
            dis.readByte(); // Reserved ...
            short isoOnTcpPacketLength = dis.readShort();

            ////////////////////////////////////////////////////
            // RFC 905 (ISO Transport Protocol)
            int isoTpHeaderLength = dis.readByte() + 1;
            TpduCode tpduCode = TpduCode.valueOf(dis.readByte());
            if(tpduCode == null) {
                throw new PlcProtocolException("Error parsing message. Unknown tpdu code.");
            }
            List<IsoTpParameter> isoTpParameters = new LinkedList<>();
            switch (tpduCode) {
                case CONNECTION_REQUEST:
                case CONNECTION_CONFIRM:
                    short destinationReference = dis.readShort();
                    short sourceReference = dis.readShort();
                    ProtocolClass protocolClass = ProtocolClass.valueOf(dis.readByte());
                    if(protocolClass != ProtocolClass.CLASS_0) {
                        throw new PlcProtocolException("Error parsing message. Only class 0 supported.");
                    }
                    int variablePartLength = isoTpHeaderLength - 7;
                    CalledTsapIsoTpParameter calledTsapIsoTpParameter = null;
                    for(int i = 0; i < variablePartLength; ) {
                        IsoTpParameter parameter = IsoTpParameterFactory.parse(dis);
                        if(parameter instanceof CalledTsapIsoTpParameter) {
                            calledTsapIsoTpParameter = (CalledTsapIsoTpParameter) parameter;
                        }
                        isoTpParameters.add(parameter);
                        i += parameter.getLength();
                    }
                    if(calledTsapIsoTpParameter == null) {
                        throw new PlcProtocolException("Error parsing message. Missing called tsap parameter.");
                    }
                    if(tpduCode == TpduCode.CONNECTION_REQUEST) {
                        return new ConnectionRequest(sourceReference,
                            calledTsapIsoTpParameter.getRackNumber(), calledTsapIsoTpParameter.getSlotNumber());
                    } else {
                        return new ConnectionResponse(sourceReference,
                            calledTsapIsoTpParameter.getRackNumber(), calledTsapIsoTpParameter.getSlotNumber());
                    }
                case DATA:
                    byte tpduNr = dis.readByte();
                    // S7 Protocol
                    if(dis.readByte() != S7RequestMessage.S7_PROTOCOL_MAGIC_NUMBER) {
                        throw new PlcProtocolException("Error parsing message. Expecting S7 protocol magic number");
                    }
                    MessageType messageType = MessageType.valueOf(dis.readByte());
                    dis.readShort(); // Reserved ...
                    short tpduReference = dis.readShort();
                    short headerParametersLength = dis.readShort();
                    short userDataLength = dis.readShort();
                    byte errorClass = dis.readByte();
                    byte errorCode = dis.readByte();
                    List<S7Parameter> s7Parameters = new LinkedList<>();
                    List<S7Payload> s7Payloads = new LinkedList<>();
                    for(int i = 0; i < headerParametersLength; ) {
                        S7Parameter parameter = S7ParameterFactory.parse(
                            dis, messageType == MessageType.ACK_DATA);
                        s7Parameters.add(parameter);
                        i += parameter.getLength();
                    }
                    if(userDataLength > 0) {
                        if(dis.readByte() != (byte) 0xFF) {
                            throw new PlcProtocolException("Error reading data");
                        }
                        DataTransportSize dataTransportSize = DataTransportSize.valueOf(dis.readByte());
                        short length = (dataTransportSize.isSizeInBits()) ?
                            (short) Math.ceil(dis.readShort() / 8) : dis.readShort();
                        byte[] data = new byte[length];
                        dis.read(data);
                        ReadVarPayload payload = new ReadVarPayload(null, dataTransportSize, data);
                        s7Payloads.add(payload);
                    }
                    return new S7ResponseMessage(messageType, s7Parameters, s7Payloads, errorClass, errorCode);
                default:
                    throw new PlcProtocolException("Unimplemented tpdu type");
            }
        } catch (IOException e) {
            throw new PlcIoException("Error parsing message", e);
        }
    }

}
