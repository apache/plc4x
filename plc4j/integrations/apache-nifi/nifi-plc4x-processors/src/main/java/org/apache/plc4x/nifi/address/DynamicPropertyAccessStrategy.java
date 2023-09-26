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

import java.util.HashMap;
import java.util.Map;

import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.components.ValidationContext;
import org.apache.nifi.components.ValidationResult;
import org.apache.nifi.components.Validator;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.nifi.BasePlc4xProcessor;


public class DynamicPropertyAccessStrategy implements AddressesAccessStrategy{

    private Map<String,String> extractAddressesFromAttributes(final ProcessContext context, final FlowFile flowFile) {
        Map<String,String> addressMap = new HashMap<>();

        context.getProperties().keySet().stream().filter(PropertyDescriptor::isDynamic).forEach(
            t -> addressMap.put(t.getName(), context.getProperty(t.getName()).evaluateAttributeExpressions(flowFile).getValue()));
        
        return addressMap; 
    }

    @Override
    public Map<String, String> extractAddresses(final ProcessContext context, final FlowFile flowFile) {
        return extractAddressesFromAttributes(context, flowFile);
    }

    public static class TagValidator implements Validator {
        @Override
        public ValidationResult validate(String subject, String input, ValidationContext context) {
            String connectionString = context.getProperty(BasePlc4xProcessor.PLC_CONNECTION_STRING).getValue();

            if (context.isExpressionLanguageSupported(subject) && context.isExpressionLanguagePresent(input) || 
                context.isExpressionLanguagePresent(connectionString)) {
                return new ValidationResult.Builder().subject(subject).input(input).explanation("Expression Language Present").valid(true).build();
            }

            try {
                DefaultPlcDriverManager manager = new DefaultPlcDriverManager();
                PlcDriver driver =  manager.getDriverForUrl(connectionString);

                if (!context.isExpressionLanguagePresent(input)) {
                    driver.prepareTag(input);
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
