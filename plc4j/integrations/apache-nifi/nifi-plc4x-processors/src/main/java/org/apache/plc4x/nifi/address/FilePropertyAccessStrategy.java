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

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;

import java.nio.file.StandardOpenOption;
import java.nio.file.Path;
import java.util.Collection;
import java.util.List;
import java.util.Map;

import org.apache.nifi.components.AllowableValue;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.processor.ProcessContext;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.plc4x.java.DefaultPlcDriverManager;

import com.fasterxml.jackson.databind.ObjectMapper;

public class FilePropertyAccessStrategy extends BaseAccessStrategy {

    @Override
    public AllowableValue getAllowableValue() {
        return AddressesAccessUtils.ADDRESS_FILE;
    }

    @Override
    public List<PropertyDescriptor> getPropertyDescriptors() {
        return List.of(AddressesAccessUtils.ADDRESS_FILE_PROPERTY);
    }

    @Override
    public Map<String, String> extractAddressesFromResources(final ProcessContext context, final FlowFile flowFile) throws ProcessException{
        try {
            return extractAddressesFromFile(context.getProperty(AddressesAccessUtils.ADDRESS_FILE_PROPERTY).evaluateAttributeExpressions(flowFile).getValue());
        } catch (Exception e) {
            throw new ProcessException(e.toString());
        }
    }

    public static Map<String,String> extractAddressesFromFile(String fileName) throws IOException {
        ObjectMapper mapper = new ObjectMapper();

        Path filePath = Path.of(fileName);
        InputStream input = Files.newInputStream(filePath, StandardOpenOption.READ);
        
        return mapper.readerForMapOf(String.class).readValue(input);
    }

    public static class TagValidator extends BaseAccessStrategy.TagValidator {
        public TagValidator(DefaultPlcDriverManager manager) {
            super(manager);
        }

        @Override
        protected Collection<String> getTags(String input) throws Exception {
            return FilePropertyAccessStrategy.extractAddressesFromFile(input).values();
        } 
    }
}
