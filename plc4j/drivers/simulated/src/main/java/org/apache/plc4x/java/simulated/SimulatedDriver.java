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
import org.apache.plc4x.java.simulated.configuration.SimulatedConfiguration;
import org.apache.plc4x.java.simulated.connection.SimulatedConnection;
import org.apache.plc4x.java.simulated.connection.SimulatedDevice;
import org.apache.plc4x.java.simulated.tag.SimulatedTag;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

/**
 * Test driver holding its state in the client process.
 *
 * <p>URL schema: {@code simulated:<device_name>}. Each call creates a fresh
 * in-memory device; values written through one connection do not survive its
 * close, and devices are not shared between connections.</p>
 *
 * <p>Tag addresses follow {@code <TYPE>/<name>:<plc-value-type>[<count>]?},
 * e.g. {@code RANDOM/foo:INT}, {@code STATE/bar:STRING[16]}, {@code STDOUT/log:STRING}.</p>
 */
public class SimulatedDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "simulated";
    }

    @Override
    public String getProtocolName() {
        return "Simulated PLC4X Datasource";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return SimulatedConfiguration.class;
    }

    @Override
    protected boolean canRead() {
        return true;
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected boolean canPing() {
        return true;
    }

    /**
     * The simulated driver runs entirely in-process and never opens a transport.
     * We override {@link #getConnection(String)} to parse the {@code simulated:<device_name>}
     * URL form directly — the inherited URI parser expects a {@code ://} delimiter
     * and a transport, neither of which apply here.
     */
    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        String prefix = getProtocolCode() + ":";
        if (connectionString == null || !connectionString.startsWith(prefix)) {
            throw new PlcConnectionException(
                "Invalid URL: expected '" + prefix + "<device_name>'");
        }
        String deviceName = connectionString.substring(prefix.length());
        if (deviceName.isEmpty()) {
            throw new PlcConnectionException("Invalid URL: no device name given.");
        }
        SimulatedDevice device = new SimulatedDevice(deviceName);
        return new SimulatedConnection(device,
            new SimulatedConfiguration(),
            AuditLog.builder().withSource(getProtocolCode()).build());
    }

    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication authentication)
        throws PlcConnectionException {
        if (authentication != null) {
            throw new PlcConnectionException("Simulated driver does not support authentication.");
        }
        return getConnection(connectionString);
    }

    /**
     * Unused — {@link #getConnection(String)} above takes the direct path that
     * doesn't go through {@link DriverBase}'s transport-aware factory. Kept so
     * the abstract contract is still satisfied.
     */
    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        throw new UnsupportedOperationException(
            "Simulated driver bypasses the transport-aware connection factory.");
    }

    @Override
    public SimulatedTag prepareTag(String tagAddress) {
        return SimulatedTag.of(tagAddress);
    }

}
