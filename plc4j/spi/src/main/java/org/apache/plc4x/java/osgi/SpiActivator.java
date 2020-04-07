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

import org.apache.plc4x.java.spi.transport.Transport;
import org.osgi.framework.*;
import org.osgi.framework.wiring.BundleWiring;
import org.osgi.util.tracker.BundleTracker;
import org.osgi.util.tracker.BundleTrackerCustomizer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Hashtable;
import java.util.List;
import java.util.ServiceLoader;

public class SpiActivator implements BundleActivator, BundleTrackerCustomizer<List<ServiceRegistration<Transport>>>{

    private Logger logger = LoggerFactory.getLogger(SpiActivator.class);
    private BundleTracker<List<ServiceRegistration<Transport>>> tracker;

    private final String TRANSPORT_CODE ="org.apache.plc4x.transport.code";
    private final String TRANSPORT_NAME ="org.apache.plc4x.transport.name";

    @Override
    public void start(BundleContext context) throws Exception {
        tracker = new BundleTracker<>(context, Bundle.ACTIVE, this);
        tracker.open();
    }

    @Override
    public void stop(BundleContext context) throws Exception {
        tracker.close();
    }

    @Override
    public List<ServiceRegistration<Transport>> addingBundle(Bundle bundle, BundleEvent event) {
        if (bundle.getBundleId() == 0) {
            return null;
        }
        try {
            ClassLoader cl = bundle.adapt(BundleWiring.class).getClassLoader();
            ServiceLoader<Transport> transports = ServiceLoader.load(Transport.class, cl);
            List<ServiceRegistration<Transport>> registrations = new ArrayList<>();
            for (Transport transport : transports) {
                Hashtable<String, String> props = new Hashtable<String, String>();
                props.put(TRANSPORT_CODE, transport.getTransportCode());
                props.put(TRANSPORT_NAME, transport.getTransportName());
                ServiceRegistration<Transport> reg = bundle.getBundleContext().registerService(Transport.class, transport, props);
                registrations.add(reg);
                logger.info("Added {}",transport.getTransportName());
            }

            return registrations.isEmpty() ? null : registrations;
        } catch (Exception ex) {
            logger.error("Could not register Transport(s) from bundle {}", bundle.getSymbolicName(), ex);
            return null;
        }
    }

    @Override
    public void modifiedBundle(Bundle bundle, BundleEvent event, List<ServiceRegistration<Transport>> object) {
    }

    @Override
    public void removedBundle(Bundle bundle, BundleEvent event, List<ServiceRegistration<Transport>> object) {
        for (ServiceRegistration<Transport> reg : object) {
            reg.unregister();
        }
    }

}