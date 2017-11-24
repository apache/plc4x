package org.apache.plc4x.camel.s7;

import org.apache.camel.builder.RouteBuilder;
import org.apache.camel.component.mock.MockEndpoint;
import org.apache.camel.test.junit4.CamelTestSupport;
import org.apache.plc4x.java.api.types.ByteValue;
import org.apache.plc4x.java.s7.mina.model.types.MemoryArea;
import org.apache.plc4x.java.s7.model.S7Address;
import org.junit.Test;

public class S7ComponentTest extends CamelTestSupport {

    @Test
    public void testAwesome() throws Exception {
        MockEndpoint mock = getMockEndpoint("mock:result");
        mock.expectedMinimumMessageCount(1);

        assertMockEndpointsSatisfied();
    }

    @Override
    protected RouteBuilder createRouteBuilder() throws Exception {
        return new RouteBuilder() {
            public void configure() {
                from("s7://10.10.10.1/0/1")
                    .setHeader(Constants.ADDRESS_HEADER, constant(new S7Address(MemoryArea.INPUTS, (short) 0x44)))
                    .setHeader(Constants.DATATYPE_HEADER, constant(ByteValue.class))
                    .setBody(constant(0x0))
                    .to("s7://10.10.10.1/1/1")
                    .to("mock:result");
            }
        };
    }

}