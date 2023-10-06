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
import java.util.concurrent.TimeoutException;

import org.apache.nifi.annotation.behavior.InputRequirement;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.SeeAlso;
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
@SeeAlso({Plc4xSinkProcessor.class})
@InputRequirement(InputRequirement.Requirement.INPUT_FORBIDDEN)
@CapabilityDescription("Processor able to read data from industrial PLCs using Apache PLC4X")
@WritesAttributes({@WritesAttribute(attribute="value", description="some value")})
public class Plc4xSourceProcessor extends BasePlc4xProcessor {

	public static final String EXCEPTION = "plc4x.read.exception";

    @Override
    public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
        
        FlowFile incomingFlowFile = null;
        if (context.hasIncomingConnection()) {
            incomingFlowFile = session.get();
            if (incomingFlowFile == null && context.hasNonLoopConnection()) {
                return;
            }
        }

        final ComponentLog logger = getLogger();
        final FlowFile flowFile = session.create();
    
        try(PlcConnection connection = getConnectionManager().getConnection(getConnectionString(context, incomingFlowFile))) {

            if (!connection.getMetadata().canRead()) {
                throw new ProcessException("Reading not supported by connection");
            }

            final Map<String,String> addressMap = getPlcAddressMap(context, incomingFlowFile);
            final Map<String, PlcTag> tags = getSchemaCache().retrieveTags(addressMap);


            PlcReadRequest readRequest = getReadRequest(logger, addressMap, tags, connection);

            try {
                final PlcReadResponse response = readRequest.execute().get(getTimeout(context, incomingFlowFile), TimeUnit.MILLISECONDS);
                
                evaluateReadResponse(session, flowFile, response);
                
            } catch (TimeoutException e) {
                logger.error("Timeout reading the data from PLC", e);
                getConnectionManager().removeCachedConnection(getConnectionString(context, incomingFlowFile));
                throw new ProcessException(e);
            } catch (Exception e) {
                logger.error("Exception reading the data from PLC", e);
                throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
            }

            
            if (incomingFlowFile != null) {
                session.remove(incomingFlowFile);
            }
            session.transfer(flowFile, REL_SUCCESS);
                
            if (tags == null){
                if (debugEnabled)
                    logger.debug("Adding PlcTypes resolution into cache with key: " + addressMap);
                getSchemaCache().addSchema(
                    addressMap, 
                    readRequest.getTagNames(),
                    readRequest.getTags(),
                    null
                );
            }
            
        } catch (Exception e) {
            session.remove(flowFile);
            if (incomingFlowFile != null){
                incomingFlowFile = session.putAttribute(incomingFlowFile, EXCEPTION, e.getLocalizedMessage());
                session.transfer(incomingFlowFile, REL_FAILURE);
            }
            session.commitAsync();
            throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
        }
    }
    
}
