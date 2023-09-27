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

import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.nifi.components.AllowableValue;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.plc4x.java.DefaultPlcDriverManager;


public class DynamicPropertyAccessStrategy extends BaseAccessStrategy{

    @Override
    public AllowableValue getAllowableValue() {
        return AddressesAccessUtils.ADDRESS_PROPERTY;
    }

    @Override
    public List<PropertyDescriptor> getPropertyDescriptors() {
        return List.of();
    }

    @Override
    public Map<String,String> extractAddressesFromResources(final ProcessContext context, final FlowFile flowFile) {
        return extractAddressesFromAttributes(context, flowFile);
    }

    private Map<String,String> extractAddressesFromAttributes(final ProcessContext context, final FlowFile flowFile) {
        Map<String,String> addressMap = new HashMap<>();

        context.getProperties().keySet().stream().filter(PropertyDescriptor::isDynamic).forEach(
            t -> addressMap.put(t.getName(), context.getProperty(t.getName()).evaluateAttributeExpressions(flowFile).getValue()));
        
        return addressMap; 
    }


    public static class TagValidator extends BaseAccessStrategy.TagValidator {
        public TagValidator(DefaultPlcDriverManager manager) {
            super(manager);
        }

        @Override
        protected Collection<String> getTags(String input) throws Exception {
            return List.of(input);
        }
    }

}
