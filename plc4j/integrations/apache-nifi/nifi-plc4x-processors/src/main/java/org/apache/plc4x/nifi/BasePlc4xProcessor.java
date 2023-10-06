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

import java.io.OutputStream;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.concurrent.atomic.AtomicLong;

import org.apache.nifi.annotation.lifecycle.OnScheduled;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.components.ValidationContext;
import org.apache.nifi.components.ValidationResult;
import org.apache.nifi.components.Validator;
import org.apache.nifi.expression.ExpressionLanguageScope;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.AbstractProcessor;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.processor.util.StandardValidators;
import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.utils.cache.CachedPlcConnectionManager;
import org.apache.plc4x.nifi.address.AddressesAccessStrategy;
import org.apache.plc4x.nifi.address.AddressesAccessUtils;
import org.apache.plc4x.nifi.address.DynamicPropertyAccessStrategy;
import org.apache.plc4x.nifi.record.Plc4xWriter;
import org.apache.plc4x.nifi.record.SchemaCache;

public abstract class BasePlc4xProcessor extends AbstractProcessor {

    protected List<PropertyDescriptor> properties;
    protected Set<Relationship> relationships;
    protected volatile boolean debugEnabled;
    protected Integer cacheSize = 0;

    protected final SchemaCache schemaCache = new SchemaCache(0);
    protected AddressesAccessStrategy addressAccessStrategy;

    private CachedPlcConnectionManager connectionManager;

    protected CachedPlcConnectionManager getConnectionManager() {
        return connectionManager;
    }

    protected void refreshConnectionManager() {
        connectionManager = CachedPlcConnectionManager.getBuilder()
            .withMaxLeaseTime(Duration.ofSeconds(1000L))
            .withMaxWaitTime(Duration.ofSeconds(500L))
            .build();
    }


	public static final PropertyDescriptor PLC_CONNECTION_STRING = new PropertyDescriptor.Builder()
        .name("plc4x-connection-string")
        .displayName("PLC connection String")
        .description("PLC4X connection string used to connect to a given PLC device.")
        .required(true)
        .expressionLanguageSupported(ExpressionLanguageScope.FLOWFILE_ATTRIBUTES)
        .addValidator(new Plc4xConnectionStringValidator())
        .build();
	
    public static final PropertyDescriptor PLC_SCHEMA_CACHE_SIZE = new PropertyDescriptor.Builder()
        .name("plc4x-record-schema-cache-size")
        .displayName("Schema Cache Size")
		.description("Maximum number of entries in the cache. Can improve performance when addresses change dynamically.")
		.defaultValue("1")
		.required(true)
        .expressionLanguageSupported(ExpressionLanguageScope.VARIABLE_REGISTRY)
		.addValidator(StandardValidators.POSITIVE_INTEGER_VALIDATOR)
		.build();

    public static final PropertyDescriptor PLC_FUTURE_TIMEOUT_MILISECONDS = new PropertyDescriptor.Builder()
		.name("plc4x-request-timeout")
		.displayName("Timeout (miliseconds)")
		.description( "Request timeout in miliseconds")
		.defaultValue("10000")
		.required(true)
        .expressionLanguageSupported(ExpressionLanguageScope.FLOWFILE_ATTRIBUTES)
		.addValidator(StandardValidators.POSITIVE_INTEGER_VALIDATOR)
		.build();

    public static final PropertyDescriptor PLC_TIMESTAMP_FIELD_NAME = new PropertyDescriptor.Builder()
        .name("plc4x-timestamp-field-name")
        .displayName("Timestamp Field Name")
        .description("Name of the field that will display the timestamp of the operation.")
        .required(true)
        .expressionLanguageSupported(ExpressionLanguageScope.VARIABLE_REGISTRY)
        .addValidator(new Plc4xTimestampFieldValidator())
        .defaultValue("ts")
        .build();


    protected static final Relationship REL_SUCCESS = new Relationship.Builder()
	    .name("success")
	    .description("Successfully processed")
	    .build();
    
    protected static final Relationship REL_FAILURE = new Relationship.Builder()
        .name("failure")
        .description("An error occurred processing")
        .build();


