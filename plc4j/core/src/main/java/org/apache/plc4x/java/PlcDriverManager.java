/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java;

import org.apache.plc4x.java.spi.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import java.net.URI;
import java.net.URISyntaxException;
import java.util.HashMap;
import java.util.Map;
import java.util.ServiceLoader;

public class PlcDriverManager {

    private Map<String, PlcDriver> driverMap = null;

    public PlcDriverManager() {
        this(Thread.currentThread().getContextClassLoader());
    }

    public PlcDriverManager(ClassLoader classLoader) {
        driverMap = new HashMap<>();
        ServiceLoader<PlcDriver> plcDriverLoader = ServiceLoader.load(PlcDriver.class, classLoader);
        for (PlcDriver driver : plcDriverLoader) {
            if (driverMap.containsKey(driver.getProtocolCode())) {
                throw new IllegalStateException(
                    "Multiple driver implementations available for protocol code '" +
                        driver.getProtocolCode() + "'");
            }
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
        PlcDriver driver = getDriver(url);
        PlcConnection connection = driver.connect(url);
        connection.connect();
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
        PlcDriver driver = getDriver(url);
        PlcConnection connection = driver.connect(url, authentication);
        connection.connect();
        return connection;
    }

    private PlcDriver getDriver(String url) throws PlcConnectionException {
        try {
            URI connectionUri = new URI(url);
            String protocol = connectionUri.getScheme();
            PlcDriver driver = driverMap.get(protocol);
            if (driver == null) {
                throw new PlcConnectionException("Unable to find driver for protocol '" + protocol + "'");
            }
            return driver;
        } catch (URISyntaxException e) {
            throw new PlcConnectionException("Invalid plc4j connection string '" + url + "'", e);
        }
    }

}
