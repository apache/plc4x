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
import org.apache.nifi.annotation.lifecycle.OnScheduled;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.flowfile.attributes.CoreAttributes;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.processor.util.StandardValidators;
import org.apache.nifi.serialization.RecordSetWriterFactory;
import org.apache.nifi.util.StopWatch;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.apache.plc4x.nifi.record.Plc4xWriter;
import org.apache.plc4x.nifi.record.RecordPlc4xWriter;
import org.apache.plc4x.nifi.util.PLC4X_PROTOCOL;
import org.apache.plc4x.nifi.util.Plc4xCommon;

@Tags({ "plc4x-source" })
@InputRequirement(InputRequirement.Requirement.INPUT_FORBIDDEN)
@CapabilityDescription("Processor able to read data from industrial PLCs using Apache PLC4X")
@WritesAttributes({ @WritesAttribute(attribute = "value", description = "some value") })
public class Plc4xSourceRecordProcessor extends BasePlc4xProcessor {

	public static final String RESULT_ROW_COUNT = "plc4x.read.row.count";
	public static final String RESULT_QUERY_DURATION = "plc4x.read.query.duration";
	public static final String RESULT_QUERY_EXECUTION_TIME = "plc4x.read.query.executiontime";
	public static final String RESULT_QUERY_FETCH_TIME = "plc4x.read.query.fetchtime";
	public static final String INPUT_FLOWFILE_UUID = "input.flowfile.uuid";
	public static final String RESULT_ERROR_MESSAGE = "plc4x.read.error.message";

	public static final PropertyDescriptor RECORD_WRITER_FACTORY = new PropertyDescriptor.Builder()
			.name("plc4x-record-writer").displayName("Record Writer")
			.description("Specifies the Controller Service to use for writing results to a FlowFile. The Record Writer may use Inherit Schema to emulate the inferred schema behavior, i.e. "
							+ "an explicit schema need not be defined in the writer, and will be supplied by the same logic used to infer the schema from the column types.")
			.identifiesControllerService(RecordSetWriterFactory.class).required(true).build();
	
	public static final PropertyDescriptor FORCE_RECONNECT = new PropertyDescriptor.Builder()
			.name("plc4x-reconnect-force").displayName("Force Reconnect every request")
			.description("Specifies if the connection to plc will be recreated on trigger event")
			.required(true).addValidator(StandardValidators.BOOLEAN_VALIDATOR).defaultValue("true").build();

	public Plc4xSourceRecordProcessor() {

	}

	private PlcConnection connection = null;
	private PooledPlcDriverManager pool;
	
	private static PLC4X_PROTOCOL PROTOCOL = null;
	
	@Override
	@OnScheduled
    public void onScheduled(final ProcessContext context) {
        super.onScheduled(context);
        //TODO: Change this to use NiFi service instead of direct connection and add @OnStopped Phase to close connection
        try {
        	Boolean force =  context.getProperty(FORCE_RECONNECT).asBoolean();
        	pool = new PooledPlcDriverManager();
        	if(!force) {
        		this.connection = pool.getConnection(getConnectionString());
        	}
			//TODO how to infer protocol within the writer?
			PROTOCOL = Plc4xCommon.getConnectionProtocol(getConnectionString());
		} catch (PlcConnectionException e) {
			if(this.connection != null)
				try {
					this.connection.close();
				} catch (Exception e1) {
					//do nothing
				}
			getLogger().error("Error while creating the connection to "+getConnectionString(), e);
		}
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
		pds.add(FORCE_RECONNECT);
		this.descriptors = Collections.unmodifiableList(pds);
	}

