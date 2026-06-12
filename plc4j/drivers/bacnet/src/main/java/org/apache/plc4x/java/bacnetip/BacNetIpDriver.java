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
package org.apache.plc4x.java.bacnetip;

import org.apache.plc4x.java.bacnetip.configuration.BacNetIpConfiguration;
import org.apache.plc4x.java.bacnetip.configuration.BacNetUdpTransportConfiguration;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.List;
import java.util.Optional;

public class BacNetIpDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "bacnet-ip";
    }

    @Override
    public String getProtocolName() {
        return "BACnet/IP";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return BacNetIpConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("udp".equals(transport.getTransportCode())) {
            return BacNetUdpTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("udp");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("udp", "test");
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new BacNetIpConnection((BacNetIpConfiguration) configuration, transportInstance, auditLog);
    }

}
