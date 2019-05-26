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

import org.apache.camel.*;
import org.apache.camel.impl.DefaultEndpoint;
import org.apache.camel.spi.Metadata;
import org.apache.camel.spi.UriEndpoint;
import org.apache.camel.spi.UriParam;
import org.apache.camel.spi.UriPath;
import org.apache.plc4x.java.PlcDriverManager;

import java.util.Objects;

@UriEndpoint(scheme = "plc4x", title = "PLC4X", syntax = "plc4x:driver", label = "plc4x")
public class Plc4XEndpoint extends DefaultEndpoint {

    /**
     * The name of the PLC4X driver
     */
    @UriPath(label = "common")
    @Metadata(required = "true")
    @SuppressWarnings("unused")
    private String driver;

    /**
     * The address for the PLC4X driver
     */
    @UriParam
    @Metadata(required = "false")
    @SuppressWarnings("unused")
    private String address;

    /**
     * TODO: document me
     */
    @UriParam
    @Metadata(required = "false")
    @SuppressWarnings("unused")
    private Class dataType;

    private final PlcDriverManager plcDriverManager;

    public Plc4XEndpoint(String endpointUri, Component component) {
        super(endpointUri, component);
        plcDriverManager = new PlcDriverManager();
    }

    @Override
    public Producer createProducer() throws Exception {
        return new Plc4XProducer(this);
    }

    @Override
    public Consumer createConsumer(Processor processor) throws Exception {
        return new Plc4XConsumer(this, processor);
    }

    @Override
    public PollingConsumer createPollingConsumer() throws Exception {
        return new Plc4XPollingConsumer(this);
    }

    @Override
    public boolean isSingleton() {
        return true;
    }

    public PlcDriverManager getPlcDriverManager() {
        return plcDriverManager;
    }

    public String getDriver() {
        return driver;
    }

    public void setDriver(String driver) {
        this.driver = driver;
    }

    public String getAddress() {
        return address;
    }

    public void setAddress(String address) {
        this.address = address;
    }

    public Class getDataType() {
        return dataType;
    }

    public void setDataType(Class dataType) {
        this.dataType = dataType;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Plc4XEndpoint)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        Plc4XEndpoint that = (Plc4XEndpoint) o;
        return Objects.equals(getDriver(), that.getDriver()) &&
            Objects.equals(getAddress(), that.getAddress()) &&
            Objects.equals(getDataType(), that.getDataType()) &&
            Objects.equals(getPlcDriverManager(), that.getPlcDriverManager());
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), getDriver(), getAddress(), getDataType(), getPlcDriverManager());
    }

}
