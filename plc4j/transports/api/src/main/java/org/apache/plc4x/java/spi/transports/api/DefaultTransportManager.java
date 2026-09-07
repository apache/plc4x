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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.ServiceLoader;

public class DefaultTransportManager implements TransportManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(DefaultTransportManager.class);

    protected final ClassLoader classLoader;

    private final Map<String, Transport> transportMap;

    public DefaultTransportManager() {
        this(Thread.currentThread().getContextClassLoader());
    }

    public DefaultTransportManager(ClassLoader classLoader) {
        LOGGER.info("Instantiating new Transport Manager with class loader {}", classLoader);
        this.classLoader = classLoader;
        transportMap = new HashMap<>();
        ServiceLoader<Transport> transportLoader = ServiceLoader.load(Transport.class, classLoader);
        LOGGER.info("Registering available transports...");
        for (Transport transport : transportLoader) {
            if (transportMap.containsKey(transport.getTransportCode())) {
                throw new IllegalStateException(
                    "Multiple transport implementations available for transport code '" +
                        transport.getTransportName() + "'");
            }
            LOGGER.info("Registering transport {} ({})", transport.getTransportCode(), transport.getTransportName());
            transportMap.put(transport.getTransportCode(), transport);
        }
    }

    @Override
    public Optional<Transport> getTransport(String transportCode) {
        return Optional.ofNullable(transportMap.get(transportCode));
    }

}
