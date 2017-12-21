/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.camel;

import org.apache.camel.Component;
import org.apache.camel.Consumer;
import org.apache.camel.Processor;
import org.apache.camel.Producer;
import org.apache.camel.impl.DefaultEndpoint;
import org.apache.camel.spi.Metadata;
import org.apache.camel.spi.UriEndpoint;
import org.apache.camel.spi.UriPath;
import org.apache.plc4x.java.PlcDriverManager;

@UriEndpoint(scheme = "plc4x", title = "PLC4X", syntax = "plc4x:driver:address", label = "plc4x")
public class PLC4XEndpoint extends DefaultEndpoint {

    /**
     * The name 0f the PLC4X driver
     */
    @UriPath
    @Metadata(required = "true")
    @SuppressWarnings("unused")
    String driver;

    /**
     * The address for the PLC4X driver
     */
    @UriPath
    @Metadata(required = "true")
    @SuppressWarnings("unused")
    String address;

    final PlcDriverManager plcDriverManager;

    public PLC4XEndpoint(String endpointUri, Component component) {
        super(endpointUri, component);
        plcDriverManager = new PlcDriverManager();
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
