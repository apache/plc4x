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
package org.apache.plc4x.java.osgi;

import org.apache.plc4x.java.DriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.BundleContext;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;

import java.net.URI;
import java.net.URISyntaxException;
import java.util.Collection;

public class OsgiDriverManager implements DriverManager {

    public static final String PROTOCOL_NAME = "org.apache.plc4x.driver.name";
    public static final String PROTOCOL_CODE = "org.apache.plc4x.driver.code";

    private final BundleContext context;

    OsgiDriverManager(BundleContext context) {
        this.context = context;
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        PlcDriver driver = getDriver(url);
        return driver.getConnection(url);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        PlcDriver driver = getDriver(url);
        return driver.getConnection(url, authentication);
    }

    private PlcDriver getDriver(String url) throws PlcConnectionException {
        try {
            URI connectionUri = new URI(url);
            String protocol = connectionUri.getScheme();

            Collection<ServiceReference<PlcDriver>> references = context.getServiceReferences(PlcDriver.class, null);

            for (ServiceReference<PlcDriver> driverReference : references) {
                Object property = driverReference.getProperty(PROTOCOL_CODE);
                if (protocol.equals(property)) {
                    return context.getService(driverReference);
                }
            }
            throw new PlcConnectionException("Unable to find driver for protocol '" + protocol + "'");
        } catch (InvalidSyntaxException e) {
            throw new PlcConnectionException("Could not locate driver for url: " + url, e);
        } catch (URISyntaxException e) {
            throw new PlcConnectionException("Malformed connection string: " + url, e);
        }
    }

    void close() {

    }
}