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
package org.apache.plc4x.java.api;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.api.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.api.model.PlcTag;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

/**
 * General interface defining the minimal methods required for adding a new type of driver to the PLC4J system.
 * <br>
 * <b>Note that each driver has to add a service file called org.apache.plc4x.java.spi.PlcDriver to
 * src/main/resources/META-INF which contains the fully qualified classname in order to get loaded
 * by the PlcDriverManager instances.</b>
 */
public interface PlcDriver {

    /**
     * @return code of the implemented protocol. This is usually a lot shorter than the String returned by @see #getProtocolName().
     */
    String getProtocolCode();

    /**
     * @return name of the implemented protocol.
     */
    String getProtocolName();

    /**
     * @return the type of the Configuration used by this driver.
     */
    Class<? extends PlcConnectionConfiguration> getConfigurationType();

    default Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationType(String transportCode) {
        return Optional.empty();
    }

    /**
     * @return Provides driver metadata.
     */
    default PlcDriverMetadata getMetadata() {
        return () -> false;
    }

    /**
     * @return Optional that allows returning the transport code of a default transport.
     */
    default Optional<String> getDefaultTransportCode() {
        return Optional.empty();
    }

    /**
     * @return List of explicitly supported transport codes.
     */
    default List<String> getSupportedTransportCodes() {
        return Collections.emptyList();
    }

    /**
     * Connects to a PLC using the given plc connection string.
     *
     * @param url plc connection string.
     * @return PlcConnection object.
     * @throws PlcConnectionException an exception if the connection attempt failed.
     */
    PlcConnection getConnection(String url) throws PlcConnectionException;

    /**
     * Connects to a PLC using the given plc connection string using given authentication credentials.
     *
     * @param url            plc connection string.
     * @param authentication authentication credentials.
     * @return PlcConnection object.
     * @throws PlcConnectionException an exception if the connection attempt failed.
     */
    PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException;

    default PlcTag prepareTag(String tagAddress) {
        throw new PlcNotImplementedException("Not implemented for " + getProtocolName());
    }

    /**
     * @return discovery request builder.
     * @throws PlcUnsupportedOperationException if the connection does not support subscription
     */
    default PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        throw new PlcNotImplementedException("Not implemented for " + getProtocolName());
    }

}
