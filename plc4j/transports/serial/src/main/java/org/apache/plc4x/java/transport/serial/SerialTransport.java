/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.transport.serial;

import org.apache.plc4x.java.api.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.transport.Transport;

public class SerialTransport implements Transport, HasConfiguration<SerialTransportConfiguration> {

    private SerialTransportConfiguration configuration;

    @Override
    public String getTransportCode() {
        return "serial";
    }

    @Override
    public String getTransportName() {
        return "Serial Port Transport";
    }

    @Override
    public void setConfiguration(SerialTransportConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public ChannelFactory createChannelFactory(String transportConfig) {
        SerialSocketAddress socketAddress = new SerialSocketAddress(transportConfig);
        SerialChannelFactory serialChannelFactory = new SerialChannelFactory(socketAddress);
        if(configuration != null) {
            serialChannelFactory.setConfiguration(configuration);
        }
        return serialChannelFactory;
    }

    @Override
    public Class<? extends PlcTransportConfiguration> getTransportConfigType() {
        return DefaultSerialTransportConfiguration.class;
    }

}
