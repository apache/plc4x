package org.apache.plc4x.camel;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.apache.camel.Component;
import org.apache.camel.Consumer;
import org.apache.camel.Processor;
import org.apache.camel.Producer;
import org.apache.camel.impl.DefaultEndpoint;
import org.apache.camel.spi.Metadata;
import org.apache.camel.spi.UriEndpoint;
import org.apache.camel.spi.UriPath;

@UriEndpoint(scheme = "plc4x", title = "PLC4X", syntax = "plc4x:driver:address", label = "plc4x")
@Setter
@Getter
@NoArgsConstructor
public class PLC4XEndpoint extends DefaultEndpoint {

    /**
     * The address for the PLC4X driver
     */
    @UriPath
    @Metadata(required = "true")
    String address;

    public PLC4XEndpoint(String endpointUri, Component component) {
        super(endpointUri, component);
    }

    @Override
    public Producer createProducer() throws Exception {
        return new PLC4XProducer(this);
    }

    @Override
    public Consumer createConsumer(Processor processor) throws Exception {
        throw new UnsupportedOperationException("The PLC4X endpoint doesn't support consumers.");
    }

    @Override
    public boolean isSingleton() {
        return true;
    }
}
