package org.apache.plc4x.camel;

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
import org.apache.plc4x.java.api.messages.GenericPlcWriteRequest;
import org.apache.plc4x.java.api.messages.GenericPlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class PLC4XProducer extends DefaultAsyncProducer implements ShutdownAware {
    private static final Logger LOG = LoggerFactory.getLogger(PLC4XProducer.class);
    private PLC4XEndpoint endpoint;
    private PlcConnection plcConnection;

    public PLC4XProducer(PLC4XEndpoint endpoint) {
        super(endpoint);
        this.endpoint = endpoint;
        try {
            plcConnection = new PlcDriverManager().getConnection(endpoint.getEndpointUri().replaceFirst("plc4x:/?/?", ""));
            plcConnection.connect();
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
    }

    @SuppressWarnings("unchecked")
    public void process(Exchange exchange) throws Exception {
        Message in = exchange.getIn();
        Address address = in.getHeader(Constants.ADDRESS_HEADER, Address.class);
        Class<?> datatype = in.getHeader(Constants.DATATYPE_HEADER, Class.class);
        Object value = in.getBody(Object.class);
        GenericPlcWriteRequest plcSimpleWriteRequest = new GenericPlcWriteRequest(datatype, address, value);
        StreamUtils.streamOf(plcConnection.getWriter())
            .map(plcWriter -> plcWriter.write(plcSimpleWriteRequest))
            .forEach(plcWriteResponseCompletableFuture -> {
                try {
                    // FIXME: If I omit the cast to CompletableFuture the java compiler complains
                    GenericPlcWriteResponse response = (GenericPlcWriteResponse)
                        ((CompletableFuture) plcWriteResponseCompletableFuture).get();
                    if (exchange.getPattern().isOutCapable()) {
                        Message out = exchange.getOut();
                        out.copyFrom(exchange.getIn());
                        out.setBody(response);
                    } else {
                        in.setBody(response);
                    }
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
