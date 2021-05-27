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

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;

import org.apache.nifi.annotation.behavior.InputRequirement;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.components.PropertyValue;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.flowfile.attributes.CoreAttributes;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.serialization.RecordSetWriterFactory;
import org.apache.nifi.util.StopWatch;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.nifi.record.Plc4xWriter;
import org.apache.plc4x.nifi.record.RecordPlc4xWriter;

@Tags({ "plc4x-source" })
@InputRequirement(InputRequirement.Requirement.INPUT_ALLOWED)
@CapabilityDescription("Processor able to read data from industrial PLCs using Apache PLC4X")
@WritesAttributes({ @WritesAttribute(attribute = "value", description = "some value") })
public class Plc4xSourceRecordProcessor extends BasePlc4xProcessor {

	public static final String RESULT_ROW_COUNT = "plc4x.read.row.count";
	public static final String RESULT_QUERY_DURATION = "plc4x.read.query.duration";
	public static final String RESULT_QUERY_EXECUTION_TIME = "plc4x.read.query.executiontime";
	public static final String RESULT_QUERY_FETCH_TIME = "plc4x.read.query.fetchtime";
	public static final String INPUT_FLOWFILE_UUID = "input.flowfile.uuid";
	public static final String RESULT_ERROR_MESSAGE = "plc4x.read.error.message";

	public static final PropertyDescriptor RECORD_WRITER_FACTORY = new PropertyDescriptor.Builder().name("plc4x-record-writer").displayName("Record Writer")
			.description("Specifies the Controller Service to use for writing results to a FlowFile. The Record Writer may use Inherit Schema to emulate the inferred schema behavior, i.e. "
					+ "an explicit schema need not be defined in the writer, and will be supplied by the same logic used to infer the schema from the column types.")
			.identifiesControllerService(RecordSetWriterFactory.class)
			.required(true)
			.build();

	public Plc4xSourceRecordProcessor() {

	}

	@Override
	protected void init(final ProcessorInitializationContext context) {
		super.init(context);
		final Set<Relationship> r = new HashSet<>();
		r.addAll(super.getRelationships());
		this.relationships = Collections.unmodifiableSet(r);

		final List<PropertyDescriptor> pds = new ArrayList<>();
		pds.addAll(super.getSupportedPropertyDescriptors());
		pds.add(RECORD_WRITER_FACTORY);
		this.properties = Collections.unmodifiableList(pds);
	}

	@Override 
	public void onScheduled(final ProcessContext context) {
        connectionString = context.getProperty(PLC_CONNECTION_STRING.getName()).getValue();
    }
	
	@Override
	public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
		FlowFile fileToProcess = null;
		// TODO: In the future the processor will be configurable to get the address and
		// the connection from incoming flowfile
		if (context.hasIncomingConnection()) {
			fileToProcess = session.get();
			// If we have no FlowFile, and all incoming connections are self-loops then we
			// can continue on.
			// However, if we have no FlowFile and we have connections coming from other
			// Processors, then we know that we should run only if we have a FlowFile.
			if (fileToProcess == null && context.hasNonLoopConnection()) {
				return;
			}
		}

		// TODO: this could be enhanced checking if address map should be updated (via a cache boolean, checking property values is a nifi expression language, etc)
		Map<String, String> addressMap = new HashMap<>();
        PropertyValue addresses = context.getProperty(PLC_ADDRESS_STRING.getName());
        for (String segment : addresses.evaluateAttributeExpressions(fileToProcess).getValue().split(";")) {
            String[] parts = segment.split("=");
            if(parts.length != 2) {
                throw new ProcessException("Invalid address format");
            }
            addressMap.put(parts[0], parts[1]);
        }
		
		
		final List<FlowFile> resultSetFlowFiles = new ArrayList<>();

		Plc4xWriter plc4xWriter = new RecordPlc4xWriter(context.getProperty(RECORD_WRITER_FACTORY).asControllerService(RecordSetWriterFactory.class), fileToProcess == null ? Collections.emptyMap() : fileToProcess.getAttributes());
		final ComponentLog logger = getLogger();
		// Get an instance of a component able to read from a PLC.
		// TODO: Change this to use NiFi service instead of direct connection
		final AtomicLong nrOfRows = new AtomicLong(0L);
		final StopWatch executeTime = new StopWatch(true);

		try (PlcConnection connection = getDriverManager().getConnection(getConnectionString())) {

			String inputFileUUID = fileToProcess == null ? null : fileToProcess.getAttribute(CoreAttributes.UUID.key());
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

			PlcReadRequest.Builder builder = connection.readRequestBuilder();
			getFields().forEach(field -> {
				String address = getAddress(field);
				if (address != null) {
					builder.addItem(field, address);
				}
			});
			PlcReadRequest readRequest = builder.build();
			PlcReadResponse readResponse = readRequest.execute().get(10, TimeUnit.SECONDS);
			resultSetFF = session.write(resultSetFF, out -> {
				try {
					nrOfRows.set(plc4xWriter.writePlcReadResponse(readResponse, out, logger, null));
				} catch (Exception e) {
					throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
				}
			});
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
			logger.info("{} contains {} records; transferring to 'success'", new Object[] { resultSetFF, nrOfRows.get() });
			// Report a FETCH event if there was an incoming flow file, or a RECEIVE event
			// otherwise
			if (context.hasIncomingConnection()) {
				session.getProvenanceReporter().fetch(resultSetFF, "Retrieved " + nrOfRows.get() + " rows", executionTimeElapsed);
			} else {
				session.getProvenanceReporter().receive(resultSetFF, "Retrieved " + nrOfRows.get() + " rows", executionTimeElapsed);
			}
			resultSetFlowFiles.add(resultSetFF);
			if (resultSetFlowFiles.size() >= 0) {
				session.transfer(resultSetFlowFiles, BasePlc4xProcessor.REL_SUCCESS);
				// Need to remove the original input file if it exists
				if (fileToProcess != null) {
					session.remove(fileToProcess);
					fileToProcess = null;
				}
				session.commit();
				resultSetFlowFiles.clear();
			}

		} catch (Exception e) {
			logger.error("Got an error while trying to get a connection", e);
			throw new ProcessException("Got an error while trying to get a connection", e);
		}
	}

}
