package org.apache.plc4x.camel.s7;

import org.apache.camel.Exchange;
import org.apache.camel.Processor;
import org.apache.camel.impl.ScheduledPollConsumer;
import org.apache.plc4x.java.api.messages.PlcSimpleWriteRequest;
import org.apache.plc4x.java.api.types.ByteValue;

/**
 * The Awesome consumer.
 */
public class S7Consumer extends ScheduledPollConsumer {
    private final S7Endpoint endpoint;

    public S7Consumer(S7Endpoint endpoint, Processor processor) {
        super(endpoint, processor);
        this.endpoint = endpoint;
    }

    @Override
    protected int poll() throws Exception {
        Exchange exchange = endpoint.createExchange();

        exchange.getIn().setBody(new PlcSimpleWriteRequest(ByteValue.class, null, new ByteValue(Byte.valueOf((byte) 0xa0))));

        try {
            // send message to next processor in the route
            getProcessor().process(exchange);
            return 1; // number of messages polled
        } finally {
            // log exception if an exception occurred and was not handled
            if (exchange.getException() != null) {
                getExceptionHandler().handleException("Error processing exchange", exchange, exchange.getException());
            }
        }
    }
}
