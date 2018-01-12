package org.apache.plc4x.camel;

import org.apache.camel.Exchange;
import org.apache.camel.ExchangePattern;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.model.Address;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Optional;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import static org.mockito.Mockito.*;

class PLC4XProducerTest {

    private PLC4XProducer SUT;

    private Exchange testExchange;

    @BeforeEach
    void setUp() throws Exception {
        PLC4XEndpoint endpointMock = mock(PLC4XEndpoint.class, RETURNS_DEEP_STUBS);
        when(endpointMock.getEndpointUri()).thenReturn("plc4x:mock:10.10.10.1/1/1");
        PlcDriverManager plcDriverManagerMock = mock(PlcDriverManager.class, RETURNS_DEEP_STUBS);
        when(plcDriverManagerMock.getConnection(anyString()).getWriter())
            .thenReturn(Optional.of(mock(PlcWriter.class, RETURNS_DEEP_STUBS)));
        when(endpointMock.getPlcDriverManager()).thenReturn(plcDriverManagerMock);
        SUT = new PLC4XProducer(endpointMock);
        testExchange = mock(Exchange.class, RETURNS_DEEP_STUBS);
        when(testExchange.getIn().getHeader(eq(Constants.ADDRESS_HEADER), eq(Address.class)))
            .thenReturn(mock(Address.class));
    }

    @Test
    void process() throws Exception {
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOnly);
        SUT.process(testExchange);
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOut);
        SUT.process(testExchange);
        when(testExchange.getPattern()).thenReturn(ExchangePattern.OutOnly);
        SUT.process(testExchange);
    }

    @Test
    void process_Async() throws Exception {
        SUT.process(testExchange, doneSync -> { });
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOnly);
        SUT.process(testExchange, doneSync -> { });
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOut);
        SUT.process(testExchange, doneSync -> { });
        when(testExchange.getPattern()).thenReturn(ExchangePattern.OutOnly);
        SUT.process(testExchange, doneSync -> { });
    }

    @Test
    void doStop() throws Exception {
        SUT.doStop();
    }

}