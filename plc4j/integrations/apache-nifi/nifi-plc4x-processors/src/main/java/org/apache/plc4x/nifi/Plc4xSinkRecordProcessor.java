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

import java.io.InputStream;
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
import org.apache.nifi.annotation.behavior.TriggerSerially;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.processor.util.StandardValidators;
import org.apache.nifi.serialization.RecordReader;
import org.apache.nifi.serialization.RecordReaderFactory;
import org.apache.nifi.serialization.record.Record;
import org.apache.nifi.util.StopWatch;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;

@TriggerSerially
@Tags({ "plc4x-sink" })
@InputRequirement(InputRequirement.Requirement.INPUT_REQUIRED)
@CapabilityDescription("Processor able to write data to industrial PLCs using Apache PLC4X")
@WritesAttributes({ @WritesAttribute(attribute = "value", description = "some value") })
public class Plc4xSinkRecordProcessor extends BasePlc4xProcessor {

	public static final String RESULT_ROW_COUNT = "plc4x.write.row.count";
	public static final String RESULT_QUERY_DURATION = "plc4x.write.query.duration";
	public static final String RESULT_QUERY_EXECUTION_TIME = "plc4x.write.query.executiontime";
	public static final String RESULT_QUERY_FETCH_TIME = "plc4x.write.query.fetchtime";
	public static final String INPUT_FLOWFILE_UUID = "input.flowfile.uuid";
	public static final String RESULT_ERROR_MESSAGE = "plc4x.write.error.message";
	
	public static final PropertyDescriptor PLC_RECORD_READER_FACTORY = new PropertyDescriptor.Builder()
			.name("record-reader").displayName("Record Reader")
			.description(
					"Specifies the Controller Service to use for reading record from a FlowFile. The Record Reader may use Inherit Schema to emulate the inferred schema behavior, i.e. "
							+ "an explicit schema need not be defined in the reader, and will be supplied by the same logic used to infer the schema from the column types.")
			.identifiesControllerService(RecordReaderFactory.class)
			.required(true)
			.build();

	public static final PropertyDescriptor PLC_WRITE_FUTURE_TIMEOUT_MILISECONDS = new PropertyDescriptor.Builder().name("plc4x-record-write-timeout").displayName("Write timeout (miliseconds)")
			.description("Write timeout in miliseconds")
			.defaultValue("10000")
			.required(true)
			.addValidator(StandardValidators.POSITIVE_INTEGER_VALIDATOR)
			.build();
	

	public Plc4xSinkRecordProcessor() {
	}

