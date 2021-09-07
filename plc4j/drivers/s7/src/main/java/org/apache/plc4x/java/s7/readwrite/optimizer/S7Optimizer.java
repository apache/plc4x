/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.readwrite.optimizer;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.s7.readwrite.field.S7Field;
import org.apache.plc4x.java.s7.readwrite.types.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.types.TransportSize;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.utils.FieldValueItem;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;

import java.util.*;

public class S7Optimizer extends BaseOptimizer {

    public static final int EMPTY_READ_REQUEST_SIZE = new S7MessageRequest(0, new S7ParameterReadVarRequest(
        new S7VarRequestParameterItem[0]), null).getLengthInBytes();
    public static final int EMPTY_READ_RESPONSE_SIZE = new S7MessageResponseData(0, new S7ParameterReadVarResponse(
        (short) 0), new S7PayloadReadVarResponse(new S7VarPayloadDataItem[0]), (short) 0, (short) 0).getLengthInBytes();
    public static final int EMPTY_WRITE_REQUEST_SIZE = new S7MessageRequest(0, new S7ParameterWriteVarRequest(
        new S7VarRequestParameterItem[0]), new S7PayloadWriteVarRequest(new S7VarPayloadDataItem[0])).getLengthInBytes();
    public static final int EMPTY_WRITE_RESPONSE_SIZE = new S7MessageResponseData(0, new S7ParameterWriteVarResponse(
        (short) 0), new S7PayloadWriteVarResponse(new S7VarPayloadStatusItem[0]), (short) 0, (short) 0).getLengthInBytes();
    public static final int S7_ADDRESS_ANY_SIZE = 2 +
        new S7AddressAny(TransportSize.INT, 1, 1, MemoryArea.DATA_BLOCKS, 1, (byte) 0).getLengthInBytes();

    @Override
    protected List<PlcRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        S7DriverContext s7DriverContext = (S7DriverContext) driverContext;
        List<PlcRequest> processedRequests = new LinkedList<>();

        // This calculates the size of the header for the request and response.
        int curRequestSize = EMPTY_READ_REQUEST_SIZE;
        // An empty response has the same size as an empty request.
        int curResponseSize = EMPTY_READ_RESPONSE_SIZE;

        // List of all items in the current request.
        LinkedHashMap<String, PlcField> curFields = new LinkedHashMap<>();

        for (String fieldName : readRequest.getFieldNames()) {
            S7Field field = (S7Field) readRequest.getField(fieldName);

            int readRequestItemSize = S7_ADDRESS_ANY_SIZE;
            int readResponseItemSize = 4 + (field.getNumberOfElements() * field.getDataType().getSizeInBytes());
            // If it's an odd number of bytes, add one to make it even
            if (readResponseItemSize % 2 == 1) {
                readResponseItemSize++;
            }

            // If adding the item would not exceed the sizes, add it to the current request.
            if (((curRequestSize + readRequestItemSize) <= s7DriverContext.getPduSize()) &&
                ((curResponseSize + readResponseItemSize) <= s7DriverContext.getPduSize())) {
                // Increase the current request sizes.
                curRequestSize += readRequestItemSize;
                curResponseSize += readResponseItemSize;

                // Add the item.
            }
            // If they would exceed, start a new request.
            else {
                // Create a new PlcReadRequest containing the current field item.
                processedRequests.add(new DefaultPlcReadRequest(
                    ((DefaultPlcReadRequest) readRequest).getReader(), curFields));

                // Reset the size and item lists.
                curRequestSize = EMPTY_READ_REQUEST_SIZE + readRequestItemSize;
                curResponseSize = EMPTY_READ_RESPONSE_SIZE + readResponseItemSize;
                curFields = new LinkedHashMap<>();

                // Splitting of huge fields not yet implemented, throw an exception instead.
                if(((curRequestSize + readRequestItemSize) > s7DriverContext.getPduSize()) &&
                    ((curResponseSize + readResponseItemSize) > s7DriverContext.getPduSize())) {
                    throw new PlcRuntimeException("Field size exceeds maximum payload for one item.");
                }
            }
            curFields.put(fieldName, field);
        }

        // Create a new PlcReadRequest from the remaining field items.
        if(!curFields.isEmpty()) {
            processedRequests.add(new DefaultPlcReadRequest(
                ((DefaultPlcReadRequest) readRequest).getReader(), curFields));
        }

        return processedRequests;
    }

    @Override
    protected List<PlcRequest> processWriteRequest(PlcWriteRequest writeRequest, DriverContext driverContext) {
        S7DriverContext s7DriverContext = (S7DriverContext) driverContext;
        List<PlcRequest> processedRequests = new LinkedList<>();

        // This calculates the size of the header for the request and response.
        int curRequestSize = EMPTY_WRITE_REQUEST_SIZE;
        // An empty response has the same size as an empty request.
        int curResponseSize = EMPTY_WRITE_RESPONSE_SIZE;

        // List of all items in the current request.
        LinkedHashMap<String, FieldValueItem> curFields = new LinkedHashMap<>();

        for (String fieldName : writeRequest.getFieldNames()) {
            S7Field field = (S7Field) writeRequest.getField(fieldName);
            PlcValue value = writeRequest.getPlcValue(fieldName);

            int writeRequestItemSize = S7_ADDRESS_ANY_SIZE + 4/* Size of Payload item header*/;
            if (field.getDataType() == TransportSize.BOOL) {
                writeRequestItemSize += Math.ceil((double) field.getNumberOfElements() / 8);
            } else {
                writeRequestItemSize += (field.getNumberOfElements() * field.getDataType().getSizeInBytes());
            }
            // If it's an odd number of bytes, add one to make it even
            if (writeRequestItemSize % 2 == 1) {
                writeRequestItemSize++;
            }
            int writeResponseItemSize = 4;

            // If adding the item would not exceed the sizes, add it to the current request.
            if (((curRequestSize + writeRequestItemSize) <= s7DriverContext.getPduSize()) &&
                ((curResponseSize + writeResponseItemSize) <= s7DriverContext.getPduSize())) {
                // Increase the current request sizes.
                curRequestSize += writeRequestItemSize;
                curResponseSize += writeResponseItemSize;

                // Add the item.
            }
            // If adding them would exceed, start a new request.
            else {
                // Create a new PlcWriteRequest containing the current field item.
                processedRequests.add(new DefaultPlcWriteRequest(
                    ((DefaultPlcWriteRequest) writeRequest).getWriter(), curFields));

                // Reset the size and item lists.
                curRequestSize = EMPTY_WRITE_REQUEST_SIZE + writeRequestItemSize;
                curResponseSize = EMPTY_WRITE_RESPONSE_SIZE + writeResponseItemSize;
                curFields = new LinkedHashMap<>();

                // Splitting of huge fields not yet implemented, throw an exception instead.
                if(((curRequestSize + writeRequestItemSize) > s7DriverContext.getPduSize()) &&
                    ((curResponseSize + writeResponseItemSize) > s7DriverContext.getPduSize())) {
                    throw new PlcRuntimeException("Field size exceeds maximum payload for one item.");
                }
            }
            curFields.put(fieldName, new FieldValueItem(field, value));
        }

        // Create a new PlcWriteRequest from the remaining field items.
        if(!curFields.isEmpty()) {
            processedRequests.add(new DefaultPlcWriteRequest(
                ((DefaultPlcWriteRequest) writeRequest).getWriter(), curFields));
        }

        return processedRequests;
    }

}