    @Override
    protected void init(final ProcessorInitializationContext context) {
    	final List<PropertyDescriptor> properties = new ArrayList<>();

    	properties.add(PLC_CONNECTION_STRING);
        properties.add(AddressesAccessUtils.PLC_ADDRESS_ACCESS_STRATEGY);
        properties.add(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY);
        properties.add(AddressesAccessUtils.ADDRESS_FILE_PROPERTY);
        properties.add(PLC_SCHEMA_CACHE_SIZE);
        properties.add(PLC_FUTURE_TIMEOUT_MILISECONDS);
        properties.add(PLC_TIMESTAMP_FIELD_NAME);
        this.properties = Collections.unmodifiableList(properties);

    	
    	final Set<Relationship> relationships = new HashSet<>();
        relationships.add(REL_SUCCESS);
        relationships.add(REL_FAILURE);
        this.relationships = Collections.unmodifiableSet(relationships);
    }

    public Map<String, String> getPlcAddressMap(ProcessContext context, FlowFile flowFile) {
        return addressAccessStrategy.extractAddresses(context, flowFile);
    }
    
    public String getConnectionString(ProcessContext context, FlowFile flowFile) {
        return context.getProperty(PLC_CONNECTION_STRING).evaluateAttributeExpressions(flowFile).getValue();
    }

    public Long getTimeout(ProcessContext context, FlowFile flowFile) {
        return context.getProperty(PLC_FUTURE_TIMEOUT_MILISECONDS).evaluateAttributeExpressions(flowFile).asLong();
    }

    public String getTimestampField(ProcessContext context) {
        return context.getProperty(PLC_TIMESTAMP_FIELD_NAME).evaluateAttributeExpressions().getValue();
    }

    public SchemaCache getSchemaCache() {
        return schemaCache;
    }
    
	@Override
    public Set<Relationship> getRelationships() {
        return this.relationships;
    }

    @Override
    public final List<PropertyDescriptor> getSupportedPropertyDescriptors() {
        return properties;
    }
    
    @Override
    protected PropertyDescriptor getSupportedDynamicPropertyDescriptor(final String propertyDescriptorName) {
        return new PropertyDescriptor.Builder()
                .name(propertyDescriptorName)
                .expressionLanguageSupported(ExpressionLanguageScope.FLOWFILE_ATTRIBUTES)
                .addValidator(StandardValidators.ATTRIBUTE_KEY_PROPERTY_NAME_VALIDATOR)
                .dependsOn(AddressesAccessUtils.PLC_ADDRESS_ACCESS_STRATEGY, AddressesAccessUtils.ADDRESS_PROPERTY)
                .addValidator(new DynamicPropertyAccessStrategy.TagValidator(AddressesAccessUtils.getManager()))
                .required(false)
                .dynamic(true)
                .build();
    }


