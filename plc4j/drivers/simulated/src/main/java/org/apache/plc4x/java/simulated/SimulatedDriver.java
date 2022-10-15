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
package org.apache.plc4x.java.simulated;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.simulated.connection.SimulatedConnection;
import org.apache.plc4x.java.simulated.connection.SimulatedDevice;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.simulated.field.SimulatedField;

/**
 * Test driver holding its state in the client process.
 * The URL schema is {@code simulated:<device_name>}.
 * Devices are created each time a connection is established and should not be reused.
 * Every device contains a random value generator accessible by address {@code random}.
 * Any value can be stored into test devices, however the state will be gone when connection is closed.
 */
public class SimulatedDriver implements PlcDriver {

    @Override
    public String getProtocolCode() {
        return "simulated";
    }

    @Override
    public String getProtocolName() {
        return "Simulated PLC4X Datasource";
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        // TODO: perform further checks
        String deviceName = url.substring(getProtocolCode().length() + 1);
        if (deviceName.isEmpty()) {
            throw new PlcConnectionException("Invalid URL: no device name given.");
        }
        SimulatedDevice device = new SimulatedDevice(deviceName);
        return new SimulatedConnection(device);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Test driver does not support authentication.");
    }

    @Override
    public SimulatedField prepareField(String query){
        return SimulatedField.of(query);
    }

}
