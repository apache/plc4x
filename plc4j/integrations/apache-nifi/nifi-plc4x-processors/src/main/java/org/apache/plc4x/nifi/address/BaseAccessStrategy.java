/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.nifi.address;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.Map;

import org.apache.nifi.components.AllowableValue;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.components.ValidationContext;
import org.apache.nifi.components.ValidationResult;
import org.apache.nifi.components.Validator;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.nifi.BasePlc4xProcessor;


public abstract class BaseAccessStrategy implements AddressesAccessStrategy{
    private boolean isInitializated = false;
    private boolean isDynamic;
    protected Map<String,String> cachedAddresses = null;

    protected AllowableValue allowableValue;
    protected List<PropertyDescriptor> propertyDescriptors = new ArrayList<>();

    protected Map<String, String> getCachedAddresses() {
        return cachedAddresses;
    }

    public Map<String,String> extractAddressesFromResources(final ProcessContext context, final FlowFile flowFile) {
        throw new UnsupportedOperationException("Method 'extractAddressesFromResources' not implemented");
    }


    @Override
    public Map<String, String> extractAddresses(final ProcessContext context, final FlowFile flowFile) {
        if (!isInitializated) {
            getPropertyDescriptors().forEach(prop -> {
                if (context.isExpressionLanguagePresent(prop)){
                    isDynamic = true;
                }
            });
            isInitializated = true;
        }

        Map<String, String> result = getCachedAddresses();
        if (result == null) {
            result = extractAddressesFromResources(context, flowFile);
            if (!isDynamic) {
                cachedAddresses = result;
            }
        }
        return result;
    }

    public static class TagValidator implements Validator {
        
        private DefaultPlcDriverManager manager;

        public TagValidator(DefaultPlcDriverManager manager) {
            this.manager = manager;
        }

        protected void checkTags(PlcDriver driver, Collection<String> tags) {
            for (String tag : tags) {
                driver.prepareTag(tag);
            }
        }

        protected Collection<String> getTags(String input) throws Exception {
            throw new UnsupportedOperationException("Method 'getTags' not implemented");
        } 

        @Override
        public ValidationResult validate(String subject, String input, ValidationContext context) {
            String connectionString = context.getProperty(BasePlc4xProcessor.PLC_CONNECTION_STRING).getValue();

            if (context.isExpressionLanguageSupported(subject) && context.isExpressionLanguagePresent(input) || 
                context.isExpressionLanguagePresent(connectionString)) {
                return new ValidationResult.Builder().subject(subject).input(input)
                    .explanation("Expression Language Present").valid(true).build();
            }

            try {
                PlcDriver driver = manager.getDriverForUrl(connectionString);

                if (!context.isExpressionLanguagePresent(input)) {
                    checkTags(driver, getTags(input));
                } 
                
            }catch (Exception e) {
                    return new ValidationResult.Builder().subject(subject)
                        .explanation(e.getLocalizedMessage())
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
