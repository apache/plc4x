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

import static org.junit.jupiter.api.Assertions.assertTrue;
import java.io.IOException;
import java.util.Map;

import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.apache.plc4x.nifi.Plc4xSourceProcessor;
import org.apache.plc4x.nifi.util.Plc4xCommonTest;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockedStatic;
import org.mockito.Mockito;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class AccessStrategyTest {

    @Mock
    FilePropertyAccessStrategy testFileObject = new FilePropertyAccessStrategy();

    private TestRunner testRunner; 

    // Tests that addresses in dynamic properties are read correctly and addresses are cached if no EL is used
    @Test
    public void testDynamicPropertyAccessStrategy() {

        DynamicPropertyAccessStrategy testObject = new DynamicPropertyAccessStrategy();
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        assert testObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_PROPERTY);
        assert testObject.getPropertyDescriptors().isEmpty();
        
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(k, v));
		
        FlowFile flowFile = testRunner.enqueue("");
        
        Map<String, String> values = testObject.extractAddresses(testRunner.getProcessContext(), flowFile);

        assertTrue(testObject.getCachedAddresses().equals(values));
        assertTrue(testObject.getCachedAddresses().equals(Plc4xCommonTest.getAddressMap()));
    }

    // Tests incorrect address detection on dynamic properties
    @Test
    public void testDynamicPropertyAccessStrategyIncorrect() {
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(k, "no an correct address"));

        testRunner.assertNotValid();
    }

    // Tests that if EL is present in dynamic properties the processor is valid
    @Test
    public void testDynamicPropertyAccessStrategyELPresent() {
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        testRunner.setProperty(Plc4xSourceProcessor.PLC_CONNECTION_STRING, "simulated://127.0.0.1");
        
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(k, "${attribute}"));

        testRunner.assertValid();
    }

    // Tests that addresses in text property are read correctly and addresses are cached if no EL is used
    @Test
    public void testTextPropertyAccessStrategy() throws JsonProcessingException {

        TextPropertyAccessStrategy testObject = new TextPropertyAccessStrategy();
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        assert testObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_TEXT);
        assert testObject.getPropertyDescriptors().contains(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY);
        
        testRunner.setProperty(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY, new ObjectMapper().writeValueAsString(Plc4xCommonTest.getAddressMap()).toString());
		
        FlowFile flowFile = testRunner.enqueue("");
        
        Map<String, String> values = testObject.extractAddresses(testRunner.getProcessContext(), flowFile);

        assertTrue(testObject.getCachedAddresses().equals(values));
        assertTrue(testObject.getCachedAddresses().equals(Plc4xCommonTest.getAddressMap()));
    }

    

    // Tests incorrect address detection on text property
    @Test
    public void testTextPropertyAccessStrategyIncorrect() {

        TextPropertyAccessStrategy testObject = new TextPropertyAccessStrategy();
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        assert testObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_TEXT);
        assert testObject.getPropertyDescriptors().contains(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY);
        
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY.getName(), "no an correct address"));

        testRunner.assertNotValid();

        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY.getName(), "{\"neither\":\"this one\"}"));

        testRunner.assertNotValid();
    }

    // Tests that if EL is present in text property the processor is valid 
    @Test
    public void testTextPropertyAccessStrategyELPresent() {

        TextPropertyAccessStrategy testObject = new TextPropertyAccessStrategy();
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);

        testRunner.setProperty(Plc4xSourceProcessor.PLC_CONNECTION_STRING, "simulated://127.0.0.1");
        
        assert testObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_TEXT);
        assert testObject.getPropertyDescriptors().contains(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY);
        
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY.getName(), "${attribute}"));

        testRunner.assertValid();
    }

    // Tests that addresses in file are read correctly and addresses are cached if no EL is used
    @Test
    public void testFilePropertyAccessStrategy() throws IOException {

        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);

        assert testFileObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_FILE);
        assert testFileObject.getPropertyDescriptors().contains(AddressesAccessUtils.ADDRESS_FILE_PROPERTY);


        testRunner.setProperty(AddressesAccessUtils.ADDRESS_FILE_PROPERTY, "file");

        try (MockedStatic<FilePropertyAccessStrategy> staticMock = Mockito.mockStatic(FilePropertyAccessStrategy.class)) {
            staticMock.when(() -> FilePropertyAccessStrategy.extractAddressesFromFile("file"))
                .thenReturn(Plc4xCommonTest.getAddressMap());


            FlowFile flowFile = testRunner.enqueue("");
            Map<String, String> values = testFileObject.extractAddresses(testRunner.getProcessContext(), flowFile);

            assertTrue(testFileObject.getCachedAddresses().equals(values));
            assertTrue(testFileObject.getCachedAddresses().equals(Plc4xCommonTest.getAddressMap()));
        }
    }

    // Tests incorrect address detection on file
    @Test
    public void testFilePropertyAccessStrategyIncorrect() throws IOException {

        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        assert testFileObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_FILE);
        assert testFileObject.getPropertyDescriptors().contains(AddressesAccessUtils.ADDRESS_FILE_PROPERTY);
        
        testRunner.setProperty(AddressesAccessUtils.ADDRESS_FILE_PROPERTY, "file");

        try (MockedStatic<FilePropertyAccessStrategy> staticMock = Mockito.mockStatic(FilePropertyAccessStrategy.class)) {
            staticMock.when(() -> FilePropertyAccessStrategy.extractAddressesFromFile("file"))
                .thenReturn(Map.of("not", "a correct address"));

            testRunner.assertNotValid();
        }
    }

    // Tests that if EL is present in file the processor is valid 
    @Test
    public void testFilePropertyAccessStrategyELPresent() throws IOException {

        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);

        testRunner.setProperty(Plc4xSourceProcessor.PLC_CONNECTION_STRING, "simulated://127.0.0.1");
        
        assert testFileObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_FILE);
        assert testFileObject.getPropertyDescriptors().contains(AddressesAccessUtils.ADDRESS_FILE_PROPERTY);
        
        testRunner.setProperty(AddressesAccessUtils.ADDRESS_FILE_PROPERTY, "file");

        try (MockedStatic<FilePropertyAccessStrategy> staticMock = Mockito.mockStatic(FilePropertyAccessStrategy.class)) {
            staticMock.when(() -> FilePropertyAccessStrategy.extractAddressesFromFile("file"))
                .thenReturn(Map.of("EL in use", "${attribute}"));

            testRunner.assertValid();
        }
    }
}
