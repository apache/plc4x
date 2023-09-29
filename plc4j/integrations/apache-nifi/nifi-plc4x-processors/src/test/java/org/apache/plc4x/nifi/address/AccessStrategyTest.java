package org.apache.plc4x.nifi.address;

import static org.junit.jupiter.api.Assertions.assertTrue;
import java.util.Map;

import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.apache.plc4x.nifi.Plc4xSourceProcessor;
import org.apache.plc4x.nifi.util.Plc4xCommonTest;
import org.junit.jupiter.api.Test;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class AccessStrategyTest {

    private TestRunner testRunner; 

    // Test correct cacheing of addresses in case 
    @Test
    public void testDynamicPropertyAccessStrategy() {

        DynamicPropertyAccessStrategy testObject = new DynamicPropertyAccessStrategy();
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        assert testObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_PROPERTY);
        assert testObject.getPropertyDescriptors().isEmpty();
        
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(k, v));
		
        FlowFile flowFile = testRunner.enqueue("");
        
        Map<String, String> values = testObject.extractAddresses(testRunner.getProcessContext(), flowFile);

        System.out.println(values);
        assertTrue(testObject.getCachedAddresses().equals(values));
        assertTrue(testObject.getCachedAddresses().equals(Plc4xCommonTest.getAddressMap()));
    }

    @Test
    public void testDynamicPropertyAccessStrategyIncorrect() {

        DynamicPropertyAccessStrategy testObject = new DynamicPropertyAccessStrategy();
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        assert testObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_PROPERTY);
        assert testObject.getPropertyDescriptors().isEmpty();
        
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(k, "no an correct address"));

        testRunner.assertNotValid();
    }

    @Test
    public void testTextPropertyAccessStrategy() throws JsonProcessingException {

        TextPropertyAccessStrategy testObject = new TextPropertyAccessStrategy();
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        
        assert testObject.getAllowableValue().equals(AddressesAccessUtils.ADDRESS_TEXT);
        assert testObject.getPropertyDescriptors().contains(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY);
        
        testRunner.setProperty(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY, new ObjectMapper().writeValueAsString(Plc4xCommonTest.getAddressMap()).toString());
		
        FlowFile flowFile = testRunner.enqueue("");
        
        Map<String, String> values = testObject.extractAddresses(testRunner.getProcessContext(), flowFile);

        System.out.println(values);
        assertTrue(testObject.getCachedAddresses().equals(values));
        assertTrue(testObject.getCachedAddresses().equals(Plc4xCommonTest.getAddressMap()));
    }

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
}
