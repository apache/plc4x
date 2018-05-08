package org.apache.plc4x.java.s7.netty.strategies;/*
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

import org.apache.plc4x.java.s7.netty.model.messages.S7CompositeRequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;

import java.util.Collection;
import java.util.Collections;
import java.util.Optional;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * This is just a dummy processor that splits up an incomming message with mutltiple ver items
 * into one message per item. It's not really useful, but might be useful when tracking down problems.
 */
public class SingleAddressPerMessageS7MessageProcessor implements S7MessageProcessor {

    private AtomicInteger tpduRefGen;

    public SingleAddressPerMessageS7MessageProcessor() {
        this.tpduRefGen = new AtomicInteger(1);
    }

    @Override
    public Collection<? extends S7Message> process(S7Message s7Message, int pduSize) {
        // The following considerations have to be taken into account:
        // - The size of all parameters and payloads of a message cannot exceed the negotiated PDU size
        // - When reading data, the size of the returned data cannot exceed the negotiated PDU size
        //
        // Examples:
        // - Size of the request exceeds the maximum
        //  When having a negotiated max PDU size of 256, the maximum size of individual addresses can be at most 18
        //  If more are sent, the S7 will respond with a frame error.
        // - Size of the response exceeds the maximum
        //  When reading two Strings of each 200 bytes length, the size of the request is ok, however the PLC would
        //  have to send back 400 bytes of String data, which would exceed the PDU size. In this case the first String
        //  is correctly returned, but for the second item the PLC will return a code of 0x03 = Access Denied

        if(s7Message instanceof S7RequestMessage) {
            // TODO: This currently splits up every address area into it's own request ... this is just for testing ...
            S7RequestMessage originalRequestMessage = (S7RequestMessage) s7Message;
            Optional<VarParameter> varParameterOptional = originalRequestMessage.getParameter(VarParameter.class);
            if (varParameterOptional.isPresent()) {
                S7CompositeRequestMessage compositeRequestMessage = new S7CompositeRequestMessage(originalRequestMessage);
                VarParameter varParameter = varParameterOptional.get();
                for (VarParameterItem varParameterItem : varParameter.getItems()) {
                    S7Parameter parameter = new VarParameter(varParameter.getType(), Collections.singletonList(varParameterItem));
                    S7RequestMessage subMessage = new S7RequestMessage(
                        s7Message.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                        Collections.singletonList(parameter), Collections.emptyList(), compositeRequestMessage);
                    compositeRequestMessage.addRequestMessage(subMessage);
                }
                return compositeRequestMessage.getRequestMessages();
            }
        }

        return Collections.singletonList(s7Message);
    }

}
