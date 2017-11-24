package org.apache.plc4x.camel.s7;

import org.apache.camel.AsyncCallback;
import org.apache.camel.Exchange;
import org.apache.camel.Message;
import org.apache.camel.ShutdownRunningTask;
import org.apache.camel.impl.DefaultAsyncProducer;
import org.apache.camel.spi.ShutdownAware;
import org.apache.plc4x.camel.util.StreamUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.Address;
import org.apache.plc4x.java.api.messages.PlcSimpleWriteRequest;
import org.apache.plc4x.java.api.messages.PlcSimpleWriteResponse;
import org.apache.plc4x.java.api.types.Value;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.ExecutionException;

/**
 * The Awesome producer.
 */
public class S7Producer extends DefaultAsyncProducer implements ShutdownAware {
    private static final Logger LOG = LoggerFactory.getLogger(S7Producer.class);
    private S7Endpoint endpoint;
    private PlcConnection plcConnection;

    public S7Producer(S7Endpoint endpoint) {
        super(endpoint);
        this.endpoint = endpoint;
        try {
            plcConnection = new PlcDriverManager().getConnection(endpoint.getEndpointUri());
            plcConnection.connect();
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
    }

    @SuppressWarnings("unchecked")
    public void process(Exchange exchange) throws Exception {
        Message in = exchange.getIn();
        Address address = in.getHeader(Constants.ADDRESS_HEADER, Address.class);
        Class<Value> datatype = in.getHeader(Constants.DATATYPE_HEADER, Class.class);
        Value value = in.getBody(Value.class);
        PlcSimpleWriteRequest<Value> plcSimpleWriteRequest = new PlcSimpleWriteRequest<Value>(datatype, address, value);
        StreamUtils.streamOf(plcConnection.getWriter())
            .map(plcWriter -> plcWriter.write(plcSimpleWriteRequest))
            .forEach(plcSimpleWriteResponseCompletableFuture -> {
                try {
                    PlcSimpleWriteResponse<Value> valuePlcSimpleWriteResponse = plcSimpleWriteResponseCompletableFuture.get();
                    in.setHeader(Constants.DATATYPE_HEADER, valuePlcSimpleWriteResponse.getDatatype());
                    in.setHeader(Constants.ADDRESS_HEADER, valuePlcSimpleWriteResponse.getAddress());
                    in.setHeader(Constants.SIZE_HEADER, valuePlcSimpleWriteResponse.getSize());
                    in.setBody(valuePlcSimpleWriteResponse.getValue());
                } catch (InterruptedException | ExecutionException e) {
                    throw new RuntimeException(e);
                }
            });
    }

    public boolean process(Exchange exchange, AsyncCallback callback) {
        try {
            process(exchange);
            Message out = exchange.getOut();
            out.copyFrom(exchange.getIn());
        } catch (Exception e) {
            exchange.setOut(null);
            exchange.setException(e);
        }
        callback.done(true);
        return true;
    }

    @Override
    public boolean deferShutdown(ShutdownRunningTask shutdownRunningTask) {
        switch (shutdownRunningTask) {
            case CompleteCurrentTaskOnly:
                break;
            case CompleteAllTasks:
                break;
        }
        try {
            plcConnection.close();
        } catch (Exception ignore) {
        }
        return false;
    }

    @Override
    public int getPendingExchangesSize() {
        return 0;
    }

    @Override
    public void prepareShutdown(boolean suspendOnly, boolean forced) {

    }
}
