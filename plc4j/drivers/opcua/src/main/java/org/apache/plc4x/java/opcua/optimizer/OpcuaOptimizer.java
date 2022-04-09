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
package org.apache.plc4x.java.opcua.optimizer;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.opcua.field.OpcuaField;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;

import java.util.LinkedHashMap;
import java.util.LinkedList;
import java.util.List;

public class OpcuaOptimizer extends BaseOptimizer{

    @Override
    protected List<PlcRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        List<PlcRequest> processedRequests = new LinkedList<>();

        // List of all items in the current request.
        LinkedHashMap<String, PlcField> curFields = new LinkedHashMap<>();

        for (String fieldName : readRequest.getFieldNames()) {
            OpcuaField field = (OpcuaField) readRequest.getField(fieldName);
            curFields.put(fieldName, field);
        }

        // Create a new PlcReadRequest from the remaining field items.
        if(!curFields.isEmpty()) {
            processedRequests.add(new DefaultPlcReadRequest(
                ((DefaultPlcReadRequest) readRequest).getReader(), curFields));
        }

        return processedRequests;
    }


}
