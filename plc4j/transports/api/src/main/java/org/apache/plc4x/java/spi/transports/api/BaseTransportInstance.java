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

package org.apache.plc4x.java.spi.transports.api;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;

import java.util.Objects;

public abstract class BaseTransportInstance<T extends TransportConfiguration> implements TransportInstance<T> {

    private final T transportConfig;
    private final AuditLog auditLog;
    // The driver-config: the URL remainder after the transport-config (host/port),
    // set by the transport that parsed the address. Empty unless the transport
    // supplies one. See TransportInstance#getDriverConfig().
    private String driverConfig = "";

    public BaseTransportInstance(T transportConfig, AuditLog auditLog) {
        this.transportConfig = Objects.requireNonNull(transportConfig);
        this.auditLog = auditLog;
        auditLog.write(AuditLogEventType.SYSTEM, "Creating Transport with config", transportConfig);
    }

    public T getConfiguration() {
        return transportConfig;
    }

    @Override
    public String getDriverConfig() {
        return driverConfig;
    }

    /**
     * Sets the driver-config parsed from the connection URL. Transports call this
     * right after constructing the instance with the portion of the address they did
     * not consume (e.g. an OPC UA {@code /milo} path).
     */
    public void setDriverConfig(String driverConfig) {
        this.driverConfig = driverConfig == null ? "" : driverConfig;
    }

    public AuditLog getAuditLog() {
        return auditLog;
    }

}
