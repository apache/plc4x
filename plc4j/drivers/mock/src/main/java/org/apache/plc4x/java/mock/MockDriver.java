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
package org.apache.plc4x.java.mock;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.mock.configuration.MockConfiguration;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.tag.MockTag;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Mocking driver used by other modules' unit tests.
 *
 * <p>URL schema: {@code mock:<device_name>}. The driver keeps a per-driver-instance
 * cache of connections keyed by device name, so two callers asking for the same
 * device get the same {@link MockConnection} back — that's how test code
 * configures a connection (via {@link MockConnection#setDevice}) and then later
 * fetches the same instance to assert on it.</p>
 *
 * <p>Like the simulated driver this bypasses {@link DriverBase}'s transport-aware
 * URI parser by overriding {@link #getConnection(String)} directly — there is no
 * transport.</p>
 */
public class MockDriver extends DriverBase {

    private final Map<String, MockConnection> connectionMap = new ConcurrentHashMap<>();

    @Override
    public String getProtocolCode() {
        return "mock";
    }

    @Override
    public String getProtocolName() {
        return "Mock Protocol Implementation";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return MockConfiguration.class;
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
    protected boolean canBrowse() {
        return true;
    }

    @Override
    protected boolean canPing() {
        return true;
    }

    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        return getConnection(connectionString, null);
    }

    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication authentication)
        throws PlcConnectionException {
        String prefix = getProtocolCode() + ":";
        if (connectionString == null || !connectionString.startsWith(prefix)) {
            throw new PlcConnectionException(
                "Invalid URL: expected '" + prefix + "<device_name>'");
        }
        String deviceName = connectionString.substring(prefix.length());
        if (deviceName.isEmpty()) {
            throw new PlcConnectionException("Invalid URL: no device name given.");
        }
        return connectionMap.computeIfAbsent(deviceName, name ->
            new MockConnection(authentication, new MockConfiguration(),
                AuditLog.builder().withSource(getProtocolCode()).build()));
    }

    /**
     * Unused — {@link #getConnection(String)} above takes the direct path that
     * doesn't go through {@link DriverBase}'s transport-aware factory.
     */
    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        throw new UnsupportedOperationException(
            "Mock driver bypasses the transport-aware connection factory.");
    }

    @Override
    public MockTag prepareTag(String tagAddress) {
        return MockTag.of(tagAddress);
    }

}
