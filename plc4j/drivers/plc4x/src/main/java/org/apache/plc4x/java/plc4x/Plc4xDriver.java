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
package org.apache.plc4x.java.plc4x;

import org.apache.plc4x.java.plc4x.config.Plc4xConfiguration;
import org.apache.plc4x.java.plc4x.config.Plc4xTcpTransportConfiguration;
import org.apache.plc4x.java.plc4x.config.Plc4xTlsTransportConfiguration;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.List;
import java.util.Optional;

public class Plc4xDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "plc4x";
    }

    @Override
    public String getProtocolName() {
        return "PLC4X (Proxy-Protocol)";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return Plc4xConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("tcp".equals(transport.getTransportCode())) {
            return Plc4xTcpTransportConfiguration.class;
        }
        if ("tls".equals(transport.getTransportCode())) {
            return Plc4xTlsTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        // TLS is the default transport - the proxy carries credentials, so the channel should
        // be encrypted by default. Plaintext TCP remains available as an explicit opt-in
        // (e.g. "plc4x:tcp://...") for trusted networks or testing.
        return Optional.of("tls");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("tls", "tcp", "test");
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
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new Plc4xConnection((Plc4xConfiguration) configuration, transportInstance, auditLog);
    }

}
