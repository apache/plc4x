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
package org.apache.plc4x.java.s7.netty.model.messages;

import org.apache.plc4x.java.netty.Message;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;

import java.util.List;

public abstract class S7Message extends Message {

    private final MessageType messageType;
    private final short tpduReference;
    private final List<S7Parameter> parameters;
    private final List<S7Payload> payloads;

    protected S7Message(MessageType messageType, short tpduReference, List<S7Parameter> parameters, List<S7Payload> payloads) {
        super(null);
        this.messageType = messageType;
        this.tpduReference = tpduReference;
        this.parameters = parameters;
        this.payloads = payloads;
    }

    public MessageType getMessageType() {
        return messageType;
    }

    public short getTpduReference() {
        return tpduReference;
    }

    public List<S7Parameter> getParameters() {
        return parameters;
    }

    public <T> T getParameter(Class<? extends T> parameterType) {
        if (parameters != null) {
            for (S7Parameter s7Parameter : parameters) {
                if (s7Parameter.getClass() == parameterType) {
                    return (T) s7Parameter;
                }
            }
        }
        return null;
    }

    public List<S7Payload> getPayloads() {
        return payloads;
    }

    public <T> T getPayload(Class<? extends T> payloadType) {
        if (payloads != null) {
            for (S7Payload s7Payload : payloads) {
                if (s7Payload.getClass() == payloadType) {
                    return (T) s7Payload;
                }
            }
        }
        return null;
    }

}