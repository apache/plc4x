package org.apache.plc4x.java.s7.netty.model.messages;/*
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

import org.apache.plc4x.java.api.messages.ProtocolMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

import java.util.Collection;
import java.util.LinkedList;
import java.util.List;

public class S7CompositeRequestMessage implements ProtocolMessage {

    private S7RequestMessage originalRequest;

    private Collection<S7RequestMessage> requests;
    private Collection<S7ResponseMessage> responses;

    public S7CompositeRequestMessage(S7RequestMessage originalRequest) {
        this.originalRequest = originalRequest;
        this.requests = new LinkedList<>();
        this.responses = new LinkedList<>();
    }

    @Override
    public ProtocolMessage getParent() {
        return originalRequest;
    }

    /**
     * A {@link S7CompositeRequestMessage} is only acknowledged, if all children are acknowledged.
     *
     * @return true if all children are acknowledged.
     */
    public boolean isAcknowledged() {
        for (S7RequestMessage requestMessage : requests) {
            if(!requestMessage.isAcknowledged()) {
                return false;
            }
        }
        return true;
    }

    public void addRequestMessage(S7RequestMessage requestMessage) {
        requests.add(requestMessage);
    }

    public Collection<S7RequestMessage> getRequestMessages() {
        return requests;
    }

    public void addResponseMessage(S7ResponseMessage responseMessage) {
        responses.add(responseMessage);
    }

    public S7ResponseMessage getResponseMessage() {
        S7ResponseMessage firstResponse = null;
        short tpduRerefence = originalRequest.getTpduReference();
        List<S7Parameter> s7Parameters = new LinkedList<>();
        List<S7Payload > s7Payloads = new LinkedList<>();
        byte errorClass = 0;
        byte errorCode = 0;
        VarParameter readVarParameter = null;
        VarParameter writeVarParameter = null;
        VarPayload readVarPayload = null;
        VarPayload writeVarPayload = null;
        for (S7ResponseMessage response : responses) {
            if(firstResponse == null) {
                firstResponse = response;
            }
            // Some parameters have to be merged. In case of read and write parameters
            // their items have to be merged into one single parameter.
            for(S7Parameter parameter : response.getParameters()) {
                if (parameter.getType() == ParameterType.READ_VAR) {
                    if (readVarParameter == null) {
                        readVarParameter = (VarParameter) parameter;
                        s7Parameters.add(parameter);
                    } else {
                        readVarParameter.mergeParameter((VarParameter) parameter);
                    }
                } else if (parameter.getType() == ParameterType.WRITE_VAR) {
                    if (writeVarParameter == null) {
                        writeVarParameter = (VarParameter) parameter;
                        s7Parameters.add(parameter);
                    } else {
                        writeVarParameter.mergeParameter((VarParameter) parameter);
                    }
                } else {
                    s7Parameters.add(parameter);
                }
            }

            // Some payloads have to be merged. In case of read and write payloads
            // their items have to be merged into one single payload.
            for(S7Payload payload : response.getPayloads()) {
                if(payload.getType() == ParameterType.READ_VAR) {
                    if (readVarPayload == null) {
                        readVarPayload = (VarPayload) payload;
                        s7Payloads.add(payload);
                    } else {
                        readVarPayload.mergePayload((VarPayload) payload);
                    }
                } else if(payload.getType() == ParameterType.WRITE_VAR) {
                    if(writeVarPayload == null) {
                        writeVarPayload = (VarPayload) payload;
                        s7Payloads.add(payload);
                    } else {
                        writeVarPayload.mergePayload((VarPayload) payload);
                    }
                } else {
                    s7Payloads.add(payload);
                }
            }
        }
        if(firstResponse != null) {
            MessageType messageType = firstResponse.getMessageType();
            return new S7ResponseMessage(messageType, tpduRerefence, s7Parameters, s7Payloads, errorClass, errorCode);
        }
        return null;
    }

}
