/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.nifi;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import org.apache.nifi.annotation.behavior.InputRequirement;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;

@Tags({"plc4x", "get", "input", "source", "attributes"})
@InputRequirement(InputRequirement.Requirement.INPUT_FORBIDDEN)
@CapabilityDescription("Processor able to read data from industrial PLCs using Apache PLC4X")
@WritesAttributes({@WritesAttribute(attribute="value", description="some value")})
public class Plc4xSourceProcessor extends BasePlc4xProcessor {

    @Override
    public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {

        final ComponentLog logger = getLogger();
        // Get an instance of a component able to read from a PLC.
        try(PlcConnection connection = getConnectionManager().getConnection(getConnectionString())) {

            // Prepare the request.
            if (!connection.getMetadata().canRead()) {
                throw new ProcessException("Writing not supported by connection");
            }

            FlowFile flowFile = session.create();
            try {
                PlcReadRequest.Builder builder = connection.readRequestBuilder();
                Map<String,String> addressMap = getPlcAddressMap(context, flowFile);
                final Map<String, PlcTag> tags = getSchemaCache().retrieveTags(addressMap);

                if (tags != null){
                    for (Map.Entry<String,PlcTag> tag : tags.entrySet()){
                        builder.addTag(tag.getKey(), tag.getValue());
                    }
                } else {
                    if (debugEnabled)
                        logger.debug("PlcTypes resolution not found in cache and will be added with key: " + addressMap.toString());
                    for (Map.Entry<String,String> entry: addressMap.entrySet()){
                        builder.addTagAddress(entry.getKey(), entry.getValue());
                    }
                }

                PlcReadRequest readRequest = builder.build();
                PlcReadResponse response = readRequest.execute().get(this.timeout, TimeUnit.MILLISECONDS);
                Map<String, String> attributes = new HashMap<>();
                for (String tagName : response.getTagNames()) {
                    for (int i = 0; i < response.getNumberOfValues(tagName); i++) {
                        Object value = response.getObject(tagName, i);
                        attributes.put(tagName, String.valueOf(value));
                    }
                }
                flowFile = session.putAllAttributes(flowFile, attributes); 
                
                if (tags == null){
                    if (debugEnabled)
                        logger.debug("Adding PlcTypes resolution into cache with key: " + addressMap.toString());
                    getSchemaCache().addSchema(
                        addressMap, 
                        readRequest.getTagNames(),
                        readRequest.getTags(),
                        null
                    );
                }

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                throw new ProcessException(e);
            } catch (ExecutionException e) {
                throw new ProcessException(e);
            } catch (TimeoutException e) {
                session.remove(flowFile);
            }
            session.transfer(flowFile, REL_SUCCESS);
        } catch (ProcessException e) {
            throw e;
        } catch (Exception e) {
            throw new ProcessException("Got an error while trying to get a connection", e);
        }
    }

}