	@Override
	public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {
		FlowFile fileToProcess = null;
		//TODO: In the future the processor will be configurable to get the address and the connection from incoming flowfile
		if (context.hasIncomingConnection()) {
			fileToProcess = session.get();
			// If we have no FlowFile, and all incoming connections are self-loops then we can continue on.
			// However, if we have no FlowFile and we have connections coming from other
			// Processors, then we know that we should run only if we have a FlowFile.
			if (fileToProcess == null && context.hasNonLoopConnection()) {
				return;
			}
		}

		final List<FlowFile> resultSetFlowFiles = new ArrayList<>();

		Plc4xWriter plc4xWriter = new RecordPlc4xWriter(context.getProperty(RECORD_WRITER_FACTORY).asControllerService(RecordSetWriterFactory.class), fileToProcess == null ? Collections.emptyMap() : fileToProcess.getAttributes());
		Boolean force =  context.getProperty(FORCE_RECONNECT).asBoolean();
		final ComponentLog logger = getLogger();
		logger.info("Te connection {} will be recreated? (Force Reconnect every request) is {}",	new Object[] { getConnectionString(), force});
		// Get an instance of a component able to read from a PLC.
		// TODO: Change this to use NiFi service instead of direct connection
		final AtomicLong nrOfRows = new AtomicLong(0L);
		final StopWatch executeTime = new StopWatch(true);
		try {
			if(force) {
				logger.debug("Recreating the connection {} because the parameter (Force Reconnect every request) is {}",	new Object[] { getConnectionString(), force});
				this.connection = pool.getConnection(getConnectionString());
			}
			if(this.connection != null) {
				// Prepare the request.
				if (!connection.getMetadata().canRead()) {
					throw new ProcessException("Reading not supported by connection");
				}
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
				try {
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
							nrOfRows.set(plc4xWriter.writePlcReadResponse(readResponse, this.getPlcAddress(), out, logger, null, PROTOCOL));
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
					logger.info("{} contains {} records; transferring to 'success'",
							new Object[] { resultSetFF, nrOfRows.get() });
					// Report a FETCH event if there was an incoming flow file, or a RECEIVE event otherwise
					if (context.hasIncomingConnection()) {
						session.getProvenanceReporter().fetch(resultSetFF, "Retrieved " + nrOfRows.get() + " rows", executionTimeElapsed);
					} else {
						session.getProvenanceReporter().receive(resultSetFF, "Retrieved " + nrOfRows.get() + " rows", executionTimeElapsed);
					}
					resultSetFlowFiles.add(resultSetFF);
					if (resultSetFlowFiles.size() >= 0) {
						session.transfer(resultSetFlowFiles, super.REL_SUCCESS);
						// Need to remove the original input file if it exists
						if (fileToProcess != null) {
							session.remove(fileToProcess);
							fileToProcess = null;
						}
						session.commit();
						resultSetFlowFiles.clear();
					}
				} catch (Exception e) {
					//if there is any error, recreate the pool because it may not recreate correctly if the error ocurred on the execute.get() 
					pool = new PooledPlcDriverManager();
					if (e instanceof InterruptedException)
						Thread.currentThread().interrupt();
					session.remove(resultSetFF);
					if (e instanceof ProcessException) {
						throw (ProcessException) e;
					} else {
						throw new ProcessException(e);
					}
				}
			} else {
				throw new ProcessException("Connection is null");
			}
		} catch (Exception e) {
			//if there is any error, recreate the pool because it may not recreate correctly if the error ocurred on the execute.get() 
			pool = new PooledPlcDriverManager();
			if (fileToProcess == null) {
				// This can happen if any exceptions occur while setting up the connection,
				// statement, etc.
				logger.error("Unable to execute PLC4X query {} due to {}. No FlowFile to route to failure", new Object[] { getConnectionString(), e });
				context.yield();
			} else {
				if (context.hasIncomingConnection()) {
					logger.error("Unable to execute processor select query {} for {} due to {}; routing to failure", new Object[] { getConnectionString(), fileToProcess, e });
					fileToProcess = session.penalize(fileToProcess);
				} else {
					logger.error("Unable to execute processor select query {} due to {}; routing to failure",	new Object[] { getConnectionString(), e });
					context.yield();
				}
				session.putAttribute(fileToProcess, RESULT_ERROR_MESSAGE, e.getMessage());
				session.transfer(fileToProcess, REL_FAILURE);
			}
		} finally {
			logger.info("Force the close"); //TODO log level, antes estaba a warn
			if(force && this.connection != null)
				try {
					this.connection.close();
				} catch (Exception e) {
					logger.warn("Unable to close connection {} due to {}",	new Object[] { getConnectionString(), e });
					//throw new ProcessException("Connection closed problems");
				}
		}
	}

}
