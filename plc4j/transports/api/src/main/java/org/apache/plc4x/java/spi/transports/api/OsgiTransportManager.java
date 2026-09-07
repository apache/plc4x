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

import org.osgi.service.component.annotations.Component;
import org.osgi.service.component.annotations.Reference;
import org.osgi.service.component.annotations.ReferenceCardinality;
import org.osgi.service.component.annotations.ReferencePolicy;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

@Component(service = TransportManager.class, immediate = true)
public class OsgiTransportManager implements TransportManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(OsgiTransportManager.class);

    private final Map<String, Transport<?>> transportMap = new ConcurrentHashMap<>();

    @Reference(
        cardinality = ReferenceCardinality.MULTIPLE,
        policy = ReferencePolicy.DYNAMIC
    )
    protected void bindTransport(Transport<?> transport) {
        String transportCode = transport.getTransportCode();
        LOGGER.info("Registering transport {} ({})", transportCode, transport.getTransportName());

        Transport<?> existing = transportMap.putIfAbsent(transportCode, transport);
        if (existing != null) {
            LOGGER.warn("Transport with code '{}' is already registered. Ignoring duplicate.", transportCode);
        }
    }

    protected void unbindTransport(Transport<?> transport) {
        String transportCode = transport.getTransportCode();
        LOGGER.info("Unregistering transport {} ({})", transportCode, transport.getTransportName());
        transportMap.remove(transportCode);
    }

    @Override
    public Optional<Transport> getTransport(String transportCode) {
        return Optional.ofNullable(transportMap.get(transportCode));
    }

}
