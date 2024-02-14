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
package org.apache.plc4x.java.modbus.ascii.config;

import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;

public class ModbusAsciiConfiguration implements PlcConnectionConfiguration {

    @ConfigurationParameter("request-timeout")
    @IntDefaultValue(5_000)
    @Description("Default timeout for all types of requests.")
    private int requestTimeout;

    @ConfigurationParameter("unit-identifier")
    @IntDefaultValue(1)
    @Description("Unit-identifier that identifies the target PLC (On RS485 multiple Modbus Devices can be listening). Defaults to 1.")
    private int unitIdentifier;

    public int getRequestTimeout() {
        return requestTimeout;
    }

    public void setRequestTimeout(int requestTimeout) {
        this.requestTimeout = requestTimeout;
    }

    public int getUnitIdentifier() {
        return unitIdentifier;
    }

    public void setUnitIdentifier(int unitIdentifier) {
        this.unitIdentifier = unitIdentifier;
    }

    @Override
    public String toString() {
        return "ModbusAsciiConfiguration{" +
            "requestTimeout=" + requestTimeout +
            ", unitIdentifier=" + unitIdentifier +
            '}';
    }

}
