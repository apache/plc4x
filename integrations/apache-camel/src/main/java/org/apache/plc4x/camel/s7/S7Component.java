package org.apache.plc4x.camel.s7;

import org.apache.camel.Endpoint;
import org.apache.camel.impl.DefaultComponent;

import java.util.Map;

public class S7Component extends DefaultComponent {

    @Override
    protected Endpoint createEndpoint(String uri, String remaining, Map<String, Object> parameters) throws Exception {
        S7Endpoint s7Endpoint = new S7Endpoint(uri, this);
        setProperties(s7Endpoint, parameters);
        return s7Endpoint;
    }
}
