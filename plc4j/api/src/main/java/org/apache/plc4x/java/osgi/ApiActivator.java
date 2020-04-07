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

import org.apache.plc4x.java.DriverManager;
import org.apache.plc4x.java.api.PlcDriver;
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

public class ApiActivator implements BundleActivator, BundleTrackerCustomizer<List<ServiceRegistration<PlcDriver>>>,
    ServiceFactory<DriverManager> {

    private Logger logger = LoggerFactory.getLogger(ApiActivator.class);
    private ServiceRegistration<DriverManager> registration;
    private BundleTracker<List<ServiceRegistration<PlcDriver>>> tracker;

    @Override
    public void start(BundleContext context) throws Exception {
        tracker = new BundleTracker<>(context, Bundle.ACTIVE, this);
        tracker.open();

        registration = context.registerService(DriverManager.class, new OsgiDriverManager(context), new Hashtable<>());
    }

    @Override
    public void stop(BundleContext context) throws Exception {
        registration.unregister();

        tracker.close();
    }

    @Override
    public List<ServiceRegistration<PlcDriver>> addingBundle(Bundle bundle, BundleEvent event) {
        if (bundle.getBundleId() == 0) {
            return null;
        }
        try {
            ClassLoader cl = bundle.adapt(BundleWiring.class).getClassLoader();
            ServiceLoader<PlcDriver> drivers = ServiceLoader.load(PlcDriver.class, cl);
            List<ServiceRegistration<PlcDriver>> registrations = new ArrayList<>();
            for (PlcDriver driver : drivers) {
                Hashtable<String, String> props = new Hashtable<String, String>();
                props.put(OsgiDriverManager.PROTOCOL_CODE, driver.getProtocolCode());
                props.put(OsgiDriverManager.PROTOCOL_NAME, driver.getProtocolName());
                ServiceRegistration<PlcDriver> reg = bundle.getBundleContext().registerService(PlcDriver.class, driver, props);
                registrations.add(reg);
                logger.info("Added {}",driver.getProtocolName());
            }

            return registrations.isEmpty() ? null : registrations;
        } catch (Exception ex) {
            logger.error("Could not register PlcDrivers from bundle {}", bundle.getSymbolicName(), ex);
            return null;
        }
    }

    @Override
    public void modifiedBundle(Bundle bundle, BundleEvent event, List<ServiceRegistration<PlcDriver>> object) {
    }

    @Override
    public void removedBundle(Bundle bundle, BundleEvent event, List<ServiceRegistration<PlcDriver>> object) {
        for (ServiceRegistration<PlcDriver> reg : object) {
            reg.unregister();
        }
    }

    @Override
    public DriverManager getService(Bundle bundle, ServiceRegistration<DriverManager> registration) {
        return new OsgiDriverManager(bundle.getBundleContext());
    }

    @Override
    public void ungetService(Bundle bundle, ServiceRegistration<DriverManager> registration, DriverManager service) {
        if (service instanceof OsgiDriverManager) {
            ((OsgiDriverManager) service).close();
        }
    }
}