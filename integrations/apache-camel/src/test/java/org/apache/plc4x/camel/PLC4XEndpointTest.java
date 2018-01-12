package org.apache.plc4x.camel;

import org.apache.camel.Component;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.mockito.Mockito.mock;

class PLC4XEndpointTest {

    PLC4XEndpoint SUT;

    @BeforeEach
    void setUp() throws Exception {
        SUT = new PLC4XEndpoint("plc4x:mock:10.10.10.1/1/1", mock(Component.class));
    }

    @Test
    void createProducer() throws Exception {
        Assertions.assertNotNull(SUT.createProducer());
    }

    @Test
    void createConsumer() throws Exception {
        Assertions.assertThrows(UnsupportedOperationException.class, () -> SUT.createConsumer(null));
    }

    @Test
    void isSingleton() throws Exception {
        Assertions.assertTrue(SUT.isSingleton());
    }

}