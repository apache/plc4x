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
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;

import org.apache.nifi.annotation.behavior.InputRequirement;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.configuration.DefaultSchedule;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.documentation.SeeAlso;
import org.apache.nifi.annotation.documentation.Tags;
import org.apache.nifi.annotation.lifecycle.OnScheduled;
import org.apache.nifi.annotation.lifecycle.OnStopped;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.components.ValidationContext;
import org.apache.nifi.components.ValidationResult;
import org.apache.nifi.components.Validator;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.processor.util.StandardValidators;
import org.apache.nifi.serialization.RecordSetWriterFactory;
import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.nifi.util.StopWatch;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.nifi.subscription.Plc4xListenerDispatcher;
import org.apache.plc4x.nifi.subscription.Plc4xSubscriptionType;
import org.apache.plc4x.nifi.record.Plc4xWriter;
import org.apache.plc4x.nifi.record.RecordPlc4xWriter;

@DefaultSchedule(period="0.1 sec")
@Tags({"plc4x", "get", "input", "source", "listen", "record"})
@SeeAlso({Plc4xSourceRecordProcessor.class, Plc4xSinkRecordProcessor.class})
@InputRequirement(InputRequirement.Requirement.INPUT_FORBIDDEN)
@CapabilityDescription("Processor able to read data from industrial PLCs using Apache PLC4X subscriptions")
@WritesAttributes({ 
	@WritesAttribute(attribute = Plc4xListenRecordProcessor.RESULT_ROW_COUNT, description = "Number of rows written into the output FlowFile"),
	@WritesAttribute(attribute = Plc4xListenRecordProcessor.RESULT_LAST_EVENT, description = "Time elapsed from last subscription event")
 })
public class Plc4xListenRecordProcessor extends BasePlc4xProcessor {

	public static final String RESULT_ROW_COUNT = "plc4x.listen.row.count";
	public static final String RESULT_LAST_EVENT = "plc4x.listen.lastEvent";

    protected Plc4xSubscriptionType subscriptionType = null;
    protected Long cyclingPollingInterval = null;
	protected final BlockingQueue<PlcSubscriptionEvent> events = new LinkedBlockingQueue<>();
	protected Plc4xListenerDispatcher dispatcher;
	protected RecordSchema recordSchema;
	protected Thread readerThread;
	protected Map<String, String> addressMap;
	final StopWatch executeTime = new StopWatch(false);

	public static final PropertyDescriptor PLC_RECORD_WRITER_FACTORY = new PropertyDescriptor.Builder()
        .name("plc4x-record-writer")
        .displayName("Record Writer")
		.description("Specifies the Controller Service to use for writing results to a FlowFile. The Record Writer may use Inherit Schema to emulate the inferred schema behavior, i.e. "
				+ "an explicit schema need not be defined in the writer, and will be supplied by the same logic used to infer the schema from the column types.")
		.identifiesControllerService(RecordSetWriterFactory.class)
		.required(true)
		.build();

    public static final PropertyDescriptor PLC_SUBSCRIPTION_TYPE = new PropertyDescriptor.Builder()
        .name("plc4x-subscription-type")
        .displayName("Subscription Type")
		.description("Sets the subscription type. The subscritpion types available for each driver are stated in the documentation.")
		.allowableValues(Plc4xSubscriptionType.values())
		.required(true)
        .defaultValue(Plc4xSubscriptionType.CHANGE.name())
		.build();

    public static final PropertyDescriptor PLC_SUBSCRIPTION_CYCLIC_POLLING_INTERVAL = new PropertyDescriptor.Builder()
        .name("plc4x-subscription-cyclic-polling-interval")
        .displayName("Cyclic polling interval")
		.description("In case of Cyclic subscription type a time interval must be provided.")
		.dependsOn(PLC_SUBSCRIPTION_TYPE, Plc4xSubscriptionType.CYCLIC.name())
		.required(true)
        .addValidator(StandardValidators.POSITIVE_INTEGER_VALIDATOR)
		.addValidator(new CyclycPollingIntervalValidator())
        .defaultValue("10000")
		.build();

	@Override
	protected void init(final ProcessorInitializationContext context) {
		super.init(context);
		final Set<Relationship> relationships = new HashSet<>();
        relationships.add(REL_SUCCESS);
		this.relationships = Collections.unmodifiableSet(relationships);

		final List<PropertyDescriptor> pds = new ArrayList<>();
		pds.addAll(super.getSupportedPropertyDescriptors());
		pds.add(PLC_RECORD_WRITER_FACTORY);
		pds.add(PLC_SUBSCRIPTION_TYPE);
		pds.add(PLC_SUBSCRIPTION_CYCLIC_POLLING_INTERVAL);
		this.properties = Collections.unmodifiableList(pds);
	}

    @Override
    @OnScheduled
    public void onScheduled(final ProcessContext context) {
		super.onScheduled(context);
		subscriptionType = Plc4xSubscriptionType.valueOf(context.getProperty(PLC_SUBSCRIPTION_TYPE).getValue());
        cyclingPollingInterval = context.getProperty(PLC_SUBSCRIPTION_CYCLIC_POLLING_INTERVAL).asLong();
		createDispatcher(context, events);
	}

