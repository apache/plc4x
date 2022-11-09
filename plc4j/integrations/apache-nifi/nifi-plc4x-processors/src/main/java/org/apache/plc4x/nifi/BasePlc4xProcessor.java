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
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Set;

import org.apache.nifi.annotation.lifecycle.OnScheduled;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.components.ValidationContext;
import org.apache.nifi.components.ValidationResult;
import org.apache.nifi.components.Validator;
import org.apache.nifi.expression.ExpressionLanguageScope;
import org.apache.nifi.processor.AbstractProcessor;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.nifi.processor.util.StandardValidators;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;

public abstract class BasePlc4xProcessor extends AbstractProcessor {

	protected static final PropertyDescriptor PLC_CONNECTION_STRING = new PropertyDescriptor
        .Builder().name("PLC_CONNECTION_STRING")
        .displayName("PLC connection String")
        .description("PLC4X connection string used to connect to a given PLC device.")
        .required(true)
        .addValidator(new Plc4xConnectionStringValidator())
        .build();
	
    protected static final Relationship REL_SUCCESS = new Relationship.Builder()
	    .name("success")
	    .description("Successfully processed")
	    .build();
    
    protected static final Relationship REL_FAILURE = new Relationship.Builder()
        .name("failure")
        .description("An error occurred processing")
        .build();


    protected List<PropertyDescriptor> properties;
    protected Set<Relationship> relationships;
  
    protected String connectionString;
    protected Map<String, String> addressMap;


    private final PooledPlcDriverManager driverManager = new PooledPlcDriverManager();

    @Override
    protected void init(final ProcessorInitializationContext context) {
    	final List<PropertyDescriptor> properties = new ArrayList<>();
    	properties.add(PLC_CONNECTION_STRING);
        this.properties = Collections.unmodifiableList(properties);

    	
    	final Set<Relationship> relationships = new HashSet<>();
        relationships.add(REL_SUCCESS);
        relationships.add(REL_FAILURE);
        this.relationships = Collections.unmodifiableSet(relationships);
    }

    public Map<String, String> getPlcAddress() {
        return addressMap;
    }
    
    public String getConnectionString() {
        return connectionString;
    }

    Collection<String> getTags() {
        return addressMap.keySet();
    }
    String getAddress(String tagName) {
        return addressMap.get(tagName);
    }
    
	@Override
    public Set<Relationship> getRelationships() {
        return this.relationships;
    }

    @Override
    public final List<PropertyDescriptor> getSupportedPropertyDescriptors() {
        return properties;
    }
    
    //dynamic prop
    @Override
    protected PropertyDescriptor getSupportedDynamicPropertyDescriptor(final String propertyDescriptorName) {
        return new PropertyDescriptor.Builder()
                .name(propertyDescriptorName)
                .expressionLanguageSupported(ExpressionLanguageScope.NONE)
                .addValidator(StandardValidators.ATTRIBUTE_KEY_PROPERTY_NAME_VALIDATOR)
                .required(false)
                .dynamic(true)
                .build();
    }


    @OnScheduled
    public void onScheduled(final ProcessContext context) {
		connectionString = context.getProperty(PLC_CONNECTION_STRING.getName()).getValue();
		addressMap = new HashMap<>();
		//variables are passed as dynamic properties
		context.getProperties().keySet().stream().filter(PropertyDescriptor::isDynamic).forEach(
				t -> addressMap.put(t.getName(), context.getProperty(t.getName()).getValue()));
		if (addressMap.isEmpty()) {
			throw new PlcRuntimeException("No address specified");
		}	
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
            Objects.equals(getRelationships(), that.getRelationships()) &&
            Objects.equals(getConnectionString(), that.getConnectionString()) &&
            Objects.equals(addressMap, that.addressMap);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), properties, getRelationships(), getConnectionString(), addressMap);
    }

    public static class Plc4xConnectionStringValidator implements Validator {
        @Override
        public ValidationResult validate(String subject, String input, ValidationContext context) {
            // TODO: Add validation here ...
            return new ValidationResult.Builder().subject(subject).explanation("").valid(true).build();
        }
    }

    public static class Plc4xAddressStringValidator implements Validator {
        @Override
        public ValidationResult validate(String subject, String input, ValidationContext context) {
            // TODO: Add validation here ...
            return new ValidationResult.Builder().subject(subject).explanation("").valid(true).build();
        }
    }

    protected PooledPlcDriverManager getDriverManager() {
        return driverManager;
    }

}
