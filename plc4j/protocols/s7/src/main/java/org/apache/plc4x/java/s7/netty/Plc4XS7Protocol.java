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
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequestContainer;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.ReadVarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyReadVarRequestItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7AnyReadVarPayload;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.apache.plc4x.java.s7.netty.model.types.SpecificationType;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

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
            ReadVarParameter readVarParameter = new ReadVarParameter();
            if (!(readRequest.getAddress() instanceof S7Address)) {
                throw new PlcException("Can only use S7Address types on S7 connection");
            }
            S7Address s7Address = (S7Address) readRequest.getAddress();
            if (s7Address instanceof S7DataBlockAddress) {
                S7DataBlockAddress s7DataBlockAddress = (S7DataBlockAddress) s7Address;
                readVarParameter.addRequestItem(new S7AnyReadVarRequestItem(
                    SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                    transportSize, (short) 1/*readRequest.getSize()*/,
                    s7DataBlockAddress.getDataBlockNumber(), s7DataBlockAddress.getByteOffset(), (byte) 0));
            } else if (s7Address instanceof S7BitAddress) {
                S7BitAddress s7BitAddress = (S7BitAddress) s7Address;
                readVarParameter.addRequestItem(new S7AnyReadVarRequestItem(
                    SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                    transportSize, (short) 1/*readRequest.getSize()*/, (short) 0,
                    s7Address.getByteOffset(), s7BitAddress.getBitOffset()));
            } else {
                readVarParameter.addRequestItem(new S7AnyReadVarRequestItem(
                    SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                    transportSize, (short) 1/*readRequest.getSize()*/, (short) 0,
                    s7Address.getByteOffset(), (byte) 0));
            }

            // Assemble the request.
            S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.JOB,
                (short) tpduGenerator.getAndIncrement(), Collections.singletonList(readVarParameter),
                Collections.emptyList());

            requests.put(s7ReadRequest.getTpduReference(), msg);

            out.add(s7ReadRequest);
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
                    PlcReadRequest plcReadRequest = (PlcReadRequest) requestContainer.getRequest();
                    S7AnyReadVarPayload payload = responseMessage.getPayload(S7AnyReadVarPayload.class);
                    byte[] data = payload.getData();
                    Object value = fromS7Data(plcReadRequest.getDatatype(), data);
                    response = plcReadRequest.createResponse(value);
                } else if (requestContainer.getRequest() instanceof PlcWriteRequest) {
                    PlcWriteRequest plcWriteRequest = (PlcWriteRequest) requestContainer.getRequest();
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
        } else if (datatype == Calendar.class) {
            return TransportSize.DATE_AND_TIME;
        } else if (datatype == Float.class) {
            return TransportSize.REAL;
        } else if (datatype == Integer.class) {
            return TransportSize.INT;
        } else if (datatype == String.class) {
            return TransportSize.CHAR;
        }
        return null;
    }

    private Object fromS7Data(Class<?> datatype, byte[] s7Data) {
        if (datatype == Boolean.class) {
            return (s7Data[0] & 0x01) == 0x01;
        } else if (datatype == Byte.class) {
            return s7Data[0];
        }
        return null;
    }

    private byte[] toS7Data(Object datatype) {
        if (datatype.getClass() == Boolean.class) {
            return new byte[]{(byte) (((Boolean) datatype) ? 0x01 : 0x00)};
        } else if (datatype.getClass() == Byte.class) {
            return new byte[]{((Byte) datatype)};
        } else if (datatype.getClass() == Calendar.class) {

        } else if (datatype.getClass() == Float.class) {

        } else if (datatype.getClass() == Integer.class) {

        } else if (datatype.getClass() == String.class) {

        }
        return null;
    }

}
