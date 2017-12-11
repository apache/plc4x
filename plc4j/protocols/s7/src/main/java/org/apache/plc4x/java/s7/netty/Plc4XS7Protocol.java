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
package org.apache.plc4x.java.s7.netty;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7AnyVarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4XS7Protocol extends MessageToMessageCodec<S7Message, PlcRequestContainer> {

    private final static Logger logger = LoggerFactory.getLogger(S7PlcConnection.class);

    private static final AtomicInteger tpduGenerator = new AtomicInteger(1);

    private Map<Short, PlcRequestContainer> requests;

    public Plc4XS7Protocol() {
        this.requests = new HashMap<>();
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
        if (msg.getRequest() instanceof PlcReadRequest) {
            PlcReadRequest readRequest = (PlcReadRequest) msg.getRequest();

            // Try to get the correct S7 transport size for the given data type.
            // (Map PLC4X data type to S7 data type)
            TransportSize transportSize = getTransportSize(readRequest.getDatatype());
            if (transportSize == null) {
                throw new PlcException("Unknown transport size for datatype " + readRequest.getDatatype());
            }

            // Depending on the address type, generate the corresponding type of request item.
            VarParameter varParameter = getVarParameter(
                ParameterType.READ_VAR, readRequest.getAddress(), transportSize);

            // Assemble the request.
            S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.JOB,
                (short) tpduGenerator.getAndIncrement(), Collections.singletonList(varParameter),
                Collections.emptyList());

            requests.put(s7ReadRequest.getTpduReference(), msg);

            out.add(s7ReadRequest);
        } else if(msg.getRequest() instanceof PlcWriteRequest) {
            PlcWriteRequest writeRequest = (PlcWriteRequest) msg.getRequest();

            // Try to get the correct S7 transport size for the given data type.
            // (Map PLC4X data type to S7 data type)
            TransportSize transportSize = getTransportSize(writeRequest.getDatatype());
            if (transportSize == null) {
                throw new PlcException("Unknown transport size for datatype " + writeRequest.getDatatype());
            }

            // Depending on the address type, generate the corresponding type of request item.
            VarParameter varParameter = getVarParameter(
                ParameterType.WRITE_VAR, writeRequest.getAddress(), transportSize);

            // Assemble the request.
            DataTransportSize dataTransportSize = getDataTransportSize(writeRequest.getDatatype());
            S7RequestMessage s7WriteRequest = new S7RequestMessage(MessageType.JOB,
                (short) tpduGenerator.getAndIncrement(), Collections.singletonList(varParameter),
                Collections.singletonList(getVarPayload(ParameterType.WRITE_VAR, varParameter.getItems().get(0),
                    dataTransportSize, toS7Data(writeRequest.getValue()))));

            requests.put(s7WriteRequest.getTpduReference(), msg);

            out.add(s7WriteRequest);
        }
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, S7Message msg, List<Object> out) throws Exception {
        if(msg instanceof S7ResponseMessage) {
            S7ResponseMessage responseMessage = (S7ResponseMessage) msg;
            short tpduReference = responseMessage.getTpduReference();
            if (requests.containsKey(tpduReference)) {
                PlcRequestContainer requestContainer = requests.remove(tpduReference);
                PlcResponse response = null;
                if (requestContainer.getRequest() instanceof PlcReadRequest) {
                    // TODO: Link the request and response objects
                    PlcReadRequest plcReadRequest = (PlcReadRequest) requestContainer.getRequest();
                    S7AnyVarPayload payload = responseMessage.getPayload(S7AnyVarPayload.class);
                    if(payload.getDataTransportErrorCode().getCode() == 0xA) {
                        response = plcReadRequest.createResponse(null);
                    } else {
                        byte[] data = payload.getData();
                        Object value = fromS7Data(plcReadRequest.getDatatype(), data);
                        response = plcReadRequest.createResponse(value);
                    }
                } else if (requestContainer.getRequest() instanceof PlcWriteRequest) {
                    // TODO: Link the request and response objects
                    PlcWriteRequest plcWriteRequest = (PlcWriteRequest) requestContainer.getRequest();
                    response = plcWriteRequest.createResponse();
                }
                if (response != null) {
                    requestContainer.getResponseFuture().complete(response);
                }
            }
        }
    }

    private TransportSize getTransportSize(Class<?> datatype) {
        if (datatype == Boolean.class) {
            return TransportSize.BIT;
        } else if (datatype == Byte.class) {
            return TransportSize.BYTE;
        } else if (datatype == Short.class) {
            return TransportSize.WORD;
        } else if (datatype == Calendar.class) {
            return TransportSize.DATE_AND_TIME;
        } else if (datatype == Float.class) {
            return TransportSize.REAL;
        } else if (datatype == Integer.class) {
            return TransportSize.DWORD;
        } else if (datatype == String.class) {
            return TransportSize.CHAR;
        }
        return null;
    }

    private DataTransportSize getDataTransportSize(Class<?> datatype) {
        if (datatype == Boolean.class) {
            return DataTransportSize.BIT;
        } else if (datatype == Byte.class) {
            return DataTransportSize.BYTE_WORD_DWORD;
        } else if (datatype == Short.class) {
            return DataTransportSize.BYTE_WORD_DWORD;
        } else if (datatype == Calendar.class) {
            // TODO: Decide what to do here ...
            return null;
        } else if (datatype == Float.class) {
            return DataTransportSize.REAL;
        } else if (datatype == Integer.class) {
            return DataTransportSize.BYTE_WORD_DWORD;
        } else if (datatype == String.class) {
            return DataTransportSize.OCTET_STRING;
        }
        return null;
    }

    private Object fromS7Data(Class<?> datatype, byte[] s7Data) {
        if(s7Data.length == 0) {
            return null;
        }
        if (datatype == Boolean.class) {
            return (s7Data[0] & 0x01) == 0x01;
        } else if (datatype == Byte.class) {
            return s7Data[0];
        } else if (datatype == Short.class) {
            return Short.valueOf((short) (((s7Data[0] & 0xff) << 8) | (s7Data[1] & 0xff)));
        } else if (datatype == Integer.class) {
            return ByteBuffer.wrap(s7Data).getInt();
        } else if (datatype == Float.class) {
            // Description of the Real number format:
            // https://www.sps-lehrgang.de/zahlenformate-step7/#c144
            // https://de.wikipedia.org/wiki/IEEE_754
            int tmp = 0;
            tmp |= s7Data[0] << 24;
            tmp |= s7Data[1] << 16;
            tmp |= s7Data[2] << 8;
            tmp |= s7Data[3];
            return Float.intBitsToFloat(tmp);
        }
        return null;
    }

    private byte[] toS7Data(Object datatype) {
        if (datatype.getClass() == Boolean.class) {
            return new byte[]{(byte) (((Boolean) datatype) ? 0x01 : 0x00)};
        } else if (datatype.getClass() == Byte.class) {
            return new byte[]{((Byte) datatype)};
        } else if (datatype.getClass() == Short.class) {
            short intValue = (short) datatype;
            return new byte[]{(byte) ((intValue & 0xff00) >> 8), (byte) (intValue & 0xff)};
        } else if (datatype.getClass() == Integer.class) {
            int intValue = (int) datatype;
            return new byte[]{(byte) ((intValue & 0xff000000) >> 24), (byte) ((intValue & 0x00ff0000) >> 16),
                (byte) ((intValue & 0x0000ff00) >> 8), (byte) (intValue & 0xff)};
        } else if (datatype.getClass() == Calendar.class) {
            return null;
        } else if (datatype.getClass() == Float.class) {
            float floatValue = (float) datatype;
            int tmp = Float.floatToIntBits(floatValue);
            return new byte[]{(byte) ((tmp & 0xff000000) >> 24), (byte) ((tmp & 0x00ff0000) >> 16),
                (byte) ((tmp & 0x0000ff00) >> 8), (byte) (tmp & 0xff)};
        } else if (datatype.getClass() == String.class) {
            return null;
        }
        return null;
    }

    private VarParameter getVarParameter(ParameterType type, Address address, TransportSize transportSize) throws PlcProtocolException {
        // Depending on the address type, generate the corresponding type of request item.
        VarParameter varParameter = new VarParameter(type);
        if (!(address instanceof S7Address)) {
            throw new PlcProtocolException("Can only use S7Address types on S7 connection");
        }
        S7Address s7Address = (S7Address) address;
        if (s7Address instanceof S7DataBlockAddress) {
            S7DataBlockAddress s7DataBlockAddress = (S7DataBlockAddress) s7Address;
            varParameter.addRequestItem(new S7AnyVarItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                transportSize, (short) 1/*readRequest.getSize()*/,
                s7DataBlockAddress.getDataBlockNumber(), s7DataBlockAddress.getByteOffset(), (byte) 0));
        } else if (s7Address instanceof S7BitAddress) {
            S7BitAddress s7BitAddress = (S7BitAddress) s7Address;
            varParameter.addRequestItem(new S7AnyVarItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                transportSize, (short) 1/*readRequest.getSize()*/, (short) 0,
                s7Address.getByteOffset(), s7BitAddress.getBitOffset()));
        } else {
            varParameter.addRequestItem(new S7AnyVarItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                transportSize, (short) 1/*readRequest.getSize()*/, (short) 0,
                s7Address.getByteOffset(), (byte) 0));
        }
        return varParameter;
    }

    private VarPayload getVarPayload(ParameterType parameterType, VarItem varItem, DataTransportSize transportSize, byte[] data)
        throws PlcProtocolException {
        return new S7AnyVarPayload(parameterType, varItem, transportSize, data);
    }

}
