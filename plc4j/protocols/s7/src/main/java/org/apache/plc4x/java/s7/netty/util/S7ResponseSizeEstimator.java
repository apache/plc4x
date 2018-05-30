package org.apache.plc4x.java.s7.netty.util;
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

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;

import java.util.List;

/**
 * When sending S7 messages we have to also pay attention to the expected size
 * of a response for a given request. This helper estimates this response size
 * and hereby helps implementing different strategies for optimizing requests.
 */
public class S7ResponseSizeEstimator {

    private S7ResponseSizeEstimator() {
        // We don't want this to be instantiated.
    }

    /**
     * Estimate the size of the response for a given request object.
     * This is needed in order to decide if a request will be able to be processed
     * as S7 PLCs will return errors if a request will produce a response that exceeds
     * the agreed upon PDU size.
     *
     * @param requestMessage the request message for which we want to estimate the response size for.
     * @return the estimated size in bytes the response for this message will have.
     */
    public static short getEstimatedResponseMessageSize(S7RequestMessage requestMessage) {
        // The fixed header size (two bytes more than the request for error class and code values).
        short size = 12;
        size += getEstimatedResponseParametersSize(requestMessage.getParameters());
        size += getEstimatedResponsePayloadsSize(requestMessage.getParameters());
        return size;
    }

    /**
     * Calculates the estimated size adding an item (parameter and eventually payload) would add to the
     * response for a current request message. While for the parameters the size of the items depends on
     * if they are requests or responses, the payload depends on if it's a read or write operation.
     *
     * @param varParameterItem parameter item.
     * @param varPayloadItem payload item.
     * @return size in bytes in the corresponding response message, this item would add to an existing response message.
     */
    public static short getEstimatedResponseReadItemTotalSize(VarParameterItem varParameterItem,
                                                              VarPayloadItem varPayloadItem) {
        // If both are provided, this is a write request.
        if(varPayloadItem != null) {
            // When writing values the item size in the response is constantly 1 byte.
            return getEstimatedResponseWriteVarPayloadItemSize(varParameterItem);
        }
        // If the payload is empty, this is a read request and the response could contain real information.
        else {
            return getEstimatedResponseReadVarPayloadItemSize(varParameterItem);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Internal Response size estimations.
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private static short getEstimatedResponseParametersSize(List<S7Parameter> parameters) {
        short l = 0;
        if (parameters != null) {
            for (S7Parameter parameter : parameters) {
                l += getEstimatedResponseParameterSize(parameter);
            }
        }
        return l;
    }

    private static short getEstimatedResponseParameterSize(S7Parameter parameter) {
        if (parameter == null) {
            return 0;
        }

        switch (parameter.getType()) {
            case READ_VAR:
            case WRITE_VAR:
                return getEstimatedResponseReadWriteVarParameterSize((VarParameter) parameter);
            case SETUP_COMMUNICATION:
                return 8;
            default:
                throw new NotImplementedException("Not implemented");
        }
    }

    private static short getEstimatedResponseReadWriteVarParameterSize(VarParameter varParameter) {
        // In the response the var parameter only consists of the type and number of items.
        return 2;
    }

    private static short getEstimatedResponsePayloadsSize(List<S7Parameter> parameters) {
        short l = 0;
        if (parameters != null) {
            for (S7Parameter parameter : parameters) {
                l += getEstimatedResponsePayloadSize(parameter);
            }
        }
        return l;
    }

    private static short getEstimatedResponsePayloadSize(S7Parameter parameter) {
        if (parameter == null) {
            return 0;
        }

        switch (parameter.getType()) {
            case READ_VAR:
                return getEstimatedResponseReadVarPayloadSize((VarParameter) parameter);
            case WRITE_VAR:
                return getEstimatedResponseWriteVarPayloadSize((VarParameter) parameter);
            case SETUP_COMMUNICATION:
                return 0;
            default:
                throw new NotImplementedException("Not implemented");
        }
    }

    private static short getEstimatedResponseReadVarPayloadSize(VarParameter varParameter) {
        short length = 0;
        // Then come the items ...
        for (VarParameterItem varParameterItem : varParameter.getItems()) {
            length += getEstimatedResponseReadVarPayloadItemSize(varParameterItem);
        }
        return length;
    }

    private static short getEstimatedResponseReadVarPayloadItemSize(VarParameterItem varParameterItem) {
        // A var payload item always has a minimum size of 4 bytes (return code, transport size, size (two bytes))
        short length = 4;
        S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItem;
        length +=
            s7AnyVarParameterItem.getNumElements() * s7AnyVarParameterItem.getTransportSize().getSizeInBytes();
        // Get the corresponding "Data Transport Size" for a given "Transport Size".
        DataTransportSize dataTransportSize =
            DataTransportSize.getForTransportSize(s7AnyVarParameterItem.getTransportSize());
        // It seems that bit payloads need a additional separating 0x00 byte.
        if((dataTransportSize != null) && dataTransportSize.isSizeInBits()) {
            length += 1;
        }
        return length;
    }

    private static short getEstimatedResponseWriteVarPayloadSize(VarParameter varParameter) {
        short length = 0;
        for (VarParameterItem varParameterItem : varParameter.getItems()) {
            length += getEstimatedResponseWriteVarPayloadItemSize(varParameterItem);
        }
        return length;
    }

    private static short getEstimatedResponseWriteVarPayloadItemSize(VarParameterItem varParameterItem) {
        // This is just one byte containing the error code.
        return 1;
    }

}
