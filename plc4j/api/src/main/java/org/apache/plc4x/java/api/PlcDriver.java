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
package org.apache.plc4x.java.api;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.api.model.PlcField;

/**
 * General interface defining the minimal methods required for adding a new type of driver to the PLC4J system.
 *
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
     * Provides driver metadata.
     */
    default PlcDriverMetadata getMetadata() {
        return new PlcDriverMetadata() {
            @Override
            public boolean canDiscover() {
                return false;
            }
        };
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

    default PlcField prepareField(String query) {
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
