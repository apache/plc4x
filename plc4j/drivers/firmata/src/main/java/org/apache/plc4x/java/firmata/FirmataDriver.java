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
package org.apache.plc4x.java.firmata;

import org.apache.plc4x.java.firmata.configuration.FirmataConfiguration;
import org.apache.plc4x.java.firmata.configuration.FirmataSerialTransportConfiguration;
import org.apache.plc4x.java.firmata.configuration.FirmataTcpTransportConfiguration;
import org.apache.plc4x.java.firmata.tag.FirmataTag;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.List;
import java.util.Optional;

public class FirmataDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "firmata";
    }

    @Override
    public String getProtocolName() {
        return "Firmata";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return FirmataConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("serial".equals(transport.getTransportCode())) {
            return FirmataSerialTransportConfiguration.class;
        }
        if ("tcp".equals(transport.getTransportCode())) {
            return FirmataTcpTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("serial");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        // serial: the classic Firmata-over-UART path
        // tcp:    StandardFirmataWiFi / StandardFirmataEthernet style boards,
        //         and a fallback when the serial transport can't open a
        //         socat-bridged PTY (e.g. on macOS)
        // test:   in-process loopback for the driver testsuite framework
        return List.of("serial", "tcp", "test");
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new FirmataConnection((FirmataConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    public FirmataTag prepareTag(String tagAddress) {
        return FirmataTag.of(tagAddress);
    }

}
