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
import org.apache.plc4x.java.s7.netty.util.S7SizeHelper;

import java.util.Collection;
import java.util.Collections;
import java.util.Optional;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * While a SetupCommunication message is no problem, when reading or writing data,
 * situations could arise that have to be handled. The following situations have to
 * be handled:
 * - The number of request items is so big, that the resulting PDU would exceed the
 *   agreed upon PDU size: The request has to be split up into multiple requests.
 * - If large blocks of data are requested by request items the result of a request
 *   could exceed the PDU size: The requests has to be split up into multiple requests
 *   where each requests response doesn't exceed the PDU size.
 *
 * The following optimizations should be implemented:
 * - If blocks are read which are in near proximity to each other it could be better
 *   to replace multiple requests by one that includes multiple blocks.
 * - Rearranging the order of request items could reduce the number of needed PDUs.
 */
public class DefaultS7MessageProcessor implements S7MessageProcessor {

    private final int pduSize;
    private AtomicInteger tpduRefGen;

    public DefaultS7MessageProcessor(int pduSize) {
        this.pduSize = pduSize;
        this.tpduRefGen = new AtomicInteger(1);
    }

    @Override
    public Collection<? extends S7Message> process(S7Message s7Message) {
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
            S7RequestMessage originalRequestMessage = (S7RequestMessage) s7Message;
            Optional<VarParameter> varParameterOptional = originalRequestMessage.getParameter(VarParameter.class);
            if (varParameterOptional.isPresent()) {
                VarParameter varParameter = varParameterOptional.get();

                // TODO: This is a work inprogress ...
                S7CompositeRequestMessage compositeRequestMessage = new S7CompositeRequestMessage(originalRequestMessage);
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
