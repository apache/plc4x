/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.nifi;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicLong;

import org.apache.nifi.annotation.behavior.InputRequirement;
import org.apache.nifi.annotation.behavior.ReadsAttribute;
import org.apache.nifi.annotation.behavior.ReadsAttributes;
import org.apache.nifi.annotation.behavior.TriggerSerially;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.SeeAlso;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.flowfile.attributes.CoreAttributes;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.serialization.RecordReader;
import org.apache.nifi.serialization.RecordReaderFactory;
import org.apache.nifi.serialization.record.Record;
import org.apache.nifi.util.StopWatch;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;

@TriggerSerially
@Tags({"plc4x", "put", "sink", "record"})
@SeeAlso({Plc4xSourceRecordProcessor.class, Plc4xListenRecordProcessor.class})
@InputRequirement(InputRequirement.Requirement.INPUT_REQUIRED)
@CapabilityDescription("Processor able to write data to industrial PLCs using Apache PLC4X")
@ReadsAttributes({@ReadsAttribute(attribute="value", description="some value")})
@WritesAttributes({ 
	@WritesAttribute(attribute = Plc4xSinkRecordProcessor.RESULT_ROW_COUNT, description = "Number of rows from the input FlowFile written into the PLC"),
	@WritesAttribute(attribute = Plc4xSinkRecordProcessor.RESULT_QUERY_EXECUTION_TIME, description = "Time between request and response from the PLC"),
	@WritesAttribute(attribute = Plc4xSinkRecordProcessor.INPUT_FLOWFILE_UUID, description = "UUID of the input FlowFile")
 })
public class Plc4xSinkRecordProcessor extends BasePlc4xProcessor {

	public static final String RESULT_ROW_COUNT = "plc4x.write.row.count";
	public static final String RESULT_QUERY_EXECUTION_TIME = "plc4x.write.query.executiontime";
	public static final String INPUT_FLOWFILE_UUID = "input.flowfile.uuid";
	public static final String EXCEPTION = "plc4x.write.exception";
	
	public static final PropertyDescriptor PLC_RECORD_READER_FACTORY = new PropertyDescriptor.Builder()
			.name("record-reader").displayName("Record Reader")
			.description(
					"Specifies the Controller Service to use for reading record from a FlowFile. The Record Reader may use Inherit Schema to emulate the inferred schema behavior, i.e. "
							+ "an explicit schema need not be defined in the reader, and will be supplied by the same logic used to infer the schema from the column types.")
			.identifiesControllerService(RecordReaderFactory.class)
			.required(true)
			.build();


	@Override
	protected void init(final ProcessorInitializationContext context) {
		super.init(context);
        final Set<Relationship> r = new HashSet<>(super.getRelationships());
		this.relationships = Collections.unmodifiableSet(r);

        final List<PropertyDescriptor> pds = new ArrayList<>(super.getSupportedPropertyDescriptors());
		pds.add(PLC_RECORD_READER_FACTORY);
		this.properties = Collections.unmodifiableList(pds);
	}

	@Override
	public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
		FlowFile fileToProcess = session.get();

        // Abort if there's nothing to do.
        if (fileToProcess == null) {
            return;
        }
			
		final ComponentLog logger = getLogger();

		// Get an instance of a component able to read from a PLC.
		final AtomicLong nrOfRows = new AtomicLong(0L);
		final StopWatch executeTime = new StopWatch(true);

		try {
			session.read(fileToProcess, in -> {
				Record record = null;
			
				try (RecordReader recordReader = context.getProperty(PLC_RECORD_READER_FACTORY)
					.asControllerService(RecordReaderFactory.class)
					.createRecordReader(fileToProcess, in, logger)){

					while ((record = recordReader.nextRecord()) != null) {
						AtomicLong nrOfRowsHere = new AtomicLong(0);
						PlcWriteRequest writeRequest;

						final Map<String,String> addressMap = getPlcAddressMap(context, fileToProcess);
						final Map<String, PlcTag> tags = getSchemaCache().retrieveTags(addressMap);

						try (PlcConnection connection = getConnectionManager().getConnection(getConnectionString(context, fileToProcess))) {
							
							writeRequest = getWriteRequest(logger, addressMap, tags, record.toMap(), connection, nrOfRowsHere);

							PlcWriteResponse plcWriteResponse = writeRequest.execute().get(getTimeout(context, fileToProcess), TimeUnit.MILLISECONDS);

							// Response check if values were written
							evaluateWriteResponse(logger, record.toMap(), plcWriteResponse);

						} catch (TimeoutException e) {
							logger.error("Timeout writting the data to the PLC", e);
							getConnectionManager().removeCachedConnection(getConnectionString(context, fileToProcess));
							throw new ProcessException(e);
						} catch (PlcConnectionException e) {
							logger.error("Error getting the PLC connection", e);
							throw new ProcessException("Got an a PlcConnectionException while trying to get a connection", e);
						} catch (Exception e) {
							logger.error("Exception writting the data to the PLC", e);
							throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
						}
							
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
						nrOfRows.getAndAdd(nrOfRowsHere.get());

						
					}
				} catch (Exception e) {
					throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
				} 
			});

		} catch (ProcessException e) {
			logger.error("Exception writing the data to the PLC", e);
			session.putAttribute(fileToProcess, EXCEPTION, e.getLocalizedMessage());
			session.transfer(fileToProcess, REL_FAILURE);
			session.commitAsync();
			throw e;
		} 


		long executionTimeElapsed = executeTime.getElapsed(TimeUnit.MILLISECONDS);
		final Map<String, String> attributesToAdd = new HashMap<>();
		attributesToAdd.put(RESULT_ROW_COUNT, String.valueOf(nrOfRows.get()));
		attributesToAdd.put(RESULT_QUERY_EXECUTION_TIME, String.valueOf(executionTimeElapsed));
		attributesToAdd.put(INPUT_FLOWFILE_UUID, fileToProcess.getAttribute(CoreAttributes.UUID.key()));
		
		session.putAllAttributes(fileToProcess, attributesToAdd);

		session.transfer(fileToProcess, REL_SUCCESS);
		
		logger.info("Writing {} fields from {} records; transferring to 'success'", nrOfRows.get(), fileToProcess);
		if (context.hasIncomingConnection()) {
			session.getProvenanceReporter().fetch(fileToProcess, "Writted " + nrOfRows.get() + " rows", executionTimeElapsed);
		} else {
			session.getProvenanceReporter().receive(fileToProcess, "Writted " + nrOfRows.get() + " rows", executionTimeElapsed);
		}
	}
}
