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
package org.apache.plc4x.java.s7.messages;

import org.apache.plc4x.java.s7.messages.s7.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.messages.s7.params.S7Parameter;
import org.apache.plc4x.java.s7.messages.s7.params.SetupCommunicationParameter;
import org.apache.plc4x.java.s7.messages.s7.types.MessageType;

import java.util.List;

public class NegotiatePduResponse extends S7ResponseMessage {

    public NegotiatePduResponse(MessageType messageType, List<S7Parameter> s7Parameters,
                                byte errorClass, byte errorCode) {
        super(messageType, s7Parameters, null, errorClass, errorCode);
    }

    public short getPduLength() {
        if(getS7Parameters() != null) {
            for (S7Parameter s7Parameter : getS7Parameters()) {
                if(s7Parameter instanceof SetupCommunicationParameter) {
                    return ((SetupCommunicationParameter) s7Parameter).getPduLength();
                }
            }
        }
        return -1;
    }

}
