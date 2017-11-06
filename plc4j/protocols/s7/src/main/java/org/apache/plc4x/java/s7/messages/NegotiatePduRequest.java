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

import org.apache.plc4x.java.s7.messages.s7.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.messages.s7.params.SetupCommunicationParameter;
import org.apache.plc4x.java.s7.messages.s7.types.MessageType;

import java.util.Collections;

public class NegotiatePduRequest extends S7RequestMessage {

    public NegotiatePduRequest(short pduLength) {
        this((short) 0x0008, (short) 0x0008, pduLength);
    }

    public NegotiatePduRequest(short maxAmqCaller, short maxAmqCallee, short pduLength) {
        super(MessageType.JOB, Collections.singletonList(
            new SetupCommunicationParameter(maxAmqCaller, maxAmqCallee, pduLength)), null);
    }

}
