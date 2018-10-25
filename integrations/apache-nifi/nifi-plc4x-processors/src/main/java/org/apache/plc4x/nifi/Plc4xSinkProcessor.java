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
import org.apache.nifi.annotation.behavior.ReadsAttribute;
import org.apache.nifi.annotation.behavior.ReadsAttributes;
import org.apache.nifi.annotation.behavior.TriggerSerially;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;

import java.util.concurrent.CompletableFuture;

@TriggerSerially
@Tags({"plc4x-sink"})
@InputRequirement(InputRequirement.Requirement.INPUT_REQUIRED)
@CapabilityDescription("Processor able to write data to industrial PLCs using Apache PLC4X")
@ReadsAttributes({@ReadsAttribute(attribute="value", description="some value")})
public class Plc4xSinkProcessor extends BasePlc4xProcessor {

    @Override
    public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
        FlowFile flowFile = session.get();

        // Abort if there's nothing to do.
        if (flowFile == null) {
            return;
        }

        // Get an instance of a component able to write to a PLC.
        PlcConnection connection = getConnection();
        if (!connection.writeRequestBuilder().isPresent()) {
            throw new ProcessException("Writing not supported by connection");
        }

        // Prepare the request.
        PlcWriteRequest.Builder builder = connection.writeRequestBuilder().get();
        flowFile.getAttributes().forEach((field, value) -> {
            String address = getAddress(field);
            if(address != null) {
                builder.addItem(field, address, value);
            }
        });
        PlcWriteRequest writeRequest = builder.build();

        // Send the request to the PLC.
        CompletableFuture<? extends PlcWriteResponse> future = writeRequest.execute();
        future.whenComplete((response, throwable) -> {
            if (throwable != null) {
                session.transfer(session.create(), FAILURE);
            } else {
                session.transfer(session.create(), SUCCESS);
            }
        });
    }

}
