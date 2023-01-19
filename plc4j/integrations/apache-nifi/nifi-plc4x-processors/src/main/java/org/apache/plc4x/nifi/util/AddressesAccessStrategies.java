/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.nifi.util;

import java.util.HashMap;
import java.util.Map;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.nifi.components.AllowableValue;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.expression.ExpressionLanguageScope;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.processor.util.JsonValidator;

public interface AddressesAccessStrategies {

    public static final PropertyDescriptor PLC_ADDRESS_ACCESS_STRATEGY = new PropertyDescriptor.Builder()
        .name("plc4x-address-access-strategy")
        .displayName("Address Access Strategy")
        .description("Strategy used to obtain the PLC addresses")
        .required(true)
        .build();

    public static final AllowableValue ADDRESS_PROPERTY = new AllowableValue(
        "property-address", 
        "Use Properties as Addresses",
        "Each property will be treated as tag-address pairs after Expression Language is evaluated.");

    public static final AllowableValue ADDRESS_TEXT = new AllowableValue(
        "text-address", 
        "Use 'Address Text' Property",
        "Addresses will be obtained from 'Address Text' Property. It's content must be a valid JSON " +
            "after Expression Language is evaluated. ");
        
    public static final PropertyDescriptor ADDRESS_TEXT_PROPERTY = new PropertyDescriptor.Builder()
        .name("text-address-property")
        .displayName("Address Text")
        .description("Must contain a valid JSON object after Expression Language is evaluated. "
            + "Each field-value is treated as tag-address.")
        .expressionLanguageSupported(ExpressionLanguageScope.FLOWFILE_ATTRIBUTES)
        .addValidator(new JsonValidator())
        .dependsOn(PLC_ADDRESS_ACCESS_STRATEGY, ADDRESS_TEXT)
        .required(true)
        .build();


    public static AddressAccessStrategy getAccessStrategy(final ProcessContext context) {
        String value = context.getProperty(PLC_ADDRESS_ACCESS_STRATEGY).getValue();
        if (ADDRESS_PROPERTY.getValue().equalsIgnoreCase(value)) 
            return new AddressPropertyAccess();
        else if (ADDRESS_TEXT.getValue().equalsIgnoreCase(value))
            return new AddressTextAccess();
        return null;
    }

    public static interface  AddressAccessStrategy {
        public Map<String,String> extractAddresses(final ProcessContext context, final FlowFile flowFile);
    }

    public static class AddressTextAccess implements AddressAccessStrategy{

        private Map<String,String> extractAddressesFromText(String input) throws JsonMappingException, JsonProcessingException {
            ObjectMapper mapper = new ObjectMapper();

            return mapper.readValue(input, Map.class);
        }

        @Override
        public Map<String, String> extractAddresses(final ProcessContext context, final FlowFile flowFile) throws ProcessException{
            try {
                return extractAddressesFromText(context.getProperty(ADDRESS_TEXT_PROPERTY).evaluateAttributeExpressions(flowFile).getValue());
            } catch (Exception e) {
                throw new ProcessException(e.toString());
            }
            
        }
    }

    public static class AddressPropertyAccess implements AddressAccessStrategy{

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
    }
}