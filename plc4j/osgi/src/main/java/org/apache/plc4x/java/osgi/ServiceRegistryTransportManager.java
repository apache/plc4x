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
package org.apache.plc4x.java.osgi;

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportManager;
import org.osgi.framework.BundleContext;
import org.osgi.util.tracker.ServiceTracker;

import java.util.Optional;

/**
 * Resolves transports from the {@link Transport} services in the OSGi service registry, which the
 * {@link TransportActivator} of each transport bundle registers.
 * <p>
 * The lookup happens per call, so it doesn't matter whether a transport bundle starts before or after the drivers.
 */
@SuppressWarnings("rawtypes")
public class ServiceRegistryTransportManager implements TransportManager {

    private final ServiceTracker<Transport, Transport> tracker;

    public ServiceRegistryTransportManager(BundleContext context) {
        this.tracker = new ServiceTracker<>(context, Transport.class, null);
    }

    public void open() {
        tracker.open();
    }

    public void close() {
        tracker.close();
    }

    @Override
    public Optional<Transport> getTransport(String transportCode) {
        // getTracked() is sorted highest service ranking first, so that one wins if two bundles share a code.
        return tracker.getTracked().values().stream()
            .filter(transport -> transport.getTransportCode().equals(transportCode))
            .findFirst();
    }

}
