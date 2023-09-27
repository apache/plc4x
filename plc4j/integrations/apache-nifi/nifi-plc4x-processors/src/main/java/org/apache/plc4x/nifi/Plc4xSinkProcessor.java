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

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

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

	public static final String EXCEPTION = "plc4x.write.exception";

    @Override
    public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
        FlowFile flowFile = session.get();
        
        // Abort if there's nothing to do.
        if (flowFile == null) {
            return;
        }

        final ComponentLog logger = getLogger();

        try(PlcConnection connection = getConnectionManager().getConnection(getConnectionString(context, flowFile))) {
            if (!connection.getMetadata().canWrite()) {
                throw new ProcessException("Writing not supported by connection");
            }
            
            final Map<String,String> addressMap = getPlcAddressMap(context, flowFile);
            final Map<String, PlcTag> tags = getSchemaCache().retrieveTags(addressMap);

            PlcWriteRequest writeRequest = getWriteRequest(flowFile, logger, addressMap, tags, connection);

            try {
                final PlcWriteResponse plcWriteResponse = writeRequest.execute().get(getTimeout(context, flowFile), TimeUnit.MILLISECONDS);

                evaluateWriteResponse(flowFile, logger, plcWriteResponse);
 
            } catch (TimeoutException e) {
                logger.error("Timeout writting the data to the PLC", e);
                getConnectionManager().removeCachedConnection(getConnectionString(context, flowFile));
                throw new ProcessException(e);
            } catch (Exception e) {
                logger.error("Exception writting the data to the PLC", e);
                throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
            }

            session.transfer(flowFile, REL_SUCCESS);

            if (tags == null){
                addTagsToCache(logger, addressMap, writeRequest);
            }


        } catch (Exception e) {
            flowFile = session.putAttribute(flowFile, EXCEPTION, e.getLocalizedMessage());
            session.transfer(flowFile, REL_FAILURE);
            session.commitAsync();
            throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
        }
    }

    private PlcWriteRequest getWriteRequest(FlowFile flowFile, final ComponentLog logger,
            final Map<String, String> addressMap, final Map<String, PlcTag> tags, final PlcConnection connection) {

        PlcWriteRequest.Builder builder = connection.writeRequestBuilder();

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
            if (debugEnabled)
                logger.debug("PlcTypes resolution not found in cache and will be added with key: " + addressMap);
            for (Map.Entry<String,String> entry: addressMap.entrySet()){
                if (flowFile.getAttributes().containsKey(entry.getKey())) {
                    builder.addTagAddress(entry.getKey(), entry.getValue(), flowFile.getAttribute(entry.getKey()));
                }
            }
        }
         
        return builder.build();
    }

    private void evaluateWriteResponse(FlowFile flowFile, final ComponentLog logger, final PlcWriteResponse plcWriteResponse) {
        boolean codeErrorPresent = false;
        List<String> tagsAtError = null;
        for (String tag : plcWriteResponse.getTagNames()) {

            PlcResponseCode code = plcWriteResponse.getResponseCode(tag);
            if (!code.equals(PlcResponseCode.OK)) {
                if (tagsAtError == null) {
                    tagsAtError = new ArrayList<>();
                }
                logger.error("Not OK code when writing the data to PLC for tag " + tag 
        				+ " with value  " + flowFile.getAttribute(tag)
        				+ " in addresss " + plcWriteResponse.getTag(tag).getAddressString());
                codeErrorPresent = true;
                tagsAtError.add(tag);
            }
        }
        if (codeErrorPresent) {
            throw new ProcessException("At least one error was found when while writting tags: " + tagsAtError.toString());
        }
    }

    private void addTagsToCache(final ComponentLog logger, final Map<String, String> addressMap,
            PlcWriteRequest writeRequest) {
        if (debugEnabled)
            logger.debug("Adding PlcTypes resolution into cache with key: " + addressMap);
        getSchemaCache().addSchema(
            addressMap, 
            writeRequest.getTagNames(),
            writeRequest.getTags(),
            null
        );
    }

}
