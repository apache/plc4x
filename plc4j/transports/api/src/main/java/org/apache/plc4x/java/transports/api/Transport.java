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

package org.apache.plc4x.java.transports.api;

import org.apache.plc4x.java.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

public interface Transport<T extends TransportConfiguration> {

    /**
     * @return the code used in a connection string to refer to this transport
     */
    String getTransportCode();

    /**
     * @return A human-readable name for this transport
     */
    String getTransportName();

    /**
     * Returns the configuration class for this transport.
     * This is mainly used for systems where drivers are programmatically included.
     *
     * @return the configuration class for this transport
     */
    default Class<T> getTransportConfigType() {
        return null;
    }

    /**
     * Returns an instance of the current transport for the given configuration.
     *
     * @param transportUrl the URL of the transport to connect to
     * @param configuration configuration for the transport
     * @return an instance of the current transport for the given configuration
     * @throws TransportException something went wrong
     */
    TransportInstance<T> createTransportInstance(String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException;

}
