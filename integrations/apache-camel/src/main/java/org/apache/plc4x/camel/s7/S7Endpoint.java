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
package org.apache.plc4x.camel.s7;

import lombok.Getter;
import lombok.Setter;
import org.apache.camel.Component;
import org.apache.camel.Consumer;
import org.apache.camel.Processor;
import org.apache.camel.Producer;
import org.apache.camel.impl.DefaultEndpoint;
import org.apache.camel.spi.Metadata;
import org.apache.camel.spi.UriEndpoint;
import org.apache.camel.spi.UriPath;

@UriEndpoint(scheme = "s7", title = "S7", syntax = "s7:address/rack/slot",
    consumerClass = S7Consumer.class, label = "s7")
@Setter
@Getter
public class S7Endpoint extends DefaultEndpoint {

    @UriPath
    @Metadata(required = "true")
    String address;
    @UriPath
    @Metadata(required = "true")
    int rack;
    @UriPath
    @Metadata(required = "true")
    int slot;


    public S7Endpoint(String endpointUri, Component component) {
        super(endpointUri, component);
    }

    @Override
    public Producer createProducer() throws Exception {
        return new S7Producer(this);
    }

    @Override
    public Consumer createConsumer(Processor processor) throws Exception {
        return new S7Consumer(this, processor);
    }

    @Override
    public boolean isSingleton() {
        return true;
    }
}