    @OnScheduled
    public void onScheduled(final ProcessContext context) {
        Integer newCacheSize = context.getProperty(PLC_SCHEMA_CACHE_SIZE).evaluateAttributeExpressions().asInteger();
        if (!newCacheSize.equals(cacheSize)){
            schemaCache.restartCache(newCacheSize);
            cacheSize = newCacheSize;
        }
        refreshConnectionManager();
        debugEnabled = getLogger().isDebugEnabled();
        addressAccessStrategy = AddressesAccessUtils.getAccessStrategy(context);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof BasePlc4xProcessor)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        BasePlc4xProcessor that = (BasePlc4xProcessor) o;
        return Objects.equals(properties, that.properties) &&
            Objects.equals(getRelationships(), that.getRelationships());
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), properties, getRelationships());
    }

    protected PlcWriteRequest getWriteRequest(final ComponentLog logger,
            final Map<String, String> addressMap, final Map<String, PlcTag> tags, final Map<String, ? extends Object> presentTags,
            final PlcConnection connection, final AtomicLong nrOfRowsHere) {

        PlcWriteRequest.Builder builder = connection.writeRequestBuilder();

        if (tags != null){
            for (Map.Entry<String,PlcTag> tag : tags.entrySet()){
                if (presentTags.containsKey(tag.getKey())) {
                    builder.addTag(tag.getKey(), tag.getValue(), presentTags.get(tag.getKey()));
                    if (nrOfRowsHere != null) {
                        nrOfRowsHere.incrementAndGet();
                    }
                } else {
                    if (debugEnabled)
                        logger.debug("PlcTag " + tag + " is declared as address but was not found on input record.");
                }
            }
        } else {
            if (debugEnabled)
                logger.debug("PlcTypes resolution not found in cache and will be added with key: " + addressMap);
            for (Map.Entry<String,String> entry: addressMap.entrySet()){
                if (presentTags.containsKey(entry.getKey())) {
                    builder.addTagAddress(entry.getKey(), entry.getValue(), presentTags.get(entry.getKey()));
                    if (nrOfRowsHere != null) {
                        nrOfRowsHere.incrementAndGet();
                    }
                }
            }
        }
         
        return builder.build();
    }

    protected PlcReadRequest getReadRequest(final ComponentLog logger, 
            final Map<String, String> addressMap, final Map<String, PlcTag> tags,
            final PlcConnection connection) {

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

    protected void evaluateWriteResponse(final ComponentLog logger, Map<String, ? extends Object> values, PlcWriteResponse plcWriteResponse) {

		boolean codeErrorPresent = false;
		List<String> tagsAtError = null;

		PlcResponseCode code = null;

		for (String tag : plcWriteResponse.getTagNames()) {
			code = plcWriteResponse.getResponseCode(tag);
			if (!code.equals(PlcResponseCode.OK)) {
				if (tagsAtError == null) {
					tagsAtError = new ArrayList<>();
				}
				logger.error("Not OK code when writing the data to PLC for tag " + tag 
					+ " with value  " + values.get(tag).toString() 
					+ " in addresss " + plcWriteResponse.getTag(tag).getAddressString());
				
			codeErrorPresent = true;
			tagsAtError.add(tag);
						
			}
		}
		if (codeErrorPresent) {
			throw new ProcessException("At least one error was found when while writting tags: " + tagsAtError.toString());
		}
	}

   protected void evaluateReadResponse(final ProcessSession session, final FlowFile flowFile, final PlcReadResponse response) {
        Map<String, String> attributes = new HashMap<>();
        for (String tagName : response.getTagNames()) {
            for (int i = 0; i < response.getNumberOfValues(tagName); i++) {
                Object value = response.getObject(tagName, i);
                attributes.put(tagName, String.valueOf(value));
            }
        }
        session.putAllAttributes(flowFile, attributes);
    }

    protected long evaluateReadResponse(final ProcessContext context, final ComponentLog logger, final FlowFile originalFlowFile,
			Plc4xWriter plc4xWriter, OutputStream out, final RecordSchema recordSchema, PlcReadResponse readResponse)
			throws Exception {

		if(originalFlowFile == null) //there is no inherit attributes to use in writer service 
			return plc4xWriter.writePlcReadResponse(readResponse, out, logger, null, recordSchema, getTimestampField(context));
		else 
			return plc4xWriter.writePlcReadResponse(readResponse, out, logger, null, recordSchema, originalFlowFile, getTimestampField(context));
	}

    protected static class Plc4xConnectionStringValidator implements Validator {
        @Override
        public ValidationResult validate(String subject, String input, ValidationContext context) {
            DefaultPlcDriverManager manager = new DefaultPlcDriverManager();
            
            if (context.isExpressionLanguageSupported(subject) && context.isExpressionLanguagePresent(input)) {
                return new ValidationResult.Builder().subject(subject).input(input).explanation("Expression Language Present").valid(true).build();
            }
            try {
                PlcDriver driver =  manager.getDriverForUrl(input);
                driver.getConnection(input);
            } catch (PlcConnectionException e) {
                return new ValidationResult.Builder().subject(subject)
                    .explanation(e.getMessage())
                    .valid(false)
                    .build();
            }
            return new ValidationResult.Builder().subject(subject)
                .explanation("")
                .valid(true)
                .build();
        }
    }

    protected static class Plc4xTimestampFieldValidator implements Validator {
        @Override
        public ValidationResult validate(String subject, String input, ValidationContext context) {

            if (context.isExpressionLanguageSupported(subject) && context.isExpressionLanguagePresent(input)) {
                return new ValidationResult.Builder().subject(subject).input(input).explanation("Expression Language Present").valid(true).build();
            }
            
            Map<String, String> allProperties = context.getAllProperties();
            allProperties.remove(subject);

            if (allProperties.containsValue(input)) {
                return new ValidationResult.Builder().subject(subject)
                    .explanation("Timestamp field must be unique")
                    .valid(false)
                    .build(); 
            }
            return new ValidationResult.Builder().subject(subject)
                .explanation("")
                .valid(true)
                .build();

        }
    }
}
