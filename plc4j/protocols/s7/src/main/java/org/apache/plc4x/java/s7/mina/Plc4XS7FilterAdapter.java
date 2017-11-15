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
package org.apache.plc4x.java.s7.mina;

import org.apache.mina.core.filterchain.IoFilterAdapter;
import org.apache.mina.core.session.IoSession;
import org.apache.mina.core.write.WriteRequest;
import org.apache.mina.core.write.WriteRequestWrapper;
import org.apache.plc4x.java.exceptions.PlcException;
import org.apache.plc4x.java.mina.PlcRequestContainer;
import org.apache.plc4x.java.model.*;
import org.apache.plc4x.java.s7.mina.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.mina.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.mina.model.params.ReadVarParameter;
import org.apache.plc4x.java.s7.mina.model.params.items.S7AnyReadVarRequestItem;
import org.apache.plc4x.java.s7.mina.model.payloads.S7AnyReadVarPayload;
import org.apache.plc4x.java.s7.mina.model.types.MessageType;
import org.apache.plc4x.java.s7.mina.model.types.SpecificationType;
import org.apache.plc4x.java.s7.mina.model.types.TransportSize;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.types.Datatype;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4XS7FilterAdapter extends IoFilterAdapter {

    private static final AtomicInteger tpduGenerator = new AtomicInteger(1);

    private Map<Short, PlcRequestContainer> requests;

    public Plc4XS7FilterAdapter() {
        this.requests = new HashMap<>();
    }

    @Override
    public void filterWrite(NextFilter nextFilter, IoSession session, WriteRequest writeRequest) throws Exception {
        if(writeRequest.getMessage() instanceof PlcRequestContainer) {
            PlcRequestContainer readRequestContainer = (PlcRequestContainer) writeRequest.getMessage();
            PlcRequest readRequest = readRequestContainer.getRequest();

            // Try to get the correct S7 transport size for the given data type.
            // (Map PLC4X data type to S7 data type)
            TransportSize transportSize = getTransportSize(readRequest.getDatatype());
            if(transportSize == null) {
                throw new PlcException("Unknown transport size for datatype " + readRequest.getDatatype());
            }

            // Depending on the address type, generate the corresponding type of request item.
            ReadVarParameter readVarParameter = new ReadVarParameter();
            S7Address s7Address = (S7Address) readRequest.getAddress();
            if(!(readRequest.getAddress() instanceof S7Address)) {
                throw new PlcException("Can only use S7Address types on S7 connection");
            }
            if(s7Address instanceof S7DataBlockAddress) {
                S7DataBlockAddress s7DataBlockAddress = (S7DataBlockAddress) s7Address;
                readVarParameter.addRequestItem(new S7AnyReadVarRequestItem(
                    SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                    transportSize, (short) readRequest.getSize(),
                    s7DataBlockAddress.getDataBlockNumber(), s7DataBlockAddress.getByteOffset(), (byte) 0));
            } else if(s7Address instanceof S7BitAddress) {
                S7BitAddress s7BitAddress = (S7BitAddress) s7Address;
                readVarParameter.addRequestItem(new S7AnyReadVarRequestItem(
                    SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                    transportSize, (short) readRequest.getSize(), (short) 0,
                    s7Address.getByteOffset(), s7BitAddress.getBitOffset()));
            } else {
                readVarParameter.addRequestItem(new S7AnyReadVarRequestItem(
                    SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                    transportSize, (short) readRequest.getSize(), (short) 0,
                    s7Address.getByteOffset(), (byte) 0));
            }

            // Assemble the request.
            S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.JOB,
                (short) tpduGenerator.getAndIncrement(), Collections.singletonList(readVarParameter),
                Collections.emptyList());

            // Replace the writeRequest with the updated one.
            writeRequest = new WriteRequestWrapper(writeRequest) {
                @Override
                public Object getMessage() {
                    return s7ReadRequest;
                }
            };

            requests.put(s7ReadRequest.getTpduReference(), readRequestContainer);
        }/* else if(writeRequest.getMessage() instanceof PlcWriteRequest) {
            // TODO: To be implemented.
        }*/
        nextFilter.filterWrite(session, writeRequest);
    }

    @Override
    public void messageReceived(NextFilter nextFilter, IoSession session, Object message) throws Exception {
        if(message instanceof S7ResponseMessage) {
            S7ResponseMessage responseMessage = (S7ResponseMessage) message;
            short tpduReference = responseMessage.getTpduReference();
            if(requests.containsKey(tpduReference)) {
                PlcRequestContainer requestContainer = requests.remove(tpduReference);
                PlcResponse response = null;
                if(requestContainer.getRequest() instanceof PlcReadRequest) {
                    PlcReadRequest plcReadRequest = (PlcReadRequest) requestContainer.getRequest();
                    S7AnyReadVarPayload payload = responseMessage.getPayload(S7AnyReadVarPayload.class);
                    byte[] data = payload.getData();
                    response = new PlcReadResponse(plcReadRequest.getDatatype(), plcReadRequest.getAddress(),
                        plcReadRequest.getSize(), data);
                } else if(requestContainer.getRequest() instanceof PlcWriteRequest) {
                    PlcWriteRequest plcWriteRequest = (PlcWriteRequest) requestContainer.getRequest();
                }
                if(response != null) {
                    requestContainer.getResponseFuture().complete(response);
                }
            }
        }
        super.messageReceived(nextFilter, session, message);
    }

    private TransportSize getTransportSize(Datatype datatype) {
        switch (datatype) {
            case BIT:
                return TransportSize.BIT;
            case BYTE:
                return TransportSize.BYTE;
            case INTEGER:
                return TransportSize.INT;
            case FLOAT:
                return TransportSize.REAL;
            case STRING:
                return TransportSize.CHAR;
            case TIME:
                return TransportSize.TIME;
            case DATE:
                return TransportSize.DATE_AND_TIME;
            case TIMESTAMP:
                return TransportSize.DATE_AND_TIME;
        }
        return null;
    }

}
