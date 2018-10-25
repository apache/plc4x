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
package org.apache.plc4x.nifi;

import org.apache.nifi.annotation.behavior.InputRequirement;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.json.simple.JSONObject;

import java.io.IOException;
import java.io.OutputStreamWriter;
import java.util.concurrent.ExecutionException;

@Tags({"plc4x-source"})
@InputRequirement(InputRequirement.Requirement.INPUT_FORBIDDEN)
@CapabilityDescription("Processor able to read data from industrial PLCs using Apache PLC4X")
@WritesAttributes({@WritesAttribute(attribute="value", description="some value")})
public class Plc4xSourceProcessor extends BasePlc4xProcessor {

    @Override
    public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
        // Get an instance of a component able to read from a PLC.
        PlcConnection connection = getConnection();

        // Prepare the request.
        if (!connection.readRequestBuilder().isPresent()) {
            throw new ProcessException("Writing not supported by connection");
        }

        FlowFile flowFile = session.create();
        session.append(flowFile, out -> {
            try {
                PlcReadRequest.Builder builder = connection.readRequestBuilder().get();
                getFields().forEach(field -> {
                    String address = getAddress(field);
                    if(address != null) {
                        builder.addItem(field, address);
                    }
                });
                PlcReadRequest readRequest = builder.build();
                PlcReadResponse response = readRequest.execute().get();
                JSONObject obj = new JSONObject();
                for (String fieldName : response.getFieldNames()) {
                    for(int i = 0; i < response.getNumberOfValues(fieldName); i++) {
                        Object value = response.getObject(fieldName, i);
                        obj.put(fieldName, value);
                    }
                }
                obj.writeJSONString(new OutputStreamWriter(out));
            } catch (InterruptedException | ExecutionException e) {
                throw new IOException(e);
            }
        });
        session.transfer(flowFile, SUCCESS);
    }

}
