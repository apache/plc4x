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

import org.apache.plc4x.java.base.messages.PlcProtocolMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;

import java.util.List;

/**
 * Container Object for Requests to S7 which additionally stores information if the request was acknowledged (by the PLC?).
 *
 * @see S7Message for the other attributes.
 */
public class S7RequestMessage extends S7Message {

    private boolean acknowledged;

    public S7RequestMessage(MessageType messageType, short tpduReference, List<S7Parameter> s7Parameters,
                            List<S7Payload> s7Payloads, PlcProtocolMessage parent) {
        super(messageType, tpduReference, s7Parameters, s7Payloads, parent);
        acknowledged = false;
    }

    public boolean isAcknowledged() {
        return acknowledged;
    }

    public void setAcknowledged(boolean acknowledged) {
        this.acknowledged = acknowledged;
    }

}
