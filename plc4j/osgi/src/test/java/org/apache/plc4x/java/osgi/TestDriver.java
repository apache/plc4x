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
package org.apache.plc4x.java.osgi;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.spi.transports.api.TransportManager;
import org.apache.plc4x.java.spi.transports.api.TransportManagerAware;

/**
 * Driver listed in this module's test META-INF/services, so the {@link DriverActivator} finds it.
 */
public class TestDriver implements PlcDriver, TransportManagerAware {

    private TransportManager transportManager;

    @Override
    public String getProtocolCode() {
        return "osgi-test";
    }

    @Override
    public String getProtocolName() {
        return "OSGi Test Driver";
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        throw new PlcConnectionException("test driver does not connect");
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("test driver does not connect");
    }

    @Override
    public void setTransportManager(TransportManager transportManager) {
        this.transportManager = transportManager;
    }

    public TransportManager getTransportManager() {
        return transportManager;
    }

}
