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
import org.apache.plc4x.java.spi.InstanceFactory;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.parser.ConnectionParser;

public abstract class GeneratedDriverBase<BASE_PACKET extends Message> implements PlcDriver {

    protected abstract int getDefaultPortIPv4();

    protected abstract PlcFieldHandler getFieldHandler();

    protected abstract Class<? extends NettyChannelFactory> getTransportChannelFactory();

    protected abstract ProtocolStackConfigurer<BASE_PACKET> getStackConfigurer();

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        ConnectionParser parser = new ConnectionParser(getProtocolCode(), url);
        InstanceFactory instanceFactory = new InstanceFactory(parser);
        // CONFIGURATION configuration = parser.createConfiguration(getConfigurationClass());

        // Create Instance of Transport
        NettyChannelFactory transport = instanceFactory.createInstance(getTransportChannelFactory());

        return new DefaultNettyPlcConnection<>(
            instanceFactory,
            parser.getSocketAddress(getDefaultPortIPv4()),
            transport,
            true,
            getFieldHandler(),
            getStackConfigurer());
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Authentication not supported.");
    }

}
