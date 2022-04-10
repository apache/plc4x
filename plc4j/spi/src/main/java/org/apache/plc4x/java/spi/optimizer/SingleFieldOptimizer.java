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
package org.apache.plc4x.java.spi.optimizer;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.utils.FieldValueItem;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;

/**
 * Optimizer variant which automatically splits a multi-item request up into multiple single item requests.
 */
public class SingleFieldOptimizer extends BaseOptimizer {

    @Override
    protected List<PlcRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        if(readRequest.getNumberOfFields() == 1) {
            return Collections.singletonList(readRequest);
        }
        List<PlcRequest> subRequests = new ArrayList<>(readRequest.getNumberOfFields());
        for (String fieldName : readRequest.getFieldNames()) {
            PlcField field = readRequest.getField(fieldName);
            PlcReadRequest subRequest = new DefaultPlcReadRequest(
                ((DefaultPlcReadRequest) readRequest).getReader(),
                new LinkedHashMap<>(Collections.singletonMap(fieldName, field)));
            subRequests.add(subRequest);
        }
        return subRequests;
    }

    @Override
    protected List<PlcRequest> processWriteRequest(PlcWriteRequest writeRequest, DriverContext driverContext) {
        if(writeRequest.getNumberOfFields() == 1) {
            return Collections.singletonList(writeRequest);
        }
        List<PlcRequest> subRequests = new ArrayList<>(writeRequest.getNumberOfFields());
        for (String fieldName : writeRequest.getFieldNames()) {
            PlcField field = writeRequest.getField(fieldName);
            PlcValue value = writeRequest.getPlcValue(fieldName);
            PlcWriteRequest subRequest = new DefaultPlcWriteRequest(
                ((DefaultPlcWriteRequest) writeRequest).getWriter(),
                new LinkedHashMap<>(Collections.singletonMap(fieldName, new FieldValueItem(field, value))));
            subRequests.add(subRequest);
        }
        return subRequests;
    }

}
