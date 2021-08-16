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
package org.apache.plc4x.java.modbus_rtu.config;

import org.apache.plc4x.java.modbus.readwrite.ModbusConstants;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.transport.tcp.TcpTransportConfiguration;

public class ModbusRTUConfiguration implements Configuration, TcpTransportConfiguration {

    @ConfigurationParameter("request-timeout")
    @IntDefaultValue(5_000)
    private int requestTimeout;

    @ConfigurationParameter("address")
    @IntDefaultValue(1)
    private int address;

    public int getRequestTimeout() {
        return requestTimeout;
    }

    public void setRequestTimeout(int requestTimeout) {
        this.requestTimeout = requestTimeout;
    }

    public int getAddress() {
        return address;
    }

    public void setAddress(int address) {
        this.address = address;
    }

    @Override
    public int getDefaultPort() {
        return ModbusConstants.MODBUSTCPDEFAULTPORT;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

}
