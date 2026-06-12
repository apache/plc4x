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
package org.apache.plc4x.java.utils.testutils.driver.internal;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Method;
import java.util.Arrays;

import java.util.Map;
import java.util.ServiceLoader;
import java.util.stream.Collectors;

/**
 * Manages PLC connections for testing purposes.
 */
public class ConnectionManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(ConnectionManager.class);

    /**
     * Creates a PLC connection with the given driver and parameters.
     *
     * @param driver     the driver protocol code
     * @param parameters the connection parameters
     * @return the PLC connection
     * @throws DriverTestsuiteException if connection fails
     */
    public PlcConnection getConnection(String driver, Map<String, String> parameters) {
        PlcConnection connection = createConnection(driver, parameters);
        try {
            connection.connect();
        } catch (PlcConnectionException e) {
            throw new DriverTestsuiteException("Failed to connect", e);
        }
        return connection;
    }

    /**
     * Creates a PLC connection without connecting.
     * This allows pre-loading the test transport before the connection sequence runs.
     *
     * @param driver     the driver protocol code
     * @param parameters the connection parameters
     * @return the PLC connection (not connected)
     * @throws DriverTestsuiteException if creation fails
     */
    public PlcConnection createConnection(String driver, Map<String, String> parameters) {
        try {
            // Extract transport parameter if present, default to "test" for driver test suites
            String transport = parameters.getOrDefault("transport", "test");
            String host = parameters.getOrDefault("host", "localhost");

            // Build connection URL
            // Format: protocol:transport://host or protocol://host?params
            String connectionUrl;
            if (transport != null) {
                // Use format: modbus-tcp:test://localhost
                connectionUrl = driver + ":" + transport + "://" + host;
            } else {
                // Build query parameters (excluding transport and host)
                connectionUrl = driver + "://" + host;
            }
            String parameterString = parameters.entrySet().stream()
                .filter(entry -> !entry.getKey().equals("transport") && !entry.getKey().equals("host"))
                .map(entry -> entry.getKey() + "=" + entry.getValue())
                .collect(Collectors.joining("&"));
            connectionUrl += (parameterString.isEmpty() ? "" : "?" + parameterString);

            LOGGER.debug("Creating connection with URL: {}", connectionUrl);

            // Load driver using ServiceLoader
            PlcDriver plcDriver = loadDriver(driver);
            if (plcDriver == null) {
                throw new DriverTestsuiteException("Driver not found for protocol: " + driver);
            }

            // Create connection without connecting
            return plcDriver.getConnection(connectionUrl);
        } catch (PlcConnectionException e) {
            throw new DriverTestsuiteException("Failed to create connection", e);
        }
    }

    /**
     * Extracts the transport instance from a PLC connection for testing.
     *
     * @param connection the PLC connection
     * @return the transport instance
     * @throws PlcRuntimeException if the connection doesn't expose a transport instance
     */
    public TransportInstance<?> getTransportInstance(PlcConnection connection) {
        try {
            // Use reflection to get the transport instance
            // Search in all public methods including inherited ones
            Method method = connection.getClass().getMethod("getTransportInstance");
            return (TransportInstance<?>) method.invoke(connection);
        } catch (NoSuchMethodException e) {
            // Log available methods for debugging
            String availableMethods = Arrays.stream(connection.getClass().getMethods())
                .map(Method::getName)
                .filter(name -> name.contains("ransport") || name.contains("Instance"))
                .collect(Collectors.joining(", "));
            throw new PlcRuntimeException("Connection must have a public getTransportInstance() method for testing. " +
                "Connection class: " + connection.getClass().getName() +
                ", Available methods containing 'ransport' or 'Instance': [" + availableMethods + "]", e);
        } catch (Exception e) {
            throw new PlcRuntimeException("Failed to call getTransportInstance() on connection. " +
                "Connection class: " + connection.getClass().getName(), e);
        }
    }

    /**
     * Loads a driver using Java's ServiceLoader mechanism.
     *
     * @param protocolCode the protocol code (e.g., "ads", "modbus-tcp")
     * @return the PlcDriver instance, or null if not found
     */
    private PlcDriver loadDriver(String protocolCode) {
        for (PlcDriver driver : ServiceLoader.load(PlcDriver.class)) {
            if (protocolCode.equals(driver.getProtocolCode())) {
                LOGGER.debug("Loaded driver: {} (protocol: {})", driver.getClass().getName(), driver.getProtocolCode());
                return driver;
            }
        }
        return null;
    }
}
