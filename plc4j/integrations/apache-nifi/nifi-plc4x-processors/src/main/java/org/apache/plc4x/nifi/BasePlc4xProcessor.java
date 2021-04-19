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

import org.apache.nifi.annotation.lifecycle.OnScheduled;
import org.apache.nifi.components.*;
import org.apache.nifi.processor.AbstractProcessor;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.ProcessorInitializationContext;
import org.apache.nifi.processor.Relationship;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;

import java.util.*;

public abstract class BasePlc4xProcessor extends AbstractProcessor {

    private static final PropertyDescriptor PLC_CONNECTION_STRING = new PropertyDescriptor
        .Builder().name("PLC_CONNECTION_STRING")
        .displayName("PLC connection String")
        .description("PLC4X connection string used to connect to a given PLC device.")
        .required(true)
        .addValidator(new Plc4xConnectionStringValidator())
        .build();
    private static final PropertyDescriptor PLC_ADDRESS_STRING = new PropertyDescriptor
        .Builder().name("PLC_ADDRESS_STRING")
        .displayName("PLC resource address String")
        .description("PLC4X address string used identify the resource to read/write on a given PLC device " +
            "(Multiple values supported). The expected format is: {name}={address}(;{name}={address})*")
        .required(true)
        .addValidator(new Plc4xAddressStringValidator())
        .build();

    static final Relationship SUCCESS = new Relationship.Builder()
        .name("SUCCESS")
        .description("Successfully processed")
        .build();
    static final Relationship FAILURE = new Relationship.Builder()
        .name("FAILURE")
        .description("An error occurred processing")
        .build();

    private List<PropertyDescriptor> descriptors;

    Set<Relationship> relationships;

    private String connectionString;
    private Map<String, String> addressMap;

    private final PooledPlcDriverManager driverManager = new PooledPlcDriverManager();

    @Override
    protected void init(final ProcessorInitializationContext context) {
        this.descriptors = Arrays.asList(PLC_CONNECTION_STRING, PLC_ADDRESS_STRING);
        this.relationships = new HashSet<>(Arrays.asList(SUCCESS, FAILURE));
    }

    public String getConnectionString() {
        return connectionString;
    }

    Collection<String> getFields() {
        return addressMap.keySet();
    }
    String getAddress(String field) {
        return addressMap.get(field);
    }

    @Override
    public Set<Relationship> getRelationships() {
        return this.relationships;
    }

    @Override
    public final List<PropertyDescriptor> getSupportedPropertyDescriptors() {
        return descriptors;
    }

    @OnScheduled
    public void onScheduled(final ProcessContext context) {
        PropertyValue property = context.getProperty(PLC_CONNECTION_STRING.getName());
        connectionString = property.getValue();
        addressMap = new HashMap<>();
        PropertyValue addresses = context.getProperty(PLC_ADDRESS_STRING.getName());
        for (String segment : addresses.getValue().split(";")) {
            String[] parts = segment.split("=");
            if(parts.length != 2) {
                throw new PlcRuntimeException("Invalid address format");
            }
            addressMap.put(parts[0], parts[1]);
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
        return Objects.equals(descriptors, that.descriptors) &&
            Objects.equals(getRelationships(), that.getRelationships()) &&
            Objects.equals(getConnectionString(), that.getConnectionString()) &&
            Objects.equals(addressMap, that.addressMap);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), descriptors, getRelationships(), getConnectionString(), addressMap);
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
