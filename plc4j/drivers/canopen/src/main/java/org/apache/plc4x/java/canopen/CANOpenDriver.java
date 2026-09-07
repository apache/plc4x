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
package org.apache.plc4x.java.canopen;

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.canopen.configuration.CANOpenConfiguration;
import org.apache.plc4x.java.canopen.tag.CANOpenTag;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.List;
import java.util.Optional;

public class CANOpenDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "canopen";
    }

    @Override
    public String getProtocolName() {
        return "CANopen";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return CANOpenConfiguration.class;
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("can-socketcan");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("can-socketcan", "can-virtualcan", "test");
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
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new CANOpenConnection((CANOpenConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    public PlcTag prepareTag(String tagAddress) {
        return CANOpenTag.of(tagAddress);
    }

}
