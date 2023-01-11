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
package org.apache.plc4x.java;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.PlcDriver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.URI;
import java.net.URISyntaxException;
import java.util.HashMap;
import java.util.Map;
import java.util.ServiceLoader;
import java.util.Set;

public class PlcDriverManager implements PlcConnectionManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcDriverManager.class);

    protected ClassLoader classLoader;

    private Map<String, PlcDriver> driverMap;

    public PlcDriverManager() {
        this(Thread.currentThread().getContextClassLoader());
    }

    public PlcDriverManager(ClassLoader classLoader) {
        LOGGER.info("Instantiating new PLC Driver Manager with class loader {}", classLoader);
        this.classLoader = classLoader;
        driverMap = new HashMap<>();
        ServiceLoader<PlcDriver> plcDriverLoader = ServiceLoader.load(PlcDriver.class, classLoader);
        LOGGER.info("Registering available drivers...");
        for (PlcDriver driver : plcDriverLoader) {
            if (driverMap.containsKey(driver.getProtocolCode())) {
                throw new IllegalStateException(
                    "Multiple driver implementations available for protocol code '" +
                        driver.getProtocolCode() + "'");
            }
            LOGGER.info("Registering driver for Protocol {} ({})", driver.getProtocolCode(), driver.getProtocolName());
            driverMap.put(driver.getProtocolCode(), driver);
        }
    }

    /**
     * Connects to a PLC using the given plc connection string.
     *
     * @param url plc connection string.
     * @return PlcConnection object.
     * @throws PlcConnectionException an exception if the connection attempt failed.
     */
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        ClassLoader originalClassLoader = Thread.currentThread().getContextClassLoader();
        Thread.currentThread().setContextClassLoader(this.classLoader);
        PlcConnection connection;
        try {
            PlcDriver driver = getDriverForUrl(url);
            connection = driver.getConnection(url);
            connection.connect();
        } finally {
            Thread.currentThread().setContextClassLoader(originalClassLoader);
        }
        return connection;
    }

    /**
     * Connects to a PLC using the given plc connection string using given authentication credentials.
     *
     * @param url            plc connection string.
     * @param authentication authentication credentials.
     * @return PlcConnection object.
     * @throws PlcConnectionException an exception if the connection attempt failed.
     */
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        ClassLoader originalClassLoader = Thread.currentThread().getContextClassLoader();
        Thread.currentThread().setContextClassLoader(this.classLoader);
        PlcConnection connection;
        try {
            PlcDriver driver = getDriverForUrl(url);
            connection = driver.getConnection(url, authentication);
            connection.connect();
        } finally {
            Thread.currentThread().setContextClassLoader(originalClassLoader);
        }
        return connection;
    }

    /**
     * Returns the codes of all of the drivers which are currently registered at the PlcDriverManager
     * @return Set of driver codes for all drivers registered
     */
    public Set<String> listDrivers() {
        return driverMap.keySet();
    }

    /**
     * Returns suitable driver for protocol or throws an Exception.
     * @param protocolCode protocol code identifying the driver
     * @return Driver instance for the given protocol
     * @throws PlcConnectionException If no Suitable Driver can be found
     */
    public PlcDriver getDriver(String protocolCode) throws PlcConnectionException {
        PlcDriver driver = driverMap.get(protocolCode);
        if (driver == null) {
            throw new PlcConnectionException("Unable to find driver for protocol '" + protocolCode + "'");
        }
        return driver;
    }

    /**
     * Returns suitable driver for a given plc4x connection url or throws an Exception.
     * @param url Uri to use
     * @return Driver instance for the given url
     * @throws PlcConnectionException If no Suitable Driver can be found
     */
    public PlcDriver getDriverForUrl(String url) throws PlcConnectionException {
        try {
            URI connectionUri = new URI(url);
            String protocol = connectionUri.getScheme();
            return getDriver(protocol);
        } catch (URISyntaxException e) {
            throw new PlcConnectionException("Invalid plc4j connection string '" + url + "'", e);
        }
    }

}