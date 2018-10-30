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
package org.apache.plc4x.java.s7.netty.util;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.VariableAddressingMode;

import java.util.List;

/**
 * When sending S7 messages we have to pay attention to the size of the message
 * in it's serialized form as we have to stay below the agreed upon maximum
 * PDU size. This helper calculates this message size and hereby helps implementing
 * different strategies for optimizing requests.
 */
public class S7RequestSizeCalculator {

    private S7RequestSizeCalculator() {
        // We don't want this to be instantiated.
    }

    /**
     * Calculate the size of a request message.
     *
     * @param requestMessage the request message we want to get the size for.
     * @return the size in bytes this message will have.
     */
    public static short getRequestMessageSize(S7RequestMessage requestMessage) {
        // The fixed header size.
        short size = 10;
        size += getRequestParametersSize(requestMessage.getParameters());
        size += getRequestPayloadsSize(requestMessage.getPayloads());
        return size;
    }

    /**
     * Calculates the size adding an item (parameter and eventually payload) would add to any existing
     * message.
     *
     * @param varParameterItem parameter item.
     * @param varPayloadItem corresponding payload (can be null).
     * @return size in bytes this item would add to an existing request message.
     */
    public static short getRequestItemTotalSize(VarParameterItem varParameterItem, VarPayloadItem varPayloadItem) {
        return (short) (getRequestReadWriteVarParameterItemSize(varParameterItem) +
            getRequestWriteVarPayloadItemSize(varPayloadItem));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Internal Request size calculations.
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private static short getRequestParametersSize(List<S7Parameter> parameters) {
        short l = 0;
        if (parameters != null) {
            for (S7Parameter parameter : parameters) {
                l += getRequestParameterSize(parameter);
            }
        }
        return l;
    }

    private static short getRequestParameterSize(S7Parameter parameter) {
        switch (parameter.getType()) {
            case READ_VAR:
            case WRITE_VAR:
                return getRequestReadWriteVarParameterSize((VarParameter) parameter);
            case SETUP_COMMUNICATION:
                return 8;
            default:
                throw new NotImplementedException("Not implemented");
        }
    }

    private static short getRequestReadWriteVarParameterSize(VarParameter varParameter) {
        // A var parameter always has a minimum size of 2 bytes (type of parameter, number of items)
        short length = 2;
        // Then come the items ...
        for (VarParameterItem varParameterItem : varParameter.getItems()) {
            length += getRequestReadWriteVarParameterItemSize(varParameterItem);
        }
        return length;
    }

    private static short getRequestReadWriteVarParameterItemSize(VarParameterItem varParameterItem) {
        VariableAddressingMode addressMode = varParameterItem.getAddressingMode();
        // S7 Any items have a fixed size of 12 bytes.
        if (addressMode == VariableAddressingMode.S7ANY) {
            return 12;
        } else {
            throw new NotImplementedException("Not implemented");
        }
    }

    private static short getRequestPayloadsSize(List<S7Payload> payloads) {
        short l = 0;
        if (payloads != null) {
            for (S7Payload payload : payloads) {
                l += getRequestPayloadSize(payload);
            }
        }
        return l;
    }

    private static short getRequestPayloadSize(S7Payload payload) {
        switch (payload.getType()) {
            case WRITE_VAR:
                return getRequestWriteVarPayloadSize((VarPayload) payload);
            default:
                throw new NotImplementedException("Not implemented");
        }
    }

    private static short getRequestWriteVarPayloadSize(VarPayload varPayload) {
        short length = 0;
        // Then come the items ...
        for (VarPayloadItem varPayloadItem : varPayload.getItems()) {
            length += getRequestWriteVarPayloadItemSize(varPayloadItem);
        }
        return length;
    }

    private static short getRequestWriteVarPayloadItemSize(VarPayloadItem varPayloadItem) {
        // A var payload item always has a minimum size of 4 bytes (return code, transport size, size (two bytes))
        short length = 4;
        length += varPayloadItem.getData().length;
        // It seems that bit payloads need a additional separating 0x00 byte.
        if(varPayloadItem.getDataTransportSize().isSizeInBits()) {
            length += 1;
        }
        return length;
    }

}
