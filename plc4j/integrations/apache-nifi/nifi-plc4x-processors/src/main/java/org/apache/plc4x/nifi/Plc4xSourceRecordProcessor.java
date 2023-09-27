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

import java.io.OutputStream;
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
import org.apache.nifi.serialization.RecordSetWriterFactory;
import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.nifi.util.StopWatch;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.nifi.record.Plc4xWriter;
import org.apache.plc4x.nifi.record.RecordPlc4xWriter;

@Tags({"plc4x", "get", "input", "source", "record"})
@SeeAlso({Plc4xSinkRecordProcessor.class, Plc4xListenRecordProcessor.class})
@InputRequirement(InputRequirement.Requirement.INPUT_ALLOWED)
@CapabilityDescription("Processor able to read data from industrial PLCs using Apache PLC4X")
@WritesAttributes({ 
	@WritesAttribute(attribute = Plc4xSourceRecordProcessor.RESULT_ROW_COUNT, description = "Number of rows written into the output FlowFile"),
	@WritesAttribute(attribute = Plc4xSourceRecordProcessor.RESULT_QUERY_EXECUTION_TIME, description = "Time between request and response from the PLC"),
	@WritesAttribute(attribute = Plc4xSourceRecordProcessor.INPUT_FLOWFILE_UUID, description = "UUID of the input FlowFile")
 })
public class Plc4xSourceRecordProcessor extends BasePlc4xProcessor {

	public static final String RESULT_ROW_COUNT = "plc4x.read.row.count";
	public static final String RESULT_QUERY_EXECUTION_TIME = "plc4x.read.query.executiontime";
	public static final String INPUT_FLOWFILE_UUID = "input.flowfile.uuid";

	public static final PropertyDescriptor PLC_RECORD_WRITER_FACTORY = new PropertyDescriptor.Builder()
		.name("plc4x-record-writer").displayName("Record Writer")
		.description("Specifies the Controller Service to use for writing results to a FlowFile. The Record Writer may use Inherit Schema to emulate the inferred schema behavior, i.e. "
				+ "an explicit schema need not be defined in the writer, and will be supplied by the same logic used to infer the schema from the column types.")
		.identifiesControllerService(RecordSetWriterFactory.class)
		.required(true)
		.build();
	
	@Override
	protected void init(final ProcessorInitializationContext context) {
		super.init(context);
        final Set<Relationship> r = new HashSet<>(super.getRelationships());
		this.relationships = Collections.unmodifiableSet(r);

        final List<PropertyDescriptor> pds = new ArrayList<>(super.getSupportedPropertyDescriptors());
		pds.add(PLC_RECORD_WRITER_FACTORY);
		this.properties = Collections.unmodifiableList(pds);
	}

	
	@Override
	public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
		
		FlowFile fileToProcess = null;
		if (context.hasIncomingConnection()) {
			fileToProcess = session.get();
			
			if (fileToProcess == null && context.hasNonLoopConnection()) {
				return;
			}
		}

		final ComponentLog logger = getLogger();
		
		
		// Get an instance of a component able to read from a PLC.
		final AtomicLong nrOfRows = new AtomicLong(0L);
		final StopWatch executeTime = new StopWatch(true);

		String inputFileUUID = fileToProcess == null ? null : fileToProcess.getAttribute(CoreAttributes.UUID.key());
		final FlowFile resultSetFF = createResultFlowFile(session, fileToProcess);

		final FlowFile originalFlowFile = fileToProcess;

		Plc4xWriter plc4xWriter = new RecordPlc4xWriter(context.getProperty(PLC_RECORD_WRITER_FACTORY).asControllerService(RecordSetWriterFactory.class), 
			fileToProcess == null ? Collections.emptyMap() : fileToProcess.getAttributes());


		try {
			session.write(resultSetFF, out -> {

			

				final Map<String,String> addressMap = getPlcAddressMap(context, originalFlowFile);
				final RecordSchema recordSchema = getSchemaCache().retrieveSchema(addressMap);
				final Map<String, PlcTag> tags = getSchemaCache().retrieveTags(addressMap);
				PlcReadRequest readRequest;
				Long nrOfRowsHere;

				try (PlcConnection connection = getConnectionManager().getConnection(getConnectionString(context, originalFlowFile))) {
					
					readRequest =  getReadRequest(logger, addressMap, tags, connection);
					
					PlcReadResponse readResponse = readRequest.execute().get(getTimeout(context, originalFlowFile), TimeUnit.MILLISECONDS);
							
					nrOfRowsHere = evaluateReadResponse(context, logger, originalFlowFile, plc4xWriter, out, recordSchema, readResponse);

				} catch (TimeoutException e) {
					logger.error("Timeout reading the data from PLC", e);
					getConnectionManager().removeCachedConnection(getConnectionString(context, originalFlowFile));
					throw new ProcessException(e);
				} catch (PlcConnectionException e) {
					logger.error("Error getting the PLC connection", e);
					throw new ProcessException("Got an a PlcConnectionException while trying to get a connection", e);
				} catch (Exception e) {
					logger.error("Exception reading the data from PLC", e);
					throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
				}

				if (recordSchema == null){
					addTagsToCache(logger, addressMap, readRequest, plc4xWriter.getRecordSchema());
				}
				nrOfRows.set(nrOfRowsHere);

			});
			
		} catch (Exception e) {
			logger.error("Exception reading the data from the PLC", e);
			if (fileToProcess != null) {
				session.putAttribute(fileToProcess, "exception", e.getLocalizedMessage());
				session.transfer(fileToProcess, REL_FAILURE);
			}
			session.remove(resultSetFF);
			session.commitAsync();
			throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
		}

