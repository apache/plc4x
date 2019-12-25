/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.spi.connection;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.parser.ConnectionParser;

public abstract class GeneratedDriverBase<CONFIGURATION, BASE_PACKET extends Message> implements PlcDriver {

    protected abstract int getDefaultPortIPv4();

    protected abstract Class<CONFIGURATION> getConfigurationClass();

    protected abstract PlcFieldHandler getFieldHandler();

    protected abstract Class<? extends NettyChannelFactory> getTransportChannelFactory();

    protected abstract ProtocolStackConfigurer<BASE_PACKET> getStackConfigurer(CONFIGURATION configuration);

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        ConnectionParser parser = new ConnectionParser(getProtocolCode(), url);
        CONFIGURATION configuration = parser.createConfiguration(getConfigurationClass());

        // Create Instance of Transport
        NettyChannelFactory transport;
        try {
            transport = getTransportChannelFactory().newInstance();
        } catch (InstantiationException | IllegalAccessException e) {
            throw new IllegalStateException("Cannot Instantiate Transport '"
                + getTransportChannelFactory().getSimpleName()
                + "'. Cannot access Default no Args Constructor.", e);
        }
        // Set all Properties
        transport.setProperties(parser.getProperties());

        return new DefaultNettyPlcConnection<>(
            parser.getSocketAddress(getDefaultPortIPv4()),
            transport,
            true,
            getFieldHandler(),
            getStackConfigurer(configuration));
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Basic S7 connections don't support authentication (NG).");
    }
}
