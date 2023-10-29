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
package org.apache.plc4x.java.cbus.configuration;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationProvider;

public class CBusConfiguration implements Configuration, TransportConfigurationProvider {

    @ConfigurationParameter("srchk")
    @BooleanDefaultValue(false)
    public boolean srchk = false;

    @ComplexConfigurationParameter(prefix = "tcp", defaultOverrides = {}, requiredOverrides = {})
    private CBusTcpTransportConfiguration tcpTransportConfiguration;

    public boolean isSrchk() {
        return srchk;
    }

    public void setSrchk(boolean srchk) {
        this.srchk = srchk;
    }

    public CBusTcpTransportConfiguration getTcpTransportConfiguration() {
        return tcpTransportConfiguration;
    }

    public void setTcpTransportConfiguration(CBusTcpTransportConfiguration tcpTransportConfiguration) {
        this.tcpTransportConfiguration = tcpTransportConfiguration;
    }

    @Override
    public TransportConfiguration getTransportConfiguration(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return tcpTransportConfiguration;
        }
        return null;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            "srchk=" + srchk +
            '}';
    }

}