		completeResultFlowFile(context, session, logger, nrOfRows, executeTime, inputFileUUID, resultSetFF,
				plc4xWriter);
		
		if (fileToProcess != null) {
			session.remove(fileToProcess);
		}
		session.transfer(resultSetFF, REL_SUCCESS);
	}

	private FlowFile createResultFlowFile(final ProcessSession session, FlowFile fileToProcess) {

		Map<String, String> inputFileAttrMap = fileToProcess == null ? null : fileToProcess.getAttributes();
		
		FlowFile resultSetFF;
		if (fileToProcess == null) {
			resultSetFF = session.create();
		} else {
			resultSetFF = session.create(fileToProcess);
		}
		if (inputFileAttrMap != null) {
			resultSetFF = session.putAllAttributes(resultSetFF, inputFileAttrMap);
		}
		return resultSetFF;
	}

	private PlcReadRequest getReadRequest(final ComponentLog logger, final Map<String, String> addressMap,
			final Map<String, PlcTag> tags, final PlcConnection connection) {

		PlcReadRequest.Builder builder = connection.readRequestBuilder();

		if (tags != null){
			for (Map.Entry<String,PlcTag> tag : tags.entrySet()){
				builder.addTag(tag.getKey(), tag.getValue());
			}
		} else {
			if (debugEnabled)
				logger.debug("Plc-Avro schema and PlcTypes resolution not found in cache and will be added with key: " + addressMap);
			for (Map.Entry<String,String> entry: addressMap.entrySet()){
				builder.addTagAddress(entry.getKey(), entry.getValue());
			}
		}
		return builder.build();
	}


	private long evaluateReadResponse(final ProcessContext context, final ComponentLog logger, final FlowFile originalFlowFile,
			Plc4xWriter plc4xWriter, OutputStream out, final RecordSchema recordSchema, PlcReadResponse readResponse)
			throws Exception {

		if(originalFlowFile == null) //there is no inherit attributes to use in writer service 
			return plc4xWriter.writePlcReadResponse(readResponse, out, logger, null, recordSchema, getTimestampField(context));
		else 
			return plc4xWriter.writePlcReadResponse(readResponse, out, logger, null, recordSchema, originalFlowFile, getTimestampField(context));
	}

	private void addTagsToCache(final ComponentLog logger, final Map<String, String> addressMap,
			PlcReadRequest readRequest, RecordSchema schema) {
		if (debugEnabled)
			logger.debug("Adding Plc-Avro schema and PlcTypes resolution into cache with key: " + addressMap);
		getSchemaCache().addSchema(
			addressMap, 
			readRequest.getTagNames(),
			readRequest.getTags(),
			schema
		);
	}

	private void completeResultFlowFile(final ProcessContext context, final ProcessSession session,
			final ComponentLog logger, final AtomicLong nrOfRows, final StopWatch executeTime, String inputFileUUID,
			FlowFile resultSetFF, Plc4xWriter plc4xWriter) {
		
		long executionTimeElapsed = executeTime.getElapsed(TimeUnit.MILLISECONDS);
		final Map<String, String> attributesToAdd = new HashMap<>();
		attributesToAdd.put(RESULT_ROW_COUNT, String.valueOf(nrOfRows.get()));
		attributesToAdd.put(RESULT_QUERY_EXECUTION_TIME, String.valueOf(executionTimeElapsed));
		if (inputFileUUID != null) {
			attributesToAdd.put(INPUT_FLOWFILE_UUID, inputFileUUID);
		}
		attributesToAdd.putAll(plc4xWriter.getAttributesToAdd());

		resultSetFF = session.putAllAttributes(resultSetFF, attributesToAdd);
		plc4xWriter.updateCounters(session);
		
		logger.info("{} contains {} records; transferring to 'success'", resultSetFF, nrOfRows.get());

		if (context.hasIncomingConnection()) {
			session.getProvenanceReporter().fetch(resultSetFF, "Retrieved " + nrOfRows.get() + " rows", executionTimeElapsed);
		} else {
			session.getProvenanceReporter().receive(resultSetFF, "Retrieved " + nrOfRows.get() + " rows", executionTimeElapsed);
		}
	}

}