    protected void createDispatcher(final ProcessContext context, final BlockingQueue<PlcSubscriptionEvent> events) {
		if (readerThread != null) {
			return;
		}

		// create the dispatcher and calls open() to start listening to the plc subscription
        dispatcher =  new Plc4xListenerDispatcher(getTimeout(context, null), subscriptionType, cyclingPollingInterval, getLogger(), events);
		try {
			addressMap = getPlcAddressMap(context, null);
			dispatcher.open(getConnectionString(context, null), addressMap);
		} catch (Exception e) {
			if (debugEnabled) {
				getLogger().debug("Error creating a the subscription event dispatcher");
				e.printStackTrace();
			}
			throw new ProcessException(e);
		}

		if (dispatcher.isRunning()) {
			readerThread = new Thread(dispatcher);
			readerThread.setName(getClass().getName() + " [" + getIdentifier() + "]");
			readerThread.setDaemon(true);
			readerThread.start();
		}
		executeTime.start();
    }

    @OnStopped
    public void closeDispatcher() throws ProcessException {
		executeTime.stop();
		if (readerThread != null) {
			readerThread.interrupt();
			try {
				readerThread.join();
			} catch (InterruptedException e) {
				throw new ProcessException(e);
			}
			if (!readerThread.isAlive()){
				readerThread = null;
			}
		}
    }

	protected PlcSubscriptionEvent getMessage(final ProcessContext context) {
		if (readerThread != null && readerThread.isAlive()) {
			return events.poll();
			
		}
		// If dispatcher is not running the connection broke or gave a time out.
		if (debugEnabled) {
			getLogger().debug("Connection to Plc broke. Trying to restart connection");
		}
		closeDispatcher();
		createDispatcher(context, events);
		throw new ProcessException("Connection to Plc broke. Trying to restart connection");
	}
	
	@Override
	public void onTrigger(final ProcessContext context, final ProcessSession session) throws ProcessException {

		DefaultPlcSubscriptionEvent event = (DefaultPlcSubscriptionEvent) getMessage(context);

		if (event == null) {
			return;
		} else {
			session.adjustCounter("Messages Received", 1L, false);
		}


		final AtomicLong nrOfRows = new AtomicLong(0L);

		FlowFile resultSetFF = session.create();

		Plc4xWriter plc4xWriter = new RecordPlc4xWriter(context.getProperty(PLC_RECORD_WRITER_FACTORY).asControllerService(RecordSetWriterFactory.class), Collections.emptyMap());

		try {
			session.write(resultSetFF, out -> {
				try {
					nrOfRows.set(plc4xWriter.writePlcReadResponse(event, out, getLogger(), null, recordSchema, getTimestampField(context)));
				}  catch (Exception e) {
					getLogger().error("Exception reading the data from PLC", e);
					throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
				}

				if (recordSchema == null){
					addTagsToCache(event, plc4xWriter);
				}
			});
			resultSetFF = completeResultFlowFile(session, nrOfRows, resultSetFF, plc4xWriter);
			session.transfer(resultSetFF, REL_SUCCESS);

			executeTime.start();

		} catch (Exception e) {
			getLogger().error("Got an error while trying to get a subscription event", e);
			throw new ProcessException("Got an error while trying to get a subscription event", e);
		}
	}

	private void addTagsToCache(DefaultPlcSubscriptionEvent event, Plc4xWriter plc4xWriter) {
		if (debugEnabled)
			getLogger().debug("Adding Plc-Avro schema and PlcTypes resolution into cache with key: " + addressMap.toString());
		
		// Add schema to the cache
		LinkedHashSet<String> addressNames = new LinkedHashSet<>();
		addressNames.addAll(event.getTagNames());
		
		List<PlcTag> addressTags = addressNames.stream().map(addr -> 
				new PlcTag() {
					@Override
					public String getAddressString() {
						return addr;
					}

					@Override
					public PlcValueType getPlcValueType() {
						return event.getPlcValue(addr).getPlcValueType();
					}
				}
			).collect(Collectors.toList()); 

		getSchemaCache().addSchema(
			addressMap, 
			addressNames,
			addressTags,
			plc4xWriter.getRecordSchema()
		);
		recordSchema = getSchemaCache().retrieveSchema(addressMap);
	}

	private FlowFile completeResultFlowFile(final ProcessSession session, final AtomicLong nrOfRows, FlowFile resultSetFF,
			Plc4xWriter plc4xWriter) {
				
		long executionTimeElapsed = executeTime.getElapsed(TimeUnit.MILLISECONDS);
		executeTime.stop();
		
		final Map<String, String> attributesToAdd = new HashMap<>();
		attributesToAdd.put(RESULT_ROW_COUNT, String.valueOf(nrOfRows.get()));
		attributesToAdd.put(RESULT_LAST_EVENT, String.valueOf(executionTimeElapsed));

		attributesToAdd.putAll(plc4xWriter.getAttributesToAdd());
		resultSetFF = session.putAllAttributes(resultSetFF, attributesToAdd);
		plc4xWriter.updateCounters(session);
		getLogger().info("{} contains {} records; transferring to 'success'", resultSetFF, nrOfRows.get());
		
		session.getProvenanceReporter().receive(resultSetFF, "Retrieved " + nrOfRows.get() + " rows from subscription", executionTimeElapsed);
		return resultSetFF;
	}


	protected static class CyclycPollingIntervalValidator implements Validator {
		@Override
		public ValidationResult validate(String subject, String input, ValidationContext context) {
			if (context.getProperty(PLC_FUTURE_TIMEOUT_MILISECONDS).asLong() > Long.valueOf(input)) {
				return new ValidationResult.Builder().valid(true).build();
			} else {
				return new ValidationResult.Builder()
				.valid(false)
				.input(input)
				.subject(PLC_SUBSCRIPTION_CYCLIC_POLLING_INTERVAL.getDisplayName())
				.explanation(String.format("it must me smaller than the value of %s", PLC_FUTURE_TIMEOUT_MILISECONDS.getDisplayName()))
				.build();
			}
		}	
	}
}
