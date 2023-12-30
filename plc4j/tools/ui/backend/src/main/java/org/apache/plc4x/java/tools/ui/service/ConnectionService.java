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

package org.apache.plc4x.java.tools.ui.service;

import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.tools.ui.model.Device;
import org.apache.plc4x.java.tools.ui.repository.ConnectionRepository;
import org.apache.plc4x.java.utils.cache.CachedPlcConnectionManager;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;

@Component
public class ConnectionService {

    private final ConnectionRepository connectionRepository;
    private final ApplicationEventPublisher publisher;
    private final CachedPlcConnectionManager cachedPlcConnectionManager;

    public ConnectionService(ConnectionRepository connectionRepository, ApplicationEventPublisher publisher, PlcDriverManager driverManager) {
        this.connectionRepository = connectionRepository;
        this.publisher = publisher;
        this.cachedPlcConnectionManager = CachedPlcConnectionManager.getBuilder(driverManager.getConnectionManager()).build();
    }

    public List<Device> getAllConnections() {
        return connectionRepository.findAll();
    }

    public Device createConnection(Device device) {
        return connectionRepository.save(device);
    }

    public Optional<Device> readConnection(Integer id) {
        return connectionRepository.findById(id);
    }

    public Device updateConnection(Device device) {
        return connectionRepository.save(device);
    }

    public void deleteConnection(Device device) {
        connectionRepository.delete(device);
    }

}
