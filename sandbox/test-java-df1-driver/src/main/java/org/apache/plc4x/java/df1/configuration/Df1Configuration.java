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
package org.apache.plc4x.java.df1.configuration;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationProvider;
import org.apache.plc4x.java.transport.serial.DefaultSerialTransportConfiguration;

public class Df1Configuration implements Configuration, TransportConfigurationProvider {

    @ConfigurationParameter("local-addr")
    private short localAddr;

    @ConfigurationParameter("remote-addr")
    private short remoteAddr;

    @ComplexConfigurationParameter(prefix = "serial", defaultOverrides = {}, requiredOverrides = {})
    private DefaultSerialTransportConfiguration serialTransportConfiguration;

    public short getLocalAddr() {
        return localAddr;
    }

    public void setLocalAddr(short localAddr) {
        this.localAddr = localAddr;
    }

    public short getRemoteAddr() {
        return remoteAddr;
    }

    public void setRemoteAddr(short remoteAddr) {
        this.remoteAddr = remoteAddr;
    }

    public DefaultSerialTransportConfiguration getSerialTransportConfiguration() {
        return serialTransportConfiguration;
    }

    public void setSerialTransportConfiguration(DefaultSerialTransportConfiguration serialTransportConfiguration) {
        this.serialTransportConfiguration = serialTransportConfiguration;
    }

    @Override
    public TransportConfiguration getTransportConfiguration(String transportCode) {
        switch (transportCode) {
            case "serial":
                return serialTransportConfiguration;
        }
        return null;
    }

}
