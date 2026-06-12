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
package org.apache.plc4x.java.iec608705104;

import org.apache.plc4x.java.iec608705104.configuration.Iec608705014Configuration;
import org.apache.plc4x.java.iec608705104.configuration.Iec608705014TcpTransportConfiguration;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104TagHandler;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import org.apache.plc4x.java.iec608705104.readwrite.Constants;

public class Iec60870514PlcDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "iec-60870-5-104";
    }

    @Override
    public String getProtocolName() {
        return "IEC 60870-5-104";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return Iec608705014Configuration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("tcp".equals(transport.getTransportCode())) {
            return Iec608705014TcpTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("tcp");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("tcp", "test");
    }

    @Override
    public Set<Integer> defaultPorts(String transportCode) {
        if ("tcp".equalsIgnoreCase(transportCode)) {
            return Set.of(Constants.DEFAULTPORT);
        }
        return Collections.emptySet();
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new Iec60870Connection((Iec608705014Configuration) configuration, transportInstance, auditLog);
    }

    @Override
    public Iec608705104Tag prepareTag(String tagAddress) {
        return (Iec608705104Tag) new Iec608705104TagHandler().parseTag(tagAddress);
    }

}
