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

import java.util.Map;
import java.util.concurrent.TimeUnit;

import org.apache.nifi.annotation.behavior.InputRequirement;
import org.apache.nifi.annotation.behavior.ReadsAttribute;
import org.apache.nifi.annotation.behavior.ReadsAttributes;
import org.apache.nifi.annotation.behavior.TriggerSerially;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.SeeAlso;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.model.PlcTag;

@TriggerSerially
@Tags({"plc4x", "put", "sink"})
@SeeAlso({Plc4xSourceProcessor.class})
@InputRequirement(InputRequirement.Requirement.INPUT_REQUIRED)
@CapabilityDescription("Processor able to write data to industrial PLCs using Apache PLC4X")
@ReadsAttributes({@ReadsAttribute(attribute="value", description="some value")})
public class Plc4xSinkProcessor extends BasePlc4xProcessor {

    @Override
    public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
        FlowFile flowFile = session.get();
        final ComponentLog logger = getLogger();

        // Abort if there's nothing to do.
        if (flowFile == null) {
            return;
        }

        // Get an instance of a component able to write to a PLC.
        try(PlcConnection connection = getConnectionManager().getConnection(getConnectionString())) {
            if (!connection.getMetadata().canWrite()) {
                throw new ProcessException("Writing not supported by connection");
            }

            // Prepare the request.
            PlcWriteRequest.Builder builder = connection.writeRequestBuilder();
            Map<String,String> addressMap = getPlcAddressMap(context, flowFile);
            final Map<String, PlcTag> tags = getSchemaCache().retrieveTags(addressMap);

            if (tags != null){
                for (Map.Entry<String,PlcTag> tag : tags.entrySet()){
                    if (flowFile.getAttributes().containsKey(tag.getKey())) {
                        builder.addTag(tag.getKey(), tag.getValue(), flowFile.getAttribute(tag.getKey()));
                    } else {
                        if (debugEnabled)
                            logger.debug("PlcTag " + tag + " is declared as address but was not found on input record.");
                    }
                }
            } else {
                for (Map.Entry<String,String> entry: addressMap.entrySet()){
                    if (flowFile.getAttributes().containsKey(entry.getKey())) {
                        builder.addTagAddress(entry.getKey(), entry.getValue(), flowFile.getAttribute(entry.getKey()));
                    }
                }
                if (debugEnabled)
                    logger.debug("PlcTypes resolution not found in cache and will be added with key: " + addressMap);
            }
           
            PlcWriteRequest writeRequest = builder.build();

            // Send the request to the PLC.
            try {
                final PlcWriteResponse plcWriteResponse = writeRequest.execute().get(this.timeout, TimeUnit.MILLISECONDS);
                PlcResponseCode code = null;

                for (String tag : plcWriteResponse.getTagNames()) {
                    code = plcWriteResponse.getResponseCode(tag);
                    if (!code.equals(PlcResponseCode.OK)) {
                        logger.error("Not OK code when writing the data to PLC for tag " + tag 
								+ " with value  " + flowFile.getAttribute(tag)
								+ " in addresss " + plcWriteResponse.getTag(tag).getAddressString());
                        throw new Exception(code.toString());
                    }
                }
                session.transfer(flowFile, REL_SUCCESS);

                if (tags == null){
                    if (debugEnabled)
                        logger.debug("Adding PlcTypes resolution into cache with key: " + addressMap);
                    getSchemaCache().addSchema(
                        addressMap, 
                        writeRequest.getTagNames(),
                        writeRequest.getTags(),
                        null
                    );
                }
            } catch (Exception e) {
                flowFile = session.putAttribute(flowFile, "exception", e.getLocalizedMessage());
                session.transfer(flowFile, REL_FAILURE);
            }

        } catch (ProcessException e) {
            throw e;
        } catch (Exception e) {
            throw new ProcessException("Got an error while trying to get a connection", e);
        }
    }

}