	@Override
	protected void init(final ProcessorInitializationContext context) {
		super.init(context);
		final Set<Relationship> r = new HashSet<>();
		r.addAll(super.getRelationships());
		this.relationships = Collections.unmodifiableSet(r);

		final List<PropertyDescriptor> pds = new ArrayList<>();
		pds.addAll(super.getSupportedPropertyDescriptors());
		pds.add(PLC_RECORD_READER_FACTORY);
		pds.add(PLC_WRITE_FUTURE_TIMEOUT_MILISECONDS);
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

		final FlowFile originalFlowFile = fileToProcess;

		InputStream in = session.read(originalFlowFile);

		RecordReader recordReader;
		try {
			recordReader = context.getProperty(PLC_RECORD_READER_FACTORY)
					.asControllerService(RecordReaderFactory.class)
					.createRecordReader(originalFlowFile, in, logger);
		} catch (Exception e) {
			throw new ProcessException(e);
		} 
		Record record = null;
		
		
		try {
			while ((record = recordReader.nextRecord()) != null) {
				long nrOfRowsHere = 0L;
				PlcWriteResponse plcWriteResponse = null;
				PlcWriteRequest writeRequest = null;

				Map<String,String> addressMap = getPlcAddressMap(context, fileToProcess);
				final Map<String, PlcTag> tags = getSchemaCache().retrieveTags(addressMap);

				try (PlcConnection connection = getConnectionManager().getConnection(getConnectionString())) {
					PlcWriteRequest.Builder builder = connection.writeRequestBuilder();
					
					
					if (tags != null){
						for (Map.Entry<String,PlcTag> tag : tags.entrySet()){
							if (record.toMap().containsKey(tag.getKey())) {
								builder.addTag(tag.getKey(), tag.getValue(), record.getValue(tag.getKey()));
								nrOfRowsHere++;
							} else {
								logger.debug("PlcTag " + tag + " is declared as address but was not found on input record.");
							}
						}
					} else {
						logger.debug("Plc-Avro schema and PlcTypes resolution not found in cache and will be added with key: " + addressMap.toString());
						for (Map.Entry<String,String> entry: addressMap.entrySet()){
							if (record.toMap().containsKey(entry.getKey())) {
								builder.addTagAddress(entry.getKey(), entry.getValue(), record.getValue(entry.getKey()));
								nrOfRowsHere++;
							} else {
								logger.debug("PlcTag " + entry.getKey() + " with address " + entry.getValue() + " was not found on input record.");
							}
						}
					}
					writeRequest = builder.build();

					plcWriteResponse = writeRequest.execute().get(
						context.getProperty(PLC_WRITE_FUTURE_TIMEOUT_MILISECONDS.getName()).asInteger(), TimeUnit.MILLISECONDS
						);
				} catch (Exception e) {
					System.out.println(e.getMessage());
					in.close();
					logger.error("Exception writing the data to PLC", e);
					session.transfer(originalFlowFile, REL_FAILURE);
					session.commitAsync();
					throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
				}

				if (tags == null){
					getLogger().debug("Adding PlcTypes resolution into cache with key: " + addressMap.toString());
					getSchemaCache().addSchema(
						addressMap, 
						writeRequest.getTagNames(),
						writeRequest.getTags(),
						null
					);
				}

				// Response check if values were written
				PlcResponseCode code = null;

				for (String tag : plcWriteResponse.getTagNames()) {
					code = plcWriteResponse.getResponseCode(tag);
					if (!code.equals(PlcResponseCode.OK)) {
						logger.error("Not OK code when writing the data to PLC for tag " + tag 
							+ " with value  " + record.getValue(tag).toString() 
							+ " in addresss " + plcWriteResponse.getTag(tag).getAddressString());
						throw new ProcessException("Writing response code for " + plcWriteResponse.getTag(tag).getAddressString() + "was " + code.name() + ", expected OK");
					}
				}
				
				nrOfRows.getAndAdd(nrOfRowsHere);
			}
			in.close();
		} catch (Exception e) {
			throw new ProcessException(e);
		} 
		
		long executionTimeElapsed = executeTime.getElapsed(TimeUnit.MILLISECONDS);
		final Map<String, String> attributesToAdd = new HashMap<>();
		attributesToAdd.put(RESULT_ROW_COUNT, String.valueOf(nrOfRows.get()));
		attributesToAdd.put(RESULT_QUERY_EXECUTION_TIME, String.valueOf(executionTimeElapsed));

		FlowFile resultSetFF = session.putAllAttributes(originalFlowFile, attributesToAdd);

		logger.info("Writing {} fields from {} records; transferring to 'success'", new Object[] { nrOfRows.get(), resultSetFF });
		// Report a FETCH event if there was an incoming flow file, or a RECEIVE event
		// otherwise
		if (context.hasIncomingConnection()) {
			session.getProvenanceReporter().fetch(resultSetFF, "Writted " + nrOfRows.get() + " rows", executionTimeElapsed);
		} else {
			session.getProvenanceReporter().receive(resultSetFF, "Writted " + nrOfRows.get() + " rows", executionTimeElapsed);
		}

		session.transfer(resultSetFF, BasePlc4xProcessor.REL_SUCCESS);
		session.commitAsync();
	}
}
